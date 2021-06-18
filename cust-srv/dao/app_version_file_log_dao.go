package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	AppVersionFileLogDaoInstance AppVersionFileLogDao
)

type AppVersionFileLogDao struct {
}

func (*AppVersionFileLogDao) CheckAppNo(AppVersionFileLogNo string) (err, filename string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := ` select file_name from upload_file_log where id = $1 `
	var fileName sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&fileName}, AppVersionFileLogNo)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return ss_err.ERR_PARAM, ""
	}

	return ss_err.ERR_SUCCESS, fileName.String
}

func (*AppVersionFileLogDao) ModifyFileName(tx *sql.Tx, AppVersionFileLogNo, filename string) (err string) {
	sqlStr := ` update upload_file_log set file_name=$2 where id = $1 `
	errT := ss_sql.ExecTx(tx, sqlStr, AppVersionFileLogNo, filename)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}
