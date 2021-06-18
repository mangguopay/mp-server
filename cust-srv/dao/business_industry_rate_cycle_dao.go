package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessIndustryRateCycleDao struct {
	Id                string
	Code              string
	BusinessChannelNo string
	Rate              string
	Cycle             string
}

var BusinessIndustryRateCycleDaoInst BusinessIndustryRateCycleDao

func (BusinessIndustryRateCycleDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		" FROM business_industry_rate_cycle birc " +
		" LEFT JOIN business_industry bi ON bi.code = birc.code " +
		" LEFT JOIN business_channel bc ON bc.business_channel_no = birc.business_channel_no "
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr+whereModel.WhereStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

func (BusinessIndustryRateCycleDao) GetDatas(whereList []*model.WhereSqlCond, page, pageSize int32) (datas []*go_micro_srv_cust.BusinessIndustryRateCycleData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY birc.code ASC, birc.create_time DESC  ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)

	sqlStr := "SELECT birc.id, birc.code,  birc.business_channel_no, birc.rate, birc.cycle," +
		" birc.create_time, birc.modify_time, bi.up_code, bi.name_ch, bc.channel_name " +
		" FROM business_industry_rate_cycle birc " +
		" LEFT JOIN business_industry bi ON bi.code = birc.code " +
		" LEFT JOIN business_channel bc ON bc.business_channel_no = birc.business_channel_no "
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	for rows.Next() {
		data := go_micro_srv_cust.BusinessIndustryRateCycleData{}
		var modifyTime sql.NullString
		err2 = rows.Scan(
			&data.Id,
			&data.Code,
			&data.BusinessChannelNo,
			&data.Rate,
			&data.Cycle,

			&data.CreateTime,
			&modifyTime,
			&data.UpCode,
			&data.IndustryName,
			&data.ChannelName,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			continue
		}
		data.ModifyTime = modifyTime.String

		datas = append(datas, &data)
	}

	return datas, nil
}

func (BusinessIndustryRateCycleDao) GetDetail(whereList []*model.WhereSqlCond) (*go_micro_srv_cust.BusinessIndustryRateCycleData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "SELECT birc.id, birc.code,  birc.business_channel_no, birc.rate, birc.cycle," +
		" birc.create_time, birc.modify_time, bi.up_code, bi.name_ch, bc.channel_name " +
		" FROM business_industry_rate_cycle birc " +
		" LEFT JOIN business_industry bi ON bi.code = birc.code " +
		" LEFT JOIN business_channel bc ON bc.business_channel_no = birc.business_channel_no "
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	data := &go_micro_srv_cust.BusinessIndustryRateCycleData{}
	var modifyTime sql.NullString
	err2 = rows.Scan(
		&data.Id,
		&data.Code,
		&data.BusinessChannelNo,
		&data.Rate,
		&data.Cycle,

		&data.CreateTime,
		&modifyTime,
		&data.UpCode,
		&data.IndustryName,
		&data.ChannelName,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}
	data.ModifyTime = modifyTime.String

	return data, nil
}

func (BusinessIndustryRateCycleDao) Add(data BusinessIndustryRateCycleDao) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()
	sqlStr := "INSERT INTO business_industry_rate_cycle(id,code,business_channel_no,rate,cycle,create_time) " +
		" values($1,$2,$3,$4,$5,current_timestamp) "
	if err := ss_sql.Exec(dbHandler, sqlStr, id, data.Code, data.BusinessChannelNo, data.Rate, data.Cycle); err != nil {
		ss_log.Error("err=[%v]", err)
		return "", err
	}

	return id, nil
}

func (BusinessIndustryRateCycleDao) Update(data BusinessIndustryRateCycleDao) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "UPDATE business_industry_rate_cycle set rate = $2, cycle = $3," +
		" modify_time =  current_timestamp " +
		" WHERE id = $1 and is_delete = '0' "
	if err := ss_sql.Exec(dbHandler, sqlStr, data.Id, data.Rate, data.Cycle); err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	return nil
}

func (BusinessIndustryRateCycleDao) Delete(id string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "UPDATE business_industry_rate_cycle SET is_delete = '1' " +
		" WHERE id = $1 AND is_delete = '0' "
	if err := ss_sql.Exec(dbHandler, sqlStr, id); err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	return nil
}

func (BusinessIndustryRateCycleDao) CheckUnique(id, code, businessChannelNo string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := ""
	var cnt sql.NullString

	if id == "" {
		sqlStr = " SELECT COUNT(1)" +
			" FROM business_industry_rate_cycle " +
			" WHERE code = $1 AND business_channel_no = $2 AND is_delete = '0' "
		if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, code, businessChannelNo); err != nil {
			ss_log.Error("err=[%v]", err)
			return false
		}
	} else { //如果是有id ,有可能是自己本身
		sqlStr = " SELECT COUNT(1)" +
			" FROM business_industry_rate_cycle " +
			" WHERE id != $1 AND code = $2 AND business_channel_no = $3 AND is_delete = '0' "
		if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, id, code, businessChannelNo); err != nil {
			ss_log.Error("err=[%v]", err)
			return false
		}
	}

	return cnt.String == "0"
}
