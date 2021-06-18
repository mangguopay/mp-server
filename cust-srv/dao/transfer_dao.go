package dao

import (
	"database/sql"
	"time"

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

type TransferDao struct {
}

var TransferDaoInst TransferDao

//accountNo 有可能是转账发起的账号也有可能是转账收款的账号
func (*TransferDao) CustTransferBillsDetailByAccount(accountNO, logNo string) (data *go_micro_srv_cust.CustTransferBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//查询传来的账号对应的两个币种的转账虚拟账户
	var usdVaccountNo, khrVaccountNo sql.NullString
	sqlquery := "select vaccount_no from vaccount where account_no = $1 and  balance_type = $2 "
	errQueryUsdVacc := ss_sql.QueryRow(dbHandler, sqlquery, []*sql.NullString{&usdVaccountNo}, accountNO, "usd")
	if errQueryUsdVacc != nil {
		ss_log.Error("errQueryUsdVacc=[%v],查询账号的usd转账虚拟账户失败", errQueryUsdVacc)
		//return nil, ss_err.ERR_SYS_DB_GET
	}

	errQueryKhrVacc := ss_sql.QueryRow(dbHandler, sqlquery, []*sql.NullString{&khrVaccountNo}, accountNO, "khr")
	if errQueryKhrVacc != nil {
		ss_log.Error("errQueryKhrVacc=[%v],查询账号的khr转账虚拟账户失败", errQueryKhrVacc)
		//return nil, ss_err.ERR_SYS_DB_GET
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
	data = &go_micro_srv_cust.CustTransferBillsDetailData{}
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

// GetTransferCountByDate 统计
func (*TransferDao) GetTransferCountByDate(beforeDay, currencyType string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := `select count(1) as tatalCount, sum(amount) as totalAmount,sum(fees) as totalFee from transfer_order 
	WHERE create_time >= $1 and create_time  < $2 and balance_type = $3 and  order_status = 3 `
	var num, amount, fee sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount, &fee}, startTime, endTime, currencyType)
	if errT != nil {
		return nil, errT
	}
	return &DataCount{
		Num:    strext.ToInt64(num.String),
		Amount: strext.ToInt64(amount.String),
		Fee:    strext.ToInt64(fee.String),
		Day:    startTime,
		CType:  currencyType,
	}, nil
}
