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
