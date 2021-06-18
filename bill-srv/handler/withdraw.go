package handler

import (
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/bill-srv/i"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	go_micro_srv_settle "a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

// 手机号取款,只有未激活的账号才可以取款,账单中产生的手续费是内扣的形式.
func (b *BillHandler) Withdrawal(ctx context.Context, req *go_micro_srv_bill.WithdrawalRequest, reply *go_micro_srv_bill.WithdrawalReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//判断登录的是什么账号类型
	servicerNo := ""
	var opAccType int // 1-服务商;2-收银员
	switch req.AccountType {
	case constants.AccountType_POS:
		servicerNo = dao.ServiceDaoInst.GetServicerNoByCashierNo(req.OpAccNo)
		opAccType = constants.OpAccType_Pos
	case constants.AccountType_SERVICER:
		servicerNo = req.OpAccNo
		opAccType = constants.OpAccType_Servicer
	default:
		ss_log.Error("AccountType异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 判断是否有取款权限
	_, outGoPermission := dao.ServiceDaoInst.GetPermissionFromSrvNo(servicerNo)
	if outGoPermission == "" || outGoPermission == "0" {
		ss_log.Error("err=[手机号取款接口,没有取款权限,当前服务商id为----->%s,当前取款权限为----->%s]", servicerNo, outGoPermission)
		reply.ResultCode = ss_err.ERR_NOT_ROLE
		return nil
	}

	// 如果收款人是平台的有效账号,就不允许进来手机号取款,应该是扫一扫取款
	isActive := dao.AccDaoInstance.GetIsActiveFromPhone(req.RecvPhone, req.RecvCountryCode)
	if isActive == "" {
		ss_log.Error("err=[---->%s]", "手机号取款, 查询收款人账号的账号是否激活错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if isActive == constants.AccountActived {
		// 已经有账号,不能用核销码取款,把核销码状态改为禁用状态
		//if req.SaveCode != "" {
		//	ss_log.Error("------------%s", "已有虚拟账户,不能使用核销码")
		//	reply.ResultCode = ss_err.ERR_CANT_USE_WRITE_OFF_CODE
		//	return nil
		//}
		reply.ResultCode = ss_err.ERR_REC_ACCOUNT_IS_EXISTS_NO_PHONE_WITHDRAWAL
		return nil
	}

	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[手机号取款, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}
	// 校验最大金额最小金额限制
	if err := CheckAmountIsMaxMinWithdraw(req.MoneyType, req.Amount); err != nil {
		ss_log.Error("提现最大最小金额校验失败,err:--->%s", err.Error())
		reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
		return nil
	}

	// 修正各币种的金额
	amount1 := common.NormalAmountByMoneyType(req.MoneyType, req.Amount)

	// 通过币种获取手续费类型
	feesType, fErr := common.FeesTypeByMoneyType(constants.Scene_Withdraw, req.MoneyType)
	if fErr != nil {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 通过币种获取虚拟账号类型
	vaType, vaErr := common.VirtualAccountTypeByMoneyType(req.MoneyType, "0")
	if vaErr != nil {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 判断支付密码是否被限制,并验证支付密码，如果支付密码正确且不是限制中则清除连续支付密码出错次数
	replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &go_micro_srv_auth.CheckPayPWDRequest{
		AccountUid:  req.AccountUid,
		AccountType: req.AccountType,
		Password:    req.Password,
		NonStr:      req.NonStr,
		IdenNo:      req.OpAccNo,
	})
	if errCheckPayPwd != nil {
		ss_log.Error("paymentPwdErrLimit 调用验证支付密码接口出错,err[%v]", errCheckPayPwd)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if replyCheckPayPwd.ResultCode != ss_err.ERR_SUCCESS {
		reply.ResultCode = replyCheckPayPwd.ResultCode
		reply.PayPasswordErrTips = replyCheckPayPwd.ErrTips //提示还可以输入几次错误支付密码
		return nil
	}

	recvAccNo := dao.AccDaoInstance.GetAccNoFromPhone(req.RecvPhone, req.RecvCountryCode)
	//sendAccNo := dao.AccDaoInstance.GetAccNoFromPhone(req.SendPhone, req.SaveCountryCode)
	if recvAccNo == "" {
		ss_log.Error("[%v]", "recvAccNo为空")
		return nil
	}

	//  获取手续费
	rate, fees, feeErr := doFees(feesType, req.Amount)
	if feeErr != nil {
		ss_log.Error("提现计算手续费失败,err:--->%s", feeErr.Error())
		reply.ResultCode = ss_err.ERR_RATE_FAILD
		return nil
	}
	ss_log.Info("req.amount为----->%s,fees为----->%s", req.Amount, fees)
	//=========================================

	recvVaccNo := InternalCallHandlerInst.ConfirmExistVAccount(recvAccNo, req.MoneyType, strext.ToInt32(vaType))

	var withdrawalAmount string
	if fees != "" && fees != "0" {
		withdrawalAmount = ss_count.Sub(req.Amount, fees).String()
	} else {
		withdrawalAmount = req.Amount
	}

	// 判断上一个操作员存进的单跟现在的时间相比是否超过5秒
	fiveSec := ss_time.Now(global.Tz).Add(-5 * time.Second).Format(ss_time.DateTimeDashFormat)
	createTime, fErr := dao.OutgoOrderDaoInst.QueryCreateTime(req.OpAccNo, fiveSec)
	if fErr != nil && fErr.Error() != ss_sql.DB_NO_ROWS_MSG {
		ss_log.Error("IncomeOrder 查询最新时间失败,err: %s", fErr.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	if createTime != "" {
		ss_log.Info("最新时间为: %s,5秒前的时间 %s", createTime, fiveSec)
		reply.ResultCode = ss_err.ERR_ORDER_SUBMITTED_FREQUENTLY
		return nil
	}

	logNo := dao.OutgoOrderDaoInst.InsertOutgoV2(recvVaccNo, req.Amount, servicerNo, req.OpAccNo, req.MoneyType,
		fees, rate, withdrawalAmount, req.Lat, req.Lng, req.Ip, constants.WITHDRAW_PHONE, strext.ToInt32(opAccType))
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_PAY_OUT_MONEY
		return nil
	}

	////风控
	//riskReply, _ := i.RiskCtrlHandleInstance.Client.GetRiskCtrlReuslt(context.TODO(), &go_micro_srv_riskctrl.GetRiskCtrlResultRequest{
	//	ApiType: "mobile_num_withdrawal",
	//	// 发起支付的账号
	//	PayerAccNo: sendAccNo,
	//	ActionTime: time.Now().String(),
	//	Amount:     req.Amount,
	//	Ip:         req.Ip,
	//	PayType:    constants.Risk_Pay_Type_Mobile_Num_Withdrawal,
	//	// 收款人账号
	//	PayeeAccNo:  recvAccNo,
	//	ProductType: "mobile_num_withdrawal",
	//	// 币种
	//	MoneyType: req.MoneyType,
	//	// 订单号
	//	OrderNo: logNo,
	//})
	//
	//ss_log.Info("手机号取款 风控返回结果,操作结果是---->%s,RiskNo为----->%s", riskReply.OpResult, riskReply.RiskNo)
	//
	//if riskReply.OpResult == constants.Risk_Result_No_Pass_Str {
	//	reply.ResultCode = ss_err.ERR_RISK_IS_RISK
	//	reply.RiskNo = riskReply.RiskNo
	//	return nil
	//}

	// 根据收款人的账号判断是否是注销账户,如果是的话就需要判断校验码是否正确,判断是否有一次性提取金额,
	var amount string
	writeInfo, err := dao.WriteoffInst.QueryOrderNo(req.SaveCode, req.RecvPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("核销码有误，查询结果为空; saveCode为--->%s,收款人手机号--->%s", req.SaveCode, req.RecvPhone)
			reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_FAILD
			return nil
		}
		ss_log.Error("根据saveCode查询订单号失败,saveCode为--->%s,收款人手机号--->%s", req.SaveCode, req.RecvPhone)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if writeInfo.DurationTime == "" {
		ss_log.Error("查询核销码期限时间为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	nowTimeStr := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	endTimeStr := ss_time.StripPostTime(writeInfo.DurationTime)
	//核销码时间必需在有效期内，过期不允许提现
	if cmp, _ := ss_time.CompareDate(ss_time.DateTimeDashFormat, nowTimeStr, endTimeStr); cmp > 0 {
		ss_log.Error("核销码[%v]已过期。", req.SaveCode)
		if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderStatus(logNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_Expired
		return nil
	}

	if writeInfo.IncomeOrderNo != "" {
		//这一步查询好像没意义，188行已经使用相同条件查过一次
		//// 判断核销码是否正确,根据状态为1,取出income_order_id
		//incomeOrderNo := dao.WriteoffInst.QueryIncomeOrderNo(req.SaveCode, req.RecvPhone)
		//if incomeOrderNo == "" { // 核销码不正确
		//	reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_FAILD
		//	return nil
		//}
		amount = dao.IncomeOrderDaoInst.QueryAmount(writeInfo.IncomeOrderNo)
	}

	if writeInfo.TransferOrderNo != "" {
		//这一步查询好像没意义，188行已经使用相同条件查过一次
		//transferOrderNo := dao.WriteoffInst.QueryTransferOrderNo(req.SaveCode, req.RecvPhone)
		//if transferOrderNo == "" { // 核销码不正确
		//	reply.ResultCode = ss_err.ERR_WRITE_OFF_CODE_FAILD
		//	return nil
		//}
		amount = dao.TransferDaoInst.QueryAmount(writeInfo.TransferOrderNo)
	}

	// 判断金额是否一次性取完
	if req.Amount != amount {
		ss_log.Error("虚拟账户里的余额为----->%s,用户提取的金额为--->%s", amount, req.Amount)
		reply.ResultCode = ss_err.ERR_WRONG_AMOUNT
		return nil
	}
	// 修改码状态为已使用
	if errCode := dao.WriteoffInst.UpdateWriteoffStatus(tx, req.SaveCode, constants.WriteOffCodeIsUse, logNo); errCode != ss_err.ERR_SUCCESS {
		reply.ResultCode = errCode
		return nil
	}

	// 修改虚账
	errCode := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, recvVaccNo, withdrawalAmount, "-", logNo, constants.VaReason_OUTGO)
	if errCode != ss_err.ERR_SUCCESS {
		if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderStatus(logNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		reply.ResultCode = errCode
		return nil
	}
	// 扣除手续费
	if fees != "" && fees != "0" {
		err := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, recvVaccNo, fees, "-", logNo, constants.VaReason_FEES)
		if err != ss_err.ERR_SUCCESS {
			if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderStatus(logNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
			reply.ResultCode = err
			return nil
		}
	}

	if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderStatusTx(tx, logNo, constants.OrderStatus_Paid); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}
	serAccNo := dao.ServiceDaoInst.GetAccNoFromSrvNo(servicerNo)
	// 调用ps 提现接口
	quotaReq := &go_micro_srv_quota.ModifyQuotaRequest{
		CurrencyType: req.MoneyType,
		//Amount:       req.Amount,
		Amount:    withdrawalAmount, // 排除手续费的金额
		AccountNo: serAccNo,
		OpType:    constants.QuotaOp_Withdraw,
		LogNo:     logNo,
	}
	quotaRepl := &go_micro_srv_quota.ModifyQuotaReply{}
	quotaRepl, _ = i.QuotaHandleInstance.Client.ModifyQuota(ctx, quotaReq)

	if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[--------------->%s]", "客户扫码提现,调用八神的服务失败,操作为客户扫码提现")
		reply.ResultCode = quotaRepl.ResultCode
		return nil
	}

	// todo 插入 billing_details_results
	if errStr := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, req.Amount, req.MoneyType, serAccNo,
		req.AccountType, logNo, "0", constants.OrderStatus_Paid, servicerNo, req.OpAccNo, constants.BillDetailTypeOut,
		fees, withdrawalAmount); errStr == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, recvAccNo); errStr != ss_err.ERR_SUCCESS {
	//	reply.ResultCode = errStr
	//	return nil
	//}

	//appLang, _ := dao.AccDaoInstance.QueryAccountLang(recvAccNo)
	//if appLang == "" {
	//	req.Lang = constants.Lang_En
	//} else {
	//	req.Lang = appLang
	//}

	ss_log.Info("用户 %s 当前的语言为--->%s", recvAccNo, req.Lang)

	//添加取款账号推送消息
	errAddMessages2 := dao.LogAppMessagesDaoInst.AddLogAppMessagesTx(tx, logNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_OUTGO_Apply, constants.VaReason_OUTGO, recvAccNo, constants.OrderStatus_Paid)
	if errAddMessages2 != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddMessages2=[%v]", errAddMessages2)
	}

	toAccountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(recvAccNo)
	//langText := dao.LangDaoInstance.GetLangTextByKey(dbHandler, "提现成功推送消息模板", req.Lang)
	//title := dao.LangDaoInstance.GetLangTextByKey(dbHandler, "提现成功", req.Lang)
	moneyType := dao.LangDaoInstance.GetLangTextByKey(req.MoneyType, req.Lang)

	//content := &go_micro_srv_push.ContentWithArgs{
	//	//Key:  "你的账户于%s 支出%s%s,请及时查看！",
	//	Key:  langText,
	//	Args: []string{time.Now().Format("2006-01-02 15:04:05"), amount1, moneyType},
	//}
	//
	//// 消息推送
	//ev := &go_micro_srv_push.SendPushMsgRequest{
	//	Phone:    toAccountType + "_" + req.RecvPhone,
	//	Lang:     req.Lang,
	//	Content:  content,
	//	Title:    title,
	//	SendType: constants.MqSendType_Jpush,
	//}

	timeString := time.Now().Format("2006-01-02 15:04:05")
	args := []string{
		timeString, amount1, moneyType,
	}
	lang, _ := dao.AccDaoInstance.QueryAccountLang(recvAccNo)
	if lang == "" || lang == constants.LangEnUS {
		args = []string{
			amount1, moneyType, timeString,
		}
	}
	// 消息推送
	ev := &go_micro_srv_push.PushReqest{
		Accounts: []*go_micro_srv_push.PushAccout{
			{
				AccountNo:   recvAccNo,
				AccountType: toAccountType,
			},
		},
		TempNo: constants.Template_WithdrawSuccess,
		Args:   args,
	}

	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}

	// 调用 ps 客户取的操作
	ss_sql.Commit(tx)

	if fees != "0" {
		// 发送手续费进MQ
		feeEv := &go_micro_srv_settle.SettleTransferRequest{
			BillNo:    logNo,
			FeesType:  constants.FEES_TYPE_WITHDRAW,
			Fees:      fees,
			MoneyType: req.MoneyType,
		}
		ss_log.Info("publishing %+v\n", feeEv)
		// publish an event
		if err := common.SettleEvent.Publish(context.TODO(), feeEv); err != nil {
			ss_log.Error("err=[pos 手机号取款接口,手续费推送到MQ失败,err----->%s]", err.Error())
		}
	}

	reply.OrderNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 扫一扫取款码(1.pos端展示二维码)
func (b *BillHandler) GenWithdrawCode(ctx context.Context, req *go_micro_srv_bill.GenWithdrawCodeRequest, reply *go_micro_srv_bill.GenWithdrawCodeReply) error {
	// 判断码是否过期
	//if !dao.GenCodeDaoInst.CheckCodeTimeExp(req.AccountNo, "", constants.CODETYPE_SWEEP, constants.CODE_USE_STATUS_IS_NO_SWEEP) {
	//	reply.ResultCode = ss_err.ERR_QR_CODE_NO_GEN
	//	return nil
	//}

	// 获取服务商的accNo
	var servicerNo string
	var opAccType int
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		servicerNo = req.OpAccNo
		opAccType = constants.OpAccType_Servicer
	case constants.AccountType_POS: // 收银员
		servicerNo = dao.CashierDaoInst.GetServicerNoFromOpAccNo(req.OpAccNo)
		opAccType = constants.OpAccType_Pos
	default:
		ss_log.Error("生成取款码,AccountType异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 判断是否有取款权限
	_, outGoPermission := dao.ServiceDaoInst.GetPermissionFromSrvNo(servicerNo)
	if outGoPermission == "" || outGoPermission == "0" {
		ss_log.Error("err=[pos端扫一扫取款码接口,没有取款权限,当前服务商id为----->%s,当前取款权限为----->%s]", servicerNo, outGoPermission)
		reply.ResultCode = ss_err.ERR_NOT_ROLE
		return nil
	}

	code := dao.GenCodeDaoInst.PosGenCode(servicerNo, "0", "", constants.CODETYPE_SWEEP, req.OpAccNo, opAccType)
	if code == "" {
		ss_log.Error("err=[%v]", code)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Code = "S." + code
	return nil
}

// 修改手机扫一扫取款的二维码状态
func (b *BillHandler) ModifyGenCodeStatus(ctx context.Context, req *go_micro_srv_bill.ModifyGenCodeStatusRequest, reply *go_micro_srv_bill.ModifyGenCodeStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var genCode string
	if strings.HasPrefix(req.GenKey, "S.") {
		split := strings.Split(req.GenKey, "S.")
		genCode = split[1]
	} else {
		genCode = req.GenKey
	}

	// 获取当前码的状态
	useStatus := dao.GenCodeDaoInst.GetCodeStatus(tx, genCode, constants.CODETYPE_SWEEP)
	switch useStatus {
	case constants.CODE_USE_STATUS_IS_NO_SWEEP:
		if req.Status != strext.ToInt32(constants.CODE_USE_STATUS_IS_SWEEP) {
			ss_log.Info("当前码的状态为----->%s,需要修改的状态为----->%v", useStatus, req.Status)
			reply.ResultCode = ss_err.ERR_PAY_CODE_STATUS
			return nil
		}
	case constants.CODE_Pendding_Confirm:
		if req.Status != strext.ToInt32(constants.CODE_USE_STATUS_IS_PAY) {
			ss_log.Info("当前码的状态为----->%s,需要修改的状态为----->%v", useStatus, req.Status)
			reply.ResultCode = ss_err.ERR_PAY_CODE_STATUS
			return nil
		}
	default:
		//ss_log.Info("当前码的状态可能有误,当前状态为----->%s,需要修改的状态为----->%v", useStatus, req.Status)
		//reply.ResultCode = ss_err.ERR_PARAM
		//return nil
	}

	// 变更成已扫的时候才需要做时间校验
	//if strext.ToStringNoPoint(req.Status) == constants.CODE_USE_STATUS_IS_SWEEP {
	// 判断码是否过期
	if errStr := dao.GenCodeDaoInst.CheckCodeTimeExp(tx, "", genCode, constants.CODETYPE_SWEEP, constants.CODE_USE_STATUS_IS_NO_SWEEP); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_QR_CODE_EXPIRED
		return nil
	}
	//}
	// 修改状态
	if errStr := dao.GenCodeDaoInst.UpdateGenCodeUseStatus(tx, req.Status, genCode, req.AccountUid, req.AccountType, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 获取扫一扫码的状态(2.pos端查询二维码状态是否有更改)
func (b *BillHandler) QuerySweepCodeStatus(ctx context.Context, req *go_micro_srv_bill.QuerySweepCodeStatusRequest, reply *go_micro_srv_bill.QuerySweepCodeStatusReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var genCode string
	if strings.HasPrefix(req.GenCode, "S.") {
		split := strings.Split(req.GenCode, "S.")
		genCode = split[1]
	} else {
		genCode = req.GenCode
	}

	// 判断码是否过期
	errStr := dao.GenCodeDaoInst.CheckCodeTimeExp(tx, "", genCode, constants.CODETYPE_SWEEP, constants.CODE_USE_STATUS_IS_NO_SWEEP)
	if errStr != ss_err.ERR_SUCCESS {
		// 修改码为已过期状态
		if errStr := dao.GenCodeDaoInst.UpdateGenCodeExp(tx, constants.CODE_EXP, genCode); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
	}

	var nickname, headURL string
	status, orderNo, sweepAccountNo := dao.GenCodeDaoInst.QeurySweepCodeStatusFromAccountNo(tx, req.IdenNo, genCode)
	if sweepAccountNo != "" { // 获取扫码人的头像和昵称
		nickname, headURL = dao.AccDaoInstance.GetNameAndURLFromUID(tx, sweepAccountNo)
		// 拼接图片url
		value := dao.GlobalParamDaoInstance.QeuryParamValue("image_base_url")
		headURL = fmt.Sprintf("%s/%s", value, headURL)
	}
	if status == "" {
		reply.ResultCode = ss_err.ERR_QUERY_SWEEP_CODE_STATUS_FAILD
		return nil
	}
	ss_sql.Commit(tx)

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Status = status
	reply.OrderNo = orderNo
	reply.SweepAccountNo = sweepAccountNo
	reply.NickName = nickname
	reply.HeadUrl = headURL
	return nil
}

// 扫一扫取款 (3.app扫码并端输入提现金额)
func (b *BillHandler) SweepWithdrawal(ctx context.Context, req *go_micro_srv_bill.SweepWithdrawRequest, reply *go_micro_srv_bill.SweepWithdrawReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[扫一扫取款, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	if req.MoneyType == constants.CURRENCY_USD {
		// 判断金额是否包含小数点
		if strings.Contains(req.Amount, ".") {
			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_NO_INTEGER
			return nil
		}
	}

	// 普通提现
	if req.SwithdrawType == 1 {
		// 判断币种,判断最大最小金额
		switch req.MoneyType {
		case constants.CURRENCY_USD:
			// 最大限额
			usdFaceSingleMax := dao.GlobalParamDaoInstance.QeuryParamValue("usd_face_single_max")
			// 最小限额
			usdFaceSingleMin := dao.GlobalParamDaoInstance.QeuryParamValue("usd_face_single_min")
			if strext.ToFloat64(req.Amount) < strext.ToFloat64(usdFaceSingleMin) || strext.ToFloat64(req.Amount) > strext.ToFloat64(usdFaceSingleMax) {
				ss_log.Error("扫码取款,币种 美元,超出金额限制,当前金额为--->%s,最大金额限制为--->%s,最小金额限制为--->%s", req.Amount, usdFaceSingleMax, usdFaceSingleMin)
				reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
				return nil
			}

		case constants.CURRENCY_KHR:
			// 最大限额
			khrFaceSingleMax := dao.GlobalParamDaoInstance.QeuryParamValue("khr_face_single_max")
			// 最小限额
			khrFaceSingleMin := dao.GlobalParamDaoInstance.QeuryParamValue("khr_face_single_min")
			if strext.ToFloat64(req.Amount) < strext.ToFloat64(khrFaceSingleMin) || strext.ToFloat64(req.Amount) > strext.ToFloat64(khrFaceSingleMax) {
				ss_log.Error("扫码取款,币种 瑞尔,超出金额限制,当前金额为--->%s,最大金额限制为--->%s,最小金额限制为--->%s", req.Amount, khrFaceSingleMax, khrFaceSingleMin)
				reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
				return nil
			}
		}
	}

	// 判断支付密码是否被限制,并验证支付密码，如果支付密码正确且不是限制中则清除连续支付密码出错次数
	replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &go_micro_srv_auth.CheckPayPWDRequest{
		AccountUid:  req.AccountUid,
		AccountType: req.AccountType,
		Password:    req.Password,
		NonStr:      req.NonStr,
		IdenNo:      req.OpAccNo,
	})
	if errCheckPayPwd != nil {
		ss_log.Error("paymentPwdErrLimit 调用验证支付密码接口出错,err[%v]", errCheckPayPwd)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if replyCheckPayPwd.ResultCode != ss_err.ERR_SUCCESS {
		reply.ResultCode = replyCheckPayPwd.ResultCode
		reply.PayPasswordErrTips = replyCheckPayPwd.ErrTips //提示还可以输入几次错误支付密码
		return nil
	}

	vaType := 0
	var feesType int32
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
		feesType = 6
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
		feesType = 11
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	var genCode string
	if strings.HasPrefix(req.GenCode, "S.") {
		split := strings.Split(req.GenCode, "S.")
		genCode = split[1]
	} else {
		genCode = req.GenCode
	}

	// 判断码是否过期
	if errStr := dao.GenCodeDaoInst.CheckCodeTimeExp(tx, "", genCode, constants.CODETYPE_SWEEP, constants.CODE_USE_STATUS_IS_SWEEP); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = ss_err.ERR_QR_CODE_EXPIRED
		return nil
	}

	// 确保虚拟账号存在
	recvVaccNo := InternalCallHandlerInst.ConfirmExistVAccount(req.AccountUid, req.MoneyType, strext.ToInt32(vaType))
	balance, _ := dao.VaccountDaoInst.GetBalance(recvVaccNo)
	var amount string
	switch req.SwithdrawType {
	case 1: //  普通提现
		amount = req.Amount
	case 2: // 全部提现
		amount = balance
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//  获取手续费
	rate, fees, feeErr := doFees(feesType, amount)
	if feeErr != nil {
		ss_log.Error("扫码取款 计算手续费失败,err:--->%s", feeErr.Error())
		reply.ResultCode = ss_err.ERR_RATE_FAILD
		return nil
	}
	ss_log.Info("币种为------->%s,金额为----->%s,手续费为-------->%s", req.MoneyType, req.Amount, fees)

	var withdrawAmount string

	switch req.SwithdrawType {
	case 1: //  普通提现(手续费外扣)
		withdrawAmount = req.Amount
	case 2: // 全部提现(手续费内扣)
		withdrawAmountDeci := ss_count.Sub(balance, fees)
		f, _ := withdrawAmountDeci.Float64()
		if f <= 0 {
			ss_log.Error("err=[全部提现手续费不够扣,当前余额为----->%s,应扣手续费为---->%s]", balance, fees)
			reply.ResultCode = ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
			return nil
		}
		withdrawAmount = withdrawAmountDeci.String()
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	ss_log.Info("提现金额为----->%s,手续费为----->%s", withdrawAmount, fees)
	// 获取 serviceNo
	srvNo, opAccNo, opAccType := dao.GenCodeDaoInst.GetSrvFromCode(tx, genCode)
	if srvNo == "" || opAccNo == "" || opAccType == "" {
		ss_log.Error("err=[查询服务商的id,操作员id,谁操作的信息失败,二维码为------>%s]", genCode)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	logNo := dao.OutgoOrderDaoInst.InsertOutgo(tx, recvVaccNo, withdrawAmount, srvNo, opAccNo, req.MoneyType, fees, rate, withdrawAmount, req.Lat, req.Lng, req.Ip, req.SwithdrawType, strext.ToInt32(opAccType))
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_PAY_OUT_MONEY
		return nil
	}

	//风控
	riskReply, _ := i.RiskCtrlHandleInstance.Client.GetRiskCtrlReuslt(context.TODO(), &go_micro_srv_riskctrl.GetRiskCtrlResultRequest{
		ApiType: constants.Risk_Ctrl_Sweep_Withdrawal,
		// 发起支付的账号
		PayerAccNo: req.AccountUid,
		ActionTime: time.Now().String(),
		Amount:     req.Amount,
		Ip:         req.Ip,
		PayType:    constants.Risk_Pay_Type_Sweep_Withdrawal,
		// 收款人账号
		PayeeAccNo:  req.AccountUid,
		ProductType: constants.Risk_Ctrl_Sweep_Withdrawal,
		// 币种
		MoneyType: req.MoneyType,
		// 订单号
		OrderNo: logNo,
	})

	ss_log.Info("扫一扫取款 风控返回结果,操作结果是---->%s,RiskNo为----->%s", riskReply.OpResult, riskReply.RiskNo)

	if riskReply.OpResult == constants.Risk_Result_No_Pass_Str {
		reply.ResultCode = ss_err.ERR_RISK_IS_RISK
		return nil
	}

	// 修改订单的 risk_no
	if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderRiskNo(tx, logNo, riskReply.RiskNo); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[用户扫码取款,新增risk_no失败,订单号为----->%s, risk_no为----->%s]", logNo, riskReply.RiskNo)
		reply.ResultCode = errStr
		return nil
	}

	// 判断输入的金额是否超额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", recvVaccNo, withdrawAmount, logNo, constants.VaReason_OUTGO, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	if fees != "0" {
		// 修改手续费
		//if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, recvVaccNo, fees, "-", logNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		//	reply.ResultCode = errStr
		//	return nil
		//}

		if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", recvVaccNo, fees, logNo, constants.VaReason_FEES, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

	}

	ss_log.Info("gencode为------------->%s,原来的gencode为----->%s", genCode, req.GenCode)
	// 修改二维码的订单号
	if errStr := dao.GenCodeDaoInst.UpdateOrderNoFormGenCode(tx, logNo, genCode); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderStatusTx(tx, logNo, constants.OrderStatus_Pending_Confirm); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	// 修改状态
	if errStr := dao.GenCodeDaoInst.UpdateGenCodeStatus(tx, strext.ToInt32(constants.CODE_Pendding_Confirm), genCode, req.AccountUid); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	// 设置key进redis
	//if _, err := cache.RedisCli.SetWithExp(constants.DefPoolName, GetExpEey(genCode), genCode, "300"); err != nil {
	//	ss_log.Error("err=[设置监听的key进redis失败,genCode为-----> %s,err------> %s]", genCode, err.Error())
	//}

	// 设置key进redis
	if err := cache.RedisClient.Set(common.GetExpEey(genCode), genCode, common.SweepWithdrawalExpireTime).Err(); err != nil {
		ss_log.Error("err=[设置监听的key进redis失败,genCode为-----> %s,err------> %s]", genCode, err.Error())
	}

	//// todo 插入 billing_details_results
	//servicerNo := dao.RelaAccIdenDaoInst.GetIdenNo(srvAccNo, constants.AccountType_SERVICER)
	srvAccNo := dao.ServiceDaoInst.GetAccNoFromSrvNo(srvNo)
	// 从二维码中获取 opAccNo,opAccType
	opAccTypeT, opAccNoT := dao.GenCodeDaoInst.GetGenRecvCodeInfo(genCode, constants.CODETYPE_SWEEP)
	var accountType string
	if opAccTypeT == strext.ToStringNoPoint(constants.OpAccType_Servicer) { // 服务商
		accountType = constants.AccountType_SERVICER
	} else if opAccTypeT == strext.ToStringNoPoint(constants.OpAccType_Pos) { // 收银员
		accountType = constants.AccountType_POS
	}

	if errStr := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, req.Amount, req.MoneyType, srvAccNo, accountType,
		logNo, "0", constants.OrderStatus_Pending_Confirm, srvNo, opAccNoT, constants.BillDetailTypeOut, fees, withdrawAmount); errStr == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_sql.Commit(tx)

	reply.OrderNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// pos端确认提现操作(4.pos端确认取款)
func (b *BillHandler) ConfirmWithdrawal(ctx context.Context, req *go_micro_srv_bill.ConfirmWithdrawRequest, reply *go_micro_srv_bill.ConfirmWithdrawReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)
	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[pos端确认提现操作, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	// 判断金额是否包含小数点
	if req.MoneyType == constants.CURRENCY_USD {
		if strings.Contains(req.Amount, ".") {
			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_NO_INTEGER
			return nil
		}
	}

	// 判断支付密码是否被限制,并验证支付密码，如果支付密码正确且不是限制中则清除连续支付密码出错次数
	replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &go_micro_srv_auth.CheckPayPWDRequest{
		AccountUid:  req.AccountUid,
		AccountType: req.AccountType,
		Password:    req.Password,
		NonStr:      req.NonStr,
		IdenNo:      req.OpAccNo,
	})
	if errCheckPayPwd != nil {
		ss_log.Error("paymentPwdErrLimit 调用验证支付密码接口出错,err[%v]", errCheckPayPwd)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if replyCheckPayPwd.ResultCode != ss_err.ERR_SUCCESS {
		reply.ResultCode = replyCheckPayPwd.ResultCode
		reply.PayPasswordErrTips = replyCheckPayPwd.ErrTips //提示还可以输入几次错误支付密码
		return nil
	}

	// 根据订单号查找 risk_no
	riskNo, err := dao.OutgoOrderDaoInst.GetRiskNoFromLogNo(req.OutOrderNo)
	if err != nil {
		ss_log.Error("err=[根据订单号查询risk_no 失败,err为----->%s]", err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_log.Info("pos 端确认取款,根据订单号查询出来的riskNo为----->%s", riskNo)
	// todo 执行风控查单
	result2Reply, _ := i.RiskCtrlHandleInstance.Client.GetRiskCtrlReuslt2(ctx, &go_micro_srv_riskctrl.GetRiskCtrlResult2Request{
		RiskNo: riskNo,
	})
	if result2Reply.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[pos端确认取款,客户信息被风控了,riskNo为----->%s]", riskNo)

		// 调用 pos 端取消接口
		cancelReq := &go_micro_srv_bill.CancelWithdrawRequest{
			OrderNo:      req.OutOrderNo,
			CancelReason: "风控不通过",
		}
		cancelRepl := &go_micro_srv_bill.CancelWithdrawReply{}
		_ = b.CancelWithdraw(ctx, cancelReq, cancelRepl)
		if cancelRepl.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[确认取款接口,调用取消确认的rpc失败,订单号为----->%s]", req.OutOrderNo)
			reply.ResultCode = cancelRepl.ResultCode
			return nil
		}

		reply.ResultCode = result2Reply.ResultCode
		return nil
	}

	vaType := 0
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	var genCode string
	if strings.HasPrefix(req.GenCode, "S.") {
		split := strings.Split(req.GenCode, "S.")
		genCode = split[1]
	} else {
		genCode = req.GenCode
	}
	ss_log.Info("%s", genCode)
	// 判断码是否过期
	if errStr := dao.GenCodeDaoInst.CheckCodeTimeExp(tx, "", genCode, constants.CODETYPE_SWEEP, constants.CODE_Pendding_Confirm); errStr != ss_err.ERR_SUCCESS {
		// 调用 pos 端取消接口
		cancelReq := &go_micro_srv_bill.CancelWithdrawRequest{
			OrderNo:      req.OutOrderNo,
			CancelReason: "二维码过期",
		}
		cancelRepl := &go_micro_srv_bill.CancelWithdrawReply{}
		_ = b.CancelWithdraw(ctx, cancelReq, cancelRepl)
		if cancelRepl.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[确认取款接口,调用取消确认的rpc失败,订单号为----->%s]", req.OutOrderNo)
			reply.ResultCode = cancelRepl.ResultCode
			return nil
		}

		// 修改订单状态为已超时,错误状态
		if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderStatusTx(tx, req.OutOrderNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		reply.ResultCode = ss_err.ERR_QR_CODE_EXPIRED
		return nil
	}

	// 根据币种,用户ID查询出虚拟账号的ID,
	vaccNo := InternalCallHandlerInst.ConfirmExistVAccount(req.UseAccountUid, req.MoneyType, strext.ToInt32(vaType))

	// 根据订单ID,pendding状态订单的金额致性(outgoorder)
	_, createTimeT := dao.OutgoOrderDaoInst.GetAmountFromLogNo(req.OutOrderNo, constants.OrderStatus_Pending_Confirm)
	//if amount != req.Amount {
	//	ss_log.Error("err=[数据库的金额为----->%s,请求的金额为----->%s]", amount, req.Amount)
	//	reply.ResultCode = ss_err.ERR_WRONG_AMOUNT
	//	return nil
	//}

	ss_log.Info("订单生成的时间为----->%s", createTimeT)
	// todo 判断订单是否超时
	if ss_time.ParseTimeFromPostgres(createTimeT, global.Tz).Add(5 * time.Minute).Before(ss_time.Now(global.Tz)) {
		reply.ResultCode = ss_err.ERR_PAY_TIMEOUT
		return nil
	}
	// 查询手续费
	_, _, _, _, fees, _, _, sweepType := dao.OutgoOrderDaoInst.QueryOutGoOrderFromLogNo(req.OutOrderNo, constants.OrderStatus_Pending_Confirm)
	if sweepType == constants.WITHDRAWAL_TYPE_ALL { // 全部提现
		amountDeci := ss_count.Sub(req.Amount, fees)
		amountF, _ := amountDeci.Float64()
		if amountF < 0 {
			ss_log.Error("全部提现,手续费超出提现金额,提现金额为--->%s,手续费为--->%s", req.Amount, fees)
			reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
			return nil
		}
		req.Amount = amountDeci.String()
	}
	// 修改冻结余额
	if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero(tx, vaccNo, req.Amount, req.OutOrderNo, constants.VaReason_OUTGO, fees); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	if fees != "" && fees != "0" {
		// 修改冻结手续费
		if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero(tx, vaccNo, fees, req.OutOrderNo, constants.VaReason_FEES, ""); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
	}

	// 更新outgo的状态为已支付状态
	if errStr := dao.OutgoOrderDaoInst.UpdateOutgoOrderStatusTx(tx, req.OutOrderNo, constants.OrderStatus_Paid); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	// 修改码为已支付状态
	if errStr := dao.GenCodeDaoInst.UpdateGenCodeExp(tx, constants.CODE_USE_STATUS_IS_PAY, req.GenCode); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	servicerNo := ""
	servicerAccNo := ""
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		servicerNo = req.OpAccNo
		servicerAccNo = dao.RelaAccIdenDaoInst.GetAccNo(servicerNo, constants.AccountType_SERVICER)
	case constants.AccountType_POS: // 收银员
		servicerNo = dao.CashierDaoInst.GetServicerNoFromOpAccNo(req.OpAccNo)
		servicerAccNo = dao.ServiceDaoInst.GetAccNoFromSrvNo(servicerNo)
	}

	// 因提交的参数都是logNo的数据，没法直接使用billing_details_results表的bill_no来查询实际金额, 待优化(需要和pos端一起)
	servicerRealAmount, err := dao.BillingDetailsResultsDaoInstance.GetRealAmountByOutOrderLogNo(req.OutOrderNo)
	if err != nil {
		ss_log.Error("通过LogNo查询实际金额失败,err:%v,LogNo:%v", err, req.OutOrderNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 调用ps 提现接口
	quotaReq := &go_micro_srv_quota.ModifyQuotaRequest{
		CurrencyType: req.MoneyType,
		//Amount:       req.Amount,
		Amount:    servicerRealAmount, // 修改为实际金额
		AccountNo: servicerAccNo,
		OpType:    constants.QuotaOp_Withdraw,
		LogNo:     req.OutOrderNo,
	}
	quotaRepl := &go_micro_srv_quota.ModifyQuotaReply{}
	quotaRepl, _ = i.QuotaHandleInstance.Client.ModifyQuota(ctx, quotaReq)

	if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[--------------->%s]", "客户扫码提现,调用八神的服务失败,操作为客户扫码提现")
		reply.ResultCode = quotaRepl.ResultCode
		return nil
	}

	// todo 插入 billing_details_results
	//if errStr := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, req.Amount, req.MoneyType, servicerAccNo, constants.AccountType_SERVICER, req.OutOrderNo, "0", constants.OrderStatus_Paid, 2); errStr == "" {
	//	reply.ResultCode = ss_err.ERR_PARAM
	//	return nil
	//}

	// todo 修改订单明细中订单状态为已支付状态
	if errStr := dao.BillingDetailsResultsDaoInstance.UpdateOrderStatusFromLogNo(tx, constants.OrderStatus_Paid, req.OutOrderNo); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	// 获取用户的语言
	appLang, _ := dao.AccDaoInstance.QueryAccountLang(req.UseAccountUid)
	if appLang == "" {
		req.Lang = constants.LangEnUS
	} else {
		req.Lang = appLang
	}

	ss_log.Info("用户 %s 当前的语言为--->%s", req.UseAccountUid, req.Lang)

	//添加pos取款推送消息
	errAddMessages := dao.LogAppMessagesDaoInst.AddLogAppMessagesTx(tx, req.OutOrderNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_OUTGO_Apply, constants.VaReason_OUTGO, req.UseAccountUid, constants.OrderStatus_Paid)
	if errAddMessages != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddMessages=[%v]", errAddMessages)
	}

	toAccountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(req.UseAccountUid)
	moneyType := dao.LangDaoInstance.GetLangTextByKey(req.MoneyType, req.Lang)
	amountB := ""
	switch req.MoneyType {
	case "usd":
		amountB = strext.ToStringNoPoint(strext.ToFloat64(req.Amount) / 100)
	case "khr":
		amountB = strext.ToStringNoPoint(strext.ToFloat64(req.Amount))
	}

	timeString := time.Now().Format("2006-01-02 15:04:05")
	args := []string{
		timeString, amountB, moneyType,
	}
	lang, _ := dao.AccDaoInstance.QueryAccountLang(req.UseAccountUid)
	if lang == "" || lang == constants.LangEnUS {
		args = []string{
			amountB, moneyType, timeString,
		}
	}
	// 消息推送
	ev := &go_micro_srv_push.PushReqest{
		Accounts: []*go_micro_srv_push.PushAccout{
			{
				AccountNo:   req.UseAccountUid,
				AccountType: toAccountType,
			},
		},
		TempNo: constants.Template_WithdrawSuccess,
		Args:   args,
	}

	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}

	// todo 因为订单在pendding状态,需要等pos端确认以后才收费
	if fees != "" && fees != "0" {
		// 发送手续费进MQ
		feeEv := &go_micro_srv_settle.SettleTransferRequest{
			BillNo:    req.OutOrderNo,
			FeesType:  constants.FEES_TYPE_WITHDRAW,
			Fees:      fees,
			MoneyType: req.MoneyType,
		}
		ss_log.Info("publishing %+v\n", feeEv)
		// publish an event
		if err := common.SettleEvent.Publish(context.TODO(), feeEv); err != nil {
			ss_log.Error("err=[pos 扫一扫取款接口,手续费推送到MQ失败,err----->%s]", err.Error())
		}
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// pos机 手机号,扫一扫取款打印小票查询
func (b *BillHandler) WithdrawReceipt(ctx context.Context, req *go_micro_srv_bill.WithdrawReceiptRequest, reply *go_micro_srv_bill.WithdrawReceiptReply) error {
	amount, serviceNo, finishTime, vaccountNo, fees, balanceType, _, withdrawType := dao.OutgoOrderDaoInst.QueryOutGoOrderFromLogNo(req.OrderNo, constants.OrderStatus_Paid)
	// 查询取款人手机号
	//phone := dao.AccDaoInstance.GetPhoneFromVAccNo(vaccountNo)
	phone, countryCode, err := dao.AccDaoInstance.GetPhoneCountryCodeFromVAccNo(vaccountNo)
	if err != nil {
		ss_log.Error("WithdrawReceipt 查询手机号和国家码失败,err: %s", err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 查询pos机号
	num := dao.ServicerTerminalDaoInstance.QueryNumberFromServiceNo(serviceNo)
	if amount == "" || num == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 判断提现类型
	var withdrawAmount, arriveAmount string
	switch withdrawType {
	case strext.ToStringNoPoint(constants.WITHDRAW_PHONE): //手机号取款
		withdrawAmount = amount
		arriveAmount = ss_count.Sub(amount, fees).String()
	case strext.ToStringNoPoint(constants.WITHDRAW_SWEEP): // 扫码提现
		withdrawAmount = amount
		arriveAmount = amount
	case strext.ToStringNoPoint(constants.WITHDRAW_SWEEP_ALL): // 全部提现
		arriveAmount = amount
		withdrawAmount = ss_count.Add(amount, fees)
	}

	data := &go_micro_srv_bill.WithdrawReceiptResult{
		// 订单号
		OrderNo: req.OrderNo,
		// 商户号
		ServiceNo: serviceNo,
		// 终端编号
		TerminalNumber: num,
		// 取款手机号
		WithdrawPhone: fmt.Sprintf("%s%s", countryCode, phone),
		// 申请金额
		ApplyAmount: withdrawAmount,
		// 手续费
		Fees: fees,
		// 到账金额
		ArriveAmount: arriveAmount,
		// 日期
		Date:      finishTime,
		MoneyType: balanceType,
	}
	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (b BillHandler) SweepWithdrawDetail(ctx context.Context, req *go_micro_srv_bill.SweepWithdrawDetailRequest, reply *go_micro_srv_bill.SweepWithdrawDetailReply) error {
	if req.OrderNo == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	var withdrawAmount, arriveAmount string
	amount, createTime, fees, moneyType, vaccountNo, withdrawType, orderStatus, _ := dao.OutgoOrderDaoInst.GetOutGoDetailFromLogNo(req.OrderNo)
	switch withdrawType {
	case strext.ToStringNoPoint(constants.WITHDRAW_PHONE): //手机号取款
		withdrawAmount = amount
		arriveAmount = ss_count.Sub(amount, fees).String()
	case strext.ToStringNoPoint(constants.WITHDRAW_SWEEP): // 扫码提现
		withdrawAmount = amount
		arriveAmount = amount
	case strext.ToStringNoPoint(constants.WITHDRAW_SWEEP_ALL): // 全部提现
		arriveAmount = amount
		withdrawAmount = ss_count.Add(amount, fees)
	}

	//获取产生该笔订单的账号
	account, errGet := dao.BillingDetailsResultsDaoInstance.GetAccountByOrderNo(req.OrderNo)
	if errGet != nil {
		ss_log.Error("获取产生该笔订单的账号失败，OrderNo= [%v],err=[%v]", req.OrderNo, errGet)
	}

	// 通过订单号查询码的信息和扫码人
	genCode, sweepAccNo := dao.GenCodeDaoInst.GetCodeAccNoFromLogNo(req.OrderNo)
	// 根据虚账查询出手机号
	phone := dao.AccDaoInstance.GetPhoneFromVAccNo(vaccountNo)
	data := &go_micro_srv_bill.SweepWithdrawDetailResult{
		OrderNo:         req.OrderNo,
		WithdrawAmount:  withdrawAmount,
		ArriveAmount:    arriveAmount,
		WithdrawPhone:   phone,
		Fees:            fees,
		Date:            createTime,
		MoneyType:       moneyType,
		OrderStatus:     orderStatus,
		GenCode:         "S." + genCode,
		SweepAccountUid: sweepAccNo,
		Account:         account,
	}
	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// pos 端取消取款
func (b *BillHandler) CancelWithdraw(ctx context.Context, req *go_micro_srv_bill.CancelWithdrawRequest, reply *go_micro_srv_bill.CancelWithdrawReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//  判断订单状态是否是待确认状态.只有在待确认状态的订单才可以取消退款操作.
	status := dao.OutgoOrderDaoInst.GetOutGoStatusFromLogNo(req.OrderNo)
	if status != constants.OrderStatus_Pending_Confirm {
		ss_log.Error("err=[订单不在待确认状态,不能取消操作,订单号为-----> %s, 状态为-----> %s]", req.OrderNo, status)
		reply.ResultCode = ss_err.ERR_ORDER_IS_NO_PENDING
		return nil
	}
	// 修改订单为失败状态
	if errStr := dao.OutgoOrderDaoInst.CancelOutgoOrder(tx, req.CancelReason, req.OrderNo); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}
	vaccNo, amount, fees := dao.OutgoOrderDaoInst.GetVaccNoTx(tx, req.OrderNo)
	if vaccNo == "" {
		ss_log.Error("err=[获取vaccNo失败,outgoOrderNo 为----->%s]", req.OrderNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 恢复余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "+", vaccNo, amount, req.OrderNo, constants.VaReason_Cancel_withdraw, constants.VaOpType_Defreeze_Add); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	if fees != "" && fees != "0" {
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "+", vaccNo, fees, req.OrderNo, constants.VaReason_FEES, constants.VaOpType_Defreeze_Add); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
	}

	accNo := dao.VaccountDaoInst.GetAccNoFromVaccNo(tx, vaccNo)
	if accNo == "" {
		ss_log.Error("err=[根据虚账id查询账号id失败---->%s]", vaccNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, accNo); errStr != ss_err.ERR_SUCCESS {
	//	reply.ResultCode = errStr
	//	return nil
	//}

	// todo 修改订单明细中 订单的状态为取消状态
	if errStr := dao.BillingDetailsResultsDaoInstance.UpdateOrderStatusFromLogNo(tx, constants.OrderStatus_Cancel, req.OrderNo); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//查询该订单的用户uid、该订单的金额、币种
	useAccountUid, orderAmount, orderBalanceType := dao.OutgoOrderDaoInst.GetAccNoByLogNo(tx, req.OrderNo)
	if useAccountUid == "" {
		ss_log.Error("获取取消取款订单信息出错,useAccountUid=[%v],orderAmount=[%v],orderBalanceType=[%v]", useAccountUid, orderAmount, orderBalanceType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//添加取款失败推送消息
	errAddLogAppMessages := dao.LogAppMessagesDaoInst.AddLogAppMessagesTx(tx, req.OrderNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_OUTGO_Fail, constants.VaReason_OUTGO, useAccountUid, constants.OrderStatus_Paid)
	if errAddLogAppMessages != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddLogAppMessages=[%v]", errAddLogAppMessages)
		reply.ResultCode = errAddLogAppMessages
		return nil
	}

	appLang, _ := dao.AccDaoInstance.QueryAccountLang(useAccountUid)
	if appLang == "" {
		req.Lang = constants.LangEnUS
	} else {
		req.Lang = appLang
	}

	ss_log.Info("用户 %s 当前的语言为--->%s", useAccountUid, req.Lang)

	accountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(useAccountUid)
	switch orderBalanceType {
	case "usd":
		orderAmount = strext.ToStringNoPoint(strext.ToFloat64(orderAmount) / 100)
	case "khr":
		orderAmount = strext.ToStringNoPoint(strext.ToFloat64(orderAmount))
	}

	timeString := time.Now().Format("2006-01-02 15:04:05")
	args := []string{
		timeString, orderAmount, orderBalanceType,
	}
	lang, _ := dao.AccDaoInstance.QueryAccountLang(useAccountUid)
	if lang == "" || lang == constants.LangEnUS {
		args = []string{
			orderAmount, orderBalanceType, timeString,
		}
	}
	// 消息推送
	ev := &go_micro_srv_push.PushReqest{
		Accounts: []*go_micro_srv_push.PushAccout{
			{
				AccountNo:   useAccountUid,
				AccountType: accountType,
			},
		},
		TempNo: constants.Template_WithdrawFail,
		Args:   args,
	}

	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("error publishing: %v", err)
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
