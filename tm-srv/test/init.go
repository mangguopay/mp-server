package test

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	"a.a/mp-server/common/ss_serv"
	"a.a/mp-server/tm-srv/i"
)

func init() {
	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)

	i.DoInitBase()

	// 初始化基础配置信息
	ss_serv.DoInitDb(constants.DB_CRM)
}
