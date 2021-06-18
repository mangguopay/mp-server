package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type relaApiEventDao struct{}

var RelaApiEventDaoInstance relaApiEventDao

func (*relaApiEventDao) GetEventNoFromApi(apiType string) string {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var eventNoT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT event_no FROM rela_api_event WHERE api_type = $1 and is_delete='0'", []*sql.NullString{&eventNoT}, apiType); err != nil {
		return ""
	}
	return eventNoT.String
}
