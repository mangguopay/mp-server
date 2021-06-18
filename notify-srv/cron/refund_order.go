package cron

import (
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"sync"
	"time"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/notify-srv/common"
	"a.a/mp-server/notify-srv/dao"
)

var RefundNotifyTask = &RefundNotify{CronBase{LogCat: "退款订单-异步通知定时任务:", LockExpire: time.Minute * 5}}

type RefundNotify struct {
	CronBase
}

// 运行定时任务
func (n *RefundNotify) Run() {
	if n.Runing { // 正在运行中
		return
	}

	lockKey := GetLockKey(n)
	lockValue := strext.NewUUID()

	// 获取分布式锁
	if !cache.GetDistributedLock(lockKey, lockValue, n.LockExpire) {
		return
	}
	n.Runing = true

	n.Doing()

	// 释放分布式锁
	cache.ReleaseDistributedLock(lockKey, lockValue)
	n.Runing = false
}

func (n *RefundNotify) Doing() {
	n.QueryNotifyOmissionOrder()
	n.QueryNotifyBreakOrder()
}

func (n *RefundNotify) QueryNotifyOmissionOrder() {
	ss_log.Info(n.LogCat + "定时检查通知遗漏的退款订单---------------------------------start")

	//查询通知遗漏订单,checkTime=15s,
	//第一次通知后等15秒才会进行一下次通知, 如果pay_time < currentTime - 15则第一次通知失败,要再通知一次
	checkTime := common.GetNotifyWaitTimeById(common.NotifyWaitTimesOne)
	finishTime := ss_time.Now(global.Tz).Add(time.Duration(-checkTime) * time.Second).Format(ss_time.DateTimeDashFormat)
	omissionOrders, err := dao.BusinessRefundOrderDaoInst.QueryNotifyOmission(constants.BusinessRefundStatusSuccess,
		constants.NotifyStatusNOT, finishTime)
	if err != nil {
		ss_log.Error(n.LogCat+"查询退款订单失败, err=[%v]", err)
		return
	}
	omissionOrderNum := len(omissionOrders)
	ss_log.Info(n.LogCat+"通知遗漏的退款订单数量:%v,omissionOrders=%v", omissionOrderNum, omissionOrders)
	//并发处理
	n.ParallelDoing(omissionOrders, omissionOrderNum)

	ss_log.Info(n.LogCat + "定时检查通知遗漏的退款订单------------------------- ------------end")
}

func (n *RefundNotify) QueryNotifyBreakOrder() {
	ss_log.Info(n.LogCat + "定时检查通知中断的退款订单---------------------------------start")

	//查询通知中断订单
	nextNotifyTime := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	breakOrders, err := dao.BusinessRefundOrderDaoInst.QueryNotifyBreak(constants.BusinessRefundStatusSuccess,
		constants.NotifyStatusDoing, nextNotifyTime)
	if err != nil {
		ss_log.Error(n.LogCat+"查询通知中断的退款订单, err=[%v]", err)
		return
	}
	breakOrderNum := len(breakOrders)
	ss_log.Info(n.LogCat+"中断通知的退款订单数量:%v, breakOrders=%v", breakOrderNum, breakOrders)
	if breakOrderNum <= 0 {
		return
	}

	//并发处理
	n.ParallelDoing(breakOrders, breakOrderNum)

	ss_log.Info(n.LogCat + "定时检查通知中断的退款订单------------------------- ------------end")
}

func (n *RefundNotify) ParallelDoing(orders []string, orderNum int) {
	parallelNum := 50 // 并发执行数量
	var wg sync.WaitGroup
	for i := 0; i < orderNum; i++ {
		if i > 0 && i%parallelNum == 0 {
			ss_log.Info(n.LogCat+"并发等待----------------------------%v", i)
			wg.Wait()
		}
		wg.Add(1)
		go func(orderNo string) {
			defer wg.Done()
			n.InsertRedis(orderNo)
		}(orders[i])
	}
	wg.Wait()
}

func (n *RefundNotify) InsertRedis(refundNo string) {
	err := cache.RedisClient.SetNX(common.GetRefundNotifyExpireKey(refundNo), common.NotifyWaitTimesOne, time.Second).Err()
	if err != nil {
		ss_log.Error(n.LogCat+"redis插入数据失败, key=%v, err=%v", common.GetTransferNotifyExpireKey(refundNo), err)
	}
}
