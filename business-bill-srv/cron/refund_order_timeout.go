package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/cache"
	"time"
)

var RefundOrderTimeoutTask = &RefundOrderTimeout{CronBase{LogCat: "退款超时订单定时任务:", LockExpire: time.Minute * 30}}

type RefundOrderTimeout struct {
	CronBase
}

func (s *RefundOrderTimeout) Run() {
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

func (s *RefundOrderTimeout) Doing() {
	ss_log.Info(s.LogCat + "退款超时订单定时任务-------------------------start\n")

	//修改退款订单订单状态为处理中(0)，且超时的订单为退款失败
	list, err := dao.BusinessRefundOrderDaoInst.UpdateTimeOutOrderStatus()
	if err != nil {
		ss_log.Error(s.LogCat+"修改订单为支付超时失败, err=%v", err)
		return
	}

	//测试用
	if list != nil {
		//ss_log.Info(s.LogCat+"退款超时订单：%v", strext.ToJson(list))
	}

	ss_log.Info(s.LogCat + "退款超时订单定时任务---------------------------end\n\n")
}
