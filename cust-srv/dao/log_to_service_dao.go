package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type LogToServiceDao struct {
}

var (
	LogToServiceDaoInstance LogToServiceDao
)

func (*LogToServiceDao) InsertLogToService(serviceNo, amount, collectType, cardNo string, orderType, vaType int) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	logNoT := strext.GetDailyId()

	err := ss_sql.Exec(dbHandler, `insert into log_to_servicer(log_no,currency_type,servicer_no,collection_type,card_no,amount,order_type,order_status,create_time) 
				values ($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)`,
		logNoT, vaType, serviceNo, collectType, cardNo, amount, orderType, "0")
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_LOG_TO_SERVICE_FAILD
	}
	return ss_err.ERR_SUCCESS
}
func (LogToServiceDao) InsertLogToService1(currencyType, servicerNo, collectionType, cardNo, amount, orderType, orderStatus string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//新增转账
	sqlStr := "insert into log_to_servicer(log_no, currency_type, servicer_no, collection_type, card_no, amount, create_time, order_type, order_status) " +
		" values ($1,$2,$3,$4,$5,$6,current_timestamp,$7,$8)"
	if err := ss_sql.Exec(dbHandler, sqlStr, strext.GetDailyId(), currencyType, servicerNo, collectionType, cardNo, amount, orderType, orderStatus); err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}
	return ss_err.ERR_SUCCESS
}
func (*LogToServiceDao) QueryOrderStatusFromlogNo(orderNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var statusT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select order_status from log_to_servicer where log_no=$1  limit 1`, []*sql.NullString{&statusT}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return statusT.String
}

func (*LogToServiceDao) QueryLogFromlogNo(orderNo string) (string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var currencyTypeT, servicerNoT, amountT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select currency_type,servicer_no,amount from log_to_servicer where log_no=$1  limit 1`, []*sql.NullString{&currencyTypeT, &servicerNoT, &amountT}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", ""
	}
	return currencyTypeT.String, servicerNoT.String, amountT.String
}

func (*LogToServiceDao) UpdateStatusFromLogNo(orderNo string, status int32) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update log_to_servicer set order_status= $1 where log_no=$2`, status, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_OPERATE_FAILD
	}
	return ss_err.ERR_SUCCESS
}
