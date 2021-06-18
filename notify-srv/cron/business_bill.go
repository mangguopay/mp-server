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

var NotifyPaymentSuccessTask = &NotifyPaymentSuccessRet{CronBase{LogCat: "交易订单-异步通知定时任务:", LockExpire: time.Minute * 5}}

type NotifyPaymentSuccessRet struct {
	CronBase
}

// 运行定时任务
func (n *NotifyPaymentSuccessRet) Run() {
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

func (n *NotifyPaymentSuccessRet) Doing() {
	n.QueryNotifyOmissionOrder()
	n.QueryNotifyBreakOrder()
}

func (n *NotifyPaymentSuccessRet) QueryNotifyOmissionOrder() {
	ss_log.Info(n.LogCat + "定时检查通知遗漏交易订单---------------------------------start")

	//查询通知遗漏订单,checkTime=15s,
	//第一次通知后等15秒才会进行一下次通知, 如果pay_time < currentTime - 15则第一次通知失败,要再通知一次
	checkTime := common.GetNotifyWaitTimeById(common.NotifyWaitTimesOne)
	payTime := ss_time.Now(global.Tz).Add(time.Duration(-checkTime) * time.Second).Format(ss_time.DateTimeDashFormat)
	omissionOrders, err := dao.BusinessBillDaoInst.QueryNotifyOmission(constants.BusinessOrderStatusPay, constants.NotifyStatusNOT, payTime)
	if err != nil {
		ss_log.Error(n.LogCat+"查询交易订单失败, err=[%v]", err)
		return
	}
	omissionOrderNum := len(omissionOrders)
	ss_log.Info(n.LogCat+"遗漏通知交易订单数量:%v,omissionOrders=%v", omissionOrderNum, omissionOrders)
	//并发处理
	n.ParallelDoing(omissionOrders, omissionOrderNum)

	ss_log.Info(n.LogCat + "定时检查通知遗漏交易订单------------------------- ------------end")
}

func (n *NotifyPaymentSuccessRet) QueryNotifyBreakOrder() {
	ss_log.Info(n.LogCat + "定时检查通知中断交易订单---------------------------------start")

	//查询通知中断订单
	nextNotifyTime := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	breakOrders, err := dao.BusinessBillDaoInst.QueryNotifyBreak(constants.BusinessOrderStatusPay, constants.NotifyStatusDoing, nextNotifyTime)
	if err != nil {
		ss_log.Error(n.LogCat+"查询通知中断交易订单, err=[%v]", err)
		return
	}
	breakOrderNum := len(breakOrders)
	ss_log.Info(n.LogCat+"中断通知交易订单数量:%v, breakOrders=%v", breakOrderNum, breakOrders)
	if breakOrderNum <= 0 {
		return
	}

	//并发处理
	n.ParallelDoing(breakOrders, breakOrderNum)

	ss_log.Info(n.LogCat + "定时检查通知中断交易订单------------------------- ------------end")
}

func (n *NotifyPaymentSuccessRet) ParallelDoing(orders []string, orderNum int) {
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

func (n *NotifyPaymentSuccessRet) InsertRedis(orderNo string) {
	err := cache.RedisClient.SetNX(common.GetPayNotifyExpireKey(orderNo), common.NotifyWaitTimesOne, time.Second).Err()
	if err != nil {
		ss_log.Error(n.LogCat+"redis插入数据失败, key=%v, err=%v", common.GetPayNotifyExpireKey(orderNo), err)
	}
}
