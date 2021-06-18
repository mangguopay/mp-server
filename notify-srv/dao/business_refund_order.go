package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessRefundOrderDao struct {
	AppId           string
	OrderNo         string
	OutOrderNo      string
	TransAmount     string
	RefundNo        string
	OutRefundNo     string
	RefundAmount    string
	CurrencyType    string
	RefundStatus    string
	CreateTime      string
	FinishTime      string
	NotifyStatus    string
	NotifyUrl       string
	NotifyFailTimes int
	NexNotifyTime   string
}

var BusinessRefundOrderDaoInst BusinessRefundOrderDao

func (BusinessRefundOrderDao) GetRefundDetail(orderNo string) (*BusinessRefundOrderDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT br.refund_no, br.out_refund_no, br.amount AS refund_amount, br.refund_status, br.pay_order_no, " +
		"br.create_time, br.finish_time, br.notify_status, br.notify_url, br.notify_fail_times, br.next_notify_time, " +
		"bb.out_order_no, bb.amount AS trans_amount, bb.currency_type, bb.app_id " +
		"FROM business_refund_order br " +
		"LEFT JOIN business_bill bb ON bb.order_no = br.pay_order_no " +
		"WHERE br.refund_no=$1 "
	var refundNo, outRefundNo, refundAmount, refundStatus, payOrderNo, createTime, finishTime, notifyStatus, notifyUrl,
		notifyFailTimes, nextNotifyTime, outPayNo, transAmount, currencyType, appId sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&refundNo, &outRefundNo, &refundAmount, &refundStatus,
		&payOrderNo, &createTime, &finishTime, &notifyStatus, &notifyUrl, &notifyFailTimes, &nextNotifyTime, &outPayNo,
		&transAmount, &currencyType, &appId}, orderNo)
	if err != nil {
		return nil, err
	}

	order := new(BusinessRefundOrderDao)
	order.RefundNo = refundNo.String
	order.OutRefundNo = outRefundNo.String
	order.RefundAmount = refundAmount.String
	order.RefundStatus = refundStatus.String
	order.OrderNo = payOrderNo.String
	order.CreateTime = createTime.String
	order.FinishTime = finishTime.String
	order.NotifyStatus = notifyStatus.String
	order.NotifyUrl = notifyUrl.String
	order.NotifyFailTimes = strext.ToInt(notifyFailTimes.String)
	order.NexNotifyTime = nextNotifyTime.String
	order.OutOrderNo = outPayNo.String
	order.TransAmount = transAmount.String
	order.CurrencyType = currencyType.String
	order.AppId = appId.String
	return order, nil
}

type UpdateRefundNotifyStatus struct {
	OrderNo         string
	NextTime        string
	NotifyStatus    string
	NotifyFailTimes int32
}

//修改订单通知次数, 通知失败次数(失败times=1, 成功times=0)
func (BusinessRefundOrderDao) UpdateNotifyStatusByOrderNo(req *UpdateRefundNotifyStatus) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var err error
	if req.NextTime == "" { //通知成功或超时处理
		sqlStr := "UPDATE business_refund_order SET notify_status=$1,notify_fail_times=notify_fail_times+$2 WHERE refund_no=$3 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.OrderNo)

	} else { //通知失败处理
		sqlStr := "UPDATE business_refund_order SET notify_status=$1,notify_fail_times=notify_fail_times+$2,next_notify_time=$3 WHERE refund_no=$4 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.NextTime, req.OrderNo)
	}
	return err
}

//查询遗漏通知的订单
func (BusinessRefundOrderDao) QueryNotifyOmission(refundStatus, notifyStatus, finishTime string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT refund_no FROM business_refund_order WHERE refund_status=$1 AND notify_status=$2 AND finish_time < $3"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, refundStatus, notifyStatus, finishTime)
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
func (BusinessRefundOrderDao) QueryNotifyBreak(refundStatus, notifyStatus, nextNotifyTime string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT refund_no FROM business_refund_order WHERE refund_status=$1 AND notify_status=$2 AND next_notify_time < $3"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, refundStatus, notifyStatus, nextNotifyTime)
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
