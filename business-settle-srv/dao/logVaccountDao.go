package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type LogVaccountDao struct {
}

var LogVaccountDaoInst LogVaccountDao

func (LogVaccountDao) InsertLogTx(tx *sql.Tx, vaccountNo, opType, amount, bizLogNo, reason string) string {
	balance, fbalance := VaccountDaoInst.GetBalance(tx, vaccountNo)

	err := ss_sql.ExecTx(tx, `insert into log_vaccount(log_no,create_time,vaccount_no,amount,op_type,frozen_balance,balance,reason,biz_log_no) values ($1,current_timestamp,$2,$3,$4,$5,$6,$7,$8)`,
		strext.GetDailyId(), vaccountNo, amount, opType, fbalance, balance, reason, bizLogNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}
