package util

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_mq"
	"fmt"
	"github.com/micro/go-plugins/broker/stan"
	stancf "github.com/nats-io/stan.go"
	"net"
	"testing"
)

func TestPushMsg(t *testing.T) {
	m := map[string]interface{}{}
	m["host"] = "127.0.0.1"
	m["port"] = "4223"
	m["adapter"] = "nats"

	clusterId := "settle"
	topic := constants.Nats_Topic_Settle

	ss_log.Info("init nats=[%v:%v@%v] client", m["host"], m["port"], clusterId)
	options := stancf.GetDefaultOptions()
	options.NatsURL = net.JoinHostPort(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
	na := stan.NewBroker(stan.ClusterID(clusterId), stan.Options(options))
	if err := na.Init(); err != nil {
		panic(fmt.Sprintf("Broker Init error: %v", err))
	}
	if err := na.Connect(); err != nil {
		panic(fmt.Sprintf("Broker Connect error: %v", err))
	}
	ss_mq.SsTopicInst.SetBlocker(topic, &na)

	msg := &go_micro_srv_settle.SettleTransferRequest{
		BillNo:    "13800138000",
		OrderType: "1",
	}

	//b, _ := proto.Marshal(msg)
	//brokerMsg := broker.Message{
	//	Header: map[string]string{
	//		"m": "settle",
	//	},
	//	Body: b,
	//}
	//ss_log.Info("-------->%s", string(b))
	for i := 0; i < 10; i++ {
		if err := ss_mq.SsTopicInst.PushMsg(topic, msg, "settle"); err != nil {
			ss_log.Info("[pub]failed: %v", err)
		} else {
			ss_log.Error("[pub]sent")
		}
	}
}
