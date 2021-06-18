package main

import (
	"fmt"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_gis "a.a/mp-server/common/proto/gis"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"a.a/mp-server/gis-srv/handler"
	"a.a/mp-server/gis-srv/i"
	"a.a/mp-server/gis-srv/subcriber"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameGis

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

	natsHost, natsPort := i.GetNatsInfo()
	ss_log.Info("natsHost-----> %s,natsPort----->%d", natsHost, natsPort)
	i.InitGisInfo()

	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {

			// 注册rpc
			if err := go_micro_srv_gis.RegisterGISHandler(s, new(handler.GisHandler)); err != nil {
				return err
			}

			// 注册订阅
			if err := subcriber.InitSubcriber(s); err != nil {
				return err
			}

			return nil
		},
		BrokerType: constants.BrokerTypeNats,
		BrokerAddr: []string{
			fmt.Sprintf("%s:%s", natsHost, strext.ToStringNoPoint(natsPort)),
		},
	}
	sa.RunAsSrv()
}
