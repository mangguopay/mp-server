package util

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/cust-srv/common"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/v2/broker"
)

// 发送异步消息
func PushMsg(msg proto.Message, reqstr string) {
	b, _ := proto.Marshal(msg)
	brokerMsg := &broker.Message{
		Header: map[string]string{
			"m": reqstr,
		},
		Body: b,
	}
	if err := (*common.MqPushMsg).Publish(common.MqTopic, brokerMsg); err != nil {
		ss_log.Info("MqTopic:%s,pubFailed: %v", common.MqTopic, err)
	} else {
		ss_log.Error("MqTopic:%s,pubSent", common.MqTopic)
	}
}
