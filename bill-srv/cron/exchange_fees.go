package cron

import (
	"context"
	"time"

	go_micro_srv_settle "a.a/mp-server/common/proto/settle"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/common"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/common/cache"
)

var ExchangeFeesTask = &ExchangeFees{CronBase{LogCat: "定时任务ExchangeFeesTask:", LockExpire: time.Hour * 2}}

type ExchangeFees struct {
	CronBase
}

// 运行定时任务
func (s *ExchangeFees) Run() {
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

func (s *ExchangeFees) Doing() {
	ss_log.Info(s.LogCat + "开始")

	exchangeResults := dao.ExchangeOrderDaoInst.ExchangeFeesTaskResult()
	if len(exchangeResults) > 0 {
		for _, v := range exchangeResults {
			if v.Fees == "0" {
				continue
			}

			// 发送手续费进MQ
			ev := &go_micro_srv_settle.SettleTransferRequest{
				BillNo:    v.LogNo,
				FeesType:  v.FeesType,
				Fees:      v.Fees,
				MoneyType: v.MoneyType,
			}
			ss_log.Info(s.LogCat+"publishing %+v\n", ev)
			// publish an event
			if err := common.SettleEvent.Publish(context.TODO(), ev); err != nil {
				ss_log.Error(s.LogCat+"err=[定时任务,兑换接口,手续费推送到MQ失败,err----->%s]", err.Error())
				continue
			}
		}
	}

	ss_log.Info(s.LogCat + "结束")
}
