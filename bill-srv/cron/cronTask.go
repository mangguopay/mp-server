package cron

import (
	"fmt"
	"strings"
	"time"

	"a.a/mp-server/common/constants"

	"a.a/cu/crond"
)

func DoInitCronTask() {
	crond.DoInitCrondWithoutFunc()

	// 定时查询漏掉统计手续费的订单
	crond.AddFunc("0 0 * * * ?", func(i ...interface{}) {
		ExchangeFeesTask.Run()
	})

	crond.AddFunc("0 0 * * * ?", func(i ...interface{}) {
		CollectionTask.Run()
	})

	crond.AddFunc("0 0 * * * ?", func(i ...interface{}) {
		SaveMoneyFeesTask.Run()
	})

	crond.AddFunc("0 0 * * * ?", func(i ...interface{}) {
		TransferFeesTask.Run()
	})

	crond.AddFunc("0 0 * * * ?", func(i ...interface{}) {
		WithdrawFeesTask.Run()
	})

	crond.AddFunc("0 */5 * * * *", func(i ...interface{}) {
		// 批量转账任务
		BusinessBatchTransferTask.Run()
	})
}

type CronBase struct {
	LogCat     string        // 用来记录日志类型
	Runing     bool          // 任务是否正在运行中
	LockExpire time.Duration // 分布式锁: 有效期
}

// 获取分布式锁的redis键名(当前服务名称+即结构体名称)
func GetLockKey(t interface{}) string {
	s := fmt.Sprintf("%sLock%T", constants.ServerNameBill, t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}
