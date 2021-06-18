package handler

import (
	"context"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/api-pay/inner_util"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type BillHandler struct {
}

/**
企业转账
*/
func (BillHandler) EnterpriseTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		req := &businessBillProto.EnterpriseTransferRequest{
			AppId:            inner_util.M(c, "app_id"),          // 应用id
			OutTransferNo:    inner_util.M(c, "out_transfer_no"), // 外部转账单号
			Amount:           inner_util.M(c, "amount"),          // 订单金额
			CurrencyType:     inner_util.M(c, "currency_type"),   // 货币类型
			PayeeCountryCode: inner_util.M(c, "country_code"),    // 国家码
			PayeePhone:       inner_util.M(c, "payee_phone"),     // 收款人手机号
			PayeeEmail:       inner_util.M(c, "payee_email"),     // 收款人邮箱
			NotifyUrl:        inner_util.M(c, "notify_url"),      // 异步通知地址
			Remark:           inner_util.M(c, "remark"),          // 备注信息
			Lang:             inner_util.M(c, "lang"),            // 语言类型
			Attach:           inner_util.M(c, "attach"),          // 原样返回字段
		}

		reply, err := BusinessBillHandlerInst.Client.EnterpriseTransfer(context.TODO(), req)
		if err != nil {
			ss_log.Error("%s|请求转账出错,err:%v,req:%v", traceNo, err, strext.ToJson(req))
			c.Set(common.RET_CODE, ss_err.ACErrSysErr)
			return
		}

		if reply.Order == nil {
			c.Set(common.RET_CODE, ss_err.ACErrSuccess)
			c.Set(common.RET_DATA, gin.H{
				"sub_code": reply.ResultCode,
				"sub_msg":  reply.Msg,
			})
			return
		}

		order := reply.Order
		c.Set(common.RET_CODE, ss_err.ACErrSuccess)
		c.Set(common.RET_DATA, gin.H{
			"transfer_no":     order.TransferNo,
			"out_transfer_no": order.OutTransferNo,
			"amount":          order.Amount,
			"currency_type":   order.CurrencyType,
			"transfer_status": order.TransferStatus,

			"sub_code": reply.ResultCode,
			"sub_msg":  reply.Msg,
		})
		return
	}
}

/**
转账查询
*/
func (BillHandler) QueryTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		req := &businessBillProto.QueryTransferRequest{
			AppId:         inner_util.M(c, "app_id"),          // 应用id
			TransferNo:    inner_util.M(c, "transfer_no"),     // 平台转账单号
			OutTransferNo: inner_util.M(c, "out_transfer_no"), // 外部转账单号
			Lang:          inner_util.M(c, "lang"),            // 语言类型
		}

		reply, err := BusinessBillHandlerInst.Client.QueryTransfer(context.TODO(), req)
		if err != nil {
			ss_log.Error("%s|请求转账出错,err:%v,req:%v", traceNo, err, strext.ToJson(req))
			c.Set(common.RET_CODE, ss_err.ACErrSysErr)
			return
		}

		if reply.Order == nil {
			c.Set(common.RET_CODE, ss_err.ACErrSuccess)
			c.Set(common.RET_DATA, gin.H{
				"sub_code": reply.ResultCode,
				"sub_msg":  reply.Msg,
			})
			return
		}

		order := reply.Order
		c.Set(common.RET_CODE, ss_err.ACErrSuccess)
		c.Set(common.RET_DATA, gin.H{
			"transfer_no":     order.TransferNo,
			"out_transfer_no": order.OutTransferNo,
			"amount":          order.Amount,
			"currency_type":   order.CurrencyType,
			"transfer_status": order.TransferStatus,
			"transfer_time":   order.TransferTime,
			"wrong_reason":    order.WrongReason,

			"sub_code": reply.ResultCode,
			"sub_msg":  reply.Msg,
		})
		return
	}
}
