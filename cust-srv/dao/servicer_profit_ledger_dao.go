package dao

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ServicerProfitLedgerDao struct {
}

var ServicerProfitLedgerDaoInst ServicerProfitLedgerDao

func (*ServicerProfitLedgerDao) GetServicerProfitLedgerCount(dbHandler *sql.DB, startTimeStr, endTimeStr string) (datas []*go_micro_srv_cust.GetServicerOrderCountData, err string) {
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "payment_time", Val: startTimeStr, EqType: ">="},
		{Key: "payment_time", Val: endTimeStr, EqType: "<="},
	})
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " group by servicer_no, currency_type ")

	sqlStr := "select sum(actual_income), servicer_no, currency_type from servicer_profit_ledger " + whereModel.WhereStr
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	var datasT []*go_micro_srv_cust.GetServicerOrderCountData
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	} else {
		for rows.Next() {
			var data go_micro_srv_cust.GetServicerOrderCountData
			errT = rows.Scan(
				&data.ProfitAmountSum,
				&data.ServicerNo,
				&data.BalanceType,
			)
			if errT != nil {
				ss_log.Error("err=[%v]", errT)
				return nil, ss_err.ERR_SYS_DB_GET
			}
			data.IncomeAmountSum = "0"
			data.IncomeTotalSum = "0"
			data.OutgoAmountSum = "0"
			data.OutgoTotalSum = "0"
			datasT = append(datasT, &data)
		}
	}

	return datasT, ss_err.ERR_SUCCESS

}
