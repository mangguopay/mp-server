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

type IncomeOrderDao struct {
}

var IncomeOrderDaoInst IncomeOrderDao

func (IncomeOrderDao) InsertIncomeOrder(tx *sql.Tx, recvAccNo, recvVaccNo, amount, actAccNo, servicerNo, code, fees, balanceType string) (logNo string) {
	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into income_order(log_no,act_acc_no,amount,servicer_no,create_time,order_status,`+
		`balance_type,fees,recv_acc_no,recv_vacc,is_count) values($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,0)`,
		logNo, actAccNo, amount, servicerNo, constants.OrderStatus_Pending, balanceType, fees, recvAccNo, recvVaccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return logNo
}

func (IncomeOrderDao) UpdateIncomeOrderOrderStatus(tx *sql.Tx, logNo, orderStatus string) (errCode string) {
	err := ss_sql.ExecTx(tx, `update income_order set order_status=$1,finish_time=current_timestamp where log_no=$2`,
		orderStatus, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (IncomeOrderDao) UpdateIsCountFromLogNo(tx *sql.Tx, logNo string, countStatus int) string {
	err := ss_sql.ExecTx(tx, `update income_order set is_count=$2,modify_time=current_timestamp where log_no=$1`, logNo, countStatus)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (*IncomeOrderDao) ConfirmIsNoCount(tx *sql.Tx, logNo string) (phone string) {

	var logNoT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select log_no from income_order where log_no=$1 and is_count='1' or is_count='2' limit 1`, []*sql.NullString{&logNoT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT.String
}

func (*IncomeOrderDao) QueryAmount(logNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// todo 把orderStatus改成常量传入
	var amount sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select amount from income_order where log_no=$1 and order_status = '3' limit 1`,
		[]*sql.NullString{&amount}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return amount.String
}
func (*IncomeOrderDao) QuerySrvAccNoFromLogNo(logNo string) (string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var srvAccNo, servicerNoT, opAccNoT, opAccTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `SELECT rai.account_no,io.servicer_no,io.op_acc_no,io.op_acc_type FROM income_order io 
			LEFT JOIN rela_acc_iden rai ON io.servicer_no = rai.iden_no  WHERE io.log_no = $1 AND rai.account_type = '3'`,
		[]*sql.NullString{&srvAccNo, &servicerNoT, &opAccNoT, &opAccTypeT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", ""
	}
	return srvAccNo.String, servicerNoT.String, opAccNoT.String, opAccTypeT.String
}

func (*IncomeOrderDao) QueryIncomeOrder(logNo, status string) (string, string, string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var amountT, serviceNoT, finishTime, saveAccountT, recAccountT, feesT, balanceTypeT, statusT sql.NullString
	if status != "" {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,finish_time,act_acc_no,recv_acc_no,fees,balance_type
			from income_order where log_no=$1 and order_status = $2 limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &saveAccountT, &recAccountT, &feesT, &balanceTypeT}, logNo, status)
		if nil != err {
			ss_log.Error("err=%v", err)
			return "", "", "", "", "", "", "", ""
		}
	} else {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,finish_time,act_acc_no,recv_acc_no,fees,balance_type,order_status
			from income_order where log_no=$1  limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &saveAccountT, &recAccountT, &feesT, &balanceTypeT, &statusT}, logNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return "", "", "", "", "", "", "", ""
		}
	}
	return amountT.String, serviceNoT.String, finishTime.String, saveAccountT.String, recAccountT.String, feesT.String, balanceTypeT.String, statusT.String
}

func (*IncomeOrderDao) CustIncomeBillsDetail(logNo string) (data *go_micro_srv_bill.CustIncomeBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ior.log_no", Val: logNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT ior.log_no, ior.amount, ior.order_status, ior.payment_type, ior.create_time, ior.finish_time, ior.balance_type,ior.fees " +
		",lv.op_type " +
		" FROM income_order ior " +
		" LEFT JOIN log_vaccount lv ON ior.log_no = lv.biz_log_no  " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.CustIncomeBillsDetailData{}
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
				&data.BalanceType,
				&data.Fees,
				&data.OpType,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
			if finishTime.String != "" {
				data.FinishTime = finishTime.String
			}
			if paymentType.String != "" {
				data.PaymentType = paymentType.String
			}
			data.OrderType = constants.VaReason_INCOME
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return data, ss_err.ERR_SUCCESS
}
