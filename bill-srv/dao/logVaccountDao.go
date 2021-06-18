package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"context"
	"database/sql"
	"errors"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	tmProto "a.a/mp-server/common/proto/tm"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/common/ss_struct"
)

type LogVaccountDao struct {
	CreateTime string
	Amount     string
	Reason     string
	Balance    string
	BizLogNo   string
	OpType     string
}

var LogVaccountDaoInst LogVaccountDao

func (LogVaccountDao) InsertLogTxV2(tmProxy *ss_struct.TmServerProxy, vaccountNo, opType, amount, bizLogNo, reason string) error {
	balance, fbalance, err1 := VaccountDaoInst.GetBalanceV2(tmProxy, vaccountNo)
	if err1 != nil {
		return err1
	}

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
func (LogVaccountDao) InsertLogTx(tx *sql.Tx, vaccountNo, opType, amount, bizLogNo, reason string) string {

	balance, fbalance := VaccountDaoInst.GetBalance2Tx(tx, vaccountNo)

	err := ss_sql.ExecTx(tx, `insert into log_vaccount(log_no,create_time,vaccount_no,amount,op_type,frozen_balance,balance,reason,biz_log_no) values ($1,current_timestamp,$2,$3,$4,$5,$6,$7,$8)`,
		strext.GetDailyId(), vaccountNo, amount, opType, fbalance, balance, reason, bizLogNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}
func (LogVaccountDao) InsertPosConfirmWithdrawLogTx(tx *sql.Tx, vaccountNo, opType, amount, bizLogNo, reason, fees string) string {
	balance, fbalance := VaccountDaoInst.GetBalance2Tx(tx, vaccountNo)
	var resultBalance, resultFbalance string
	if fees != "" && fees != "0" { // 还没扣除手续费的时候,需要把手续费加进balance
		// 当前余额需要相加
		resultBalance = ss_count.Add(balance, fees)
		resultFbalance = fbalance
	} else {
		resultBalance = balance
		resultFbalance = fbalance
	}
	err := ss_sql.ExecTx(tx, `insert into log_vaccount(log_no,create_time,vaccount_no,amount,op_type,frozen_balance,balance,reason,biz_log_no) values ($1,current_timestamp,$2,$3,$4,$5,$6,$7,$8)`,
		strext.GetDailyId(), vaccountNo, amount, opType, resultFbalance, resultBalance, reason, bizLogNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (LogVaccountDao) GetLogVAccountByBizLogNo(account, bizLogNo, reason string) (*LogVaccountDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "SELECT lv.create_time, lv.amount, lv.reason, lv.balance, lv.biz_log_no, lv.op_type " +
		"FROM log_vaccount lv " +
		"LEFT JOIN vaccount v ON v.vaccount_no = lv.vaccount_no " +
		"WHERE v.account_no = $1 AND lv.biz_log_no = $2 AND lv.reason = $3 LIMIT 1"

	var createTime, amount, reasonT, balance, bizLogNoT, opType sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&createTime, &amount, &reasonT, &balance, &bizLogNoT, &opType},
		account, bizLogNo, reason)
	if err != nil {
		return nil, err
	}

	obj := new(LogVaccountDao)
	obj.CreateTime = createTime.String
	obj.Amount = amount.String
	obj.Reason = reasonT.String
	obj.Balance = balance.String
	obj.BizLogNo = bizLogNoT.String
	obj.OpType = opType.String

	return obj, nil
}

type WriteOffCodeVAccLog struct {
	CreateTime   string
	Amount       string
	OpType       string
	OrderNo      string
	OrderType    string
	CurrencyType string
}

func (LogVaccountDao) GetLogVAccountJoinWriteOff(accountNo, LogNo, reason string) (*WriteOffCodeVAccLog, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "SELECT lv.create_time, lv.amount, lv.op_type, wo.income_order_no, wo.transfer_order_no, v.balance_type " +
		"FROM log_vaccount lv " +
		"LEFT JOIN vaccount v ON v.vaccount_no = lv.vaccount_no " +
		"LEFT JOIN writeoff wo ON wo.code = lv.biz_log_no " +
		"WHERE v.account_no = $1 AND lv.log_no = $2 AND lv.reason = $3 LIMIT 1"

	var createTime, amount, opType, incomeOrderNo, transferOrderNo, currencyType sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&createTime, &amount, &opType, &incomeOrderNo,
		&transferOrderNo, &currencyType}, accountNo, LogNo, reason)
	if err != nil {
		return nil, err
	}

	obj := new(WriteOffCodeVAccLog)
	obj.CreateTime = createTime.String
	obj.Amount = amount.String
	obj.OpType = opType.String
	obj.OrderNo = incomeOrderNo.String
	obj.CurrencyType = currencyType.String
	obj.OrderType = "充值"
	if obj.OrderNo == "" {
		obj.OrderNo = transferOrderNo.String
		obj.OrderType = "转账"
	}

	return obj, nil
}
