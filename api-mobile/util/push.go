package util

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/api-mobile/common"
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
	if err := (*common.MqPushMsg).Publish("a", brokerMsg); err != nil {
		ss_log.Info("[pub]failed: %v", err)
	} else {
		ss_log.Error("[pub]sent")
	}
}
