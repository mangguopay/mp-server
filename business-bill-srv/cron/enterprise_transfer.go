package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/business-bill-srv/handler"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"sync"
	"time"
)

var EnterpriseTransferTask = &EnterpriseTransfer{CronBase{LogCat: "企业付款未处理订单订单定时任务:", LockExpire: time.Minute * 30}}

type EnterpriseTransfer struct {
	CronBase
}

func (s *EnterpriseTransfer) Run() {
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

	ss_log.Info(s.LogCat + "企业付款未处理订单订单定时任务-------------------------start\n")

	s.Doing()

	ss_log.Info(s.LogCat + "企业付款未处理订单订单定时任务---------------------------end\n\n")

	// 释放分布式锁
	cache.ReleaseDistributedLock(lockKey, lockValue)
	s.Runing = false
}

func (s *EnterpriseTransfer) Doing() {
	//查询30分钟前的未处理的订单号
	endTime := ss_time.Now(global.Tz).Add(-30 * time.Minute).Format(ss_time.DateTimeDashFormat)
	transferNoList, err := dao.BusinessTransferOrderDaoInst.GetOrderNoList(endTime, constants.BusinessTransferOrderStatusPending, 500)
	if nil != err {
		ss_log.Error(s.LogCat+"修改订单为支付超时失败, err=%v", err)
		return
	}

	if transferNoList == nil {
		return
	}

	s.ParallelDoing(transferNoList)

	ss_log.Info("单号：%v", strext.ToJson(transferNoList))
}

func (s *EnterpriseTransfer) ParallelDoing(transferNoList []string) {
	//并发量
	var parallelNum = 50
	var wg sync.WaitGroup
	for i := 0; i < len(transferNoList); i++ {
		if i > 0 && i%parallelNum == 0 {
			ss_log.Info(s.LogCat+"并发等待----------------------------%v", i)
			wg.Wait()
		}
		wg.Add(1)
		go func(transferNo string) {
			defer wg.Done()

			handler.SyncTransfer(&handler.SyncTransferRequest{
				TransferNo: transferNo,
				Lang:       constants.DefaultLang,
			})
		}(transferNoList[i])
	}
	wg.Wait()
}
