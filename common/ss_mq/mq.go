package ss_mq

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/v2/broker"
	"sync"
)

var (
	SsTopicInst SsTopic
)

type SsTopic struct {
	//          topic    msger
	TopicM sync.Map //[string]*broker.Broker
}

func (r *SsTopic) SetBlocker(topic string, b *broker.Broker) {
	r.TopicM.Store(topic, b)
}

func (r *SsTopic) PushMsg(topic string, msg proto.Message, reqstr string) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	brokerMsg := &broker.Message{
		Header: map[string]string{
			"m": reqstr,
		},
		Body: b,
	}
	brokerT, ok := r.TopicM.Load(topic)
	if !ok {
		ss_log.Error("broker is nil, topic: %s", topic)
		return ss_err.ErrBrokerIsNil
	}

	switch brokerT.(type) {
	case *broker.Broker:
		// do nothing...
	default:
		ss_log.Error("broker is nil")
		return ss_err.ErrBrokerIsNil
	}

	if err := (*brokerT.(*broker.Broker)).Publish(topic, brokerMsg); err != nil {
		ss_log.Info("[pub]failed: %v", err)
	} else {
		ss_log.Error("[pub]sent")
	}

	return err
}
