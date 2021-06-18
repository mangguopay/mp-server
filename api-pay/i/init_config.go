package i

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/global"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/etcd"
)

func DoInitConfig() {
	// 加载配置文件
	if err := config.Load(etcd.NewSource(
		//etcd.WithAddress("127.0.0.1:2379"),
		// 限定键的前缀，默认是 /micro/config
		etcd.WithPrefix("/micro/config/payserv"),
		// optionally strip the provided prefix from the keys, defaults to false
		etcd.StripPrefix(true),
	)); err != nil {
		fmt.Println(err)
		return
	}
}

func DoInitDb() {
	l := []string{"crm"}
	for _, v := range l {
		m := map[string]interface{}{}
		err := config.Get("database", v).Scan(&m)
		if err != nil {
			ss_log.Error("err=%v", err)
		}

		host := strext.ToStringNoPoint(m["host"])
		port := strext.ToStringNoPoint(m["port"])
		user := strext.ToStringNoPoint(m["user"])
		password := strext.ToStringNoPoint(m["password"])
		name := strext.ToStringNoPoint(m["name"])
		alias := strext.ToStringNoPoint(v)
		driver := strext.ToStringNoPoint(m["driver"])
		switch driver {
		case "postgres":
			db.DoDBInitPostgres(alias, host, port, user, password, name)
		default:
			fmt.Printf("not support database|driver=[%v]\n", driver)
		}
	}
}

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

	ss_log.SetShowSql(strext.ToBool(m["show_sql"]))

	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}
