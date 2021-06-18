package dao

import (
	"database/sql"
	"time"

	"a.a/mp-server/common/ss_count"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type BillingDetailsResultsDao struct {
}

var BillingDetailsResultsDaoInst BillingDetailsResultsDao

// 按日期统计总账户数量(统计某一天有多少服务商有账单信息)
func (BillingDetailsResultsDao) CountServicerByDate(date string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := date
	endTime, tErr := ss_time.TimeAfter(date, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return 0, tErr
	}

	var total sql.NullString
	sqlStr := "SELECT COUNT(DISTINCT(servicer_no)) FROM billing_details_results "
	sqlStr += " WHERE order_status=$1 AND create_time >= $2 AND create_time < $3 "

	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total},
		constants.OrderStatus_Paid, startTime, endTime,
	)

	if err != nil {
		return 0, err
	}

	return strext.ToInt(total.String), nil
}

// 按日期分页获取服务商编号
func (BillingDetailsResultsDao) GetServicerNoByDate(date string, page, pageSize int) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := date
	endTime, tErr := ss_time.TimeAfter(date, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	startNum := (page - 1) * pageSize

	sqlStr := "SELECT servicer_no FROM billing_details_results "
	sqlStr += " WHERE order_status=$1 AND create_time >= $2 AND create_time < $3 "
	sqlStr += " GROUP BY servicer_no ORDER BY servicer_no ASC LIMIT $4 OFFSET $5"

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, constants.OrderStatus_Paid, startTime, endTime, pageSize, startNum)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if err != nil {
		ss_log.Error("err=%v", err)
		return nil, err
	}

	list := []string{}

	for rows.Next() {
		var d sql.NullString
		err := rows.Scan(&d)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		list = append(list, d.String)
	}

	return list, nil
}

// 按日期统计指定服务商
func (BillingDetailsResultsDao) GetServicerStatis(servicerNo, currencyType, date string) (ServicerCheckListStatis, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var data ServicerCheckListStatis

	startTime := date
	endTime, tErr := ss_time.TimeAfter(date, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return data, tErr
	}

	sqlStr := "SELECT "

	// 存款
	sqlStr += " SUM(case WHEN bill_type=1 THEN 1 else 0 END) as in_num, "
	sqlStr += " SUM(case WHEN bill_type=1 THEN amount else 0 END) as in_amount, "

	// 取款
	sqlStr += " SUM(case WHEN bill_type=2 THEN 1 else 0 END) as out_num, "
	sqlStr += " SUM(case WHEN bill_type=2 THEN amount else 0 END) as out_amount, "

	// 收益
	sqlStr += " SUM(case WHEN bill_type=3 THEN 1 else 0 END) as profit_num, "
	sqlStr += " SUM(case WHEN bill_type=3 THEN amount else 0 END) as profit_amount, "

	// 充值
	sqlStr += " SUM(case WHEN bill_type=4 THEN 1 else 0 END) as recharge_num, "
	sqlStr += " SUM(case WHEN bill_type=4 THEN amount else 0 END) as recharge_amount, "

	// 提现
	sqlStr += " SUM(case WHEN bill_type=5 THEN 1 else 0 END) as withdraw_num, "
	sqlStr += " SUM(case WHEN bill_type=5 THEN amount else 0 END) as withdraw_amount "

	sqlStr += " FROM billing_details_results AS b "
	sqlStr += " WHERE order_status=$1 AND create_time >= $2 AND create_time < $3 AND servicer_no=$4 AND currency_type=$5 "
	sqlStr += " GROUP BY servicer_no "
	var inNum, inAmount, outNum, outAmount, profitNum, profitAmount, rechargeNum, rechargeAmount, withdrawNum, withdrawAmount sql.NullString

	qErr := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{
			&inNum, &inAmount,
			&outNum, &outAmount,
			&profitNum, &profitAmount,
			&rechargeNum, &rechargeAmount,
			&withdrawNum, &withdrawAmount,
		}, constants.OrderStatus_Paid, startTime, endTime, servicerNo, currencyType)

	if qErr != nil {
		if qErr.Error() != ss_sql.DB_NO_ROWS_MSG {
			return data, qErr
		}
	}

	data.ServicerNo = servicerNo
	data.CurrencyType = currencyType
	data.Dates = date
	data.InNum = strext.ToInt64(inNum.String)
	data.InAmount = strext.ToInt64(inAmount.String)
	data.OutNum = strext.ToInt64(outNum.String)
	data.OutAmount = strext.ToInt64(outAmount.String)
	data.ProfitNum = strext.ToInt64(profitNum.String)
	data.ProfitAmount = strext.ToInt64(profitAmount.String)
	data.RechargeNum = strext.ToInt64(rechargeNum.String)
	data.RechargeAmount = strext.ToInt64(rechargeAmount.String)
	data.WithdrawNum = strext.ToInt64(withdrawNum.String)
	data.WithdrawAmount = strext.ToInt64(withdrawAmount.String)

	return data, nil
}

func (BillingDetailsResultsDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM billing_details_results bdr " +
		" left join log_change_balance_order lcbo on lcbo.log_no = bdr.order_no " + whereModel.WhereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}
	return totalT.String
}

func (BillingDetailsResultsDao) GetSumAmt(whereStr string, whereArgs []interface{}) (sum string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//全部数量统计
	var sumT sql.NullString
	sqlCnt := "SELECT sum(amount) " +
		"FROM billing_details_results bdr " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&sumT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0"
	}

	if sumT.String == "" {
		return "0"
	}
	return sumT.String
}

func (BillingDetailsResultsDao) GetHaveDataTime(dbHandler *sql.DB, whereModelStr string, whereModelArgs []interface{}) (timeR []string, errR string) {

	sqlGetTime := " select distinct to_char(bdr.create_time,'yyyy-MM-dd') " +
		" from billing_details_results bdr " + whereModelStr
	var timeS []string
	rowsGetTime, stmtGetTime, errGetTime := ss_sql.Query(dbHandler, sqlGetTime, whereModelArgs...)
	if stmtGetTime != nil {
		defer stmtGetTime.Close()
	}
	defer rowsGetTime.Close()
	if errGetTime == nil {
		for rowsGetTime.Next() {
			var time sql.NullString
			errGetTime = rowsGetTime.Scan(
				&time,
			)
			if errGetTime != nil {
				ss_log.Error("errGetTime=[%v]", errGetTime)
				continue
			}
			timeS = append(timeS, time.String)
		}
	} else {
		ss_log.Error("err=[%v]", errGetTime)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	return timeS, ss_err.ERR_SUCCESS

}

func (BillingDetailsResultsDao) GetBillingDetailsResults(whereList []*model.WhereSqlCond, pageSize, page int) (datas []*go_micro_srv_cust.ServicerBillsData, errR string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, `order by bdr.create_time desc`)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)

	sqlStr := "SELECT bdr.bill_type, bdr.create_time, bdr.currency_type, bdr.order_no, bdr.amount," +
		" bdr.order_status, bdr.account_type, bdr.op_acc_no, lcbo.op_type  " +
		" FROM billing_details_results bdr " +
		" left join log_change_balance_order lcbo on lcbo.log_no = bdr.order_no " + whereModel.WhereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	for rows.Next() {
		data := &go_micro_srv_cust.ServicerBillsData{}
		var opType sql.NullString
		err = rows.Scan(
			&data.BillType,
			&data.CreateTime,
			&data.CurrencyType,
			&data.OrderNo,
			&data.Amount,
			&data.OrderStatus,
			&data.AccountType,
			&data.OpNo,
			&opType,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}

		switch data.BillType {
		case constants.BILL_TYPE_INCOME: //         用户存款
			data.OpType = constants.VaOpType_Minus
		case constants.BILL_TYPE_OUTGO: //          用户取款
			data.OpType = constants.VaOpType_Add
		case constants.BILL_TYPE_PROFIT: //         收益(佣金)
			data.OpType = constants.VaOpType_Add
		case constants.BILL_TYPE_RECHARGE: //       充值
			data.OpType = constants.VaOpType_Add
		case constants.BILL_TYPE_WITHDRAWALS: //    提现
			data.OpType = constants.VaOpType_Minus
		case constants.BILL_TYPE_ChangeBalance: //  平台修改余额:
			data.OpType = opType.String
		}

		datas = append(datas, data)
	}

	return datas, ss_err.ERR_SUCCESS
}

type BillsDaySum struct {
	UsdOutgoSum  string
	KhrOutgoSum  string
	UsdIncomeSum string
	KhrIncomeSum string
}

func (b *BillingDetailsResultsDao) GetServicerBillsDaySum(date string, whereList []*model.WhereSqlCond, billType string) BillsDaySum {
	ret := BillsDaySum{}
	ret.UsdOutgoSum = "0"
	ret.KhrOutgoSum = "0"
	ret.UsdIncomeSum = "0"
	ret.KhrIncomeSum = "0"

	if billType == constants.BILL_TYPE_OUTGO || billType == "" {
		//usd取款统计
		usdOutgoWhereModel := ss_sql.SsSqlFactoryInst.InitWhereList(append(whereList,
			&model.WhereSqlCond{Key: "bill_type", Val: constants.BILL_TYPE_OUTGO, EqType: "="},
			&model.WhereSqlCond{Key: "currency_type", Val: constants.CURRENCY_USD, EqType: "="},
			&model.WhereSqlCond{Key: "to_char(create_time,'yyyy-MM-dd')", Val: date, EqType: "="},
		))
		usdOutgoSum := b.GetSumAmt(usdOutgoWhereModel.WhereStr, usdOutgoWhereModel.Args)
		ret.UsdOutgoSum = ss_count.Add(ret.UsdOutgoSum, usdOutgoSum)

		//khr取款统计
		khrOutgoWhereModel := ss_sql.SsSqlFactoryInst.InitWhereList(append(whereList,
			&model.WhereSqlCond{Key: "bill_type", Val: constants.BILL_TYPE_OUTGO, EqType: "="},
			&model.WhereSqlCond{Key: "currency_type", Val: constants.CURRENCY_KHR, EqType: "="},
			&model.WhereSqlCond{Key: "to_char(create_time,'yyyy-MM-dd')", Val: date, EqType: "="},
		))
		khrOutgoSum := b.GetSumAmt(khrOutgoWhereModel.WhereStr, khrOutgoWhereModel.Args)
		ret.KhrOutgoSum = ss_count.Add(ret.KhrOutgoSum, khrOutgoSum)
	}

	if billType == constants.BILL_TYPE_INCOME || billType == "" {
		//usd收款统计
		usdIncomWhereModel := ss_sql.SsSqlFactoryInst.InitWhereList(append(whereList,
			&model.WhereSqlCond{Key: "bill_type", Val: constants.BILL_TYPE_INCOME, EqType: "="},
			&model.WhereSqlCond{Key: "currency_type", Val: constants.CURRENCY_USD, EqType: "="},
			&model.WhereSqlCond{Key: "to_char(create_time,'yyyy-MM-dd')", Val: date, EqType: "="},
		))
		usdIncomeSum := b.GetSumAmt(usdIncomWhereModel.WhereStr, usdIncomWhereModel.Args)
		ret.UsdIncomeSum = ss_count.Add(ret.UsdIncomeSum, usdIncomeSum)

		//khr收款统计
		khrIncomWhereModel := ss_sql.SsSqlFactoryInst.InitWhereList(append(whereList,
			&model.WhereSqlCond{Key: "bill_type", Val: constants.BILL_TYPE_INCOME, EqType: "="},
			&model.WhereSqlCond{Key: "currency_type", Val: constants.CURRENCY_KHR, EqType: "="},
			&model.WhereSqlCond{Key: "to_char(create_time,'yyyy-MM-dd')", Val: date, EqType: "="},
		))
		khrIncomeSum := b.GetSumAmt(khrIncomWhereModel.WhereStr, khrIncomWhereModel.Args)
		ret.KhrIncomeSum = ss_count.Add(ret.KhrIncomeSum, khrIncomeSum)

	}

	return ret
}

func (BillingDetailsResultsDao) InsertResult(amount, balanceType, accountNo, accountType, orderNo, balance, orderStatus string, billType int, fees string, realAmount string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	logNo := strext.GetDailyId()
	err := ss_sql.Exec(dbHandler, `insert into billing_details_results(bill_no,amount,currency_type,bill_type,account_no,account_type,`+
		`order_no,balance,order_status,create_time, fees, real_amount) values($1,$2,$3,$4,$5,$6,$7,$8,$9,current_timestamp, $10, $11)`,
		logNo, amount, balanceType, billType, accountNo, accountType, orderNo, balance, orderStatus, fees, realAmount)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNo
}
