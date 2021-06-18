package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type LogAppMessagesDao struct {
	LogNo       string
	OrderNo     string
	AppMessType string
	OrderType   string
	AccountNo   string
	OrderStatus string
}

var LogAppMessagesDaoInst LogAppMessagesDao

func (*LogAppMessagesDao) AddLogAppMessages(d *LogAppMessagesDao) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "insert into log_app_messages(log_no, order_no, order_type, account_no, app_mess_type, order_status, create_time) " +
		"values($1,$2,$3,$4,$5,$6,current_timestamp)"
	return ss_sql.Exec(dbHandler, sqlStr, strext.GetDailyId(), d.OrderNo, d.OrderType, d.AccountNo, d.AppMessType, d.OrderStatus)
}
