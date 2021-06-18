package main

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/admin-auth-srv/handler"
	"a.a/mp-server/admin-auth-srv/i"
	"a.a/mp-server/common/constants"
	adminAuth "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameAdminAuth

	// 初始化日志
	ss_log2.DoInit(servname)
	ss_log.Info(servname)

	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)

	// 初始化基础配置信息
	host, portFrom, portTo := i.DoInitBase()
	ss_serv.DoInitCache()
	ss_serv.DoInitDb("crm")
	i.InitSrv()

	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			return adminAuth.RegisterAdminAuthHandler(s, new(handler.AdminAuth))
		},
	}
	sa.RunAsSrv()
}
