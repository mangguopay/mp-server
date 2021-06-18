package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/cust-srv/dao"
	"time"
)

// 生成服务商对账总计信息
var ServicerCountTask = &ServicerCount{CronBase{LogCat: "服务商对账总计定时任务:", LockExpire: time.Hour * 2}}

type ServicerCount struct {
	CronBase
}

// 运行定时任务
func (s *ServicerCount) Run() {
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

func (s *ServicerCount) Doing() {
	date := ss_time.GetDayBefore() // 获取昨天的日期
	s.HandleByDate(date)
}

// 按天处理
func (s *ServicerCount) HandleByDate(date string) {
	pageSize := 10 // 每次处理多少个记录
	index := 1

	ss_log.Info(s.LogCat+"生成服务商的对账总计信息,日期:%s, 每次处理数量:%d", date, pageSize)

	for {
		// 只到执行结束才跳出循环
		if isEnd := s.pageHandle(date, index, pageSize); isEnd {
			break
		}
		index++
	}

	ss_log.Info(s.LogCat+"生成服务商的对账总计信息-执行结束-执行了%d页", index)
}

// 分页进行处理数据
func (s *ServicerCount) pageHandle(date string, index, pageSize int) bool {
	ss_log.Info(s.LogCat+"index:%d", index)

	// 将一页的数据一次性全部取出(1个服务商现在对应只有两条记录(usd和khr))
	statisList, err := dao.ServicerCheckListDaoInst.GetCheckListStatis(date, pageSize)
	if err != nil {
		ss_log.Error(s.LogCat+"查询1页中的所有数据失败,err:%v", err)
		return true
	}

	if len(statisList) == 0 { // 没有数据了
		return true
	}

	//var wg sync.WaitGroup

	// 每页并发处理
	for _, v := range statisList {
		//wg.Add(1)
		//go func(data dao.ServicerCheckListStatis) {
		//defer wg.Done()
		if err := dao.ServicerCountDaoInst.UpdateCountData(v); err != nil {
			ss_log.Error(s.LogCat+"更新数据失败,err:%v, value:%+v", err, v)
		}
		//}(v)
	}

	//wg.Wait()

	return false
}
