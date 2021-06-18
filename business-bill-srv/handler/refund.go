package handler

import (
	notifyProto "a.a/mp-server/common/proto/notify"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/business-bill-srv/common"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_sql"
)

//商家中心-支付完成订单退款
func (b *BusinessBillHandler) BusinessBillRefund(ctx context.Context, req *businessBillProto.BusinessBillRefundRequest, reply *businessBillProto.BusinessBillRefundReply) error {
	//查询订单详情
	order, err := dao.BusinessBillDaoInst.GetBillByLogNoAndBusinessNo(req.OrderNo, req.BusinessNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单不存在，orderNo=%v", req.OrderNo)
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询订单失败，orderNo=%v， err=%v", req.OrderNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	if !util.InSlice(order.OrderStatus, []string{constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRebatesRefund}) {
		ss_log.Error("订单状态不对，不能退款")
		reply.ResultCode = ss_err.OrderNotRefundable
		reply.Msg = ss_err.GetMsg(ss_err.OrderNotRefundable, req.Lang)
		return nil
	}

	//当前只支持全额退款
	if req.RefundAmount != order.Amount {
		ss_log.Error("退款金额与交易订单金额不一致，refundAmount=%v, transAmount=%v", req.RefundAmount, order.Amount)
		reply.ResultCode = ss_err.RefundAmountDisagree
		reply.Msg = ss_err.GetMsg(ss_err.RefundAmountDisagree, req.Lang)
		return nil
	}

	checkResp := RefundMonitor(req.RefundAmount, order)
	if checkResp.ResultCode != ss_err.Success {
		ss_log.Error("未通过退款检测，result=%v", ss_err.GetMsg(checkResp.ResultCode, req.Lang))
		reply.Msg = ss_err.GetMsg(checkResp.ResultCode, req.Lang)
		reply.ResultCode = checkResp.ResultCode
		return nil
	}

	//插入商家退款订单日志
	d := new(dao.BusinessRefundOrderDao)
	d.Amount = checkResp.BusinessRefundAmount
	d.PayOrderNo = order.OrderNo
	d.RefundStatus = constants.BusinessRefundStatusPending
	d.Remarks = req.RefundReason
	refundNo, err := dao.BusinessRefundOrderDaoInst.Insert(d)
	if err != nil {
		ss_log.Error("插入退款日志失败，data:=%v, err=%v", strext.ToJson(d), err)
		reply.ResultCode = ss_err.SystemErr
		return nil
	}

	//退款-减扣虚账金额
	refundedReq := &RefundedAmountRequest{
		PlatformVAccNo: checkResp.PlatformVAccNo,
		BusinessVAccNo: checkResp.BusinessVAccNo,
		CustVAccNo:     checkResp.PayeeVAccNo,
		Amount:         checkResp.BusinessRefundAmount,
		CurrencyType:   order.CurrencyType,
		Fee:            checkResp.PlatformRefundAmount,
		RefundOrderNo:  refundNo,
		RefundType:     checkResp.RefundType,
		PayOrderNo:     order.OrderNo,
	}
	resultCode, userAmount := refundedAmount(refundedReq)
	if resultCode != ss_err.Success {
		//退款失败，修改退款订单状态
		err := dao.BusinessRefundOrderDaoInst.UpdatePendingOrderByOrderNo(constants.BusinessRefundStatusFail, refundNo)
		if nil != err {
			ss_log.Error("修改退款订单状态失败，refundNo=%v, err=%v", refundNo, err)
			reply.ResultCode = ss_err.ERR_SYSTEM
			return nil
		}
	}

	//添加退款消息
	msg := new(dao.LogAppMessagesDao)
	msg.OrderNo = refundNo
	msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Business_Refund_Cust
	msg.OrderType = constants.VaReason_BusinessRefund
	msg.AccountNo = order.AccountNo
	msg.OrderStatus = constants.OrderStatus_Paid
	err = dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
	if err != nil {
		ss_log.Error("aAddMessagesErr=[%v]", err)
	}

	//推送消息给用户
	pushMessageToUser(order.AccountNo, userAmount, order.CurrencyType, constants.Template_RefundSuccess, req.Lang)

	reply.ResultCode = ss_err.Success
	reply.RefundNo = refundNo
	reply.RefundAmount = checkResp.BusinessRefundAmount
	return nil
}

//接口-支付完成订单退款
func (b *BusinessBillHandler) ApiPayRefund(ctx context.Context, req *businessBillProto.ApiPayRefundRequest, reply *businessBillProto.ApiPayRefundReply) error {
	resultCode, err := CheckApiPayRefundReq(req)
	if err != nil {
		ss_log.Error("参数有误, err=%v", err)
		reply.ResultCode = resultCode
		reply.Msg = ss_err.GetMsg(resultCode, req.Lang)
		return nil
	}

	//查询订单详情
	order, err := dao.BusinessBillDaoInst.AppQueryOrder(req.AppId, req.OrderNo, req.OutOrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单不存在，orderNo=%v", req.OrderNo)
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询订单失败，orderNo=%v， err=%v", req.OrderNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//检查订单状态
	if !util.InSlice(order.OrderStatus, []string{constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRebatesRefund}) {
		ss_log.Error("订单状态不满足退款条件")
		switch order.OrderStatus {
		case constants.BusinessOrderStatusPending:
			fallthrough
		case constants.BusinessOrderStatusPayTimeOut:
			reply.ResultCode = ss_err.OrderUnpaid
		case constants.BusinessOrderStatusRefund:
			reply.ResultCode = ss_err.OrderFullRefund
		}
		reply.Msg = ss_err.GetMsg(reply.ResultCode, req.Lang)
		return nil
	}

	checkResp := RefundMonitor(req.RefundAmount, order)
	if checkResp.ResultCode != ss_err.Success {
		reply.ResultCode = checkResp.ResultCode
		reply.Msg = ss_err.GetMsg(checkResp.ResultCode, req.Lang)
		return nil
	}

	//插入退款订单日志
	d := new(dao.BusinessRefundOrderDao)
	d.Amount = checkResp.BusinessRefundAmount
	d.PayOrderNo = order.OrderNo
	d.RefundStatus = constants.BusinessRefundStatusPending
	d.Remarks = req.RefundReason
	d.OutRefundNo = req.OutRefundNo
	d.NotifyUrl = req.NotifyUrl
	d.NotifyStatus = constants.NotifyStatusNOT
	refundNo, err := dao.BusinessRefundOrderDaoInst.Insert(d)
	if err != nil {
		ss_log.Error("插入退款日志失败，data:=%v, err=%v", strext.ToJson(d), err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	go func() {
		//退款-减扣虚账金额
		refundedReq := &RefundedAmountRequest{
			PlatformVAccNo: checkResp.PlatformVAccNo,
			BusinessVAccNo: checkResp.BusinessVAccNo,
			CustVAccNo:     checkResp.PayeeVAccNo,
			Amount:         checkResp.BusinessRefundAmount,
			CurrencyType:   order.CurrencyType,
			Fee:            checkResp.PlatformRefundAmount,
			RefundOrderNo:  refundNo,
			RefundType:     checkResp.RefundType,
			PayOrderNo:     order.OrderNo,
		}
		resultCode, userAmount := refundedAmount(refundedReq)

		//异步通知
		notifyEv := &notifyProto.PaySystemResultNotify{
			OrderNo:   refundNo,
			OrderType: constants.VaReason_BusinessRefund,
		}
		err = common.PayResultNotifyEvent.Publish(context.TODO(), notifyEv)
		if err != nil {
			ss_log.Error("notify-srv.Publish()失败, event=[%v], err=%v", strext.ToJson(notifyEv), err)
		}

		if resultCode != ss_err.Success {
			//退款失败，修改退款订单状态
			err := dao.BusinessRefundOrderDaoInst.UpdatePendingOrderByOrderNo(constants.BusinessRefundStatusFail, refundNo)
			if nil != err {
				ss_log.Error("修改退款订单状态失败，refundNo=%v, err=%v", refundNo, err)
			}
			return
		}

		//添加退款到的账号推送消息
		msg := new(dao.LogAppMessagesDao)
		msg.OrderNo = refundNo
		msg.AppMessType = constants.LOG_APP_MESSAGES_ORDER_TYPE_Business_Refund_Cust
		msg.OrderType = constants.VaReason_BusinessRefund
		msg.AccountNo = order.AccountNo
		msg.OrderStatus = constants.OrderStatus_Paid
		aAddMessagesErr := dao.LogAppMessagesDaoInst.AddLogAppMessages(msg)
		if aAddMessagesErr != nil {
			ss_log.Error("errAddMessages2=[%v]", aAddMessagesErr)
		}
		//推成功退款消息给用户
		pushMessageToUser(order.AccountNo, userAmount, order.CurrencyType, constants.Template_RefundSuccess, req.Lang)
	}()

	time.Sleep(1 * time.Second)

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = order.OrderNo
	reply.OutOrderNo = order.OutOrderNo
	reply.RefundNo = refundNo
	reply.OutRefundNo = req.OutRefundNo
	reply.RefundAmount = req.RefundAmount
	reply.RealRefundAmount = checkResp.BusinessRefundAmount
	reply.RefundStatus = constants.BusinessRefundStatusPending
	return nil
}

type MonitorResponse struct {
	ResultCode           string //结果码
	PayeeVAccNo          string //用户虚账
	BusinessVAccNo       string //商家虚账
	PlatformVAccNo       string //平台虚账
	BusinessRefundAmount string //商家实退金额
	PlatformRefundAmount string //平台实退金额
	RefundType           string //退款类型(部分退款，全额退款)
}

/**
1.退款交易-商家、用户监测
*/
func RefundMonitor(refundAmount string, order *dao.BusinessBillDao) *MonitorResponse {
	reply := new(MonitorResponse)
	//查询付款人虚账
	payeeVaType := global.GetUserVAccType(order.CurrencyType, true)
	payeeVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(order.AccountNo, payeeVaType)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("用户虚账不存在，AccountNo=%v", order.AccountNo)
			reply.ResultCode = ss_err.VirtualAccountNotExist
			return reply
		}
		ss_log.Error("查询用户虚账失败，AccountNo=%v, err=%v", order.AccountNo, err)
		reply.ResultCode = ss_err.SystemErr
		return reply
	}

	//检查退款金额，识别退款类型(部分退款，全额退款)
	retCode, refundType, err := checkRefundAmount(refundAmount, order.Amount, order.OrderNo)
	if err != nil {
		ss_log.Error("退款金额有误，err=%v", err)
		reply.ResultCode = retCode
		return reply
	}

	/**
	查商家虚账, 计算出商家应退金额
	并根据订单是否已结算，计算出平台应退金额
	*/
	var businessVAccNo, deductAmount, platVAccNo, fee string
	if order.SettleId == "" {
		businessVAccNo = order.BusinessVaccountNo
		deductAmount = order.Amount
	} else {
		//已结算订单 todo 目前已结算订单考虑全额退款
		if refundType != constants.BusinessOrderStatusRefund {
			ss_log.Error("目前已结算订单不支持部分退款")
			reply.ResultCode = ss_err.OrderNotRefundable
			return reply
		} else {
			//查询平台虚账和平台应退金额
			platAccNo, err := dao.GlobalParamDaoInstance.GetParamValue(constants.GlobalParamKeyAccPlat)
			if err != nil {
				ss_log.Error("查询平台账号失败，paramKey=%v, err=%v", constants.GlobalParamKeyAccPlat, platAccNo)
				reply.ResultCode = ss_err.SystemErr
				return reply
			}
			platVAccType := global.GetPlatFormVAccType(order.CurrencyType)
			platVAccNo, err = dao.VaccountDaoInst.GetVaccountNo(platAccNo, platVAccType)
			if err != nil {
				ss_log.Error("查询平台虚账失败，paramKey=%v, err=%v", constants.GlobalParamKeyAccPlat, platAccNo)
				reply.ResultCode = ss_err.SystemErr
				return reply
			}
			fee = order.Fee

			//查询商家虚账和商家应退金额
			if order.BusinessAccountNo == "" {
				order.BusinessAccountNo, err = dao.BusinessDaoInst.QueryAccNoByBusinessNo(order.BusinessNo)
				if err != nil {
					ss_log.Error("商户虚账不存在，BusinessNo=%v", order.BusinessNo)
					reply.ResultCode = ss_err.BusinessNotExist
					return reply
				}
			}
			businessVaType := global.GetBusinessVAccType(order.CurrencyType, true)
			businessVAccNo, err = dao.VaccountDaoInst.GetVaccountNo(order.BusinessAccountNo, businessVaType)
			if err != nil {
				if err == sql.ErrNoRows {
					ss_log.Error("商户虚账不存在，BusinessAccNo=%v", order.BusinessAccountNo)
					reply.ResultCode = ss_err.VirtualAccountNotExist
					return reply
				}
				ss_log.Error("查询商户虚账失败，BusinessAccNo=%v, err=%v", order.BusinessAccountNo, err)
				reply.ResultCode = ss_err.SystemErr
				return reply
			}
			deductAmount = order.RealAmount
		}
	}

	//查询商家余额
	balance, err := dao.VaccountDaoInst.GetBalanceByVAccNo(businessVAccNo)
	if err != nil {
		ss_log.Error("查询商户未结算虚账余额失败, BusinessVAccNo=%v, err=%v", businessVAccNo, err)
		reply.ResultCode = ss_err.SystemErr
		return reply
	}
	//判断余额是否足够退款， 计算结果小于0，Sign()返回-1；等于0，返回0；大于0，返回1
	if ss_count.Sub(strext.ToString(balance), deductAmount).Sign() < 0 {
		ss_log.Error("商户虚账余额不足, BusinessVAccNo=%v, CurrencyType=%v, balance=%v", businessVAccNo, order.CurrencyType, balance)
		reply.ResultCode = ss_err.BalanceNotEnough
		return reply
	}

	reply.ResultCode = ss_err.Success
	reply.PayeeVAccNo = payeeVAccNo
	reply.BusinessVAccNo = businessVAccNo
	reply.BusinessRefundAmount = deductAmount
	reply.PlatformVAccNo = platVAccNo
	reply.PlatformRefundAmount = fee
	reply.RefundType = refundType
	return reply
}

/**
2.返还金额
*/
type RefundedAmountRequest struct {
	PlatformVAccNo string
	BusinessVAccNo string
	CustVAccNo     string
	Amount         string
	CurrencyType   string
	Fee            string
	RefundOrderNo  string
	RefundType     string
	PayOrderNo     string
}

//退还金额
func refundedAmount(req *RefundedAmountRequest) (errCode, userAmount string) {
	//开启事务
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, txErr := dbHandler.BeginTx(context.TODO(), nil)
	if txErr != nil {
		ss_log.Error("开启事务失败，err=%v", txErr)
		return ss_err.SystemErr, ""
	}

	//平台退还手续费
	if req.PlatformVAccNo != "" && strext.ToInt64(req.Fee) > 0 {
		//减少商平台虚账的余额
		platformBalance, platformFrozenBalance, err := dao.VaccountDaoInst.MinusBalance(tx, req.PlatformVAccNo, req.Fee)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("减少平台账户余额, PlatformVAccNo=%v, err=%v", req.PlatformVAccNo, err)
			return ss_err.SystemErr, ""
		}

		// 判断余额是否为负数
		if r, err := ss_func.JudgeAmountPositiveOrNegative(platformBalance); err != nil || r < 0 {
			ss_sql.Rollback(tx)
			ss_log.Error("平台账户余额不足,err:%v,platformBalance:%v, result:%v", err, platformBalance, r)
			return ss_err.BalanceNotEnough, ""
		}

		//记录账户变动日志
		log1 := dao.LogVaccountDao{
			VaccountNo:    req.PlatformVAccNo,
			OpType:        constants.VaOpType_Minus,
			Amount:        req.Fee,
			Balance:       platformBalance,
			FrozenBalance: platformFrozenBalance,
			Reason:        constants.VaReason_BusinessRefund,
			BizLogNo:      req.RefundOrderNo,
		}
		if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log1); err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("插入平台账户变动日志失败,err:%v,data:%+v", err, log1)
			return ss_err.SystemErr, ""
		}

		//插入手续费盈利扣减记录
		d := &dao.HeadquartersProfit{
			GeneralLedgerNo: req.RefundOrderNo,
			Amount:          req.Fee,
			OrderStatus:     constants.OrderStatus_Paid,
			BalanceType:     strings.ToLower(req.CurrencyType),
			ProfitSource:    constants.ProfitSource_ModernPayOrderRefund,
			OpType:          constants.PlatformProfitMinus,
		}
		_, err = dao.HeadquartersProfitDao.InsertHeadquartersProfit(tx, d)
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("插入平台盈利失败, data=%v, err=%v", strext.ToJson(d), err)
			return ss_err.SystemErr, ""
		}

		//同步总部虚账的余额(等于收益表中的可提现余额)
		err = dao.HeadquartersProfitDao.SyncHeadquartersProfit(tx, req.PlatformVAccNo, req.Fee, strings.ToLower(req.CurrencyType))
		if err != nil {
			ss_sql.Rollback(tx)
			ss_log.Error("同步收益余额失败, headVacc=%v, err=%v", req.PlatformVAccNo, err)
			return ss_err.SystemErr, ""
		}
	}

	//减少商家未结算虚账的余额
	businessBalance, businessFrozenBalance, err := dao.VaccountDaoInst.MinusBalance(tx, req.BusinessVAccNo, req.Amount)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("减少商户账户余额, businessVAccNo=%v, err=%v", req.BusinessVAccNo, err)
		return ss_err.SystemErr, ""
	}

	// 判断余额是否为负数
	if r, err := ss_func.JudgeAmountPositiveOrNegative(businessBalance); err != nil || r < 0 {
		ss_sql.Rollback(tx)
		ss_log.Error("商户账户余额不足,err:%v,businessBalance:%v, result:%v", err, businessBalance, r)
		return ss_err.BalanceNotEnough, ""
	}

	//记录账户变动日志
	log1 := dao.LogVaccountDao{
		VaccountNo:    req.BusinessVAccNo,
		OpType:        constants.VaOpType_Minus,
		Amount:        req.Amount,
		Balance:       businessBalance,
		FrozenBalance: businessFrozenBalance,
		Reason:        constants.VaReason_BusinessRefund,
		BizLogNo:      req.RefundOrderNo,
	}
	if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log1); err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入商家账户变动日志失败,err:%v,data:%+v", err, log1)
		return ss_err.SystemErr, ""
	}

	//增加用户虚账余额，并记录账户变动日志
	userAmount = ss_count.Add(req.Amount, req.Fee)
	userBalance, userFrozenBalance, err := dao.VaccountDaoInst.PlusBalance(tx, req.CustVAccNo, userAmount)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("增加用户虚账余额失败，CustVAccNo=%v, err=%v", req.CustVAccNo, err)
		return ss_err.SystemErr, ""
	}

	//记录账户变动日志
	log2 := dao.LogVaccountDao{
		VaccountNo:    req.CustVAccNo,
		OpType:        constants.VaOpType_Add,
		Amount:        userAmount,
		Balance:       userBalance,
		FrozenBalance: userFrozenBalance,
		Reason:        constants.VaReason_BusinessRefund,
		BizLogNo:      req.RefundOrderNo,
	}
	if err := dao.LogVaccountDaoInst.InsertLogTx(tx, log2); err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("插入用户账户变动日志失败,err:%v,data:%+v", err, log2)
		return ss_err.SystemErr, ""
	}

	//退款成功，修改退款订单状态
	err = dao.BusinessRefundOrderDaoInst.UpdatePendingOrderByOrderNoTx(tx, constants.BusinessRefundStatusSuccess, req.RefundOrderNo)
	if nil != err {
		ss_sql.Rollback(tx)
		ss_log.Error("修改退款订单状态失败，orderNo=%v, err=%v", req.RefundOrderNo, err)
		return ss_err.SystemErr, ""
	}

	//退款成功，修改原支付订单状态
	err = dao.BusinessBillDaoInst.UpdateStatusByOrderNoTx(tx, req.RefundType, req.PayOrderNo)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error("修改支付订单状态失败，orderNo=%v, err=%v", req.PayOrderNo, err)
		return ss_err.SystemErr, ""
	}

	ss_sql.Commit(tx)
	return ss_err.Success, userAmount
}

//检测并获取商家退款金额：已退金额 + 退款金额 <= 交易金额
type BusinessRefundAmount struct {
	ResultCode string
	RefundType string
}

//计算退款金额和退款类型
func checkRefundAmount(refundAmount, transAmount, orderNo string) (resultCode, refundType string, err error) {
	//退剩下的所有金额
	if refundAmount == "" {
		return ss_err.ParamErr, "", errors.New("refundAmount参数为空")
	}

	//订单已退金额
	totalAmount, err := dao.BusinessRefundOrderDaoInst.GetTotalAmountByPayOrderNo(orderNo)
	if err != nil {
		if err != sql.ErrNoRows {
			return ss_err.SystemErr, "", errors.New(fmt.Sprintf("查询订单[%v]已退金额失败,dbErr=%v", orderNo, err))
		}
	}

	//退款金额大
	subRet := ss_count.Sub(transAmount, refundAmount).Sign()
	if subRet == -1 {
		return ss_err.RefundAmountExcessBalance, "", errors.New(fmt.Sprintf("退款金额大于订单金额|交易金额:%v|退款金额:%v", transAmount, refundAmount))
	}

	//金额参数相等(全额退款)
	if subRet == 0 {
		//退金额超出可退金额
		if ss_count.Sub(transAmount, ss_count.Add(refundAmount, totalAmount)).Sign() < 0 {
			return ss_err.RefundAmountExcessBalance, "", errors.New(fmt.Sprintf("退款金额大于订单剩余金额|交易金额:%v|退款金额:%v|已退金额:%v", transAmount, refundAmount, totalAmount))
		}
		return ss_err.Success, constants.BusinessOrderStatusRefund, nil
	}

	//退款金额小(部分退款)
	if subRet == 1 {
		//退还部分金额
		subRet = ss_count.Sub(transAmount, ss_count.Add(refundAmount, totalAmount)).Sign()
		if subRet < 0 {
			//金额有误
			return ss_err.RefundAmountExcessBalance, "", errors.New(fmt.Sprintf("退款金额超出可退金额|交易金额:%v|退款金额:%v|已退金额:%v", transAmount, refundAmount, totalAmount))
		} else if subRet == 0 {
			//金额退完
			return ss_err.Success, constants.BusinessOrderStatusRefund, nil
		} else {
			//金额有剩余
			return ss_err.Success, constants.BusinessOrderStatusRebatesRefund, nil
		}
	}

	return ss_err.ParamErr, "", errors.New("refundAmount参数为空")
}
