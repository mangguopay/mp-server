package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"strings"
)

type LogBusinessMessage struct {
	LogNo       string
	IsRead      string
	AccountNo   string
	CreateTime  string
	Content     string
	AccountType string
}

var LogBusinessMessageDao LogBusinessMessage

func (LogBusinessMessage) Insert(d *LogBusinessMessage) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "INSERT INTO log_business_messages(log_no, is_read, account_no, content, account_type, create_time) " +
		"VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)"
	return ss_sql.Exec(dbHandler, sqlStr, d.LogNo, d.IsRead, d.AccountNo, d.Content, d.AccountType)
}

type BusinessMessageTemp struct {
	PushNoList []string
	TitleKey   string
	ContentKey string
	LenArgs    int32
}

func (LogBusinessMessage) GetTemplate(tempNo string) (*BusinessMessageTemp, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select push_nos,title_key,content_key,len_args from push_temp where temp_no=$1 and is_delete = '0' limit 1"
	var pushNos, titleKey, contentKey, lenArgs sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&pushNos, &titleKey, &contentKey, &lenArgs}, tempNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	pushNoList := strings.Split(pushNos.String, ",")
	obj := new(BusinessMessageTemp)
	obj.PushNoList = pushNoList
	obj.TitleKey = titleKey.String
	obj.ContentKey = contentKey.String
	obj.LenArgs = strext.ToInt32(lenArgs.String)
	return obj, nil
}
