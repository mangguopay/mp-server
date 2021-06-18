package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type PushRecordDao struct{}

var PushRecordDaoInst PushRecordDao

func (PushRecordDao) Insert(business, phone, content, pushNo, tempNo, message string, status int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()
	err := ss_sql.Exec(dbHandler, `insert into push_record(id,business,phone,content,status,push_no,temp_no,message,create_time )
			values($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp )`,
		id, business, phone, content, status, pushNo, tempNo, message)
	if err != nil {
		return err
	}

	return nil
}
