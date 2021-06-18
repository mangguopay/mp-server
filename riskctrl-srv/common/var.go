package common

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_serv"
	"github.com/micro/go-micro/v2"
)

// 推送事件
var RiskEvent = micro.NewEvent(constants.Nats_Broker_Header_Risk, ss_serv.GetRpcDefClient())
