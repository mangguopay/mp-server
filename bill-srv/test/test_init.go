package test

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/i"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_etcd"
	ss_log2 "a.a/mp-server/common/ss_log"
	"a.a/mp-server/common/ss_serv"
	"fmt"
)

func DoInitDb() {
	l := []string{constants.DB_CRM}
	for _, v := range l {
		host := "10.41.1.241"
		port := "5432"
		user := "postgres"
		password := "123"
		name := "mp_crm"
		alias := strext.ToStringNoPoint(v)
		driver := "postgres"
		switch driver {
		case "postgres":
			db.DoDBInitPostgres(alias, host, port, user, password, name)
		default:
			fmt.Printf("not support database|driver=[%v]\n", driver)
		}
	}
}

func DoInitRedis() {
	err := cache.InitRedis("10.41.1.241", "6379", "123456a", 2)
	fmt.Printf("InitRedis, err:%v\n", err)
}

func init() {
	ss_log.InitLog(ss_log2.CheckDevEnvironment(constants.ServerNameBill))
	// 获取etcd地址
	etcAddrList := ss_etcd.GetEtcdAddr()
	ss_serv.DoInitConfigFromEtcd(constants.ETCDPrefix, etcAddrList)
	DoInitDb()
	DoInitRedis()
	i.InitSrv()

}
