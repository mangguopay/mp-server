package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessSettleOneDao struct {
	SettleId string
}

var BusinessSettleOneDaoInst BusinessSettleOneDao

func (BusinessSettleOneDao) InsertTx(tx *sql.Tx, d *BusinessSettleOneDao) (string, error) {
	settleId := strext.GetDailyId()
	sqlStr := `INSERT INTO business_bill_settle_one (settle_id, create_time) VALUES ($1, current_timestamp) `
	if err := ss_sql.ExecTx(tx, sqlStr, settleId); err != nil {
		return "", err
	}
	return settleId, nil
}
