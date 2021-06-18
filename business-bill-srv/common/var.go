package common

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_serv"
	"errors"
	micro "github.com/micro/go-micro/v2"
)

//支付订单结果异步通知事件
var PayResultNotifyEvent = micro.NewEvent(constants.PaySystemResultNotify, ss_serv.GetRpcDefClient())

var RedisValueNilErr = errors.New("redis: nil")

// 推送事件
var PushEvent = micro.NewEvent(constants.Nats_Broker_Header_Send_Push_Msg, ss_serv.GetRpcDefClient())
