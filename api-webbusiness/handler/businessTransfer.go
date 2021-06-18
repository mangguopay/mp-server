package handler

import (
	"a.a/cu/container"
	"a.a/mp-server/api-webbusiness/verify"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	billProto "a.a/mp-server/common/proto/bill"
	"context"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webbusiness/inner_util"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"

	"github.com/gin-gonic/gin"
)

type BusinessTransferHandler struct {
}

/**
商家单笔付款
*/
func (*BusinessTransferHandler) AddBusinessTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &billProto.AddBusinessTransferRequest{
				BusinessAccNo: inner_util.GetJwtDataString(c, "account_uid"),
				BusinessNo:    inner_util.GetJwtDataString(c, "iden_no"),
				AccountType:   inner_util.GetJwtDataString(c, "account_type"),
				PayeeNo:       container.GetValFromMapMaybe(params, "payee_no").ToStringNoPoint(), //转账账号
				Amount:        container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				PaymentPwd:    container.GetValFromMapMaybe(params, "payment_pwd").ToStringNoPoint(),
				NonStr:        container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				Remarks:       container.GetValFromMapMaybe(params, "remarks").ToStringNoPoint(),
				CountryCode:   container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(),
				CurrencyType:  container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				Lang:          ss_net.GetCommonData(c).Lang,
				TransferType:  constants.BusinessTransferOrderTypeOrdinary,
			}

			if errStr := verify.AddBusinessTransferVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := BillHandlerInst.Client.AddBusinessTransfer(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, gin.H{
				"log_no":       reply.LogNo,
				"order_status": reply.OrderStatus,
			}, nil
		})
	}
}

/**
商家转账订单列表
*/
func (*BusinessTransferHandler) GetBusinessTransferOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessTransferOrderListRequest{
				Page:          strext.ToStringNoPoint(params[0]),
				PageSize:      strext.ToStringNoPoint(params[1]),
				StartTime:     strext.ToStringNoPoint(params[2]),
				EndTime:       strext.ToStringNoPoint(params[3]),
				CurrencyType:  strext.ToStringNoPoint(params[4]),
				LogNo:         strext.ToStringNoPoint(params[5]),
				OrderStatus:   strext.ToStringNoPoint(params[6]),
				ToAccountNo:   strext.ToStringNoPoint(params[7]),
				BatchNo:       strext.ToStringNoPoint(params[8]),
				TransferType:  strext.ToStringNoPoint(params[9]),
				OutTransferNo: strext.ToStringNoPoint(params[10]),
				BusinessNo:    inner_util.GetJwtDataString(c, "iden_no"),
				BusinessAccNo: inner_util.GetJwtDataString(c, "account_uid"),
				Lang:          ss_net.GetCommonData(c).Lang,
			}

			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}
			if req.BusinessAccNo == "" {
				ss_log.Error("IdenNo参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}
			reply, err := CustHandlerInst.Client.GetBusinessTransferOrderList(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用服务失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.List, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "currency_type", "log_no", "order_status",
			"to_account_no", "batch_no", "transfer_type", "out_transfer_no")
	}
}

/**
商家转账订单详情
*/
func (BusinessTransferHandler) GetBusinessTransferOrderDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessTransferOrderDetail(context.TODO(), &custProto.GetBusinessTransferOrderDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			})
			if err != nil {
				ss_log.Error("调用cust服务失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Data, 0, nil
		}, "log_no")
	}
}

/**
商家转账批次列表
*/
func (*BusinessTransferHandler) GetBusinessTransferBatchList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessTransferBatchListRequest{
				Page:          strext.ToStringNoPoint(params[0]),
				PageSize:      strext.ToStringNoPoint(params[1]),
				StartTime:     strext.ToStringNoPoint(params[2]),
				EndTime:       strext.ToStringNoPoint(params[3]),
				CurrencyType:  strext.ToStringNoPoint(params[4]),
				BatchNo:       strext.ToStringNoPoint(params[5]),
				Status:        strext.ToStringNoPoint(params[6]),
				BusinessNo:    inner_util.GetJwtDataString(c, "iden_no"),
				BusinessAccNo: inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.BusinessNo == "" {
				ss_log.Error("IdenNo参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessTransferBatchList(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用服务失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "currency_type", "batch_no", "order_status")
	}
}

/**
商家转账批次列表
*/
func (*BusinessTransferHandler) GetBusinessTransferBatchDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessTransferBatchDetailRequest{
				BatchNo: strext.ToStringNoPoint(params[0]),
			}
			if req.BatchNo == "" {
				ss_log.Error("BatchNo参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessTransferBatchDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用服务失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.BatchData, 0, nil
		}, "batch_no")
	}
}

/**
获取商家转账批次文件分析结果
*/
func (*BusinessTransferHandler) GetBatchAnalysisResult() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &billProto.GetBatchAnalysisResultRequest{
				FileId:     container.GetValFromMapMaybe(params, "file_id").ToStringNoPoint(),
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
				BusinessNo: inner_util.GetJwtDataString(c, "iden_no"),
				Lang:       ss_net.GetCommonData(c).Lang,
			}
			if req.Lang == "" {
				req.Lang = constants.DefaultLang
			}

			if errStr := verify.GetBatchAnalysisResult(req); errStr != "" {
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := BillHandlerInst.Client.GetBatchAnalysisResult(context.TODO(), req, global.RequestTimeoutOptions)

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"data":        reply.Data,
				"wrong_datas": reply.WrongDatas,
				"batch_no":    reply.BatchNo,
			}, err
		}, "params")
	}
}

/**
商家确认分析结果，开始执行转账
*/
func (*BusinessTransferHandler) BatchConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &billProto.BusinessBatchTransferConfirmRequest{
				BusinessNo:    inner_util.GetJwtDataString(c, "iden_no"),
				BusinessAccNo: inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:   inner_util.GetJwtDataString(c, "account_type"),
				PayPwd:        container.GetValFromMapMaybe(params, "payment_pwd").ToStringNoPoint(),
				NonStr:        container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				BatchNo:       container.GetValFromMapMaybe(params, "batch_no").ToStringNoPoint(),
			}

			if errStr := verify.BusinessBatchTransferConfirmVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := BillHandlerInst.Client.BusinessBatchTransferConfirm(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, nil, nil
		})
	}
}
