package dao

import (
	"database/sql"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

var LogVaccountDaoInst LogVaccountDao

type LogVaccountDao struct {
	VaccountNo    string
	Amount        string
	OpType        string
	Balance       string
	FrozenBalance string
	Reason        string
	BizLogNo      string
}

func (l *LogVaccountDao) InsertLogTx(tx *sql.Tx, logVacc LogVaccountDao) error {
	sqlStr := `insert into log_vaccount (log_no, create_time, vaccount_no, amount, op_type, balance, frozen_balance, reason, biz_log_no) `
	sqlStr += ` values ($1, current_timestamp, $2, $3, $4, $5, $6, $7, $8)`

	err := ss_sql.ExecTx(tx, sqlStr, strext.GetDailyId(), logVacc.VaccountNo, logVacc.Amount,
		logVacc.OpType, logVacc.Balance, logVacc.FrozenBalance, logVacc.Reason, logVacc.BizLogNo,
	)

	return err
}
