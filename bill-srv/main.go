package main

import (
	"net"

	billCommon "a.a/mp-server/bill-srv/common"
	"a.a/mp-server/common/ss_func"

	go_micro_srv_bill "a.a/mp-server/common/proto/bill"

	"a.a/cu/ss_log"
	"a.a/mp-server/bill-srv/cron"
	"a.a/mp-server/bill-srv/handler"
	"a.a/mp-server/bill-srv/i"
	"a.a/mp-server/bill-srv/subcriber"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameBill

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
	i.InitAwss3()

	go cron.DoInitCronTask()

	natsHost, natsPort := i.GetPrimaryNatsAddr()
	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			billCommon.BillServerFullId = ss_func.GetServerFullId(s.Options().Name, s.Options().Id)

			// 监听订阅
			if err := subcriber.InitSubcriber(s); err != nil {
				return err
			}

			// 监听服务
			if err := go_micro_srv_bill.RegisterBillHandler(s, new(handler.BillHandler)); err != nil {
				return err
			}

			return nil
		},
		BrokerType: constants.BrokerTypeNats,
		BrokerAddr: []string{net.JoinHostPort(natsHost, natsPort)},
	}
	sa.RunAsSrv()
}
