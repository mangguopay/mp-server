package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type AppFingerprintDao struct{}

var AppFingerprintDaoInstance AppFingerprintDao

//确认指纹支付标识是存在并且是有开启的
func (AppFingerprintDao) CheckSignKey(accountNo, deviceUuid, signKey string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT COUNT(1) FROM app_fingerprint_sign WHERE account_no = $1 AND device_uuid = $2 AND sign_key = $3 AND use_status = $4 `
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, accountNo, deviceUuid, signKey, constants.AppFingerprintUseStatus_Enable); err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	return cnt.String != "0"
}
