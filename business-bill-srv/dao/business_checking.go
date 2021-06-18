package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type BusinessCheckingDao struct {
}

var BusinessCheckingDaoInst BusinessCheckingDao

type BusinessChecking struct {
	CheckingId         string
	BusinessNo         string
	BusinessAccountNo  string
	CurrencyType       string
	BusinessBillAmount int64
	AccountBalance     int64
	SettledId          string
}

func (b *BusinessCheckingDao) Insert(data BusinessChecking) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "INSERT INTO business_checking(checking_id, business_no, business_account_no,currency_type,business_bill_amount,"
	sqlStr += "account_balance, settle_id, create_time) VALUES($1,$2,$3,$4,$5,$6,$7,CURRENT_TIMESTAMP)"

	return ss_sql.Exec(dbHandler, sqlStr, data.CheckingId, data.BusinessNo, data.BusinessAccountNo, data.CurrencyType,
		data.BusinessBillAmount, data.AccountBalance, data.SettledId)
}
