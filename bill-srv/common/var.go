package common

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_serv"
	micro "github.com/micro/go-micro/v2"
)

// 推送事件
var PushEvent = micro.NewEvent(constants.Nats_Broker_Header_Send_Push_Msg, ss_serv.GetRpcDefClient())

// 清分事件
var SettleEvent = micro.NewEvent(constants.Settle_Type, ss_serv.GetRpcDefClient())

// 推送核销码短信事件
var WriteOffEvent = micro.NewEvent(constants.Nats_Broker_Header_Write_Off, ss_serv.GetRpcDefClient())

// 当前的服务id
var BillServerFullId string
