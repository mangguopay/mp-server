package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type ServicerprofitledgerDao struct{}

var ServicerprofitledgerDaoInstance ServicerprofitledgerDao

// 插入总部利润
func (*ServicerprofitledgerDao) InsertServicerProfitLedger(tx *sql.Tx, orderNo, amountOrder, serviceFeeAmountSum, splitProportion, actualIncome, servicerNo, currencyType, orderType string) string {
	err := ss_sql.ExecTx(tx, `insert into servicer_profit_ledger(log_no,amount_order,servicefee_amount_sum,split_proportion,actual_income,servicer_no,currency_type,order_type,payment_time)
        values($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)`,
		orderNo, amountOrder, serviceFeeAmountSum, splitProportion, actualIncome, servicerNo, currencyType, orderType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS
}
