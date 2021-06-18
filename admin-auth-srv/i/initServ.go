package i

import (
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	"a.a/mp-server/common/ss_serv"
)

type RiskCtrlHandle struct {
	Client go_micro_srv_riskctrl.RiskCtrlService
}

type CustHandler struct {
	Client go_micro_srv_cust.CustService
}

var (
	CustHandlerInst    CustHandler
	RiskCtrlHandleInst RiskCtrlHandle
)

func InitSrv() {
	CustHandlerInst.Client = go_micro_srv_cust.NewCustService(constants.ServerNameCust, ss_serv.GetRpcDefClient())
	RiskCtrlHandleInst.Client = go_micro_srv_riskctrl.NewRiskCtrlService(constants.ServerNameRiskctrl, ss_serv.GetRpcDefClient())
}
