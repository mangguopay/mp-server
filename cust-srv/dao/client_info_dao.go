package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ClientInfoDao struct {
}

var ClientInfoDaoInst ClientInfoDao

func (*ClientInfoDao) GetCustCnt(whereModelStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from client_info_app cli " + whereModelStr
	var totalT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereArgs...)

	return totalT.String, errT
}

func (*ClientInfoDao) GetCustClientInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.ClientInfoData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT cli.id,cli.device_brand,cli.device_model,cli.resolution,cli.screen_size" +
		",cli.imei1,cli.imei2,cli.system_ver,cli.create_time,cli.upload_point" +
		",cli.user_agent,cli.platform,cli.app_ver,cli.account,cli.uuid " +
		" FROM client_info_app cli " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.ClientInfoData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.ClientInfoData{}
		err2 := rows.Scan(
			&data.Id,
			&data.DeviceBrand,
			&data.DeviceModel,
			&data.Resolution,
			&data.ScreenSize,

			&data.Imei1,
			&data.Imei2,
			&data.SystemVer,
			&data.CreateTime,
			&data.UploadPoint,

			&data.UserAgent,
			&data.Platform,
			&data.AppVer,
			&data.Account,
			&data.Uuid,
		)
		if err2 != nil {
			ss_log.Error("client_info_app表查询出错。id：[%v],err:[%v]", data.Id, err2)
			continue
		}
		datas = append(datas, data)
	}

	return datas, err
}

//servicer
func (*ClientInfoDao) GetSerCnt(whereModelStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from client_info_pos cli  " + whereModelStr
	var totalT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereArgs...)

	return totalT.String, errT
}

func (*ClientInfoDao) GetSerClientInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.ClientInfoData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT cli.id,cli.device_brand,cli.device_model,cli.resolution,cli.screen_size" +
		",cli.imei1,cli.imei2,cli.system_ver,cli.create_time,cli.upload_point" +
		",cli.user_agent,cli.platform,cli.app_ver,cli.account,cli.uuid " +
		" FROM client_info_pos cli " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.ClientInfoData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.ClientInfoData{}
		err2 := rows.Scan(
			&data.Id,
			&data.DeviceBrand,
			&data.DeviceModel,
			&data.Resolution,
			&data.ScreenSize,

			&data.Imei1,
			&data.Imei2,
			&data.SystemVer,
			&data.CreateTime,
			&data.UploadPoint,

			&data.UserAgent,
			&data.Platform,
			&data.AppVer,
			&data.Account,
			&data.Uuid,
		)
		if err2 != nil {
			ss_log.Error("client_info_pos表查询出错。id：[%v],err:[%v]", data.Id, err2)
			continue
		}
		datas = append(datas, data)
	}

	return datas, err
}
