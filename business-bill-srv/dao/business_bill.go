package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
)

var BusinessBillDaoInst BusinessBillDao

type BusinessBillDao struct {
	OrderNo            string
	OutOrderNo         string
	OrderStatus        string
	Rate               string
	Fee                string
	Amount             string
	RealAmount         string
	CurrencyType       string
	NotifyUrl          string
	RreturnUrl         string
	AppId              string
	AppName            string
	SimplifyName       string
	BusinessNo         string
	BusinessId         string
	BusinessVaccountNo string
	BusinessAccountNo  string
	AccountNo          string
	VAccountNo         string
	CreateTime         string
	PayTime            string
	ExpireTime         int64
	Remark             string
	Subject            string
	SceneNo            string
	SettleId           string
	Cycle              int
	SettleDate         int64

	TradeType         string
	BusinessChannelNo string

	// 附加字段
	PayAccount string // 付款人账号(account表中的account字段)
}

// 下单入库
func (b *BusinessBillDao) InsertOrderTx(tx *sql.Tx, order BusinessBillDao) error {
	insertSql, insertData, _, err := ss_sql.MkInsertSql("business_bill", map[string]string{
		"order_no":            order.OrderNo,
		"out_order_no":        order.OutOrderNo,
		"rate":                order.Rate,
		"fee":                 order.Fee,
		"amount":              order.Amount,
		"real_amount":         order.RealAmount,
		"currency_type":       order.CurrencyType,
		"order_status":        order.OrderStatus,
		"notify_url":          order.NotifyUrl,
		"return_url":          order.RreturnUrl,
		"remark":              order.Remark,
		"business_no":         order.BusinessNo,
		"app_id":              order.AppId,
		"subject":             order.Subject,
		"scene_no":            order.SceneNo,
		"business_account_no": order.BusinessAccountNo,
		"time_expire":         strext.ToString(order.ExpireTime),
		"account_no":          order.AccountNo,
		"trade_type":          order.TradeType,
		"business_channel_no": order.BusinessChannelNo,
		"create_time":         ss_sql.CURRENT_TIMESTAMP,
	})
	if err != nil {
		return err
	}

	return ss_sql.ExecTx(tx, insertSql, insertData...)
}
func (b *BusinessBillDao) InsertOrder(order BusinessBillDao) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr, param, _, err := ss_sql.MkInsertSql("business_bill", map[string]string{
		"order_no":            order.OrderNo,
		"fee":                 order.Fee,
		"amount":              order.Amount,
		"real_amount":         order.RealAmount,
		"order_status":        order.OrderStatus,
		"remark":              order.Remark,
		"notify_url":          order.NotifyUrl,
		"return_url":          order.RreturnUrl,
		"out_order_no":        order.OutOrderNo,
		"rate":                order.Rate,
		"business_no":         order.BusinessNo,
		"app_id":              order.AppId,
		"currency_type":       order.CurrencyType,
		"subject":             order.Subject,
		"scene_no":            order.SceneNo,
		"business_account_no": order.BusinessAccountNo,
		"time_expire":         strext.ToString(order.ExpireTime),
		"trade_type":          order.TradeType,
		"business_channel_no": order.BusinessChannelNo,
		"create_time":         ss_sql.CURRENT_TIMESTAMP,
	})
	if err != nil {
		return err
	}

	return ss_sql.Exec(dbHandler, sqlStr, param...)
}

type UpdateOrderPaidData struct {
	OrderNo            string
	AccountNo          string
	VaccountNo         string
	BusinessVaccountNo string
	PayTime            string
	Cycle              string
	SettleDate         int64
	PaymentMethod      string
}

// 修改支付订单为支付成功
func (b *BusinessBillDao) UpdateOrderPaid(tx *sql.Tx, data UpdateOrderPaidData) error {
	sqlStr := `update business_bill set account_no=$1, vaccount_no=$2, order_status=$3, business_vaccount_no=$4, pay_time=$5,
				cycle=$6, settle_date=$7, payment_method=$8 `
	sqlStr += `where order_no=$9 and order_status=$10`

	return ss_sql.ExecTx(tx, sqlStr,
		data.AccountNo, data.VaccountNo, constants.BusinessOrderStatusPay, data.BusinessVaccountNo, data.PayTime,
		data.Cycle, data.SettleDate, data.PaymentMethod,
		data.OrderNo, constants.BusinessOrderStatusPending,
	)
}

// 根据订单号修改支付订单为超时
func (b *BusinessBillDao) UpdateOrderOutTimeById(orderNo string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `update business_bill set order_status=$1 where order_no=$2 and order_status=$3`

	return ss_sql.Exec(dbHandler, sqlStr, constants.BusinessOrderStatusPayTimeOut, orderNo, constants.BusinessOrderStatusPending)
}

//修改所有支付超时的订单为超时
func (b *BusinessBillDao) UpdateOrderOutTime() error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `update business_bill set order_status=$1 where order_status=$2 AND time_expire <= $3 `

	return ss_sql.Exec(dbHandler, sqlStr, constants.BusinessOrderStatusPayTimeOut, constants.BusinessOrderStatusPending, ss_time.Now(global.Tz).Unix())
}

// 商户的订单号是否已经存在
func (b *BusinessBillDao) OutOrderNoExist(businessNo, outOrderNo string) (bool, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return false, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	sqlStr := "SELECT count(1) FROM business_bill WHERE business_no=$1 and out_order_no = $2 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&count}, businessNo, outOrderNo)
	return strext.ToInt(count.String) > 0, err
}

func (b *BusinessBillDao) GetOrderInfoByOrderNo(orderNo string) (*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var (
		accountNo, orderStatus, amount, realAmount, currencyType, createTime, payTime, businessNo,
		subject, remark, rate, fee, businessAccNo, expireTime, appId, sceneNo, tradeType sql.NullString
	)
	err := ss_sql.QueryRow(dbHandler, "SELECT account_no, order_status, amount, real_amount, currency_type, create_time, pay_time,"+
		"business_no, subject, remark, rate, fee, business_account_no, time_expire, app_id, scene_no, trade_type "+
		"FROM business_bill WHERE order_no= $1 limit 1 ",
		[]*sql.NullString{&accountNo, &orderStatus, &amount, &realAmount, &currencyType, &createTime, &payTime, &businessNo, &subject,
			&remark, &rate, &fee, &businessAccNo, &expireTime, &appId, &sceneNo, &tradeType}, orderNo)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessBillDao)
	obj.OrderNo = orderNo
	obj.AccountNo = accountNo.String
	obj.OrderStatus = orderStatus.String
	obj.Amount = amount.String
	obj.RealAmount = realAmount.String
	obj.CurrencyType = currencyType.String
	obj.CreateTime = createTime.String
	obj.PayTime = payTime.String
	obj.BusinessNo = businessNo.String
	obj.Subject = subject.String
	obj.Remark = remark.String
	obj.Rate = rate.String
	obj.Fee = fee.String
	obj.BusinessAccountNo = businessAccNo.String
	obj.ExpireTime = strext.ToInt64(expireTime.String)
	obj.AppId = appId.String
	obj.SceneNo = sceneNo.String
	obj.TradeType = tradeType.String
	return obj, nil
}

func (b *BusinessBillDao) GetOrderInfo(orderNo, outOrderNo string) (*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var (
		i    = 1
		args []interface{}
	)

	whereStr := "WHERE 1=1 "
	if orderNo != "" {
		whereStr += fmt.Sprintf("AND bl.order_no=$%v", i)
		args = append(args, orderNo)
		i++
	}
	if outOrderNo != "" {
		whereStr += fmt.Sprintf("AND bl.out_order_no=$%v", i)
		args = append(args, outOrderNo)
	}

	sqlStr := "SELECT bl.business_no, bl.app_id, bl.order_no, bl.out_order_no, bl.order_status, bl.amount," +
		"bl.currency_type, bl.subject, b.business_id, app.app_name, bl.account_no " +
		"FROM business_bill bl " +
		"LEFT JOIN business b ON b.business_no=bl.business_no " +
		"LEFT JOIN business_app app ON app.app_id = bl.app_id "
	var businessNo, appId, orderNoT, outOrderNoT, orderStatus, amount, currencyType, subject, businessId, appName, accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr+whereStr,
		[]*sql.NullString{&businessNo, &appId, &orderNoT, &outOrderNoT,
			&orderStatus, &amount, &currencyType, &subject, &businessId, &appName, &accountNo},
		args...)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessBillDao)
	obj.BusinessId = businessId.String
	obj.BusinessNo = businessNo.String
	obj.AppId = appId.String
	obj.AppName = appName.String
	obj.OrderNo = orderNoT.String
	obj.OutOrderNo = outOrderNoT.String
	obj.OrderStatus = orderStatus.String
	obj.Amount = amount.String
	obj.Subject = subject.String
	obj.CurrencyType = currencyType.String
	obj.AccountNo = accountNo.String

	return obj, nil
}

type BusinessBillSettleData struct {
	BusinessNo      string
	BusinessAccNo   string
	AppId           string
	CurrencyType    string
	TotalAmount     int64
	TotalRealAmount int64
	TotalFees       int64
	TotalOrder      int64
}

// 获取商户App的结算金额
func (b *BusinessBillDao) GetSettleData(startTime, endTime int64) ([]*BusinessBillSettleData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT business_no, business_account_no, app_id, currency_type,SUM(amount) AS t_amount,
		SUM(real_amount) AS t_real_amount, SUM(fee) AS t_fee,COUNT(order_no) AS t_order 
		FROM business_bill 
		WHERE settle_date >= $1 AND settle_date <= $2 AND order_status = $3 AND settle_id = '' 
		GROUP BY business_no, business_account_no, app_id, currency_type `
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, startTime, endTime, constants.BusinessOrderStatusPay)
	if nil != err {
		return nil, err
	}

	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*BusinessBillSettleData
	for rows.Next() {
		var businessNo, businessAccNo, appId, currencyType, tAmount, tRealAmount, tFee, tOrder sql.NullString
		err := rows.Scan(&businessNo, &businessAccNo, &appId, &currencyType, &tAmount, &tRealAmount, &tFee, &tOrder)
		if err != nil {
			return nil, err
		}
		ret := new(BusinessBillSettleData)
		ret.BusinessNo = businessNo.String
		ret.BusinessAccNo = businessAccNo.String
		ret.AppId = appId.String
		ret.CurrencyType = currencyType.String
		ret.TotalAmount = strext.ToInt64(tAmount.String)
		ret.TotalRealAmount = strext.ToInt64(tRealAmount.String)
		ret.TotalFees = strext.ToInt64(tFee.String)
		ret.TotalOrder = strext.ToInt64(tOrder.String)

		dataList = append(dataList, ret)
	}

	return dataList, nil
}

type UpdateOrderSettleId struct {
	SettleId   string
	BusinessNo string
	AppId      string
	StartTime  int64
	EndTime    int64
}

//修改订单settle_id
func (b *BusinessBillDao) UpdateSettleIdTx(tx *sql.Tx, data *UpdateOrderSettleId) error {
	sqlStr := "UPDATE business_bill SET settle_id=$1 " +
		"WHERE business_no = $2 AND app_id = $3 AND settle_date >= $4 AND settle_date <= $5 AND order_status=$6 AND settle_id='' "
	return ss_sql.ExecTx(tx, sqlStr, data.SettleId, data.BusinessNo, data.AppId, data.StartTime, data.EndTime, constants.BusinessOrderStatusPay)

}
func (b *BusinessBillDao) UpdateOrderSettleIdTx(tx *sql.Tx, settleId, orderNo string) error {
	sqlStr := "UPDATE business_bill SET settle_id=$1 WHERE order_no=$2 AND settle_id='' "
	return ss_sql.ExecTx(tx, sqlStr, settleId, orderNo)

}

type BusinessTransData struct {
	TotalAmount     int64
	TotalRealAmount int64
	TotalFees       int64
	TotalOrder      int64
}

//查询对账数据, 参数：是否已结算
func (b *BusinessBillDao) GetBusinessTransData(isSettled bool, businessNo, currencyType string) (*BusinessTransData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT SUM(amount) AS t_amount, SUM(real_amount) AS t_real_amount FROM business_bill `
	sqlStr += `WHERE  order_status=$1 AND business_no=$2 AND currency_type=$3 `
	if isSettled {
		sqlStr += `AND settle_id != '' GROUP BY currency_type `
	} else {
		sqlStr += `AND settle_id='' GROUP BY currency_type `
	}
	var tamount, trealAmount sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&tamount, &trealAmount},
		constants.BusinessOrderStatusPay, businessNo, currencyType,
	)
	if err != nil {
		return nil, err
	}

	data := new(BusinessTransData)
	data.TotalAmount = strext.ToInt64(tamount.String)
	data.TotalRealAmount = strext.ToInt64(trealAmount.String)

	return data, nil
}

func (*BusinessBillDao) GetBillByLogNoAndBusinessNo(logNo, businessNo string) (*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := ` SELECT order_no, fee, amount, real_amount, rate, order_status, business_no, business_vaccount_no,
		account_no, vaccount_no, currency_type, pay_time, settle_id, notify_url
		FROM business_bill
		WHERE order_no=$1 and business_no=$2 `

	var orderNo, fee, amount, realAmount, rate, orderStatus, businessNoT, businessVAccNo, accountNo, vAccNo, currencyType,
		payTime, settleId, notifyUrl sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&orderNo, &fee, &amount, &realAmount, &rate, &orderStatus,
		&businessNoT, &businessVAccNo, &accountNo, &vAccNo, &currencyType, &payTime, &settleId, &notifyUrl}, logNo, businessNo)
	if err != nil {
		return nil, err
	}

	data := new(BusinessBillDao)
	data.OrderNo = orderNo.String
	data.Fee = fee.String
	data.Amount = amount.String
	data.RealAmount = realAmount.String
	data.Rate = rate.String
	data.OrderStatus = orderStatus.String
	data.BusinessNo = businessNoT.String
	data.BusinessVaccountNo = businessVAccNo.String
	data.AccountNo = accountNo.String
	data.VAccountNo = vAccNo.String
	data.CurrencyType = currencyType.String
	data.PayTime = payTime.String
	data.SettleId = settleId.String
	return data, nil
}

func (*BusinessBillDao) GetBillByLogNoAndAppId(logNo, appId string) (*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := ` SELECT order_no, fee, amount, real_amount, rate, order_status, business_no, business_vaccount_no,
		account_no, vaccount_no, currency_type, pay_time, settle_id, notify_url
		FROM business_bill
		WHERE order_no=$1 and app_id=$2 `

	var orderNo, fee, amount, realAmount, rate, orderStatus, businessNoT, businessVAccNo, accountNo, vAccNo, currencyType,
		payTime, settleId, notifyUrl sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&orderNo, &fee, &amount, &realAmount, &rate, &orderStatus,
		&businessNoT, &businessVAccNo, &accountNo, &vAccNo, &currencyType, &payTime, &settleId, &notifyUrl}, logNo, appId)
	if err != nil {
		return nil, err
	}

	data := new(BusinessBillDao)
	data.OrderNo = orderNo.String
	data.Fee = fee.String
	data.Amount = amount.String
	data.RealAmount = realAmount.String
	data.Rate = rate.String
	data.OrderStatus = orderStatus.String
	data.BusinessNo = businessNoT.String
	data.BusinessVaccountNo = businessVAccNo.String
	data.AccountNo = accountNo.String
	data.VAccountNo = vAccNo.String
	data.CurrencyType = currencyType.String
	data.PayTime = payTime.String
	data.SettleId = settleId.String
	return data, nil
}

func (*BusinessBillDao) UpdateStatusByOrderNoTx(tx *sql.Tx, targetStatus, orderNo string) error {
	sqlStr := "UPDATE business_bill SET order_status=$1 WHERE order_no=$2 AND (order_status=$3 OR order_status=$4) "
	return ss_sql.ExecTx(tx, sqlStr, targetStatus, orderNo, constants.BusinessOrderStatusPay, constants.BusinessOrderStatusRebatesRefund)
}

type PendingPayOrder struct {
	OrderNo      string
	Amount       string
	CurrencyType string
	AppName      string
	Subject      string
	SimplifyName string
}

//查询用户待支付订单
func (*BusinessBillDao) GetCustPendingPayOrder(accountNo, orderNo string) (*PendingPayOrder, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT bb.order_no, bb.amount, bb.currency_type, bb.subject, app.app_name, bu.simplify_name
		FROM business_bill bb 
		LEFT JOIN business_app app ON app.app_id = bb.app_id
		LEFT JOIN business bu ON bu.business_no = bb.business_no
		WHERE bb.account_no=$1 AND bb.order_no=$2 AND bb.order_status=$3 AND bb.time_expire>$4 
		ORDER BY bb.create_time DESC LIMIT 1 `
	var orderNoT, amount, currencyType, subject, appName, simplifyName sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&orderNoT, &amount, &currencyType, &subject, &appName, &simplifyName},
		accountNo, orderNo, constants.BusinessOrderStatusPending, ss_time.Now(global.Tz).Add(5*time.Second).Unix())
	if err != nil {
		return nil, err
	}
	order := new(PendingPayOrder)
	order.OrderNo = orderNoT.String
	order.Amount = amount.String
	order.CurrencyType = currencyType.String
	order.Subject = subject.String
	order.AppName = appName.String
	order.SimplifyName = simplifyName.String
	return order, nil
}

// 应用-通过订单号查询订单(内部订单号或外部订单号，二选一)
func (b *BusinessBillDao) AppQueryOrder(appId, orderNoParam, outOrderNoParam string) (*BusinessBillDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var (
		orderNo, outOrderNo,
		accountNo, orderStatus, amount, realAmount, currencyType, createTime, payTime, businessNo, subject, remark, rate,
		fee, businessAccNo, businessVaccountNo, expireTime, settleId, payAccount sql.NullString
	)

	sqlStr := "SELECT b.order_no, b.out_order_no, b.account_no, b.order_status, b.amount, b.real_amount, b.currency_type, " +
		"b.create_time, b.pay_time, b.business_no, b.subject, b.remark, b.rate, b.fee, b.business_account_no, " +
		"b.business_vaccount_no, b.time_expire, b.settle_id, a.account AS pay_account "
	sqlStr += "FROM business_bill AS b LEFT JOIN account AS a ON b.account_no=a.uid WHERE b.app_id=$1 "

	queryParam := ""
	if orderNoParam != "" {
		sqlStr += " AND b.order_no=$2 LIMIT 1 "
		queryParam = orderNoParam
	} else if outOrderNoParam != "" {
		sqlStr += " AND b.out_order_no=$2 LIMIT 1 "
		queryParam = outOrderNoParam
	} else {
		return nil, errors.New("缺少必须参数")
	}

	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&orderNo, &outOrderNo, &accountNo, &orderStatus, &amount, &realAmount, &currencyType, &createTime,
			&payTime, &businessNo, &subject, &remark, &rate, &fee, &businessAccNo, &businessVaccountNo, &expireTime,
			&settleId, &payAccount},
		appId, queryParam,
	)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessBillDao)
	obj.OrderNo = orderNo.String
	obj.OutOrderNo = outOrderNo.String
	obj.AccountNo = accountNo.String
	obj.OrderStatus = orderStatus.String
	obj.Amount = amount.String
	obj.RealAmount = realAmount.String
	obj.CurrencyType = currencyType.String
	obj.CreateTime = createTime.String
	obj.PayTime = payTime.String
	obj.BusinessNo = businessNo.String
	obj.Subject = subject.String
	obj.Remark = remark.String
	obj.Rate = rate.String
	obj.Fee = fee.String
	obj.BusinessAccountNo = businessAccNo.String
	obj.BusinessVaccountNo = businessVaccountNo.String
	obj.ExpireTime = strext.ToInt64(expireTime.String)
	obj.SettleId = settleId.String
	obj.PayAccount = payAccount.String

	return obj, nil
}

func (b *BusinessBillDao) UpdateOrderNotifyStatus(orderNo, notifyStatus string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "UPDATE business_bill SET notify_status=$1 WHERE order_no=$2 AND order_status=$3 "
	return ss_sql.Exec(dbHandler, sqlStr, notifyStatus, orderNo, constants.BusinessOrderStatusPay)
}

func (b *BusinessBillDao) GetOrderChannelNo(orderNo string) (channelNo, channelType string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bc.business_channel_no, bc.channel_type " +
		"FROM business_bill bb " +
		"LEFT JOIN business_channel bc ON bc.business_channel_no = bb.business_channel_no " +
		"WHERE bb.order_no = $1 "

	var channelNoT, channelTypeT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&channelNoT, &channelTypeT}, orderNo); err != nil {
		return "", "", err
	}

	return channelNoT.String, channelTypeT.String, nil
}
