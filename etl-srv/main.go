package main

import (
	"flag"
	"fmt"
	"os"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/etl-srv/i"
	"a.a/mp-server/etl-srv/m"
	"a.a/mp-server/etl-srv/task"
)

const (
	LOG_DIR_NAME = "log_dir"
)

// 日志存放目录
var logDir = flag.String(LOG_DIR_NAME, "../runlog/etl-srv", "指定日志存放目录")

func main() {
	flag.Parse()
	fmt.Println("logDir:", *logDir)
	// 初始化日志
	ss_log.InitLog(*logDir)

	// 排除参数，防止go-micro内部进行解析报错
	os.Args = ss_func.ExcludeOsArgs(os.Args, LOG_DIR_NAME)

	i.DoInitConfig()

	ss_log.Info("go.micro.srv.etl")
	_, _, _ = i.DoInitBase()

	i.DoInitDb()
	for _, v := range []m.TaskGroup{task.FetchAccountTaskInst} {
		doTask(v)
	}
}

func doTask(tg m.TaskGroup) {
	// extract
	ctx := tg.Do()
	ctx.Extract.Task.Do(ctx)
	ss_log.Info("extract")

	// transform
	for _, v := range ctx.Transform {
		v.Task.Do(ctx)
		ctx.TransformPc++
	}
	ss_log.Info("transform")

	//load
	for _, v := range ctx.Load {
		v.Task.Do(ctx)
		ctx.LoadPc++
	}
	ss_log.Info("load")
}
