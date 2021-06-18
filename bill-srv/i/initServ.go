package i

import (
	"a.a/mp-server/common/constants"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	custProto "a.a/mp-server/common/proto/cust"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	"a.a/mp-server/common/ss_serv"
)

type AuthHandler struct {
	Client go_micro_srv_auth.AuthService
}

var (
	AuthHandlerInst AuthHandler
)

type QuotaHandle struct {
	Client go_micro_srv_quota.QuotaService
}

var QuotaHandleInstance QuotaHandle

type RiskCtrlHandle struct {
	Client go_micro_srv_riskctrl.RiskCtrlService
}

var CustHandleInstance CustHandle

type CustHandle struct {
	Client custProto.CustService
}

var RiskCtrlHandleInstance RiskCtrlHandle

func InitSrv() {
	AuthHandlerInst.Client = go_micro_srv_auth.NewAuthService(constants.ServerNameAuth, ss_serv.GetRpcDefClient())
	QuotaHandleInstance.Client = go_micro_srv_quota.NewQuotaService(constants.ServerNameQuota, ss_serv.GetRpcDefClient())
	RiskCtrlHandleInstance.Client = go_micro_srv_riskctrl.NewRiskCtrlService(constants.ServerNameRiskctrl, ss_serv.GetRpcDefClient())
	CustHandleInstance.Client = custProto.NewCustService(constants.ServerNameCust, ss_serv.GetRpcDefClient())
}
