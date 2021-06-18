package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type AppFingerprintSignDao struct {
	Id         string
	SignKey    string
	AccountNo  string
	Account    string
	DeviceUuid string
	NonStr     string
	OpenTime   string
	CreateTime string
	ModifyTime string
	UseStatus  string
}

var AppFingerprintSignDaoInst AppFingerprintSignDao

func (AppFingerprintSignDao) Count(whereList string, args []interface{}) (int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT  COUNT(1) " +
		"FROM app_fingerprint_sign afs " +
		"LEFT JOIN account acc ON acc.uid = afs.account_no "

	sqlStr += whereList
	var total sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, args...)
	if err != nil {
		return 0, nil
	}

	return strext.ToInt32(total.String), nil
}

func (AppFingerprintSignDao) GetList(whereList string, args []interface{}) ([]*AppFingerprintSignDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT afs.id, afs.device_uuid, acc.account, afs.open_time, afs.use_status " +
		"FROM app_fingerprint_sign afs " +
		"LEFT JOIN account acc ON acc.uid = afs.account_no "

	sqlStr += whereList
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*AppFingerprintSignDao
	for rows.Next() {
		var id, deviceNo, account, openTime, useStatus sql.NullString
		err := rows.Scan(&id, &deviceNo, &account, &openTime, &useStatus)
		if err != nil {
			return nil, err
		}
		data := new(AppFingerprintSignDao)
		data.Id = id.String
		data.DeviceUuid = deviceNo.String
		data.Account = account.String
		data.OpenTime = openTime.String
		data.UseStatus = useStatus.String
		dataList = append(dataList, data)
	}

	return dataList, nil
}

func (AppFingerprintSignDao) UpdateUseStatusSingle(tx *sql.Tx, useStatus, id string) error {
	sqlStr := "UPDATE app_fingerprint_sign set use_status = $1, modify_time=CURRENT_TIMESTAMP WHERE id = $2 "
	return ss_sql.ExecTx(tx, sqlStr, useStatus, id)
}

func (AppFingerprintSignDao) BanAll(tx *sql.Tx) error {
	sqlStr := "UPDATE app_fingerprint_sign set use_status = $1, modify_time=CURRENT_TIMESTAMP WHERE use_status = $2 "
	return ss_sql.ExecTx(tx, sqlStr, constants.AppFingerprintUseStatus_Disable, constants.AppFingerprintUseStatus_Enable)
}
