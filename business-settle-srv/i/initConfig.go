package i

import (
	"context"
	"fmt"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/server"
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

// implements the server.HandlerWrapper
func LogWrapperS(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		ss_log.Info("method|%s|req=[%v]", req.Method(), req.Body())
		defer func() {
			ss_log.Info("method|%s|rsp=[%v]", req.Method(), rsp)
			if r := recover(); r != nil {
				ss_log.Info("recovered")
			}
		}()
		return fn(ctx, req, rsp)
	}
}

type LogWrapper struct {
	client.Client
}

func LogWrap(c client.Client) client.Client {
	return &LogWrapper{c}
}

func DoInitDb() {
	l := []string{constants.DB_CRM}
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
		//alias := strext.ToStringNoPoint(m["alias"])
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

func DoInitNats(r func(*broker.Broker)) {
	m := map[string]interface{}{}
	err := config.Get("mq", "settle").Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	switch m["adapter"] {
	case "nats":
		//ss_log.Info("init nats=[%v:%v@%v] client", m["host"], m["port"], "settle")
		//options := stancf.GetDefaultOptions()
		//options.NatsURL = net.JoinHostPort(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
		//na := stan.NewBroker(stan.ClusterID("settle"), stan.Options(options))
		//if err := na.Init(); err != nil {
		//	panic(fmt.Sprintf("Broker Init error: %v", err))
		//}
		//if err := na.Connect(); err != nil {
		//	panic(fmt.Sprintf("Broker Connect error: %v", err))
		//}
		//r(&na)
	}
}

func GetPrimaryNatsAddr() (host, port string) {
	clusterId := constants.Nats_Cluster_Primary

	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	return strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"])
}
