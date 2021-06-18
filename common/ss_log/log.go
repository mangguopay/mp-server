package ss_log

import (
	"flag"
	"log"
	"runtime"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
)

const (
	LOG_DIR_NAME = "log_dir"

	// 开发环境固定日志存放目录
	DevLogDir = "/var/logs"
)

/**
 * 初始化日志
 */
func DoInit(app string) {
	// 开发环境固定日志目录
	app = CheckDevEnvironment(app)

	// 解析日志目录
	var logDir = flag.String(LOG_DIR_NAME, app, "指定日志存放目录")
	flag.Parse()

	log.Println("logDir:", *logDir)

	// 初始化日志
	ss_log.InitLog(*logDir)

	// 排除参数，防止go-micro内部进行解析报错
	strext.ExcludeFlags(LOG_DIR_NAME)
}

// 检测是否是开发环境，是开发环境重新设置目录
func CheckDevEnvironment(appName string) string {
	if runtime.GOOS == "windows" {
		return DevLogDir + "/" + appName
	}

	return appName
}
