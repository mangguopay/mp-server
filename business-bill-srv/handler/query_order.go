package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/common"
	"a.a/mp-server/common/cache"

	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
)

// 查询订单详情
func (b *BusinessBillHandler) Query(ctx context.Context, req *businessBillProto.QueryRequest, reply *businessBillProto.QueryReply) error {
	if req.AppId == "" {
		ss_log.Error("AppId参数为空")
		reply.ResultCode = ss_err.AppIdIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.AppIdIsEmpty, req.Lang)
		return nil
	}

	if req.OrderNo == "" && req.OutOrderNo == "" {
		ss_log.Error("OrderNo和OutOrderNo两个参数为空")
		reply.ResultCode = ss_err.OrderNoIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.OrderNoIsEmpty, req.Lang)
		return nil
	}

	order, err := dao.BusinessBillDaoInst.AppQueryOrder(req.AppId, req.OrderNo, req.OutOrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("Query查询订单不存在,appId:%v, orderNo=%v, outOrderNo=%v, err:%s", req.AppId, req.OrderNo, req.OutOrderNo, err.Error())
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}

		ss_log.Error("Query查询订单信息失败,appId:%v, orderNo=%v, outOrderNo=%v, err:%s", req.AppId, req.OrderNo, req.OutOrderNo, err.Error())
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OutOrderNo = order.OutOrderNo
	reply.OrderNo = order.OrderNo
	reply.PayAccount = order.PayAccount
	reply.OrderStatus = order.OrderStatus
	reply.Amount = order.Amount
	reply.CurrencyType = order.CurrencyType
	reply.CreateTime = fmt.Sprintf("%v", ss_time.ParseTimeFromPostgres(order.CreateTime, global.Tz).Unix())
	reply.Subject = order.Subject
	reply.Remark = order.Remark
	reply.Rate = order.Rate
	reply.Fee = order.Fee
	if order.PayTime != "" {
		reply.PayTime = fmt.Sprintf("%v", ss_time.ParseTimeFromPostgres(order.PayTime, global.Tz).Unix())
	}

	return nil
}

//根据订单号查询订单
func (b *BusinessBillHandler) QueryOrder(ctx context.Context, req *businessBillProto.QueryOrderRequest, reply *businessBillProto.QueryOrderReply) error {
	// 0.检查
	if req.OrderNo == "" && req.OutOrderNo == "" {
		ss_log.Error("OrderNo和OutOrderNo参数都为空")
		reply.ResultCode = ss_err.OrderNoIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.OrderNoIsEmpty, req.Lang)
		return nil
	}

	order, err := dao.BusinessBillDaoInst.GetOrderInfo(req.OrderNo, req.OutOrderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单不存在,orderNo=%v, outOrderNo=%v", req.OrderNo, req.OutOrderNo)
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询订单信息失败,orderNo=%v ,err: %s", req.OrderNo, err.Error())
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = order.OrderNo
	reply.OrderStatus = order.OrderStatus
	reply.Amount = order.Amount
	reply.CurrencyType = order.CurrencyType
	reply.Subject = order.Subject
	reply.AppName = order.AppName

	return nil
}

//根据二维码id查询订单交易信息(金额, 订单标题, 商家名称)
func (b *BusinessBillHandler) QueryTransInfo(ctx context.Context, req *businessBillProto.QueryTransInfoRequest, reply *businessBillProto.QueryTransInfoReply) error {
	// 0.检查
	if req.QrCodeId == "" {
		ss_log.Error("QrCodeId参数为空")
		reply.ResultCode = ss_err.QrCodeIdIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.QrCodeIdIsEmpty, req.Lang)
		return nil
	}

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

	order, err := dao.BusinessBillDaoInst.GetOrderInfoByOrderNo(orderNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单不存在,orderNo=%v", orderNo)
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询订单信息失败,orderNo=%v ,err: %s", orderNo, err.Error())
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//订单状态为已支付
	if order.OrderStatus == constants.BusinessOrderStatusPay {
		ss_log.Error("订单已支付，orderNo=%v", order.OrderNo)
		reply.ResultCode = ss_err.OrderPaid
		reply.Msg = ss_err.GetMsg(ss_err.OrderPaid, req.Lang)
		return nil
	}

	//订单状态为已退款
	if order.OrderStatus == constants.BusinessOrderStatusRefund {
		ss_log.Error("订单已退款，orderNo=%v", order.OrderNo)
		reply.ResultCode = ss_err.OrderFullRefund
		reply.Msg = ss_err.GetMsg(ss_err.OrderFullRefund, req.Lang)
		return nil
	}

	//订单状态为已超时
	if order.OrderStatus == constants.BusinessOrderStatusPayTimeOut {
		ss_log.Error("订单已过期,qr:%v", req.QrCodeId)
		reply.ResultCode = ss_err.QrCodeExpired
		reply.Msg = ss_err.GetMsg(ss_err.QrCodeExpired, req.Lang)
		return nil
	}

	//不满足以上状态且又不是待支付状态，返回未知
	if order.OrderStatus != constants.BusinessOrderStatusPending {
		ss_log.Error("订单状态未知，orderNo=%v, status=%v", order.OrderNo, order.OrderStatus)
		reply.ResultCode = ss_err.OrderStatusUnknown
		reply.Msg = ss_err.GetMsg(ss_err.OrderStatusUnknown, req.Lang)
		return nil
	}

	//检查订单是否已过期(过期时间小于当前时间)
	if order.ExpireTime < ss_time.NowTimestamp(global.Tz) {
		//修更新订单为已超时
		updateErr := dao.BusinessBillDaoInst.UpdateOrderOutTimeById(order.OrderNo)
		if updateErr != nil {
			ss_log.Error("修改订单为已过期失败,OrderNo:%v, err:%v", order.OrderNo, updateErr)
		}
		ss_log.Error("订单已过期,OrderNo:%v", order.OrderNo)
		reply.ResultCode = ss_err.QrCodeExpired
		reply.Msg = ss_err.GetMsg(ss_err.QrCodeExpired, req.Lang)
		return nil
	}

	// 查询商户名称
	_, simplifyName, err := dao.BusinessDaoInst.QueryNameByBusinessNo(order.BusinessNo)
	if err != nil {
		ss_log.Error("查询商户名称失败,businessNo=%v, err: %s", order.BusinessNo, err.Error())
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.Amount = order.Amount
	reply.Subject = order.Subject
	reply.BusinessName = simplifyName
	reply.CurrencyType = order.CurrencyType
	return nil
}

//生成用户付款码
func (b *BusinessBillHandler) GetPaymentCode(ctx context.Context, req *businessBillProto.GetPaymentCodeRequest, reply *businessBillProto.GetPaymentCodeReply) error {
	if req.AccountNo == "" {
		ss_log.Error("AccountNo参数为空")
		reply.ResultCode = ss_err.AccountNoIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.AccountNoIsEmpty, req.Lang)
		return nil
	}

	//查询用户
	_, err := dao.CustDaoInst.QueryCustNo(req.AccountNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("用户不存在，AccountNo=%v, err=%v", req.AccountNo, err)
			reply.ResultCode = ss_err.AccountNoNotExist
			reply.Msg = ss_err.GetMsg(ss_err.AccountNoNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询用户失败，AccountNo=%v, err=%v", req.AccountNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//生成付款码
	paymentCode := strext.GetDailyId()
	redisKey := constants.GetCustPaymentCodeRedisKey(paymentCode)
	//付款码过期时间暂设为150秒，前端目前是每60秒请求刷新一次
	err = cache.RedisClient.SetNX(redisKey, req.AccountNo, 150*time.Second).Err()
	if err != nil {
		ss_log.Error("刷新用户付款码失败，key=%v, err=%v", redisKey, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}
	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.PaymentCode = paymentCode
	return nil
}

//根据付款码查询用户账号id
func GetAccountNoByCode(paymentCode string) (string, error) {
	if paymentCode == "" {
		return "", errors.New("paymentCode参数为空")
	}

	//根据付款码获取付款人账号
	redisKey := constants.GetCustPaymentCodeRedisKey(paymentCode)
	accountNo, err := cache.RedisClient.Get(redisKey).Result()
	if err != nil {
		if err.Error() == common.RedisValueNilErr.Error() {
			return "", errors.New("付款码已过期")
		}
		return "", errors.New(fmt.Sprintf("redis查询付款人账号失败，PaymentCode=%v, err=%v", paymentCode, err))
	}

	return accountNo, nil

}

//查询用户待支付订单
func (b *BusinessBillHandler) QueryPendingPayOrder(ctx context.Context, req *businessBillProto.QueryPendingPayOrderRequest, reply *businessBillProto.QueryPendingPayOrderReply) error {
	if req.AccountNo == "" {
		ss_log.Error("AccountNo参数为空")
		reply.ResultCode = ss_err.AccountNoIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.AccountNoIsEmpty, req.Lang)
		return nil
	}

	if req.PaymentCode == "" {
		ss_log.Error("PaymentCode参数为空")
		reply.ResultCode = ss_err.PaymentCodeIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.PaymentCodeIsEmpty, req.Lang)
		return nil
	}

	var order *dao.PendingPayOrder
	var err error
	ss_log.Info("查询用户待支付订单---------------------------------start")
	for i := 0; i < 3; i++ {
		ss_log.Info("第%v次查询", i+1)
		//目前是一个付款码码一个订单，付款码即是订单号
		order, err = dao.BusinessBillDaoInst.GetCustPendingPayOrder(req.AccountNo, req.PaymentCode)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("用户没有对应的待支付订单, AccountNo=%v, PaymentCode=%v", req.AccountNo, req.PaymentCode)
				time.Sleep(1000 * time.Millisecond)
			} else {
				ss_log.Error("查询用户待支付订单失败, AccountNo=%v, PaymentCode=%v, err=%v", req.AccountNo, req.PaymentCode, err)
				reply.ResultCode = ss_err.SystemErr
				reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
				return nil
			}
		}
		if order != nil {
			continue
		}
	}
	ss_log.Info("查询用户待支付订单--------------------------------------end")
	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	if order != nil {
		reply.OrderNo = order.OrderNo
		reply.Amount = order.Amount
		reply.CurrencyType = order.CurrencyType
		reply.Subject = order.Subject
		reply.BusinessName = order.SimplifyName
	}
	return nil
}

//退款查询
func (b *BusinessBillHandler) QueryRefundOrder(ctx context.Context, req *businessBillProto.QueryRefundOrderRequest, reply *businessBillProto.QueryRefundOrderReply) error {

	if req.AppId == "" {
		reply.ResultCode = ss_err.AppIdIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.AppIdIsEmpty, req.Lang)
		return nil
	}

	if req.RefundNo == "" && req.OutRefundNo == "" {
		reply.ResultCode = ss_err.ParamErr
		reply.Msg = ss_err.GetMsg(ss_err.ParamErr, req.Lang)
		ss_log.Error("退款查询参数有误：%v", strext.ToJson(req))
		return nil
	}

	//退款详情
	refund, err := dao.BusinessRefundOrderDaoInst.GetRefundByRefundNo(req.AppId, req.RefundNo, req.OutRefundNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("订单退款记录不存在, refundNo=%v, outRefundNo=%v", req.RefundNo, req.OutRefundNo)
			reply.ResultCode = ss_err.RefundNoNotExist
			reply.Msg = ss_err.GetMsg(ss_err.RefundNoNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询订单退款记录失败, refundNo=%v, outRefundNo=%v, err=%v", req.RefundNo, req.OutRefundNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//交易订单详情
	order, err := dao.BusinessBillDaoInst.GetOrderInfo(refund.PayOrderNo, "")
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("交易订单不存在，OrderNo=%v", refund.PayOrderNo)
			reply.ResultCode = ss_err.OrderNotExist
			reply.Msg = ss_err.GetMsg(ss_err.OrderNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询交易订单信息失败，OrderNo=%v, err=%v", refund.PayOrderNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	//查询收款人
	payeeAcc, err := dao.AccountDaoInst.GetAccountById(order.AccountNo)
	if err != nil {
		ss_log.Error("查询退款接收账号失败，AccountNo=%v, err=%v", order.AccountNo, err)
		reply.ResultCode = ss_err.PayeeAccountNoNotExist
		reply.Msg = ss_err.GetMsg(ss_err.PayeeAccountNoNotExist, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.OrderNo = order.OrderNo
	reply.OutOrderNo = order.OutOrderNo
	reply.CurrencyType = order.CurrencyType
	reply.RefundNo = refund.RefundNo
	reply.OutRefundNo = refund.OutRefundNo
	reply.Amount = order.Amount
	reply.RefundAmount = refund.Amount
	reply.RefundStatus = refund.RefundStatus
	//reply.RefundTime = refund.FinishTime
	reply.RefundTime = fmt.Sprintf("%v", ss_time.ParseTimeFromPostgres(refund.FinishTime, global.Tz).Unix())

	reply.RefundPayeeAcc = payeeAcc

	return nil
}

//转账查询
func (b *BusinessBillHandler) QueryTransfer(ctx context.Context, req *businessBillProto.QueryTransferRequest, reply *businessBillProto.QueryTransferReply) error {
	if req.AppId == "" {
		reply.ResultCode = ss_err.AppIdIsEmpty
		reply.Msg = ss_err.GetMsg(ss_err.AppIdIsEmpty, req.Lang)
		return nil
	}

	if req.TransferNo == "" && req.OutTransferNo == "" {
		reply.ResultCode = ss_err.ParamErr
		reply.Msg = ss_err.GetMsg(ss_err.ParamErr, req.Lang)
		ss_log.Error("转账查询参数有误：%v", strext.ToJson(req))
		return nil
	}

	transfer, err := dao.BusinessTransferOrderDaoInst.GetOrderByTransferNo(req.AppId, req.TransferNo, req.OutTransferNo)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("转账订单不存在, transferNo=%v, outTransferNo=%v", req.TransferNo, req.OutTransferNo)
			reply.ResultCode = ss_err.TransferNoNotExist
			reply.Msg = ss_err.GetMsg(ss_err.TransferNoNotExist, req.Lang)
			return nil
		}
		ss_log.Error("查询转账订单失败, transferNo=%v, outTransferNo=%v, err=%v", req.TransferNo, req.OutTransferNo, err)
		reply.ResultCode = ss_err.SystemErr
		reply.Msg = ss_err.GetMsg(ss_err.SystemErr, req.Lang)
		return nil
	}

	reply.ResultCode = ss_err.Success
	reply.Msg = ss_err.GetMsg(ss_err.Success, req.Lang)
	reply.Order = &businessBillProto.EnterpriseTransferOrder{
		TransferNo:     transfer.LogNo,
		OutTransferNo:  transfer.OutTransferNo,
		Amount:         transfer.Amount,
		CurrencyType:   transfer.CurrencyType,
		TransferStatus: transfer.OrderStatus,
		//TransferTime:   transfer.FinishTime,
		TransferTime: fmt.Sprintf("%v", ss_time.ParseTimeFromPostgres(transfer.FinishTime, global.Tz).Unix()),
		WrongReason:  transfer.WrongReason,
	}
	return nil
}
