package main

import (
	_ "database/sql"

	"a.a/cu/ss_log"
	"a.a/mp-server/api-cb/i"
	"a.a/mp-server/api-cb/router"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	_ "github.com/lib/pq"
)

func main() {
	servname := constants.ServerNameApiApiCb

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
	//i.DoInitNatsCli()
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

//func main() {
//	i.DoInitConfig()
//	// 日志
//	ss_log.InitLog(path.Join("..", "runlog", "api-cb", "logs"))
//	defer ss_log.CloseLog()
//	//
//	ss_log.Info("go.micro.api.cb")
//	host, portFrom, portTo := i.DoInitBase()
//	i.DoInitRedis()
//	i.DoInitDb()
//	i.DoInitNatsCli()
//	i.InitSrv()
//
//	etcReg := etcdv3.NewRegistry()
//	var service web.Service
//	var err error = errors.New("start")
//	port := portFrom
//	for err != nil && port <= portTo {
//		ss_log.Info("try port=[%v]", port)
//		service = web.NewService(
//			web.Name("go.micro.api.cb"),
//			web.Version("latest"),
//			web.Registry(etcReg),
//			web.Address(net.JoinHostPort(host, strext.ToStringNoPoint(port))),
//		)
//		service.Init()
//
//		// Register Handler
//		service.Handle("/", router.InitRouter())
//
//		// Run service
//		err = service.Run()
//		if err != nil {
//			port++
//		}
//	}
//}
