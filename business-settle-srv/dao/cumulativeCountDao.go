package dao

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type CumulativeCountDao struct{}

var CumulativeCountDaoInst CumulativeCountDao

// 成功交易总金额(amount),平台利润,商户收入总金额
func (CumulativeCountDao) ModifyCumulative1(tx *sql.Tx, successAmountSum, headquartersProfit, merchantAmountSum string) string {

	err := ss_sql.ExecTx(tx, `update business_cumulative_count set success_amount_sum = success_amount_sum + $1, 
					headquarters_profit = headquarters_profit + $2, merchant_amount_sum = merchant_amount_sum +$3,modify_time = current_timestamp `,
		successAmountSum, headquartersProfit, merchantAmountSum)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PAY_UPDATE_ORDER
	}

	return ss_err.ERR_SUCCESS
}

// 修改平台利润
func (CumulativeCountDao) ModifyCumulative2(tx *sql.Tx, agencyProfit string) string {

	err := ss_sql.ExecTx(tx, `update business_cumulative_count set agency_profit = agency_profit + $1,modify_time = current_timestamp`,
		agencyProfit)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PAY_UPDATE_ORDER
	}

	return ss_err.ERR_SUCCESS
}

func (CumulativeCountDao) ModifyHeadquartersProfit(tx *sql.Tx, HeadquartersProfit string) string {

	err := ss_sql.ExecTx(tx, `update business_cumulative_count set headquarters_profit = headquarters_profit + $1,modify_time = current_timestamp`,
		HeadquartersProfit)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PAY_UPDATE_ORDER
	}

	return ss_err.ERR_SUCCESS
}
