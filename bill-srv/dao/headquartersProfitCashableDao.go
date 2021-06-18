package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type HeadquartersProfitCashableDao struct {
}

var HeadquartersProfitCashableDaoInstance HeadquartersProfitCashableDao

func (*HeadquartersProfitCashableDao) SyncAccProfit(tx *sql.Tx, headVaccNo, amount, balanceType string) string {
	err := ss_sql.ExecTx(tx, `update headquarters_profit_cashable set revenue_money=revenue_money+$2, cashable_balance=(select sum(balance) from vaccount where vaccount_no=$1  and is_delete='0') where currency_type=$3`, headVaccNo, amount, balanceType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	//err = ss_sql.ExecTx(tx, `update headquarters_profit_cashable set revenue_money=revenue_money+$2, cashable_balance=(select sum(balance) from vaccount where account_no=$1  and is_delete='0') where currency_type=$3`, headVaccNo, amount,balanceType)
	//if nil != err {
	//	ss_log.Error("err=%v", err)
	//	return ""
	//}
	return ss_err.ERR_SUCCESS
}
