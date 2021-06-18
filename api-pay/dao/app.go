package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var (
	AppDaoInstance AppDao
)

type AppDao struct {
	AppId           string
	BusinessNo      string
	IpWhiteList     string
	Status          string
	SignMethod      string
	BusinessPubKey  string
	PlatformPubKey  string
	PlatformPrivKey string
}

func (a *AppDao) GetSignInfo(app_id string) (*AppDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var appId, signMethod, ipWhiteList, businessPubKey, platformPrivKey, status sql.NullString

	sqlStr := `SELECT app_id, sign_method, ip_white_list, business_pub_key, platform_priv_key, status FROM business_app WHERE app_id=$1 LIMIT 1`
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&appId, &signMethod, &ipWhiteList, &businessPubKey, &platformPrivKey, &status}, app_id)
	if err != nil {
		return nil, err
	}

	obj := new(AppDao)
	obj.AppId = appId.String
	obj.SignMethod = signMethod.String
	obj.IpWhiteList = ipWhiteList.String
	obj.BusinessPubKey = businessPubKey.String
	obj.PlatformPrivKey = platformPrivKey.String
	obj.Status = status.String

	return obj, nil
}
