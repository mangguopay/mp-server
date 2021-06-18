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

type BulletinDao struct {
}

var BulletinDaoInst BulletinDao

func (*BulletinDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlCnt := "select count(1) " +
		" from bulletin  " + whereModel.WhereStr
	var totalT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "0"
	}

	return totalT.String
}

func (*BulletinDao) GetBulletins(whereList []*model.WhereSqlCond, page, pageSize int, extraStr string) (datas []*go_micro_srv_cust.BulletinData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	//先展示未发布的，再展示发布里面置顶的，最后商家按发布时间倒序，管理员、运营按创建时间倒序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, extraStr)
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)

	sqlStr := "select bulletin_id, title, content, use_status, create_time, bulletin_time, top_status  " +
		" from bulletin " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		return nil, err2
	}

	var datasT []*go_micro_srv_cust.BulletinData
	for rows.Next() {
		data := go_micro_srv_cust.BulletinData{}
		var bulletinTime sql.NullString
		err2 = rows.Scan(
			&data.BulletinId,
			&data.Title,
			&data.Content,
			&data.UseStatus,
			&data.CreateTime,

			&bulletinTime,
			&data.TopStatus,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			return nil, err2
		}
		data.BulletinTime = bulletinTime.String
		datasT = append(datasT, &data)
	}

	return datasT, nil
}

func (*BulletinDao) GetBulletinDetail(whereList []*model.WhereSqlCond) (*go_micro_srv_cust.BulletinData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "select bulletin_id, title, content, use_status, create_time, bulletin_time, top_status  " +
		" from bulletin " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		return nil, err2
	}

	dataT := &go_micro_srv_cust.BulletinData{}
	var bulletinTime sql.NullString
	err2 = rows.Scan(
		&dataT.BulletinId,
		&dataT.Title,
		&dataT.Content,
		&dataT.UseStatus,
		&dataT.CreateTime,

		&bulletinTime,
		&dataT.TopStatus,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}
	dataT.BulletinTime = bulletinTime.String

	return dataT, nil
}

func (*BulletinDao) AddBulletin(title, content string) (logNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.NewUUID()
	sqlCnt := "insert into bulletin(bulletin_id, title, content, use_status, top_status, create_time) " +
		" values($1,$2,$3,$4,$5,current_timestamp) "
	if err := ss_sql.Exec(dbHandler, sqlCnt, id, title, content, constants.BulletinUseStatus_UnBulletin, constants.BulletinTopStatus_False); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", err
	}

	return logNo, nil
}

func (*BulletinDao) UpdateBulletin(bulletinId, title, content string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlCnt := "update bulletin set title = $2, content = $3 where bulletin_id = $1  "
	if err := ss_sql.Exec(dbHandler, sqlCnt, bulletinId, title, content); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return err
	}

	return nil
}

func (*BulletinDao) DeleteBulletin(bulletinId string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "delete from bulletin where bulletin_id = $1 "
	if err := ss_sql.Exec(dbHandler, sqlStr, bulletinId); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return err
	}

	return nil
}

func (*BulletinDao) GetBulletinUseStatus(bulletinId string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select use_status from bulletin where bulletin_id = $1 "
	var useStatus sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&useStatus}, bulletinId); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return ""
	}

	return useStatus.String
}

func (*BulletinDao) UpdateUseStatus(bulletinId, useStatus string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update bulletin set use_status = $3, bulletin_time = current_timestamp where bulletin_id = $1 and use_status = $2 "
	if err := ss_sql.Exec(dbHandler, sqlStr, bulletinId, constants.BulletinUseStatus_UnBulletin, useStatus); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return err
	}

	return nil
}

func (*BulletinDao) UpdateTopStatus(bulletinId, topStatus string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update bulletin set top_status = $2 where bulletin_id = $1  "
	if err := ss_sql.Exec(dbHandler, sqlStr, bulletinId, topStatus); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return err
	}

	return nil
}
