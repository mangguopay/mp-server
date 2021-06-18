package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type LogAppFingerprintPayDao struct{}

var LogAppFingerprintPayDaoInstance LogAppFingerprintPayDao

type LogAppFingerprintPayData struct {
	AccountNo    string
	DeviceUuid   string
	SignKey      string
	OrderNo      string
	OrderType    string
	Amount       string
	CurrencyType string
}

func (LogAppFingerprintPayDao) AddTx(TX *sql.Tx, data *LogAppFingerprintPayData) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	logNo := strext.GetDailyId()

	sqlStr := `INSERT INTO log_app_fingerprint_pay(log_no, account_no, device_uuid, sign_key, order_no,
		order_type, amount, currency_type, create_time) VALUES($1,$2,$3,$4,$5,$6,$7,$8,CURRENT_TIMESTAMP) `
	if err := ss_sql.ExecTx(TX, sqlStr, logNo, data.AccountNo, data.DeviceUuid, data.SignKey, data.OrderNo,
		data.OrderType, data.Amount, data.CurrencyType); err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	return nil
}
