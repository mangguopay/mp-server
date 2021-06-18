package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BillDao struct{}

var BillDaoInst BillDao

// 下单入库
func (BillDao) InsertOrder(reqNo, innerOrderNo, amount, fee, channelNo, rateNo, rate, accNo, supplierCode, productType, callbackUrl, callbackSelf string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `insert into business_bill(order_no,amount,req_no,create_time,channel_no,rate_no,rate,fee,order_status,acc_no,supplier_code,product_type,callback_url,callback_self) `+
		`values($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		innerOrderNo, amount, reqNo, channelNo, rateNo, rate, fee, constants.OrderStatus_Pending, accNo, supplierCode, productType, callbackUrl, callbackSelf)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "bill_req_no_acc_no_key"`:
			return ss_err.ERR_PAY_DUP_ORDER
		default:
		}
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS
}

// 订单信息修改
func (BillDao) UpdateRet(innerOrderNo, retMsg, orderStatus string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update business_bill set ret_msg=$2,order_status=$3,finish_time=current_timestamp where order_no=$1`,
		innerOrderNo, retMsg, orderStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PAY_UPDATE_ORDER
	}

	return ss_err.ERR_SUCCESS
}

// 订单信息
func (BillDao) GetBillOrderInfoFromReqNo(reqNo, accNo string) (orderStatus, supplierCode, innerOrderNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var orderStatusT, supplierCodeT, innerOrderNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select order_status,supplier_code,order_no from business_bill where req_no=$1 and acc_no=$2 limit 1`,
		[]*sql.NullString{&orderStatusT, &supplierCodeT, &innerOrderNoT}, reqNo, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return orderStatusT.String, supplierCodeT.String, innerOrderNoT.String
	}

	return orderStatusT.String, supplierCodeT.String, innerOrderNoT.String
}
func (BillDao) GetBillOrderInfoFromReqNo1(reqNo, accNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var orderNoT, amountT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select  order_no,amount from business_bill where req_no=$1 and acc_no=$2 limit 1`,
		[]*sql.NullString{&orderNoT, &amountT}, reqNo, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ""
	}

	return orderNoT.String, amountT.String
}

// 订单信息
func (BillDao) GetBillOrderInfo(innerOrderNo string) (orderStatus, createTime string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var orderStatusT, createTimeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select order_status,create_time from business_bill where order_no=$1 limit 1`,
		[]*sql.NullString{&orderStatusT, &createTimeT}, innerOrderNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return orderStatusT.String, createTimeT.String
	}

	return orderStatusT.String, createTimeT.String
}

func (BillDao) GetBillChannelNoFromOrderNo(innerOrderNo string) (channelNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var channelNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select channel_no from business_bill where order_no=$1 limit 1`,
		[]*sql.NullString{&channelNoT}, innerOrderNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return channelNoT.String
	}

	return channelNoT.String
}

// 订单信息
func (BillDao) GetBillFee(innerOrderNo string) (fee string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var feeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select fee from business_bill where order_no=$1 limit 1`,
		[]*sql.NullString{&feeT}, innerOrderNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return feeT.String
	}

	return feeT.String
}
func (BillDao) GetBillAccNo(innerOrderNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select acc_no from business_bill where order_no=$1 limit 1`,
		[]*sql.NullString{&accNoT}, innerOrderNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return accNoT.String
}

// 订单信息
func (BillDao) GetBillNotifyUrl(innerOrderNo string) (notifyUrl string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var notifyUrlT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select callback_url from business_bill where order_no=$1 limit 1`,
		[]*sql.NullString{&notifyUrlT}, innerOrderNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return notifyUrlT.String
	}

	return notifyUrlT.String
}

func (BillDao) GetBillChannelNo(innerOrderNo string) (channelNo, amount string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var channelNoT, amountT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select channel_no,amount from business_bill where order_no=$1 limit 1`,
		[]*sql.NullString{&channelNoT, &amountT}, innerOrderNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return channelNoT.String, amountT.String
	}

	return channelNoT.String, amountT.String
}
