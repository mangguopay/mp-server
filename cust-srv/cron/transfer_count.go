package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"time"
)

// 生成服务商对账总计信息
var TransferCountTask = &TransferCount{CronBase{LogCat: "转账统计定时任务:", LockExpire: time.Hour * 2}}

type TransferCount struct {
	CronBase
}

// 运行定时任务
func (s *TransferCount) Run() {
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

func (s *TransferCount) Doing() {
	date := ss_time.GetDayBefore() // 获取昨天的日期
	s.HandleByDate(date)
}

func (s *TransferCount) HandleByDate(date string) {
	ss_log.Info(s.LogCat+"CreateData,date:%s", date)
	for _, currencyType := range []string{constants.CURRENCY_USD, constants.CURRENCY_KHR} {
		// log_cust_to_headquarters
		data, err := dao.TransferDaoInst.GetTransferCountByDate(date, currencyType)
		if err != nil && err.Error() != ss_sql.DB_NO_ROWS_MSG {
			ss_log.Error(s.LogCat+"统计 %s 转账统计操作失败,err: %s", currencyType, err.Error())
			continue
		}
		// 插入数据
		if data != nil {
			if err := dao.StatisticUserTransferDaoInst.Insert(data); err != nil {
				ss_log.Error(s.LogCat+"统计 %s 转账统计插入数据操作失败,err: %s", currencyType, err.Error())
			}
		}
	}
	ss_log.Info(s.LogCat + "CreateData,end")
}
