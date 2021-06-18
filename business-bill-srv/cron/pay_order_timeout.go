package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/cache"
	"time"
)

var PayOrderTimeoutTask = &PayOrderTimeout{CronBase{LogCat: "支付超时订单定时任务:", LockExpire: time.Minute * 30}}

type PayOrderTimeout struct {
	CronBase
}

func (s *PayOrderTimeout) Run() {
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

func (s *PayOrderTimeout) Doing() {
	ss_log.Info(s.LogCat + "支付超时订单定时任务-------------------------start\n")
	//修改订单状态为待支付，且支付过期时间到期的订单状态为超时
	if err := dao.BusinessBillDaoInst.UpdateOrderOutTime(); nil != err {
		ss_log.Error(s.LogCat+"修改订单为支付超时失败, err=%v", err)
	}
	ss_log.Info(s.LogCat + "支付超时订单定时任务---------------------------end\n\n")
}
