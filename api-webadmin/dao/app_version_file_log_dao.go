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
	UploadFileLogDaoInstance UploadFileLogDao
)

type UploadFileLogDao struct {
}

func (*UploadFileLogDao) AddUploadFileLog(accNo, accType, filename, fileType string) (idStr, err string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()

	sqlStr := ` insert into upload_file_log(id,account_no,file_name,file_type,account_type,create_time) values($1,$2,$3,$4,$5,current_timestamp)`
	errT := ss_sql.Exec(dbHandler, sqlStr, id, accNo, filename, fileType, accType)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return "", ss_err.ERR_PARAM
	}

	return id, ss_err.ERR_SUCCESS
}
