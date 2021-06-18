package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/bill-srv/i"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	go_micro_srv_settle "a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type BillHandler struct{}

var BillHandlerInst BillHandler

func (b *BillHandler) ExchangeAmount(ctx context.Context, req *go_micro_srv_bill.ExchangeAmountRequest, reply *go_micro_srv_bill.ExchangeAmountReply) error {
	if strext.ToFloat64(req.Amount) <= 0 {
		ss_log.Error("ExchangeAmount 交易金额不能小于0,req.amount: %s", req.Amount)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	var amount string
	switch req.Type {
	case constants.ExchangeAmountUsdToKhr:
		_, usdToKhr, err := cache.ApiDaoInstance.GetGlobalParam("usd_to_khr")
		if err != nil {
			ss_log.Error("ExchangeAmount 获取 usd_to_khr 费率失败, err: %s", err.Error())
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		ss_log.Info("ExchangeAmount usd_to_khr 兑换费率为: %s", usdToKhr)
		amount = InternalCallHandlerInst.ExchangeUsdToKhr(req.Amount, usdToKhr)
	case constants.ExchangeAmountKhrToUsd:
		_, khrToUsd, err := cache.ApiDaoInstance.GetGlobalParam("khr_to_usd")
		if err != nil {
			ss_log.Error("ExchangeAmount 获取 khr_to_usd 费率失败, err: %s", err.Error())
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		ss_log.Info("ExchangeAmount khr_to_usd 兑换费率为: %s", khrToUsd)
		syncAmount, exErr := InternalCallHandlerInst.ExchangeKhrToUsd(req.Amount, khrToUsd)
		if exErr != nil {
			ss_log.Error(exErr.Error())
			reply.ResultCode = ss_err.ERR_PAY_EXCHANGE_MIN_AMOUNT
			return nil
		}
		amount = syncAmount
	default:
		ss_log.Error("兑换类型错误,type: %d", req.Type)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	reply.Amount = amount
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (b *BillHandler) QueryCustHasPwd(ctx context.Context, req *go_micro_srv_bill.QueryCustHasPwdRequest, reply *go_micro_srv_bill.QueryCustHasPwdReply) error {
	// 判断是否有支付密码
	if req.AccountType == constants.AccountType_USER {
		if pwd := dao.CustDaoInstance.QueryPwdFromOpAccNo(req.IdenNo); pwd == "" {
			reply.ResultCode = ss_err.ERR_PAY_PWD_IS_NULL
			return nil
		}
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//todo 此处代码已移至InternalCallHandlerInst，后期没有问题就删了
func (b BillHandler) ExchangeUsdToKhr(amount, rate string) string {
	a := ss_count.Multiply(amount, rate)
	result := ss_count.Div(a.String(), "100") // 结果需要处于100然后取整
	// 取整数
	return ss_big.SsBigInst.ToRound(result, 0, ss_big.RoundingMode_HALF_EVEN).String()
}

//todo 此处代码已移至InternalCallHandlerInst，后期没有问题就删了
func (b BillHandler) ExchangeKhrToUsd(amount, rate string) (string, error) {
	reqAmount := ss_count.Multiply(amount, "100") // 先把khr *100后再比较,因为 usd 的单位是分

	// 判断是否能足够兑换1美元
	_, exchangeRateAmount, _ := cache.ApiDaoInstance.GetGlobalParam(constants.Exchange_Khr_To_Usd)
	reqAmountF := strext.StringToFloat64(amount)
	exchangeRateF := strext.StringToFloat64(exchangeRateAmount) / 100
	if reqAmountF < exchangeRateF { // 最低兑换金额为0.01USD
		return "", errors.New(fmt.Sprintf("khr_to_usd, khr 的余额不足,khr的余额为---> %v,最低需要为---> %v", reqAmountF, exchangeRateF))
	}

	fromString := ss_count.Div(reqAmount.String(), rate)
	// 美元  四舍六入五成双
	return ss_big.SsBigInst.ToRound(fromString, 0, ss_big.RoundingMode_HALF_EVEN).String(), nil
}

// 兑换
func (b BillHandler) Exchange(ctx context.Context, req *go_micro_srv_bill.ExchangeRequest, reply *go_micro_srv_bill.ExchangeReply) error {
	ss_log.Info("Exchange 请求参数: %+v", req)
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// todo 判断该用户是否有交易的权限,cust表中的 payment_status 字段
	tradingAuthority, _, _, _, _ := dao.CustDaoInstance.QueryRateRoleFrom(req.OpAccNo)
	if strext.ToInt(tradingAuthority) == constants.TradingAuthorityForbid {
		ss_log.Error("Exchange 输错密码超过限制,禁止交易,accountNo: %s,custNo: %s", req.AccountNo, req.OpAccNo)
		reply.ResultCode = ss_err.ERR_Payment_Pwd_Count_Limit
		return nil
	}

	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[兑换, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_PAY_AMOUNT_NOT_MIN_ZERO
		return nil
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if !dao.AppFingerprintDaoInstance.CheckSignKey(req.AccountNo, req.DeviceUuid, req.SignKey) {
			ss_log.Error("查询不到指纹支付标识，或者指纹支付标识状态为无效, AccountNo[%v], DeviceUuid[%v], SignKey[%v]", req.AccountNo, req.DeviceUuid, req.SignKey)
			reply.ResultCode = ss_err.ERR_AppFingerprint_FAILD
			return nil
		}

	} else {
		//验证支付密码
		replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &go_micro_srv_auth.CheckPayPWDRequest{
			AccountUid:  req.AccountNo,
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
	}

	x := fmt.Sprintf("%s_to_%s", req.InType, req.OutType)
	_, rate, _ := cache.ApiDaoInstance.GetGlobalParam(x)
	rateI := strext.ToInt64(rate)
	if rateI <= 0 {
		reply.ResultCode = ss_err.ERR_PAY_MISSING_EXCHANGE_RATE
		return nil
	}

	khrVaccountNo := InternalCallHandlerInst.ConfirmExistVAccount(req.AccountNo, constants.CURRENCY_KHR, constants.VaType_KHR_DEBIT)
	usdVaccoutnNo := InternalCallHandlerInst.ConfirmExistVAccount(req.AccountNo, constants.CURRENCY_USD, constants.VaType_USD_DEBIT)

	var moneyType string
	var fromVacc, toVacc string
	var syncAmount string
	var feeType int32

	switch x {
	case constants.Exchange_Usd_To_Khr:
		moneyType = req.InType
		fromVacc = usdVaccoutnNo
		toVacc = khrVaccountNo
		feeType = constants.Fees_Type_Usd_To_Khr_Count_Fee

		// 取整数
		syncAmount = InternalCallHandlerInst.ExchangeUsdToKhr(req.Amount, rate)

	case constants.Exchange_Khr_To_Usd:
		moneyType = req.OutType
		fromVacc = khrVaccountNo
		toVacc = usdVaccoutnNo
		feeType = constants.Fees_Type_Khr_To_Usd_Count_Fee

		var exErr error
		// 美元  四舍六入五成双
		syncAmount, exErr = InternalCallHandlerInst.ExchangeKhrToUsd(req.Amount, rate)
		if exErr != nil {
			ss_log.Error(exErr.Error())
			reply.ResultCode = ss_err.ERR_PAY_EXCHANGE_MIN_AMOUNT
			return nil
		}
	default:
		reply.ResultCode = ss_err.ERR_PAY_PAY_TYPE_NOT_SUPPORT
		return nil
	}

	//  获取手续费
	_, fees, feeErr := doFees(feeType, req.Amount)
	if feeErr != nil {
		ss_log.Error("兑换 计算手续费失败,err:--->%s", feeErr.Error())
		reply.ResultCode = ss_err.ERR_RATE_FAILD
		return nil
	}
	ss_log.Info("%s 的手续费为----------->%s", x, fees)

	logNo := dao.ExchangeOrderDaoInst.InsertExchangeOrder(tx, req.AccountNo, req.InType, req.OutType, req.Amount, rate, "app", syncAmount, "", fees, req.Lat, req.Lng, req.Ip)
	//风控
	riskReply, _ := i.RiskCtrlHandleInstance.Client.GetRiskCtrlReuslt(context.TODO(), &go_micro_srv_riskctrl.GetRiskCtrlResultRequest{
		ApiType: constants.Risk_Ctrl_Exchange,
		// 发起支付的账号
		PayerAccNo: req.AccountNo,
		ActionTime: time.Now().String(),
		Amount:     req.Amount,
		Ip:         req.Ip,
		PayType:    constants.Risk_Pay_Type_Exchange, //
		// 收款人账号
		PayeeAccNo:  req.AccountNo,
		ProductType: constants.Risk_Ctrl_Exchange,
		// 币种
		MoneyType: req.InType,
		// 订单号
		OrderNo: logNo,
	})

	ss_log.Info("%s 兑换 风控返回结果,操作结果是---->%s,RiskNo为----->%s", x, riskReply.OpResult, riskReply.RiskNo)

	if riskReply.OpResult == constants.Risk_Result_No_Pass_Str {
		reply.ResultCode = ss_err.ERR_RISK_IS_RISK
		reply.RiskNo = riskReply.RiskNo
		return nil
	}

	if retCode := dao.VaccountDaoInst.SameAccFromAToBUpperZero(tx, req.Amount, req.AccountNo, fromVacc, toVacc, syncAmount, logNo, constants.VaReason_Exchange); retCode != ss_err.ERR_SUCCESS {
		reply.ResultCode = retCode
		return nil
	}

	if fees != "0" {
		// 修改手续费
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, fromVacc, fees, "-", logNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, req.AccountNo); errStr != ss_err.ERR_SUCCESS {
		//	reply.ResultCode = errStr
		//	return nil
		//}
	}

	//if reply.ResultCode == ss_err.ERR_SUCCESS {
	if errStr := dao.ExchangeOrderDaoInst.UpdateExchangeOrderStatus(tx, logNo, constants.OrderStatus_Paid, "0"); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if err := dao.LogAppFingerprintPayDaoInstance.AddTx(tx, &dao.LogAppFingerprintPayData{
			AccountNo:    req.AccountNo,
			DeviceUuid:   req.DeviceUuid,
			SignKey:      req.SignKey,
			OrderNo:      logNo,
			OrderType:    constants.VaReason_Exchange,
			Amount:       req.Amount,
			CurrencyType: req.InType,
		}); err != nil {
			ss_log.Error("指纹无密码支付插入日志失败,err[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	ss_sql.Commit(tx)

	if fees != "0" {

		// 发送手续费进MQ
		ev := &go_micro_srv_settle.SettleTransferRequest{
			BillNo:    logNo,
			FeesType:  constants.FEES_TYPE_EXCHANGE,
			Fees:      fees,
			MoneyType: moneyType,
		}
		ss_log.Info("publishing %+v\n", ev)
		// publish an event
		if err := common.SettleEvent.Publish(context.TODO(), ev); err != nil {
			ss_log.Error("err=[兑换接口,手续费推送到MQ失败,err----->%s]", err.Error())
		}
	}
	reply.RiskNo = riskReply.RiskNo
	reply.OrderNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 转账
func (b BillHandler) Transfer(ctx context.Context, req *go_micro_srv_bill.TransferRequest, reply *go_micro_srv_bill.TransferReply) error {
	if req.PaymentMethod == constants.PayMethodBankCard {
		ss_log.Error("暂不支持银行卡支付")
		reply.ResultCode = ss_err.ERR_Bank_Card_Not_Supported
		return nil
	}

	// 判断转账金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[转账金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	// 检验最大最小金额是否超出限制金额
	if err := CheckAmountIsMaxMinTransfer(req.MoneyType, req.Amount); err != nil {
		ss_log.Error("Transfer 转账最大最小金额校验失败,err:--->%s", err.Error())
		reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
		return nil
	}

	custInfo, err := dao.CustDaoInstance.QueryCustInfo(req.FromAccountNo)
	if err != nil {
		ss_log.Error("查询用户信息失败, accountNo=%v, err=%v", req.FromAccountNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	//检查用户是否已实名
	if custInfo.IndividualAuthStatus != constants.AuthMaterialStatus_Passed {
		ss_log.Error("用户账号未实名认证, AccountNo:%v, IndividualAuthStatus=%v", req.FromAccountNo, custInfo.IndividualAuthStatus)
		reply.ResultCode = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
		return nil
	}

	//判断该用户是否有交易的权限
	if strext.ToInt(custInfo.TradingAuthority) == constants.TradingAuthorityForbid {
		ss_log.Error("输错密码超过限制,禁止交易,accountNo: %s,custNo: %s", req.FromAccountNo, req.IdenNo)
		reply.ResultCode = ss_err.ERR_Payment_Pwd_Count_Limit
		return nil
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if !dao.AppFingerprintDaoInstance.CheckSignKey(req.FromAccountNo, req.DeviceUuid, req.SignKey) {
			ss_log.Error("查询不到指纹支付标识，或者指纹支付标识状态为无效, FromAccountNo[%v], DeviceUuid[%v], SignKey[%v]", req.FromAccountNo, req.DeviceUuid, req.SignKey)
			reply.ResultCode = ss_err.ERR_AppFingerprint_FAILD
			return nil
		}

	} else {
		if req.Password == "" || req.NonStr == "" {
			ss_log.Error("Transfer 验证支付密码的必要参数为空 Password: %s,  NonStr: %s", req.Password, req.NonStr)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//支付密码校验，如果支付密码正确且不是限制中则清除连续支付密码出错次数
		replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &go_micro_srv_auth.CheckPayPWDRequest{
			AccountUid:  req.FromAccountNo,
			AccountType: req.AccountType,
			Password:    req.Password,
			NonStr:      req.NonStr,
			IdenNo:      req.IdenNo,
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

	}

	// 转账的人必须是有账号的,必须是激活的,此处的isActived为激活状态1
	vaType, vaErr1 := common.VirtualAccountTypeByMoneyType(req.MoneyType, "1")
	if vaErr1 != nil {
		ss_log.Error("Transfer 转账的人必须是有账号的 err: %s", vaErr1.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 判断转账人的账号是否有转出权限
	fromCustNo := dao.RelaAccIdenDaoInst.GetIdenNo(req.FromAccountNo, constants.AccountType_USER)
	if _, _, _, _, outTransferAuthorizationT := dao.CustDaoInstance.QueryRateRoleFrom(fromCustNo); outTransferAuthorizationT == "" || outTransferAuthorizationT == "0" {
		ss_log.Error("Transfer 该用户没转出权限,用户id为----->%s", fromCustNo)
		reply.ResultCode = ss_err.ERR_PAY_NO_OUT_GO_PERMISSION
		return nil
	}

	// 判断收款人是否存在账户
	toAccountNo := dao.AccDaoInstance.ConfirmAccIsExit(req.ToPhone, req.CountryCode)
	if toAccountNo == "" {
		ss_log.Error("Transfer 查询用户的账号uuid为空,用户ToPhone为[%v], CountryCode[%v]", req.ToPhone, req.CountryCode)
		reply.ResultCode = ss_err.ERR_PhoneUnRegisteredUser
		return nil
	}

	isActive := dao.AccDaoInstance.GetIsActiveFromPhone(req.ToPhone, req.CountryCode)
	if isActive == constants.AccountActived {
		toCustNo := dao.RelaAccIdenDaoInst.GetIdenNo(toAccountNo, constants.AccountType_USER)
		if _, _, _, inTransferAuthorizationT, _ := dao.CustDaoInstance.QueryRateRoleFrom(toCustNo); inTransferAuthorizationT == "" || inTransferAuthorizationT == "0" {
			ss_log.Error("Transfer 该用户没转入权限,用户id为----->%s", toCustNo)
			reply.ResultCode = ss_err.ERR_PAY_NO_IN_COME_PERMISSION
			return nil
		}
	}

	//收款人虚账类型
	toVaType, vaErr := common.VirtualAccountTypeByMoneyType(req.MoneyType, isActive)
	if vaErr != nil {
		ss_log.Error("Transfer 通过币种获取虚拟账号类型失败 err: %s", vaErr.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	//收款人虚账
	toVacc := InternalCallHandlerInst.ConfirmExistVAccount(toAccountNo, req.MoneyType, strext.ToInt32(toVaType))

	//付款人虚账
	fromVacc := InternalCallHandlerInst.ConfirmExistVAccount(req.FromAccountNo, req.MoneyType, strext.ToInt32(vaType))

	if fromVacc == toVacc {
		reply.ResultCode = ss_err.ERR_ACCOUNT_TRANSFER_TO_SELF
		return nil
	}

	// 通过币种获取手续费类型
	feesType, fErr := common.FeesTypeByMoneyType(constants.Scene_Transfer, req.MoneyType)
	//  获取手续费
	if fErr != nil {
		ss_log.Error("Transfer 通过币种获取手续费类型失败 err: %s", fErr.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	rate, fees, feeErr := doFees(feesType, req.Amount)
	if feeErr != nil {
		ss_log.Error("转账 计算手续费失败,err:--->%s", feeErr.Error())
		reply.ResultCode = ss_err.ERR_RATE_FAILD
		return nil
	}
	ss_log.Info("Transfer  req.amount为----->%s,fees为----->%s", req.Amount, fees)

	// 计算本金加手续费的余额
	allAmount := ss_count.Add(req.Amount, fees) // 加上手续费后的amount
	balance, _ := dao.VaccountDaoInst.GetBalance(fromVacc)
	if strext.ToInt64(balance) < strext.ToInt64(allAmount) {
		ss_log.Info("转账余额不足,balance:%v,allAmount:%v", balance, allAmount)
		reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 正式下单
	log := new(dao.TransferDao)
	log.FromVacc = fromVacc
	log.ToVacc = toVacc
	log.Amount = req.Amount
	log.ExchangeType = req.ExchangeType
	log.Fees = fees
	log.MoneyType = req.MoneyType
	log.FeeRate = rate
	log.RealAmount = req.Amount
	log.Lat = req.Lat
	log.Lng = req.Lng
	log.Ip = req.Ip
	logNo, err := dao.TransferDaoInst.InsertTransfer(tx, log) // ExchangeType 0-扫码1-支付
	if err != nil {
		ss_log.Error("Transfer 插入订单失败,err: %s", err.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	fromPhone, fromCountryCode, pcErr := dao.AccDaoInstance.GetPhoneCountryCodeFromAccNo(req.FromAccountNo)
	if pcErr != nil {
		ss_log.Error("Transfer GetPhoneCountryCodeFromAccNo 查询转账人的手机号和国家码失败,err: %s", pcErr.Error())
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_log.Info("转账人手机号：%v, 国家码: %v", fromPhone, fromCountryCode)

	var writeoffCode string
	if isActive != constants.AccountActived { // 账户不存在才生成码
		nowTimeStr := ss_time.Now(global.Tz).Format(ss_time.DateTimeSlashFormat)

		endTimeStr, err := dao.WriteoffInst.GetCodeEndTimeStr(nowTimeStr)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		// 生成存款码
		writeoffCode = util.RandomDigitStrOnlyNum(constants.Len_SaveMoneyCode)
		log := new(dao.Writeoff)
		log.Code = writeoffCode
		log.IncomeOrderNo = ""
		log.OutGoOrder = ""
		log.TransferOrderNo = logNo
		log.UseStatus = constants.WriteOffCodeWaitUse
		log.SendAccNo = req.FromAccountNo
		log.SendPhone = fromPhone
		log.RecvPhone = req.ToPhone
		log.CreateTime = nowTimeStr
		log.DurationTime = endTimeStr
		log.RecvAccNo = toAccountNo
		errCode := dao.WriteoffInst.InitWriteoff(tx, log)
		if errCode != ss_err.ERR_SUCCESS {
			if errStr := dao.IncomeOrderDaoInst.UpdateIncomeOrderOrderStatus(tx, logNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
		}
	}

	//风控
	riskReply, _ := i.RiskCtrlHandleInstance.Client.GetRiskCtrlReuslt(context.TODO(), &go_micro_srv_riskctrl.GetRiskCtrlResultRequest{
		ApiType: "transfer",
		// 发起支付的账号
		PayerAccNo: req.FromAccountNo,
		ActionTime: time.Now().String(),
		Amount:     req.Amount,
		Ip:         req.Ip,
		PayType:    constants.Risk_Pay_Type_Transfer,
		// 收款人账号
		PayeeAccNo:  toAccountNo,
		ProductType: "transfer",
		// 币种
		MoneyType: req.MoneyType,
		// 订单号
		OrderNo: logNo,
	})

	ss_log.Info("转账风控返回结果,操作结果是---->%s,RiskNo为----->%s", riskReply.OpResult, riskReply.RiskNo)

	if riskReply.OpResult == constants.Risk_Result_No_Pass_Str {
		reply.ResultCode = ss_err.ERR_RISK_IS_RISK
		reply.RiskNo = riskReply.RiskNo
		return nil
	}

	errCode := dao.VaccountDaoInst.AccFromAToBUpperZero(tx, fromVacc, toVacc, req.Amount, logNo, constants.VaReason_TRANSFER)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", errCode)
		if errStr := dao.TransferDaoInst.UpdateTransferOrderStatus(tx, logNo, constants.OrderStatus_Err); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		reply.ResultCode = errCode
		return nil
	}

	if fees != "0" && fees != "" {
		// 修改手续费
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, fromVacc, fees, "-", logNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
	}

	// 下单成功
	if errStr := dao.TransferDaoInst.UpdateTransferOrderStatus(tx, logNo, constants.OrderStatus_Paid); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}

	//添加转账到的账号推送消息
	errAddMessages2 := dao.LogAppMessagesDaoInst.AddLogAppMessagesTx(tx, logNo, constants.LOG_APP_MESSAGES_ORDER_TYPE_TRANSFER_Apply, constants.VaReason_TRANSFER, toAccountNo, constants.OrderStatus_Paid)
	if errAddMessages2 != ss_err.ERR_SUCCESS {
		ss_log.Error("errAddMessages2=[%v]", errAddMessages2)
	}

	appLang, _ := dao.AccDaoInstance.QueryAccountLang(toAccountNo)
	if appLang == "" {
		req.Lang = constants.LangEnUS
	} else {
		req.Lang = appLang
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if err := dao.LogAppFingerprintPayDaoInstance.AddTx(tx, &dao.LogAppFingerprintPayData{
			AccountNo:    req.FromAccountNo,
			DeviceUuid:   req.DeviceUuid,
			SignKey:      req.SignKey,
			OrderNo:      logNo,
			OrderType:    constants.VaReason_TRANSFER,
			Amount:       req.Amount,
			CurrencyType: req.MoneyType,
		}); err != nil {
			ss_log.Error("指纹无密码支付插入日志失败,err[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	ss_log.Info("用户 %s 当前的语言为--->%s", toAccountNo, req.Lang)
	ss_sql.Commit(tx)

	if isActive == constants.AccountActived {
		toAccountType := dao.AccDaoInstance.GetAccountTypeFromAccNo(toAccountNo)
		moneyType := dao.LangDaoInstance.GetLangTextByKey(req.MoneyType, req.Lang)
		timeString := time.Now().Format("2006-01-02 15:04:05")
		// 修正各币种的金额
		amount := common.NormalAmountByMoneyType(req.MoneyType, req.Amount)

		args := []string{
			timeString, amount, moneyType,
		}
		lang, _ := dao.AccDaoInstance.QueryAccountLang(toAccountNo)
		if lang == "" || lang == constants.LangEnUS {
			args = []string{
				amount, moneyType, timeString,
			}
		}

		// 消息推送
		ev := &go_micro_srv_push.PushReqest{
			Accounts: []*go_micro_srv_push.PushAccout{
				{
					AccountNo:   toAccountNo,
					AccountType: toAccountType,
				},
			},
			TempNo: constants.Template_TransferSuccess,
			Args:   args,
		}

		ss_log.Info("publishing %+v\n", ev)
		// publish an event
		if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
			ss_log.Error("error publishing: %v", err)
		}
	}

	if fees != "0" && fees != "" {
		// 发送手续费进MQ
		feeEv := &go_micro_srv_settle.SettleTransferRequest{
			BillNo:    logNo,
			FeesType:  constants.FEES_TYPE_TRANSFER,
			Fees:      fees,
			MoneyType: req.MoneyType,
		}
		ss_log.Info("publishing %+v\n", feeEv)
		// publish an event
		if err := common.SettleEvent.Publish(context.TODO(), feeEv); err != nil {
			ss_log.Error("err=[转账接口,手续费推送到MQ失败,err----->%s]", err.Error())
		}
	}

	reply.OrderNo = logNo
	reply.ResultCode = errCode
	reply.RiskNo = riskReply.RiskNo
	return nil
}

func (b BillHandler) GenRecvCode(ctx context.Context, req *go_micro_srv_bill.GenRecvCodeRequest, reply *go_micro_srv_bill.GenRecvCodeReply) error {
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
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

	code := dao.GenCodeDaoInst.GenCode(req.AccountNo, req.Amount, req.MoneyType, constants.CodeType_Recv)
	if code == "" {
		ss_log.Error("err=[%v]", code)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Code = "B." + code
	return nil
}

// 生成收款码
func (b BillHandler) ScanRecvCode(ctx context.Context, req *go_micro_srv_bill.ScanRecvCodeRequest, reply *go_micro_srv_bill.ScanRecvCodeReply) error {
	// 对二维码进行切割
	split := strings.Split(req.Code, ".")
	// 判断码的类型
	if strings.HasPrefix(req.Code, "A.") {
		// 获取用户信息
		uid, phone, _ := dao.AccDaoInstance.QeuryNamePhoneFromGenKey(split[1])
		data := &go_micro_srv_bill.RecvCodeData{
			AccountNo: uid,
			RecvPhone: phone,
		}
		reply.ResultCode = ss_err.ERR_SUCCESS
		reply.Data = data
		return nil
	}

	accountNo, amount, moneyType, createTime, useStatus := dao.GenCodeDaoInst.GetRecvCode(split[1], constants.CodeType_Recv)

	if createTime == "" {
		reply.ResultCode = ss_err.ERR_PAY_NO_QRCODE
		return nil
	}
	if useStatus == "0" {
		reply.ResultCode = ss_err.ERR_PAY_CANNOT_PAY_CODE
		return nil
	}
	// todo 下面代码测试时候注释,后续要释放
	if ss_time.ParseTimeFromPostgres(createTime, global.Tz).Add(5 * time.Minute).Before(ss_time.Now(global.Tz)) {
		reply.ResultCode = ss_err.ERR_PAY_TIMEOUT
		return nil
	}

	_, usdRecvRate, _ := cache.ApiDaoInstance.GetGlobalParam("usd_recv_rate")
	_, khrRecvRate, _ := cache.ApiDaoInstance.GetGlobalParam("khr_recv_rate")

	rate := "0"
	switch moneyType {
	case "khr":
		rate = khrRecvRate
	case "usd":
		rate = usdRecvRate
	}

	phone := dao.AccDaoInstance.GetPhoneFromAccNo(accountNo)
	if phone == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS

	data := &go_micro_srv_bill.RecvCodeData{
		MoneyType: moneyType,
		AccountNo: accountNo,
		RecvPhone: phone,
		FeeRate:   rate,
	}
	// amount 不返回
	if strings.HasPrefix(req.Code, "A.") {
		reply.Data = data

	} else {
		data.Amount = amount
		reply.Data = data
	}

	return nil
}

// 用户转账到总部
func (b *BillHandler) CustTransferToHeadquarters(ctx context.Context, req *go_micro_srv_bill.CustTransferToHeadquartersRequest, reply *go_micro_srv_bill.CustTransferToHeadquartersReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// todo 判断该用户是否有交易的权限,cust表中的 payment_status 字段
	tradingAuthority, _, _, _, _ := dao.CustDaoInstance.QueryRateRoleFrom(req.OpAccNo)
	if strext.ToInt(tradingAuthority) == constants.TradingAuthorityForbid {
		ss_log.Error("CustTransferToHeadquarters 输错密码超过限制,禁止交易,accountNo: %s,custNo: %s", req.AccountUid, req.OpAccNo)
		reply.ResultCode = ss_err.ERR_Payment_Pwd_Count_Limit
		return nil
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if !dao.AppFingerprintDaoInstance.CheckSignKey(req.AccountUid, req.DeviceUuid, req.SignKey) {
			ss_log.Error("查询不到指纹支付标识，或者指纹支付标识状态为无效, AccountUid[%v], DeviceUuid[%v], SignKey[%v]", req.AccountUid, req.DeviceUuid, req.SignKey)
			reply.ResultCode = ss_err.ERR_AppFingerprint_FAILD
			return nil
		}

	} else {
		if req.Password == "" || req.NonStr == "" {
			ss_log.Error("CustTransferToHeadquarters 验证支付密码的必要参数为空 Password: %s,  NonStr: %s", req.Password, req.NonStr)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		// 判断支付密码是否被限制,并验证支付密码，如果支付密码正确且不是限制中则清除连续支付密码出错次数
		//验证支付密码
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
	}

	// 判断转账金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[转账到总部, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	// 校验收款人账户
	name, cardNum, balanceType, channelNo := dao.CardHeadquartersDaoInst.QueryNameAndNumFromNo(req.CardNo, constants.AccountType_USER)
	if name == "" || name != req.RecName || cardNum == "" || cardNum != req.RecCarNum || balanceType != req.MoneyType {
		reply.ResultCode = ss_err.ERR_REC_CARD_NUM_FAILD
		return nil
	}
	// 判断通道是否对客户充值开放
	// 存款手续费率  单笔存款最大金额  存款单笔手续费    存款计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	saveRate, saveMaxAmount, saveSingleMinFee, saveChargeType := dao.ChannelCustDaoInst.QueryChannelCustSaveInfoFromNo(channelNo, balanceType)
	if saveChargeType == "" {
		ss_log.Error("根据 channelNo 查询 渠道信息失败,channelNo: %s", channelNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 判断限额
	if saveMaxAmount != "" {
		if strext.ToFloat64(req.Amount) > strext.ToFloat64(saveMaxAmount) {
			ss_log.Error("用户向总部存款,交易金额超出单笔最大金额,交易金额为: %s,单笔存款限制最大金额为: %s", req.Amount, saveMaxAmount)
			reply.ResultCode = ss_err.ERR_LOCAL_RULE_EXCEED_AMOUNT
			return nil
		}
	}

	var fees, arriveAmount string
	switch saveChargeType {
	case constants.Charge_Type_Rate: // 手续费比例收取
		feesDeci := ss_count.CountFees(req.Amount, saveRate, "0")
		// 取整
		fees = ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
		arriveAmount = ss_count.Sub(req.Amount, fees).String()
	case constants.Charge_Type_Count: // 按单笔手续费收取
		fees = saveSingleMinFee
		arriveAmount = ss_count.Sub(req.Amount, fees).String()
	}
	ss_log.Info("用户向总部存款,手续费类型为: %s,存款金额为: %s, 计算手续费率为: %s, 手续费为: %s,实际到账金额为: %s", saveChargeType, req.Amount, saveRate, fees, arriveAmount)
	// 插入日志表
	logNo := dao.LogCustToHeadquartersDaoInst.Insert(tx, req.OpAccNo, req.MoneyType, req.Amount, constants.AuditOrderStatus_Pending,
		constants.COLLECTION_TYPE_BANK_TRANSFER, req.CardNo, strext.ToStringNoPoint(constants.TRANSFER_TYPE_BILL), req.ImageId, arriveAmount,
		fees, req.Lat, req.Lng, req.Ip)
	if logNo == "" {
		ss_log.Error("用户向总部充值,记录数据库失败")
		reply.ResultCode = ss_err.ERR_TRANSFER_TO_HEAD_QUARTERS_FAILD
		return nil
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if err := dao.LogAppFingerprintPayDaoInstance.AddTx(tx, &dao.LogAppFingerprintPayData{
			AccountNo:    req.AccountUid,
			DeviceUuid:   req.DeviceUuid,
			SignKey:      req.SignKey,
			OrderNo:      logNo,
			OrderType:    constants.VaReason_Cust_Save,
			Amount:       req.Amount,
			CurrencyType: req.MoneyType,
		}); err != nil {
			ss_log.Error("指纹无密码支付插入日志失败,err[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.OrderNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 转账到总部
func (b *BillHandler) TransferToHeadquarters(ctx context.Context, req *go_micro_srv_bill.TransferToHeadquartersRequest, reply *go_micro_srv_bill.TransferToHeadquartersReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断转账金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[转账到总部, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
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

	// 校验收款人账户
	name, cardNum, balanceType, _ := dao.CardHeadquartersDaoInst.QueryNameAndNumFromNo(req.CardNo, constants.AccountType_SERVICER)
	if name == "" || name != req.RecName || cardNum == "" || cardNum != req.RecCarNum || balanceType != req.MoneyType {
		reply.ResultCode = ss_err.ERR_REC_CARD_NUM_FAILD
		return nil
	}

	// 获取service_no
	var serviceNo string
	switch req.AccountType {
	case constants.AccountType_POS:
		sNo := dao.CashierDaoInst.GetServicerNoFromOpAccNo(req.OpAccNo)
		if sNo == "" {
			reply.ResultCode = ss_err.ERR_TRANSFER_TO_HEAD_QUARTERS_FAILD
			return nil
		}
		serviceNo = sNo
	case constants.AccountType_SERVICER:
		serviceNo = req.OpAccNo
	}
	// 存记录log_to_head
	logNo := dao.LogToHeadquartersDaoInst.InsertLogToHeadquarters(tx, serviceNo, req.ImageId, req.CardNo, constants.CASH_COLLECTION_TYPE,
		req.Amount, req.MoneyType, constants.TRANSFER_TYPE_BILL)
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_TRANSFER_TO_HEAD_QUARTERS_FAILD
		return nil
	}

	// 用户充值100,数据库存的时候是存-100,跟信用卡的形式一样
	//saveAmount := 0 - strext.ToFloat64(req.Amount)
	//ss_log.Info("已修改符号的充值金额为--->%v", saveAmount)
	//// 调用ps服务商预存
	//quotaReq := &go_micro_srv_quota.ModifyQuotaRequest{
	//	CurrencyType: req.MoneyType,
	//	Amount:       strext.ToStringNoPoint(saveAmount),
	//	AccountNo:    servicerAccNo,
	//	OpType:       constants.QuotaOp_SvrPreSave,
	//	LogNo:        logNo,
	//}
	//quotaRepl := &go_micro_srv_quota.ModifyQuotaReply{}
	//quotaRepl, _ = i.QuotaHandleInstance.Client.ModifyQuota(ctx, quotaReq)
	//
	//if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
	//	ss_log.Error("err=[--------------->%s]", "服务商存款,调用八神的服务失败,操作为服务商预存款")
	//	reply.ResultCode = quotaRepl.ResultCode
	//	return nil
	//}
	//
	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 向总部请款,不需要调用服务商取的操作,在后台管理系统里边才需要调用
func (b *BillHandler) ApplyMoney(ctx context.Context, req *go_micro_srv_bill.ApplyMoneyRequest, reply *go_micro_srv_bill.ApplyMoneyReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 判断转账金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[向总部请款, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	// 判断金额是否是10的倍数
	//if strext.ToInt(req.Amount)%10 != 0 {
	//	reply.ResultCode = ss_err.ERR_APPLY_MONEY_NO_TEN_MUTIL
	//	return nil
	//}
	// 获取service_no
	var serviceNo string
	switch req.AccountType {
	case constants.AccountType_POS:
		sNo := dao.CashierDaoInst.GetServicerNoFromOpAccNo(req.OpAccNo)
		if sNo == "" {
			reply.ResultCode = ss_err.ERR_LOG_TO_SERVICE_FAILD
			return nil
		}
		serviceNo = sNo
	case constants.AccountType_SERVICER:
		serviceNo = req.OpAccNo
	}
	// 判断是哪家银行的,判断服务商卡号是否存在
	cardNo, balanceType := dao.CardDaoInst.QueryNameFromNumAndChennel(req.RecCarNum, req.ChannelName, "2")
	if cardNo == "" || balanceType == "" {
		reply.ResultCode = ss_err.ERR_REC_CARD_NUM_FAILD
		return nil
	}

	if balanceType != req.MoneyType {
		reply.ResultCode = ss_err.ERR_REC_CARD_NUM_FAILD
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

	// 插入一条日志进数据库 log_to_service
	logNo := dao.LogToServiceDaoInstance.InsertLogToService(tx, serviceNo, req.Amount, constants.BANK_COLLECTION_TYPE, cardNo, constants.TRANSFER_TYPE_BILL, req.MoneyType)
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_LOG_TO_SERVICE_FAILD
		return nil
	}

	//查询服务商的accountUid 根据serviceNo
	srvAccNo := dao.ServiceDaoInst.GetAccNoFromSrvNo(serviceNo)
	if srvAccNo == "" {
		ss_log.Error("获取服务商账号uid出错")
		reply.ResultCode = ss_err.ERR_LOG_TO_SERVICE_FAILD
		return nil
	}

	//实时额度冻结+，实时额度+
	quotaReq := &go_micro_srv_quota.ModifyQuotaRequest{
		CurrencyType: req.MoneyType,
		Amount:       req.Amount,
		AccountNo:    srvAccNo,
		OpType:       constants.QuotaOp_SvrPreWithdraw,
		LogNo:        logNo,
	}
	quotaRepl, err2 := i.QuotaHandleInstance.Client.ModifyQuota(ctx, quotaReq)
	if err2 != nil || quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[--------------->%s]", "服务商取款申请,调用八神的服务失败,操作为服务商取款申请")
		reply.ResultCode = quotaRepl.ResultCode
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 客户向总部提现
func (b *BillHandler) CustWithdraw(ctx context.Context, req *go_micro_srv_bill.CustWithdrawRequest, reply *go_micro_srv_bill.CustWithdrawReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// todo 判断该用户是否有交易的权限,cust表中的 payment_status 字段
	tradingAuthority, _, _, _, _ := dao.CustDaoInstance.QueryRateRoleFrom(req.IdenNo)
	if strext.ToInt(tradingAuthority) == constants.TradingAuthorityForbid {
		ss_log.Error("CustWithdraw 输错密码超过限制,禁止交易,accountNo: %s,custNo: %s", req.AccountUid, req.IdenNo)
		reply.ResultCode = ss_err.ERR_Payment_Pwd_Count_Limit
		return nil
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if !dao.AppFingerprintDaoInstance.CheckSignKey(req.AccountUid, req.DeviceUuid, req.SignKey) {
			ss_log.Error("CustWithdraw 查询不到指纹支付标识，或者指纹支付标识状态为无效, AccountUid[%v], DeviceUuid[%v], SignKey[%v]", req.AccountUid, req.DeviceUuid, req.SignKey)
			reply.ResultCode = ss_err.ERR_AppFingerprint_FAILD
			return nil
		}

	} else {
		if req.Password == "" || req.NonStr == "" {
			ss_log.Error("CustWithdraw 验证支付密码的必要参数为空 Password: %s,  NonStr: %s", req.Password, req.NonStr)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		// 判断支付密码是否被限制,并验证支付密码，如果支付密码正确且不是限制中则清除连续支付密码出错次数
		replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &go_micro_srv_auth.CheckPayPWDRequest{
			AccountUid:  req.AccountUid,
			AccountType: req.AccountType,
			Password:    req.Password,
			NonStr:      req.NonStr,
			IdenNo:      req.IdenNo,
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

	}

	// 判断转账金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[向总部请款, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	// 判断用户的卡是否存在
	_, _, balanceType, channelNo := dao.CardDaoInst.QueryNameAndNumFromNo(req.RecCarNo)
	if balanceType != req.MoneyType {
		ss_log.Error("银行卡的币种和取款的币种不一致,数据库的币种为: %s,取款的币种为: %s", balanceType, req.MoneyType)
		reply.ResultCode = ss_err.ERR_REC_CARD_NUM_FAILD
		return nil
	}
	// 获取渠道手续费率
	// 提现手续费率      单笔提现最大金额      提现单笔手续费      提现计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	withdrawRate, withdrawMaxAmount, withdrawSingleMinFee, withdrawChargeType := dao.ChannelCustDaoInst.QueryChannelCustWithdrawInfoFromNo(channelNo, balanceType)
	if withdrawChargeType == "" {
		ss_log.Error("根据 channelNo 查询 渠道信息失败,channelNo: %s", channelNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 判断限额
	if withdrawMaxAmount != "" {
		if strext.ToFloat64(req.Amount) > strext.ToFloat64(withdrawMaxAmount) {
			ss_log.Error("用户向总部提现,交易金额超出单笔最大金额,交易金额为: %s,单笔提现最大金额为: %s", req.Amount, withdrawMaxAmount)
			reply.ResultCode = ss_err.ERR_LOCAL_RULE_EXCEED_AMOUNT
			return nil
		}
	}
	var fees string
	switch withdrawChargeType {
	case constants.Charge_Type_Rate: // 手续费比例收取
		feesDeci := ss_count.CountFees(req.Amount, withdrawRate, "0")
		// 取整
		fees = ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
	case constants.Charge_Type_Count: // 按单笔手续费收取
		fees = withdrawSingleMinFee
	}
	// 根据币种获取虚账类型
	var vaType int32
	switch req.MoneyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_DEBIT
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_DEBIT
	default:
		ss_log.Error("用户向总部提现,币种错误,MoneyType: %s", req.MoneyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 确保虚拟账号存在
	recvVaccNo := InternalCallHandlerInst.ConfirmExistVAccount(req.AccountUid, req.MoneyType, strext.ToInt32(vaType))
	// 判断是普通提现还是全部提现,获取金额,计算手续费
	var withdrawAmount string
	switch req.WithdrawType {
	case constants.WITHDRAWAL_TYPE_ORDINARY: // 普通提现
		withdrawAmount = req.Amount
	case constants.WITHDRAWAL_TYPE_ALL: // 全部提现

		balance, _ := dao.VaccountDaoInst.GetBalance(recvVaccNo)
		if req.Amount != balance {
			ss_log.Error("全部提现,提交的金额和用户账户余额对应不上,提交的金额为: %s,账户余额为: %s", req.Amount, balance)
			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_FAILD
			return nil
		}
		withdrawAmountDeci := ss_count.Sub(balance, fees)
		f, _ := withdrawAmountDeci.Float64()
		if f <= 0 {
			ss_log.Error("err=[全部提现手续费不够扣,当前余额为----->%s,应扣手续费为---->%s]", balance, fees)
			reply.ResultCode = ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
			return nil
		}
		withdrawAmount = withdrawAmountDeci.String()
	default:
		ss_log.Error("用户向总部提现,提现类型错误,WithdrawType: %s", req.WithdrawType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_log.Info("用户向总部提现,手续费类型为: %s,提现金额为: %s, 计算手续费率为: %s, 手续费为: %s", withdrawChargeType, req.Amount, withdrawRate, fees)

	logNo := dao.LogToCustDaoInst.Insert(tx, req.MoneyType, req.IdenNo, constants.COLLECTION_TYPE_BANK_TRANSFER,
		req.RecCarNo, withdrawAmount, constants.Order_Type_In_Come, constants.AuditOrderStatus_Pending, req.Lat, req.Lng, fees, req.Ip)
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_LOG_TO_SERVICE_FAILD
		return nil
	}

	// 判断输入的金额是否超额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", recvVaccNo, withdrawAmount, logNo, constants.VaReason_Cust_Withdraw, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
		reply.ResultCode = errStr
		return nil
	}
	if fees != "0" && fees != "" {
		// 修改手续费
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", recvVaccNo, fees, logNo, constants.VaReason_FEES, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
	}
	//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, req.AccountUid); errStr != ss_err.ERR_SUCCESS {
	//	reply.ResultCode = errStr
	//	return nil
	//}

	if req.SignKey != "" { //如果是指纹无密码支付
		if err := dao.LogAppFingerprintPayDaoInstance.AddTx(tx, &dao.LogAppFingerprintPayData{
			AccountNo:    req.AccountUid,
			DeviceUuid:   req.DeviceUuid,
			SignKey:      req.SignKey,
			OrderNo:      logNo,
			OrderType:    constants.VaReason_Cust_Withdraw,
			Amount:       req.Amount,
			CurrencyType: req.MoneyType,
		}); err != nil {
			ss_log.Error("指纹无密码支付插入日志失败,err[%v]", err)
			reply.ResultCode = ss_err.ERR_SYS_DB_OP
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.OrderNo = logNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// POS获取转账至总部的转账记录
func (b *BillHandler) GetTransferToHeadquartersLog(ctx context.Context, req *go_micro_srv_bill.GetTransferToHeadquartersLogRequest, reply *go_micro_srv_bill.GetTransferToHeadquartersLogReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//判断登录的是什么账号类型
	servicerNo := ""
	switch req.AccountType {
	case constants.AccountType_POS:
		servicerNo = dao.ServiceDaoInst.GetServicerNoByCashierNo(req.AccountNo)
	case constants.AccountType_SERVICER:
		servicerNo = dao.ServiceDaoInst.GetSerNoBySerAcc(req.AccountNo)
	default:
		ss_log.Error("AccountType异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if servicerNo == "" {
		ss_log.Error("账号[%v]无对应服务商", req.AccountNo)
		reply.ResultCode = ss_err.ERR_ACC_NO_SER_FAILD
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "to_char(lth.create_time,'yyyy-MM-dd')", Val: req.StartTime, EqType: ">="},
		{Key: "to_char(lth.create_time,'yyyy-MM-dd')", Val: req.EndTime, EqType: "<="},
		{Key: "lth.order_type", Val: "1", EqType: "="},
		{Key: "lth.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "lth.servicer_no", Val: servicerNo, EqType: "="},
		{Key: "lth.currency_type", Val: req.CurrencyType, EqType: "="},
	})
	//统计
	total := dao.ServiceDaoInst.GetCnt(dbHandler, "log_to_headquarters lth", whereModel.WhereStr, whereModel.Args)

	//添加limit
	whereModelM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModelM, `order by lth.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModelM, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	datas, err := dao.ServiceDaoInst.GetTransferToHeadquartersLog(dbHandler, whereModelM.WhereStr, whereModelM.Args)

	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

func (b *BillHandler) CustQeuryFees(ctx context.Context, req *go_micro_srv_bill.CustQeuryRateRequest, reply *go_micro_srv_bill.CustQeuryRateReply) error {
	if req.OpType != constants.OpAccType_Count_Fee_Save && req.OpType != constants.OpAccType_Count_Fee_Withdraw {
		ss_log.Error("计算手续费操作类型有误,req.opType: %v", req.OpType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	saveRate, withdrawRate, withdrawMaxAmount, saveSingleMinFee, withdrawSingleMinFee, saveChargeType, withdrawChargeType,
		supportType, saveMaxAmount := dao.ChannelCustDaoInst.QueryCountFeeInfoFromNo(req.ChannelNo, req.MoneyType)

	if supportType == "" {
		ss_log.Error("获取计算费率的数据信息失败,req.channelNo: %s", req.ChannelNo)
		reply.ResultCode = ss_err.ERR_PAY_QUERY_FEE_FAILD
		return nil
	}
	var rate, fees, chargeType string
	// 判断是存款还是取款
	switch req.OpType {
	case constants.OpAccType_Count_Fee_Save: // 存款
		// 判断该渠道是否有存取款权限
		if supportType != constants.Channel_Support_Type_Save && supportType != constants.Channel_Support_Type_Save_Withdraw {
			reply.ResultCode = ss_err.ERR_PAY_Channel_No_Support_Save
			return nil
		}
		// 判断最大金额限制
		if strext.ToFloat64(req.Amount) > strext.ToFloat64(saveMaxAmount) {
			ss_log.Error("存款,操作金额为: %s,限制的金额为: %s", req.Amount, saveMaxAmount)
			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
			return nil
		}
		// 根据计算手续费类型来计算手续费
		switch saveChargeType {
		case constants.Fee_Charge_Type_Rate: // 按比例
			feeDeci := ss_count.CountFees(req.Amount, saveRate, "0")
			//取整
			fees = ss_big.SsBigInst.ToRound(feeDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
			rate = saveRate
			//chargeType = saveChargeType
		case constants.Fee_Charge_Type_Count: // 单笔
			fees = saveSingleMinFee
			rate = ""
		default:
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		chargeType = saveChargeType
	case constants.OpAccType_Count_Fee_Withdraw: // 取款
		// 判断该渠道是否有存取款权限
		if supportType != constants.Channel_Support_Type_Withdraw && supportType != constants.Channel_Support_Type_Save_Withdraw {
			reply.ResultCode = ss_err.ERR_PAY_Channel_No_Support_Withdraw
			return nil
		}
		// 判断最大金额限制
		if strext.ToFloat64(req.Amount) > strext.ToFloat64(withdrawMaxAmount) {
			ss_log.Error("取款,操作金额为: %s,限制的金额为: %s", req.Amount, withdrawMaxAmount)
			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
			return nil
		}

		// 根据计算手续费类型来计算手续费
		switch withdrawChargeType {
		case constants.Fee_Charge_Type_Rate: // 按比例
			feeDeci := ss_count.CountFees(req.Amount, withdrawRate, "0")
			fees = ss_big.SsBigInst.ToRound(feeDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
			rate = withdrawRate
			//chargeType = withdrawChargeType
		case constants.Fee_Charge_Type_Count: // 单笔
			fees = withdrawSingleMinFee
			rate = ""
		default:
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		chargeType = withdrawChargeType
	}
	reply.Data = &go_micro_srv_bill.CustRateData{
		Rate:       rate,
		Fees:       fees,
		ChargeType: chargeType,
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (b *BillHandler) QueryMinMaxAmount(ctx context.Context, req *go_micro_srv_bill.QueryMinMaxAmountRequest, reply *go_micro_srv_bill.QueryMinMaxAmountReply) error {
	if req.MoneyType == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	var minAmount, maxAmount string
	switch req.Type {
	case constants.TRANSACTION_TYPE_MOBILE_PHONE_WITHDRAW:
		switch req.MoneyType {
		case constants.CURRENCY_USD:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_phone_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_phone_single_min")

		case constants.CURRENCY_KHR:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_phone_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_phone_single_min")
		}
	case constants.TRANSACTION_TYPE_SAVE_MONEY:
		switch req.MoneyType {
		case constants.CURRENCY_USD:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_deposit_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_deposit_single_min")

		case constants.CURRENCY_KHR:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_deposit_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_deposit_single_min")
		}
	case constants.TRANSACTION_TYPE_TRANSFER:
		switch req.MoneyType {
		case constants.CURRENCY_USD:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_transfer_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_transfer_single_min")

		case constants.CURRENCY_KHR:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_transfer_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_transfer_single_min")
		}
	case constants.TRANSACTION_TYPE_SWEEP_WITHDRAW:
		switch req.MoneyType {
		case constants.CURRENCY_USD:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_face_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("usd_face_single_min")

		case constants.CURRENCY_KHR:
			// 最大限额
			maxAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_face_single_max")
			// 最小限额
			minAmount = dao.GlobalParamDaoInstance.QeuryParamValue("khr_face_single_min")
		}

	default:
		ss_log.Error("需要查询最大最小金额的类型失败,请求的类型为--->%s", req.Type)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	reply.Data = &go_micro_srv_bill.MinMaxAmountData{
		MaxAmount: maxAmount,
		MinAmount: minAmount,
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// POS获取转账至服务商的转账记录
func (b *BillHandler) GetTransferToServicerLogs(ctx context.Context, req *go_micro_srv_bill.GetTransferToServicerLogsRequest, reply *go_micro_srv_bill.GetTransferToServicerLogsReply) error {
	//判断登录的是什么账号类型
	servicerNo := ""
	switch req.AccountType {
	case constants.AccountType_POS:
		servicerNo = dao.ServiceDaoInst.GetServicerNoByCashierNo(req.AccountNo)
	case constants.AccountType_SERVICER:
		servicerNo = dao.ServiceDaoInst.GetSerNoBySerAcc(req.AccountNo)
	default:
		ss_log.Error("AccountType异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if servicerNo == "" {
		ss_log.Error("账号[%v]无对应服务商", req.AccountNo)
		reply.ResultCode = ss_err.ERR_ACC_NO_SER_FAILD
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "to_char(lts.create_time,'yyyy-MM-dd')", Val: req.StartTime, EqType: ">="},
		{Key: "to_char(lts.create_time,'yyyy-MM-dd')", Val: req.EndTime, EqType: "<="},
		{Key: "lts.order_type", Val: "1", EqType: "="},
		{Key: "lts.order_status", Val: req.OrderStatus, EqType: "="},
		{Key: "lts.servicer_no", Val: servicerNo, EqType: "="},
		{Key: "lts.currency_type", Val: req.CurrencyType, EqType: "="},
	})

	//统计
	total := dao.ServiceDaoInst.GetCnt(dbHandler, "log_to_servicer lts", whereModel.WhereStr, whereModel.Args)

	whereModelM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModelM, `order by lts.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModelM, strext.ToInt(req.PageSize), strext.ToInt(req.Page))
	datas, err2 := dao.ServiceDaoInst.GetTransferToServicerLogs(dbHandler, whereModelM.WhereStr, whereModelM.Args)
	if err2 != ss_err.ERR_SUCCESS {
		ss_log.Error("err2=[%v]", err2)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

// POS获取对账单
func (b *BillHandler) GetServicerCheckList(ctx context.Context, req *go_micro_srv_bill.GetServicerCheckListRequest, reply *go_micro_srv_bill.GetServicerCheckListReply) error {
	// 检查日期格式是否正确
	if !ss_time.CheckDateIsRight(req.StartTime) || !ss_time.CheckDateIsRight(req.EndTime) {
		ss_log.Error("参数错误:日期格式错误,StartTime:%s,EndTime:%s", req.StartTime, req.EndTime)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 检查结束时间是否大于等于开始时间
	if cmp, _ := ss_time.CompareDate(ss_time.DateFormat, req.StartTime, req.EndTime); cmp > 0 {
		ss_log.Error("参数错误:开始时间大于结束时间,StartTime:%s,EndTime:%s", req.StartTime, req.EndTime)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Page < 1 || req.PageSize < 1 {
		ss_log.Error("参数错误:Page:%d,PageSize:%d", req.Page, req.PageSize)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//判断登录的是什么账号类型
	servicerNo := ""
	switch req.AccountType {
	case constants.AccountType_POS:
		servicerNo = dao.ServiceDaoInst.GetServicerNoByCashierNo(req.AccountNo)
	case constants.AccountType_SERVICER:
		servicerNo = dao.ServiceDaoInst.GetSerNoBySerAcc(req.AccountNo)
	default:
		ss_log.Error("AccountType异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if servicerNo == "" {
		ss_log.Error("账号[%v]无对应服务商", req.AccountNo)
		reply.ResultCode = ss_err.ERR_ACC_NO_SER_FAILD
		return nil
	}

	datas, total, err := dao.ServiceDaoInst.GetServicerCheckList(req.StartTime, req.EndTime, servicerNo, req.Page, req.PageSize)
	if err != nil {
		ss_log.Error("GetServicerCheckList = [%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = total
	return nil
}

// POS获取佣金统计（服务商利润，实际所得）
func (b *BillHandler) GetServicerProfitLedgers(ctx context.Context, req *go_micro_srv_bill.GetServicerProfitLedgersRequest, reply *go_micro_srv_bill.GetServicerProfitLedgersReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//判断登录的是什么账号类型
	servicerNo := ""
	switch req.AccountType {
	case constants.AccountType_POS:
		servicerNo = dao.ServiceDaoInst.GetServicerNoByCashierNo(req.AccountNo)
	case constants.AccountType_SERVICER:
		servicerNo = dao.ServiceDaoInst.GetSerNoBySerAcc(req.AccountNo)
	default:
		ss_log.Error("AccountType异常")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if servicerNo == "" {
		ss_log.Error("账号[%v]无对应服务商", req.AccountNo)
		reply.ResultCode = ss_err.ERR_ACC_NO_SER_FAILD
		return nil
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "to_char(spl.payment_time,'yyyy-MM-dd')", Val: req.StartTime, EqType: ">="},
		{Key: "to_char(spl.payment_time,'yyyy-MM-dd')", Val: req.EndTime, EqType: "<="},
		{Key: "spl.servicer_no", Val: servicerNo, EqType: "="},
	})

	usdCountSumM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(usdCountSumM, "spl.currency_type", "usd", "=")
	usdCountSum := dao.ServiceDaoInst.GetSumAmount(dbHandler, usdCountSumM.WhereStr, usdCountSumM.Args)

	khrCountSumM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhere(khrCountSumM, "spl.currency_type", "khr", "=")
	khrCountSum := dao.ServiceDaoInst.GetSumAmount(dbHandler, khrCountSumM.WhereStr, khrCountSumM.Args)

	ss_sql.SsSqlFactoryInst.AppendWhere(whereModel, "spl.currency_type", req.CurrencyType, "=")
	//获取计数
	total := dao.ServiceDaoInst.GetCntServicerProfitLedgers(dbHandler, whereModel.WhereStr, whereModel.Args)

	//获取服务商利润拥有数据的时间(去重后的年月日)
	whereModelM := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	whereModelB := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModelM, `order by spl.payment_time desc`)
	// 分页
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModelM, strext.ToInt(req.PageSize), strext.ToInt(req.Page))

	datas, err := dao.ServiceDaoInst.GetServicerProfitLedgers(dbHandler, whereModelM.WhereStr, whereModelM.Args, whereModelB)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[%v]", err)

	}

	reply.UsdCountSum = usdCountSum
	reply.KhrCountSum = khrCountSum

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Datas = datas
	reply.Total = strext.ToInt32(total)
	return nil
}

// POS获取对账单明细
func (b *BillHandler) GetServicerProfitLedgerDetail(ctx context.Context, req *go_micro_srv_bill.GetServicerProfitLedgerDetailRequest, reply *go_micro_srv_bill.GetServicerProfitLedgerDetailReply) error {
	if req.LogNo == "" {
		ss_log.Error("%s", "请求参数logNo为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	data, err := dao.ServiceDaoInst.GetServicerProfitLedgerDetail(req.LogNo)

	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("GetServicerChecksDetails | err= [%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (b *BillHandler) CustIncomeBillsDetail(ctx context.Context, req *go_micro_srv_bill.CustIncomeBillsDetailRequest, reply *go_micro_srv_bill.CustIncomeBillsDetailReply) error {
	ss_log.Info("TradingHandler | CustIncomeBillsDetail req==[%v]", req)
	data, err := dao.IncomeOrderDaoInst.CustIncomeBillsDetail(req.LogNo)

	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("CustIncomeBillsDetail | err= [%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (b *BillHandler) CustOutgoBillsDetail(ctx context.Context, req *go_micro_srv_bill.CustOutgoBillsDetailRequest, reply *go_micro_srv_bill.CustOutgoBillsDetailReply) error {
	ss_log.Info("TradingHandler | CustOutgoBillsDetail req==[%v]", req)

	data, err := dao.OutgoOrderDaoInst.CustOutgoBillsDetail(req.LogNo)

	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("CustIncomeBillsDetail | err= [%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (b *BillHandler) CustTransferBillsDetail(ctx context.Context, req *go_micro_srv_bill.CustTransferBillsDetailRequest, reply *go_micro_srv_bill.CustTransferBillsDetailReply) error {
	ss_log.Info("TradingHandler | CustOutgoBillsDetail req==[%v]", req)

	data, err := dao.TransferDaoInst.CustTransferBillsDetail(req.LogNo)

	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("CustTransferBillsDetail | err= [%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (b *BillHandler) CustCollectionBillsDetail(ctx context.Context, req *go_micro_srv_bill.CustCollectionBillsDetailRequest, reply *go_micro_srv_bill.CustCollectionBillsDetailReply) error {
	data, err := dao.CollectionDaoInst.CustCollectionBillsDetail(req.LogNo)

	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("CustCollectionBillsDetail | err= [%v]", err)
		reply.ResultCode = ss_err.ERR_SYS_DB_GET
		return nil
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

//用户查询自己的消息详细消息
func (b *BillHandler) CustOrderBillDetail(ctx context.Context, req *go_micro_srv_bill.CustOrderBillDetailRequest, reply *go_micro_srv_bill.CustOrderBillDetailReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	data := new(go_micro_srv_bill.CustOrderBillDetailData)
	if req.LogNo != "" && req.OrderType == constants.VaReason_PlatformFreeze {
		/**
		查询因操作核销码产生的虚账金额变动日志
		*/
		if req.LogNo == "" {
			ss_log.Error("LogNo 为空")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		log, err := dao.LogVaccountDaoInst.GetLogVAccountJoinWriteOff(req.AccountNo, req.LogNo, constants.VaReason_PlatformFreeze)
		if err != nil {
			ss_log.Error("GetLogVAccountByLogNo() err=[%v]", err)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		data.CreateTime = log.CreateTime
		data.OpType = log.OpType
		data.Amount = log.Amount
		data.LogNo = log.OrderNo
		data.OrderType = constants.VaReason_PlatformFreeze
		data.BalanceType = log.CurrencyType

	} else {
		if req.OrderNo == "" {
			ss_log.Error("orderNo 为空")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//根据订单号、和订单类型查询内容
		var errCode string
		data, errCode = dao.LogAppMessagesDaoInst.CustOrderBillDetail(req.AccountNo, req.OrderNo, req.OrderType)
		if errCode != ss_err.ERR_SUCCESS {
			ss_log.Error("GetLogAppMessages errCode=[%v]", errCode)
			reply.ResultCode = errCode
			return nil
		}

		//查询核销码
		switch req.OrderType {
		case constants.VaReason_INCOME: // 存款、充值
			fallthrough
		case constants.VaReason_TRANSFER: // 转账
			code, errGetCode := dao.WriteoffInst.GetCode(dbHandler, req.OrderNo, data.OrderType)
			if errGetCode != ss_err.ERR_SUCCESS {
				ss_log.Error("获取核销码Code失败,err=[%v]", errGetCode)
			}
			data.Code = code
		default:

		}
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = data
	return nil
}

func (*BillHandler) GetLogAppMessagesCnt(ctx context.Context, req *go_micro_srv_bill.GetLogAppMessagesCntRequest, reply *go_micro_srv_bill.GetLogAppMessagesCntReply) error {
	//统计未读的消息数量
	total := dao.LogAppMessagesDaoInst.GetNoReadCnt(req.AccountUid)
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total)
	return nil
}

//实时报表
func (*BillHandler) RealTimeCount(ctx context.Context, req *go_micro_srv_bill.RealTimeCountRequest, reply *go_micro_srv_bill.RealTimeCountReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var srvAccNo string
	// 判断是服务商还是收银员
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		srvAccNo = req.AccountUid
	case constants.AccountType_POS:
		//查询店员服务商的uid
		serAccNo := dao.CashierDaoInst.GetSerAccNoByCashierAccNo(dbHandler, req.AccountUid)
		if serAccNo == "" {
			ss_log.Error("查询不到账号对应的服务商账号id")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		//req.AccountUid = serAccNo
		srvAccNo = serAccNo
	}

	//只统计当天的
	nowTime := time.Now()
	nowTimeStr := nowTime.Format("2006-01-02")

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "to_char(bdr.create_time,'yyyy-MM-dd')", Val: nowTimeStr, EqType: "="},
		//{Key: "bdr.account_no", Val: req.AccountUid, EqType: "="},
		{Key: "bdr.account_no", Val: srvAccNo, EqType: "="},
		{Key: "bdr.account_type", Val: constants.AccountType_SERVICER, EqType: "="},
		{Key: "bdr.order_status", Val: constants.OrderStatus_Paid, EqType: "="}, //只统计已完成的
	})

	var datas []*go_micro_srv_bill.RealTimeCountData

	billTypes := []string{
		constants.BILL_TYPE_INCOME,
		constants.BILL_TYPE_OUTGO,
		constants.BILL_TYPE_PROFIT,
		constants.BILL_TYPE_RECHARGE,
		constants.BILL_TYPE_WITHDRAWALS,
	}

	for _, billType := range billTypes {
		data := &go_micro_srv_bill.RealTimeCountData{}
		whereModelB := ss_sql.SsSqlFactoryInst.DeepClone(whereModel)
		ss_sql.SsSqlFactoryInst.AppendWhere(whereModelB, "bdr.bill_type", billType, "=")

		whereModelUsd := ss_sql.SsSqlFactoryInst.DeepClone(whereModelB)
		ss_sql.SsSqlFactoryInst.AppendWhere(whereModelUsd, "bdr.currency_type", "usd", "=")
		usdSum := dao.BillingDetailsResultsDaoInstance.GetSum(dbHandler, whereModelUsd.WhereStr, whereModelUsd.Args)

		whereModelKhr := ss_sql.SsSqlFactoryInst.DeepClone(whereModelB)
		ss_sql.SsSqlFactoryInst.AppendWhere(whereModelKhr, "bdr.currency_type", "khr", "=")
		khrSum := dao.BillingDetailsResultsDaoInstance.GetSum(dbHandler, whereModelKhr.WhereStr, whereModelKhr.Args)

		data.Type = billType
		data.UsdSum = usdSum
		data.KhrSum = khrSum

		datas = append(datas, data)
	}

	//// usd 授权额度
	//usdAuthBalance, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountUid, "usd")
	//// khr 授权额度
	//khrAuthBalance, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountUid, "khr")
	//
	////实时额度和可用额度
	//usdRealBalance, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountUid, "usd_spent")
	//usdNoSpent := strext.ToStringNoPoint(strext.ToInt64(usdAuthBalance) - strext.ToInt64(usdRealBalance))
	//
	//khrRealBalance, _ := dao.ServiceDaoInst.GetSerAuthCollectLimit(dbHandler, req.AccountUid, "khr_spent")
	//khrNoSpent := strext.ToStringNoPoint(strext.ToInt64(khrAuthBalance) - strext.ToInt64(khrRealBalance))
	posReminReq := &go_micro_srv_auth.GetPosRemainRequest{
		// 账号
		AccountUid:  req.AccountUid,
		AccountType: req.AccountType,
	}

	remainReply, _ := i.AuthHandlerInst.Client.GetPosRemain(ctx, posReminReq)
	if remainReply.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("------>%s", "调用获取pos 可用额度rpc失败")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//授权收款额度
	authBalanceData := &go_micro_srv_bill.RealTimeCountData{
		Type:   "auth_balance",
		UsdSum: remainReply.Data.AuthUsd,
		KhrSum: remainReply.Data.AuthKhr,
	}
	datas = append(datas, authBalanceData)
	usdUseSum := ss_count.Add(remainReply.Data.AuthUsd, remainReply.Data.UseUsd)
	khrUseSum := ss_count.Add(remainReply.Data.AuthKhr, remainReply.Data.UseKhr)
	//未使用授权收款额度
	NoSpentData := &go_micro_srv_bill.RealTimeCountData{
		Type:   "no_spent_balance",
		UsdSum: usdUseSum,
		KhrSum: khrUseSum,
	}
	datas = append(datas, NoSpentData)

	reply.Data = datas
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

/**
添加改变余额订单
*/
func (b *BillHandler) AddChangeBalanceOrder(ctx context.Context, req *go_micro_srv_bill.AddChangeBalanceOrderRequest, reply *go_micro_srv_bill.AddChangeBalanceOrderReply) error {
	if req.Amount == "" || req.Amount == "0" {
		ss_log.Error("订单金额为空或为0")
		reply.ResultCode = ss_err.ERR_PAY_AMOUNT_NOT_MIN_ZERO
		return nil
	}

	//验证登录密码
	if !dao.AccDaoInstance.CheckAdminLoginPWD(req.LoginUid, req.LoginPwd, req.NonStr) {
		ss_log.Error("登录密码错误。")
		reply.ResultCode = ss_err.ERR_DB_PWD
		return nil
	}

	beforeBalance := "0" //改变前余额
	afterBalance := "0"  //改变后余额
	quotaType := ""
	switch req.AccountType {
	case constants.AccountType_SERVICER:
		switch req.CurrencyType {
		case constants.CURRENCY_USD:
			_, usdBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.AccUid, constants.VaType_QUOTA_USD_REAL)
			beforeBalance = usdBalance
		case constants.CURRENCY_KHR:
			_, khrBalance := dao.VaccountDaoInst.GetBalanceFromAccNo(req.AccUid, constants.VaType_QUOTA_KHR_REAL)
			beforeBalance = khrBalance
		default:
			ss_log.Error("币种类型[%v]错误", req.CurrencyType)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		switch req.OpType {
		case constants.VaOpType_Add:
			quotaType = constants.QuotaOp_ChangeSvrBalanceAdd
			amountB := ss_count.Sub("0", req.Amount).String()
			afterBalance = ss_count.Add(beforeBalance, amountB)
		case constants.VaOpType_Minus:
			quotaType = constants.QuotaOp_ChangeSvrBalanceMinus
			amountB := ss_count.Sub("0", req.Amount).String()
			afterBalance = ss_count.Sub(beforeBalance, amountB).String()
		default:
			ss_log.Error("操作类型[%v]错误", req.OpType)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	case constants.AccountType_USER:
		//查询当前用户的余额
		khrBalance, usdBalance := dao.AccDaoInstance.GetRemain(req.AccUid)
		switch req.CurrencyType {
		case constants.CURRENCY_USD:
			beforeBalance = usdBalance
		case constants.CURRENCY_KHR:
			beforeBalance = khrBalance
		default:
			ss_log.Error("币种类型[%v]错误", req.CurrencyType)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}

		switch req.OpType {
		case constants.VaOpType_Add:
			quotaType = constants.QuotaOp_ChangeCustBalanceAdd
			afterBalance = ss_count.Add(beforeBalance, req.Amount)
		case constants.VaOpType_Minus:
			quotaType = constants.QuotaOp_ChangeCustBalanceMinus
			afterBalance = ss_count.Sub(beforeBalance, req.Amount).String()
		default:
			ss_log.Error("操作类型[%v]错误", req.OpType)
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	default:
		ss_log.Error("操作类型错误")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}
	defer ss_sql.Rollback(tx)

	//添加日志
	logNo, err := dao.ChangeBalanceOrderDaoInst.AddChangeBalanceOrder(tx, dao.ChangeBalanceOrderDao{
		AccountNo:      req.AccUid,
		CurrencyType:   req.CurrencyType,
		BeforeBalance:  beforeBalance,
		ChangeAmount:   req.Amount,
		AfterBalance:   afterBalance,
		ChangeReason:   req.ChangeReason,
		OrderStatus:    constants.OrderStatus_Paid,
		AccountType:    req.AccountType,
		AdminAccountNo: req.LoginUid,
		OpType:         req.OpType,
	})
	if err != nil {
		ss_log.Error("新增改变余额订单失败，err[%v]", err)
		reply.ResultCode = ss_err.ERR_OPERATE_FAILD
		return nil
	}

	//发起rpc 修改余额
	quotaRepl, quotaErr := i.QuotaHandleInstance.Client.ModifyQuota(context.TODO(), &go_micro_srv_quota.ModifyQuotaRequest{
		CurrencyType: req.CurrencyType,
		Amount:       req.Amount,
		AccountNo:    req.AccUid,
		OpType:       quotaType,
		LogNo:        logNo,
	})
	if quotaErr != nil {
		ss_log.Error("调用远程ModifyQuota rpc失败,err: %v", quotaErr)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}

	if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[--------------->%s]", "调用服务失败,操作为改变余额")
		reply.ResultCode = quotaRepl.ResultCode
		return nil
	}

	switch req.AccountType {
	case constants.AccountType_USER:
		// 同步用户余额
		//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, req.AccUid); errStr != ss_err.ERR_SUCCESS {
		//	reply.ResultCode = errStr
		//	return nil
		//}
	case constants.AccountType_SERVICER:
		idenNo := dao.RelaAccIdenDaoInst.GetIdenNo(req.AccUid, req.AccountType)

		// todo 插入 billing_details_results
		if logNo2 := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, req.Amount, req.CurrencyType, req.AccUid,
			req.AccountType, logNo, "0", constants.OrderStatus_Paid, idenNo, idenNo, constants.BillDetailTypeChangeBalance,
			"0", req.Amount); logNo2 == "" {
			ss_log.Error("插入服务商结果表billing_details_results失败。")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
