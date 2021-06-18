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

type ExchangeOrderDao struct {
}

var ExchangeOrderDaoInst ExchangeOrderDao

func (ExchangeOrderDao) InsertExchangeOrder(tx *sql.Tx, accountNo, inType, outType, amount, rate, transFrom, transAmount, errReason, fees string) (id string) {
	id = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into exchange_order(log_no,in_type,out_type,amount,create_time,rate,order_status,account_no,trans_from,trans_amount,err_reason,fees,is_count) values($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8,$9,$10,$11,0)`,
		id, inType, outType, amount, rate, constants.OrderStatus_Init, accountNo, transFrom, transAmount, errReason, fees)
	if nil != err {
		ss_log.Error("err=%v", err)
		return id
	}
	return id
}

func (*ExchangeOrderDao) ConfirmIsNoCount(tx *sql.Tx, logNo string) (phone string) {

	var logNoT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select log_no from exchange_order where log_no=$1 and is_count='1' or is_count='2' limit 1`, []*sql.NullString{&logNoT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT.String
}

func (ExchangeOrderDao) UpdateExchangeOrderStatus(tx *sql.Tx, logNo, orderStatus, errReason string) string {
	err := ss_sql.ExecTx(tx, `update exchange_order set order_status=$2, err_reason=$3 where log_no=$1`,
		logNo, orderStatus, errReason)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}
func (ExchangeOrderDao) UpdateIsCountFromLogNo(tx *sql.Tx, logNo string, countStatus int) string {
	err := ss_sql.ExecTx(tx, `update exchange_order set is_count=$2,modify_time=current_timestamp where log_no=$1`, logNo, countStatus)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (ExchangeOrderDao) CustExchangeBillsDetail(logNo string) (data *go_micro_srv_bill.ExchangeOrderData, errR string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "log_no", Val: logNo, EqType: "="},
	})

	sqlStr := "select log_no, in_type, out_type, amount, create_time, rate, order_status, finish_time, account_no, trans_from" +
		" ,trans_amount, err_reason, fees " +
		" from exchange_order " + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.ExchangeOrderData{}
	if err == nil {
		for rows.Next() {
			var finishTime, errReason sql.NullString
			err = rows.Scan(
				&data.LogNo,
				&data.InType,
				&data.OutType,
				&data.Amount,
				&data.CreateTime,
				&data.Rate,
				&data.OrderStatus,
				&finishTime,
				&data.AccountNo,
				&data.TransFrom,
				&data.TransAmount,
				&errReason,
				&data.Fees,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
			if finishTime.String != "" {
				data.FinishTime = finishTime.String
			}
			if errReason.String != "" {
				data.ErrReason = errReason.String
			}

			data.OrderType = constants.VaReason_Exchange
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return data, ss_err.ERR_SUCCESS
}
