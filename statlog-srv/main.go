package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/statlog-srv/i"
	"a.a/mp-server/statlog-srv/router"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/etcdv3"
)

const (
	LOG_DIR_NAME = "log_dir"
)

// 日志存放目录
var logDir = flag.String(LOG_DIR_NAME, "../runlog/statlog-srv", "指定日志存放目录")

func main() {
	flag.Parse()
	fmt.Println("logDir:", *logDir)
	// 初始化日志
	ss_log.InitLog(*logDir)
	// 排除参数，防止go-micro内部进行解析报错
	os.Args = ss_func.ExcludeOsArgs(os.Args, LOG_DIR_NAME)

	i.DoInitConfig()

	ss_log.Info("go.micro.srv.statlog")
	host, portFrom, portTo := i.DoInitBase()

	i.DoInitDb()
	defer db.SsInfluxDBInst.Close()
	i.DoInitNats(router.R)
	i.DoInitRedis()
	etcReg := etcdv3.NewRegistry()
	var service micro.Service
	var err error = errors.New("start")
	service = micro.NewService(
		micro.Name("go.micro.srv.statlog"),
		micro.Version("latest"),
		micro.Registry(etcReg),
		micro.WrapClient(i.LogWrap),
		micro.WrapHandler(i.LogWrapperS),
	)
	port := portFrom
	for err != nil && port <= portTo {
		ss_log.Info("try port=[%v]", port)
		service = micro.NewService(
			micro.Address(net.JoinHostPort(host, strext.ToStringNoPoint(port))),
		)
		service.Init()

		// Run service
		err = service.Run()
		if err != nil {
			port++
		}
	}
}
