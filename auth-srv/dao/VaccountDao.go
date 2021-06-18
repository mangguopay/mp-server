package dao

import (
	"database/sql"
	"errors"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type VaccountDao struct {
	VAccountNo string
	AccountNo  string
	VAccType   string
	Balance    string
}

var VaccountDaoInst VaccountDao

//获取虚账余额
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
	return ss_err.ERR_SUCCESS, balanceT.String
}

//查询未激活的虚账
func (VaccountDao) GetFreezeVAccountNo(accountNo string) ([]*VaccountDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select vaccount_no,va_type,balance from vaccount where account_no=$1 and va_type in ($2, $3) and is_delete='0' "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, accountNo, constants.VaType_FREEZE_USD_DEBIT, constants.VaType_FREEZE_KHR_DEBIT)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*VaccountDao
	for rows.Next() {
		var vAccountNo, vAccType, balance sql.NullString
		err := rows.Scan(&vAccountNo, &vAccType, &balance)
		if err != nil {
			return nil, err
		}
		data := new(VaccountDao)
		data.VAccountNo = vAccountNo.String
		data.VAccType = vAccType.String
		data.Balance = balance.String
		dataList = append(dataList, data)
	}

	return dataList, nil
}

//查询账号所有的虚账
func (VaccountDao) GetVAccNoByAccountNo(accountNo string) ([]*VaccountDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT vaccount_no,va_type FROM vaccount WHERE account_no=$1 AND is_delete='0' "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, accountNo)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*VaccountDao
	for rows.Next() {
		var vAccNo, vAccType sql.NullString
		err := rows.Scan(&vAccNo, &vAccType)
		if err != nil {
			return nil, err
		}
		data := new(VaccountDao)
		data.VAccountNo = vAccNo.String
		data.VAccType = vAccType.String

		dataList = append(dataList, data)
	}
	return dataList, nil
}

//创建虚账
func (VaccountDao) InitVAccountNoTx(tx *sql.Tx, accountNo string, vaType int32, balanceType string) (string, error) {
	vAccountNo := strext.NewUUID()
	sqlStr := "insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance,balance_type) "
	sqlStr += "values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8)"
	err := ss_sql.ExecTx(tx, sqlStr, vAccountNo, accountNo, vaType, "0", "0", "1", "0", balanceType)
	if err != nil {
		return "", err
	}
	return vAccountNo, nil
}

//创建虚账
func (VaccountDao) InitVAccountNo(accountNo string, vaType int32, balanceType string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	vAccountNo := strext.NewUUID()
	sqlStr := "insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance,balance_type) "
	sqlStr += "values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8)"
	err := ss_sql.Exec(dbHandler, sqlStr, vAccountNo, accountNo, vaType, "0", "0", "1", "0", balanceType)
	if err != nil {
		return "", err
	}
	return vAccountNo, nil
}

//减少虚账余额，返回修改后的余额
func (v *VaccountDao) MinusBalance(tx *sql.Tx, vAccountNo string, amount string) (string, string, error) {
	var balance, frozenBalance sql.NullString
	sqlStr := `UPDATE vaccount SET balance=balance-$1, modify_time=current_timestamp `
	sqlStr += ` WHERE vaccount_no=$2 AND is_delete='0' RETURNING balance, frozen_balance `
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&balance, &frozenBalance}, amount, vAccountNo)
	if err != nil {
		return "", "", err
	}
	return balance.String, frozenBalance.String, err
}

//增加虚账余额,返回修改后的余额
func (v *VaccountDao) PlusBalance(tx *sql.Tx, vAccountNo string, amount string) (string, string, error) {
	var balance, frozenBalance sql.NullString
	sqlStr := `UPDATE vaccount SET balance=balance+$1, modify_time=current_timestamp `
	sqlStr += ` WHERE vaccount_no=$2 AND is_delete='0' RETURNING balance, frozen_balance `
	err := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&balance, &frozenBalance}, amount, vAccountNo)
	if err != nil {
		return "", "", err
	}
	return balance.String, frozenBalance.String, err
}

//同步虚账金额和账号余额
func (v *VaccountDao) SyncAccRemain(tx *sql.Tx, accNo string) error {
	sqlStr1 := "update account set usd_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0' "
	err := ss_sql.ExecTx(tx, sqlStr1, accNo, constants.VaType_USD_DEBIT, constants.VaType_FREEZE_USD_DEBIT)
	if nil != err {
		return err
	}

	sqlStr2 := "update account set khr_balance=(select sum(balance) from vaccount where account_no=$1 and va_type in($2,$3) and is_delete='0') where uid=$1 and is_delete='0' "
	err = ss_sql.ExecTx(tx, sqlStr2, accNo, constants.VaType_KHR_DEBIT, constants.VaType_FREEZE_KHR_DEBIT)
	if nil != err {
		return err
	}

	return nil
}
