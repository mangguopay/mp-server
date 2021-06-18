package dao

import (
	"a.a/cu/strext"
	"database/sql"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type BusinessTransferOrderDao struct {
	LogNo           string
	Amount          string
	CurrencyType    string
	Fee             string
	RealAmount      string
	OrderStatus     string
	CreateTime      string
	FinishTime      string
	WrongReason     string
	OutTransferNo   string
	AppId           string
	NotifyUrl       string
	NotifyStatus    string
	NotifyFailTimes int
}

var BusinessTransferOrderDaoInst BusinessTransferOrderDao

func (BusinessTransferOrderDao) GetTransferOrderByLogNo(logNo string) (*BusinessTransferOrderDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select bt.out_transfer_no, bt.amount, bt.rate, bt.fee, bt.currency_type, bt.order_status, bt.wrong_reason, et.app_id, " +
		"bt.finish_time, et.notify_url, et.notify_status, et.notify_fail_times " +
		"from business_transfer_order bt " +
		"left join enterprise_transfer_order et on et.transfer_log_no = bt.log_no " +
		"where log_no = $1 "

	var outOrderNo, amount, rate, fee, currencyType, orderStatus, appId, finishTime, wrongReason, notifyUrl, notifyStatus,
		notifyFailTimes sql.NullString

	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&outOrderNo, &amount, &rate, &fee, &currencyType,
		&orderStatus, &wrongReason, &appId, &finishTime, &notifyUrl, &notifyStatus, &notifyFailTimes}, logNo)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessTransferOrderDao)
	obj.LogNo = logNo
	obj.OutTransferNo = outOrderNo.String
	obj.Amount = amount.String
	obj.Fee = fee.String
	obj.CurrencyType = currencyType.String
	obj.OrderStatus = orderStatus.String
	obj.AppId = appId.String
	obj.FinishTime = finishTime.String
	obj.WrongReason = wrongReason.String
	obj.NotifyUrl = notifyUrl.String
	obj.NotifyStatus = notifyStatus.String
	obj.NotifyFailTimes = strext.ToInt(notifyFailTimes.String)
	return obj, nil
}

//修改订单通知次数, 通知失败次数(失败times=1, 成功times=0)
func (BusinessTransferOrderDao) UpdateNotifyStatusByOrderNo(req *UpdateTransferNotifyStatus) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var err error
	if req.NextTime == "" { //通知成功或超时处理
		sqlStr := "UPDATE enterprise_transfer_order SET notify_status=$1,notify_fail_times=notify_fail_times+$2 WHERE transfer_log_no=$3 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.OrderNo)

	} else { //通知失败处理
		sqlStr := "UPDATE enterprise_transfer_order SET notify_status=$1,notify_fail_times=notify_fail_times+$2,next_notify_time=$3 WHERE transfer_log_no=$4 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.NextTime, req.OrderNo)
	}
	return err
}

//查询遗漏通知的订单
func (BusinessTransferOrderDao) QueryNotifyOmission(orderStatus, notifyStatus, finishTime string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT log_no " +
		"FROM business_transfer_order bt " +
		"LEFT JOIN enterprise_transfer_order et ON et.transfer_log_no = bt.log_no " +
		"WHERE bt.transfer_type= 2 AND bt.order_status=$1 AND et.notify_status=$2 AND et.finish_time < $3"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, orderStatus, notifyStatus, finishTime)
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
func (BusinessTransferOrderDao) QueryNotifyBreak(orderStatus, notifyStatus, nextNotifyTime string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT log_no " +
		"FROM business_transfer_order bt " +
		"LEFT JOIN enterprise_transfer_order et ON et.transfer_log_no = bt.log_no " +
		"WHERE bt.transfer_type= 2 AND bt.order_status=$1 AND et.notify_status=$2 AND et.next_notify_time < $3"
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
