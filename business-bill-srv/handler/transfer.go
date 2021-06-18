package handler

import (
	"a.a/cu/db"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/common"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/business-bill-srv/i"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	billProto "a.a/mp-server/common/proto/bill"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	notifyProto "a.a/mp-server/common/proto/notify"
	pushProto "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"github.com/shopspring/decimal"
)

// 企业转账(单笔) todo 未完成
func (b *BusinessBillHandler) EnterpriseTransfer(ctx context.Context, req *businessBillProto.EnterpriseTransferRequest, reply *businessBillProto.EnterpriseTransferReply) error {
	resultCode, err := CheckEnterpriseTransferReq(req)
	if err != nil {
		ss_log.Error("请求参数有误，err=%v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
	}

	app, err := dao.BusinessAppDaoInst.GetAppInfoByAppId(req.AppId)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("app[%v]不存在", req.AppId)
			reply.ResultCode = ss_err.AppNotExist
			reply.Msg = ss_err.GetMsg(ss_err.AppNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询app[%v]信息失败，err=%v", req.AppId, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	signed, err := dao.BusinessSceneSignedDaoInst.GetSignedByTradeType(app.BusinessNo, constants.TradeTypeEnterprisePay)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("应用[%v]未签约[%v]产品", req.AppId, constants.TradeTypeEnterprisePay)
			reply.ResultCode = ss_err.ProductUnsigned
			reply.Msg = ss_err.GetMsg(ss_err.ProductUnsigned, req.Lang)
			return nil
		}
		ss_log.Error("查询app[%v]签约[%v]信息失败，err=%v", req.AppId, constants.TradeTypeEnterprisePay, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	//检查商户签约
	resultCode, err = CheckBusinessSigned(signed)
	if err != nil {
		ss_log.Error("app[%v]签约检测失败, err=%v", req.AppId, err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//检查产品是否可用
	isEnabled, err := dao.BusinessSceneDao.GetSceneIsEnabled(signed.SceneNo)
	if err != nil {
		ss_log.Error("查询产品是否可用失败，SceneNo=%v, err=%v", signed.SceneNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	if !isEnabled {
		ss_log.Error("产品[%v]已被系统禁用，暂时不可交易", signed.SceneNo)
		reply.ResultCode = ss_err.SceneDisabled
		reply.Msg = ss_err.GetMsg(ss_err.SceneDisabled, req.Lang)
		return nil
	}

	// 查询商户交易配置（商户账号，商户是否启用，收款权限）
	business, err := dao.BusinessDaoInst.GetTransConfig(app.BusinessNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("商户不存在,BusinessNo:%v", app.BusinessNo)
			reply.ResultCode = ss_err.QrCodeNotInvalid
			reply.Msg = ss_err.GetMsg(ss_err.QrCodeNotInvalid, req.Lang)
			return nil
		}
		ss_log.Error("查询商户交易配置失败,err:%v,BusinessNo:%v", err, app.BusinessNo)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	// 商户状态不可用
	if !business.IsEnabled {
		ss_log.Error("商户[%v]状态为[%v]不可用", business.BusinessNo, business.IsEnabled)
		reply.ResultCode = ss_err.BusinessNotAvailable
		return nil

	}
	//商户是否有出款权限
	if business.OutgoAuthorization == constants.BusinessOutGoAuthDisabled {
		ss_log.Error("商户[%v]没有收款权限, outGoAuthorization:%v", business.BusinessNo, business.IncomeAuthorization)
		reply.ResultCode = ss_err.AccountNoNotTradeForbid
		return nil
	}

	//查询收款人是否存在
	var payeeInfo *dao.AccountDao
	if req.PayeePhone != "" {
		payeeAccount := ss_func.ComposeAccountByPhoneCountryCode(ss_func.PrePhone(req.PayeeCountryCode, req.PayeePhone), req.PayeeCountryCode)
		if business.BusinessAcc == payeeAccount {
			ss_log.Error("不能转账给自己, payer=%v, payee=%v", business.BusinessAcc, payeeAccount)
			reply.ResultCode = ss_err.PayeeAccountErr
			reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountErr, req.Lang)
			return nil
		}
		payeeInfo, err = dao.AccountDaoInst.GetInfoByAccount(payeeAccount)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("收款人[%v]不存在", payeeAccount)
				reply.ResultCode = ss_err.PayeeAccountNoNotExist
				reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountNoNotExist, req.Lang)
				return nil
			}
			ss_log.Error("查询收款人[%v]失败, err=%v", payeeAccount, err)
			reply.ResultCode = ss_err.SystemErr
			reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
			return nil
		}
	} else if req.PayeeEmail != "" {
		if business.BusinessAcc == req.PayeeEmail {
			ss_log.Error("不能转账给自己, payer=%v, payee=%v", business.BusinessAcc, req.PayeeEmail)
			reply.ResultCode = ss_err.PayeeAccountErr
			reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountErr, req.Lang)
			return nil
		}
		payeeInfo, err = dao.AccountDaoInst.GetInfoByAccount(req.PayeeEmail)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("收款人[%v]不存在", req.PayeeEmail)
				reply.ResultCode = ss_err.PayeeAccountNoNotExist
				reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountNoNotExist, req.Lang)
				return nil
			}
			ss_log.Error("查询收款人[%v]失败, err=%v", req.PayeeEmail, err)
			reply.ResultCode = ss_err.SystemErr
			reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
			return nil
		}
	}

	businessAccNo, err := dao.BusinessDaoInst.QueryAccNoByBusinessNo(app.BusinessNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("商家[%v]不存在", app.BusinessNo)
			reply.ResultCode = ss_err.BusinessNotExist
			reply.Msg = ss_err.GetMsg(ss_err.BusinessNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询商家[%v]信息失败，err=%v", app.BusinessNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//不能自己转给自己
	if businessAccNo == payeeInfo.AccountNo {
		ss_log.Error("不能转账给自己, payer=%v, payee=%v", businessAccNo, payeeInfo.AccountNo)
		reply.ResultCode = ss_err.PayeeAccountErr
		return nil
	}

	//计算手续费
	fee, rate, getFeeCode := getBusinessFee(req.Amount, req.CurrencyType)
	if getFeeCode != ss_err.Success {
		reply.ResultCode = getFeeCode
		return nil
	}

	//判断付款商余额是否足够
	businessVAccType := global.GetBusinessVAccType(req.CurrencyType, true)
	balance, err := dao.VaccountDaoInst.GetBalanceByAccNo(businessAccNo, businessVAccType)
	if err != nil {
		ss_log.Error("查询商家虚账失败, err=%v", err)
		reply.ResultCode = ss_err.ParamErr
		return nil
	}
	if strext.ToInt(balance) < strext.ToInt(ss_count.Add(req.Amount, fee)) {
		ss_log.Error("商家余额不足，BusinessAccNo=%v, CurrencyType=%v,", businessAccNo, req.CurrencyType)
		reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败, err=%v", err)
		reply.ResultCode = ss_err.SystemErr
		return nil
	}

	//插入转账订单
	order := new(dao.BusinessTransferOrderDao)
	order.FromBusinessNo = app.BusinessNo
	order.FromAccountNo = businessAccNo
	order.ToAccountNo = payeeInfo.AccountNo
	order.Amount = req.Amount
	order.RealAmount = req.Amount
	order.CurrencyType = req.CurrencyType
	order.PaymentType = constants.ORDER_PAYMENT_BALANCE
	order.Fee = fee
	order.Rate = rate
	order.OrderStatus = constants.BusinessTransferOrderStatusPending
	order.Remarks = req.Remark
	order.TransferType = constants.BusinessTransferOrderTypeEnterprise
	order.OutTransferNo = req.OutTransferNo
	orderNo, err := dao.BusinessTransferOrderDaoInst.InsertTx(tx, order)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入转账订单失败,data:%v, err=%v", strext.ToJson(order), err)
		reply.ResultCode = ss_err.SystemErr
		return nil
	}
	//插入转账订单部署-企业转账
	log := new(dao.EnterpriseTransferOrderDao)
	log.AppId = req.AppId
	log.NotifyUrl = req.NotifyUrl
	log.TransferLogNo = orderNo
	_, err = dao.EnterpriseTransferOrderDaoInst.InsertTx(tx, log)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入企业转账日志失败，data:%v, err=%v", strext.ToJson(log), err)
		reply.ResultCode = ss_err.SystemErr
		return nil
	}
	ss_sql.Commit(tx)

	go SyncTransfer(&SyncTransferRequest{
		TransferNo: orderNo,
		Lang:       req.Lang,
	})

	data := new(businessBillProto.EnterpriseTransferOrder)
	data.TransferNo = orderNo
	data.OutTransferNo = req.OutTransferNo
	data.Amount = req.Amount
	data.CurrencyType = req.CurrencyType
	data.TransferStatus = constants.BusinessTransferOrderStatusPending

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.Order = data
	return nil

}

type SyncTransferRequest struct {
	TransferNo string
	Lang       string
}

//异步处理转账
func SyncTransfer(req *SyncTransferRequest) {
	//调用内部转账接口
	resp, err := i.BillHandlerInst.Client.EnterpriseTransferToUser(context.TODO(), &billProto.EnterpriseTransferToUserRequest{
		TransferNo:   req.TransferNo,
		Lang:         req.Lang,
		TransferType: constants.BusinessTransferOrderTypeEnterprise,
	})
	if err != nil {
		ss_log.Error("调用bill-src.AddBusinessTransfer()失败, err=%v", err)
		return
	}

	//结果异步通知
	notifyEv := &notifyProto.PaySystemResultNotify{
		OrderNo:   req.TransferNo,
		OrderType: constants.VaReason_BusinessTransferToBusiness,
	}
	err = common.PayResultNotifyEvent.Publish(context.TODO(), notifyEv)
	if err != nil {
		ss_log.Error("notify-srv.Publish()失败, event=[%v], err=%v", strext.ToJson(notifyEv), err)
	}

	resp.ResultCode = ss_err.BillSrvRetCode(resp.ResultCode)
	if resp.ResultCode != ss_err.Success {
		ss_log.Error("转账失败, errCode=%v", resp.ResultCode)
		return
	}

	transferOrder := resp.Order
	//消息中心
	if transferOrder.PayeeAccType == constants.AccountType_USER {
		//添加APP消息
		msg := new(dao.LogAppMessagesDao)
		msg.OrderNo = transferOrder.OrderNo
		msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_TRANSFER_Apply
		msg.OrderType = constants.VaReason_BusinessTransferToBusiness
		msg.AccountNo = transferOrder.PayeeAccNo
		msg.OrderStatus = constants.OrderStatus_Paid
		err = dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
		if err != nil {
			ss_log.Error("aAddMessagesErr=[%v]", err)
		}

		//推送消息给用户——只有收款方为普通用户是才有APP消息推送
		pushMessageToUser(transferOrder.PayeeAccNo, transferOrder.Amount, transferOrder.CurrencyType, constants.Template_TransferSuccess, req.Lang)

	} else {
		//消息模板变量
		timeString := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
		moneyType := getCurrencySign(transferOrder.CurrencyType)
		amount := common.NormalAmountByMoneyType(transferOrder.CurrencyType, transferOrder.Amount)

		var message string
		temp, err := dao.LogBusinessMessageDao.GetTemplate(constants.Template_TransferSuccess)
		if err != nil {
			ss_log.Error("查询商家中心消息模板[%v]失败", constants.Template_TransferSuccess)
			message = fmt.Sprintf("您于%s收到一笔转账%s%s, 请注意查收！", timeString, moneyType, amount)
		} else {
			message = dao.LangDaoInstance.GetLangTextByKey(temp.ContentKey, req.Lang)
			args := []interface{}{timeString, moneyType, amount}
			if req.Lang == constants.LangEnUS {
				args = []interface{}{moneyType, amount, timeString}
			}
			//填充内容
			message = fmt.Sprintf(message, args...)
		}

		d := new(dao.LogBusinessMessage)
		d.LogNo = transferOrder.OrderNo
		d.IsRead = "0"
		d.AccountType = transferOrder.PayeeAccType
		d.AccountNo = transferOrder.PayeeAccNo
		d.Content = message
		errAddMessages := dao.LogBusinessMessageDao.Insert(d)
		if err != nil {
			ss_log.Error("插入商家中心消息日志失败，err[%v]", errAddMessages)
		}
	}

	return
}

func pushMessageToUser(accountNo, amount, currencyType, tempNo, lang string) {
	timeString := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	// 修正金额
	amount = common.NormalAmountByMoneyType(currencyType, amount)
	moneyType := dao.LangDaoInstance.GetLangTextByKey(strings.ToLower(currencyType), lang)

	args := []string{
		timeString, amount, moneyType,
	}
	if lang == constants.LangEnUS {
		args = []string{
			amount, moneyType, timeString,
		}
	}
	// 消息推送
	ev := &pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				AccountNo:   accountNo,
				AccountType: constants.AccountType_USER,
			},
		},
		TempNo: tempNo,
		Args:   args,
	}
	ss_log.Info("publishing %+v\n", ev)
	if err := common.PushEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("消息推送到用户[%v]出错。error : %v", accountNo, err)
	}
}

//获取币种符号
func getCurrencySign(currencyType string) string {
	currencyType = strings.ToUpper(currencyType)
	if currencyType == constants.CURRENCY_UP_USD {
		return "＄"
	} else if currencyType == constants.CURRENCY_UP_KHR {
		return "៛"
	}
	return ""
}

//获取商家转账手续费
func getBusinessFee(amount, currencyType string) (fee, rate, errCode string) {
	currencyType = strings.ToUpper(currencyType)

	//查询商家转账配置
	transferConf, err := dao.GlobalParamDaoInstance.GetBusinessTransferParamValue()
	if err != nil {
		ss_log.Error("查询商家转账配置失败, err=%v", err)
		return "", "", ss_err.SystemErr
	}

	//判断交易金额是否超出限制,并计算手续费
	var feesDeci decimal.Decimal
	switch currencyType {
	case constants.CURRENCY_UP_USD:
		if strext.ToInt64(amount) >= transferConf.USDMinAmount && strext.ToInt64(amount) <= transferConf.USDMaxAmount {
			rate = strext.ToString(transferConf.USDRate)
			feesDeci = ss_count.CountFees(amount, rate, strext.ToStringNoPoint(transferConf.USDMinFee))
		} else {
			ss_log.Error("转账金额超出限制，")
			return "", "", ss_err.TransactionAmountLimit
		}
	case constants.CURRENCY_UP_KHR:
		if strext.ToInt64(amount) >= transferConf.KHRMinAmount && strext.ToInt64(amount) <= transferConf.KHRMaxAmount {
			rate = strext.ToString(transferConf.KHRRate)
			feesDeci = ss_count.CountFees(amount, rate, strext.ToStringNoPoint(transferConf.KHRMinFee))
		} else {
			ss_log.Error("转账金额超出限制，")
			return "", "", ss_err.TransactionAmountLimit
		}
	}
	// 取整
	fee = ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()

	return fee, rate, ss_err.Success
}
