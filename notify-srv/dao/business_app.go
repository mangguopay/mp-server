package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type BusinessAppDao struct {
}

type BusinessApp struct {
	AppId           string
	BusinessNo      string
	IpWhiteList     string
	Status          int
	SignMethod      string
	BusinessPubKey  string
	PlatformPubKey  string
	PlatformPrivKey string
}

var BusinessAppDaoInst BusinessAppDao

//查询商家签名信息
func (*BusinessAppDao) GetSignInfo(appId string) (*BusinessApp, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var signMethod, ipWhiteList, businessPubKey, platformPrivKey, isOpen sql.NullString

	sqlStr := `SELECT sign_method, ip_white_list, business_pub_key, platform_priv_key, status FROM business_app WHERE app_id=$1 LIMIT 1`
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&signMethod, &ipWhiteList, &businessPubKey, &platformPrivKey, &isOpen}, appId)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessApp)
	obj.AppId = appId
	obj.SignMethod = signMethod.String
	obj.IpWhiteList = ipWhiteList.String
	obj.BusinessPubKey = businessPubKey.String
	obj.PlatformPrivKey = platformPrivKey.String
	obj.Status = strext.ToInt(isOpen.String)

	return obj, nil
}
