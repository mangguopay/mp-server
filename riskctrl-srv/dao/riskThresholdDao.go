package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type riskThresholdDao struct{}

var RiskThresholdDaoInstance riskThresholdDao

func (*riskThresholdDao) GetThresholdFromOpNo(opNo string) string {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	var riskThresholdT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT risk_threshold FROM risk_threshold WHERE op_no = $1 and is_delete='0'", []*sql.NullString{&riskThresholdT}, opNo); err != nil {
		return ""
	}
	return riskThresholdT.String
}
