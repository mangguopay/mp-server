package dao

import (
	"context"
	"database/sql"
	"errors"

	"a.a/mp-server/bill-srv/m"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	tmProto "a.a/mp-server/common/proto/tm"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/common/ss_struct"
)

type IncomeOrderDao struct {
	LogNo       string
	ActAccNo    string
	Amount      string
	ServicerNo  string
	OrderStatus string
	BalanceType string
	Fees        string
	RecvAccNo   string
	RecvVAccNo  string
	PaymentType string
	ReeRate     string
	RealAmount  string
	OpAccNo     string
	FinishTime  string
	OpAccType   int
}

var IncomeOrderDaoInst IncomeOrderDao

func (IncomeOrderDao) InsertIncomeOrder(tx *sql.Tx, recvAccNo, recvVaccNo, amount, actAccNo, servicerNo, fees, balanceType, paymentType, reeRate, realAmount, opAccNo string, opAccType int) (logNo string) {
	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into income_order(log_no,act_acc_no,amount,servicer_no,create_time,order_status,`+
		`balance_type,fees,recv_acc_no,recv_vacc,is_count,payment_type,ree_rate,real_amount,op_acc_no,op_acc_type) values($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,0,$10,$11,$12,$13,$14)`,
		logNo, actAccNo, amount, servicerNo, constants.OrderStatus_Pending, balanceType, fees, recvAccNo, recvVaccNo, paymentType, reeRate, realAmount, opAccNo, opAccType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}

	return logNo
}

func (IncomeOrderDao) InsertIncomeOrderV3(data *IncomeOrderDao) (logNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	logNo = strext.GetDailyId()
	err := ss_sql.Exec(dbHandler, `insert into income_order(log_no,act_acc_no,amount,servicer_no,create_time,order_status,`+
		`balance_type,fees,recv_acc_no,recv_vacc,is_count,payment_type,ree_rate,real_amount,op_acc_no,op_acc_type) values($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,0,$10,$11,$12,$13,$14)`,
		logNo, data.ActAccNo, data.Amount, data.ServicerNo, constants.OrderStatus_Pending, data.BalanceType, data.Fees,
		data.RecvAccNo, data.RecvVAccNo, data.PaymentType, data.ReeRate, data.RealAmount, data.OpAccNo, data.OpAccType)
	if err != nil {
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
func (IncomeOrderDao) UpdateIncomeOrderOrderStatusV2(tmProxy *ss_struct.TmServerProxy, logNo, orderStatus string) error {
	sql := `update income_order set order_status=$1,finish_time=current_timestamp where log_no=$2`
	args := []string{orderStatus, logNo}
	rsp, err := tmProxy.GetTmServer().TxExec(context.TODO(), &tmProto.TxExecRequest{FromServerId: tmProxy.GetFromServerId(), TxNo: tmProxy.GetTxNo(), Sql: sql, Args: args})
	if err != nil {
		return err
	}
	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}
	return nil
}

func (*IncomeOrderDao) QueryAmount(logNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// todo 把orderStatus改成常量传入
	var amount sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select real_amount from income_order where log_no=$1 and order_status = '3' limit 1`,
		[]*sql.NullString{&amount}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return amount.String
}

func (*IncomeOrderDao) QueryIncomeOrder(logNo, status string) (*IncomeOrderDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var amountT, serviceNoT, finishTime, saveAccountT, recAccountT, feesT, balanceTypeT, statusT, realAmount sql.NullString
	if status != "" {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,finish_time,act_acc_no,recv_acc_no,fees,balance_type,real_amount
			from income_order where log_no=$1 and order_status = $2 limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &saveAccountT, &recAccountT, &feesT, &balanceTypeT, &realAmount}, logNo, status)
		if nil != err {
			ss_log.Error("err=%v", err)
			return nil, err
		}
	} else {
		err := ss_sql.QueryRow(dbHandler, `select amount,servicer_no,finish_time,act_acc_no,recv_acc_no,fees,balance_type,order_status,real_amount
			from income_order where log_no=$1  limit 1`,
			[]*sql.NullString{&amountT, &serviceNoT, &finishTime, &saveAccountT, &recAccountT, &feesT, &balanceTypeT, &statusT, &realAmount}, logNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return nil, err
		}
	}
	data := new(IncomeOrderDao)
	data.Amount = amountT.String
	data.ServicerNo = serviceNoT.String
	data.FinishTime = finishTime.String
	data.ActAccNo = saveAccountT.String
	data.RecvAccNo = recAccountT.String
	data.Fees = feesT.String
	data.BalanceType = balanceTypeT.String
	data.RealAmount = realAmount.String
	data.OrderStatus = statusT.String

	return data, nil
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

// 获取未统计手续费的订单信息
func (*IncomeOrderDao) SaveMoneyFeesTaskResult() (saveMoneyResults []*m.CommonFeesResult) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT  log_no,  balance_type,  fees  FROM income_order where is_count = '0' and  create_time >  current_timestamp+interval  '-1 hour'")
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[查询 ---------> %s 结果失败,err -----> %s]", "SaveMoneyFeesTaskResult", err.Error())
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
		data.FeesType = constants.FEES_TYPE_SAVEMONEY
		saveMoneyResults = append(saveMoneyResults, data)
	}
	return saveMoneyResults
}

func (*IncomeOrderDao) QueryCreateTime(opAccNo, createTime string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var createTimeT sql.NullString
	sqlCnt := "select create_time from income_order where op_acc_no= $1 AND create_time > $2 and order_status = $3 order by create_time desc  limit 1"
	return createTimeT.String, ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&createTimeT}, opAccNo, createTime, constants.OrderStatus_Pending)
}
