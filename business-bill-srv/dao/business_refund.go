package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type BusinessRefundOrderDao struct {
	RefundNo     string
	OutRefundNo  string
	PayOrderNo   string
	Amount       string
	RefundStatus string
	Remarks      string
	CreateTime   string
	FinishTime   string
	NotifyUrl    string
	NotifyStatus string
}

var BusinessRefundOrderDaoInst BusinessRefundOrderDao

func (*BusinessRefundOrderDao) Insert(d *BusinessRefundOrderDao) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "INSERT INTO business_refund_order(refund_no, amount, refund_status, pay_order_no, remark, out_refund_no," +
		"notify_status, notify_url, create_time) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)"
	refundNo := strext.GetDailyId()
	err := ss_sql.Exec(dbHandler, sqlStr, refundNo, d.Amount, d.RefundStatus, d.PayOrderNo, d.Remarks, d.OutRefundNo,
		d.NotifyStatus, d.NotifyUrl)
	if err != nil {
		return "", err
	}
	return refundNo, nil
}

func (*BusinessRefundOrderDao) UpdatePendingOrderByOrderNo(orderStatus, refundNo string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "UPDATE business_refund_order SET refund_status=$1, finish_time=CURRENT_TIMESTAMP WHERE refund_no=$2 AND refund_status=$3 "
	return ss_sql.Exec(dbHandler, sqlStr, orderStatus, refundNo, constants.BusinessRefundStatusPending)
}

func (*BusinessRefundOrderDao) UpdatePendingOrderByOrderNoTx(tx *sql.Tx, orderStatus, refundNo string) error {
	sqlStr := "UPDATE business_refund_order SET refund_status=$1, finish_time=CURRENT_TIMESTAMP WHERE refund_no=$2 AND refund_status=$3 "
	return ss_sql.ExecTx(tx, sqlStr, orderStatus, refundNo, constants.BusinessRefundStatusPending)
}

//统计交易订单已退款金额
func (*BusinessRefundOrderDao) GetTotalAmountByPayOrderNo(payOrderNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT SUM(amount) FROM business_refund_order WHERE pay_order_no=$1 AND refund_status=$2 "
	var totalAmount sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalAmount}, payOrderNo, constants.BusinessRefundStatusSuccess)
	if err != nil {
		return "", err
	}
	return totalAmount.String, nil
}

//统计交易订单已退记录
func (*BusinessRefundOrderDao) CountRefundByPayOrderNo(payOrderNo string) (num, amount string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT COUNT(refund_no) AS total_num, SUM(amount) AS total_amount  " +
		"FROM business_refund_order WHERE pay_order_no=$1 AND refund_status=$2 "
	var totalAmount, totalNum sql.NullString
	err = ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalNum, &totalAmount}, payOrderNo, constants.BusinessRefundStatusSuccess)
	if err != nil {
		return "", "", err
	}
	return totalNum.String, totalAmount.String, nil
}

//查询交易订单退款详情
func (*BusinessRefundOrderDao) GetRefundByPayOrderNo(payOrderNo string) ([]*BusinessRefundOrderDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT refund_no, out_refund_no, amount, refund_status, finish_time " +
		"FROM business_refund_order " +
		"WHERE pay_order_no=$1  "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, payOrderNo)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var list []*BusinessRefundOrderDao
	for rows.Next() {
		var refundNo, outRefundNo, amount, refundStatus, finishTime sql.NullString
		err := rows.Scan(&refundNo, &outRefundNo, &amount, &refundStatus, &finishTime)
		if err != nil {
			return nil, err
		}
		data := &BusinessRefundOrderDao{
			RefundNo:     refundNo.String,
			OutRefundNo:  outRefundNo.String,
			Amount:       amount.String,
			RefundStatus: refundStatus.String,
			FinishTime:   finishTime.String,
		}
		list = append(list, data)
	}

	return list, nil
}
func (*BusinessRefundOrderDao) GetRefundByRefundNo(appId, refundNo, outRefundNo string) (*BusinessRefundOrderDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var (
		i    = 1
		args []interface{}
	)

	whereStr := "WHERE 1=1 "
	if appId != "" {
		whereStr += fmt.Sprintf("AND bb.app_id =$%v", i)
		args = append(args, appId)
		i++
	}

	if refundNo != "" {
		whereStr += fmt.Sprintf("AND refund_no=$%v", i)
		args = append(args, refundNo)
		i++
	}
	if outRefundNo != "" {
		whereStr += fmt.Sprintf("AND out_refund_no=$%v", i)
		args = append(args, outRefundNo)
	}

	if len(args) == 0 {
		return nil, errors.New("refundNo，outRefundNo")
	}

	sqlStr := "SELECT br.refund_no, br.out_refund_no, br.amount, br.refund_status, br.finish_time, br.pay_order_no " +
		"FROM business_refund_order br " +
		"LEFT JOIN business_bill bb ON bb.order_no = br.pay_order_no "
	sqlStr += whereStr
	var refundNoT, outRefundNoT, amount, refundStatus, finishTime, payOrderNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&refundNoT, &outRefundNoT, &amount,
		&refundStatus, &finishTime, &payOrderNo}, args...)
	if err != nil {
		return nil, err
	}
	data := &BusinessRefundOrderDao{
		RefundNo:     refundNoT.String,
		OutRefundNo:  outRefundNoT.String,
		Amount:       amount.String,
		RefundStatus: refundStatus.String,
		FinishTime:   finishTime.String,
		PayOrderNo:   payOrderNo.String,
	}

	return data, nil
}

//修改超时订单为支付失败
func (*BusinessRefundOrderDao) UpdateTimeOutOrderStatus() ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "UPDATE business_refund_order SET refund_status=$1,finish_time=CURRENT_TIMESTAMP WHERE refund_status=$2 AND create_time <$3 RETURNING refund_no "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, constants.BusinessRefundStatusFail, constants.BusinessRefundStatusPending,
		ss_time.Now(global.Tz).Add(-constants.BusinessOrderExpireTime*time.Minute).Format(ss_time.DateTimeDashFormat))
	if err != nil {
		return nil, err
	}

	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var refundNo sql.NullString
		if err := rows.Scan(&refundNo); nil != err {
			return nil, err
		}
		list = append(list, refundNo.String)
	}

	return list, nil
}
