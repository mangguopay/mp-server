package ss_etcd

import (
	"a.a/cu/ss_log"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"os"
	"strings"
)

// 获取etcd地址信息
func getEtcdAddrFromEnv() []string {
	str := os.Getenv(constants.ETCDAddrEnvName)

	list := []string{}

	for _, v := range strings.Split(str, ",") {
		vv := strings.TrimSpace(v)

		if vv != "" {
			list = append(list, vv)
		}
	}

	return util.UniqueString(list)
}

func GetEtcdAddr() []string {
	// 从环境变量中读取etcd地址
	etcAddrList := getEtcdAddrFromEnv()
	if len(etcAddrList) == 0 {
		panic("请配置环境变量[MP_ETCD_ADDR_LIST]")
	}
	ss_log.Info("etcdAddrList: %v", etcAddrList)
	return etcAddrList
}
