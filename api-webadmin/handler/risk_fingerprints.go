package handler

import (
	"a.a/cu/container"
	"a.a/mp-server/api-webadmin/inner_util"
	"context"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"

	"github.com/gin-gonic/gin"
)

func (CustHandler) GetAppFingerprintOn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAppFingerprintOn(context.TODO(), &custProto.GetAppFingerprintOnRequest{})
			if err != nil {
				ss_log.Error("调用cust-srv.GetFingerprintList()失败，err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			data := gin.H{
				"is_open": reply.IsOpen,
			}
			return reply.ResultCode, data, 0, err
		})
	}
}

func (CustHandler) GetFingerprintList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetFingerprintList(context.TODO(), &custProto.GetFingerprintListRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				Account:   strext.ToString(params[4]), //账号
				DeviceNo:  strext.ToString(params[5]), //设备号
				UseStatus: strext.ToString(params[6]),
			})
			if err != nil {
				ss_log.Error("调用cust-srv.GetFingerprintList()失败，err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.List, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "account", "device_no", "use_status")
	}
}

func (s *CustHandler) CloseFingerprintInputFunction() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.CloseFingerprintFunctionRequest{
				IsOpen:   container.GetValFromMapMaybe(params, "is_open").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}

			reply, err := CustHandlerInst.Client.CloseFingerprintFunction(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.CloseFingerprintFunction()失败，err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, nil, err
		})
	}
}

func (s *CustHandler) CleanFingerprintData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.CleanFingerprintDataRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				OpType:   container.GetValFromMapMaybe(params, "op_type").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}

			reply, err := CustHandlerInst.Client.CleanFingerprintData(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.CleanFingerprintDataRequest()失败，err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, nil, err
		})
	}
}
