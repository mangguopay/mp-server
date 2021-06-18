package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessBillDao struct {
}

var BusinessBillDaoInst BusinessBillDao

type QueryOrderInfo struct {
	OutOrderNo      string
	OrderNo         string
	CustAccount     string
	OrderStatus     string
	Amount          string
	CurrencyType    string
	CreateTime      string
	PayTime         string
	AppId           string
	BusinessAccount string
	Subject         string
	Remark          string
	Rate            string
	Fee             string
	NotifyUrl       string
	NotifyStatus    string
	NotifyFailTimes int
	TradeType       string
}

//查询订单信息
func (BusinessBillDao) QueryOrderInfoByOrderNo(innerOrderNo string) (*QueryOrderInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var (
		outOrderNo, orderNo, orderStatus, amount, currencyType, payTime, appId, subject, notifyStatus, notifyUrl,
		notifyFailTimes, tradeType, userAcc, businessAcc sql.NullString
	)
	sqlStr := "SELECT b.out_order_no, b.order_no, b.order_status, b.amount, b.currency_type, b.pay_time, b.app_id," +
		"b.subject, b.notify_status, b.notify_url, b.notify_fail_times, b.trade_type, " +
		"acc1.account AS user_acc, acc2.account AS business_acc " +
		"FROM business_bill b " +
		"LEFT JOIN account acc1 ON acc1.uid=b.account_no " +
		"LEFT JOIN account acc2 ON acc2.uid=b.business_account_no " +
		"WHERE b.order_no=$1 LIMIT 1"
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{
		&outOrderNo, &orderNo, &orderStatus, &amount, &currencyType, &payTime, &appId, &subject, &notifyStatus, &notifyUrl,
		&notifyFailTimes, &tradeType, &userAcc, &businessAcc}, innerOrderNo)

	if err != nil {
		return nil, err
	}
	orderInfo := &QueryOrderInfo{
		OutOrderNo:      outOrderNo.String,
		OrderNo:         orderNo.String,
		OrderStatus:     orderStatus.String,
		Amount:          amount.String,
		CurrencyType:    currencyType.String,
		PayTime:         ss_time.ParseTimeFromPostgres(payTime.String, global.Tz).Format(ss_time.DateTimeDashFormat),
		AppId:           appId.String,
		Subject:         subject.String,
		NotifyUrl:       notifyUrl.String,
		NotifyStatus:    notifyStatus.String,
		NotifyFailTimes: strext.ToInt(notifyFailTimes.String),
		CustAccount:     userAcc.String,
		BusinessAccount: businessAcc.String,
		TradeType:       tradeType.String,
	}

	return orderInfo, nil
}

type UpdateBillNotifyStatus struct {
	OrderNo         string
	NextTime        string
	NotifyStatus    string
	NotifyFailTimes int32
}

//修改订单通知次数, 通知失败次数(失败times=1, 成功times=0)
func (BusinessBillDao) UpdateNotifyStatusByOrderNo(req *UpdateBillNotifyStatus) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var err error
	if req.NextTime == "" { //通知成功或超时处理
		sqlStr := "UPDATE business_bill SET notify_status=$1,notify_fail_times=notify_fail_times+$2 WHERE order_no=$3 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.OrderNo)

	} else { //通知失败处理
		sqlStr := "UPDATE business_bill SET notify_status=$1,notify_fail_times=notify_fail_times+$2,next_notify_time=$3 WHERE order_no=$4 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.NextTime, req.OrderNo)
	}
	return err
}

//查询遗漏通知的订单
func (BusinessBillDao) QueryNotifyOmission(orderStatus, notifyStatus, payTime string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT order_no FROM business_bill WHERE order_status=$1 AND notify_status=$2 AND pay_time < $3"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, orderStatus, notifyStatus, payTime)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var orderNoList []string
	for rows.Next() {
		var orderNo sql.NullString
		err := rows.Scan(&orderNo)
		if err != nil {
			return nil, err
		}
		orderNoList = append(orderNoList, orderNo.String)
	}
	return orderNoList, nil
}

//查询通知中断的订单
func (BusinessBillDao) QueryNotifyBreak(orderStatus, notifyStatus, nextNotifyTime string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT order_no FROM business_bill "
	sqlStr += "WHERE order_status=$1 AND notify_status=$2 AND next_notify_time < $3"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, orderStatus, notifyStatus, nextNotifyTime)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var orderNoList []string
	for rows.Next() {
		var orderNo sql.NullString
		err := rows.Scan(&orderNo)
		if err != nil {
			return nil, err
		}
		orderNoList = append(orderNoList, orderNo.String)
	}
	return orderNoList, nil
}
