package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var BusinessAppDaoInst BusinessAppDao

type BusinessAppDao struct {
	AppId        string
	AppName      string
	BusinessNo   string
	Status       string
	SignMethod   string
	SimplifyName string //商家简称
}

func (b *BusinessAppDao) GetAppInfoByAppId(appId string) (*BusinessAppDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var appIdT, appName, businessNo, signMethod, status, simplifyName sql.NullString
	sqlStr := "SELECT app.app_id, app.app_name, app.business_no, app.sign_method, app.status, bu.simplify_name " +
		" FROM business_app app " +
		" LEFT JOIN business bu ON bu.business_no = app.business_no " +
		" WHERE app.app_id = $1 LIMIT 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&appIdT, &appName, &businessNo, &signMethod, &status, &simplifyName},
		appId,
	)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessAppDao)
	obj.AppId = appIdT.String
	obj.AppName = appName.String
	obj.BusinessNo = businessNo.String
	obj.SignMethod = signMethod.String
	obj.Status = status.String
	obj.SimplifyName = simplifyName.String
	return obj, nil
}

// 通过固定二维码查询应用基本信息
func (b *BusinessAppDao) GetAppInfoByFixedQrCode(fixedQrCode string) (*BusinessAppDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var appId, businessNo, appName, signMethod, status, simplifyName sql.NullString

	sqlStr := "SELECT app.app_id, app.business_no, app.app_name, app.sign_method, app.status, bu.simplify_name  " +
		" FROM business_app app " +
		" LEFT JOIN business bu ON bu.business_no = app.business_no " +
		"WHERE app.fixed_qrcode=$1  LIMIT 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&appId, &businessNo, &appName, &signMethod, &status, &simplifyName}, fixedQrCode)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessAppDao)
	obj.AppId = appId.String
	obj.AppName = appName.String
	obj.SignMethod = signMethod.String
	obj.BusinessNo = businessNo.String
	obj.Status = status.String
	obj.SimplifyName = simplifyName.String

	return obj, nil
}
