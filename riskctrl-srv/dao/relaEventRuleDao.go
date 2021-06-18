package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type relaEventRuleDao struct{}

var RelaEventRuleDaoInstance relaEventRuleDao

func (*relaEventRuleDao) GetRuleNoFromEventNo(eventNo string) string {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var ruleNoT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT rule_no FROM rela_event_rule WHERE event_no = $1 and is_delete='0'", []*sql.NullString{&ruleNoT}, eventNo); err != nil {
		return ""
	}
	return ruleNoT.String
}
