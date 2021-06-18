package common

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_serv"
	"github.com/micro/go-micro/v2"
)

// 二维码过期监听
var ListenExpKey = micro.NewEvent(constants.Nats_Listen_Exp_key, ss_serv.GetRpcDefClient())
