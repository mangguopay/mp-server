package cron

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/crond"
	"a.a/mp-server/common/constants"
)

func DoInitCronTask() {
	crond.DoInitCrondWithoutFunc()

	//每天6点执行一次
	crond.AddFunc("0 0 6 * * *", func(i ...interface{}) {
		BusinessSettleTask.Run()
	})

	//修改支付超时的交易订单
	crond.AddFunc("0 */10 * * * *", func(i ...interface{}) {
		PayOrderTimeoutTask.Run()
	})

	//修改退款超时的退款订单
	crond.AddFunc("0 */3 * * * *", func(i ...interface{}) {
		RefundOrderTimeoutTask.Run()
	})

	//企业付款订单
	crond.AddFunc("0 */15 * * * *", func(i ...interface{}) {
		EnterpriseTransferTask.Run()
	})
}

type CronBase struct {
	LogCat     string        // 用来记录日志类型
	Runing     bool          // 任务是否正在运行中
	LockExpire time.Duration // 分布式锁: 有效期
}

// 获取分布式锁的redis键名(当前服务名称+即结构体名称)
func GetLockKey(t interface{}) string {
	s := fmt.Sprintf("%sLock%T", constants.ServerNameBusinessBill, t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}
