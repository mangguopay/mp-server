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
func (*LogDao) InsertAccountLog(description, accountUID string, typeInt int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, "insert into log_account(log_no,description,account_uid,log_time,type) values ($1,$2,$3,current_timestamp,$4)",
		strext.GetDailyId(), description, accountUID, typeInt)
}

//此表只记WEB,即后台管理系统的关键操作日志
func (*LogDao) InsertWebAccountLog(description, accountUID string, typeStr string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, "insert into log_account_web(log_no,description,account_uid,type,create_time) values ($1,$2,$3,$4,current_timestamp)",
		strext.GetDailyId(), description, accountUID, typeStr)
}
