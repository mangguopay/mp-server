package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type logLoginDao struct{}

var LogLoginDaoInstance logLoginDao

func (*logLoginDao) GetCountFromTime(t string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var countT sql.NullString
	//err := ss_sql.QueryRow(dbHandler, `SELECT count(log_time) FROM log_login WHERE log_time >  current_timestamp+interval  '-1 day'`,
	sqlStr := "SELECT count(log_time) FROM log_login WHERE log_time >  current_timestamp+interval " + "'" + t + "'"
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&countT})
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return countT.String
}
