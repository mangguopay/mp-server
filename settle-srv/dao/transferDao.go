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

type TransferDao struct {
}

var TransferDaoInst TransferDao

func (TransferDao) InsertTransfer(tx *sql.Tx, fromVacc, toVacc, amount, exchangeType, fees string) (logNo string) {

	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into transfer_order(log_no,from_vaccount_no,to_vaccount_no,amount,create_time,order_status,exchange_type,fees,is_count)
values($1,$2,$3,$4,current_timestamp,$5,$6,$7,0)`,
		logNo, fromVacc, toVacc, amount, constants.OrderStatus_Pending, exchangeType, fees)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return logNo
}

func (TransferDao) UpdateTransferOrderStatus(tx *sql.Tx, logNo, orderStatus string) string {
	err := ss_sql.ExecTx(tx, `update transfer_order set order_status=$1 where log_no=$2`,
		orderStatus, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (*TransferDao) ConfirmIsNoCount(tx *sql.Tx, logNo string) (phone string) {

	var logNoT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select log_no from transfer_order where log_no=$1 and is_count='1' or is_count='2' limit 1`, []*sql.NullString{&logNoT}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT.String
}

func (TransferDao) UpdateIsCountFromLogNo(tx *sql.Tx, logNo string, countStatus int) string {
	err := ss_sql.ExecTx(tx, `update transfer_order set is_count=$2,modify_time=current_timestamp where log_no=$1`, logNo, countStatus)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (*TransferDao) CustTransferBillsDetail(logNo string) (data *go_micro_srv_bill.CustTransferBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "log_no", Val: logNo, EqType: "like"},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT tro.log_no, tro.amount, tro.order_status, tro.payment_type, tro.create_time, tro.finish_time, tro.fees, tro.balance_type " +
		", acc.account " +
		" FROM transfer_order tro " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = tro.to_vaccount_no" +
		" LEFT JOIN account acc ON acc.uid= vacc.account_no " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.CustTransferBillsDetailData{}
	if err == nil {
		for rows.Next() {
			var finishTime sql.NullString
			err = rows.Scan(
				&data.LogNo,
				&data.Amount,
				&data.OrderStatus,
				&data.PaymentType,
				&data.CreateTime,
				&finishTime,
				&data.Fees,
				&data.BalanceType,

				&data.ToPhone,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
			if finishTime.String != "" {
				data.FinishTime = finishTime.String
			}

			data.OrderType = constants.VaReason_TRANSFER
		}
	} else {
		ss_log.Error("OutgoOrderDao | CustTransferBillsDetail | err=%v\n\nsql=[%v]", err, sqlStr)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return data, ss_err.ERR_SUCCESS
}

//accountNo 有可能是转账发起的账号也有可能是转账收款的账号
func (*TransferDao) CustTransferBillsDetailByAccount(accountNO, logNo string) (data *go_micro_srv_bill.CustTransferBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//查询传来的账号对应的两个币种的转账虚拟账户
	var usdVaccountNo, khrVaccountNo sql.NullString
	sqlquery := "select vaccount_no from vaccount where account_no = $1 and  balance_type = $2 "
	errQueryUsdVacc := ss_sql.QueryRow(dbHandler, sqlquery, []*sql.NullString{&usdVaccountNo}, accountNO, "usd")
	if errQueryUsdVacc != nil {
		ss_log.Error("errQueryUsdVacc=[%v],查询账号的usd转账虚拟账户失败", errQueryUsdVacc)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	errQueryKhrVacc := ss_sql.QueryRow(dbHandler, sqlquery, []*sql.NullString{&khrVaccountNo}, accountNO, "khr")
	if errQueryKhrVacc != nil {
		ss_log.Error("errQueryKhrVacc=[%v],查询账号的khr转账虚拟账户失败", errQueryKhrVacc)
		return nil, ss_err.ERR_SYS_DB_GET
	}

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "log_no", Val: logNo, EqType: "like"},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT tro.log_no, tro.amount, tro.order_status, tro.payment_type, tro.create_time, tro.finish_time, tro.fees, tro.balance_type, tro.from_vaccount_no, tro.to_vaccount_no " +
		", acc.account " +
		" FROM transfer_order tro " +
		//以下两条LEFT JOIN 是为了查询转账至谁的账号
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = tro.to_vaccount_no " +
		" LEFT JOIN account acc ON acc.uid= vacc.account_no " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.CustTransferBillsDetailData{}
	if err == nil {
		for rows.Next() {
			var finishTime, fromVaccountNo, toVaccountNo, paymentType sql.NullString
			err = rows.Scan(
				&data.LogNo,
				&data.Amount,
				&data.OrderStatus,
				&paymentType,
				&data.CreateTime,
				&finishTime,
				&data.Fees,
				&data.BalanceType,
				&fromVaccountNo,
				&toVaccountNo,
				&data.ToPhone, //转账至谁的账号
				//&data.OpType,
			)
			if finishTime.String != "" {
				data.FinishTime = finishTime.String
			}
			if paymentType.String != "" {
				data.PaymentType = paymentType.String
			}

			if data.BalanceType == "usd" {
				if fromVaccountNo.String == usdVaccountNo.String {
					data.OpType = constants.VaOpType_Add
				} else if toVaccountNo.String == usdVaccountNo.String {
					data.OpType = constants.VaOpType_Minus
				}
			} else if data.BalanceType == "khr" {
				if fromVaccountNo.String == khrVaccountNo.String {
					data.OpType = constants.VaOpType_Minus
				} else if toVaccountNo.String == khrVaccountNo.String {
					data.OpType = constants.VaOpType_Add
				}
			}

			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	return data, ss_err.ERR_SUCCESS
}
