package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var (
	SMSDaoInstance SMSDao
)

type SMSDao struct {
}

// 记录发型短信
func (*SMSDao) InsertSMSRecord(msgID, account, business, mobile, msg string, status int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	id := strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into sms_send_record(id,msgid,account,business,mobile,msg,status,created_at)values($1,$2,$3,$4,$5,$6,$7,current_timestamp)`,
		id, msgID, account, business, mobile, msg, status)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
