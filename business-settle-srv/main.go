package main

import (
	"net"

	"a.a/cu/ss_log"
	"a.a/mp-server/business-settle-srv/i"
	"a.a/mp-server/business-settle-srv/subcriber"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameBusinessSettle

	// 初始化日志
	ss_log2.DoInit(servname)
	ss_log.Info(servname)

	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)

	// 初始化基础配置信息
	host, portFrom, portTo := i.DoInitBase()
	ss_serv.DoInitCache()
	ss_serv.DoInitDb(constants.DB_CRM)

	natsHost, natsPort := i.GetPrimaryNatsAddr()

	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			return subcriber.InitSubcriber(s)
		},
		BrokerType: constants.BrokerTypeNats,
		BrokerAddr: []string{net.JoinHostPort(natsHost, natsPort)},
	}
	sa.RunAsSrv()
}
