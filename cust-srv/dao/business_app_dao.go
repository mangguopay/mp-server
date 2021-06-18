package dao

import (
	"database/sql"
	"errors"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_rsa"
	"a.a/mp-server/common/ss_sql"
)

type BusinessAppDao struct {
	AppId      string //应用id
	BusinessNo string //商户id
	//SceneNo         string //场景id
	IpwhiteList string //ip白名单(逗号分隔)
	//Rate            string //费率
	Status          string //状态 0 -未审核 1-审核通过 ，2审核未通过
	PlatformPubKey  string //平台公钥
	PlatformPrivKey string //平台私钥
	BusinessPubKey  string //商户公钥
	SignMethod      string //签名方式
	ApplyType       string //应用类型 1-移动应用，2-网页应用
	SmallImgNo      string //小图标id
	BigImgNo        string //大图标id
	AppName         string //应用名称
	Describe        string //应用描述
	//Cycle           string //结算周期(1-T+1,以此类推)
	//IndustryNo      string //经营类目id
	//SignedNo        string //签约编号
}

var BusinessAppDaoInst BusinessAppDao

func (BusinessAppDao) GetBusinessAppCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		" FROM business_app ba " +
		" left join business b on b.business_no = ba.business_no " +
		" left join account acc on acc.uid = b.account_no " + whereModel.WhereStr
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

func (BusinessAppDao) GetBusinessAppList(whereList []*model.WhereSqlCond, page, pageSize string) (datasT []*go_micro_srv_cust.BusinessAppData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	orderByStr := " ORDER BY case ba.status when " + constants.BusinessAppStatus_Pending + " then 1 end, ba.create_time desc "
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, orderByStr)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))
	sqlStr := "SELECT  ba.app_id, ba.business_no, ba.ip_white_list, ba.create_time, ba.update_time " +
		", ba.status, ba.sign_method, ba.apply_type, ba.small_img_no, ba.big_img_no" +
		", ba.app_name, ba.describe, ba.notes, acc.account " +
		" FROM business_app ba " +
		" left join business b on b.business_no = ba.business_no " +
		" left join account acc on acc.uid = b.account_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.BusinessAppData
	for rows.Next() {
		data := go_micro_srv_cust.BusinessAppData{}
		err2 = rows.Scan(
			&data.AppId,
			&data.BusinessNo,
			&data.IpWhiteList,
			&data.CreateTime,
			&data.UpdateTime,

			&data.Status,
			&data.SignMethod,
			&data.ApplyType,
			&data.SmallImgNo,
			&data.BigImgNo,

			&data.AppName,
			&data.Describe,
			&data.Notes,
			&data.Account,
		)

		if err2 != nil {
			ss_log.Error("appId[%v],err=[%v]", data.AppId, err2)
			return nil, err2
		}

		datas = append(datas, &data)
	}

	return datas, nil
}

func (BusinessAppDao) GetBusinessAppDetail(whereList []*model.WhereSqlCond) (dataT *go_micro_srv_cust.BusinessAppData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "SELECT  ba.app_id, ba.business_no, ba.ip_white_list, ba.create_time, ba.update_time" +
		", ba.fixed_qrcode, ba.status, ba.platform_pub_key, ba.business_pub_key, ba.sign_method" +
		", ba.apply_type, ba.small_img_no, ba.big_img_no, ba.app_name, ba.describe " +
		", acc.account, acc.uid " +
		" FROM business_app ba " +
		" left join business b on b.business_no = ba.business_no " +
		" left join account acc on acc.uid = b.account_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	data := &go_micro_srv_cust.BusinessAppData{}
	var platformPubKey, fixedQrCode sql.NullString
	err2 = rows.Scan(
		&data.AppId,
		&data.BusinessNo,
		&data.IpWhiteList,
		&data.CreateTime,
		&data.UpdateTime,

		&fixedQrCode,
		&data.Status,
		&platformPubKey,
		&data.BusinessPubKey,
		&data.SignMethod,

		&data.ApplyType,
		&data.SmallImgNo,
		&data.BigImgNo,
		&data.AppName,
		&data.Describe,

		&data.Account,
		&data.AccUid,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}
	data.FixedQrcode = fixedQrCode.String
	if data.FixedQrcode != "" { // 有生成固定二维码
		data.FixedQrcode = constants.GetBusinessFixedQrCodeUrl(data.FixedQrcode)
	}

	data.PlatformPubKey = platformPubKey.String

	return data, nil
}

func (b BusinessAppDao) AddBusinessApp(data BusinessAppDao) (appIdT string, err error) {
	ss_log.Info("开始插入应用 data[%+v]", data)
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	appId := strext.GetDailyId()
	sqlStr := `insert into business_app(app_id,	business_no, ip_white_list, status, business_pub_key,
	sign_method, apply_type, small_img_no, big_img_no, app_name, 
	describe, create_time, update_time)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,current_timestamp,current_timestamp) `

	if err2 := ss_sql.Exec(dbHandler, sqlStr,
		appId, data.BusinessNo, data.IpwhiteList, data.Status, data.BusinessPubKey,
		data.SignMethod, data.ApplyType, data.SmallImgNo, data.BigImgNo, data.AppName,
		data.Describe); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return "", err2
	}

	return appId, nil
}

//修改未通过的应用信息
func (b BusinessAppDao) UpdateBusinessApp(data BusinessAppDao) error {
	ss_log.Info("开始修改应用 data[%+v]", data)
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `update business_app
			set ip_white_list = $3, status = $4, business_pub_key = $5, sign_method = $6, apply_type = $7,
			small_img_no = $8, big_img_no = $9, app_name = $10, describe = $11, notes = $12, 
			update_time = current_timestamp 
			where app_id = $1 and status = $2 `

	if err2 := ss_sql.Exec(dbHandler, sqlStr,
		data.AppId, constants.BusinessAppStatus_Deny, data.IpwhiteList, constants.BusinessAppStatus_Pending, data.BusinessPubKey,
		data.SignMethod, data.ApplyType, data.SmallImgNo, data.BigImgNo, data.AppName,
		data.Describe, ""); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

func (b BusinessAppDao) UpdateBusinessAppStatusTx(tx *sql.Tx, appId, status, notes, appQrCodeId string) (err error) {
	if status == constants.BusinessAppStatus_Passed {
		//平台私钥、平台公钥
		platformPrivKey, platformPubKey, err := ss_rsa.GenRsaKeyPairPKCS1(2048)
		if err != nil {
			ss_log.Error("产生密钥对出错,err[%v]", err)
			return err
		}

		platformPrivKey = ss_rsa.StripRSAKey(platformPrivKey)
		platformPubKey = ss_rsa.StripRSAKey(platformPubKey)

		sqlStr := "update business_app set status = $3, platform_pub_key = $4, platform_priv_key = $5, " +
			" notes = $6, fixed_qrcode = $7, update_time = current_timestamp " +
			" where app_id = $1 and status = $2 "
		if err2 := ss_sql.ExecTx(tx, sqlStr, appId, constants.BusinessAppStatus_Pending, status, platformPubKey, platformPrivKey, notes, appQrCodeId); err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return err2
		}
	} else if status == constants.BusinessAppStatus_Deny {
		sqlStr := "update business_app set status = $3, notes = $4, update_time = current_timestamp where app_id = $1 and status = $2 "
		if err2 := ss_sql.ExecTx(tx, sqlStr, appId, constants.BusinessAppStatus_Pending, status, notes); err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return err2
		}
	} else if status == constants.BusinessAppStatus_Invalid {
		sqlStr := "update business_app set status = $4, notes = $5, update_time = current_timestamp where app_id = $1 and status in($2,$3) "
		if err2 := ss_sql.ExecTx(tx, sqlStr, appId, constants.BusinessAppStatus_Up, constants.BusinessAppStatus_Passed, status, notes); err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return err2
		}
	} else {
		return errors.New("参数不合法")
	}

	return nil
}

//商家修改应用上下架状态使用该接口
func (b BusinessAppDao) UpdateBusinessAppStatus(appId, oldStatus, status string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_app set status = $3, update_time = current_timestamp " +
		" where app_id = $1 and status = $2 "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, appId, oldStatus, status); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

func (b BusinessAppDao) GetBusinessAppStatus(appId string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select status from business_app where app_id = $1 "
	var status sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&status}, appId); err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return status.String
}

//确认一个应用是否是该商家的
func (b BusinessAppDao) CheckBusinessApp(appId, idenNo string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	sqlStr := "select count(1) from business_app where app_id = $1 and business_no = $2 and status != $3 "
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&count}, appId, idenNo, constants.BusinessAppStatus_Delete); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return false
	}

	return strext.ToInt(count.String) != 0
}

func (b BusinessAppDao) DelBusinessApp(appId string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_app set status = $2 where app_id = $1 and status != $2 "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, appId, constants.BusinessAppStatus_Delete); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

//修改商家应用部分信息
func (b BusinessAppDao) ModifyAppBusinessPartial(appId, pubKey, signMethod, ipWhiteList string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_app set business_pub_key = $3, sign_method = $4, ip_white_list = $5 where app_id = $1 and status != $2 "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, appId, constants.BusinessAppStatus_Delete, pubKey, signMethod, ipWhiteList); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}
