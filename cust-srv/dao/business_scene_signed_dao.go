package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"errors"
	"time"
)

type BusinessSceneSignedDao struct {
}

type BusinessSceneSignedData struct {
	SignedNo         string
	BusinessNo       string
	BusinessAccNo    string
	Account          string //商家账号
	StartTime        string
	EndTime          string
	SceneNo          string
	SceneName        string
	Rate             string
	Cycle            string
	IndustryNo       string //行业id
	ParentIndustryNo string //行业的父id
	IndustryNameCh   string
	IndustryNameEn   string
	IndustryNameKm   string
	LastSignedNo     string
	Status           string
	CreateTime       string
	Notes            string
}

var BusinessSceneSignedDaoInst BusinessSceneSignedDao

func (b BusinessSceneSignedDao) GetBusinessSceneSignedDetail(signedNo string) (data *BusinessSceneSignedData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "signed_no", Val: signedNo, EqType: "="},
	})

	sqlStr := `select business_no, scene_no, rate, cycle   
				from business_scene_signed `
	sqlStr += whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("SignedNo[%v],err=[%v]", signedNo, err2)
		return nil, err2
	}

	var businessNo, sceneNo, rate, cycle sql.NullString
	err2 = rows.Scan(
		&businessNo,
		&sceneNo,
		&rate,
		&cycle,
	)
	if err2 != nil {
		ss_log.Error("SignedNo[%v],err=[%v]", signedNo, err2)
		return nil, err
	}

	dataT := &BusinessSceneSignedData{
		BusinessNo: businessNo.String,
		SceneNo:    sceneNo.String,
		Rate:       rate.String,
		Cycle:      cycle.String,
	}

	return dataT, nil

}

//确认该商家的一种产品同种状态的签约记录只有一条
func (b BusinessSceneSignedDao) CheckSceneSignedUnique(businessNo, sceneNo, status string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select count(1) 
			from business_scene_signed
			where business_no = $1 
				and scene_no = $2
				and status = $3 `
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, businessNo, sceneNo, status); err != nil {
		ss_log.Error("err=[%v]", err)
		if err == sql.ErrNoRows {
			return true
		}

		return false
	}

	return cnt.String == "0"

}

func (b BusinessSceneSignedDao) AddBusinessSigned(accountUid, businessNo, sceneNo, industryNo, rate, cycle string) (signedNoT string, err error) {
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

	sqlStr := `insert into business_scene_signed(signed_no, business_account_no, business_no, create_time,
			start_time, end_time, scene_no, rate, cycle, industry_no) 
			values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) `
	if err2 := ss_sql.Exec(dbHandler, sqlStr, signedNo, accountUid, businessNo, startTime,
		startTime, endTime, sceneNo, rate, cycle, industryNo); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return "", err2
	}

	return signedNo, nil
}

func (b BusinessSceneSignedDao) AddBusinessSignedTx(tx *sql.Tx, accountUid, businessNo, sceneNo, industryNo, rate, cycle, status string) (signedNoT string, err error) {
	//获取服务期限(单位:天)
	serviceTerm := GlobalParamDaoInstance.QeuryParamValue(constants.GlobalParamKeyAppSignedTerm)
	if serviceTerm == "" {
		ss_log.Error("获取全局参数失败，paramKey=%v", constants.GlobalParamKeyAppSignedTerm)
		return "", errors.New("获取服务期限失败")
	}

	signedNo := strext.GetDailyId()
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, strext.ToInt(serviceTerm))

	sqlStr := `insert into business_scene_signed(signed_no, business_account_no, business_no, create_time,
			start_time, end_time, scene_no, rate, cycle, industry_no, status) 
			values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) `
	if err2 := ss_sql.ExecTx(tx, sqlStr, signedNo, accountUid, businessNo, startTime,
		startTime, endTime, sceneNo, rate, cycle, industryNo, status); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return "", err2
	}

	return signedNo, nil
}

func (b BusinessSceneSignedDao) UpdateStatusTx(tx *sql.Tx, signedNo, oldStatus, status, notes string) (err error) {
	sqlStr := `update business_scene_signed set status = $3, notes = $4 where signed_no = $1 and status = $2 `
	if err2 := ss_sql.ExecTx(tx, sqlStr, signedNo, oldStatus, status, notes); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

//作废正使用的产品签约（如果通过一条审核的签约，并且原来有通过审核的签约的话）
func (b BusinessSceneSignedDao) SetStatusInvalidTx(tx *sql.Tx, businessNo, sceneNo string) (err error) {
	sqlStr := `update business_scene_signed set status = $4 where business_no = $1 and scene_no = $2 and status = $3 `
	if err2 := ss_sql.ExecTx(tx, sqlStr, businessNo, sceneNo, constants.SignedStatusPassed, constants.SignedStatusInvalid); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

//修改商家产品签约的结算周期、费率(通过的才允许修改)
func (b BusinessSceneSignedDao) UpdateInfoTx(tx *sql.Tx, signedNo, cycle, rate string) (err error) {
	sqlStr := "update business_scene_signed set cycle = $3, rate = $4, update_time = current_timestamp " +
		" where signed_no = $1 and status = $2 "

	if err2 := ss_sql.ExecTx(tx, sqlStr, signedNo, constants.SignedStatusPassed, cycle, rate); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

func (BusinessSceneSignedDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		" FROM business_scene_signed bss " +
		" LEFT JOIN business_industry bi ON bi.code = bss.industry_no " +
		" LEFT JOIN business_scene bs ON bs.scene_no = bss.scene_no " +
		" LEFT JOIN account acc ON acc.uid = bss.business_account_no " + whereModel.WhereStr
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

func (BusinessSceneSignedDao) GetList(whereList []*model.WhereSqlCond, page, pageSize int32) (datasT []*BusinessSceneSignedData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY case bss.status when "+constants.SignedStatusPending+" then 1 end, bss.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)
	sqlStr := "SELECT bss.signed_no, bs.scene_name, bi.name_ch, bi.name_en, bi.name_km, bi.up_code, bss.cycle, bss.rate " +
		", acc.account, bss.industry_no, bss.scene_no, bss.status, bss.create_time, bss.start_time, bss.end_time, bss.notes " +
		" FROM business_scene_signed bss " +
		" LEFT JOIN business_industry bi ON bi.code = bss.industry_no " +
		" LEFT JOIN business_scene bs ON bs.scene_no = bss.scene_no " +
		" LEFT JOIN account acc ON acc.uid = bss.business_account_no " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*BusinessSceneSignedData
	for rows.Next() {
		data := BusinessSceneSignedData{}
		err2 = rows.Scan(
			&data.SignedNo,
			&data.SceneName,
			&data.IndustryNameCh,
			&data.IndustryNameEn,
			&data.IndustryNameKm,

			&data.ParentIndustryNo,
			&data.Cycle,
			&data.Rate,
			&data.Account,
			&data.IndustryNo,

			&data.SceneNo,
			&data.Status,
			&data.CreateTime,
			&data.StartTime,
			&data.EndTime,

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
