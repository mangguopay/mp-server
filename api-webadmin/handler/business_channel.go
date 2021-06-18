package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/inner_util"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

//全部渠道（不加任何条件）
func (*CustHandler) GetAllPaymentChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetAllPaymentChannelRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				ChannelName: strext.ToStringNoPoint(params[4]),
			}
			reply, err := CustHandlerInst.Client.GetAllPaymentChannel(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, 0, nil
			}

			var list []interface{}
			for _, v := range reply.Channel {
				data := gin.H{
					"channel_no":   v.ChannelNo,
					"channel_name": v.ChannelName,
					"channel_type": v.ChannelType,
					"upstream_no":  v.UpstreamNo,
					"create_time":  v.CreateTime,
				}
				list = append(list, data)
			}

			return reply.ResultCode, list, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "search")
	}
}

func (*CustHandler) AddPaymentChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.AddPaymentChannelRequest{
				LoginAccount: inner_util.GetJwtDataString(c, "account_uid"),
				ChannelName:  container.GetValFromMapMaybe(params, "channel_name").ToStringNoPoint(),
				ChannelType:  container.GetValFromMapMaybe(params, "channel_type").ToStringNoPoint(),
				UpstreamNo:   container.GetValFromMapMaybe(params, "upstream_no").ToStringNoPoint(),
			}
			reply, err := CustHandlerInst.Client.AddPaymentChannel(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.AddPaymentChannel()失败，err=%v", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, nil
			}
			return reply.ResultCode, gin.H{
				"channel_no": reply.ChannelNo,
			}, nil
		})
	}
}

func (*CustHandler) UpdatePaymentChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdatePaymentChannelRequest{
				LoginAccount: inner_util.GetJwtDataString(c, "account_uid"),
				ChannelNo:    container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				ChannelName:  container.GetValFromMapMaybe(params, "channel_name").ToStringNoPoint(),
				ChannelType:  container.GetValFromMapMaybe(params, "channel_type").ToStringNoPoint(),
				UpstreamNo:   container.GetValFromMapMaybe(params, "upstream_no").ToStringNoPoint(),
			}
			reply, err := CustHandlerInst.Client.UpdatePaymentChannel(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.UpdatePaymentChannel()失败，err=%v", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, nil, nil
		})
	}
}
