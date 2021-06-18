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

type TransferDao struct {
	FromVacc     string
	ToVacc       string
	Amount       string
	ExchangeType string
	Fees         string
	MoneyType    string
	FeeRate      string
	RealAmount   string
	Lat          string
	Lng          string
	Ip           string
}

var TransferDaoInst TransferDao

func (TransferDao) InsertTransfer(tx *sql.Tx, log *TransferDao) (string, error) {
	sqlStr := `insert into transfer_order(log_no, from_vaccount_no, to_vaccount_no, amount, create_time, order_status, 
		exchange_type, fees, is_count, balance_type, ree_rate, real_amount, lat, lng, ip)
		values($1, $2, $3, $4, current_timestamp, $5, $6, $7, 0, $8, $9, $10, $11, $12, $13)`

	logNo := strext.GetDailyId()
	err := ss_sql.ExecTx(tx, sqlStr, logNo, log.FromVacc, log.ToVacc, log.Amount, constants.OrderStatus_Pending,
		log.ExchangeType, log.Fees, log.MoneyType, log.FeeRate, log.RealAmount, log.Lat, log.Lng, log.Ip)

	return logNo, err
}

func (TransferDao) UpdateTransferOrderStatus(tx *sql.Tx, logNo, orderStatus string) string {
	err := ss_sql.ExecTx(tx, `update transfer_order set order_status=$1,finish_time=current_timestamp,modify_time=current_timestamp where log_no=$2`,
		orderStatus, logNo)
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
		{Key: "log_no", Val: logNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT tro.log_no, tro.amount, tro.order_status, tro.payment_type, tro.create_time, tro.finish_time, tro.fees, tro.balance_type " +
		", acc.phone, acc.country_code, acc2.phone, acc2.country_code " +
		" FROM transfer_order tro " +
		" LEFT JOIN vaccount vacc ON vacc.vaccount_no = tro.to_vaccount_no" +
		" LEFT JOIN account acc ON acc.uid= vacc.account_no " +
		" LEFT JOIN vaccount vacc2 ON vacc2.vaccount_no = tro.from_vaccount_no" +
		" LEFT JOIN account acc2 ON acc2.uid= vacc2.account_no " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.CustTransferBillsDetailData{}
	if err == nil {
		for rows.Next() {
			var finishTime, toPhoneCountryCode, fromPhoneCountryCode sql.NullString
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
				&toPhoneCountryCode,
				&data.FromPhone,
				&fromPhoneCountryCode,
			)

			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
			data.FinishTime = finishTime.String
			data.ToPhoneCountryCode = toPhoneCountryCode.String
			data.FromPhoneCountryCode = fromPhoneCountryCode.String
			data.OrderType = constants.VaReason_TRANSFER
		}
	} else {
		ss_log.Error("err=[%v]", err)
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
					data.OpType = constants.VaOpType_Minus
				} else if toVaccountNo.String == usdVaccountNo.String {
					data.OpType = constants.VaOpType_Add
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

// 获取未统计手续费的订单信息
func (*TransferDao) TransferFeesTaskResult() (transferResults []*m.CommonFeesResult) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT  log_no,  balance_type,  fees  FROM transfer_order where is_count = '0' and  create_time >  current_timestamp+interval  '-1 hour'")
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[查询 ---------> %s 结果失败,err -----> %s]", "TransferFeesTaskResult", err.Error())
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
		data.FeesType = constants.FEES_TYPE_TRANSFER
		transferResults = append(transferResults, data)
	}
	return transferResults
}

func (*TransferDao) QueryAmount(logNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var amount sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select real_amount from transfer_order where log_no=$1 and order_status = '3' limit 1`,
		[]*sql.NullString{&amount}, logNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return amount.String
}
