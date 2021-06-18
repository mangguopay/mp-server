package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/common"
	"database/sql"
	"errors"
	"time"
)

type BusinessSignedDao struct {
	SignedNo      string
	AppId         string
	BusinessNo    string
	BusinessAccNo string
	StartTime     string
	EndTime       string
	SceneNo       string
	Rate          string
	Cycle         string
	IndustryNo    string
	LastSignedNo  string
}

var BusinessSignedDaoInst BusinessSignedDao

func (BusinessSignedDao) GetBusinessSignedCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		" FROM business_signed bs " +
		" LEFT JOIN business_industry bi ON bi.code = bs.industry_no " +
		" LEFT JOIN business_scene bsn ON bsn.scene_no = bs.scene_no " +
		" LEFT JOIN business_app app ON app.app_id = bs.app_id " +
		" LEFT JOIN account acc ON acc.uid = bs.business_account_no " + whereModel.WhereStr
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

func (BusinessSignedDao) GetBusinessSignedList(whereList []*model.WhereSqlCond, page, pageSize string) (datasT []*go_micro_srv_cust.BusinessSignedData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY case bs.status when "+constants.SignedStatusPending+" then 1 end, bs.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))
	sqlStr := "SELECT bs.signed_no, bsn.scene_name, bi.name_ch, bi.name_en, bi.up_code, bs.cycle, bs.rate, bs.app_id" +
		", acc.account, bs.industry_no, bs.scene_no, bs.status, bs.create_time, bs.start_time, bs.end_time, app.app_name, bs.notes " +
		" FROM business_signed bs " +
		" LEFT JOIN business_industry bi ON bi.code = bs.industry_no " +
		" LEFT JOIN business_scene bsn ON bsn.scene_no = bs.scene_no " +
		" LEFT JOIN business_app app ON app.app_id = bs.app_id " +
		" LEFT JOIN account acc ON acc.uid = bs.business_account_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.BusinessSignedData
	for rows.Next() {
		data := go_micro_srv_cust.BusinessSignedData{}
		err2 = rows.Scan(
			&data.SignedNo,
			&data.SceneName,
			&data.IndustryNameCh,
			&data.IndustryNameEn,
			&data.ParentIndustryNo,

			&data.Cycle,
			&data.Rate,
			&data.AppId,
			&data.Account,
			&data.IndustryNo,
			&data.SceneNo,

			&data.Status,
			&data.CreateTime,
			&data.StartTime,
			&data.EndTime,
			&data.AppName,
			&data.Notes,
		)

		if err2 != nil {
			ss_log.Error("SignedNo[%v],err=[%v]", data.SignedNo, err2)
			return nil, err
		}

		datas = append(datas, &data)
	}

	return datas, nil
}

func (b BusinessSignedDao) AddBusinessSigned(appId, accountUid, businessNo, sceneNo, rate, cycle, industryNo string) (signedNoT string, err error) {
	//获取服务期限(单位:天)
	serviceTerm := GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyAppSignedTerm)
	if serviceTerm == "" {
		ss_log.Error("获取全局参数失败，paramKey=%v", constants.GlobalParamKeyAppSignedTerm)
		return "", errors.New("获取服务期限失败")
	}

	signedNo := strext.GetDailyId()
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, strext.ToInt(serviceTerm))

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `insert into business_signed(signed_no, app_id, business_account_no, business_no, create_time,
			start_time, end_time, scene_no, rate, cycle, industry_no) 
			values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) `
	if err2 := ss_sql.Exec(dbHandler, sqlStr, signedNo, appId, accountUid, businessNo, startTime,
		startTime, endTime, sceneNo, rate, cycle, industryNo); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return "", err2
	}

	return signedNo, nil
}

//确认一个签约是否是该商家的
func (b BusinessSignedDao) CheckBusinessSigned(signedNo, idenNo string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	sqlStr := "select count(1) from business_signed where signed_no = $1 and business_no = $2 "
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&count}, signedNo, idenNo); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return false
	}

	return strext.ToInt(count.String) != 0
}

//一个应用可有多个产品,但产品只能出现签一次
func (b BusinessSignedDao) CheckAppAndSceneUnique(appId, sceneNo string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	sqlStr := "select count(1) from business_signed where app_id = $1 and scene_no = $2 and status in( $3 , $4) "
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&count}, appId, sceneNo, constants.SignedStatusPending, constants.SignedStatusPassed); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return false
	}

	return strext.ToInt(count.String) != 0
}

//查询签约快过期应用
func (b BusinessSignedDao) GetExpireSignedList(expireTime string) ([]*BusinessSignedDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, common.GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	currentTime := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	sqlStr := "SELECT s.signed_no, s.app_id,s.business_no,s.business_account_no,s.end_time, s.scene_no, s.rate, s.cycle, s.industry_no " +
		"FROM business_signed s  " +
		"LEFT JOIN business_app app ON app.app_id = s.app_id " +
		"WHERE app.status in ($1, $2) AND s.status=$3 AND s.end_time > $4 AND s.end_time <= $5  "

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr,
		constants.BusinessAppStatus_Passed, constants.BusinessAppStatus_Up, constants.SignedStatusPassed, currentTime, expireTime)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*BusinessSignedDao
	for rows.Next() {
		var signedNo, appId, businessNo, businessAccNo, signedEndTime, sceneNo, rate, cycle, industryNo sql.NullString
		if err := rows.Scan(&signedNo, &appId, &businessNo, &businessAccNo, &signedEndTime, &sceneNo, &rate, &cycle, &industryNo); err != nil {
			return nil, err
		}
		data := new(BusinessSignedDao)
		data.SignedNo = signedNo.String
		data.AppId = appId.String
		data.BusinessNo = businessNo.String
		data.BusinessAccNo = businessAccNo.String
		data.EndTime = signedEndTime.String
		data.SceneNo = sceneNo.String
		data.Rate = rate.String
		data.Cycle = cycle.String
		data.IndustryNo = industryNo.String
		dataList = append(dataList, data)
	}
	return dataList, nil
}

//续签
func (b BusinessSignedDao) AutoSigned(d *BusinessSignedDao) (signedNoT string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", common.GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	signedNo := strext.GetDailyId()
	sqlStr := "insert into business_signed(signed_no, app_id, business_account_no, business_no, start_time, end_time, status, " +
		"scene_no, rate, cycle, industry_no, last_signed_no, create_time)" +
		" values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP) "
	if err2 := ss_sql.Exec(dbHandler, sqlStr,
		signedNo, d.AppId, d.BusinessAccNo, d.BusinessNo, d.StartTime, d.EndTime, constants.SignedStatusPassed,
		d.SceneNo, d.Rate, d.Cycle, d.IndustryNo, d.LastSignedNo,
	); err2 != nil {
		return "", err2
	}

	return signedNo, nil
}

//修改签约过期记录的状态
func (b BusinessSignedDao) UpdateStatusBySignedNo(expireTime, status string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, common.GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "UPDATE business_signed SET status=$1,update_time=CURRENT_TIMESTAMP WHERE status=$2 AND end_time <= $3 RETURNING signed_no"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, status, constants.SignedStatusPassed, expireTime)
	if err != nil {
		return nil, err
	}

	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var idArr []string
	for rows.Next() {
		var id sql.NullString
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		idArr = append(idArr, id.String)
	}
	return idArr, nil
}

//修改商家应用签约产品的结算周期、费率(通过的才允许修改)
func (b BusinessSignedDao) UpdateInfo(signedNo, cycle, rate string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_signed set cycle = $3, rate = $4, update_time = current_timestamp " +
		" where signed_no = $1 and status = $2 "

	if err2 := ss_sql.Exec(dbHandler, sqlStr, signedNo, constants.SignedStatusPassed, cycle, rate); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

//修改商家应用签约状态（申请中的变通过或驳回）
func (b BusinessSignedDao) UpdateStatus(signedNo, status, notes string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_signed set status = $3, notes = $4, update_time = current_timestamp " +
		" where signed_no = $1 and status = $2  "

	if err2 := ss_sql.Exec(dbHandler, sqlStr, signedNo, constants.SignedStatusPending, status, notes); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}
