package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ServicerTerminalDao struct {
}

var ServicerTerminalDaoInst ServicerTerminalDao

//查看pos属于哪个服务商
func (ServicerTerminalDao) GetSerPosServicerNoByPosNo(terminalNo string) (servicerNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select servicer_no from servicer_terminal where pos_sn = $1 and is_delete= $2 and use_status = $3"
	var servicerNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&servicerNoT}, terminalNo, 0, 1)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return servicerNoT.String

}
