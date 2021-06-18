package handler

import (
	"context"

	"a.a/cu/strext"

	"a.a/cu/ss_log"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/api-pay/inner_util"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type BusinessBillHandler struct {
	Client businessBillProto.BusinessBillService
}

var BusinessBillHandlerInst BusinessBillHandler

// 下单接口
func (BusinessBillHandler) Prepay() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		req := &businessBillProto.PrepayRequest{
			Amount:       inner_util.M(c, "amount"),        // 订单金额
			Remark:       inner_util.M(c, "remark"),        // 备注信息
			NotifyUrl:    inner_util.M(c, "notify_url"),    // 异步通知地址
			ReturnUrl:    inner_util.M(c, "return_url"),    // 同步跳转地址
			AppId:        inner_util.M(c, "app_id"),        // 应用id
			OutOrderNo:   inner_util.M(c, "out_order_no"),  // 外部订单号(商户系统内部的订单号)
			CurrencyType: inner_util.M(c, "currency_type"), // 货币类型
			Subject:      inner_util.M(c, "subject"),       // 商品的标题/交易标题/订单标题/订单关键字等
			TimeExpire:   inner_util.M(c, "time_expire"),   // 订单过期时间(时间戳)
			PaymentCode:  inner_util.M(c, "payment_code"),  // 付款码, 不为空时场景为商家扫用户付款码
			TradeType:    inner_util.M(c, "trade_type"),    // 交易类型
			Lang:         inner_util.M(c, "lang"),          // 语言类型
			Attach:       inner_util.M(c, "attach"),        // 原样返回字段
		}

		reply, err := BusinessBillHandlerInst.Client.Prepay(context.TODO(), req)
		if err != nil {
			ss_log.Error("%s|请求下单出错,err:%v,req:%v", traceNo, err, strext.ToJson(req))
			c.Set(common.RET_CODE, ss_err.ACErrSysErr)
			return
		}

		ss_log.Error("%s|下单成功|OrderNo:%v|OutOrderNo:%v", traceNo, reply.OrderNo, reply.OutOrderNo)

		c.Set(common.RET_CODE, ss_err.ACErrSuccess)
		c.Set(common.RET_DATA, gin.H{
			"order_no":        reply.OrderNo,
			"out_order_no":    reply.OutOrderNo,
			"qr_code":         reply.QrCodeId,
			"app_pay_content": reply.AppPayContent,

			"sub_code": reply.ResultCode,
			"sub_msg":  reply.Msg,
		})
		return
	}
}

//订单结果查询
func (BusinessBillHandler) Query() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================

		req := &businessBillProto.QueryRequest{
			AppId:      inner_util.M(c, "app_id"), // 应用id
			OrderNo:    inner_util.M(c, "order_no"),
			OutOrderNo: inner_util.M(c, "out_order_no"),
			Lang:       inner_util.M(c, "lang"), // 语言类型
		}

		reply, err := BusinessBillHandlerInst.Client.Query(context.TODO(), req)
		if err != nil || reply == nil || reply.String() == "" {
			ss_log.Info("reply=[%v],err=[%v]", reply, err)
			c.Set(common.RET_CODE, ss_err.ACErrSysBusy)
			return
		}

		c.Set(common.RET_CODE, ss_err.ACErrSuccess)
		c.Set(common.RET_DATA, gin.H{
			"out_order_no": reply.OutOrderNo,
			"order_no":     reply.OrderNo,
			//"pay_account":   reply.PayAccount,
			"order_status":  reply.OrderStatus,
			"amount":        reply.Amount,
			"currency_type": reply.CurrencyType,
			"create_time":   reply.CreateTime,
			"pay_time":      reply.PayTime,
			"subject":       reply.Subject,
			"remark":        reply.Remark,
			"rate":          reply.Rate,
			"fee":           reply.Fee,

			"sub_code": reply.ResultCode,
			"sub_msg":  reply.Msg,
		})
		return
	}
}

//订单退款
func (BusinessBillHandler) Refund() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		req := &businessBillProto.ApiPayRefundRequest{
			AppId:        inner_util.M(c, "app_id"),   // 账号类型
			OrderNo:      inner_util.M(c, "order_no"), // 平台订单号
			OutOrderNo:   inner_util.M(c, "out_order_no"),
			OutRefundNo:  inner_util.M(c, "out_refund_no"), // 外部退款单号
			RefundAmount: inner_util.M(c, "refund_amount"), // 退款金额
			RefundReason: inner_util.M(c, "refund_reason"), // 备注信息
			NotifyUrl:    inner_util.M(c, "notify_url"),    // 异步通知地址
			Lang:         inner_util.M(c, "lang"),          // 语言类型
			Attach:       inner_util.M(c, "attach"),        // 原样返回字段
		}

		reply, err := BusinessBillHandlerInst.Client.ApiPayRefund(context.TODO(), req)
		if err != nil {
			ss_log.Error("%s|请求退款出错,err:%v,req:%v", traceNo, err, strext.ToJson(req))
			c.Set(common.RET_CODE, ss_err.ACErrSysErr)
			return
		}

		ss_log.Error("%s|退款单成功|OrderNo:%v", traceNo, reply.RefundNo)

		c.Set(common.RET_CODE, ss_err.ACErrSuccess)
		c.Set(common.RET_DATA, gin.H{
			"order_no":      reply.OrderNo,
			"out_order_no":  reply.OutOrderNo,
			"refund_no":     reply.RefundNo,
			"out_refund_no": reply.OutRefundNo,
			"refund_status": reply.RefundStatus,
			"refund_amount": reply.RefundAmount,

			"sub_code": reply.ResultCode,
			"sub_msg":  reply.Msg,
		})
		return
	}
}

//退款查询
func (BusinessBillHandler) QueryRefund() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		req := &businessBillProto.QueryRefundOrderRequest{
			AppId:       inner_util.M(c, "app_id"),        // 应用id
			RefundNo:    inner_util.M(c, "refund_no"),     // 退款id
			OutRefundNo: inner_util.M(c, "out_refund_no"), // 外部退款id
			Lang:        inner_util.M(c, "lang"),          // 语言类型
		}

		reply, err := BusinessBillHandlerInst.Client.QueryRefundOrder(context.TODO(), req)
		if err != nil {
			ss_log.Error("%s|请求退款查询出错,err:%v,req:%v", traceNo, err, strext.ToJson(req))
			c.Set(common.RET_CODE, ss_err.ACErrSysErr)
			return
		}

		ss_log.Error("%s|退款单成功|OrderNo:%v", traceNo, reply.OrderNo)

		c.Set(common.RET_CODE, ss_err.ACErrSuccess)
		c.Set(common.RET_DATA, gin.H{
			"order_no":      reply.OrderNo,
			"out_order_no":  reply.OutOrderNo,
			"currency_type": reply.CurrencyType,
			"refund_no":     reply.RefundNo,
			"out_refund_no": reply.OutRefundNo,
			"trans_amount":  reply.Amount,       // 订单金额
			"refund_amount": reply.RefundAmount, // 退款金额
			"refund_status": reply.RefundStatus,
			"refund_time":   reply.RefundTime,

			"sub_code": reply.ResultCode,
			"sub_msg":  reply.Msg,
		})
		return
	}
}
