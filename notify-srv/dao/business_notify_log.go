package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type BusinessNotifyLogDao struct {
}

var BusinessNotifyLogDaoInst BusinessNotifyLogDao

type BusinessNotifyLog struct {
	LogId      string
	OrderNo    string
	OutOrderNo string
	OrderType  string
	Status     int
	NotifyTime string
	Result     string
}

//插入日志
func (BusinessNotifyLogDao) InsertLog(req *BusinessNotifyLog) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	logId := strext.NewUUIDNoSplit()
	sqlStr := "INSERT INTO business_notify_log(id, order_no, out_order_no, status, result, order_type, create_time) " +
		"VALUES($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP) "

	err := ss_sql.Exec(dbHandler, sqlStr, logId, req.OrderNo, req.OutOrderNo, req.Status, req.Result, req.OrderType)
	if err != nil {
		return "", err
	}
	return logId, err
}

//修改通知日志结果
func (BusinessNotifyLogDao) UpdateNotifyResultById(req *BusinessNotifyLog) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "UPDATE business_notify_log SET result=$1,status=$2,notify_time=CURRENT_TIMESTAMP WHERE id=$3 "
	err := ss_sql.Exec(dbHandler, sqlStr, req.Result, req.Status, req.LogId)
	return err
}
