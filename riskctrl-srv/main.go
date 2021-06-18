package main

import (
	"net"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"a.a/mp-server/riskctrl-srv/handler"
	"a.a/mp-server/riskctrl-srv/i"
	"a.a/mp-server/riskctrl-srv/subcriber"
	"github.com/micro/go-micro/v2/server"
)

func main() {
	servname := constants.ServerNameRiskctrl

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
	//i.DoInitNats(router.R)
	//i.DoInitNatsCli()
	natsHost, natsPort := i.GetPrimaryNatsAddr()
	// 初始化服务
	sa := ss_serv.ServAgent{
		ServName:     servname,
		HttpHost:     host,
		HttpPortList: []int{portFrom, portTo},
		EtcdHosts:    etcAddrList,
		GrpcRegister: func(s server.Server) error {
			// 注册rpc接口
			if err := go_micro_srv_riskctrl.RegisterRiskCtrlHandler(s, new(handler.RiskCtrlHandler)); err != nil {
				return err
			}

			// 注册订阅
			if err := subcriber.InitSubcriber(s); err != nil {
				return err
			}

			return nil
		},
		BrokerType: constants.BrokerTypeNats,
		BrokerAddr: []string{net.JoinHostPort(natsHost, natsPort)},
	}
	sa.RunAsSrv()
}
