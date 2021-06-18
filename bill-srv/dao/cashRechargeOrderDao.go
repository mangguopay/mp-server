package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

type CashRechargeOrderDao struct {
	LogNo        string
	AccNo        string
	Amount       string
	IdenNo       string
	CreateTime   string
	OrderStatus  string
	CurrencyType string
	OpAccNo      string
	PaymentType  string
	Notes        string
}

var CashRechargeOrderDaoInst CashRechargeOrderDao

func (CashRechargeOrderDao) InsertCashRecharge(tx *sql.Tx, data CashRechargeOrderDao) (logNo string) {
	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into servicer_cash_recharge_order(log_no, acc_no, amount, iden_no, create_time, order_status, currency_type, op_acc_no, payment_type, notes)
		values($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9)`,
		logNo, data.AccNo, data.Amount, data.IdenNo, data.OrderStatus, data.CurrencyType,
		data.OpAccNo, data.PaymentType, data.Notes)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return logNo
}
