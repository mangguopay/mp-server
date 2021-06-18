package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type VaccountDao struct {
}

var VaccountDaoInst VaccountDao

func (VaccountDao) GetVaccountNo(accountNo string, vaType int32) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var vaccountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select vaccount_no from vaccount where account_no=$1 and va_type=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&vaccountNo}, accountNo, vaType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo.String
}
func (VaccountDao) GetVaccountNoFromMoneyType(accountNo, moneyType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var vaccountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select vaccount_no from vaccount where account_no=$1 and balance_type=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&vaccountNo}, accountNo, moneyType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo.String
}

func (VaccountDao) InitVaccountNo(accountNo, balanceType string, vaType int32) (vaccountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	vaccountNo = strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance,balance_type) values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8)`,
		vaccountNo, accountNo, vaType, "0", "0", "1", "0", balanceType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo
}

func (VaccountDao) GetBalance(tx *sql.Tx, vaccountNo string) (balance, frozenBalance string) {
	var balanceT, frozenBalanceT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select balance,frozen_balance from vaccount where vaccount_no=$1 limit 1`,
		[]*sql.NullString{&balanceT, &frozenBalanceT}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "0", "0"
	}
	return balanceT.String, frozenBalanceT.String
}

// 同名进出，余额必须正
func (VaccountDao) SameAccFromAToBUpperZero(tx *sql.Tx, toAmount, accountNo, vaccountNoFrom, vaccountNoTo, amount, logNo, reason string) (errCode string) {
	// 出
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and account_no=$3`, toAmount, vaccountNoFrom, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, toAmount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	// 进
	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and account_no=$3`, amount, vaccountNoTo, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	var tmp sql.NullString
	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and account_no=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoFrom, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	tmp.String = "-1"
	tmp.Valid = false
	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and account_no=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoTo, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

// 虚拟账号进出，余额必须正
func (r VaccountDao) AccFromAToBUpperZero(tx *sql.Tx, vaccountNoFrom, vaccountNoTo, amount, logNo, reason string) (errCode string) {

	// 出
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	// 进
	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoTo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode, balance = r.GetBalanceTx(tx, vaccountNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 存款,虚拟账号进出，余额必须正
func (r VaccountDao) SaveMoneyAccFromAToBUpperZero(tx *sql.Tx, vaccountNoFrom, vaccountNoTo, fromAmount, toAmount, logNo, reason string) (errCode string) {

	// 出
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1,modify_time=current_timestamp where vaccount_no=$2 and is_delete='0'`, fromAmount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, fromAmount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	// 进
	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1,modify_time=current_timestamp where vaccount_no=$2 and is_delete='0'`, toAmount, vaccountNoTo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, toAmount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode, balance = r.GetBalanceTx(tx, vaccountNoTo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

func (VaccountDao) SyncAccRemain(tx *sql.Tx, accNo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update account set usd_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_USD_DEBIT, constants.VaType_FREEZE_USD_DEBIT)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	err = ss_sql.ExecTx(tx, `update account set khr_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_KHR_DEBIT, constants.VaType_FREEZE_KHR_DEBIT)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 修改虚拟账户余额，余额必须正
func (r VaccountDao) ModifyVaccRemainUpperZero(tx *sql.Tx, vaccountNo, amount, op, logNo, reason string) (errCode string) {
	var opType string
	switch op {
	case "+":
		opType = constants.VaOpType_Add
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		opType = constants.VaOpType_Minus
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	default:
		return ss_err.ERR_PAY_VACCOUNT_OP_MISSING
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNo, opType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", errCode)
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=%v", errCode)
		return errCode
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 提现修改虚拟账户余额,冻结资金，余额必须正
func (r VaccountDao) ModifyVaccRemainAndFrozenUpperZero(tx *sql.Tx, vaccountNoFrom, amount, logNo, reason string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1,frozen_balance=frozen_balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Freeze, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if ss_err.ERR_SUCCESS != errCode {
		ss_log.Error("err=%v", err)
		return errCode
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

// 提现修改虚拟账,解冻，余额必须正
func (VaccountDao) ModifyVaccFrozenUpperZero(tx *sql.Tx, vaccountNoFrom, amount, logNo, reason string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update vaccount set frozen_balance=frozen_balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Defreeze, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	var tmp sql.NullString
	err = ss_sql.QueryRowTx(tx, `select frozen_balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		ss_log.Error("----->%s", "超出解冻金额")
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

func (VaccountDao) GetBalanceTx(tx *sql.Tx, vaccountNo string) (errCode, balance string) {
	var tmp sql.NullString
	err := ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH, ""
	}
	return ss_err.ERR_SUCCESS, tmp.String
}

func (va *VaccountDao) ConfirmExistVaccount(accountNo, balanceType string, vaType int32) (vaccountNo string) {
	vaccountNo = va.GetVaccountNo(accountNo, vaType)
	if vaccountNo == "" {
		vaccountNo = va.InitVaccountNo(accountNo, balanceType, vaType)
		return vaccountNo
	}
	return vaccountNo
}
