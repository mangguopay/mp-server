package main

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	protoTm "a.a/mp-server/common/proto/tm"
	"a.a/mp-server/common/ss_etcd"
	"a.a/mp-server/common/ss_func"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"a.a/mp-server/tm-srv/common"
	"a.a/mp-server/tm-srv/handler"
	"a.a/mp-server/tm-srv/i"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameTm

	// 初始化日志
	ss_log2.DoInit(servname)
	ss_log.Info(servname)

	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)

	// 初始化基础配置信息
	host, portFrom, portTo := i.DoInitBase()
	ss_serv.DoInitDb(constants.DB_CRM)

	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			common.TmServerFullid = ss_func.GetServerFullId(s.Options().Name, s.Options().Id)
			return protoTm.RegisterTmHandler(s, &handler.TmHandler{Server: s})
		},
	}
	sa.RunAsSrv()
}
