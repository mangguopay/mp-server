package handler

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_func"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webbusiness/inner_util"
	"a.a/mp-server/api-webbusiness/verify"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"

	"github.com/gin-gonic/gin"
)

type BusinessBillHandler struct {
	Client businessBillProto.BusinessBillService
}

var BusinessBillHandlerInst BusinessBillHandler

func (*CustHandler) GetBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessBillsRequest{
				BusinessNo:   inner_util.GetJwtDataString(c, "iden_no"),
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				OrderStatus:  strext.ToStringNoPoint(params[2]), //订单状态
				IsSettled:    strext.ToStringNoPoint(params[3]), //是否已结算0-未结算1-已结算
				CurrencyType: strext.ToStringNoPoint(params[4]), //币种
				StartTime:    strext.ToStringNoPoint(params[5]),
				EndTime:      strext.ToStringNoPoint(params[6]),
				OrderNo:      strext.ToStringNoPoint(params[7]),
				SceneNo:      strext.ToStringNoPoint(params[8]),
				Lang:         ss_net.GetCommonData(c).Lang,
			}
			if req.BusinessNo == "" {
				ss_log.Error("BusinessNo参数不能为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.GetBusinessBills(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, 0, nil
			}

			var list []interface{}
			for _, v := range reply.Datas {
				data := gin.H{
					"create_time":   v.CreateTime,
					"order_no":      v.OrderNo,
					"out_order_no":  v.OutOrderNo,
					"subject":       v.Subject,
					"amount":        v.Amount,
					"real_amount":   v.RealAmount,
					"currency_type": v.CurrencyType,
					"fee":           v.Fee,
					"order_status":  v.OrderStatus,
					"settle_id":     v.SettleId,
					"scene_name":    v.SceneName,
					"app_name":      v.AppName,
				}
				list = append(list, data)
			}
			return reply.ResultCode, gin.H{
				"datas":   list,         //订单列表
				"usd_cnt": reply.UsdCnt, //usd订单数量统计
				"usd_sum": reply.UsdSum, //usd订单金额统计
				"khr_cnt": reply.KhrCnt,
				"khr_sum": reply.KhrSum,
			}, reply.Total, nil
		}, "page", "page_size", "order_status", "is_settled", "currency_type", "start_time", "end_time", "order_no", "scene_no")
	}
}

func (*CustHandler) GetBillDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessBillDetailRequest{
				OrderNo: strext.ToStringNoPoint(params[0]),
				Lang:    ss_net.GetCommonData(c).Lang,
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.GetBusinessBillDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			//账号进行脱敏
			reply.Data.Account = ss_func.GetDesensitizationAccount(reply.Data.Account)
			return reply.ResultCode, reply.Data, 0, nil
		}, "order_no")
	}
}

//下载交易订单文件
func (*CustHandler) DownloadBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.CreateBillFile(context.TODO(), &custProto.CreateBillFileRequest{
				Uid:          inner_util.GetJwtDataString(c, "account_uid"),
				IdenNo:       inner_util.GetJwtDataString(c, "iden_no"),
				Page:         0,
				PageSize:     1000,
				OrderStatus:  container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),  //订单状态
				IsSettled:    container.GetValFromMapMaybe(params, "is_settled").ToStringNoPoint(),    //是否已结算0-未结算1-已结算
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(), //币种
				StartTime:    container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:      container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				OrderNo:      container.GetValFromMapMaybe(params, "order_no").ToStringNoPoint(),
				SceneName:    container.GetValFromMapMaybe(params, "scene_name").ToStringNoPoint(),
				Lang:         ss_net.GetCommonData(c).Lang,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, nil
			}

			file, err := os.Open(reply.FilePath)
			if err != nil {
				ss_log.Error("err=[%v]\n", err)
				return ss_err.ERR_SYS_IO_ERR, nil, nil
			}

			buff, err := ioutil.ReadAll(file)
			if err != nil {
				ss_log.Error("err=[%v]\n", err)
				return ss_err.ERR_SYS_IO_ERR, nil, nil
			}

			c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", reply.FilePath))
			c.Writer.Header().Add("Content-Type", "application/octet-stream")
			c.Writer.Header().Add("Content-Length", strext.ToString(len(buff)))
			c.Data(http.StatusOK, "application/octet-stream", buff)
			//c.File(reply.FilePath)
			_ = file.Close()

			//删除临时文件
			errDel := os.Remove(reply.FilePath)
			if errDel != nil {
				// 删除失败
				ss_log.Error("删除临时文件失败,err=[%v],fileName=[%v]", errDel, reply.FilePath)
			} else {
				ss_log.Info("已删除临时文件fileName[%v]", reply.FilePath)
			}
			return ss_err.ERR_SUCCESS, nil, nil
		})
	}
}

/**
支付完成订单退款
*/
func (*BillHandler) BusinessBillRefund() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &businessBillProto.BusinessBillRefundRequest{
				AccountType:   inner_util.GetJwtDataString(c, "account_type"),
				BusinessAccNo: inner_util.GetJwtDataString(c, "account_uid"),
				BusinessNo:    inner_util.GetJwtDataString(c, "iden_no"),
				PaymentPwd:    container.GetValFromMapMaybe(params, "payment_pwd").ToStringNoPoint(),
				NonStr:        container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				OrderNo:       container.GetValFromMapMaybe(params, "order_no").ToStringNoPoint(),
				RefundAmount:  container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				RefundReason:  container.GetValFromMapMaybe(params, "remarks").ToStringNoPoint(),
				Lang:          ss_net.GetCommonData(c).Lang,
			}
			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}
			if errStr := verify.BusinessBillRefundVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := BusinessBillHandlerInst.Client.BusinessBillRefund(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, nil
			}

			return ss_err.ERR_SUCCESS, gin.H{
				"log_no":       reply.RefundNo,
				"order_status": reply.RefundStatus,
			}, nil
		})
	}
}

/**
退款订单列表
*/
func (*BillHandler) GetRefundBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetRefundBillsRequest{
				BusinessNo:   inner_util.GetJwtDataString(c, "iden_no"),
				Page:         strext.ToStringNoPoint(params[0]),
				PageSize:     strext.ToStringNoPoint(params[1]),
				OrderStatus:  strext.ToStringNoPoint(params[2]), //订单状态
				CurrencyType: strext.ToStringNoPoint(params[3]), //币种
				StartTime:    strext.ToStringNoPoint(params[4]),
				EndTime:      strext.ToStringNoPoint(params[5]),
				RefundNo:     strext.ToStringNoPoint(params[6]),
				TransOrderNo: strext.ToStringNoPoint(params[7]),
				SceneNo:      strext.ToStringNoPoint(params[8]),
				Lang:         ss_net.GetCommonData(c).Lang,
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			reply, err := CustHandlerInst.Client.GetRefundBills(context.TODO(), req)
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
					"refund_no":      v.RefundNo,
					"out_refund_no":  v.OutRefundNo,
					"order_status":   v.OrderStatus,
					"trans_order_no": v.TransOrderNo,
					"amount":         v.Amount,
					"currency_type":  v.CurrencyType,
					"create_time":    v.CreateTime,
					"finish_time":    v.FinishTime,
					"subject":        v.Subject,
					"app_name":       v.AppName,
					"payee_account":  v.PayeeAccount,
					"scene_name":     v.SceneName,
				}
				list = append(list, data)
			}

			return reply.ResultCode, gin.H{
				"datas":   list,         //订单列表
				"usd_cnt": reply.UsdCnt, //usd订单数量统计
				"usd_sum": reply.UsdSum, //usd订单金额统计
				"khr_cnt": reply.KhrCnt,
				"khr_sum": reply.KhrSum,
			}, reply.Total, nil
		}, "page", "page_size", "order_status", "currency_type", "start_time", "end_time", "refund_no",
			"trans_order_no", "scene_no")
	}
}

/**
退款详情
*/
func (*BillHandler) GetRefundDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetRefundDetail(context.TODO(), &custProto.GetRefundDetailRequest{
				RefundNo: strext.ToStringNoPoint(params[0]),
				Lang:     ss_net.GetCommonData(c).Lang,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			//账号脱敏
			reply.Data.PayeeAccount = ss_func.GetDesensitizationAccount(reply.Data.PayeeAccount)
			return reply.ResultCode, reply.Data, 0, nil
		}, "refund_no")
	}
}

//下载退款订单文件
func (*BillHandler) DownloadRefundBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.CreateRefundFile(context.TODO(), &custProto.CreateRefundFileRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				BusinessNo:   inner_util.GetJwtDataString(c, "iden_no"),
				Page:         0,
				PageSize:     1000,
				OrderStatus:  container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(), //订单状态
				RefundNo:     container.GetValFromMapMaybe(params, "refund_no").ToStringNoPoint(),
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				StartTime:    container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:      container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				TransOrderNo: container.GetValFromMapMaybe(params, "trans_order_no").ToStringNoPoint(),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, nil
			}

			file, err := os.Open(reply.FilePath)
			if err != nil {
				ss_log.Error("err=[%v]\n", err)
				return ss_err.ERR_SYS_IO_ERR, nil, nil
			}

			buff, err := ioutil.ReadAll(file)
			if err != nil {
				ss_log.Error("err=[%v]\n", err)
				return ss_err.ERR_SYS_IO_ERR, nil, nil
			}
			c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", reply.FilePath))
			c.Writer.Header().Add("Content-Type", "application/octet-stream")
			c.Data(http.StatusOK, "application/octet-stream", buff)
			//c.File(reply.FilePath)

			_ = file.Close()

			//删除临时文件
			errDel := os.Remove(reply.FilePath)
			if errDel != nil {
				// 删除失败
				ss_log.Error("删除临时文件失败,err=[%v],fileName=[%v]", errDel, reply.FilePath)
			} else {
				ss_log.Info("已删除临时文件fileName[%v]", reply.FilePath)
			}

			return ss_err.ERR_SUCCESS, nil, nil
		})
	}
}
