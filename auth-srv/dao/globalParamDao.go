package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type GlobalParamDao struct{}

var GlobalParamDaoInstance GlobalParamDao

func (*GlobalParamDao) QeuryParamValue(paramKey string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var paramValueT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select param_value from global_param where param_key=$1   limit 1`, []*sql.NullString{&paramValueT}, paramKey)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return paramValueT.String
}
