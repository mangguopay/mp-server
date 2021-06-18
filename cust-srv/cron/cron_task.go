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

	//暂时停止使用, 当签约模式切换会APP模式的时候可以使用
	//crond.AddFunc("0 0 */4 * * *", func(i ...interface{}) {
	//	// 自动续签即将过期的商户应用
	//	BusinessAppAutoSignedTask.Run()
	//})
	//
	//crond.AddFunc("0 */5 * * * *", func(i ...interface{}) {
	//	// 修改过期签约的状态为过期不可使用
	//	BusinessAppRelieveSignedTask.Run()
	//})

	//todo 因改动签约结构新增的产品签约自动续签
	crond.AddFunc("0 0 */4 * * *", func(i ...interface{}) {
		// 自动续签即将过期的商户产品签约
		BusinessSceneAutoSignedTask.Run()
	})

	//todo 因改动签约结构新增的自动修改过期产品签约状态
	crond.AddFunc("0 */5 * * * *", func(i ...interface{}) {
		// 修改过期签约的状态为过期不可使用
		BusinessAppRelieveSignedTask.Run()
	})

	crond.AddFunc("0 0 2 * * *", func(i ...interface{}) {
		// 生成服务商的对账列表信息
		ServicerCheckListTask.Run()
	})

	crond.AddFunc("0 0 3 * * *", func(i ...interface{}) {
		//crond.AddFunc("*/20 * * * * *", func(i ...interface{}) {
		// 服务商累计统计
		ServicerCountTask.Run()
	})

	// 清楚登录错误统计
	crond.AddFunc("0 0 0 * * *", func(i ...interface{}) {
		//crond.AddFunc("*/30 * * * * *", func(i ...interface{}) {
		// 清楚登录密码错误次数统计
		ErrPwdCountTask.Run()
	})

	// 提现统计
	crond.AddFunc("0 0 3 * * *", func(i ...interface{}) {
		// 提现数据统计
		WithdrawCountTask.Run()
	})
	// 充值统计
	crond.AddFunc("0 0 3 * * *", func(i ...interface{}) {
		// 充值数据统计
		SaveCountTask.Run()
	})
	// 转账统计
	crond.AddFunc("0 0 3 * * *", func(i ...interface{}) {
		// 转账数据统计
		TransferCountTask.Run()
	})
	// 兑换统计
	crond.AddFunc("0 0 3 * * *", func(i ...interface{}) {
		// 兑换数据统计
		ExchangeCountTask.Run()
	})
	// 注册统计和新增服务商统计
	crond.AddFunc("0 0 3 * * *", func(i ...interface{}) {
		// 注册数据统计和新增服务商统计
		RegCountTask.Run()
	})

	// 用户资金总留存统计（每小时统计一次）
	crond.AddFunc("0 5 */1 * * *", func(i ...interface{}) {
		//crond.AddFunc("*/30 * * * * *", func(i ...interface{}) {
		UserMoneyCountTask.Run()
	})

	// 清楚支付密码错误统计
	crond.AddFunc("0 0 0 * * *", func(i ...interface{}) {
		//crond.AddFunc("*/70 * * * * *", func(i ...interface{}) {
		// 清楚支付密码错误统计
		ErrPaymentPwdCountTask.Run()
	})

	//核销码过期修改状态
	crond.AddFunc("0 0 10 * * *", func(i ...interface{}) {
		//crond.AddFunc("*/30 * * * * *", func(i ...interface{}) {
		WriteoffCodeTask.Run()
	})

}

type CronBase struct {
	LogCat     string        // 用来记录日志类型
	Runing     bool          // 任务是否正在运行中
	LockExpire time.Duration // 分布式锁: 有效期
}

// 获取分布式锁的redis键名(当前服务名称+即结构体名称)
func GetLockKey(t interface{}) string {
	s := fmt.Sprintf("%sLock%T", constants.ServerNameCust, t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}
