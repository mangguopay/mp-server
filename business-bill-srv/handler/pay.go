package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/common"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/business-bill-srv/i"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	authProto "a.a/mp-server/common/proto/auth"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	notifyProto "a.a/mp-server/common/proto/notify"
	pushProto "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
	"strings"
	"time"
)

// 用户扫码-带金额二维码支付
func (b *BusinessBillHandler) QrCodeAmountPay(ctx context.Context, req *businessBillProto.QrCodeAmountPayRequest, reply *businessBillProto.QrCodeAmountPayReply) error {
	if req.QrCodeId == "" {
		ss_log.Info("QrCodeId参数为空")
		reply.ResultCode = ss_err.OrderNoOrQrCodeIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.OrderNoOrQrCodeIsEmpty, req.Lang)
		return nil
	}

	if req.PaymentMethod == constants.PayMethodBankCard {
		if req.BankCardNo == "" {
			ss_log.Error("BankCardNumber参数为空")
			reply.ResultCode = ss_err.BankCardNumberIsEmpty
			reply.Msg = ss_err.GetMsg(ss_err.BankCardNumberIsEmpty, req.Lang)
		}
	}

	// 通过二维码id获取订单编号
	orderNo, err := dao.BusinessBillQrCodeInst.QueryOrderNoByQrCodeId(req.QrCodeId)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("二维码不存在,qrCodeId:%v", req.QrCodeId)
			reply.ResultCode = ss_err.QrCodeNotExist
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotExist, req.Lang)
			return nil
		}
		ss_log.Error("二维码id查询订单号失败, qrCodeId=%v, err=%v", req.QrCodeId, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//调用支付接口
	payRequest := PayRequest{
		AccountNo:       req.AccountNo,
		AccountType:     req.AccountType,
		PaymentPassword: req.PaymentPassword,
		NonStr:          req.NonStr,
		OrderNo:         orderNo,
		PaymentMethod:   req.PaymentMethod,
		BankCardNo:      req.BankCardNo,
		SignKey:         req.SignKey,
		DeviceUuid:      req.DeviceUuid,
	}
	payReply, resultCode := b.pay(ctx, payRequest)
	if resultCode != ss_err.Success {
		ss_log.Error("调用支付接口失败，req=%v, resultCode=%v", strext.ToJson(payRequest), resultCode)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//查询订单信息，为推送消息做准备
	billData, err := dao.BusinessBillDaoInst.GetOrderInfoByOrderNo(orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询订单失败,订单用户不存在,OrderNo:%v", orderNo)
		}
		ss_log.Error("查询订单失败,err:%v,OrderNo:%v", err, orderNo)
	}

	//添加推送消息
	msg := new(dao.LogAppMessagesDao)
	msg.OrderNo = orderNo
	msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Pay
	msg.OrderType = constants.VaReason_Cust_Pay_Order
	msg.AccountNo = req.AccountNo
	msg.OrderStatus = constants.OrderStatus_Paid
	err = dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
	if err != nil {
		ss_log.Error("aAddMessagesErr=[%v]", err)
	}

	ss_log.Info("用户 %s 当前的语言为--->%s", req.AccountNo, req.Lang)
	moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(billData.CurrencyType), req.Lang)
	timeString := time.Now().Format("2006-01-02 15:04:05")
	// 修正各币种的金额
	amount := common.NormalAmountByMoneyType(billData.CurrencyType, billData.Amount)

	args := []string{
		timeString, amount, moneyType,
	}

	if req.Lang == constants.LangEnUS {
		args = []string{
			amount, moneyType, timeString,
		}
	}

	// 消息推送
	ev := &pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				AccountNo:   req.AccountNo,
				AccountType: constants.AccountType_USER,
			},
		},
		TempNo: constants.Template_PaySuccess,
		Args:   args,
	}

	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("消息推送到用户toAccountNo[%v]出错。error : %v", req.AccountNo, err)
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = payReply.OrderNo
	reply.OutOrderNo = payReply.OutOrderNo
	reply.OrderStatus = payReply.OrderStatus
	reply.CreateTime = payReply.CreateTime
	reply.PayTime = payReply.PayTime
	reply.Subject = payReply.Subject
	reply.UserOrderType = constants.VaReason_Cust_Pay_Order
	return nil
}

// 用户扫码-固定二维码支付
func (b *BusinessBillHandler) QrCodeFixedPay(ctx context.Context, req *businessBillProto.QrCodeFixedPayRequest, reply *businessBillProto.QrCodeFixedPayReply) error {
	//请求参数检查
	if resultCode, err := CheckQrCodeFixedPayReq(req); err != nil {
		ss_log.Error("QrCodeFixedPay() 参数错误，err=%v,", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//调用支付接口
	payRequest := PayRequest{
		AccountNo:       req.AccountNo,
		AccountType:     req.AccountType,
		PaymentPassword: req.PaymentPassword,
		NonStr:          req.NonStr,
		OrderNo:         req.OrderNo,
		PaymentMethod:   req.PaymentMethod,
		BankCardNo:      req.BankCardNo,
		SignKey:         req.SignKey,
		DeviceUuid:      req.DeviceUuid,
	}
	payReply, resultCode := b.pay(ctx, payRequest)
	if resultCode != ss_err.Success {
		ss_log.Error("调用支付接口失败，req=%v, resultCode=%v", strext.ToJson(payRequest), resultCode)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	if resultCode == ss_err.Success {
		//用户扫商家固码订单不需要通知, 直接修改通知状态为成功
		err := dao.BusinessBillDaoInst.UpdateOrderNotifyStatus(payReply.OrderNo, constants.NotifyStatusSuccess)
		if err != nil {
			ss_log.Error("修改订单通知状态失败, orderNo=%v, err=%v", payReply.OrderNo, err)
		}

		//查询订单信息，为推送消息给用户做准备
		billData, err := dao.BusinessBillDaoInst.GetOrderInfoByOrderNo(payReply.OrderNo)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("查询订单失败,订单用户不存在,OrderNo:%v", payReply.OrderNo)
			}
			ss_log.Error("查询订单失败,err:%v,OrderNo:%v", err, payReply.OrderNo)
		}

		//添加推送消息
		msg := new(dao.LogAppMessagesDao)
		msg.OrderNo = payReply.OrderNo
		msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Pay
		msg.OrderType = constants.VaReason_Cust_Pay_Order
		msg.AccountNo = req.AccountNo
		msg.OrderStatus = constants.OrderStatus_Paid
		err = dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
		if err != nil {
			ss_log.Error("aAddMessagesErr=[%v]", err)
		}

		ss_log.Info("用户 %s 当前的语言为--->%s", req.AccountNo, req.Lang)
		moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(billData.CurrencyType), req.Lang)
		timeString := time.Now().Format("2006-01-02 15:04:05")
		// 修正各币种的金额
		amount := common.NormalAmountByMoneyType(billData.CurrencyType, billData.Amount)

		args := []string{
			timeString, amount, moneyType,
		}

		if req.Lang == constants.LangEnUS {
			args = []string{
				amount, moneyType, timeString,
			}
		}

		// 消息推送
		ev := &pushProto.PushReqest{
			Accounts: []*pushProto.PushAccout{
				{
					AccountNo:   req.AccountNo,
					AccountType: constants.AccountType_USER,
				},
			},
			TempNo: constants.Template_PaySuccess,
			Args:   args,
		}

		ss_log.Info("publishing %+v\n", ev)
		// publish an event
		if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
			ss_log.Error("消息推送到用户toAccountNo[%v]出错。error : %v", req.AccountNo, err)
		}

	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = payReply.OrderNo
	reply.OrderStatus = payReply.OrderStatus
	reply.OutOrderNo = payReply.OutOrderNo
	reply.CreateTime = payReply.CreateTime
	reply.Subject = payReply.Subject
	reply.PayTime = payReply.PayTime
	reply.UserOrderType = constants.VaReason_Cust_Pay_Order
	return nil
}

// 商家扫用户付款码支付-用户输入支付密码支付
func (b *BusinessBillHandler) OrderPay(ctx context.Context, req *businessBillProto.OrderPayRequest, reply *businessBillProto.OrderPayReply) error {
	//请求参数检查
	if resultCode, err := CheckOrderPayReq(req); err != nil {
		ss_log.Error("QrCodeFixedPay() 参数错误，err=%v,", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//调用支付接口
	payRequest := PayRequest{
		AccountNo:       req.AccountNo,
		AccountType:     req.AccountType,
		PaymentPassword: req.PaymentPwd,
		NonStr:          req.NonStr,
		OrderNo:         req.OrderNo,
		PaymentMethod:   req.PaymentMethod,
		BankCardNo:      req.BankCardNo,
		SignKey:         req.SignKey,    //指纹支付标识
		DeviceUuid:      req.DeviceUuid, //设备号
	}

	payReply, resultCode := b.pay(ctx, payRequest)
	if resultCode != ss_err.Success {
		ss_log.Error("调用支付接口失败，req=%v, resultCode=%v", strext.ToJson(payRequest), resultCode)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//查询订单信息，为推送消息做准备
	billData, err := dao.BusinessBillDaoInst.GetOrderInfoByOrderNo(req.OrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询订单失败,订单用户不存在,OrderNo:%v", req.OrderNo)
		}
		ss_log.Error("查询订单失败,err:%v,OrderNo:%v", err, req.OrderNo)
	}

	//添加推送消息
	msg := new(dao.LogAppMessagesDao)
	msg.OrderNo = req.OrderNo
	msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Pay
	msg.OrderType = constants.VaReason_Cust_Pay_Order
	msg.AccountNo = req.AccountNo
	msg.OrderStatus = constants.OrderStatus_Paid
	err = dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
	if err != nil {
		ss_log.Error("aAddMessagesErr=[%v]", err)
	}

	ss_log.Info("用户 %s 当前的语言为--->%s", req.AccountNo, req.Lang)
	moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(billData.CurrencyType), req.Lang)
	timeString := time.Now().Format("2006-01-02 15:04:05")
	// 修正各币种的金额
	amount := common.NormalAmountByMoneyType(billData.CurrencyType, billData.Amount)

	args := []string{
		timeString, amount, moneyType,
	}

	if req.Lang == constants.LangEnUS {
		args = []string{
			amount, moneyType, timeString,
		}
	}

	// 消息推送
	ev := &pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				AccountNo:   req.AccountNo,
				AccountType: constants.AccountType_USER,
			},
		},
		TempNo: constants.Template_PaySuccess,
		Args:   args,
	}

	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("消息推送到用户toAccountNo[%v]出错。error : %v", req.AccountNo, err)
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = payReply.OrderNo
	reply.OrderStatus = payReply.OrderStatus
	reply.CreateTime = payReply.CreateTime
	reply.Subject = payReply.Subject
	reply.PayTime = payReply.PayTime
	reply.UserOrderType = constants.VaReason_Cust_Pay_Order
	return nil
}

//App支付
func (b *BusinessBillHandler) AppPay(ctx context.Context, req *businessBillProto.AppPayRequest, reply *businessBillProto.AppPayReply) error {
	//请求参数检查
	if resultCode, err := CheckAppPayReq(req); err != nil {
		ss_log.Error("AppPay() 参数错误，err=%v,", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//验证AppPayContent参数是否是平台发出的
	ss_log.Info("AppPayContent=%v", req.AppPayContent)
	contentMap := strext.Json2Map(req.AppPayContent)
	if contentMap == nil {
		ss_log.Error("AppPayContent参数不是一个json字符串")
		reply.ResultCode = ss_err.VerifySignFail
		reply.Msg = ss_err.GetMsg(ss_err.VerifySignFail, req.Lang)
		return nil
	}
	reqSign := contentMap[common.SignField]
	verifySign := AppPayContentMakeSign(contentMap)
	if verifySign != reqSign {
		ss_log.Error("验签失败，sign=[%v], verifySign=[%v]", reqSign, verifySign)
		reply.ResultCode = ss_err.VerifySignFail
		reply.Msg = ss_err.GetMsg(ss_err.VerifySignFail, req.Lang)
		return nil
	}

	//调用支付接口
	payRequest := PayRequest{
		AccountNo:       req.AccountNo,
		AccountType:     req.AccountType,
		PaymentPassword: req.PaymentPwd,
		NonStr:          req.NonStr,
		OrderNo:         strext.ToString(contentMap["order_no"]),
		PaymentMethod:   req.PaymentMethod,
		BankCardNo:      req.BankCardNo,
		SignKey:         req.SignKey,
		DeviceUuid:      req.DeviceUuid,
	}

	payReply, resultCode := b.pay(ctx, payRequest)
	if resultCode != ss_err.Success {
		ss_log.Error("调用支付接口失败，req=%v, resultCode=%v", strext.ToJson(payRequest), resultCode)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//查询订单信息，为推送消息给用户做准备
	billData, err := dao.BusinessBillDaoInst.GetOrderInfoByOrderNo(payReply.OrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询订单失败,订单用户不存在,OrderNo:%v", payReply.OrderNo)
		}
		ss_log.Error("查询订单失败,err:%v,OrderNo:%v", err, payReply.OrderNo)
	}

	//添加推送消息
	msg := new(dao.LogAppMessagesDao)
	msg.OrderNo = payReply.OrderNo
	msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Cust_Pay
	msg.OrderType = constants.VaReason_Cust_Pay_Order
	msg.AccountNo = req.AccountNo
	msg.OrderStatus = constants.OrderStatus_Paid
	err = dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
	if err != nil {
		ss_log.Error("aAddMessagesErr=[%v]", err)
	}

	ss_log.Info("用户 %s 当前的语言为--->%s", req.AccountNo, req.Lang)
	moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(billData.CurrencyType), req.Lang)
	timeString := time.Now().Format("2006-01-02 15:04:05")
	// 修正各币种的金额
	amount := common.NormalAmountByMoneyType(billData.CurrencyType, billData.Amount)

	args := []string{
		timeString, amount, moneyType,
	}

	if req.Lang == constants.LangEnUS {
		args = []string{
			amount, moneyType, timeString,
		}
	}

	// 消息推送
	ev := &pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				AccountNo:   req.AccountNo,
				AccountType: constants.AccountType_USER,
			},
		},
		TempNo: constants.Template_PaySuccess,
		Args:   args,
	}

	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("消息推送到用户toAccountNo[%v]出错。error : %v", req.AccountNo, err)
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = payReply.OrderNo
	reply.OrderStatus = payReply.OrderStatus
	reply.CreateTime = payReply.CreateTime
	reply.Subject = payReply.Subject
	reply.PayTime = payReply.PayTime
	reply.UserOrderType = constants.VaReason_Cust_Pay_Order
	return nil
}

//==============================================================
type PayRequest struct {
	OrderNo         string
	AccountNo       string
	PaymentPassword string
	NonStr          string
	AccountType     string
	PaymentMethod   string
	BankCardNo      string
	SignKey         string //指纹支付标识
	DeviceUuid      string //设备uuid
}

type PayReply struct {
	OrderNo     string
	OutOrderNo  string
	OrderStatus string
	Subject     string
	CreateTime  string
	PayTime     string
}

// 支付接口
func (b *BusinessBillHandler) pay(ctx context.Context, req PayRequest) (*PayReply, string) {
	resultCode, err := CheckPayReq(req)
	if err != nil {
		ss_log.Error("pay() 请求参数错误，err=%v", err)
		return nil, resultCode
	}

	bill, err := dao.BusinessBillDaoInst.GetOrderInfoByOrderNo(req.OrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询订单失败,订单用户不存在,OrderNo:%v", req.OrderNo)
			return nil, ss_err.OrderNotExist
		}
		ss_log.Error("查询订单失败,err:%v,OrderNo:%v", err, req.OrderNo)
		return nil, ss_err.SystemErr
	}

	ss_log.Info("订单信息:%v", strext.ToJson(bill))

	// 只有待支付的订单才能支付
	if bill.OrderStatus != constants.BusinessOrderStatusPending {
		ss_log.Error("订单已支付，不允许再次支付,OrderNo:%v, OrderStatus:%v", req.OrderNo, bill.OrderStatus)
		return nil, ss_err.OrderPaid
	}

	//检查订单是否已过期(过期时间小于当前时间)
	if bill.ExpireTime < ss_time.Now(global.Tz).Unix() {
		//修更新订单为已超时
		updateErr := dao.BusinessBillDaoInst.UpdateOrderOutTimeById(bill.OrderNo)
		if updateErr != nil {
			ss_log.Error("修改订单为已过期失败,OrderNo:%v, err:%v", req.OrderNo, updateErr)
		}
		ss_log.Error("订单已过期,OrderNo:%v", req.OrderNo)
		return nil, ss_err.OrderExpired
	}

	if bill.BusinessAccountNo == req.AccountNo {
		ss_log.Error("不能自己付款给自己，存在刷单行为，account=%v", req.AccountNo)
		return nil, ss_err.PayeeAccountErr
	}

	//用户基本信息
	custInfo, err := dao.CustDaoInst.QueryCustInfo(req.AccountNo, "")
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("查询用户账号信息失败,用户不存在,AccountNo:%v", req.AccountNo)
			return nil, ss_err.AccountNoNotExist
		}
		ss_log.Error("查询用户账号信息失败,err:%v,AccountNo:%v", err, req.AccountNo)
		return nil, ss_err.SystemErr
	}

	//检查用户是否已实名
	if custInfo.IndividualAuthStatus != constants.AuthMaterialStatus_Passed {
		ss_log.Error("用户账号未实名认证, AccountNo:%v, IndividualAuthStatus=%v", req.AccountNo, custInfo.IndividualAuthStatus)
		return nil, ss_err.UserHasNoRealName

	}

	//判断用户是否有交易权限
	if custInfo.TradingAuthority == constants.TradingAuthorityForbid {
		ss_log.Error("用户账号被禁止交易,accountNo: %s", req.AccountNo)
		return nil, ss_err.AccountNoNotTradeForbid
	}

	//判断用户是否有出款权限
	if custInfo.OutgoAuthorization == constants.CustOutTransferAuthorizationDisabled {
		ss_log.Error("用户[%s]没有出款权限", req.AccountNo)
		return nil, ss_err.AccountNoNotTradeForbid
	}

	if req.SignKey != "" { //如果是指纹无密码支付
		if !dao.AppFingerprintDaoInstance.CheckSignKey(req.AccountNo, req.DeviceUuid, req.SignKey) {
			ss_log.Error("查询不到指纹支付标识，或者指纹支付标识状态为无效, AccountNo[%v], DeviceUuid[%v], SignKey[%v]", req.AccountNo, req.DeviceUuid, req.SignKey)
			return nil, ss_err.ERR_AppFingerprint_FAILD
		}

	} else {
		//支付密码校验
		authReq := authProto.CheckPayPWDRequest{
			AccountUid:  req.AccountNo,
			AccountType: req.AccountType,
			Password:    req.PaymentPassword,
			IdenNo:      custInfo.CustNo,
			NonStr:      req.NonStr,
		}
		authRet, authErr := i.AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &authReq)
		if authErr != nil {
			ss_log.Error("调用auth-srv服务失败,err:%v", authErr)
			return nil, ss_err.SystemErr
		}
		authRet.ResultCode = ss_err.AuthSrvRetCode(authRet.ResultCode)
		if authRet.ResultCode != ss_err.Success {
			ss_log.Error("支付密码校验失败, resultCode=%v, CheckPayPWDRequest=%v", authRet.ResultCode, strext.ToJson(authReq))
			return nil, authRet.ResultCode
		}
	}

	reply := new(PayReply)

	//支付方式：账户余额
	if req.PaymentMethod == constants.PayMethodBalance || req.PaymentMethod == "" {
		// 获取用户和服务商的虚拟账户类型
		custVaType := global.GetUserVAccType(bill.CurrencyType, true)
		businessVaType := global.GetBusinessVAccType(bill.CurrencyType, false)
		if custVaType == 0 || businessVaType == 0 {
			ss_log.Error("获取用户和服务商的虚拟账户类型失败,CurrencyType:%v", bill.CurrencyType)
			return nil, ss_err.SystemErr
		}

		// 查询用户虚拟账户是否存在
		custVaccNo, err := dao.VaccountDaoInst.GetVaccountNo(req.AccountNo, custVaType)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("用户虚拟账户不存在,AccountNo:%v,custVaType:%v", req.AccountNo, custVaType)
				return nil, ss_err.VirtualAccountNotExist
			}
			ss_log.Error("查询用户虚拟账户失败,err:%v,AccountNo:%v,custVaType:%v", err, req.AccountNo, custVaType)
			return nil, ss_err.SystemErr
		}

		// 查询商户虚拟账户是否存在
		businessVaccNo, err := dao.VaccountDaoInst.GetVaccountNo(bill.BusinessAccountNo, businessVaType)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("商户虚拟账户不存在,BusinessAccountNo:%v,businessVaType:%v", bill.BusinessAccountNo, businessVaType)
				return nil, ss_err.VirtualAccountNotExist
			}
			ss_log.Error("查询用户虚拟账户失败,err:%v,BusinessAccountNo:%v,businessVaType:%v", err, bill.BusinessAccountNo, businessVaType)
			return nil, ss_err.SystemErr
		}

		cycle, err := dao.BusinessSceneSignedDaoInst.GetBusinessSettleCycle(bill.BusinessNo, bill.SceneNo, constants.SignedStatusPassed)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("没有查询到产品[%v]对应结算周期, businessNo=%v, err=%v", bill.TradeType, bill.BusinessNo, err)
				return nil, ss_err.ProductUnsigned
			}
			ss_log.Error("查询商家[%v]签约产品[%v]的结算周期失败, err=%v", bill.BusinessNo, bill.TradeType, err)
			return nil, ss_err.SystemErr
		}

		// 获取数据库连接
		dbHandler := db.GetDB(constants.DB_CRM)
		if dbHandler == nil {
			ss_log.Error("获取数据连接失败")
			return nil, ss_err.SystemErr
		}
		defer db.PutDB(constants.DB_CRM, dbHandler)

		// 开启事物
		tx, err := dbHandler.BeginTx(ctx, nil)
		if err != nil {
			ss_log.Error("开启事物失败,err:%v", err)
			return nil, ss_err.SystemErr
		}

		// 1.减少用户虚拟账户的金额
		custBalance, custFrozenBalance, err := dao.VaccountDaoInst.MinusBalance(tx, custVaccNo, bill.Amount)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("减少用户虚拟账户的金额失败,err:%v,custVaccNo:%v,amount:%v", err, custVaccNo, bill.Amount)
			return nil, ss_err.SystemErr
		}

		// 判断余额是否为负数
		if r, err := ss_func.JudgeAmountPositiveOrNegative(custBalance); err != nil || r < 0 {
			ss_sql.Rollback(tx)
			ss_log.Error("用户账号余额不足,err:%v,custBalance:%v, result:%v", err, custBalance, r)
			return nil, ss_err.BalanceNotEnough
		}

		// 2.记录用户账户变动日志
		log1 := dao.LogVaccountDao{
			VaccountNo:    custVaccNo,
			OpType:        constants.VaOpType_Minus,
			Amount:        bill.Amount,
			Balance:       custBalance,
			FrozenBalance: custFrozenBalance,
			Reason:        constants.VaReason_Cust_Pay_Order,
			BizLogNo:      req.OrderNo,
		}
		if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log1); err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("插入用户账户变动日志失败,err:%v,data:%+v", err, log1)
			return nil, ss_err.SystemErr
		}

		// 3.增加商家虚拟账户(未结算)的金额
		businessBalance, businessFrozenBalance, err := dao.VaccountDaoInst.PlusBalance(tx, businessVaccNo, bill.Amount)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("增加商家虚拟账户的金额失败,err:%v,businessVaccNo:%v,amount:%v", err, businessVaccNo, bill.Amount)
			return nil, ss_err.SystemErr
		}

		// 4.记录商家账户变动日志
		log2 := dao.LogVaccountDao{
			VaccountNo:    businessVaccNo,
			OpType:        constants.VaOpType_Add,
			Amount:        bill.Amount,
			Balance:       businessBalance,
			FrozenBalance: businessFrozenBalance,
			Reason:        constants.VaReason_Business_Payee,
			BizLogNo:      req.OrderNo,
		}
		if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log2); err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("插入商家账户变动日志失败,err:%v,data:%+v", err, log2)
			return nil, ss_err.SystemErr
		}

		//订单结算时间
		settleTime := getSettleDate(cycle)
		if settleTime < 0 {
			ss_sql.Rollback(tx)
			ss_log.Error("获取订单结算时间失败，cycle=%v", cycle)
			return nil, ss_err.SystemErr
		}

		//5.更新订单为已支付(更改订单状态,更新付款人账号id,虚账id,收款人虚账id,结算周期，结算日期)和支付方式
		paidData := dao.UpdateOrderPaidData{
			OrderNo:            req.OrderNo,
			AccountNo:          req.AccountNo,
			VaccountNo:         custVaccNo,
			BusinessVaccountNo: businessVaccNo,
			PaymentMethod:      req.PaymentMethod,
			PayTime:            ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat),
			Cycle:              cycle,
			SettleDate:         settleTime,
		}
		updateErr := dao.BusinessBillDaoInst.UpdateOrderPaid(tx, paidData)
		if updateErr != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("更新订单信息失败,err:%v, paidData:%+v", updateErr, paidData)
			return nil, ss_err.SystemErr
		}

		if req.SignKey != "" { //如果用户是指纹无密码支付
			if err := dao.LogAppFingerprintPayDaoInstance.AddTx(tx, &dao.LogAppFingerprintPayData{
				AccountNo:    req.AccountNo,
				DeviceUuid:   req.DeviceUuid,
				SignKey:      req.SignKey,
				OrderNo:      bill.OrderNo,
				OrderType:    constants.VaReason_Cust_Pay_Order,
				Amount:       bill.Amount,
				CurrencyType: bill.CurrencyType,
			}); err != nil {
				ss_sql.Rollback(tx)
				ss_log.Error("指纹无密码支付插入日志失败,err[%v],paidData[%v]", err, paidData)
				return nil, ss_err.SystemErr
			}
		}

		ss_sql.Commit(tx)

		//发布订单支付结果异步通知事件
		pushMsg := &notifyProto.PaySystemResultNotify{
			OrderNo:   req.OrderNo,
			OrderType: constants.VaReason_Cust_Pay_Order,
		}
		publishEventErr := common.PayResultNotifyEvent.Publish(context.TODO(), pushMsg)
		if publishEventErr != nil {
			ss_log.Error("支付结果异步通知事件推送失败, err=[%v], topic=[%v], msg=[%v]", publishEventErr, constants.PaySystemResultNotify, strext.ToJson(pushMsg))
		}

		if cycle == common.SettleToD0 {
			settleId, errCode := b.SingleOrderSettle(bill.OrderNo)
			if errCode != ss_err.Success {
				ss_log.Error("D0订单结算失败，orderNo=%v, err=%v", bill.OrderNo, errCode)
			} else {
				ss_log.Info("D0订单结算成功，orderNo=%v, settleId=%v", bill.OrderNo, settleId)
			}
		}

		reply.PayTime = paidData.PayTime
	} else if req.PaymentMethod == constants.PayMethodBankCard {
		ss_log.Error("目前不支持银行卡支付")
		return nil, ss_err.BankCardNotSupported
	}

	reply.OrderNo = bill.OrderNo
	reply.OutOrderNo = bill.OutOrderNo
	reply.OrderStatus = constants.BusinessOrderStatusPay
	reply.CreateTime = bill.CreateTime
	reply.Subject = bill.Subject
	return reply, ss_err.Success
}

func getSettleDate(cycle string) int64 {
	if cycle == "" {
		return -1
	}
	switch cycle {
	case common.SettleToD0:
		//钱实际到账时间得看程序，所以减少不必要麻烦d0结算时间往后延迟5秒
		return ss_time.Now(global.Tz).Add(5 * time.Second).Unix()
	case common.SettleToT1:
		fallthrough
	case common.SettleToT2:
		fallthrough
	case common.SettleToT3:
		fallthrough
	case common.SettleToT4:
		fallthrough
	case common.SettleToT5:
		fallthrough
	case common.SettleToT6:
		fallthrough
	case common.SettleToT7:
		slice := strings.Split(cycle, "+")
		if slice == nil {
			return -1
		}
		return ss_time.Now(global.Tz).AddDate(0, 0, strext.ToInt(slice[1])).Unix()
	default:
		return -1
	}
}
