package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"

	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
}

var AccountHandlerInst AccountHandler

func (*AccountHandler) CheckAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(account interface{}) (string, interface{}, error) {
			reply, err := AuthHandlerInst.Client.CheckAccount(context.TODO(), &go_micro_srv_auth.CheckAccountRequest{
				Account: strext.ToString(account),
			})
			return reply.ResultCode, reply.Data, err
		}, "account")
	}
}
func (*AccountHandler) CheckMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(email interface{}) (string, interface{}, error) {
			req := &go_micro_srv_cust.CheckMailRequest{
				Email: strext.ToString(email),
			}
			if req.Email == "" {
				ss_log.Error("参数Email为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.CheckMail(context.TODO(), req)
			return reply.ResultCode, nil, err
		}, "email")
	}
}
