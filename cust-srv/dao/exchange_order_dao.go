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

type ExchangeOrderDao struct {
}

var ExchangeOrderDaoInst ExchangeOrderDao

func (ExchangeOrderDao) CustExchangeBillsDetail(logNo string) (data *go_micro_srv_cust.ExchangeOrderData, errR string) {
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
	data = &go_micro_srv_cust.ExchangeOrderData{}
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

// GetExchangeCountByDate 统计
func (*ExchangeOrderDao) GetExchangeCountByDate(beforeDay string) (*DataCount, error) {
	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	data := &DataCount{}
	data.Day = startTime
	totalCount, totalAmount, totalFee, tErr := usd2khrCount(startTime, endTime)
	if tErr != nil && tErr.Error() != ss_sql.DB_NO_ROWS_MSG {
		return nil, tErr
	}
	data.Usd2khrNum = totalCount
	data.Usd2khrAmount = totalAmount
	data.Usd2khrFee = totalFee

	khr2usdCount, khr2usdTotalAmount, khr2usdTotalFee, khr2usdErr := khr2usdCount(startTime, endTime)
	if khr2usdErr != nil && khr2usdErr.Error() != ss_sql.DB_NO_ROWS_MSG {
		return nil, khr2usdErr
	}
	data.Khr2usdNum = khr2usdCount
	data.Khr2usdAmount = khr2usdTotalAmount
	data.Khr2usdFee = khr2usdTotalFee
	return data, nil
}

func usd2khrCount(startTime, endTime string) (int64, int64, int64, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var totalCount, totalAmount, totalFee sql.NullString

	sqlStr := `select count(1) as totalCount,sum(amount) as totalAmount,sum(fees) as totalFee from exchange_order 
	WHERE create_time >= $1 and create_time  < $2 and in_type = $3 and out_type = $4 and order_status = 3`

	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalCount, &totalAmount, &totalFee}, startTime, endTime,
		constants.CURRENCY_USD, constants.CURRENCY_KHR)

	if err != nil {
		return strext.ToInt64(totalCount.String), strext.ToInt64(totalAmount.String), strext.ToInt64(totalFee.String), err
	}
	return strext.ToInt64(totalCount.String), strext.ToInt64(totalAmount.String), strext.ToInt64(totalFee.String), nil
}
func khr2usdCount(startTime, endTime string) (int64, int64, int64, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var totalCount, totalAmount, totalFee sql.NullString

	sqlStr := `select count(1) as totalCount,sum(amount) as totalAmount,sum(fees) as totalFee from exchange_order 
	WHERE create_time >= $1 and create_time  < $2 and in_type = $3 and out_type = $4 and order_status = 3`

	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalCount, &totalAmount, &totalFee}, startTime, endTime,
		constants.CURRENCY_KHR, constants.CURRENCY_USD)

	if err != nil {
		return strext.ToInt64(totalCount.String), strext.ToInt64(totalAmount.String), strext.ToInt64(totalFee.String), err
	}
	return strext.ToInt64(totalCount.String), strext.ToInt64(totalAmount.String), strext.ToInt64(totalFee.String), nil
}
