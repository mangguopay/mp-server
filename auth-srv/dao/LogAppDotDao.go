package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type LogAppDotDao struct {
}

var LogAppDotDaoInst LogAppDotDao

func (LogAppDotDao) InsertLogAppDot(opType, uuid string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "insert into log_app_dot(op_type,uuid,create_time) values($1,$2,current_timestamp)"
	err := ss_sql.Exec(dbHandler, sqlStr, opType, uuid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return err
	}
	return nil
}
