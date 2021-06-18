package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

type HeadquartersProfit struct {
	LogNo        string
	OrderNo      string
	Amount       string
	OrderStatus  string
	BalanceType  string
	ProfitSource string
	OpType       string
	CreateTime   string
	FinishTime   string
}

var HeadquartersProfitDao HeadquartersProfit

//插入总部利润
func (*HeadquartersProfit) InsertHeadquartersProfit(tx *sql.Tx, d *HeadquartersProfit) (string, error) {
	logNo := strext.GetDailyId()
	sqlStr := "insert into headquarters_profit " +
		"(log_no, general_ledger_no, amount,order_status, balance_type, profit_source, op_type, create_time, finish_time) " +
		"values($1,$2,$3,$4,$5,$6,$7,current_timestamp,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr,
		logNo, d.OrderNo, d.Amount, d.OrderStatus, d.BalanceType, d.ProfitSource, d.OpType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", err
	}

	return logNo, nil
}

//查询日志列表
func (*HeadquartersProfit) GetList(whereStr string, args []interface{}) ([]*HeadquartersProfit, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT log_no, general_ledger_no, amount, create_time, order_status, finish_time, balance_type, profit_source, op_type " +
		" FROM headquarters_profit " + whereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var datas []*HeadquartersProfit
	for rows.Next() {
		var LogNo, OrderNo, Amount, CreateTime, OrderStatus, FinishTime, BalanceType, ProfitSource, OpType sql.NullString
		err = rows.Scan(&LogNo, &OrderNo, &Amount, &CreateTime, &OrderStatus, &FinishTime, &BalanceType, &ProfitSource, &OpType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return nil, err
		}
		data := new(HeadquartersProfit)
		data.LogNo = LogNo.String
		data.OrderNo = OrderNo.String
		data.Amount = Amount.String
		data.CreateTime = CreateTime.String
		data.OrderStatus = OrderStatus.String
		data.FinishTime = FinishTime.String
		data.BalanceType = BalanceType.String
		data.ProfitSource = ProfitSource.String
		data.OpType = OpType.String
		datas = append(datas, data)
	}
	return datas, nil
}

//统计日志数量
func (*HeadquartersProfit) CountLog(whereStr string, args []interface{}) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var num sql.NullString
	sqlStr := "SELECT count(1) FROM headquarters_profit" + whereStr
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num}, args...)
	if err != nil {
		return "", err
	}
	return num.String, nil
}

//收益统计
func (*HeadquartersProfit) CountProfit(whereStr string, args []interface{}) (totalNum, totalAmount string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var num, amount sql.NullString
	sqlStr := "SELECT count(1),sum(amount) FROM headquarters_profit" + whereStr
	err = ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount}, args...)
	if err != nil {
		return "", "", err
	}
	return num.String, amount.String, nil
}
