package dao

import (
	"database/sql"

	"a.a/mp-server/common/ss_sql"
)

var BusinessBillSettDaoInst BusinessBillSettleDao

type BusinessBillSettleDao struct {
	SettleId        string
	BusinessNo      string
	StartTime       string
	EndTime         string
	TotalAmount     int64
	TotalRealAmount int64
	TotalFees       int64
	TotalOrder      int64
	CurrencyType    string
	AppId           string
}

func (*BusinessBillSettleDao) InsertTx(tx *sql.Tx, d *BusinessBillSettleDao) error {
	sqlStr := `INSERT INTO business_bill_settle (settle_id, create_time, business_no, start_time, end_time, `
	sqlStr += ` total_amount, total_real_amount, total_fees, total_order, currency_type, app_id) `
	sqlStr += ` VALUES ($1, current_timestamp, $2, $3, $4, $5, $6, $7, $8, $9, $10) `

	err := ss_sql.ExecTx(tx, sqlStr, d.SettleId, d.BusinessNo, d.StartTime, d.EndTime,
		d.TotalAmount, d.TotalRealAmount, d.TotalFees, d.TotalOrder, d.CurrencyType, d.AppId,
	)

	return err
}
