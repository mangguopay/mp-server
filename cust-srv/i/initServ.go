package i

import (
	"a.a/mp-server/common/constants"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	go_micro_srv_gis "a.a/mp-server/common/proto/gis"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	"a.a/mp-server/common/ss_serv"
)

type AuthHandler struct {
	Client go_micro_srv_auth.AuthService
}

var (
	AuthHandlerInst     AuthHandler
	QuotaHandleInstance QuotaHandle
	GisHandleInst       GisHandle
)

type QuotaHandle struct {
	Client go_micro_srv_quota.QuotaService
}

type GisHandle struct {
	Client go_micro_srv_gis.GISService
}

func InitSrv() {
	AuthHandlerInst.Client = go_micro_srv_auth.NewAuthService(constants.ServerNameAuth, ss_serv.GetRpcDefClient())
	QuotaHandleInstance.Client = go_micro_srv_quota.NewQuotaService(constants.ServerNameQuota, ss_serv.GetRpcDefClient())
	GisHandleInst.Client = go_micro_srv_gis.NewGISService(constants.ServerNameGis, ss_serv.GetRpcDefClient())
}
