package main

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/api-pay/i"
	"a.a/mp-server/api-pay/router"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
)

func main() {
	servname := constants.ServerNameApiPay

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
	i.InitSrv()

	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		Router:       router.InitRouter(),
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
	}
	sa.RunAsWeb()
}
