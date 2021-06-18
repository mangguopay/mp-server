package i

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"github.com/micro/go-micro/v2/config"
)

func DoInitBase() (host string, portFrom, portTo int) {
	m := map[string]interface{}{}
	err := config.Get("base", "base").Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}
	z, zErr := time.LoadLocation(strext.ToStringNoPoint(m["timezone"]))
	if zErr != nil {
		panic(fmt.Sprintf("解析时区出错,err: %v", zErr))
	}
	// 设置time包中的默认时区
	time.Local = z

	global.Tz = z
	p := strings.Split(strext.ToStringNoPoint(m["port"]), "-")
	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}

func GetNatsInfo() (host string, port int) {
	clusterId := constants.Nats_Cluster_Primary
	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	return strext.ToStringNoPoint(m["host"]), strext.ToInt(m["port"])
}
