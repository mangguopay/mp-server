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

func (VaccountDao) InitVaccountNo(accountNo string, vaType int32) (vaccountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	vaccountNo = strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance) values ($1,$2,$3,$4,current_timestamp,$5,$6,$7)`,
		vaccountNo, accountNo, vaType, "0", "0", "1", "0")
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo
}

func (VaccountDao) GetVaccountNoTx(tx *sql.Tx, accountNo string, vaType int32) string {
	var vaccountNo sql.NullString
	err := ss_sql.QueryRowTx(tx, `select vaccount_no from vaccount where account_no=$1 and va_type=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&vaccountNo}, accountNo, vaType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo.String
}

func (VaccountDao) InitVaccountNoTx(tx *sql.Tx, accountNo string, vaType int32) (vaccountNo string) {
	vaccountNo = strext.NewUUID()
	var balanceType string
	switch vaType {
	case constants.VaType_USD_DEBIT:
		fallthrough
	case constants.VaType_FREEZE_USD_DEBIT:
		fallthrough
	case constants.VaType_QUOTA_USD:
		fallthrough
	case constants.VaType_QUOTA_USD_REAL:
		fallthrough
	//case constants.VaType_USD_DEBIT_SRV:
	//	fallthrough
	case constants.VaType_USD_FEES:
		balanceType = "usd"
	case constants.VaType_KHR_DEBIT:
		fallthrough
	case constants.VaType_FREEZE_KHR_DEBIT:
		fallthrough
	case constants.VaType_QUOTA_KHR:
		fallthrough
	case constants.VaType_QUOTA_KHR_REAL:
		fallthrough
	//case constants.VaType_KHR_DEBIT_SRV:
	//	fallthrough
	case constants.VaType_KHR_FEES:
		balanceType = "khr"
	default:
	}

	err := ss_sql.ExecTx(tx, `insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance,balance_type) `+
		` values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8)`,
		vaccountNo, accountNo, vaType, "0", "0", "1", "0", balanceType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo
}

func (VaccountDao) GetBalance(vaccountNo string) (balance, frozenBalance string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var balanceT, frozenBalanceT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select balance,frozen_balance from vaccount where vaccount_no=$1 limit 1`,
		[]*sql.NullString{&balanceT, &frozenBalanceT}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "0", "0"
	}
	return balanceT.String, frozenBalanceT.String
}

func (VaccountDao) GetBalanceTx(tx *sql.Tx, vaccountNo string) (balance, frozenBalance string) {
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
//func (VaccountDao) SameAccFromAToBUpperZero(accountNo, vaccountNoFrom, vaccountNoTo, amount, logNo string) (errCode string) {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//	tx, _ := dbHandler.BeginTx(ctx, nil)
//	defer ss_sql.Rollback(tx)
//
//	// 出
//	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and account_no=$3`, amount, vaccountNoFrom, accountNo)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, amount, logNo)
//	if errCode != ss_err.ERR_SUCCESS {
//		ss_log.Error("err=%v", err)
//		return errCode
//	}
//
//	// 进
//	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and account_no=$3`, amount, vaccountNoTo, accountNo)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, amount, logNo)
//	if errCode != ss_err.ERR_SUCCESS {
//		ss_log.Error("err=%v", err)
//		return errCode
//	}
//
//	var tmp sql.NullString
//	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and account_no=$2 and is_delete='0' limit 1`,
//		[]*sql.NullString{&tmp}, vaccountNoFrom, accountNo)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//	if strext.ToInt64(tmp.String) < 0 {
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	tmp.String = "-1"
//	tmp.Valid = false
//	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and account_no=$2 and is_delete='0' limit 1`,
//		[]*sql.NullString{&tmp}, vaccountNoTo, accountNo)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//	if strext.ToInt64(tmp.String) < 0 {
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	ss_sql.Commit(tx)
//	return ss_err.ERR_SUCCESS
//}

// 虚拟账号进出，余额必须正
//func (VaccountDao) AccFromAToBUpperZero(vaccountNoFrom, vaccountNoTo, amount, logNo string) (errCode string) {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//	tx, _ := dbHandler.BeginTx(ctx, nil)
//	defer ss_sql.Rollback(tx)
//
//	// 出
//	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Minus, amount, logNo)
//	if errCode != ss_err.ERR_SUCCESS {
//		ss_log.Error("err=%v", err)
//		return errCode
//	}
//
//	// 进
//	err = ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoTo)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoTo, constants.VaOpType_Add, amount, logNo)
//	if errCode != ss_err.ERR_SUCCESS {
//		ss_log.Error("err=%v", err)
//		return errCode
//	}
//
//	var tmp sql.NullString
//	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
//		[]*sql.NullString{&tmp}, vaccountNoFrom)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//	if strext.ToInt64(tmp.String) < 0 {
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	tmp.String = "-1"
//	tmp.Valid = false
//	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
//		[]*sql.NullString{&tmp}, vaccountNoTo)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//	if strext.ToInt64(tmp.String) < 0 {
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	ss_sql.Commit(tx)
//	return ss_err.ERR_SUCCESS
//}

func (VaccountDao) SyncAccRemain(accNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update account set usd_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_USD_DEBIT, constants.VaType_FREEZE_USD_DEBIT)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	err = ss_sql.Exec(dbHandler, `update account set khr_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_KHR_DEBIT, constants.VaType_FREEZE_KHR_DEBIT)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 修改虚拟账户余额，余额必须正
func (VaccountDao) ModifyVaccRemainUpperZero(tx *sql.Tx, vaccountNo, amount, op, logNo string) (errCode string) {
	switch op {
	case "+":
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	default:
		return ss_err.ERR_PAY_VACCOUNT_OP_MISSING
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNo, constants.VaOpType_Add, amount, logNo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", errCode)
		return errCode
	}

	var tmp sql.NullString
	err := ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

// 提现修改虚拟账户余额,冻结资金，余额必须正
func (VaccountDao) ModifyVaccRemainAndFrozenUpperZero(tx *sql.Tx, vaccountNoFrom, amount, logNo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1,frozen_balance=frozen_balance+$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Freeze, amount, logNo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", err)
		return errCode
	}

	//errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Freeze, amount)
	//if errCode != ss_err.ERR_SUCCESS {
	//	ss_log.Error("err=%v", err)
	//	return errCode
	//}

	var tmp sql.NullString
	err = ss_sql.QueryRowTx(tx, `select balance from vaccount where vaccount_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&tmp}, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	if strext.ToInt64(tmp.String) < 0 {
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

// 提现修改虚拟账,解冻，余额必须正
func (VaccountDao) ModifyVaccFrozenUpperZero(tx *sql.Tx, vaccountNoFrom, amount, logNo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update vaccount set frozen_balance=frozen_balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, constants.VaOpType_Defreeze, amount, logNo)
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

// 修改虚拟账户余额，余额不必须为0
func (VaccountDao) ModifyVacc(tx *sql.Tx, vaccountNo, amount, op, logNo string) (errCode string) {
	switch op {
	case constants.VaOpType_Add:
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1, modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case constants.VaOpType_Minus:
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1, modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case constants.VaOpType_Freeze:
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1,frozen_balance=frozen_balance+$1, modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case constants.VaOpType_Defreeze_But_Minus:
		err := ss_sql.ExecTx(tx, `update vaccount set frozen_balance=frozen_balance+$1, modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case constants.VaOpType_Defreeze:
		err := ss_sql.ExecTx(tx, `update vaccount set frozen_balance=frozen_balance-$1, modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case constants.VaOpType_Defreeze_Minus:
		err := ss_sql.ExecTx(tx, `update vaccount set frozen_balance=frozen_balance-$1, modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case constants.VaOpType_Defreeze_Add:
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1,frozen_balance=frozen_balance-$1, modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	default:
		return ss_err.ERR_PAY_VACCOUNT_OP_MISSING
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNo, op, amount, logNo)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("err=%v", errCode)
		return errCode
	}

	return ss_err.ERR_SUCCESS
}
