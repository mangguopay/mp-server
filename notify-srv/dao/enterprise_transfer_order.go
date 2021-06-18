package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type EnterpriseTransferOrder struct {
	TransferNo        string
	OutTransferNo     string
	Amount            int64
	RealAmount        int64
	Rate              int64
	Fee               int64
	CurrencyType      string
	OrderStatus       string
	PaymentType       string
	BusinessAccountNo string
	AppId             string
	PayeeAccountNo    string
	PayeeAccount      string
	NotifyUrl         string
	NotifyStatus      string
	NotifyFailTimes   int
	NextNotifyTime    string
	Remark            string
	WrongReason       string
	CreateTime        string
	FinishTime        string
}

var EnterpriseTransferOrderDao EnterpriseTransferOrder

func (EnterpriseTransferOrder) GetOrderByOrderNo(orderNo string) (*EnterpriseTransferOrder, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select transfer_no, out_transfer_no, amount, rate, fee, currency_type, order_status, app_id, payee_account," +
		"finish_time, wrong_reason, notify_url, notify_status, notify_fail_times " +
		"from enterprise_transfer_order  " +
		"where transfer_no = $1 "

	var transferOrderNo, outOrderNo, amount, rate, fee, currencyType, orderStatus, appId, payeeAccount,
		finishTime, wrongReason, notifyUrl, notifyStatus, notifyFailTimes sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{
		&transferOrderNo, &outOrderNo, &amount, &rate, &fee, &currencyType, &orderStatus, &appId, &payeeAccount,
		&finishTime, &wrongReason, &notifyUrl, &notifyStatus, &notifyFailTimes,
	}, orderNo)
	if err != nil {
		return nil, err
	}

	obj := new(EnterpriseTransferOrder)
	obj.TransferNo = transferOrderNo.String
	obj.OutTransferNo = outOrderNo.String
	obj.Amount = strext.ToInt64(amount.String)
	obj.Rate = strext.ToInt64(rate.String)
	obj.Fee = strext.ToInt64(fee.String)
	obj.CurrencyType = currencyType.String
	obj.OrderStatus = orderStatus.String
	obj.AppId = appId.String
	obj.PayeeAccount = payeeAccount.String
	obj.FinishTime = finishTime.String
	obj.WrongReason = wrongReason.String
	obj.NotifyUrl = notifyUrl.String
	obj.NotifyStatus = notifyStatus.String
	obj.NotifyFailTimes = strext.ToInt(notifyFailTimes.String)
	return obj, nil
}

type UpdateTransferNotifyStatus struct {
	OrderNo         string
	NextTime        string
	NotifyStatus    string
	NotifyFailTimes int32
}

//修改订单通知次数, 通知失败次数(失败times=1, 成功times=0)
func (EnterpriseTransferOrder) UpdateNotifyStatusByOrderNo(req *UpdateTransferNotifyStatus) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var err error
	if req.NextTime == "" { //通知成功或超时处理
		sqlStr := "UPDATE enterprise_transfer_order SET notify_status=$1,notify_fail_times=notify_fail_times+$2 WHERE transfer_no=$3 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.OrderNo)

	} else { //通知失败处理
		sqlStr := "UPDATE enterprise_transfer_order SET notify_status=$1,notify_fail_times=notify_fail_times+$2,next_notify_time=$3 WHERE transfer_no=$4 "
		err = ss_sql.Exec(dbHandler, sqlStr, req.NotifyStatus, req.NotifyFailTimes, req.NextTime, req.OrderNo)
	}
	return err
}
