package i

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/cache"
	"strings"
	"time"

	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"github.com/micro/go-micro/v2/config"
)

func DoInitBase() (host string, portFrom, portTo int) {
	m := map[string]interface{}{}
	err := config.Get("base", "base").Scan(&m)
	if err != nil {
		panic("读取数据库配置失败,err:" + err.Error())
		return
	}
	z, _ := time.LoadLocation(strext.ToStringNoPoint(m["timezone"]))
	global.Tz = z
	p := strings.Split(strext.ToStringNoPoint(m["port"]), "-")
	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}

func GetPrimaryNatsAddr() (host, port string) {
	clusterId := constants.Nats_Cluster_Primary

	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		panic("读取Nats配置失败,err:" + err.Error())
		return
	}

	return strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"])
}

/**
 * 初始化cache
 */
func DoInitCache() {
	m := map[string]interface{}{}
	err := config.Get("cache", "a").Scan(&m)
	if err != nil {
		panic("读取cache配置失败,err:" + err.Error())
		return
	}

	switch m["adapter"] {
	case "redis":
		err := cache.InitRedis(strext.ToStringNoPoint(m["host"]),
			strext.ToStringNoPoint(m["port"]),
			strext.ToStringNoPoint(m["password"]),
			constants.NotifySrvRedisDb)

		if err != nil {
			panic("初始化redis失败,err:" + err.Error())
			return
		}
		ss_log.Info("reids初始化成功[%s:%s]", strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
	default:
		ss_log.Info("not support cache|driver=[%v]", m["adapter"])
	}
}
