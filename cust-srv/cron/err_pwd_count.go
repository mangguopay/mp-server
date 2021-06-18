package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/dao"
	"strings"
	"time"
)

// 生成服务商对账总计信息
var ErrPwdCountTask = &ErrPwdCount{CronBase{LogCat: "服务商对账总计定时任务:", LockExpire: time.Hour * 2}}

type ErrPwdCount struct {
	CronBase
}

// 运行定时任务
func (s *ErrPwdCount) Run() {
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

func (s *ErrPwdCount) Doing() {
	ss_log.Info("---------------------错误清除------------------")
	// 获取keys
	result, err := cache.RedisClient.Keys("err_pwd_*").Result()
	if err != nil {
		ss_log.Error("err: %s", err.Error())
		return
	}
	ss_log.Info("errResultKeys------------------> %+v", result)
	if len(result) <= 0 {
		return
	}
	err = cache.RedisClient.Del(result...).Err()
	if err != nil {
		ss_log.Error("定时任务 Del err: %s", err.Error())
		return
	}
	// 解除禁用
	for _, v := range result {
		account := strings.Split(v, cache.PrePwdErrCountKey)[1]
		if errStr := dao.AccDaoInstance.UpdateAccountStatusByAccount(account, constants.AccountUseStatusNormal); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("密码错误解除禁用失败,账号为: %s", account)
			continue
		}
	}
}
