package i

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/config"
)

func DoInitBase() (host string, portFrom, portTo int) {
	m := map[string]interface{}{}
	err := config.Get("base", "base").Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}
	z, zErr := time.LoadLocation(strext.ToStringNoPoint(m["timezone"]))
	if zErr != nil {
		panic(fmt.Sprintf("解析时区出错,err: %v", zErr))
	}
	// 设置time包中的默认时区
	time.Local = z

	global.Tz = z
	p := strings.Split(strext.ToStringNoPoint(m["port"]), "-")
	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}

func DoInitNatsCli() {
	doInitNatsCliSingle(constants.Nats_Cluster_Primary, constants.Nats_Topic_Risk)
}

func doInitNatsCliSingle(clusterId, topic string) {
	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	switch m["adapter"] {
	case "nats":
		//ss_log.Info("init nats=[%v:%v@%v] client", m["host"], m["port"], clusterId)
		//options := stancf.GetDefaultOptions()
		//options.NatsURL = net.JoinHostPort(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
		//na := stan.NewBroker(stan.ClusterID(clusterId), stan.Options(options))
		//if err := na.Init(); err != nil {
		//	panic(fmt.Sprintf("Broker Init error: %v", err))
		//}
		//if err := na.Connect(); err != nil {
		//	panic(fmt.Sprintf("Broker Connect error: %v", err))
		//}
		//ss_mq.SsTopicInst.SetBlocker(topic, &na)
	}
}

func DoInitNats(r func(*broker.Broker)) {
	clusterId := constants.Nats_Cluster_Primary

	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	switch m["adapter"] {
	case "nats":
		//ss_log.Info("init nats=[%v:%v@%v] client", m["host"], m["port"], clusterId)
		//options := stancf.GetDefaultOptions()
		//options.NatsURL = net.JoinHostPort(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
		//na := stan.NewBroker(stan.ClusterID(clusterId), stan.Options(options))
		//if err := na.Init(); err != nil {
		//	panic(fmt.Sprintf("Broker Init error: %v", err))
		//}
		//if err := na.Connect(); err != nil {
		//	panic(fmt.Sprintf("Broker Connect error: %v", err))
		//}
		//r(&na)
	}
}
func GetPrimaryNatsAddr() (host, port string) {
	clusterId := constants.Nats_Cluster_Primary

	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	return strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"])
}
