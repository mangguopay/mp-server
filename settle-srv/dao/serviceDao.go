package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	ss_sql2 "a.a/cu/ss_sql"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type ServiceDao struct {
}

var ServiceDaoInst ServiceDao

func (ServiceDao) GetServicerPWDFromOpAccNo(sid string) (servierNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select password from servicer where servicer_no =$1 and is_delete='0' limit 1`, []*sql.NullString{&pwdT}, sid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return pwdT.String
}
func (ServiceDao) GetAccNoFromSrvNo(sid string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_no from servicer where servicer_no =$1 and is_delete='0' limit 1`, []*sql.NullString{&accountNoT}, sid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountNoT.String
}

func (ServiceDao) GetSharingFromSrvNo(srvAccNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var incomeSharing, outGoSharing sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select income_sharing,commission_sharing from servicer where account_no =$1 and is_delete='0' limit 1`, []*sql.NullString{&incomeSharing, &outGoSharing}, srvAccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return incomeSharing.String, outGoSharing.String
}

func (ServiceDao) GetServicerNoByCashierNo(cid string) (servicerNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var Sno sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select servicer_no from cashier where uid =$1 and is_delete='0' limit 1`, []*sql.NullString{&Sno}, cid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return Sno.String
}

// todo 没有测试，新版本接口改写
func (ServiceDao) GetTransferToServicerLogs(pageSize, page int32, startTime, endTime, orderStatus, servicerNo string) (returnDatas []*go_micro_srv_bill.TransferToServicerLogData, returnTotal string, returnErr string) {
	dbHandler := ss_sql2.NewDbInst(constants.DB_CRM)
	defer dbHandler.Close()

	datas := []*go_micro_srv_bill.TransferToServicerLogData{}
	dbHandler.InitWhereList([]*ss_sql2.WhereSqlCond{
		{Key: "lts.create_time", Val: startTime, EqType: ">="},
		{Key: "lts.create_time", Val: endTime, EqType: "<="},
		{Key: "lts.order_type", Val: "1", EqType: "="},
		{Key: "lts.order_status", Val: orderStatus, EqType: "="},
		{Key: "lts.servicer_no", Val: servicerNo, EqType: "="},
	})

	total := strext.ToStringNoPoint(dbHandler.GetCnt("log_to_servicer lts"))

	dbHandler.AppendWhereExtra(`order by lts.create_time desc`)
	dbHandler.AppendWhereLimitI32(pageSize, page)
	sqlStr := "SELECT lts.currency_type, lts.amount, lts.order_status, lts.card_no, lts.finish_time, ch.channel_name, c.name, c.card_number" +
		" FROM log_to_servicer lts " +
		" LEFT JOIN card c ON c.card_no = lts.card_no " +
		" LEFT JOIN channel ch ON ch.channel_no = c.channel_no "
	rows, err := dbHandler.QueryWhere(sqlStr)
	if err != nil {
		ss_log.Error("err=%v", err)
		return nil, "0", ss_err.ERR_SYS_DB_GET
	}

	for rows.Next() {
		var data go_micro_srv_bill.TransferToServicerLogData
		err = rows.Scan(
			&data.CurrencyType,
			&data.Amount,
			&data.OrderStatus,
			&data.CardNo,
			&data.FinishTime,

			&data.ChannelName,
			&data.Name,
			&data.CardNumber,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		datas = append(datas, &data)
	}

	return datas, total, ss_err.ERR_SUCCESS
}

//func (ServiceDao) GetServicerProfitLedgers(startTime, endTime, servicerNo, currencyType string) (datas []*go_micro_srv_bill.ServicerProfitLedgersData, returnTotals string, returnErr string) {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//
//	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*dao.WhereSqlCond{
//		{Key: "payment_time", Val: startTime, EqType: ">="},
//		{Key: "payment_time", Val: endTime, EqType: "<="},
//		{Key: "servicer_no", Val: servicerNo, EqType: "="},
//		{Key: "currency_type", Val: currencyType, EqType: "="},
//	})
//	//统计
//	where := whereModel.WhereStr
//	args := whereModel.Args
//	var total sql.NullString
//	sqlCnt := "SELECT count(1) FROM servicer_profit_ledger " + where
//	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, args...)
//	if err != nil {
//		ss_log.Error("err=%v", err)
//	}
//	//添加limit
//	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by payment_time desc`)
//
//	where2 := whereModel.WhereStr
//	args2 := whereModel.Args
//	sqlStr := "SELECT log_no, amount_order, actual_income, payment_time, currency_type" +
//		" FROM servicer_profit_ledger " + where2
//	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
//	if stmt != nil {
//		defer stmt.Close()
//	}
//	defer rows.Close()
//
//	if err == nil {
//		for rows.Next() {
//			var data go_micro_srv_bill.ServicerProfitLedgersData
//			err = rows.Scan(
//				&data.LogNo,
//				&data.AmountOrder,
//				&data.ActualIncome,
//				&data.PaymentTime,
//				&data.CurrencyType,
//			)
//			if err != nil {
//				ss_log.Error("err=%v", err)
//				continue
//			}
//			datas = append(datas, &data)
//		}
//	} else {
//		ss_log.Error("ServiceDao | GetServicerProfitLedgers | err=%v\n\nsql=[%v]", err, sqlStr)
//		return nil, "0", ss_err.ERR_SYS_DB_GET
//	}
//
//	return datas, total.String, ss_err.ERR_SUCCESS
//}

//func (ServiceDao) GetServicerProfitLedgerDetail(logNo string) (data *go_micro_srv_bill.ServicerProfitLedgersData, returnErr string) {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//
//	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*dao.WhereSqlCond{
//		{Key: "log_no", Val: logNo, EqType: "="},
//	})
//
//	where2 := whereModel.WhereStr
//	args2 := whereModel.Args
//	sqlStr := "SELECT amount_order, actual_income, split_proportion, payment_time, log_no " +
//		" FROM servicer_profit_ledger " + where2
//	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
//	if stmt != nil {
//		defer stmt.Close()
//	}
//	defer rows.Close()
//	data = &go_micro_srv_bill.ServicerProfitLedgersData{}
//	if err == nil {
//		for rows.Next() {
//			err = rows.Scan(
//				&data.AmountOrder,
//				&data.ActualIncome,
//				&data.SplitProportion,
//				&data.PaymentTime,
//				&data.LogNo,
//			)
//			if err != nil {
//				ss_log.Error("err=%v", err)
//				continue
//			}
//		}
//	} else {
//		ss_log.Error("ServiceDao | GetServicerProfitLedgerDetail | err=%v\n\nsql=[%v]", err, sqlStr)
//		return nil, ss_err.ERR_SYS_DB_GET
//	}
//	return data, ss_err.ERR_SUCCESS
//}
