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

type IncomeOrderDao struct {
}

var IncomeOrderDaoInst IncomeOrderDao

func (*IncomeOrderDao) CustIncomeBillsDetail(logNo string) (data *go_micro_srv_cust.CustIncomeBillsDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "ior.log_no", Val: logNo, EqType: "="},
	})

	where2 := whereModel.WhereStr
	args2 := whereModel.Args
	sqlStr := "SELECT ior.log_no, ior.amount, ior.order_status, ior.payment_type, ior.create_time, ior.finish_time, ior.balance_type,ior.fees " +
		" FROM income_order ior " + where2
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args2...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_cust.CustIncomeBillsDetailData{}
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

//获取服务商存款统计
func (*IncomeOrderDao) GetIncomeBillsCount(dbHandler *sql.DB, whereModelStr string, whereModelAgrs []interface{}) (datas []*go_micro_srv_cust.GetServicerOrderCountData, err string) {
	sqlStr := "select sum(amount), count(1) servicer_no, balance_type from income_order " + whereModelStr
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
			data := go_micro_srv_cust.GetServicerOrderCountData{}
			errT = rows.Scan(
				&data.IncomeAmountSum,
				&data.IncomeTotalSum,
				&data.ServicerNo,
				&data.BalanceType,
			)
			if errT != nil {
				ss_log.Error("err=[%v]", errT)
				return nil, ss_err.ERR_SYS_DB_GET
			}

			data.OutgoAmountSum = "0"
			data.OutgoTotalSum = "0"
			data.ProfitAmountSum = "0"

			datasT = append(datasT, &data)
		}
	}

	return datasT, ss_err.ERR_SUCCESS
}

// GetIncomeOrderCountByDate 统计
func (*IncomeOrderDao) GetIncomeOrderCountByDate(beforeDay, currencyType string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := `select count(1) tatalCount,sum(amount) as totalAmount,sum(fees) as totalFees FROM income_order WHERE 
	create_time >= $1 and create_time < $2 and balance_type = $3 and order_status = 3`
	var num, amount, fee sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount, &fee}, startTime, endTime, currencyType)
	if errT != nil {
		return nil, errT
	}
	return &DataCount{
		Num:    strext.ToInt64(num.String),
		Amount: strext.ToInt64(amount.String),
		Fee:    strext.ToInt64(fee.String),
		Type:   2, // 向服务商充值
		Day:    startTime,
		CType:  currencyType,
	}, nil
}
