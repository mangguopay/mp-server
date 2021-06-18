package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	LogAccDaoInstance LogAccDao
)

type LogAccDao struct {
}

func (*LogAccDao) InsertAccountLog(description, accountUID string, t int) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, "insert into log_account(log_no,description,account_uid,log_time,type) values ($1,$2,$3,current_timestamp,$4)",
		strext.GetDailyId(), description, accountUID, t)

	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS

}
