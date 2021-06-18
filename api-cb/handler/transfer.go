package handler

import (
	"context"
	"net/http"

	"a.a/cu/ss_log"
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/api-cb/m"
	"a.a/mp-server/api-cb/poly"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"github.com/gin-gonic/gin"
)

/**
 *
 */
func (BusinessBillHandler) TransferCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		//解析上游数据
		supplierCode := c.Param("supplierCode")
		// 找对应的supplierCode
		data, b := c.Get(common.INNER_PARAM_MAP)
		if !b || data == nil {
			ss_log.Error("接受上游报文失败,b: %v,data :%v", b, data)
			return
		}
		resp := poly.PolyWrapperInst.TransferCallback(supplierCode, &m.PolyTransferCallbackReq{
			RecvMap: data.(map[string]interface{}),
		})

		//===============================================

		req := &businessBillProto.TransferCallbackRequest{
			InnerOrderNo: resp.InnerOrderNo,
			Amount:       resp.Amount,
			UpperOrderNo: resp.UpperOrderNo,
			OrderStatus:  resp.OrderStatus,
			UpdateTime:   resp.UpdateTime,
		}

		reply, err := BusinessBillHandlerInst.Client.TransferCallback(context.TODO(), req)
		ss_log.Info("reply=[%v],err=[%v]", reply, err)
		if err != nil {
			ss_log.Error("err=%v", err)
			c.String(http.StatusOK, resp.RetBody)
			return
		}

		// 返回上游特定信息
		c.String(http.StatusOK, resp.RetBody)
		return
	}
}
