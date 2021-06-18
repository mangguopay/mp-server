package i

import (
	"a.a/mp-server/common/constants"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_serv"
)

type AuthHandler struct {
	Client go_micro_srv_auth.AuthService
}

var (
	AuthHandlerInst AuthHandler
)

func InitSrv() {
	AuthHandlerInst.Client = go_micro_srv_auth.NewAuthService(constants.ServerNameAuth, ss_serv.GetRpcDefClient())
}
