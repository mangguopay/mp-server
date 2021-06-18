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

//查询终端列表
func (*CustHandler) GetTerminalList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetTerminalList(context.TODO(), &custProto.GetTerminalListRequest{
				Page:            strext.ToInt32(params[0]),
				PageSize:        strext.ToInt32(params[1]),
				TerminalNumber:  strext.ToString(params[2]),
				PosSn:           strext.ToString(params[3]),
				ServicerAccount: strext.ToString(params[4]),
				UseStatus:       strext.ToString(params[5]),
			})
			return reply.ResultCode, reply.List, reply.Total, err
		}, "page", "page_size", "terminal_number", "pos_sn", "servicer_account", "use_status")
	}
}

//添加终端
func (*CustHandler) AddTerminal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.AddTerminalRequest{
				ServicerAccount: container.GetValFromMapMaybe(params, "servicer_account").ToStringNoPoint(),
				CountryCode:     container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(),
				TerminalNumber:  container.GetValFromMapMaybe(params, "terminal_number").ToStringNoPoint(),
				PosSn:           container.GetValFromMapMaybe(params, "pos_sn").ToStringNoPoint(),
				UseStatus:       container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginAccount:    inner_util.GetJwtDataString(c, "account_uid"),
			}

			reply, err := CustHandlerInst.Client.AddTerminal(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.AddTerminal()失败, err=%v", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, nil, nil
		})
	}
}

//修改终端使用状态
func (*CustHandler) UpdateTerminalStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateTerminalRequest{
				TerminalNo:   container.GetValFromMapMaybe(params, "terminal_no").ToStringNoPoint(),
				UseStatus:    container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginAccount: inner_util.GetJwtDataString(c, "account_uid"),
			}

			reply, err := CustHandlerInst.Client.UpdateTerminal(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.AddTerminal()失败, err=%v", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, nil, nil
		})
	}
}

//修改终端使用状态
func (*CustHandler) DelTerminal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteTerminalRequest{
				TerminalNo:   container.GetValFromMapMaybe(params, "terminal_no").ToStringNoPoint(),
				LoginAccount: inner_util.GetJwtDataString(c, "account_uid"),
			}

			reply, err := CustHandlerInst.Client.DeleteTerminal(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用cust-srv.AddTerminal()失败, err=%v", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, nil, nil
		})
	}
}
