package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/cust-srv/common"
	"database/sql"
)

type SceneSignedDao struct {
	SignedNo      string
	Status        string
	CreateTime    string
	StartTime     string
	EndTime       string
	Rate          string
	Cycle         string
	SceneNo       string
	SceneName     string
	IndustryNo    string
	IndustryName  string
	BusinessNo    string
	BusinessAccNo string
	LastSignedNo  string
}

var SceneSignedDaoInst SceneSignedDao

/**
签约列表
*/
func (SceneSignedDao) GetList(whereStr, lang string, args []interface{}) ([]*SceneSignedDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select ss.signed_no, ss.status, ss.create_time, ss.start_time, ss.end_time, ss.rate, ss.cycle," +
		"bs.scene_name, bi.name_ch, bi.name_en, bi.name_km " +
		"from business_scene_signed ss " +
		"left join business_scene bs on bs.scene_no = ss.scene_no " +
		"left join business_industry bi on bi.code = ss.industry_no "

	sqlStr += whereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*SceneSignedDao
	for rows.Next() {
		var signedNo, status, createTime, startTime, endTime, rate, cycle, sceneName, nameCh, nameEn, nameKm sql.NullString
		err := rows.Scan(&signedNo, &status, &createTime, &startTime, &endTime, &rate, &cycle, &sceneName, &nameCh, &nameEn, &nameKm)
		if err != nil {
			return nil, err
		}

		industryName := nameCh.String
		switch lang {
		case constants.LangEnUS:
			industryName = nameEn.String
		case constants.LangKmKH:
			industryName = nameKm.String
		case constants.LangZhCN:
			industryName = nameCh.String
		}
		data := &SceneSignedDao{
			SignedNo:     signedNo.String,
			Status:       status.String,
			CreateTime:   createTime.String,
			StartTime:    startTime.String,
			EndTime:      endTime.String,
			Rate:         rate.String,
			Cycle:        cycle.String,
			SceneName:    sceneName.String,
			IndustryName: industryName,
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

/**
签约数量统计
*/
func (SceneSignedDao) Count(whereStr string, args []interface{}) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) " +
		"from business_scene_signed ss " +
		"left join business_scene bs on bs.scene_no = ss.scene_no " +
		"left join business_industry bi on bi.code = ss.industry_no "

	sqlStr += whereStr
	var total sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, args...)
	if err != nil {
		return "", err
	}
	return total.String, nil
}

/**
查询过期签约
*/
func (SceneSignedDao) GetExpireSignedList(expireTime string) ([]*SceneSignedDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, common.GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	currentTime := ss_time.Now(global.Tz).Format(ss_time.DateTimeDashFormat)
	sqlStr := "SELECT ss.signed_no, ss.business_no, ss.business_account_no, ss.end_time, ss.scene_no, ss.rate, ss.cycle, ss.industry_no " +
		"FROM business_scene_signed ss " +
		"LEFT JOIN business_scene_signed ss2 ON ss2.scene_no = ss.scene_no " +
		"LEFT JOIN business_scene bs ON bs.scene_no = ss.scene_no " +
		"LEFT JOIN business s ON s.business_no = ss.business_no " +
		"WHERE s.use_status=$1 AND ss.status=$2 AND ss.end_time > $3 AND ss.end_time <= $4 " +
		"AND ss2.last_signed_no != ss.signed_no "

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr,
		constants.BusinessUseStatusEnable, constants.SignedStatusPassed, currentTime, expireTime)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*SceneSignedDao
	for rows.Next() {
		var signedNo, businessNo, businessAccNo, signedEndTime, sceneNo, rate, cycle, industryNo sql.NullString
		if err := rows.Scan(&signedNo, &businessNo, &businessAccNo, &signedEndTime, &sceneNo, &rate, &cycle, &industryNo); err != nil {
			return nil, err
		}
		data := new(SceneSignedDao)
		data.SignedNo = signedNo.String
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

/**
续签
*/
func (SceneSignedDao) AutoSigned(d *SceneSignedDao) (signedNoT string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", common.GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	signedNo := strext.GetDailyId()
	sqlStr := "insert into business_scene_signed(signed_no, business_account_no, business_no, start_time, end_time, status, " +
		"scene_no, rate, cycle, industry_no, last_signed_no, create_time)" +
		" values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP) "
	if err2 := ss_sql.Exec(dbHandler, sqlStr,
		signedNo, d.BusinessAccNo, d.BusinessNo, d.StartTime, d.EndTime, constants.SignedStatusPassed,
		d.SceneNo, d.Rate, d.Cycle, d.IndustryNo, d.LastSignedNo,
	); err2 != nil {
		return "", err2
	}

	return signedNo, nil
}

/**
修改签约过期记录的状态
*/
func (SceneSignedDao) UpdateStatusBySignedNo(expireTime, status string) ([]string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, common.GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "UPDATE business_scene_signed SET status=$1,update_time=CURRENT_TIMESTAMP WHERE status=$2 AND end_time <= $3 RETURNING signed_no"
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
