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
}

var ChangeBalanceOrderDaoInst ChangeBalanceOrderDao

func (*ChangeBalanceOrderDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlCnt := "select count(1) " +
		" from log_change_balance_order cbo " +
		" left join account acc on acc.uid = cbo.account_no " + whereModel.WhereStr
	var totalT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "0"
	}

	return totalT.String
}

func (*ChangeBalanceOrderDao) GetChangeBalanceOrders(whereList []*model.WhereSqlCond, page, pageSize int) (datas []*go_micro_srv_cust.ChangeBalanceOrderData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY cbo.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)

	sqlStr := "select cbo.log_no, cbo.account_no, cbo.currency_type, cbo.before_balance, cbo.change_amount," +
		" cbo.after_balance, cbo.change_reason, cbo.order_status, cbo.account_type, cbo.op_type, " +
		" cbo.create_time, acc.account  " +
		" from log_change_balance_order cbo " +
		" left join account acc on acc.uid = cbo.account_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		return nil, err2
	}

	var datasT []*go_micro_srv_cust.ChangeBalanceOrderData
	for rows.Next() {
		data := go_micro_srv_cust.ChangeBalanceOrderData{}
		var account sql.NullString
		err2 = rows.Scan(
			&data.LogNo,
			&data.AccountNo,
			&data.CurrencyType,
			&data.BeforeBalance,
			&data.ChangeAmount,

			&data.AfterBalance,
			&data.ChangeReason,
			&data.OrderStatus,
			&data.AccountType,
			&data.OpType,

			&data.CreateTime,
			&account,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return nil, err2
		}
		data.Account = account.String
		datasT = append(datasT, &data)
	}

	return datasT, nil
}

func (*ChangeBalanceOrderDao) GetChangeBalanceOrderDetail(whereList []*model.WhereSqlCond) (data *go_micro_srv_cust.ChangeBalanceOrderData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

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

func (*ChangeBalanceOrderDao) AddChangeBalanceOrder(tx *sql.Tx, data ChangeBalanceOrderDao) (logNo string, err error) {
	logNo = strext.GetDailyId()
	sqlCnt := "insert into log_change_balance_order(log_no, account_no, currency_type, before_balance, change_amount," +
		" after_balance, change_reason, order_status, account_type, admin_account_no, create_time) " +
		" values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,current_timestamp) "
	if err := ss_sql.ExecTx(tx, sqlCnt, logNo, data.AccountNo, data.CurrencyType, data.BeforeBalance, data.ChangeAmount,
		data.AfterBalance, data.ChangeReason, data.OrderStatus, data.AccountType, data.AdminAccountNo); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", err
	}

	return logNo, nil
}
