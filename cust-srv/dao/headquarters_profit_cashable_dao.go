package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type HeadquartersProfitCashableDao struct {
}

var HeadquartersProfitCashableDaoInstance HeadquartersProfitCashableDao

func (*HeadquartersProfitCashableDao) SyncAccProfit(tx *sql.Tx, headVaccNo, amount, balanceType string) string {
	err := ss_sql.ExecTx(tx, `update headquarters_profit_cashable set  modify_time=current_timestamp,revenue_money=revenue_money+$2, cashable_balance=(select sum(balance) from vaccount where vaccount_no=$1  and is_delete='0') where currency_type=$3`, headVaccNo, amount, balanceType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}

func (*HeadquartersProfitCashableDao) GetCashableBalance(balanceType string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT cashable_balance FROM headquarters_profit_cashable where currency_type = $1 "
	var cashableBalance sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cashableBalance}, balanceType)
	if err != nil {
		return "", err
	}
	return cashableBalance.String, nil
}
