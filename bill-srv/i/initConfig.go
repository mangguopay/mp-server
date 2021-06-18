package i

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_struct"
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

// 初始化 aws s3 配置
func InitAwss3() {
	var s3Conf ss_struct.Awss3Conf

	// 获取配置信息
	err := config.Get("aws", "s3").Scan(&s3Conf)
	if err != nil {
		ss_log.Error("aws-s3初始化失败,err:%v", err)
	}

	// 验证配置信息
	if s3Conf.AccessKeyId == "" || s3Conf.SecretAccessKey == "" || s3Conf.Region == "" || s3Conf.Bucket == "" {
		ss_log.Error("aws-s3配置信息不完整,s3Conf:%+v", s3Conf)
		panic(fmt.Sprintf("aws-s3配置信息不完整,s3Conf:%+v", s3Conf))
	}

	ss_log.Info("s3Conf:%+v \n", s3Conf)

	// 初始化s3操作类
	common.InitUploadS3(s3Conf)
}

func DoInitNatsCli() {
	doInitNatsCliSingle(constants.Nats_Cluster_Primary, constants.Nats_Topic_Settle)
	doInitNatsCliSingle(constants.Nats_Cluster_Secondary, constants.Nats_Topic_Push)
}

func doInitNatsCliSingle(clusterId, topic string) {
	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	switch m["adapter"] {
	case "nats":
		// FIXME
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
		// FIXME
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
