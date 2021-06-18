package ss_net

type NatsCliConf struct {
	ClusterId string
	Topic     string
}

func DoInitNatsCli(nas []*NatsCliConf) {
	for _, v := range nas {
		doInitNatsCliSingle(v.ClusterId, v.Topic)
	}
}

func doInitNatsCliSingle(clusterId, topic string) {
	// FIXME 不兼容
	//m := map[string]interface{}{}
	//err := config.Get("mq", clusterId).Scan(&m)
	//if err != nil {
	//	ss_log.Error("err=%v", err)
	//}
	//
	//switch m["adapter"] {
	//case "nats":
	//	ss_log.Info("init nats=[%v:%v@%v] client", m["host"], m["port"], clusterId)
	//	options := stancf.GetDefaultOptions()
	//	options.NatsURL = net.JoinHostPort(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
	//	na := stan.NewBroker(stan.ClusterID(clusterId), stan.Options(options))
	//	if err := na.Init(); err != nil {
	//		panic(fmt.Sprintf("Broker Init error: %v", err))
	//	}
	//	if err := na.Connect(); err != nil {
	//		panic(fmt.Sprintf("Broker Connect error: %v", err))
	//	}
	//	ss_mq.SsTopicInst.SetBlocker(topic, &na)
	//}
}
