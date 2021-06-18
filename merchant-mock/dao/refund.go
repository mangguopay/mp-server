package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"a.a/cu/strext"

	"a.a/cu/db"
	"a.a/mp-server/common/ss_sql"
)

var RefundInstance Refund

type Refund struct {
	OutRefundNo  string
	CurrencyType string
	Amount       int64
	CreateTime   string
	RefundNo     string
	RefundTime   string
	AppId        string
	OrderSn      string
	Status       int64
	StatusStr    string
}

const (
	RefundStatusPending = 1 // 申请中
	RefundStatusSuccess = 2 // 成功
	RefundStatusFail    = 3 // 失败
)

func GetRefundStatusString(refundStatus int64) string {
	switch refundStatus {
	case RefundStatusPending:
		return "申请中"
	case RefundStatusSuccess:
		return "成功"
	case RefundStatusFail:
		return "失败"
	}

	return fmt.Sprintf("%d", refundStatus)
}

// 插入一条记录
func (o *Refund) Insert(refund *Refund) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	refund.OutRefundNo = "refund" + strext.GetDailyId()

	sqlStr := "INSERT INTO refund (create_time, out_refund_no, amount, currency_type, status, app_id)"
	sqlStr += " VALUES (current_timestamp, $1, $2, $3, $4, $5)"

	execErr := ss_sql.Exec(dbHandler, sqlStr,
		refund.OutRefundNo, refund.Amount, refund.CurrencyType, RefundStatusPending, refund.AppId)

	return execErr
}

func (o *Refund) GetRefundList(page, pageSize int) ([]Refund, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	offset := (page - 1) * pageSize

	sqlStr := "SELECT out_refund_no, amount, currency_type, status, app_id, create_time, order_sn, refund_no, refund_time FROM refund  "
	sqlStr += " ORDER BY create_time DESC LIMIT $1 OFFSET $2 "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, pageSize, offset)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		return nil, qErr
	}

	list := []Refund{}

	for rows.Next() {
		var outRefundNo, amount, currencyType, status, appId, createTime, orderSn, refundNo, refundTime sql.NullString

		err := rows.Scan(&outRefundNo, &amount, &currencyType, &status, &appId, &createTime, &orderSn, &refundNo, &refundTime)
		if err != nil {
			return nil, err
		}
		list = append(list, Refund{
			OutRefundNo:  outRefundNo.String,
			Amount:       strext.ToInt64(amount.String),
			CurrencyType: currencyType.String,
			AppId:        appId.String,
			CreateTime:   createTime.String,
			Status:       strext.ToInt64(status.String),
			StatusStr:    GetRefundStatusString(strext.ToInt64(status.String)),
			OrderSn:      orderSn.String,
			RefundNo:     refundNo.String,
			RefundTime:   refundTime.String,
		})
	}

	return list, nil
}

// 更新退款单号
func (o *Refund) UpdateRefundNo(outRefundNo string, refundNo string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE refund SET refund_no=$1 WHERE out_refund_no=$2`

	err := ss_sql.Exec(dbHandler, sqlStr, refundNo, outRefundNo)
	if nil != err {
		return err
	}

	return nil
}

// 通过订单号获取一条记录
func (o *Refund) GetOneByOutRefundNo(paramOutRefundNo string) (*Refund, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := "SELECT out_refund_no, amount, currency_type, status, app_id, create_time, order_sn, refund_no, refund_time  FROM refund  "
	sqlStr += " WHERE out_refund_no=$1 "

	var outRefundNo, amount, currencyType, status, appId, createTime, orderSn, refundNo, refundTime sql.NullString

	qErr := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&outRefundNo, &amount, &currencyType, &status, &appId, &createTime, &orderSn, &refundNo, &refundTime},
		paramOutRefundNo,
	)

	if qErr != nil {
		return nil, qErr
	}

	order := &Refund{
		OutRefundNo:  outRefundNo.String,
		Amount:       strext.ToInt64(amount.String),
		CurrencyType: currencyType.String,
		AppId:        appId.String,
		CreateTime:   createTime.String,
		Status:       strext.ToInt64(status.String),
		StatusStr:    GetRefundStatusString(strext.ToInt64(status.String)),
		OrderSn:      orderSn.String,
		RefundNo:     refundNo.String,
		RefundTime:   refundTime.String,
	}

	return order, nil
}

// 更新订单退款成功
func (o *Refund) UpdateRefundSuccess(outRefundferNo string, refundTime string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE refund SET status=$1, refund_time=$2 WHERE out_refund_no=$3`

	err := ss_sql.Exec(dbHandler, sqlStr, RefundStatusSuccess, refundTime, outRefundferNo)
	if nil != err {
		return err
	}

	return nil
}

// 更新订单退款失败
func (o *Refund) UpdateRefundFail(outRefundferNo string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE refund SET status=$1 WHERE out_refund_no=$2`

	err := ss_sql.Exec(dbHandler, sqlStr, RefundStatusFail, outRefundferNo)
	if nil != err {
		return err
	}

	return nil
}
