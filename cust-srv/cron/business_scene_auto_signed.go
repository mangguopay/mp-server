package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/dao"
	"database/sql"
	"strings"
	"time"
)

// 商家应用自动续签
var BusinessSceneAutoSignedTask = &BusinessSceneAutoSigned{CronBase{LogCat: "商家产品自动续签:", LockExpire: time.Minute * 5}}

type BusinessSceneAutoSigned struct {
	CronBase
}

// 运行定时任务
func (s *BusinessSceneAutoSigned) Run() {
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

func (s *BusinessSceneAutoSigned) Doing() {
	//日志
	ss_log.Info(s.LogCat + "商家产品自动续签定时任务---------------------------start\n")

	s.AutoSigned()

	ss_log.Info(s.LogCat + "商家产品自动续签定时任务---------------------------end\n\n")

}

func (s *BusinessSceneAutoSigned) AutoSigned() {
	//当前时间加上三天，即提前3天自动续约
	expireTime := ss_time.Now(global.Tz).AddDate(0, 0, 3).Format(ss_time.DateTimeDashFormat)
	signedList, err := dao.SceneSignedDaoInst.GetExpireSignedList(expireTime)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error(s.LogCat + "没有需要续约的签约")
			return
		}
		ss_log.Error(s.LogCat+"查询即将过期的签约失败，err=%v", err)
		return
	}

	ss_log.Info(s.LogCat+"需要续签的签约：%v", strext.ToJson(signedList))

	//获取服务期限
	serviceTerm := dao.GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyAppSignedTerm)
	if serviceTerm == "" {
		ss_log.Error(s.LogCat+"获取全局参数失败, paramKey=%v", constants.GlobalParamKeyAppSignedTerm)
		return
	}

	for _, signed := range signedList {
		timeStr := ss_time.PostgresTimeToTime(signed.EndTime, global.Tz)
		startTime, err := time.Parse(ss_time.DateTimeDashFormat, timeStr)
		if err != nil {
			ss_log.Error(s.LogCat+"时间解析失败，timeStr=%v, err=%v", signed.EndTime, err)
			continue
		}
		//新加一条签约记录
		d := new(dao.SceneSignedDao)
		d.StartTime = startTime.Format(ss_time.DateTimeDashFormat)
		d.EndTime = startTime.AddDate(0, 0, strext.ToInt(serviceTerm)).Format(ss_time.DateTimeDashFormat)
		d.BusinessNo = signed.BusinessNo
		d.BusinessAccNo = signed.BusinessAccNo
		d.SceneNo = signed.SceneNo
		d.Rate = signed.Rate
		d.Cycle = signed.Cycle
		d.IndustryNo = signed.IndustryNo
		d.LastSignedNo = signed.SignedNo
		signedNo, addErr := dao.SceneSignedDaoInst.AutoSigned(d)
		if addErr != nil {
			if strings.Contains(addErr.Error(), ss_sql.DbDuplicateKey) {
				ss_log.Error(s.LogCat+"重复签约,err:%v, signed:%v", err, strext.ToJson(d))
				continue
			}
			ss_log.Error(s.LogCat+"增加新的签约记录失败，signed=%v, err=%v", strext.ToJson(d), addErr)
			continue
		}
		ss_log.Info(s.LogCat+"商户(%v)自动续签成功; signedNo=%v", signed.BusinessNo, signedNo)
	}

}
