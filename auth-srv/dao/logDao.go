package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var (
	LogDaoInstance LogDao
)

type LogDao struct {
}

//此接口记APP、POS的操作日志
func (*LogDao) InsertAccountLog(description, accountUID, accountType string, typeInt int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, "insert into log_account(log_no,description,account_uid,type,account_type,log_time) values ($1,$2,$3,$4,$5,current_timestamp)",
		strext.GetDailyId(), description, accountUID, typeInt, accountType)
}

//此表只记WEB,即后台管理系统的关键操作日志
func (*LogDao) InsertWebAccountLog(description, accountUID string, typeStr string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, "insert into log_account_web(log_no,description,account_uid,type,create_time) values ($1,$2,$3,$4,current_timestamp)",
		strext.GetDailyId(), description, accountUID, typeStr)
}

//此表为注册成功和失败的结果表log_app_register
func (*LogDao) InsertLogAppRegister(description, uuid, opType string) error { //类型（1-注册成功，2注册失败）
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, "insert into log_app_register(description, uuid, op_type, create_time) values ($1,$2,$3,current_timestamp)",
		description, uuid, opType)
}
