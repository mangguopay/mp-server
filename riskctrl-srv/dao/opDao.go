package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type opDao struct{}

var OpDaoInstance opDao

func (*opDao) GetOpFromNo(opNo string) (opName string, scriptName string, param string, score int32) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var opNameT, scriptNameT, paramT, scoreT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT op_name,script_name,param,score FROM op WHERE op_no = $1 and is_delete='0'",
		[]*sql.NullString{&opNameT, &scriptNameT, &paramT, &scoreT}, opNo); err != nil {
		return "", "", "", 0
	}
	return opNameT.String, scriptNameT.String, paramT.String, strext.ToInt32(scoreT.String)
}

func (*opDao) GetOpNoFromScriptName(scriptName string) string {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var opNoT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT op_no  FROM op WHERE script_name = $1 and is_delete='0'", []*sql.NullString{&opNoT}, scriptName); err != nil {
		return ""
	}
	return opNoT.String
}
