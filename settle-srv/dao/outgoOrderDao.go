package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type OutgoOrderDao struct {
}

var OutgoOrderDaoInst OutgoOrderDao

func (OutgoOrderDao) InsertOutgo(tx *sql.Tx, recvVaccNo, amount, servicerNo, opAccNo, moneyType, fees, rate string) (logNo string) {
	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into outgo_order(log_no,vaccount_no,amount,create_time,order_status,`+
		`balance_type,fees,servicer_no,op_acc_no,rate,is_count) values($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8,$9,0)`,
		logNo, recvVaccNo, amount, constants.OrderStatus_Pending, moneyType, fees, servicerNo, opAccNo, rate)
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

func (OutgoOrderDao) UpdateOutgoOrderStatus(tx *sql.Tx, logNo, orderStatus string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update outgo_order set order_status=$1,modify_time=current_timestamp where log_no=$2`,
		orderStatus, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PAY_OUT_MONEY
	}
	return ss_err.ERR_SUCCESS
}

func (*OutgoOrderDao) ConfirmIsNoCount(tx *sql.Tx, logNo string) (phone string) {

	var logNoT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select log_no from outgo_order where log_no=$1 and is_count='1' or is_count='2' limit 1`, []*sql.NullString{&logNoT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT.String
}
func (*OutgoOrderDao) QuerySrvAccNoFromLogNo(logNo string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var srvAccNo, servicerNoT, opAccNo, opAccType sql.NullString
	err := ss_sql.QueryRow(dbHandler, `SELECT rai.account_no,io.servicer_no,io.op_acc_no,io.op_acc_type FROM outgo_order io 
		LEFT JOIN rela_acc_iden rai ON io.servicer_no = rai.iden_no  WHERE io.log_no = $1 AND rai.account_type = '3'`,
		[]*sql.NullString{&srvAccNo, &servicerNoT, &opAccNo, &opAccType}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return srvAccNo.String, servicerNoT.String, opAccNo.String, opAccType.String
}
func (OutgoOrderDao) UpdateIsCountFromLogNo(tx *sql.Tx, logNo string, countStatus int) string {
	err := ss_sql.ExecTx(tx, `update outgo_order set is_count=$2,modify_time=current_timestamp where log_no=$1`, logNo, countStatus)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
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

func (OutgoOrderDao) QueryOutGoOrderFromLogNo(logNo, status string) (string, string, string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var amountT, serviceNoT, finishTime, vaccountNoT, feesT, balanceTypeT, withdrawTypeT, statusT sql.NullString
	if status != "" {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,modify_time,vaccount_no,fees,balance_type,withdraw_type from outgo_order where log_no=$1 and order_status = $2 limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &vaccountNoT, &feesT, &balanceTypeT, &withdrawTypeT}, logNo, status)
		if nil != err {
			ss_log.Error("err=%v", err)
			return "", "", "", "", "", "", "", ""
		}
	} else {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,finish_time,vaccount_no,fees,balance_type,withdraw_type,order_status from outgo_order where log_no=$1 limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &vaccountNoT, &feesT, &balanceTypeT, &withdrawTypeT, &statusT}, logNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return "", "", "", "", "", "", "", ""
		}
	}

	return amountT.String, serviceNoT.String, finishTime.String, vaccountNoT.String, feesT.String, balanceTypeT.String, withdrawTypeT.String, statusT.String
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
