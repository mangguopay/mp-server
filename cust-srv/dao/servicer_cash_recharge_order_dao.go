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

type CashRechargeOrderDao struct {
}

var CashRechargeOrderDaoInst CashRechargeOrderDao

func (*CashRechargeOrderDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlCnt := "select count(1) " +
		" from servicer_cash_recharge_order scro " +
		" left join account acc on acc.uid = scro.acc_no" +
		" left join account acc2 on acc2.uid = scro.op_acc_no" + whereModel.WhereStr
	var totalT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "0"
	}

	return totalT.String
}

func (*CashRechargeOrderDao) GetCashRechargeOrders(whereList []*model.WhereSqlCond, page, pageSize int) (datas []*go_micro_srv_cust.CashRechargeOrderData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY scro.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))
	sqlStr := "select scro.log_no, scro.amount, scro.create_time, scro.order_status, scro.currency_type, scro.payment_type, scro.notes " +
		",scro.acc_no, acc.account, scro.op_acc_no, acc2.account " +
		" from servicer_cash_recharge_order scro " +
		" left join account acc on acc.uid = scro.acc_no " +
		" left join admin_account acc2 on acc2.uid = scro.op_acc_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		return nil, err2
	}

	var datasT []*go_micro_srv_cust.CashRechargeOrderData
	for rows.Next() {
		data := go_micro_srv_cust.CashRechargeOrderData{}
		err2 = rows.Scan(
			&data.LogNo,
			&data.Amount,
			&data.CreateTime,
			&data.OrderStatus,
			&data.CurrencyType,

			&data.PaymentType,
			&data.Notes,

			&data.AccNo,
			&data.AccAccount,
			&data.OpAccNo,
			&data.OpAccAccount,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			continue
		}

		datasT = append(datasT, &data)
	}

	return datasT, nil
}
