package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type EventDao struct{}

var EventDaoInstance EventDao

func (*EventDao) GetEventNoFromName(eventName string) string {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var eventNoT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT event_no FROM event WHERE event_name = $1 and is_delete='0'", []*sql.NullString{&eventNoT}, eventName); err != nil {
		return ""
	}
	return eventNoT.String
}
