package main

import (
	"net"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"a.a/mp-server/cust-srv/cron"
	"a.a/mp-server/cust-srv/handler"
	"a.a/mp-server/cust-srv/i"
	_ "github.com/lib/pq"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameCust

	// 初始化日志
	ss_log2.DoInit(servname)
	ss_log.Info(servname)
	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)
	// 初始化基础配置信息
	host, portFrom, portTo := i.DoInitBase()
	ss_serv.DoInitCache()
	ss_serv.DoInitDb("risk", constants.DB_CRM)
	i.DoInitIDW()
	i.InitSrv()
	i.InitAwss3()

	//定时器
	cron.DoInitCronTask()

	// 事件
	natsHost, natsPort := i.GetPrimaryNatsAddr()
	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			return cust.RegisterCustHandler(s, new(handler.CustHandler))
		},
		BrokerType: constants.BrokerTypeNats,
		BrokerAddr: []string{net.JoinHostPort(natsHost, natsPort)},
	}
	sa.RunAsSrv()
}
