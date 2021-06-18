package dao

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	tmProto "a.a/mp-server/common/proto/tm"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/common/ss_struct"
	"context"
	"database/sql"
	"errors"
)

type LogVaccountDao struct {
}

var LogVaccountDaoInst LogVaccountDao

func (LogVaccountDao) InsertLogTx(tx *sql.Tx, vaccountNo, opType, amount, bizLogNo, reason string) string {
	//balance, fbalance := VaccountDaoInst.GetBalance(vaccountNo)
	balance, fbalance := VaccountDaoInst.GetBalanceTx(tx, vaccountNo)
	if fbalance == "" {
		fbalance = "0"
	}
	err := ss_sql.ExecTx(tx, `insert into log_vaccount(log_no,create_time,vaccount_no,amount,op_type,frozen_balance,balance,reason,biz_log_no) values ($1,current_timestamp,$2,$3,$4,$5,$6,$7,$8)`,
		strext.GetDailyId(), vaccountNo, amount, opType, fbalance, balance, reason, bizLogNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (LogVaccountDao) InsertLogTxV2(tmProxy *ss_struct.TmServerProxy, vaccountNo, opType, amount, bizLogNo, reason string) error {
	balance, fbalance, err := VaccountDaoInst.GetBalanceV2(tmProxy, vaccountNo)
	if err != nil {
		return err
	}
	if fbalance == "" {
		fbalance = "0"
	}

	// step-2
	sql := `insert into log_vaccount(log_no,create_time,vaccount_no,amount,op_type,frozen_balance,balance,reason,biz_log_no) values ($1,current_timestamp,$2,$3,$4,$5,$6,$7,$8)`
	args := []string{strext.GetDailyId(), vaccountNo, amount, opType, fbalance, balance, reason, bizLogNo}

	rsp, err := tmProxy.GetTmServer().TxExec(context.TODO(), &tmProto.TxExecRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return err
	}

	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}

	return nil
}
