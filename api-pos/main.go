package main

import (
	_ "database/sql"

	"a.a/cu/ss_log"
	"a.a/mp-server/api-pos/i"
	"a.a/mp-server/api-pos/router"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	_ "github.com/lib/pq"
)

func main() {
	servname := constants.ServerNameApiPos

	// 初始化日志
	ss_log2.DoInit(servname)
	ss_log.Info(servname)
	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)

	host, portFrom, portTo := i.DoInitBase()
	ss_serv.DoInitCache()
	ss_serv.DoInitDb(constants.DB_CRM)
	ss_serv.DoInitMultiFromDB(constants.DB_CRM) // 初始化多语言，在数据库初始化后

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
