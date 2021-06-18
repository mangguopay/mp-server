package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"time"
)

// 生成服务商对账总计信息
var ExchangeCountTask = &ExchangeCount{CronBase{LogCat: "兑换统计定时任务:", LockExpire: time.Hour * 2}}

type ExchangeCount struct {
	CronBase
}

// 运行定时任务
func (s *ExchangeCount) Run() {
	if s.Runing { // 正在运行中
		return
	}

	lockKey := GetLockKey(s)
	lockValue := strext.NewUUID()

	// 获取分布式锁
	if !cache.GetDistributedLock(lockKey, lockValue, s.LockExpire) {
		return
	}
	s.Runing = true

	s.Doing()

	// 释放分布式锁
	cache.ReleaseDistributedLock(lockKey, lockValue)
	s.Runing = false
}

func (s *ExchangeCount) Doing() {
	date := ss_time.GetDayBefore() // 获取昨天的日期
	// log_cust_to_headquarters
	s.HandleByDate(date)
}

func (s *ExchangeCount) HandleByDate(date string) {
	ss_log.Info(s.LogCat+"HandleByDate,date:%s", date)
	data, err := dao.ExchangeOrderDaoInst.GetExchangeCountByDate(date)
	if err != nil && err.Error() != ss_sql.DB_NO_ROWS_MSG {
		ss_log.Error(s.LogCat+"兑换统计操作失败,err: %s", err.Error())
		return
	}
	// 插入数据
	if data != nil {
		if err := dao.StatisticUserExchangeDaoInst.Insert(data); err != nil {
			ss_log.Error(s.LogCat+"兑换统计插入数据操作失败,err: %s", err.Error())
		}
	}
	ss_log.Info(s.LogCat + "HandleByDate,end")
}
