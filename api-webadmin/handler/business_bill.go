package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

type BusinessBillHandler struct {
	Client businessBillProto.BusinessBillService
}

var BusinessBillHandlerInst BusinessBillHandler

//手动结算上游订单
func (*BusinessBillHandler) ManualSettle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			orderNoArrStr := container.GetValFromMapMaybe(params, "order_no").ToStringNoPoint()
			orderNoArr := strext.Json2StrList(orderNoArrStr)
			if orderNoArr == nil {
				ss_log.Error("参数[order_no_array]不是一个json字符串，orderNoArrStr=[%v]", orderNoArrStr)
				return ss_err.ERR_PARAM, "", nil
			}
			req := &businessBillProto.ManualSettleRequest{
				OrderNos: orderNoArr,
			}
			reply, err := BusinessBillHandlerInst.Client.ManualSettle(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用business-bill-srv.ManualSettle()失败，req=[%v], err=[%v]", strext.ToJson(req), err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			var list []gin.H
			if ss_err.PayRetCode(reply.ResultCode) == ss_err.ERR_SUCCESS {
				for k, v := range reply.FailOrder {
					failOrder := gin.H{
						"order_no":   k,
						"fail_cause": v,
					}
					list = append(list, failOrder)
				}
			}

			return ss_err.PayRetCode(reply.ResultCode), gin.H{
				"fail_total": reply.FailNum,
				"fail_order": list,
			}, err
		})
	}
}

//获取商家Modenpay交易订单列表
func (*CustHandler) GetModPayBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessBills(context.TODO(), &custProto.GetBusinessBillsRequest{
				BusinessId:      strext.ToStringNoPoint(params[0]),
				BusinessName:    strext.ToStringNoPoint(params[1]),
				SceneNo:         strext.ToStringNoPoint(params[2]),
				OrderStatus:     strext.ToStringNoPoint(params[3]), //订单状态
				IsSettled:       strext.ToStringNoPoint(params[4]), //是否已结算0-未结算1-已结算
				CurrencyType:    strext.ToStringNoPoint(params[5]), //币种
				Page:            strext.ToInt32(params[6]),
				PageSize:        strext.ToInt32(params[7]),
				StartTime:       strext.ToStringNoPoint(params[8]),
				EndTime:         strext.ToStringNoPoint(params[9]),
				BusinessAccount: strext.ToStringNoPoint(params[10]),
				Subject:         strext.ToStringNoPoint(params[11]),
				OrderNo:         strext.ToStringNoPoint(params[12]),
				SettleId:        strext.ToStringNoPoint(params[13]),
				ChannelType:     constants.ChannelTypeInner,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, 0, nil
			}

			/**
			商家中心也调用了这个查询接口，但两者返回字段有不同
			不建议将全部字段返回出去
			*/
			var list []interface{}
			for _, v := range reply.Datas {
				data := gin.H{
					"out_order_no":     v.OutOrderNo,
					"order_no":         v.OrderNo,
					"business_id":      v.BusinessId,
					"business_name":    v.BusinessName,
					"business_account": v.BusinessAccount,
					"subject":          v.Subject,
					"currency_type":    v.CurrencyType,
					"amount":           v.Amount,
					"scene_name":       v.SceneName,
					"create_time":      v.CreateTime,
					"settle_date":      v.SettleDate,
					"cycle":            v.Cycle,
					"rate":             v.Rate,
					"real_amount":      v.RealAmount,
					"fee":              v.Fee,
					"order_status":     v.OrderStatus,
					"app_name":         v.AppName,
				}
				list = append(list, data)
			}

			return reply.ResultCode, list, reply.Total, nil
		}, "business_id", "business_name", "scene_no", "order_status", "is_settled", "currency_type",
			"page", "page_size", "start_time", "end_time", "business_account", "subject", "order_no", "settle_id")
	}
}

//获取商家Modenpay退款订单列表
func (*CustHandler) GetModRefundBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetRefundBills(context.TODO(), &custProto.GetRefundBillsRequest{
				BusinessId:   strext.ToStringNoPoint(params[0]),
				BusinessName: strext.ToStringNoPoint(params[1]),
				SceneNo:      strext.ToStringNoPoint(params[2]),
				OrderStatus:  strext.ToStringNoPoint(params[3]), //订单状态
				CurrencyType: strext.ToStringNoPoint(params[4]), //币种
				StartTime:    strext.ToStringNoPoint(params[5]),
				EndTime:      strext.ToStringNoPoint(params[6]),
				Page:         strext.ToStringNoPoint(params[7]),
				PageSize:     strext.ToStringNoPoint(params[8]),
				PayeeAccount: strext.ToStringNoPoint(params[9]),
				RefundNo:     strext.ToStringNoPoint(params[10]),
				ChannelType:  constants.ChannelTypeInner,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, 0, nil
			}

			/**
			商家中心也调用了这个查询接口，但两者返回字段有不同
			不建议将全部字段返回出去
			*/
			var list []interface{}
			for _, v := range reply.List {
				data := gin.H{
					"refund_no":      v.RefundNo,
					"out_refund_no":  v.OutRefundNo,
					"business_id":    v.BusinessId,
					"business_name":  v.BusinessName,
					"order_status":   v.OrderStatus,
					"trans_order_no": v.TransOrderNo,
					"amount":         v.Amount,
					"currency_type":  v.CurrencyType,
					"create_time":    v.CreateTime,
					"finish_time":    v.FinishTime,
					"subject":        v.Subject,
					"app_name":       v.AppName,
					"scene_name":     v.SceneName,
					"payee_account":  v.PayeeAccount,
				}
				list = append(list, data)
			}

			return reply.ResultCode, list, reply.Total, nil
		}, "business_id", "business_name", "scene_no", "order_status", "currency_type",
			"start_time", "end_time", "page", "page_size", "payee_account", "refund_no")
	}
}

//上游渠道交易订单
func (*CustHandler) GetChannelBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetChannelBills(context.TODO(), &custProto.GetChannelBillsRequest{
				Account:      strext.ToStringNoPoint(params[0]),
				OrderStatus:  strext.ToStringNoPoint(params[1]), //订单状态
				SceneNo:      strext.ToStringNoPoint(params[2]),
				CurrencyType: strext.ToStringNoPoint(params[3]), //币种
				IsSettled:    strext.ToStringNoPoint(params[4]),
				StartTime:    strext.ToStringNoPoint(params[5]),
				EndTime:      strext.ToStringNoPoint(params[6]),
				Page:         strext.ToInt32(params[7]),
				PageSize:     strext.ToInt32(params[8]),
				ChannelNo:    strext.ToStringNoPoint(params[9]),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, 0, nil
			}

			var list []interface{}
			for _, v := range reply.List {
				data := gin.H{
					"business_account": v.BusinessAccount,
					"business_id":      v.BusinessId,
					"business_name":    v.BusinessName,
					"out_order_no":     v.OutOrderNo,
					"order_no":         v.OrderNo,
					"channel_name":     v.ChannelName,
					"channel_rate":     v.ChannelRate,
					"scene_name":       v.SceneName,
					"rate":             v.Rate,
					"cycle":            v.Cycle,
					"pay_time":         v.PayTime,
					"amount":           v.Amount,
					"real_amount":      v.RealAmount,
					"currency_type":    v.CurrencyType,
					"order_status":     v.OrderStatus,
				}
				list = append(list, data)
			}

			return reply.ResultCode, list, reply.Total, nil
		}, "business_account", "order_status", "scene_no", "currency_type", "is_settled",
			"start_time", "end_time", "page", "page_size", "channel_no")
	}
}

//上游渠道交易退款订单列表
func (*CustHandler) GetChannelRefundBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {

			reply, err := CustHandlerInst.Client.GetChannelRefundBills(context.TODO(), &custProto.GetChannelRefundBillsRequest{
				BusinessAccount: strext.ToStringNoPoint(params[0]),
				BusinessId:      strext.ToStringNoPoint(params[1]),
				RefundNo:        strext.ToStringNoPoint(params[2]),
				OutRefundNo:     strext.ToStringNoPoint(params[3]),
				RefundStatus:    strext.ToStringNoPoint(params[4]),
				CurrencyType:    strext.ToStringNoPoint(params[5]),
				StartTime:       strext.ToStringNoPoint(params[6]),
				EndTime:         strext.ToStringNoPoint(params[7]),
				Page:            strext.ToInt32(params[8]),
				PageSize:        strext.ToInt32(params[9]),
				SceneName:       strext.ToStringNoPoint(params[10]),
				ChannelNo:       strext.ToStringNoPoint(params[11]),
			})
			if err != nil {
				ss_log.Error("调用cust-srv.GetChannelRefundBills()失败， err=%v", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.List, reply.Total, nil
		}, "business_account", "business_id", "refund_no", "out_refund_no", "refund_status", "currency_type",
			"start_time", "end_time", "page", "page_size", "scene_name", "channel_no")
	}
}

//订单详情
func (*CustHandler) GetBillsDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessBillDetailRequest{
				OrderNo: strext.ToStringNoPoint(params[0]),
			}
			reply, err := CustHandlerInst.Client.GetBusinessBillDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.GetBusinessBillDetail()失败, err=%v", err)
				return ss_err.ERR_SYSTEM, nil, 0, nil
			}

			return reply.ResultCode, reply.Data, 1, nil
		}, "order_no")
	}
}

//退款订单详情
func (*CustHandler) GetRefundBillsDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetRefundDetailRequest{
				RefundNo: strext.ToStringNoPoint(params[0]),
			}
			reply, err := CustHandlerInst.Client.GetRefundDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.GetBusinessBillDetail()失败, err=%v", err)
				return ss_err.ERR_SYSTEM, nil, 0, nil
			}

			return reply.ResultCode, reply.Data, 1, nil
		}, "refund_no")
	}
}
