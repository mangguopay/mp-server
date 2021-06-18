package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	UploadFileLogDaoInstance UploadFileLogDao
)

type UploadFileLogDao struct {
	Id         string
	CreateTime string
	AccountNo  string
	FileName   string
	FileType   string
}

func (*UploadFileLogDao) GetUploadFileInfo(fileId string) (*UploadFileLogDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := ` select id, account_no, file_name, file_type from upload_file_log where id = $1 `
	row, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlStr, fileId)
	if stmt != nil {
		defer stmt.Close()
	}
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}

	var id, accountNo, fileName, fileType sql.NullString
	errT = row.Scan(&id, &accountNo, &fileName, &fileType)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}

	return &UploadFileLogDao{
		Id:        id.String,
		AccountNo: accountNo.String,
		FileName:  fileName.String,
		FileType:  fileType.String,
	}, nil
}
