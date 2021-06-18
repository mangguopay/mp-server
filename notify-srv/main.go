package main

import (
	"a.a/mp-server/notify-srv/cron"
	"net"

	"a.a/mp-server/notify-srv/subcriber"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"a.a/mp-server/notify-srv/i"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNamePayNotify

	// 初始化日志
	ss_log2.DoInit(servname)
	ss_log.Info(servname)

	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)

	// 初始化基础配置信息
	host, portFrom, portTo := i.DoInitBase()
	i.DoInitCache()
	ss_serv.DoInitDb(constants.DB_CRM)
	i.InitSrv()

	//监听异步通知失败等待到期的订单
	subcriber.DoListenRedisExpKey()

	//初始化定时任务
	cron.DoInitCronTask()

	// 定时查询漏掉统计手续费的订单
	natsHost, natsPort := i.GetPrimaryNatsAddr()

	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			return subcriber.InitSubscriber(s)
		},
		BrokerType: constants.BrokerTypeNats,
		BrokerAddr: []string{net.JoinHostPort(natsHost, natsPort)},
	}
	sa.RunAsSrv()
}
