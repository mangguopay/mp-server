package cron

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/cust-srv/dao"
	"math"
	"time"
)

// 生成服务商对账列表信息
var ServicerCheckListTask = &ServicerCheckList{CronBase{LogCat: "服务商对账列表定时任务:", LockExpire: time.Hour * 2}}

type ServicerCheckList struct {
	CronBase
}

// 运行定时任务
func (s *ServicerCheckList) Run() {
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

func (s *ServicerCheckList) Doing() {
	date := ss_time.GetDayBefore() // 获取昨天的日期
	s.HandleByDate(date)
}

// 按天处理
func (s *ServicerCheckList) HandleByDate(date string) {
	ss_log.Info(s.LogCat+"开始生成服务商的对账信息,日期:%s", date)

	// 统计某一天有多少服务商有账单信息
	total, totalErr := dao.BillingDetailsResultsDaoInst.CountServicerByDate(date)
	if totalErr != nil {
		ss_log.Error(s.LogCat+"按日期统计总服务商数量失败,date:%s,err:%v", date, totalErr)
		return
	}

	if total < 1 {
		ss_log.Info(s.LogCat + "总服务商数量为0,无须处理")
		return
	}

	pageSize := 10                                                  // 每次查询数量
	totalPage := int(math.Ceil(float64(total) / float64(pageSize))) // 总页数

	ss_log.Info(s.LogCat+"总服务商数量:%d, 每页数量:%d, 分成总页数:%d", total, pageSize, totalPage)

	// 分页进行统计数据
	for p := 1; p <= totalPage; p++ {
		s.countByPage(date, p, pageSize)
	}

	ss_log.Info(s.LogCat + "执行结束")
}

// 分页进行统计数据
func (s *ServicerCheckList) countByPage(date string, page, pageSize int) {
	ss_log.Info(s.LogCat+"处理第%d页, 每页数量:%d", page, pageSize)

	// 分页获取服务商编号
	servicerNoList, err := dao.BillingDetailsResultsDaoInst.GetServicerNoByDate(date, page, pageSize)
	if err != nil {
		ss_log.Error(s.LogCat+"分页进行统计数据失败,err:%v", err)
		return
	}

	ss_log.Info(s.LogCat+"处理第%d页, len=%d, servicerNoList:%v", page, len(servicerNoList), servicerNoList)

	//var wg sync.WaitGroup

	// 每页并发处理
	for _, v := range servicerNoList {
		//wg.Add(1)
		//go func(servicerNo string) {
		//defer wg.Done()
		s.countByAccountNo(v, date)
		//}(v)
	}

	//wg.Wait()
}

// 统计每个账户的数据
func (s *ServicerCheckList) countByAccountNo(servicerNo, date string) {
	if servicerNo == "" {
		ss_log.Error(s.LogCat+"servicerNo为空,servicerNo:%s,date:%s", servicerNo, date)
		return
	}

	// 循环插入各种币种的数据
	for _, currencyType := range []string{constants.CURRENCY_USD, constants.CURRENCY_KHR} {

		// 1.查询统计数据
		retData, qErr := dao.BillingDetailsResultsDaoInst.GetServicerStatis(servicerNo, currencyType, date)
		if qErr != nil {
			ss_log.Error(s.LogCat+"统计账户数据失败,servicerNo:%s,currencyType:%s,date:%s,err:%v", servicerNo, currencyType, date, qErr)
			return
		}

		ss_log.Info(s.LogCat+"统计单个账户,date:%s,servicerNo:%s,retData:%+v", date, servicerNo, retData)

		// 2.插入统计数据
		if err := dao.ServicerCheckListDaoInst.InsertByCurrency(retData); err != nil {
			ss_log.Error(s.LogCat+"插入统计数据失败,servicerNo:%s, retData:%+v,err:%v", servicerNo, retData, err)
			return
		}
	}
}
