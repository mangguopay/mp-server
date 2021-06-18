package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/cust-srv/dao"
	"time"
)

// 生成服务商对账总计信息
var UserMoneyCountTask = &UserMoneyCount{CronBase{LogCat: "平台用户资金总留存数据统计:", LockExpire: time.Hour * 2}}

type UserMoneyCount struct {
	CronBase
}

// 运行定时任务
func (s *UserMoneyCount) Run() {
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

func (s *UserMoneyCount) Doing() {
	//timeT := time.Now().Format("2006-01-02 15:04:05")
	//获得上一整点的时间，比如 14:15:02 ---> 14:00:00
	timeT := ss_time.GetTimeZeroOclock(time.Now()).Format("2006-01-02 15:04:05")
	s.HandleByTime(timeT)
}

func (s *UserMoneyCount) HandleByTime(time string) {
	ss_log.Info(s.LogCat+"CreateData,time:%s", time)
	// 获取用户总余额
	data := dao.AccDaoInstance.GetUserMoneyCount()
	if err := dao.StatisticUserMoneyDaoInst.Insert(data, time); err != nil {
		ss_log.Error(s.LogCat+"统计用户资金总留存 插入数据操作失败,err: %s", err.Error())
		return
	}

	ss_log.Info(s.LogCat + "CreateData,end")
}
