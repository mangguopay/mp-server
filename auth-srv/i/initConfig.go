package i

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/auth-srv/common"
	"a.a/mp-server/common/global"
	"github.com/micro/go-micro/v2/config"
)

//func DoInitDb() {
//	l := []string{constants.DB_CRM}
//	for _, v := range l {
//		m := map[string]interface{}{}
//		err := config.Get("database", v).Scan(&m)
//		if err != nil {
//			ss_log.Error("err=%v", err)
//		}
//
//		host := strext.ToStringNoPoint(m["host"])
//		port := strext.ToStringNoPoint(m["port"])
//		user := strext.ToStringNoPoint(m["user"])
//		password := strext.ToStringNoPoint(m["password"])
//		name := strext.ToStringNoPoint(m["name"])
//		//alias := strext.ToStringNoPoint(m["alias"])
//		alias := strext.ToStringNoPoint(v)
//		driver := strext.ToStringNoPoint(m["driver"])
//		switch driver {
//		case "postgres":
//			db.DoDBInitPostgres(alias, host, port, user, password, name)
//		default:
//			fmt.Printf("not support database|driver=[%v]\n", driver)
//		}
//	}
//}

//var EncryptMap map[string]interface{}

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

	common.EncryptMap = make(map[string]interface{})
	common.EncryptMap["back_pri_key"] = m["pri_key"]
	common.EncryptMap["back_pub_key"] = m["pub_key"]
	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}
