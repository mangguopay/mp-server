package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ChangeBalanceOrderDao struct {
	LogNo          string
	AccountNo      string
	CurrencyType   string
	BeforeBalance  string
	ChangeAmount   string
	AfterBalance   string
	ChangeReason   string
	OrderStatus    string
	AccountType    string
	AdminAccountNo string
	OpType         string
}

var ChangeBalanceOrderDaoInst ChangeBalanceOrderDao

func (*ChangeBalanceOrderDao) AddChangeBalanceOrder(tx *sql.Tx, data ChangeBalanceOrderDao) (logNo string, err error) {
	logNo = strext.GetDailyId()
	sqlCnt := "insert into log_change_balance_order(log_no, account_no, currency_type, before_balance, change_amount," +
		" after_balance, change_reason, order_status, account_type, admin_account_no, op_type, create_time) " +
		" values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,current_timestamp) "
	if err := ss_sql.ExecTx(tx, sqlCnt, logNo, data.AccountNo, data.CurrencyType, data.BeforeBalance, data.ChangeAmount,
		data.AfterBalance, data.ChangeReason, data.OrderStatus, data.AccountType, data.AdminAccountNo, data.OpType); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", err
	}

	return logNo, nil
}

func (*ChangeBalanceOrderDao) GetChangeBalanceOrderDetail(logNo string) (data *go_micro_srv_cust.ChangeBalanceOrderData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "cbo.log_no", Val: logNo, EqType: "="},
	})

	sqlStr := "select cbo.log_no, cbo.account_no, cbo.currency_type, cbo.before_balance, cbo.change_amount," +
		" cbo.after_balance, cbo.change_reason, cbo.order_status, cbo.account_type, cbo.op_type, " +
		" cbo.create_time, acc.account  " +
		" from log_change_balance_order cbo " +
		" left join account acc on acc.uid = cbo.account_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		return nil, err2
	}

	dataT := &go_micro_srv_cust.ChangeBalanceOrderData{}
	var account sql.NullString
	err2 = rows.Scan(
		&dataT.LogNo,
		&dataT.AccountNo,
		&dataT.CurrencyType,
		&dataT.BeforeBalance,
		&dataT.ChangeAmount,

		&dataT.AfterBalance,
		&dataT.ChangeReason,
		&dataT.OrderStatus,
		&dataT.AccountType,
		&dataT.OpType,

		&dataT.CreateTime,
		&account,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}
	dataT.Account = account.String

	return dataT, nil
}
