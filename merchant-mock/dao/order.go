package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"a.a/mp-server/common/constants"

	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

var OrderInstance Order

type Order struct {
	OrderSn      string
	Title        string
	CurrencyType string
	Amount       int64
	Status       int64
	StatusStr    string
	CreateTime   string
	PayTime      string
	PayAccount   string
	PayOrderSn   string
	QrCode       string
	AppId        string
	TradeType    string

	// 附加字段
	TradeTypeName string
}

const (
	OrderStatusUnpay    = 1 // 待支付
	OrderStatusPaying   = 2 // 支付中
	OrderStatusPaid     = 3 // 已支付
	OrderStatusPaidFail = 4 // 支付失败
)

func GetOrderStatusString(orderStatus int64) string {
	switch orderStatus {
	case OrderStatusUnpay:
		return "待支付"
	case OrderStatusPaying:
		return "支付中"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusPaidFail:
		return "支付失败"
	}

	return fmt.Sprintf("%d", orderStatus)
}

// 通过订单号获取一条记录
func (o *Order) GetOneByOrderSn(paramsOrderSn string) (*Order, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := "SELECT order_sn, title, currency_type, amount, status, create_time, pay_time, pay_account, pay_order_sn,qr_code, trade_type, app_id FROM orders  "
	sqlStr += " WHERE order_sn=$1 "

	var orderSn, title, currencyType, amount, status, createTime, payTime, payAccount, payOrderSn, qrCode, tradeType, appId sql.NullString

	qErr := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&orderSn, &title, &currencyType, &amount, &status, &createTime, &payTime, &payAccount, &payOrderSn, &qrCode, &tradeType, &appId},
		paramsOrderSn,
	)

	if qErr != nil {
		return nil, qErr
	}

	order := &Order{
		OrderSn:       orderSn.String,
		Title:         title.String,
		CurrencyType:  currencyType.String,
		Amount:        strext.ToInt64(amount.String),
		Status:        strext.ToInt64(status.String),
		StatusStr:     GetOrderStatusString(strext.ToInt64(status.String)),
		CreateTime:    createTime.String,
		PayTime:       payTime.String,
		PayAccount:    payAccount.String,
		PayOrderSn:    payOrderSn.String,
		QrCode:        qrCode.String,
		TradeType:     tradeType.String,
		TradeTypeName: GetTradeTypeName(tradeType.String),
		AppId:         appId.String,
	}

	return order, nil
}

func (o *Order) GetOrderList(page, pageSize int) ([]Order, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	offset := (page - 1) * pageSize

	sqlStr := "SELECT order_sn, title, currency_type, amount, status, create_time, pay_time, pay_account, pay_order_sn, trade_type FROM orders  "
	sqlStr += " ORDER BY create_time DESC LIMIT $1 OFFSET $2 "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, pageSize, offset)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		return nil, qErr
	}

	list := []Order{}

	for rows.Next() {
		var orderSn, title, currencyType, amount, status, createTime, payTime, payAccount, payOrderSn, tradeType sql.NullString

		err := rows.Scan(&orderSn, &title, &currencyType, &amount, &status, &createTime, &payTime, &payAccount, &payOrderSn, &tradeType)
		if err != nil {
			return nil, err
		}
		list = append(list, Order{
			OrderSn:      orderSn.String,
			Title:        title.String,
			CurrencyType: currencyType.String,
			Amount:       strext.ToInt64(amount.String),
			Status:       strext.ToInt64(status.String),
			StatusStr:    GetOrderStatusString(strext.ToInt64(status.String)),
			CreateTime:   createTime.String,

			PayTime: payTime.String,

			PayAccount:    payAccount.String,
			PayOrderSn:    payOrderSn.String,
			TradeType:     tradeType.String,
			TradeTypeName: GetTradeTypeName(tradeType.String),
		})
	}

	return list, nil
}

func GetTradeTypeName(tradeType string) string {
	switch tradeType {
	case constants.TradeTypeModernpayFaceToFace:
		return "当面付"
	case constants.TradeTypeModernpayAPP:
		return "APP支付"
	}

	return tradeType
}

// 插入一条记录
func (o *Order) Insert(order *Order) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	orderSn := "merchant" + order.OrderSn

	sqlStr := "INSERT INTO orders (order_sn, title, currency_type, amount, status, create_time, app_id, trade_type)"
	sqlStr += " VALUES ($1, $2, $3, $4, $5, current_timestamp, $6, $7)"

	execErr := ss_sql.Exec(dbHandler, sqlStr, orderSn, order.Title, order.CurrencyType, order.Amount, OrderStatusUnpay, order.AppId, order.TradeType)

	return execErr
}

// 更新支付中的状态
func (o *Order) UpdatePayingStatus(orderNo string, payOrderSn, qrCode string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE orders SET status=$1, pay_order_sn=$2, qr_code=$3 WHERE order_sn=$4`

	err := ss_sql.Exec(dbHandler, sqlStr, OrderStatusPaying, payOrderSn, qrCode, orderNo)
	if nil != err {
		return err
	}

	return nil
}

// 更新支付中的状态
func (o *Order) UpdateStatus(orderNo string, orderStatus int) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE orders SET status=$1 WHERE order_sn=$2`

	err := ss_sql.Exec(dbHandler, sqlStr, orderStatus, orderNo)
	if nil != err {
		return err
	}

	return nil
}

// 更新订单支付成功
func (o *Order) UpdatePaidOk(orderNo string, payTime string, payAccount string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE orders SET status=$1, pay_time=$2, pay_account=$3 WHERE order_sn=$4`

	err := ss_sql.Exec(dbHandler, sqlStr, OrderStatusPaid, payTime, payAccount, orderNo)
	if nil != err {
		return err
	}

	return nil
}
