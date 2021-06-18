package i

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_etcd"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-plugins/broker/stan"
	stancf "github.com/nats-io/stan.go"
)

func DoInitConfig() {
	// 从环境变量中读取etcd地址
	addrList := ss_etcd.GetEtcdAddrFromEnv()
	if len(addrList) == 0 {
		panic("从环境变量中读取不到etcd地址")
	}
	ss_log.Info("etcdAddrList: %v", addrList)
	// 加载配置文件
	if err := config.Load(etcd.NewSource(
		etcd.WithAddress(addrList...),
		//etcd.WithAddress("127.0.0.1:2379"),
		// 限定键的前缀，默认是 /micro/config
		etcd.WithPrefix("/micro/config/mp-server"),
		// optionally strip the provided prefix from the keys, defaults to false
		etcd.StripPrefix(true),
	)); err != nil {
		fmt.Println(err)
		return
	}
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

func DoInitRedis() {
	m := map[string]interface{}{}
	config.Get("cache", "a").Scan(&m)
	switch m["adapter"] {
	case "redis":
		err := cache.InitRedis(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]), strext.ToStringNoPoint(m["password"]))
		if err != nil {
			panic("初始化redis失败,err:" + err.Error())
		}
		ss_log.Info("reids初始化成功[%s:%s]", strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
	default:
		fmt.Printf("not support cache|driver=[%v]\n", m["adapter"])
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
	return strext.ToStringNoPoint(m["host"]), strext.ToInt(p[0]), strext.ToInt(p[1])
}

func DoInitDb() {
	l := []string{constants.DB_CRM, constants.DbStat}
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
		case "influxdb":
			db.SsInfluxDBInst.DoInit()
			db.SsInfluxDBInst.InitClient(host, port, "0", alias)
		default:
			fmt.Printf("not support database|driver=[%v]\n", driver)
		}
	}
}

func DoInitNats(r func(*broker.Broker)) {
	clusterId := "etl"
	m := map[string]interface{}{}
	err := config.Get("mq", clusterId).Scan(&m)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	switch m["adapter"] {
	case "nats":
		ss_log.Info("init nats=[%v:%v@%v] client", m["host"], m["port"], clusterId)
		options := stancf.GetDefaultOptions()
		options.NatsURL = net.JoinHostPort(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
		na := stan.NewBroker(stan.ClusterID(clusterId), stan.Options(options))
		if err := na.Init(); err != nil {
			panic(fmt.Sprintf("Broker Init error: %v", err))
		}
		if err := na.Connect(); err != nil {
			panic(fmt.Sprintf("Broker Connect error: %v", err))
		}
		r(&na)
	}
}
