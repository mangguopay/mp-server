package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/merchant-mock/conf"
	"a.a/mp-server/merchant-mock/dao"
	"a.a/mp-server/merchant-mock/pay"
	"github.com/gin-gonic/gin"
)

const (
	//通知类型
	NotifyTypeToPayment  = "PAYMENT"
	NotifyTypeToRefund   = "REFUND"
	NotifyTypeToTransfer = "TRANSFER"

	NotifyResponseSuccess = "success"
)

const (
	RefundStatusSuccess = "1" // 1-成功
	RefundStatusFail    = "2" // 2-失败
)

const (
	TransferStatusSuccess = "1" // 1-成功
	TransferStatusFail    = "2" // 2-失败
)

// 异步回调处理
func OrderNotify(c *gin.Context) {

	log.Printf("收到异步通知------start-----")
	//c.String(http.StatusOK, "请求失败")
	//return

	// 只处理post
	if http.MethodPost != c.Request.Method {
		c.String(http.StatusOK, "不是post请求")
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err != nil {
		c.String(http.StatusOK, "读取body失败,err:"+err.Error())
		return
	}

	if len(body) == 0 {
		c.String(http.StatusOK, "请求body为空")
		return
	}

	log.Printf("body:%v\n", string(body))

	reqMap := strext.FormString2Map(string(body))
	if reqMap == nil {
		log.Printf("解析body内容失败, body:%v", string(body))
		c.String(http.StatusOK, "解析body内容失败")
		return
	}

	notifyType, exists := reqMap["notify_type"]
	if !exists {
		log.Printf("notify_type字段不存在, body:%v", string(body))
		c.String(http.StatusOK, "notify_type字段不存在")
		return
	}

	log.Printf("notifyType:%+v\n", notifyType)

	retStr := NotifyResponseSuccess

	switch notifyType {
	case NotifyTypeToPayment: // 支付
		retStr = NotifyTypeToPaymentHandler(body)
	case NotifyTypeToRefund: // 退款
		retStr = NotifyTypeToRefundHandler(body)
	case NotifyTypeToTransfer: // 企业付款
		retStr = NotifyTypeToTransferHandler(body)
	default:
		retStr = "通知类型错误"
	}
	log.Printf("收到异步通知------end-----")

	c.String(http.StatusOK, retStr)
	return
}

// 异步通知处理-支付
func NotifyTypeToPaymentHandler(body []byte) string {
	ret, err := pay.NotifyPaymentVerify(body)
	if err != nil {
		log.Printf("验证异步通知数据失败")
		return "验证异步通知数据失败"
	}

	if ret.AppId != conf.AppId {
		log.Printf("AppId不匹配")
		return "AppId不匹配"
	}

	order, err := dao.OrderInstance.GetOneByOrderSn(ret.OutOrderNo)
	if err != nil {
		log.Printf("查询订单失败:%v", err)
		return "查询订单失败"
	}

	if ret.OrderNo != order.PayOrderSn {
		log.Printf("订单的支付号不匹配")
		return "订单的支付号不匹配"
	}

	if ret.CurrencyType != order.CurrencyType {
		log.Printf("订单的币种不匹配")
		return "订单的支付号不匹配"
	}

	if strext.ToInt64(ret.Amount) != order.Amount {
		log.Printf("订单的金额不匹配")
		return "订单的金额不匹配"
	}

	if order.Status == dao.OrderStatusPaid {
		log.Printf("订单已支付成功,重复通知")
		return NotifyResponseSuccess
	}

	if order.Status != dao.OrderStatusPaying {
		log.Printf("订单的状态不允许")
		return "订单的状态不允许"
	}

	timestamp, perr := strconv.ParseInt(ret.PayTime, 10, 64)
	if perr != nil {
		log.Printf("解析时间参数失败, payTime:%v, err:%v \n", ret.PayTime, perr)
		return "解析时间参数失败"
	}

	// 转换时间格式
	payTimeStr := ss_time.ForPostgres(time.Unix(timestamp, 0))

	// 更新订单为已支付
	if err := dao.OrderInstance.UpdatePaidOk(order.OrderSn, payTimeStr, ret.PayAccount); err != nil {
		log.Printf("更新订单支付成功失败, orderNo:%v, payTime:%v, payAccount:%v, err:%v", order.OrderSn, ret.PayTime, ret.PayAccount, err)
		return "更新订单支付成功失败"
	}

	return NotifyResponseSuccess
}

// 异步通知处理-退款
func NotifyTypeToRefundHandler(body []byte) string {
	ret, err := pay.NotifyRefundVerify(body)
	if err != nil {
		log.Printf("验证异步通知数据失败")
		return "验证异步通知数据失败"
	}

	if ret.AppId != conf.AppId {
		log.Printf("AppId不匹配")

		return "AppId不匹配"
	}

	refund, err := dao.RefundInstance.GetOneByOutRefundNo(ret.OutRefundNo)
	if err != nil {
		log.Printf("查询退款记录失败:%v", err)
		return "查询退款记录失败"
	}

	if ret.OutRefundNo != refund.OutRefundNo {
		log.Printf("退款转账单号不匹配")
		return "退款转账单号不匹配"
	}

	if ret.CurrencyType != refund.CurrencyType {
		log.Printf("退款订单的币种不匹配")
		return "退款订单的支付号不匹配"
	}

	if strext.ToInt64(ret.RefundAmount) != refund.Amount {
		log.Printf("退款订单的金额不匹配")
		return "退款订单的金额不匹配"
	}

	if refund.Status == dao.TransferStatusSuccess {
		log.Printf("订单已退款成功,重复通知")
		return NotifyResponseSuccess
	}

	if refund.Status != dao.TransferStatusPending {
		log.Printf("退款订单的状态不允许")
		return "退款订单的状态不允许"
	}

	if ret.RefundStatus != RefundStatusSuccess && ret.RefundStatus != RefundStatusFail {
		log.Printf("退款订单的状态不允许")
		return "退款订单的状态不允许"
	}

	if ret.RefundStatus == RefundStatusSuccess {
		timestamp, perr := strconv.ParseInt(ret.RefundTime, 10, 64)
		if perr != nil {
			log.Printf("解析时间参数失败, payTime:%v, err:%v \n", ret.RefundTime, perr)
			return "解析时间参数失败"
		}

		// 转换时间格式
		refundTimeStr := ss_time.ForPostgres(time.Unix(timestamp, 0))

		// 更新退款成功
		if err := dao.RefundInstance.UpdateRefundSuccess(refund.OutRefundNo, refundTimeStr); err != nil {
			log.Printf("更新退款成功失败, OutRefundNo:%v, RefundTime:%v, err:%v", refund.OutRefundNo, ret.RefundTime, err)
			return "更新退款成功失败"
		}
	} else {
		// 更新退款失败
		if err := dao.RefundInstance.UpdateRefundFail(refund.OutRefundNo); err != nil {
			log.Printf("更新退款失败失败, OutRefundNo:%v, err:%v", refund.OutRefundNo, err)
			return "更新退款失败失败"
		}
	}

	return NotifyResponseSuccess
}

// 异步通知处理-企业付款
func NotifyTypeToTransferHandler(body []byte) string {
	ret, err := pay.NotifyTransferVerify(body)
	if err != nil {
		log.Printf("验证异步通知数据失败")
		return "验证异步通知数据失败"
	}

	if ret.AppId != conf.AppId {
		log.Printf("AppId不匹配")
		return "AppId不匹配"
	}

	transfer, err := dao.TransferInstance.GetOneByOutTransferNo(ret.OutTransferNo)
	if err != nil {
		log.Printf("查询转账记录失败:%v", err)
		return "查询转账记录失败"
	}

	if ret.TransferNo != transfer.TransferNo {
		log.Printf("转账单号不匹配")
		return "转账单号不匹配"
	}

	if ret.CurrencyType != transfer.CurrencyType {
		log.Printf("订单的币种不匹配")
		return "订单的支付号不匹配"
	}

	if strext.ToInt64(ret.Amount) != transfer.Amount {
		log.Printf("订单的金额不匹配")
		return "订单的金额不匹配"
	}

	if transfer.Status == dao.TransferStatusSuccess {
		log.Printf("订单已转账成功,重复通知")
		return NotifyResponseSuccess
	}

	if transfer.Status != dao.TransferStatusPending {
		log.Printf("转账订单的状态不允许")
		return "转账订单的状态不允许"
	}

	if ret.TransferStatus != TransferStatusSuccess && ret.TransferStatus != TransferStatusFail {
		log.Printf("转账订单的状态不允许")
		return "转账订单的状态不允许"
	}

	if ret.TransferStatus == TransferStatusSuccess {
		timestamp, perr := strconv.ParseInt(ret.TransferTime, 10, 64)
		if perr != nil {
			log.Printf("解析时间参数失败, payTime:%v, err:%v \n", ret.TransferTime, perr)
			return "解析时间参数失败"
		}
		// 转换时间格式
		transferTimeStr := ss_time.ForPostgres(time.Unix(timestamp, 0))

		// 更新转账成功
		if err := dao.TransferInstance.UpdateTransferSuccess(transfer.OutTransferNo, transferTimeStr); err != nil {
			log.Printf("更新转账成功失败, OutTransferNo:%v, TransferTime:%v, err:%v", transfer.OutTransferNo, ret.TransferTime, err)
			return "更新转账成功失败"
		}
	} else {
		// 更新转账失败
		if err := dao.TransferInstance.UpdateTransferFail(transfer.OutTransferNo); err != nil {
			log.Printf("更新转账失败失败, OutTransferNo:%v, err:%v", transfer.OutTransferNo, err)
			return "更新转账失败失败"
		}
	}

	return NotifyResponseSuccess
}

// 同步跳转处理
func OrderJumpBack(c *gin.Context) {
	orderSn := c.Query("order_sn")

	if orderSn == "" {
		RedirectTipsError(c, "订单号为空")
		return
	}

	order, err := dao.OrderInstance.GetOneByOrderSn(orderSn)
	if err != nil {
		log.Printf("查询订单失败, orderSn:%v, err:%v", orderSn, err)
		RedirectTipsError(c, "查询订单失败:"+err.Error())
		return
	}

	result := 0
	if order.Status == dao.OrderStatusPaid {
		result = 1
	}

	c.HTML(http.StatusOK, "order/pay_result.html", gin.H{
		"title":  "支付结果",
		"order":  order,
		"result": result,
	})
}
