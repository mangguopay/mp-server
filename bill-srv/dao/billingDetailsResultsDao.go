package dao

import (
	"context"
	"database/sql"
	"errors"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	tmProto "a.a/mp-server/common/proto/tm"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/common/ss_struct"
)

type BillingDetailsResultsDao struct{}

var BillingDetailsResultsDaoInstance BillingDetailsResultsDao

func (BillingDetailsResultsDao) InsertResultV2(tmProxy *ss_struct.TmServerProxy, amount, balanceType, accountNo, accountType, orderNo, balance, orderStatus, servicerNo, opAccNo string, billType int, fees string, realAmount string) error {
	logNo := strext.GetDailyId()
	sql := `insert into billing_details_results(bill_no,amount,currency_type,bill_type,account_no,account_type,` +
		`order_no,balance,order_status,create_time,servicer_no,op_acc_no, fees, real_amount) values($1,$2,$3,$4,$5,$6,$7,$8,$9,current_timestamp,$10,$11, $12, $13)`
	args := []string{logNo, amount, balanceType, strext.ToStringNoPoint(billType), accountNo, accountType, orderNo, balance, orderStatus, servicerNo, opAccNo, fees, realAmount}
	rsp, err := tmProxy.GetTmServer().TxExec(context.TODO(), &tmProto.TxExecRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return err
	}
	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}

	return nil
}
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

// 通过logNo查询实际金额
func (BillingDetailsResultsDao) GetRealAmountByOutOrderLogNo(logNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select real_amount from billing_details_results where order_no=$1 AND bill_type=$2"
	var realAmount sql.NullString

	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&realAmount}, logNo, constants.BillDetailTypeOut)
	if err != nil {
		return "", err
	}

	return realAmount.String, nil
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

//获取产生该笔订单的账号
func (*BillingDetailsResultsDao) GetAccountByOrderNo(orderNo string) (account string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "order_no", Val: orderNo, EqType: "="},
		{Key: "bill_type", Val: "('" + constants.BILL_TYPE_INCOME + "','" + constants.BILL_TYPE_OUTGO + "') ", EqType: "in"},
	})

	//sqlStr := "select bdr.account_type, op_acc_no " +
	sqlStr := "select acc.account " +
		" from billing_details_results bdr " +
		" left join rela_acc_iden rai on rai.iden_no = bdr.op_acc_no " + //and rai.account_type = bdr.account_type" +
		" left join account acc on acc.uid = rai.account_no " + whereModel.WhereStr
	var accountT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accountT}, whereModel.Args...)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", err
	}

	return accountT.String, nil
}
