package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type LogLoginDao struct{}

var LogLoginDaoInstance LogLoginDao

func (LogLoginDao) InsertLogLogin(accNo, ip, result, lag, lng string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, "insert into log_login(log_time,acc_no,ip,result,client,log_no,lat,lng) values (current_timestamp,$1,$2,$3,$4,$5,$6,$7)",
		accNo, ip, result, "", strext.GetDailyId(), lag, lng)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return ss_err.ERR_SUCCESS
}
