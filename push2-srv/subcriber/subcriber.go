package subcriber

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/push2-srv/adapter"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
)

func InitSubcriber(s server.Server) error {
	var err error

	err = micro.RegisterSubscriber(constants.Nats_Broker_Header_Send_Push_Msg, s, adapter.Push, server.SubscriberQueue("queue.pubsub"))
	if err != nil {
		return err
	}

	err = micro.RegisterSubscriber(constants.Nats_Broker_Header_Reg_SMS, s, adapter.Push, server.SubscriberQueue("queue.pubsub"))
	if err != nil {
		return err
	}

	err = micro.RegisterSubscriber(constants.Nats_Broker_Header_Write_Off, s, adapter.Push, server.SubscriberQueue("queue.pubsub"))
	if err != nil {
		return err
	}

	return nil
}
