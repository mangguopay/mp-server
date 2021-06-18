package test

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"

	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
)

func init() {
	ss_log.InitLog(ss_log2.CheckDevEnvironment(constants.ServerNameApiPay))

	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)
	ss_serv.DoInitDb(constants.DB_CRM)
	ss_serv.DoInitCache()
}
