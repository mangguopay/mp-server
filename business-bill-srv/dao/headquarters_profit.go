package dao

import (
	"database/sql"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

type HeadquartersProfit struct {
	LogNo           string
	GeneralLedgerNo string
	Amount          string
	OrderStatus     string
	BalanceType     string
	ProfitSource    string
	OpType          string
}

var HeadquartersProfitDao HeadquartersProfit

// 插入总部利润
func (*HeadquartersProfit) InsertHeadquartersProfit(tx *sql.Tx, d *HeadquartersProfit) (string, error) {
	logNo := strext.GetDailyId()
	sqlStr := "insert into headquarters_profit " +
		"(log_no, general_ledger_no, amount, order_status, balance_type, profit_source, op_type, create_time, finish_time) " +
		"values($1,$2,$3,$4,$5,$6,$7,current_timestamp,current_timestamp) "
	err := ss_sql.ExecTx(tx, sqlStr,
		logNo, d.GeneralLedgerNo, d.Amount, d.OrderStatus, d.BalanceType, d.ProfitSource, d.OpType)
	if nil != err {
		return "", err
	}

	return logNo, nil
}

//同步余额
func (*HeadquartersProfit) SyncHeadquartersProfit(tx *sql.Tx, headVaccNo, amount, balanceType string) error {
	sqlStr := "update headquarters_profit_cashable set revenue_money = revenue_money+$2, cashable_balance = " +
		"(select sum(balance) from vaccount where vaccount_no=$1 and is_delete='0' ) where currency_type=$3 "
	if err := ss_sql.ExecTx(tx, sqlStr, headVaccNo, amount, balanceType); nil != err {
		return err
	}
	return nil
}
