package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type ServicerTerminalDao struct{}

var ServicerTerminalDaoInstance ServicerTerminalDao

// todo 此方法意义不明
func (ServicerTerminalDao) QueryNumberFromServiceNo(servicerNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var number sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select terminal_number from servicer_terminal where servicer_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&number}, servicerNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return number.String
}
