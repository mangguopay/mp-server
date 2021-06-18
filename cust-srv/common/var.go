package common

import (
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_serv"
	"errors"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
)

var (
	MqPushMsg *broker.Broker
	MqTopic   string
)

// 存放服务商坐标的map
//var SrvCoordinateMap = make(map[string]string)

// 存放服务商坐标的map
var SrvCoordinates = make([]*go_micro_srv_cust.NearbyServicerData, 0)

var SrvGisPub = micro.NewEvent(constants.Topic_Event_Srv_Gis, ss_serv.GetRpcDefClient())

var PushEvent = micro.NewEvent(constants.Nats_Broker_Header_Send_Push_Msg, ss_serv.GetRpcDefClient())

var SmsEvent = micro.NewEvent(constants.Nats_Broker_Header_Reg_SMS, ss_serv.GetRpcDefClient())

var GetDBConnectFailedErr = errors.New("获取数据库连接失败")
