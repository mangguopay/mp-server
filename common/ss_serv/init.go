package ss_serv

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"a.a/mp-server/common/ss_err"

	"a.a/cu/db"
	"a.a/cu/ss_lang"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	http2 "github.com/micro/go-micro/v2/broker/http"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/registry"
	etcd2 "github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/web"
)

func DoInitConfigFromEtcd(prefix string, addrList []string) {
	// 加载配置文件
	if err := config.Load(etcd.NewSource(
		etcd.WithAddress(addrList...),
		// 限定键的前缀，默认是 /micro/config
		etcd.WithPrefix(prefix),
		// optionally strip the provided prefix from the keys, defaults to false
		etcd.StripPrefix(true),
	)); err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}
}

/**
 * 初始化cache
 */
func DoInitCache() {
	m := map[string]interface{}{}
	err := config.Get("cache", "a").Scan(&m)
	if err != nil {
		log.Printf("err=[%v]", err)
		return
	}

	switch m["adapter"] {
	case "redis":
		err := cache.InitRedis(strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]), strext.ToStringNoPoint(m["password"]), strext.ToInt(m["db"]))

		if err != nil {
			panic("初始化redis失败,err:" + err.Error())
		}
		ss_log.Info("reids初始化成功[%s:%s]", strext.ToStringNoPoint(m["host"]), strext.ToStringNoPoint(m["port"]))
	default:
		ss_log.Info("not support cache|driver=[%v]", m["adapter"])
	}
}

/**
 * 初始化数据库
 */
func DoInitDb(l ...string) {
	for _, v := range l {
		m := map[string]interface{}{}
		err := config.Get("database", v).Scan(&m)
		if err != nil {
			ss_log.Error("err=%v", err)
			return
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
			ss_log.Info("not support database|driver=[%v]", driver)
		}
	}
}

// 初始化多语言
func DoInitMultiFromDB(dbName string) {
	dbHandler := db.GetDB(dbName)
	if dbHandler == nil {
		panic("初始化多语言:获取数据库连接失败")
	}
	defer db.PutDB(dbName, dbHandler)

	rows, stmt, err := ss_sql.Query(dbHandler,
		"SELECT key,lang_ch,lang_en,lang_km FROM lang WHERE type=$1 ORDER BY key ASC", "3")
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if err != nil {
		panic("初始化多语言:" + err.Error())
	}

	zh_CNMap := make(map[string]string)
	en_USMap := make(map[string]string)
	km_KHMap := make(map[string]string)

	for rows.Next() {
		var code sql.NullString
		var langCh sql.NullString
		var langEn sql.NullString
		var langKm sql.NullString

		err := rows.Scan(&code, &langCh, &langEn, &langKm)
		if err != nil {
			fmt.Printf("初始化多语言,err:%v\n", err)
			continue
		}
		//fmt.Printf("code:%v, langCh:%v, langEn:%v, langKm:%v \n", code, langCh, langEn, langKm)

		if code.Valid {
			keyStr := code.String

			// 中文语言数据
			if langCh.Valid {
				zh_CNMap[keyStr] = langCh.String
			} else {
				zh_CNMap[keyStr] = ""
			}

			// 英语语言数据
			if langEn.Valid {
				en_USMap[keyStr] = langEn.String
			} else {
				en_USMap[keyStr] = ""
			}

			// 柬埔寨语数据
			if langKm.Valid {
				km_KHMap[keyStr] = langKm.String
			} else {
				km_KHMap[keyStr] = ""
			}
		}
	}

	//fmt.Println(zh_CNMap)
	//fmt.Println(en_USMap)
	//fmt.Println(km_KHMap)

	ss_err.I18nInstance = ss_lang.DoInit()
	ss_err.I18nInstance.RegisteBundle(&ss_lang.SsLangBundle{Name: constants.LangZhCN, Data: &zh_CNMap})
	ss_err.I18nInstance.RegisteBundle(&ss_lang.SsLangBundle{Name: constants.LangEnUS, Data: &en_USMap})
	ss_err.I18nInstance.RegisteBundle(&ss_lang.SsLangBundle{Name: constants.LangKmKH, Data: &km_KHMap})
}

type ServAgent struct {
	ServName     string
	Router       http.Handler
	HttpHost     string
	HttpPortList []int
	EtcdHosts    []string
	Version      string // default=latest
	GrpcRegister func(s server.Server) error
	BrokerType   string
	BrokerAddr   []string
}

func (r ServAgent) RunAsWeb() {
	// 初始化registry
	ss_log.Info("etcdhosts=[%v]", r.EtcdHosts)
	etcReg := etcd2.NewRegistry(func(op *registry.Options) {
		op.Addrs = r.EtcdHosts
	})

	var service web.Service
	port := r.HttpPortList[0]
	portTo := r.HttpPortList[1]
	if r.Version == "" {
		r.Version = "latest"
	}
	if r.HttpHost == "" {
		r.HttpHost = "0.0.0.0"
	}
	for port <= portTo {
		ss_log.Info("try port=[%v]", port)
		addr := net.JoinHostPort(r.HttpHost, strext.ToStringNoPoint(port))
		service = web.NewService(
			web.Name(r.ServName),
			web.Version(r.Version),
			web.Registry(etcReg),
			web.Address(addr),
			web.Handler(r.Router),
		)
		err2 := service.Init()
		if err2 != nil {
			ss_log.Error("err|service.Init=[%v]", err2)
			port++
			continue
		}

		// Run service
		err := service.Run()
		if err != nil {
			port++
		}
	}
}

func (r ServAgent) RunAsSrv() {
	// 初始化registry
	etcReg := etcd2.NewRegistry(func(op *registry.Options) {
		op.Addrs = r.EtcdHosts
	})

	port := r.HttpPortList[0]
	portTo := r.HttpPortList[1]
	if r.Version == "" {
		r.Version = "latest"
	}
	if r.HttpHost == "" {
		r.HttpHost = "0.0.0.0"
	}
	var br broker.Broker
	switch r.BrokerType {
	case constants.BrokerTypeHTTP:
		br = http2.NewBroker()
	case constants.BrokerTypeNats:
		fallthrough
	default:
		br = nats.NewBroker()
		br.Init(broker.Addrs(r.BrokerAddr...))
	}

	service := micro.NewService(
		micro.Name(r.ServName),
		micro.Version(r.Version),
		micro.Registry(etcReg),
		micro.WrapClient(LogWrap),
		micro.WrapHandler(logWrapperS),
		micro.Broker(br),
	)
	fNeedInitGrpc := true
	for port <= portTo {
		ss_log.Info("try port=[%v]", port)
		addr := net.JoinHostPort(r.HttpHost, strext.ToStringNoPoint(port))
		service = micro.NewService(
			micro.Address(addr),
		)
		service.Init()
		global.Port = strext.ToStringNoPoint(port)
		// Register Handler
		if fNeedInitGrpc {
			err := r.GrpcRegister(service.Server())
			fNeedInitGrpc = false
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}
		}

		// Run service
		err := service.Run()
		if err != nil {
			port++
		}
	}
}

//==========================================================log
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ss_log.Info("[rpc_cli]send=>[%s][%v]", req.Method(), req.Body())
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			ss_log.Info("[rpc_cli]%s|recovered|r=%+v", req.Method(), r)
		}
	}()
	// 调用options
	var opss []client.CallOption
	opss = append(opss, func(o *client.CallOptions) {
		o.RequestTimeout = time.Second * 30
		o.DialTimeout = time.Second * 30
	})
	opss = append(opss, opts...)

	tb := time.Now()
	err := l.Client.Call(ctx, req, rsp, opss...)
	te := time.Now()
	diff := te.Sub(tb).Milliseconds()
	if err != nil {
		ss_log.Error("[rpc_cli]recv<=|%s|err=[%v]|cost=[%v]ms", req.Method(), err, diff)
	} else {
		ss_log.Info("[rpc_cli]recv<=|%s|[%v]|cost=[%v]ms", req.Method(), rsp, diff)
	}
	return err
}

// implements client.Wrapper as logWrapper
func LogWrap(c client.Client) client.Client {
	return &logWrapper{c}
}

// implements the server.HandlerWrapper
func logWrapperS(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		tb := time.Now()
		ss_log.Info("[rpc_srv]send<=|%s|[%v]", req.Method(), req.Body())
		defer func() {
			te := time.Now()
			diff := te.Sub(tb).Milliseconds()
			ss_log.Info("[rpc_srv]recv=>|%s|[%v]|cost=[%v]ms", req.Method(), rsp, diff)
			if r := recover(); r != nil {
				debug.PrintStack()
				ss_log.Error("[rpc_srv]%s|recovered:%+v", req.Method(), r)
			}
		}()
		return fn(ctx, req, rsp)
	}
}

//==========================================================log
func GetRpcDefClient() client.Client {
	c := client.DefaultClient
	c = LogWrap(c)
	return c
}
