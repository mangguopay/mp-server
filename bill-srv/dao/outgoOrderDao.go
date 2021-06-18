package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/m"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type OutgoOrderDao struct {
}

var OutgoOrderDaoInst OutgoOrderDao

func (OutgoOrderDao) InsertOutgo(tx *sql.Tx, recvVaccNo, amount, servicerNo, opAccNo, moneyType, fees, rate, realAmount, lat, lng, ip string, withdrawType, opAccType int32) (logNo string) {
	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into outgo_order(log_no,vaccount_no,amount,create_time,order_status,`+
		`balance_type,fees,servicer_no,op_acc_no,rate,is_count,withdraw_type,real_amount,op_acc_type,lat, lng, ip) values($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8,$9,0,$10,$11,$12,$13,$14,$15)`,
		logNo, recvVaccNo, amount, constants.OrderStatus_Pending, moneyType, fees, servicerNo, opAccNo, rate, withdrawType, realAmount, opAccType, lat, lng, ip)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return logNo
}
func (OutgoOrderDao) InsertOutgoV2(recvVaccNo, amount, servicerNo, opAccNo, moneyType, fees, rate, realAmount, lat, lng, ip string, withdrawType, opAccType int32) (logNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	logNo = strext.GetDailyId()
	err := ss_sql.Exec(dbHandler, `insert into outgo_order(log_no,vaccount_no,amount,create_time,order_status,`+
		`balance_type,fees,servicer_no,op_acc_no,rate,is_count,withdraw_type,real_amount,op_acc_type,lat, lng, ip) values($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8,$9,0,$10,$11,$12,$13,$14,$15)`,
		logNo, recvVaccNo, amount, constants.OrderStatus_Pending, moneyType, fees, servicerNo, opAccNo, rate, withdrawType, realAmount, opAccType, lat, lng, ip)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return logNo
}

func (OutgoOrderDao) GetLogNoFromCode(code string) (logNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var logNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select log_no from outgo_order where write_off=$1 limit 1`,
		[]*sql.NullString{&logNoT}, code)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return logNoT.String
}

func (OutgoOrderDao) UpdateOutgoOrderStatusTx(tx *sql.Tx, logNo, orderStatus string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update outgo_order set order_status=$1,modify_time=current_timestamp where log_no=$2`,
		orderStatus, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_OUT_MONEY
	}
	return ss_err.ERR_SUCCESS
}
func (OutgoOrderDao) UpdateOutgoOrderStatus(logNo, orderStatus string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update outgo_order set order_status=$1,modify_time=current_timestamp where log_no=$2`,
		orderStatus, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_OUT_MONEY
	}
	return ss_err.ERR_SUCCESS
}

func (OutgoOrderDao) UpdateOutgoOrderRiskNo(tx *sql.Tx, logNo, riskNo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update outgo_order set risk_no=$1,modify_time=current_timestamp where log_no=$2`,
		riskNo, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_OUT_MONEY
	}
	return ss_err.ERR_SUCCESS
}
func (OutgoOrderDao) CancelOutgoOrder(tx *sql.Tx, cancelReason, logNo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update outgo_order set order_status='6',cancel_reason = $1,modify_time=current_timestamp where log_no=$2`,
		cancelReason, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_OUT_MONEY
	}
	return ss_err.ERR_SUCCESS
}

func (OutgoOrderDao) GetAmountFromLogNo(logNo, status string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var amountT, createTimeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select amount,create_time from outgo_order where log_no=$1 and order_status = $2 limit 1`,
		[]*sql.NullString{&amountT, &createTimeT}, logNo, status)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", ""
	}

	return amountT.String, createTimeT.String
}

func (OutgoOrderDao) GetRiskNoFromLogNo(logNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var riskNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select risk_no  from outgo_order where log_no=$1  limit 1`,
		[]*sql.NullString{&riskNoT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", err
	}

	return riskNoT.String, nil
}
func (OutgoOrderDao) GetOutGoDetailFromLogNo(logNo string) (string, string, string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var amountT, createTimeT, feesT, moneyTypeT, vaccountNoT, withdrawTypeT, orderStatusT, riskNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select amount,create_time,fees,balance_type,vaccount_no,withdraw_type,order_status,risk_no from 
						outgo_order where log_no=$1   limit 1`,
		[]*sql.NullString{&amountT, &createTimeT, &feesT, &moneyTypeT, &vaccountNoT, &withdrawTypeT, &orderStatusT, &riskNoT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", "", "", "", ""
	}

	return amountT.String, createTimeT.String, feesT.String, moneyTypeT.String, vaccountNoT.String, withdrawTypeT.String, orderStatusT.String, riskNoT.String
}

func (OutgoOrderDao) GetOutGoStatusFromLogNo(logNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var orderStatusT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select  order_status  from  outgo_order where log_no=$1   limit 1`,
		[]*sql.NullString{&orderStatusT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return orderStatusT.String
}

func (OutgoOrderDao) QueryOutGoOrderFromLogNo(logNo, status string) (string, string, string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var amountT, serviceNoT, finishTime, vaccountNoT, feesT, balanceTypeT, statusT, withdrawTypeT sql.NullString
	if status != "" {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,create_time,vaccount_no,fees,balance_type,withdraw_type from outgo_order where log_no=$1 and order_status = $2 limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &vaccountNoT, &feesT, &balanceTypeT, &withdrawTypeT}, logNo, status)
		if nil != err {
			ss_log.Error("err=%v", err)
			return "", "", "", "", "", "", "", ""
		}
	} else {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,finish_time,vaccount_no,fees,balance_type,order_status,withdraw_type from outgo_order where log_no=$1 limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &vaccountNoT, &feesT, &balanceTypeT, &statusT, &withdrawTypeT}, logNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return "", "", "", "", "", "", "", ""
		}
	}

	return amountT.String, serviceNoT.String, finishTime.String, vaccountNoT.String, feesT.String, balanceTypeT.String, statusT.String, withdrawTypeT.String
}

func (*OutgoOrderDao) CustOutgoBillsDetail(logNo string) (data *go_micro_srv_bill.CustOutgoBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "oor.log_no", Val: logNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT oor.log_no, oor.amount, oor.order_status, oor.payment_type, oor.create_time, oor.modify_time, oor.fees, oor.balance_type " +
		", lv.op_type " +
		" FROM outgo_order oor " +
		" LEFT JOIN log_vaccount lv ON lv.biz_log_no = oor.log_no  " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.CustOutgoBillsDetailData{}
	if err == nil {
		for rows.Next() {
			var finishTime, paymentType sql.NullString
			err = rows.Scan(
				&data.LogNo,
				&data.Amount,
				&data.OrderStatus,
				&paymentType,
				&data.CreateTime,
				&finishTime,
				&data.Fees,
				&data.BalanceType,
				&data.OpType,
			)
			if finishTime.String != "" {
				data.FinishTime = finishTime.String
			}
			if paymentType.String != "" {
				data.PaymentType = paymentType.String
			}
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
			data.OrderType = constants.VaReason_OUTGO

		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return data, ss_err.ERR_SUCCESS
}

// 获取未统计手续费的订单信息
func (*OutgoOrderDao) WithdrawFeesTaskResult() (withdrawResults []*m.CommonFeesResult) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT  log_no,  balance_type,  fees  FROM outgo_order where is_count = '0' and order_status = '3' and  create_time >  current_timestamp+interval  '-1 hour'")
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[查询 ---------> %s 结果失败,err -----> %s]", "WithdrawFeesTaskResult", err.Error())
		return nil
	}
	for rows.Next() {
		data := &m.CommonFeesResult{}
		err := rows.Scan(
			&data.LogNo,
			&data.MoneyType,
			&data.Fees,
		)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.FeesType = constants.FEES_TYPE_WITHDRAW
		withdrawResults = append(withdrawResults, data)
	}
	return withdrawResults
}

func (OutgoOrderDao) GetVaccNoTx(tx *sql.Tx, logNo string) (string, string, string) {
	var vaccNo, amountT, feesT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select vaccount_no,amount,fees from outgo_order where log_no=$1 limit 1`,
		[]*sql.NullString{&vaccNo, &amountT, &feesT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", ""
	}
	return vaccNo.String, amountT.String, feesT.String
}

//查询取款订单的用户accNo（本接口用于推送消息时要查询用户账号uuid）
func (OutgoOrderDao) GetAccNoByLogNo(tx *sql.Tx, logNo string) (useAccountUid, amount, balanceType string) {
	var useAccountUidT, amountT, balanceTypeT sql.NullString
	sqlStr := "select vacc.account_no, oor.amount, oor.balance_type from outgo_order oor " +
		" left join vaccount vacc on vacc.vaccount_no = oor.vaccount_no and vacc.is_delete = '0' " +
		" where oor.log_no = $1 "

	errT := ss_sql.QueryRowTx(tx, sqlStr, []*sql.NullString{&useAccountUidT, &amountT, &balanceTypeT}, logNo)
	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return "", "", ""
	}

	return useAccountUidT.String, amountT.String, balanceTypeT.String
}

func (*OutgoOrderDao) QueryCreateTime(opAccNo, createTime string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var createTimeT sql.NullString
	sqlCnt := "select create_time from outgo_order where op_acc_no= $1 AND create_time > $2 and order_status = $3 order by create_time desc  limit 1"
	return createTimeT.String, ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&createTimeT}, opAccNo, createTime, constants.OrderStatus_Pending)
}
