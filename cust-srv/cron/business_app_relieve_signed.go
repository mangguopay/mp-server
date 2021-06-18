package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/cust-srv/dao"
	"time"
)

// 解除签约过期的商家应用
var BusinessAppRelieveSignedTask = &BusinessAppRelieveSigned{CronBase{LogCat: "解除过期的产品签约:", LockExpire: time.Minute * 5}}

type BusinessAppRelieveSigned struct {
	CronBase
}

// 运行定时任务
func (s *BusinessAppRelieveSigned) Run() {
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

func (s *BusinessAppRelieveSigned) Doing() {
	//日志
	ss_log.Info(s.LogCat + "解除过期的APP签约定时任务---------------------------start\n")

	//即将过期的签约
	expireTime := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	//修改过期的签约状态为已过期 4
	idList, err := dao.BusinessSignedDaoInst.UpdateStatusBySignedNo(expireTime, constants.SignedStatusInvalid)
	if err != nil {
		ss_log.Error(s.LogCat+"修改过期签约状态失败, err=%v", err)
	}
	ss_log.Info("过期的签约：%v", strext.ToJson(idList))

	ss_log.Info(s.LogCat + "解除过期的APP签约定时任务---------------------------end\n\n")

}
