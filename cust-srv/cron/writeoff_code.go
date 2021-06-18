package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/cust-srv/dao"
	"time"
)

// 生成服务商对账总计信息
var WriteoffCodeTask = &writeoffCode{CronBase{LogCat: "过期核销码处理:", LockExpire: time.Hour * 2}}

type writeoffCode struct {
	CronBase
}

// 运行定时任务
func (s *writeoffCode) Run() {
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

func (s *writeoffCode) Doing() {
	s.HandleByTime()
}

func (s *writeoffCode) HandleByTime() {
	ss_log.Info(s.LogCat+"过期核销码处理开始。time:%s", time.Now())

	//确认有过期的核销码
	codeArr := dao.WriteoffDaoInst.GetExpiredCodeArr()
	if codeArr != nil {
		//修改核销码的状态
		for _, code := range codeArr {
			if err := dao.WriteoffDaoInst.UpdateExpiredCodeStatus(code, constants.WriteOffCodeExpired); err != nil {
				ss_log.Error("修改过期核销码[%v]的状态为过期失败,err[%v]", code, err)
			}
		}
	} else {
		ss_log.Info("无过期的核销码需处理")
	}

	ss_log.Info(s.LogCat + "过期核销码处理结束，end")
}
