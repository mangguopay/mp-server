package main

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	"a.a/mp-server/common/ss_etcd"
	"a.a/mp-server/common/ss_func"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"a.a/mp-server/quota-srv/common"
	"a.a/mp-server/quota-srv/handler"
	"a.a/mp-server/quota-srv/i"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameQuota

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
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			common.QuotaServerFullid = ss_func.GetServerFullId(s.Options().Name, s.Options().Id)
			return go_micro_srv_quota.RegisterQuotaHandler(s, new(handler.QuotaHandler))
		},
	}
	sa.RunAsSrv()
}
