package test

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	"a.a/mp-server/common/ss_serv"
	"path"
)

func init() {
	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)
	ss_serv.DoInitDb(constants.DB_CRM)
	ss_serv.DoInitCache()
	ss_log.InitLog(path.Join("runlog", "cust-srv", "logs"))
}
