package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type BillingDetailsResultsDao struct{}

var BillingDetailsResultsDaoInstance BillingDetailsResultsDao

func (BillingDetailsResultsDao) InsertResult(tx *sql.Tx, amount, balanceType, accountNo, accountType, orderNo, balance, orderStatus, servicerNo, opAccNo string, billType int, fees string, realAmount string) string {
	logNo := strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into billing_details_results(bill_no,amount,currency_type,bill_type,account_no,account_type,`+
		`order_no,balance,order_status,create_time,servicer_no,op_acc_no, fees, real_amount) values($1,$2,$3,$4,$5,$6,$7,$8,$9,current_timestamp,$10,$11, $12, $13)`,
		logNo, amount, balanceType, billType, accountNo, accountType, orderNo, balance, orderStatus, servicerNo, opAccNo, fees, realAmount)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNo
}

func (*BillingDetailsResultsDao) UpdateOrderStatusFromLogNo(tx *sql.Tx, orderStatus, logNo string) (errCode string) {

	err := ss_sql.ExecTx(tx, `update billing_details_results set order_status = $1,  modify_time =  current_timestamp where order_no=$2`, orderStatus, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (BillingDetailsResultsDao) GetSum(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}) (sum string) {
	sqlStr := "select sum(bdr.amount) from billing_details_results bdr " + whereModelStr
	var sumT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&sumT}, whereModelArgs...)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "0"
	}

	if sumT.String == "" {
		return "0"
	} else {
		return sumT.String
	}

}
