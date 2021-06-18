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

type OutgoOrderDao struct {
}

var OutgoOrderDaoInst OutgoOrderDao

func (*OutgoOrderDao) CustOutgoBillsDetail(logNo string) (data *go_micro_srv_cust.CustOutgoBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "oor.log_no", Val: logNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT oor.log_no, oor.amount, oor.order_status, oor.payment_type, oor.create_time, oor.modify_time, oor.fees, oor.balance_type " +
		" FROM outgo_order oor  " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_cust.CustOutgoBillsDetailData{}
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

func (*OutgoOrderDao) GetOutgoBillsCount(dbHandler *sql.DB, whereModelStr string, whereModelAgrs []interface{}) (datas []*go_micro_srv_cust.GetServicerOrderCountData, err string) {
	sqlStr := "select sum(amount), count(1) servicer_no, balance_type from outgo_order " + whereModelStr
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereModelAgrs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	var datasT []*go_micro_srv_cust.GetServicerOrderCountData
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return nil, ss_err.ERR_SYS_DB_GET
	} else {
		for rows.Next() {
			var data go_micro_srv_cust.GetServicerOrderCountData
			errT = rows.Scan(
				&data.OutgoAmountSum,
				&data.OutgoTotalSum,
				&data.ServicerNo,
				&data.BalanceType,
			)
			if errT != nil {
				ss_log.Error("err=[%v]", errT)
				return nil, ss_err.ERR_SYS_DB_GET
			}
			data.IncomeAmountSum = "0"
			data.IncomeTotalSum = "0"

			datasT = append(datasT, &data)
		}
	}

	return datasT, ss_err.ERR_SUCCESS
}

// GetToCustInfoCount 统计(核销码取款)
func (*OutgoOrderDao) GetToCustInfoCountWriteOffByDate(beforeDay, currencyType string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := `select count(1) as totalCount,sum(amount) as totalAmount,sum(fees) as totalFees  from outgo_order 
WHERE create_time >= $1 AND create_time < $2 and balance_type = $3 and withdraw_type = 0 and order_status = 3  `

	//====================================
	var num, amount, fee sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount, &fee}, startTime, endTime, currencyType)

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}
	data := &DataCount{
		Num:    strext.ToInt64(num.String),
		Amount: strext.ToInt64(amount.String),
		Fee:    strext.ToInt64(fee.String),
	}

	data.Type = 2 // 银行卡提现
	data.Day = beforeDay
	data.CType = currencyType
	return data, nil
}

// GetToCustInfoCount 统计(核销码取款)
func (*OutgoOrderDao) GetToCustInfoCountSweepByDate(beforeDay, currencyType string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := `select count(1) as totalCount,sum(amount) as totalAmount,sum(fees) as totalFees  from outgo_order 
WHERE create_time >= $1 AND create_time < $2 and balance_type = $3 and withdraw_type in (1,2) and order_status = 3  `

	var num, amount, fee sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount, &fee}, startTime, endTime, currencyType)

	if errT != nil {
		return nil, errT
	}
	data := &DataCount{
		Num:    strext.ToInt64(num.String),
		Amount: strext.ToInt64(amount.String),
		Fee:    strext.ToInt64(fee.String),
		Type:   3, // 银行卡提现
		Day:    beforeDay,
		CType:  currencyType,
	}
	return data, nil
}

//获取提现订单的手续费
func (OutgoOrderDao) GetOutgoOrderFees(orderNo string) (fess string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var feesT sql.NullString
	sqlStr := "select fees from outgo_order where log_no = $1 "
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&feesT}, orderNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "0"
	}
	return feesT.String
}
