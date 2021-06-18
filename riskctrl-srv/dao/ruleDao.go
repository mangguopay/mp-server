package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ruleDao struct{}

var RuleDaoInstance ruleDao

func (ruleDao) GetRuleFromNo(ruleNo string) string {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var ruleT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT rule FROM rule WHERE rule_no = $1 and is_delete='0'", []*sql.NullString{&ruleT}, ruleNo); err != nil {
		return ""
	}
	return ruleT.String
}

func (ruleDao) GetRuleNoFromApiType(apiType string) string {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var ruleNoT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT r.rule_no FROM rela_api_event rav LEFT JOIN rela_event_rule rer ON rav.event_no = rer.event_no"+
		" LEFT JOIN rule r ON rer.rule_no = r.rule_no WHERE rav.api_type = $1 and r.is_delete='0'", []*sql.NullString{&ruleNoT}, apiType); err != nil {
		return ""
	}
	return ruleNoT.String
}

// 获取规则阈值
func (ruleDao) GetRuleThreshold(ruleNo string) int32 {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var threshold sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "select risk_threshold from risk_threshold where rule_no=$1 and is_delete='0' limit 1",
		[]*sql.NullString{&threshold}, ruleNo); err != nil {
		return 0
	}
	return strext.ToInt32(threshold.String)
}
