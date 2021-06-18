package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessBillQrCode struct {
}

var BusinessBillQrCodeInst BusinessBillQrCode

//查询订单号. 二维码id
func (*BusinessBillQrCode) InsertOrderQrCode(tx *sql.Tx, orderNo, payQrCodeId string) error {
	sqlStr := `INSERT INTO  business_bill_qrcode (order_no, pay_qrcode_id) VALUES($1,$2)`
	return ss_sql.ExecTx(tx, sqlStr, orderNo, payQrCodeId)
}

//查询订单号
func (*BusinessBillQrCode) QueryOrderNoByQrCodeId(qrCodeId string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var orderNo sql.NullString
	sqlStr := `SELECT order_no FROM business_bill_qrcode WHERE pay_qrcode_id=$1 `
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&orderNo}, qrCodeId)
	return orderNo.String, err
}
