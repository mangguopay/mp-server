package dao

import (
	"database/sql"
	"errors"

	"a.a/cu/strext"

	"a.a/cu/db"
	"a.a/mp-server/common/ss_sql"
)

var AppInstance App

type App struct {
	AppId              string
	AppName            string
	MerchantPrivateKey string
	MerchantPublicKey  string
	PlatformPublicKey  string
	MerchantKeyType    string
	IsUse              string
	CreateTime         string
}

const (
	IsUseOk = "1"
	IsUseNo = "0"
)

// 插入一条记录
func (a *App) Insert(app *App) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := "INSERT INTO app (id, app_id, app_name, merchant_private_key, merchant_public_key, platform_public_key, merchant_key_type, is_use, create_time)"
	sqlStr += " VALUES ($1, $2, $3, $4, $5,  $6, $7, $8, current_timestamp)"

	execErr := ss_sql.Exec(dbHandler, sqlStr,
		strext.NewUUID(),
		app.AppId,
		app.AppName,
		app.MerchantPrivateKey,
		app.MerchantPublicKey,
		app.PlatformPublicKey,
		app.MerchantKeyType,
		IsUseNo,
	)

	return execErr
}

// 获取应用列表
func (a *App) GetList(page, pageSize int) ([]App, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	offset := (page - 1) * pageSize

	sqlStr := "SELECT app_id, app_name, merchant_private_key, merchant_public_key, platform_public_key, merchant_key_type, is_use, create_time FROM app  "
	sqlStr += " ORDER BY create_time DESC LIMIT $1 OFFSET $2 "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, pageSize, offset)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		return nil, qErr
	}

	list := []App{}

	for rows.Next() {
		var appId, appName, merchantPrivateKey, merchantPublicKey, platformPublicKey, merchantKeyType, isUse, createTime sql.NullString

		err := rows.Scan(&appId, &appName, &merchantPrivateKey, &merchantPublicKey, &platformPublicKey, &merchantKeyType, &isUse, &createTime)
		if err != nil {
			return nil, err
		}
		list = append(list, App{
			AppId:              appId.String,
			AppName:            appName.String,
			MerchantPrivateKey: merchantPrivateKey.String,
			MerchantPublicKey:  merchantPublicKey.String,
			PlatformPublicKey:  platformPublicKey.String,
			MerchantKeyType:    merchantKeyType.String,
			IsUse:              isUse.String,
			CreateTime:         createTime.String,
		})
	}

	return list, nil
}

// 获取一条配置可用的应用
func (a *App) GetUsingRow() (*App, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := "SELECT app_id, app_name, merchant_private_key, platform_public_key FROM app  "
	sqlStr += " WHERE is_use=1 ORDER BY create_time DESC LIMIT 1 "

	var appId, appName, merchantPrivateKey, platformPublicKey sql.NullString

	qErr := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&appId, &appName, &merchantPrivateKey, &platformPublicKey})
	if qErr != nil {
		return nil, qErr
	}

	app := &App{
		AppId:              appId.String,
		AppName:            appName.String,
		MerchantPrivateKey: merchantPrivateKey.String,
		PlatformPublicKey:  platformPublicKey.String,
	}

	return app, nil
}

// 设置app为可用
func (a *App) SetAppUsing(appId string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	// 将其所有应用修改为不可用
	err := ss_sql.Exec(dbHandler, `UPDATE app SET is_use=$1`, IsUseNo)
	if nil != err {
		return err
	}

	// 将指定的应该修改为可用
	err2 := ss_sql.Exec(dbHandler, `UPDATE app SET is_use=$1 WHERE app_id=$2`, IsUseOk, appId)
	if nil != err2 {
		return err2
	}

	return nil
}

// 通过订单号获取一条记录
func (a *App) GetOneByAppId(appIdParam string) (*App, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := "SELECT app_id, app_name, merchant_private_key, merchant_public_key, platform_public_key, merchant_key_type, is_use, create_time FROM app  "
	sqlStr += " WHERE app_id=$1 "

	var appId, appName, merchantPrivateKey, merchantPublicKey, platformPublicKey, merchantKeyType, isUse, createTime sql.NullString

	qErr := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&appId, &appName, &merchantPrivateKey, &merchantPublicKey, &platformPublicKey, &merchantKeyType, &isUse, &createTime},
		appIdParam,
	)

	if qErr != nil {
		return nil, qErr
	}

	app := &App{
		AppId:              appId.String,
		AppName:            appName.String,
		MerchantPrivateKey: merchantPrivateKey.String,
		MerchantPublicKey:  merchantPublicKey.String,
		PlatformPublicKey:  platformPublicKey.String,
		MerchantKeyType:    merchantKeyType.String,
		IsUse:              isUse.String,
		CreateTime:         createTime.String,
	}

	return app, nil
}
