package dao

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type HeadquartersProfitWithdrawDao struct {
}

var HeadquartersProfitWithdrawDaoInst HeadquartersProfitWithdrawDao

//查询平台利润可提现余额
func (HeadquartersProfitWithdrawDao) GetProfitCashable(tx *sql.Tx, currencyType string) (amount, err string) {

	sqlStr := "select cashable_balance from headquarters_profit_cashable where currency_type = $1 "
	var amountT sql.NullString
	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&amountT}, currencyType)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", ss_err.ERR_PARAM
	}
	return amountT.String, ss_err.ERR_SUCCESS
}

func (HeadquartersProfitWithdrawDao) ModifyProfitCashable(tx *sql.Tx, amount, currencyType string) (err string) {
	errT := ss_sql.ExecTx(tx, `update headquarters_profit_cashable set cashable_balance = cashable_balance-$1,modify_time = current_timestamp where currency_type = $2 `, amount, currencyType)
	if nil != errT {
		ss_log.Error("errT=[%v]", errT)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	return ss_err.ERR_SUCCESS
}
