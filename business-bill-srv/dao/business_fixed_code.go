package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessFixedCodeDao struct {
	FixedCode  string
	BusinessNo string
}

var BusinessFixedCodeDaoInst BusinessFixedCodeDao

func (BusinessFixedCodeDao) GetBusinessByCode(fixedCode string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT business_no FROM business_fixed_code WHERE fixed_code = $1 "

	var businessNo sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&businessNo}, fixedCode); err != nil {
		return "", err
	}
	return businessNo.String, nil
}
