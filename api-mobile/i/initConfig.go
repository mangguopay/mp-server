package i

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/common"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_net"
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
	// 初始化密钥组
	common.EncryptMap = make(map[string]interface{})
	common.EncryptMap[common.SKEY_PlatPri] = m["pri_key"]
	common.EncryptMap[common.SKEY_PlatPub] = m["pub_key"]
	common.EncryptMap[common.SKEY_DefPri] = m["default_pri_key"]
	common.EncryptMap[common.SKEY_DefPub] = m["default_pub_key"]

	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}

func DoInitNatsCli() {
	l := []*ss_net.NatsCliConf{
		{
			ClusterId: constants.Nats_Cluster_Secondary,
			Topic:     constants.Nats_Topic_Statlog,
		},
	}
	ss_net.DoInitNatsCli(l)
}
