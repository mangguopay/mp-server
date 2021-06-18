package cron

import (
	"a.a/cu/crond"
	"a.a/mp-server/common/constants"
	"fmt"
	"strings"
	"time"
)

func DoInitCronTask() {
	crond.DoInitCrondWithoutFunc()
	//隔10秒执行一次
	crond.AddFunc("*/10 * * * * *", func(i ...interface{}) {
		go NotifyPaymentSuccessTask.Run()
	})

	//隔10秒执行一次
	crond.AddFunc("*/10 * * * * *", func(i ...interface{}) {
		go TransferNotifyTask.Run()
	})

	//隔10秒执行一次
	crond.AddFunc("*/10 * * * * *", func(i ...interface{}) {
		go RefundNotifyTask.Run()
	})

}

type CronBase struct {
	LogCat     string        // 用来记录日志类型
	Runing     bool          // 任务是否正在运行中
	LockExpire time.Duration // 分布式锁: 有效期
}

// 获取分布式锁的redis键名(当前服务名称+即结构体名称)
func GetLockKey(t interface{}) string {
	s := fmt.Sprintf("%sLock%T", constants.ServerNamePayNotify, t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}
