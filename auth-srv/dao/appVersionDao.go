package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type appVersionDao struct{}

var AppVersionDaoInstance appVersionDao

func (appVersionDao) QueryAppVersion(version, vsType string) (string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var descriptionT, versionT, appURLT, vsCodeT, isForce sql.NullString
	if version == "" {
		err := ss_sql.QueryRow(dbHandler, "SELECT description,version,app_url,vs_code,is_force from app_version where vs_type = $1 and status = 1 ORDER BY create_time desc LIMIT 1",
			[]*sql.NullString{&descriptionT, &versionT, &appURLT, &vsCodeT, &isForce}, vsType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return "", "", "", "", ""
		}
	} else {
		err := ss_sql.QueryRow(dbHandler, "SELECT description,version,app_url,vs_code,is_force from app_version where version=$1 and vs_type = $2 status = 1 ",
			[]*sql.NullString{&descriptionT, &versionT, &appURLT, &vsCodeT, &isForce}, version, vsType)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return "", "", "", "", ""
		}
	}

	return descriptionT.String, versionT.String, appURLT.String, vsCodeT.String, isForce.String
}
func (appVersionDao) QueryAppVersion1(vsType, system string) (string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var descriptionT, versionT, appURLT, vsCodeT, isForce sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT description,version,app_url,vs_code,is_force from app_version where vs_type = $1 and system  = $2 and status = 1 ORDER BY create_time desc LIMIT 1",
		[]*sql.NullString{&descriptionT, &versionT, &appURLT, &vsCodeT, &isForce}, vsType, system)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", "", "", "", ""
	}

	return descriptionT.String, versionT.String, appURLT.String, vsCodeT.String, isForce.String
}
