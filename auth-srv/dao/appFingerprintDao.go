package dao

import (
	"a.a/cu/db"
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type AppFingerprintDao struct {
}

var AppFingerprintDaoInstance AppFingerprintDao

//	A = sha256(设备id+账号id+开启时间+随机字符串A)
//	B = sha256(A+随机字符串B)
//	用户指纹支付标识id: C = A+B+10位随机数
//
func createSignKey(accountNo, devicUuid, openTime string) string {
	nonStrA := util.RandomDigitStrOnlyAlphabet(16)
	nonStrB := util.RandomDigitStrOnlyAlphabet(16)
	nonStrC := util.RandomDigitStrOnlyAlphabet(10)

	strBase := devicUuid + accountNo + openTime + nonStrA

	strA := encrypt.DoSha256(strBase, nonStrA)
	strB := encrypt.DoSha256(strA+nonStrB, nonStrB)

	signKey := strA + strB + nonStrC
	return signKey
}

func (*AppFingerprintDao) AddTx(tx *sql.Tx, accountNo, devicUuid string) (signKey string, err error) {
	id := strext.GetDailyId()
	//openTime := time.Now()
	openTime := ss_time.Now(global.Tz)
	signKey = createSignKey(accountNo, devicUuid, openTime.Format("20060102150405"))

	sqlStr := "INSERT INTO app_fingerprint_sign(id, sign_key, account_no, device_uuid, " +
		" use_status, open_time, create_time) VALUES($1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP)"
	if err := ss_sql.ExecTx(tx, sqlStr, id, signKey, accountNo, devicUuid, constants.AppFingerprintUseStatus_Enable, openTime); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", err
	}
	return signKey, nil
}

//设置该账号、设备的可能存在的旧指纹标识无效、禁用
func (*AppFingerprintDao) SetUseStatusDisableByAccountTx(tx *sql.Tx, accountNo, devicUuid string) error {
	sqlStr := "UPDATE app_fingerprint_sign SET use_status = $4, modify_time = CURRENT_TIMESTAMP WHERE account_no = $1 AND device_uuid = $2 AND use_status = $3 "
	if err := ss_sql.ExecTx(tx, sqlStr, accountNo, devicUuid, constants.AppFingerprintUseStatus_Enable, constants.AppFingerprintUseStatus_Disable); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return err
	}
	return nil
}

//设置该账号、设备的可能存在的旧指纹标识无效、禁用
func (*AppFingerprintDao) SetUseStatusDisableByAccount(accountNo, devicUuid string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "UPDATE app_fingerprint_sign SET use_status = $4, modify_time = CURRENT_TIMESTAMP WHERE account_no = $1 AND device_uuid = $2 AND use_status = $3 "
	if err := ss_sql.Exec(dbHandler, sqlStr, accountNo, devicUuid, constants.AppFingerprintUseStatus_Enable, constants.AppFingerprintUseStatus_Disable); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return err
	}
	return nil
}
