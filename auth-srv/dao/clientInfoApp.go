package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type ClientInfoAppDao struct {
	Id          string
	DeviceBrand string
	DeviceModel string
	Resolution  string
	ScreenSize  string
	Imei1       string
	Imei2       string
	SystemVer   string
	CreateTime  string
	UploadPoint string
	UserAgent   string
	Platform    string
	AppVer      string
	Account     string
	Uuid        string
}

// 插入一条记录
func (c *ClientInfoAppDao) Insert() error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "INSERT INTO client_info_app (id, device_brand, device_model, resolution, screen_size, "
	sqlStr += " imei1, imei2, system_ver, upload_point, user_agent, platform, app_ver, account, uuid, create_time) "
	sqlStr += " values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, $11, $12, $13, $14, current_timestamp)"

	execErr := ss_sql.Exec(dbHandler, sqlStr,
		strext.GetDailyId(), c.DeviceBrand, c.DeviceModel, c.Resolution, c.ScreenSize,
		c.Imei1, c.Imei2, c.SystemVer, c.UploadPoint, c.UserAgent, c.Platform, c.AppVer, c.Account, c.Uuid,
	)

	return execErr
}
