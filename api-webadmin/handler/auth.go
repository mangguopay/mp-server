package handler

import (
	colloection "a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/mp-server/api-webadmin/inner_util"
	"a.a/mp-server/api-webadmin/verify"
	AuthProto "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

var (
	AuthHandlerInst AuthHandler
)

type AuthHandler struct {
	Client AuthProto.AuthService
}

func (a *AuthHandler) AddCashier() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &AuthProto.AddCashierRequest{
				Phone:       colloection.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				CountryCode: colloection.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(),
				ServicerNo:  colloection.GetValFromMapMaybe(params, "servicer_no").ToString(),
				LoginUid:    inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.CheckAddCashierVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			req.Phone = ss_func.PrePhone(req.CountryCode, req.Phone)

			reply, err := AuthHandlerInst.Client.AddCashier(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, reply.Uid, err
		})
	}
}
