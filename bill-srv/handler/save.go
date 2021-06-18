package handler

import (
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"context"
	"database/sql"
	"fmt"
	"time"

	"a.a/cu/db"
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/bill-srv/i"
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
	"a.a/mp-server/common/ss_struct"
)

// 存钱 操作是线下现金存款.
func (b *BillHandler) SaveMoney(ctx context.Context, req *go_micro_srv_bill.SaveMoneyRequest, reply *go_micro_srv_bill.SaveMoneyReply) error {
	servicerNo := ""
	servicerAccNo := ""
	var opAccType int // 1-服务商;2-收银员
	switch req.AccountType {
	case constants.AccountType_SERVICER: //服务商
		servicerNo = req.OpAccNo
		opAccType = constants.OpAccType_Servicer
		servicerAccNo = dao.RelaAccIdenDaoInst.GetAccNo(servicerNo, constants.AccountType_SERVICER)
	case constants.AccountType_POS: // 收银员
		servicerNo = dao.CashierDaoInst.GetServicerNoFromOpAccNo(req.OpAccNo)
		servicerAccNo = dao.ServiceDaoInst.GetAccNoFromSrvNo(servicerNo)
		opAccType = constants.OpAccType_Pos
	}
	// 判断是否有存款权限
	imcomePermission, _ := dao.ServiceDaoInst.GetPermissionFromSrvNo(servicerNo)
	if imcomePermission == "" || imcomePermission == "0" {
		ss_log.Error("err=[存款接口,没有存款权限,当前服务商id为----->%s,当前收款权限为----->%s]", servicerNo, imcomePermission)
		reply.ResultCode = ss_err.ERR_NOT_ROLE
		return nil
	}

	// 校验最大金额最小金额限制
	if err := CheckAmountIsMaxMinSave(req.MoneyType, req.Amount); err != nil {
		ss_log.Error("存款 最大最小金额校验失败,err:--->%s", err.Error())
		reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
		return nil
	}

	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[存钱, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
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

	// 确保虚拟账号存在
	var recvAccNo, isActived string
	recvAccNo, isActived = dao.AccDaoInstance.GetIsActiveAccNoFromPhone(req.RecvPhone, req.RecvCountryCode)

	if recvAccNo == "" {
		if err := dao.CountryCodePhoneDaoInst.Insert1(req.RecvCountryCode, req.RecvPhone); err != nil {
			ss_log.Error("新增或更改账号, 新增手机号和国家码进唯一表失败,accountUid: %s,err: %s", req.AccountUid, err.Error())
			reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
			return nil
		}
		recvAccNo = dao.AccDaoInstance.InsertEmptyAccount(req.RecvPhone, req.RecvCountryCode)

		//初始化钱包3、4
		if vaccNo := dao.VaccountDaoInst.InitVaccountNo(recvAccNo, constants.CURRENCY_USD, constants.VaType_FREEZE_USD_DEBIT); vaccNo == "" {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v", recvAccNo, constants.VaType_FREEZE_USD_DEBIT)
			//return err
		}
		if vaccNo := dao.VaccountDaoInst.InitVaccountNo(recvAccNo, constants.CURRENCY_KHR, constants.VaType_FREEZE_KHR_DEBIT); vaccNo == "" {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v", recvAccNo, constants.VaType_FREEZE_KHR_DEBIT)
			//return err
		}
	}

	// 通过币种获取手续费类型
	feesType, fErr := common.FeesTypeByMoneyType(constants.Scene_Save, req.MoneyType)
	if fErr != nil {
		ss_log.Error("FeesTypeByMoneyType err: %v", fErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 通过币种获取虚拟账号类型
	vaType, vaErr := common.VirtualAccountTypeByMoneyType(req.MoneyType, isActived)
	if vaErr != nil {
		ss_log.Error("VirtualAccountTypeByMoneyType err: %v", vaErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	var sendAccNo string
	sendAccNo = dao.AccDaoInstance.GetAccNoFromPhone(req.SendPhone, req.SaveCountryCode)
	recAccNo := dao.AccDaoInstance.GetAccNoFromPhone(req.RecvPhone, req.RecvCountryCode)

	if sendAccNo == "" {
		if err := dao.CountryCodePhoneDaoInst.Insert1(req.SaveCountryCode, req.SendPhone); err != nil {
			ss_log.Error("新增或更改账号, 新增手机号和国家码进唯一表失败,accountUid: %s,err: %s", req.AccountUid, err.Error())
			reply.ResultCode = ss_err.ERR_ACCOUNT_ALREADY_EXISTS
			return nil
		}

		sendAccNo = dao.AccDaoInstance.InsertEmptyAccount(req.SendPhone, req.SaveCountryCode)

		//初始化钱包3、4
		if vaccNo := dao.VaccountDaoInst.InitVaccountNo(sendAccNo, constants.CURRENCY_USD, constants.VaType_FREEZE_USD_DEBIT); vaccNo == "" {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v", sendAccNo, constants.VaType_FREEZE_USD_DEBIT)
			//return err
		}
		if vaccNo := dao.VaccountDaoInst.InitVaccountNo(sendAccNo, constants.CURRENCY_KHR, constants.VaType_FREEZE_KHR_DEBIT); vaccNo == "" {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v", sendAccNo, constants.VaType_FREEZE_KHR_DEBIT)
			//return err
		}
	}

	recvVaccNo := InternalCallHandlerInst.ConfirmExistVAccount(recvAccNo, req.MoneyType, strext.ToInt32(vaType))

	//  获取手续费
	rate, fees, feeErr := doFees(feesType, req.Amount)
	if feeErr != nil {
		ss_log.Error("存款计算手续费失败,err:--->%s", feeErr.Error())
		reply.ResultCode = ss_err.ERR_RATE_FAILD
		return nil
	}
	ss_log.Info("req.amount为----->%s,fees为----->%s", req.Amount, fees)

	if strext.ToFloat64(req.Amount) < strext.ToFloat64(fees) { // 判断手续费是否大于存款金额
		reply.ResultCode = ss_err.ERR_PAY_MISSING_EXCHANGE_RATE
		return nil
	}
	toAmount := ss_count.Sub(req.Amount, fees).String()

	// 判断上一个操作员存进的单跟现在的时间相比是否超过5秒
	fiveSec := ss_time.Now(global.Tz).Add(-5 * time.Second).Format(ss_time.DateTimeDashFormat)
	createTime, fErr := dao.IncomeOrderDaoInst.QueryCreateTime(req.OpAccNo, fiveSec)
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

	incomeOrder := new(dao.IncomeOrderDao)
	incomeOrder.RecvAccNo = recvAccNo
	incomeOrder.RecvVAccNo = recvVaccNo
	incomeOrder.Amount = req.Amount
	incomeOrder.ActAccNo = sendAccNo
	incomeOrder.ServicerNo = servicerNo
	incomeOrder.Fees = fees
	incomeOrder.BalanceType = req.MoneyType
	incomeOrder.PaymentType = constants.ORDER_PAYMENT_TYPE_CASH
	incomeOrder.ReeRate = rate
	incomeOrder.RealAmount = toAmount
	incomeOrder.OpAccNo = req.OpAccNo
	incomeOrder.OpAccType = opAccType
	logNo := dao.IncomeOrderDaoInst.InsertIncomeOrderV3(incomeOrder)
	if logNo == "" {
		ss_log.Error("---------->%s", "incomeOrder 插入失败")
		reply.ResultCode = ss_err.ERR_PAY_SAVE_MONEY
		return nil
	}

	//风控
	riskReply, riskErr := i.RiskCtrlHandleInstance.Client.GetRiskCtrlReuslt(context.TODO(), &go_micro_srv_riskctrl.GetRiskCtrlResultRequest{
		ApiType: constants.Risk_Ctrl_Save_Money,
		// 发起支付的账号
		PayerAccNo: sendAccNo,
		ActionTime: time.Now().String(),
		Amount:     req.Amount,
		Ip:         req.Ip,
		PayType:    constants.Risk_Pay_Type_Save_Money, //
		// 收款人账号
		PayeeAccNo:  recvAccNo,
		ProductType: constants.Risk_Ctrl_Save_Money,
		// 币种
		MoneyType: req.MoneyType,
		// 订单号
		OrderNo: logNo,
	})

	if riskErr != nil {
		ss_log.Error("GetRiskCtrlReuslt 返回错误,err: %v", riskErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	ss_log.Info("存款 风控返回结果,操作结果是---->%s,RiskNo为----->%s", riskReply.OpResult, riskReply.RiskNo)

	if riskReply.OpResult == constants.Risk_Result_No_Pass_Str {
		reply.ResultCode = ss_err.ERR_RISK_IS_RISK
		reply.RiskNo = riskReply.RiskNo
		return nil
	}

	// 获取tm服务代理
	tmProxy, fErr := ss_struct.NewTmServerProxy(common.BillServerFullId)
	if fErr != nil {
		ss_log.Error("ss_struct.NewTmServer err: %v", fErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 开启tm事务
	if fErr := tmProxy.TxBegin(); fErr != nil {
		ss_log.Error("ss_struct.TxBegin err: %v", fErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	//============================================================================
	// 调用ps rpc,判断金额是否超额,预存
	quotaReq := &go_micro_srv_quota.ModifyQuotaRequest{
		CurrencyType: req.MoneyType,
		Amount:       req.Amount,
		AccountNo:    servicerAccNo,
		OpType:       constants.QuotaOp_PreSave,
		LogNo:        logNo,
		TxNo:         tmProxy.GetTxNo(),
	}
	var quotaErr error
	quotaRepl := &go_micro_srv_quota.ModifyQuotaReply{}
	quotaRepl, quotaErr = i.QuotaHandleInstance.Client.ModifyQuota(context.TODO(), quotaReq)
	if quotaErr != nil {
		ss_log.Error("调用远程ModifyQuota rpc失败,err: %v", quotaErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
		tmProxy.TxRollback()
		ss_log.Error("err=[--------------->%s]", "客户存钱,调用八神的服务失败,操作为客户预存款")

		reply.ResultCode = quotaRepl.ResultCode
		return nil
	}

	// 调用 rpc 存款
	quotaReq.OpType = constants.QuotaOp_Save
	quotaRepl, quotaErr = i.QuotaHandleInstance.Client.ModifyQuota(context.TODO(), quotaReq)
	if quotaErr != nil {
		ss_log.Error("调用远程 ModifyQuota rpc失败,err: %v", quotaErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
		tmProxy.TxRollback()
		ss_log.Error("err=[--------------->%s]", "客户存钱,调用八神的服务失败,操作为客户存款")

		reply.ResultCode = quotaRepl.ResultCode
		return nil
	}
	//============================================================================
	// 生成存款码
	code := util.RandomDigitStrOnlyNum(constants.Len_SaveMoneyCode)

	ss_log.Info("客户预存款成功")
	// 判断收款人是否存在账户
	isActive := dao.AccDaoInstance.GetIsActiveFromPhone(req.RecvPhone, req.RecvCountryCode)
	if isActive != constants.AccountActived { // 账户不存在才生成码
		nowTimeStr := ss_time.Now(global.Tz).Format(ss_time.DateTimeSlashFormat)

		endTimeStr, err := dao.WriteoffInst.GetCodeEndTimeStr(nowTimeStr)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		log := new(dao.Writeoff)
		log.Code = code
		log.IncomeOrderNo = logNo
		log.UseStatus = constants.WriteOffCodeWaitUse
		log.SendAccNo = sendAccNo
		log.SendPhone = req.SendPhone
		log.RecvPhone = req.RecvPhone
		log.CreateTime = nowTimeStr
		log.DurationTime = endTimeStr
		log.RecvAccNo = recAccNo
		initWriteErr := dao.WriteoffInst.InitWriteoffV2(tmProxy, log)
		if initWriteErr != nil {
			ss_log.Error("InitWriteoffV2 ,initWriteErr : %v", initWriteErr)
			//if errStr := dao.IncomeOrderDaoInst.UpdateIncomeOrderOrderStatus(tx, logNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
			if err := dao.IncomeOrderDaoInst.UpdateIncomeOrderOrderStatusV2(tmProxy, logNo, constants.OrderStatus_Err); err != nil {
				ss_log.Error("UpdateIncomeOrderOrderStatusV2,err : %v", err)
				reply.ResultCode = ss_err.ERR_PARAM
				return nil
			}
			// 回滚
			tmProxy.TxRollback()

			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	// 修改虚账
	//err := dao.VaccountDaoInst.SaveMoneyAccFromAToBUpperZero(tx, fromVaccNo, recvVaccNo, req.Amount, toAmount, logNo, constants.VaReason_INCOME)
	err := dao.VaccountDaoInst.ModifyVaccRemainUpperZeroV2(tmProxy, recvVaccNo, toAmount, "+", logNo, constants.VaReason_INCOME)
	if err != nil {
		ss_log.Error("ModifyVaccRemainUpperZeroV2 ,err : %v", err)
		if errStr := dao.IncomeOrderDaoInst.UpdateIncomeOrderOrderStatusV2(tmProxy, logNo, constants.OrderStatus_Err); errStr != nil {
			ss_log.Error("UpdateIncomeOrderOrderStatusV2,err : %v", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		// 回滚
		tmProxy.TxRollback()
		reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
		return nil
	}

	if errStr := dao.IncomeOrderDaoInst.UpdateIncomeOrderOrderStatusV2(tmProxy, logNo, constants.OrderStatus_Paid); errStr != nil {
		ss_log.Error("UpdateIncomeOrderOrderStatusV2,err : %v", errStr)
		// 回滚
		tmProxy.TxRollback()
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// todo 插入 billing_details_results
	if err := dao.BillingDetailsResultsDaoInstance.InsertResultV2(tmProxy, req.Amount, req.MoneyType, servicerAccNo, req.AccountType,
		logNo, "0", constants.OrderStatus_Paid, servicerNo, req.OpAccNo, constants.BillDetailTypeIn, fees, req.Amount); err != nil {
		ss_log.Error("InsertResultV2,err : %v", err)
		// 回滚
		tmProxy.TxRollback()
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//if err := dao.VaccountDaoInst.SyncAccRemainV2(tmProxy, recvAccNo); err != nil {
	//	ss_log.Error("SyncAccRemainV2,err : %v", err)
	//	// 回滚
	//	tmProxy.TxRollback()
	//
	//	reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
	//	return nil
	//}

	// 提交tm事务
	if commitErr := tmProxy.TxCommit(); commitErr != nil {
		ss_log.Error("TxCommit err: %v", commitErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	dbHandler.SetConnMaxLifetime(-1)
	tx, txErr := dbHandler.BeginTx(ctx, nil)
	if txErr != nil {
		ss_log.Error("dbHandler.BeginTx 开启事务失败,err: %v", txErr)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	defer ss_sql.Rollback(tx)
	appLang, _ := dao.AccDaoInstance.QueryAccountLang(recAccNo)
	if appLang == "" {
		req.Lang = constants.LangEnUS
	} else {
		req.Lang = appLang
	}
	ss_log.Info("用户 %s 当前的语言为--->%s", recAccNo, req.Lang)

	// 收款人账号是激活的,推送消息
	if isActive == constants.AccountActived {
		//添加收款账号的推送消息
		if errStr := dao.LogAppMessagesDaoInst.AddLogAppMessagesTx(tx, logNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_INCOME_Apply, constants.VaReason_INCOME, recAccNo, constants.OrderStatus_Paid); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}

		recAccountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(recAccNo)
		moneyType := dao.LangDaoInstance.GetLangTextByKey(req.MoneyType, req.Lang)

		// 修正各币种的金额
		amount := common.NormalAmountByMoneyType(req.MoneyType, req.Amount)

		timeString := time.Now().Format("2006-01-02 15:04:05")
		ss_time.GetDayBefore()
		args := []string{
			timeString, amount, moneyType,
		}
		lang, _ := dao.AccDaoInstance.QueryAccountLang(recvAccNo)
		if lang == "" || lang == constants.LangEnUS {
			args = []string{
				amount, moneyType, timeString,
			}
		}
		// 消息推送
		ev := &go_micro_srv_push.PushReqest{
			Accounts: []*go_micro_srv_push.PushAccout{
				{
					AccountNo:   recvAccNo,
					AccountType: recAccountType,
				},
			},
			TempNo: constants.Template_AddSuccess,
			Args:   args,
		}

		ss_log.Info("publishing %+v\n", ev)
		// publish an event
		if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
			ss_log.Error("error publishing: %v", err)
		}
	}

	ss_sql.Commit(tx)

	if fees != "0" && fees != "" {
		// 发送手续费进MQ
		feeEv := &go_micro_srv_settle.SettleTransferRequest{
			BillNo:    logNo,
			FeesType:  constants.FEES_TYPE_SAVEMONEY,
			Fees:      fees,
			MoneyType: req.MoneyType,
			VaType:    strext.ToInt32(vaType),
		}
		ss_log.Info("publishing %+v\n", feeEv)
		// publish an event
		if err := common.SettleEvent.Publish(context.TODO(), feeEv); err != nil {
			ss_log.Error("err=[pos 存款接口,手续费推送到MQ失败,err----->%s]", err.Error())
		}
	}

	////发送短信核销码
	//if isActive != constants.AccountActived {
	//	// 消息推送
	//	ev := &go_micro_srv_push.PushReqest{
	//		Accounts: []*go_micro_srv_push.PushAccout{
	//			{
	//				Phone:       req.SendPhone,
	//				Lang:        req.Lang,
	//				CountryCode: req.SaveCountryCode,
	//			},
	//		},
	//		TempNo: constants.Template_SmsWriteOff,
	//		Args: []string{
	//			code,
	//		},
	//	}
	//	ss_log.Info("publishing %+v\n", ev)
	//	// publish an event
	//	if err := common.WriteOffEvent.Publish(context.TODO(), ev); err != nil {
	//		ss_log.Error("err=[pos 存款接口,核销码推送到MQ失败,err----->%s]", err.Error())
	//	}
	//	ss_log.Info("短信核销码推送到队列成功,code: %s", code)
	//}
	reply.RiskNo = riskReply.RiskNo
	reply.OrderNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//这个校验的是支付密码
func checkoutPWD(accountType, opAccNo, nonStr, reqPWD string) string {
	var pwd string // 数据库密码
	switch accountType {
	case constants.AccountType_SERVICER: //服务商
		pwd = dao.ServiceDaoInst.GetServicerPWDFromOpAccNo(opAccNo)
	case constants.AccountType_POS: // 收银员
		pwd = dao.CashierDaoInst.GetCashierPwdFromOpAccNo(opAccNo)
	case constants.AccountType_USER:
		pwd = dao.CustDaoInstance.QueryPwdFromOpAccNo(opAccNo)
	}
	if pwd == "" {
		ss_log.Error("----->%s", "数据库密码为空")
		return ss_err.ERR_DB_PWD
	}
	// 校验操作员密码
	pwdMD5FixedDB := encrypt.DoMd5Salted(pwd, nonStr)
	if reqPWD != pwdMD5FixedDB {
		ss_log.Error("----->%s", "密码错误")
		return ss_err.ERR_DB_PWD
	}
	return ss_err.ERR_SUCCESS
}

// 查询存款小票信息
func (b *BillHandler) QuerySaveReceipt(ctx context.Context, req *go_micro_srv_bill.QuerySaveReceiptRequest, reply *go_micro_srv_bill.QuerySaveReceiptReply) error {
	if req.OrderNo == "" {
		ss_log.Error("OrderNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	incomeOrder, err := dao.IncomeOrderDaoInst.QueryIncomeOrder(req.OrderNo, "3")
	if err != nil {
		ss_log.Error("查询存款小票信息失败，log_no:%v, err:%v", req.OrderNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 查询手机号
	savePhone, saveCountryCode, pcErr := dao.AccDaoInstance.GetPhoneCountryCodeFromAccNo(incomeOrder.ActAccNo)
	if pcErr != nil {
		ss_log.Error("QuerySaveReceipt GetPhoneCountryCodeFromAccNo 查询存款人的手机号和国家码失败,err: %s", pcErr.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	recPhone, recCountryCode, recErr := dao.AccDaoInstance.GetPhoneCountryCodeFromAccNo(incomeOrder.RecvAccNo)
	if pcErr != nil {
		ss_log.Error("QuerySaveReceipt GetPhoneCountryCodeFromAccNo 查询收款人的手机号和国家码失败,err: %s", recErr.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 查询pos机号
	num := dao.ServicerTerminalDaoInstance.QueryNumberFromServiceNo(incomeOrder.ServicerNo)

	//查询存款核销码
	saveCode, err := dao.WriteoffInst.GetCodeByIncomeOrderNo(req.OrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			saveCode = ""
		} else {
			ss_log.Error("查询存款核销码失败, IncomeOrderNo:%v, err:%v", req.OrderNo, err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	// 计算到账金额
	//arriveAmount := amount
	//applyAmount := ss_count.Add(amount, fees)

	// todo 向上取整数
	data := &go_micro_srv_bill.QuerySaveReceiptResult{
		OrderNo: req.OrderNo,
		// 商户号
		ServiceNo: incomeOrder.ServicerNo,
		// 终端编号
		TerminalNumber: num,
		// 存款手机号
		SavePhone: fmt.Sprintf("%s%s", saveCountryCode, savePhone),
		//核销码
		SaveCode: saveCode,
		// 收款手机号
		RecPhone: fmt.Sprintf("%s%s", recCountryCode, recPhone),
		// 申请金额
		ApplyAmount: incomeOrder.Amount,
		// 手续费
		Fees: incomeOrder.Fees,
		// 到账金额
		ArriveAmount: incomeOrder.RealAmount,
		// 日期
		Date:      incomeOrder.FinishTime,
		MoneyType: incomeOrder.BalanceType,
	}
	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 存款明细
func (b *BillHandler) SaveDetail(ctx context.Context, req *go_micro_srv_bill.SaveDetailRequest, reply *go_micro_srv_bill.SaveDetailReply) error {
	var amount, finishTime, saveAccount, recAccount, fees, balanceType, savePhone, recPhone, status string
	switch req.Type {
	case 1: // 存款 income
		order, err := dao.IncomeOrderDaoInst.QueryIncomeOrder(req.OrderNo, "")
		if err != nil {
			ss_log.Error("查询存款明细失败，OrderNo=%v, err=%v", req.OrderNo, err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		amount = order.Amount
		finishTime = order.FinishTime
		saveAccount = order.ActAccNo
		recAccount = order.RecvAccNo
		fees = order.Fees
		balanceType = order.BalanceType
		status = order.OrderStatus
		// 查询手机号
		savePhone = dao.AccDaoInstance.GetPhoneFromAccNo(saveAccount)
		recPhone = dao.AccDaoInstance.GetPhoneFromAccNo(recAccount)
	case 2: // 取款 outgo
		var vaccountNo string
		amount, _, finishTime, vaccountNo, fees, balanceType, status, _ = dao.OutgoOrderDaoInst.QueryOutGoOrderFromLogNo(req.OrderNo, "")
		// 查询取款人手机号
		recPhone = dao.AccDaoInstance.GetPhoneFromVAccNo(vaccountNo)
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	data := &go_micro_srv_bill.SaveDetailResult{
		OrderNo: req.OrderNo,
		// 存款手机号
		SavePhone: savePhone,
		// 收款手机号
		RecPhone: recPhone, // 取款的时候也是这个手机号字段
		// 申请金额
		Amount: amount,
		// 手续费
		Fees: fees,
		// 日期
		Date:      finishTime,
		MoneyType: balanceType,
		Status:    status,
	}
	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (b *BillHandler) SaveMoneyDetail(ctx context.Context, req *go_micro_srv_bill.SaveMoneyDetailRequest, reply *go_micro_srv_bill.SaveMoneyDetailReply) error {
	if req.OrderNo == "" {
		ss_log.Error("OrderNo参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	incomeOrder, err := dao.IncomeOrderDaoInst.QueryIncomeOrder(req.OrderNo, "")
	if err != nil {
		ss_log.Error("查询存款订单详情失败, log_no=%v, err=%v", req.OrderNo, err)
	}
	if incomeOrder.Amount == "" {
		ss_log.Error("err=[存款查询订单详情失败,订单号为----->%s]", req.OrderNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//获取产生该笔订单的账号
	account, errGet := dao.BillingDetailsResultsDaoInstance.GetAccountByOrderNo(req.OrderNo)
	if errGet != nil {
		ss_log.Error("获取产生该笔订单的账号失败，OrderNo=[%v],err=[%v]", req.OrderNo, errGet)
	}

	savePhone := dao.AccDaoInstance.GetPhoneFromAccNo(incomeOrder.ActAccNo)
	recPhone := dao.AccDaoInstance.GetPhoneFromAccNo(incomeOrder.RecvAccNo)

	saveCode, err := dao.WriteoffInst.GetCodeByIncomeOrderNo(req.OrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			saveCode = ""
		} else {
			ss_log.Error("查询存款核销码失败，OrderNo=%v, err=%v", req.OrderNo, err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	data := &go_micro_srv_bill.SaveMoneyDetailResult{
		OrderNo:      req.OrderNo,
		SaveAmount:   incomeOrder.Amount,
		ArriveAmount: incomeOrder.RealAmount,
		SendPhone:    savePhone,
		RecvPhone:    recPhone,
		Fees:         incomeOrder.Fees,
		Date:         incomeOrder.FinishTime,
		MoneyType:    incomeOrder.BalanceType,
		OrderStatus:  incomeOrder.OrderStatus,
		Account:      account,
		SaveCode:     saveCode,
	}
	reply.Data = data
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//服务商现金充值
func (b *BillHandler) AddCashRecharge(ctx context.Context, req *go_micro_srv_bill.AddCashRechargeRequest, reply *go_micro_srv_bill.AddCashRechargeReply) error {

	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[存钱, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	//验证登录密码
	if !dao.AccDaoInstance.CheckAdminLoginPWD(req.Uid, req.Password, req.NonStr) {
		ss_log.Error("登录密码错误")
		reply.ResultCode = ss_err.ERR_DB_PWD
		return nil
	}

	//根据账号查询uid
	serAccNo, errGet := dao.AccDaoInstance.GetAccNoFromAccount(req.AccAccount)
	if errGet != nil {
		ss_log.Error("err=[%v]", errGet)
		reply.ResultCode = ss_err.ERR_ACC_NO_SER_FAILD
		return nil
	}

	idenNo := dao.RelaAccIdenDaoInst.GetIdenNo(serAccNo, constants.AccountType_SERVICER)
	if idenNo == "" {
		ss_log.Error("查询账号[%v]服务商id出错", serAccNo)
		reply.ResultCode = ss_err.ERR_ACC_NO_SER_FAILD
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//添加现金充值日志
	logNo := dao.CashRechargeOrderDaoInst.InsertCashRecharge(tx, dao.CashRechargeOrderDao{
		AccNo:        serAccNo,
		Amount:       req.Amount,
		IdenNo:       idenNo,
		OrderStatus:  constants.OrderStatus_Paid,
		CurrencyType: req.CurrencyType,
		OpAccNo:      req.Uid,
		PaymentType:  constants.ORDER_PAYMENT_TYPE_CASH,
		Notes:        req.Notes,
	})

	if logNo == "" {
		ss_log.Error("插入服务商现金充值表错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 调用ps rpc,进行加服务商余额
	quotaRepl, quotaErr := i.QuotaHandleInstance.Client.ModifyQuota(context.TODO(), &go_micro_srv_quota.ModifyQuotaRequest{
		CurrencyType: req.CurrencyType,
		Amount:       req.Amount,
		AccountNo:    serAccNo,
		OpType:       constants.QuotaOp_SvrCashRecharge,
		LogNo:        logNo,
	})
	if quotaErr != nil {
		ss_log.Error("调用远程ModifyQuota rpc失败,err: %v", quotaErr)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[--------------->%s]", "调用服务失败,操作为服务商现金充值")
		reply.ResultCode = quotaRepl.ResultCode
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
