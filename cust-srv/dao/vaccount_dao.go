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

func (VaccountDao) GetFrozenBalanceFromAccNo(accNo string, vaType int) (errCode, frozenBalance string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var frozenBalanceT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select frozen_balance from vaccount where account_no=$1 and va_type = $2 and is_delete='0' limit 1`,
		[]*sql.NullString{&frozenBalanceT}, accNo, vaType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM, "0"
	}
	if frozenBalanceT.String == "" {
		frozenBalanceT.String = "0"
	}
	return ss_err.ERR_SUCCESS, frozenBalanceT.String
}

func (VaccountDao) GetBalanceFromAccNo(accNo string, vaType int) (errCode, balance string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var balanceT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select balance from vaccount where account_no=$1 and va_type = $2 and is_delete='0' limit 1`,
		[]*sql.NullString{&balanceT}, accNo, vaType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM, "0"
	}

	if balanceT.String == "" {
		balanceT.String = "0"
	}

	return ss_err.ERR_SUCCESS, balanceT.String
}

func (VaccountDao) GetVaccountNo(accountNo string, vaType int32) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var vaccountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select vaccount_no from vaccount where account_no=$1 and va_type=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&vaccountNo}, accountNo, vaType)
	if err != nil {
		return "", err
	}
	return vaccountNo.String, nil
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

func (VaccountDao) InitVaccountNoTx(tx *sql.Tx, accountNo, balanceType string, vaType int32) (vaccountNo string) {
	vaccountNo = strext.NewUUID()
	err := ss_sql.ExecTx(tx, `insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance,balance_type) values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8)`,
		vaccountNo, accountNo, vaType, "0", "0", "1", "0", balanceType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo
}

func (VaccountDao) ConfirmExistVaccount(accountNo, balanceType string, vaType int32) (vaccountNo string) {
	vAccountNo, _ := VaccountDaoInst.GetVaccountNo(accountNo, vaType)
	if vAccountNo == "" {
		vAccountNo = VaccountDaoInst.InitVaccountNo(accountNo, balanceType, vaType)
		return vAccountNo
	}
	return vAccountNo
}

func (VaccountDao) ModifyVaccFrozenUpperZero(tx *sql.Tx, vaccountNoFrom, amount, logNo, reason, fees string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update vaccount set frozen_balance=frozen_balance-$1 where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	if fees != "" && fees != "0" {
		errCode = LogVaccountDaoInst.InsertPosConfirmWithdrawLogTx(tx, vaccountNoFrom, constants.VaOpType_Defreeze, amount, logNo, reason, fees)
		if errCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=%v", err)
			return errCode
		}
	} else {
		errCode = LogVaccountDaoInst.InsertPosConfirmWithdrawLogTx(tx, vaccountNoFrom, constants.VaOpType_Defreeze, amount, logNo, reason, "")
		if errCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=%v", err)
			return errCode
		}
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

// 提现修改虚拟账户余额,冻结资金，余额必须正
func (r VaccountDao) ModifyVaccRemainAndFrozenUpperZero(tx *sql.Tx, op, vaccountNoFrom, amount, logNo, reason, logVaccOpType string) (errCode string) {
	switch op {
	case "+":
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance+$1,frozen_balance=frozen_balance-$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0' `, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		err := ss_sql.ExecTx(tx, `update vaccount set balance=balance-$1,frozen_balance=frozen_balance+$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, logVaccOpType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if ss_err.ERR_SUCCESS != errCode {
		//ss_log.Error("err=%v", err)
		return errCode
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

// 提现修改虚拟账户余额,解冻,余额减,冻结减
func (r VaccountDao) ModifyVaccRemainAndFrozenUpperZero1(tx *sql.Tx, op, vaccountNoFrom, amount, logNo, reason, logVaccOpType string) (errCode string) {
	switch op {
	case "+":
		sqStr := `update vaccount set balance=balance+$1, frozen_balance=frozen_balance-$1, modify_time=current_timestamp where vaccount_no=$2 and is_delete='0' `
		err := ss_sql.ExecTx(tx, sqStr, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		sqlStr := `update vaccount set balance=balance-$1, frozen_balance=frozen_balance+$1, modify_time=current_timestamp where vaccount_no=$2 and is_delete='0'`
		err := ss_sql.ExecTx(tx, sqlStr, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, logVaccOpType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}

	errCode, balance := r.GetBalanceTx(tx, vaccountNoFrom)
	if ss_err.ERR_SUCCESS != errCode {
		//ss_log.Error("err=%v", err)
		return errCode
	}
	if strext.ToInt64(balance) < 0 {
		return ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}

//  冻结资金，余额必须正
func (VaccountDao) ModifyVaccFrozenUpperZero1(tx *sql.Tx, op, vaccountNoFrom, amount, logNo, reason, logVaccOpType string) (errCode string) {
	switch op {
	case "+":
		err := ss_sql.ExecTx(tx, `update vaccount set  frozen_balance=frozen_balance+$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0' `, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		err := ss_sql.ExecTx(tx, `update vaccount set  frozen_balance=frozen_balance-$1,modify_time = current_timestamp where vaccount_no=$2 and is_delete='0'`, amount, vaccountNoFrom)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vaccountNoFrom, logVaccOpType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}

	return ss_err.ERR_SUCCESS
}

//修改余额
func (VaccountDao) ModifyVAccBalance(tx *sql.Tx, op, vAccNo, amount, logNo, reason, logVaccOpType string) (errCode string) {
	switch op {
	case "+":
		sqlStr := `update vaccount set balance=balance+$1, modify_time=current_timestamp where vaccount_no=$2 and is_delete='0' `
		err := ss_sql.ExecTx(tx, sqlStr, amount, vAccNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		sqlStr := `update vaccount set balance=balance-$1, modify_time=current_timestamp where vaccount_no=$2 and is_delete='0'`
		err := ss_sql.ExecTx(tx, sqlStr, amount, vAccNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	}

	errCode = LogVaccountDaoInst.InsertLogTx(tx, vAccNo, logVaccOpType, amount, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		return errCode
	}

	return ss_err.ERR_SUCCESS
}

//获取当前个人冻结总金额
func (VaccountDao) GetUserFrozenBalanceCount() (usdFrozenBalance, khrFrozenBalance string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var usdFrozenBalanceT, khrFrozenBalanceT sql.NullString
	sqlStr := "select sum(frozen_balance) from vaccount where va_type = $1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&usdFrozenBalanceT}, constants.VaType_USD_DEBIT)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	sqlStr2 := "select sum(frozen_balance) from vaccount where va_type = $1 "
	err2 := ss_sql.QueryRow(dbHandler, sqlStr2, []*sql.NullString{&khrFrozenBalanceT}, constants.VaType_KHR_DEBIT)
	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
	}

	if usdFrozenBalanceT.String == "" {
		usdFrozenBalanceT.String = "0"
	}
	if khrFrozenBalanceT.String == "" {
		khrFrozenBalanceT.String = "0"
	}

	return usdFrozenBalanceT.String, khrFrozenBalanceT.String

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

type VAccountBalance struct {
	Balance      string
	CurrencyType string
	VAccType     int
}

func (VaccountDao) GetAllVAccBalanceByAccNo(accNo string) ([]*VAccountBalance, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select balance, balance_type, va_type from vaccount where account_no=$1 AND is_delete='0' `
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, accNo)
	if nil != err {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var list []*VAccountBalance
	for rows.Next() {
		var balance, balanceType, vAType sql.NullString
		if err := rows.Scan(&balance, &balanceType, &vAType); err != nil {
			return nil, err
		}
		if balance.String == "" {
			balance.String = "0"
		}
		data := &VAccountBalance{
			Balance:      balance.String,
			CurrencyType: balanceType.String,
			VAccType:     strext.ToInt(vAType.String),
		}
		list = append(list, data)
	}

	return list, nil
}
