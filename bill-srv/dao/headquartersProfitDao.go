package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
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

var HeadquartersProfitDaoInstance HeadquartersProfit

// 插入总部利润
func (*HeadquartersProfit) InsertHeadquartersProfit(tx *sql.Tx, d *HeadquartersProfit) string {
	logNo := strext.GetDailyId()
	sqlStr := "insert into headquarters_profit " +
		"(log_no, general_ledger_no, amount, order_status, balance_type, profit_source, op_type, create_time, finish_time) " +
		"values($1,$2,$3,$4,$5,$6,$7,current_timestamp,current_timestamp) "
	err := ss_sql.ExecTx(tx, sqlStr,
		logNo, d.GeneralLedgerNo, d.Amount, d.OrderStatus, d.BalanceType, d.ProfitSource, d.OpType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return ss_err.ERR_SUCCESS
}
