package cron

import (
	"a.a/cu/ss_time"
	"a.a/mp-server/business-bill-srv/common"
	"a.a/mp-server/business-bill-srv/handler"
	"a.a/mp-server/common/global"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	//"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
)

// 商家结算定时任务
var BusinessSettleTask = &BusinessSettle{CronBase{LogCat: "商家结算定时任务:", LockExpire: time.Hour * 2}}

type BusinessSettle struct {
	CronBase
}

type CheckUnSettledBalanceReq struct {
	SettledId     string
	BusinessNo    string
	BusinessAccNo string
	CurrencyType  string
}

// 运行定时任务
func (s *BusinessSettle) Run() {
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

func (s *BusinessSettle) Doing() {
	ss_log.Info(s.LogCat + "结算任务---------start----------------")

	s.BusinessAppTransSettle()

	ss_log.Info(s.LogCat + "结算任务----------end-----------------\n")
}

func (s *BusinessSettle) BusinessAppTransSettle() {
	currentDate := ss_time.Now(global.Tz).Format(ss_time.DateFormat)
	startTime, _ := time.Parse(ss_time.DateTimeDashFormat, currentDate+common.StartTimePostfix)
	endTime, _ := time.Parse(ss_time.DateTimeDashFormat, currentDate+common.EndTimePostfix)
	ss_log.Info(s.LogCat+"\n结算区间：[%v——%v]", startTime.Format(ss_time.DateTimeDashFormat), endTime.Format(ss_time.DateTimeDashFormat))

	//结算数据(币种分组)
	dataList, err := dao.BusinessBillDaoInst.GetSettleData(startTime.Unix(), endTime.Unix())
	if err != nil {
		ss_log.Error(s.LogCat+"查询算数据失败, startTime:%v, endTime:%v err:%v", startTime, endTime, err)
		return
	}

	if dataList == nil {
		ss_log.Info(s.LogCat + "没有需要结算的订单")
		return
	}

	for _, data := range dataList {
		//订单数量为0,不需要结算
		if data.TotalOrder <= 0 {
			continue
		}
		s.SettleOne(startTime, endTime, data)
	}
}

func (s *BusinessSettle) SettleOne(startTime, endTime time.Time, data *dao.BusinessBillSettleData) {
	//获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		ss_log.Error(s.LogCat + "获取数据库连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//开启事务
	tx, err := dbHandler.Begin()
	if err != nil {
		ss_log.Error(s.LogCat+"开启事务失败, err:%v", err)
		return
	}

	//记录结算批次
	settleId := strext.GetDailyId()
	logData := &dao.BusinessBillSettleDao{
		SettleId:        settleId,
		AppId:           data.AppId,
		BusinessNo:      data.BusinessNo,
		StartTime:       startTime.Format(ss_time.DateTimeDashFormat),
		EndTime:         endTime.Format(ss_time.DateTimeDashFormat),
		TotalAmount:     data.TotalAmount,
		TotalRealAmount: data.TotalRealAmount,
		TotalFees:       data.TotalFees,
		TotalOrder:      data.TotalOrder,
		CurrencyType:    data.CurrencyType,
	}
	err = dao.BusinessBillSettDaoInst.InsertTx(tx, logData)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error(s.LogCat+"插入结算批次失败, data:%v, err:%v", strext.ToJson(logData), err)
		return
	}

	//开始结算
	req := new(handler.DisPoseAmountReq)
	req.SettleId = settleId
	req.BusinessNo = data.BusinessNo
	req.BusinessAccNo = data.BusinessAccNo
	req.CurrencyType = data.CurrencyType
	req.TotalAmount = data.TotalAmount
	req.TotalRealAmount = data.TotalRealAmount
	req.TotalFees = data.TotalFees
	settledErr := handler.DisPoseAmount(tx, s.LogCat, req)
	if settledErr != nil {
		ss_sql.Rollback(tx)
		ss_log.Error(s.LogCat+"结算失败(%v),req:%v, err:%v", data.CurrencyType, strext.ToJson(req), settledErr)
		return
	}

	//修改订单结算批次
	update := new(dao.UpdateOrderSettleId)
	update.SettleId = settleId
	update.BusinessNo = data.BusinessNo
	update.AppId = data.AppId
	update.StartTime = startTime.Unix()
	update.EndTime = endTime.Unix()
	err = dao.BusinessBillDaoInst.UpdateSettleIdTx(tx, update)
	if err != nil {
		ss_sql.Rollback(tx)
		ss_log.Error(s.LogCat+"修改订单settle_id失败,where:%v,err:%v", strext.ToJson(update), err)
		return
	}

	ss_log.Info(s.LogCat+"商户号:%v|币种:%v|结算批次:%v, \n结算数据:%v", data.BusinessNo, data.CurrencyType,
		settleId, strext.ToJson(data),
	)

	ss_sql.Commit(tx)
}

//检查结算后商户未结算余额是否和交易订单的总交易金额相等
func (s *BusinessSettle) checkUnSettledBalance(data *CheckUnSettledBalanceReq) {
	//查询商户未结算虚账
	vAccType := handler.GetVAccType(data.CurrencyType)
	businessVAccNo, err := dao.VaccountDaoInst.GetVaccountNo(data.BusinessAccNo, vAccType.BusinessUnSettledVaType)
	if err != nil {
		ss_log.Error(s.LogCat+"查询商户虚账失败, businessNo:%v, currencyType:%v, err:%v", data.BusinessNo, data.CurrencyType, err)
		return
	}
	//商户未结算余额
	balance, err := dao.VaccountDaoInst.GetBalanceByVAccNo(businessVAccNo)
	if err != nil {
		ss_log.Error(s.LogCat+"查询商户虚账余额失败, vAccountNo:%v, err:%v", businessVAccNo, err)
		return
	}

	//商户未结算交易总金额
	transData, err := dao.BusinessBillDaoInst.GetBusinessTransData(false, data.BusinessNo, data.CurrencyType)
	if err != nil {
		ss_log.Error("查询商户交易未结算总数据失败，businessNo=%v, currencyType=%v, err=%v", data.BusinessNo, data.CurrencyType, err)
		return
	}

	if balance != transData.TotalAmount {
		ss_log.Info(s.LogCat+"商户未结算虚账账金额有误, businessNo:%v, currencyType:%v, transAmount:%v, balance:%v",
			data.BusinessNo, data.CurrencyType, transData.TotalAmount, balance,
		)
		data := dao.BusinessChecking{
			CheckingId:         strext.GetDailyId(),
			BusinessNo:         data.BusinessNo,
			BusinessAccountNo:  data.BusinessAccNo,
			CurrencyType:       data.CurrencyType,
			BusinessBillAmount: transData.TotalAmount,
			AccountBalance:     balance,
			SettledId:          data.SettledId,
		}
		err := dao.BusinessCheckingDaoInst.Insert(data)
		if err != nil {
			if err.Error() == dao.InsertBusinessCheckingLogErr.Error() {
				ss_log.Error(s.LogCat+"重复数据，未结算对账日志:%v", strext.ToJson(data))
				return
			}
			ss_log.Error(s.LogCat+"插入未结算对账日志失败, err:%v, data:%v", err, strext.ToJson(data))
			return
		}
	}
}
