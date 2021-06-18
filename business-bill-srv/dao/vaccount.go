package dao

import (
	"a.a/cu/strext"
	"database/sql"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type VaccountDao struct {
}

var VaccountDaoInst VaccountDao

func (VaccountDao) GetVaccountNo(accountNo string, vaType int) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var vaccountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select vaccount_no from vaccount where account_no=$1 and va_type=$2 and is_delete='0' limit 1`,
		[]*sql.NullString{&vaccountNo}, accountNo, vaType)

	return vaccountNo.String, err
}

//------------------------------------------------------------------------------------------------------

// 冻结用户余额(即将balance的余额转到frozen_balance)
// 返回修改后的balance和frozen_balance的值
func (v *VaccountDao) FreezeBalance(tx *sql.Tx, vAccountNo string, amount string) (string, string, error) {
	sqlStr := `UPDATE vaccount SET balance=balance-$1, frozen_balance=frozen_balance+$1, modify_time=current_timestamp `
	sqlStr += ` WHERE vaccount_no=$2 AND is_delete='0' RETURNING balance, frozen_balance `

	var balance, frozenBalance sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&balance, &frozenBalance}, amount, vAccountNo)
	if nil != err {
		return "", "", err
	}

	return balance.String, frozenBalance.String, nil
}

// 解冻用户余额(即将frozen_balance的余额转到balance)
// 返回修改后的balance和frozen_balance的值
func (v *VaccountDao) UnfreezeBalance(tx *sql.Tx, vAccountNo string, amount string) (string, string, error) {
	sqlStr := `UPDATE vaccount SET balance=balance+$1, frozen_balance=frozen_balance-$1, modify_time=current_timestamp `
	sqlStr += ` WHERE vaccount_no=$2 AND is_delete='0' RETURNING balance, frozen_balance `

	var balance, frozenBalance sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&balance, &frozenBalance}, amount, vAccountNo)
	if nil != err {
		return "", "", err
	}

	return balance.String, frozenBalance.String, nil
}

// 增加用户余额(即将balance的余额增加对应金额)
// 返回修改后的balance和frozen_balance的值
func (v *VaccountDao) PlusBalance(tx *sql.Tx, vAccountNo string, amount string) (string, string, error) {
	sqlStr := `UPDATE vaccount SET balance=balance+$1, modify_time=current_timestamp `
	sqlStr += ` WHERE vaccount_no=$2 AND is_delete='0' RETURNING balance, frozen_balance `

	var balance, frozenBalance sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&balance, &frozenBalance}, amount, vAccountNo)
	if nil != err {
		return "", "", err
	}

	return balance.String, frozenBalance.String, nil
}

// 减少用户余额(即将balance的余额减少对应金额)
// 返回修改后的balance和frozen_balance的值
func (v *VaccountDao) MinusBalance(tx *sql.Tx, vAccountNo string, amount string) (string, string, error) {
	sqlStr := `UPDATE vaccount SET balance=balance-$1, modify_time=current_timestamp `
	sqlStr += ` WHERE vaccount_no=$2 AND is_delete='0' RETURNING balance, frozen_balance `

	var balance, frozenBalance sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&balance, &frozenBalance}, amount, vAccountNo)
	if nil != err {
		return "", "", err
	}

	return balance.String, frozenBalance.String, nil
}

// 减少用户冻结余额(即将frozen_balance的余额减少对应金额)
// 返回修改后的balance和frozen_balance的值
func (v *VaccountDao) MinusFrozenBalance(tx *sql.Tx, vAccountNo string, amount string) (string, string, error) {
	sqlStr := `UPDATE vaccount SET frozen_balance=frozen_balance-$1, modify_time=current_timestamp `
	sqlStr += ` WHERE vaccount_no=$2 AND is_delete='0' RETURNING balance, frozen_balance `

	var balance, frozenBalance sql.NullString
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&balance, &frozenBalance}, amount, vAccountNo)
	if nil != err {
		return "", "", err
	}

	return balance.String, frozenBalance.String, nil
}

//------------------------------------------------------------------------------------------------------

//获取账户余额
func (v *VaccountDao) GetBalanceByVAccNo(vAccountNo string) (int64, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return -1, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT balance FROM vaccount WHERE vaccount_no=$1 AND is_delete='0' LIMIT 1"
	var balance sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&balance}, vAccountNo)
	if err != nil {
		return -1, nil
	}

	return strext.ToInt64(balance.String), nil
}
func (v *VaccountDao) GetBalanceByAccNo(accountNo string, vAccType int) (int64, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return -1, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT balance FROM vaccount WHERE account_no=$1 AND va_type= $2 AND is_delete='0' LIMIT 1"
	var balance sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&balance}, accountNo, vAccType)
	if err != nil {
		return -1, nil
	}

	return strext.ToInt64(balance.String), nil
}

//同步虚账余额
func (v *VaccountDao) SyncAccRemain(tx *sql.Tx, accNo string) error {
	err := ss_sql.ExecTx(tx, `update account set usd_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_USD_DEBIT, constants.VaType_FREEZE_USD_DEBIT)
	if nil != err {
		return err
	}
	err = ss_sql.ExecTx(tx, `update account set khr_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0'`, accNo, constants.VaType_KHR_DEBIT, constants.VaType_FREEZE_KHR_DEBIT)
	if nil != err {
		return err
	}

	return nil
}
