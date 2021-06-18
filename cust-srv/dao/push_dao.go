package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	PushDaoInstance PushDao
)

type PushDao struct {
}

func (PushDao) GetPushConfCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM push_conf  " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (PushDao) GetPushRecordCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM push_record  " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (PushDao) GetPushTempsCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM push_temp  " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*PushDao) GetPushConfInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.PushConfData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT pusher,config,use_status,create_time,update_time,pusher_no, condition_type, condition_value " +
		" FROM push_conf  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.PushConfData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.PushConfData{}
		var updateTime sql.NullString
		err2 := rows.Scan(
			&data.Pusher,
			&data.Config,
			&data.UseStatus,
			&data.CreateTime,
			&updateTime,

			&data.PusherNo,
			&data.ConditionType,
			&data.ConditionValue,
		)
		if err2 != nil {
			ss_log.Error("push_conf表查询出错。PusherNo：[%v],err:[%v]", data.PusherNo, err2)
			continue
		}
		if updateTime.String != "" {
			data.UpdateTime = updateTime.String
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*PushDao) GetPushConfDetail(whereModelStr string, whereArgs []interface{}) (dataT *go_micro_srv_cust.PushConfData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT pusher,config,use_status,create_time,update_time,pusher_no, condition_type, condition_value " +
		" FROM push_conf  " + whereModelStr
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	data := &go_micro_srv_cust.PushConfData{}
	var updateTime sql.NullString
	err2 := rows.Scan(
		&data.Pusher,
		&data.Config,
		&data.UseStatus,
		&data.CreateTime,
		&updateTime,

		&data.PusherNo,
		&data.ConditionType,
		&data.ConditionValue,
	)
	if err2 != nil {
		ss_log.Error("push_conf表查询出错。PusherNo：[%v],err:[%v]", data.PusherNo, err2)
	}
	data.UpdateTime = updateTime.String

	return data, err
}

func (*PushDao) GetPushRecordInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.PushRecordData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT id,business,phone,content,status,create_time,push_no,temp_no,message " +
		" FROM push_record  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.PushRecordData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.PushRecordData{}
		var pushNo, tempNo, message sql.NullString
		err2 := rows.Scan(
			&data.Id,
			&data.Business,
			&data.Phone,
			&data.Content,
			&data.Status,

			&data.CreateTime,
			&pushNo,
			&tempNo,
			&message,
		)
		if err2 != nil {
			ss_log.Error("push_record表查询出错。Id：[%v],err:[%v]", data.Id, err2)
			continue
		}
		data.PushNo = pushNo.String
		data.TempNo = tempNo.String
		data.Message = message.String

		datas = append(datas, data)
	}

	return datas, err
}

func (*PushDao) GetPushTempsInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.PushTempData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT temp_no,push_nos,title_key,content_key,len_args " +
		" FROM push_temp  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.PushTempData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.PushTempData{}
		var lenArgs sql.NullString
		err2 := rows.Scan(
			&data.TempNo,
			&data.PushNos,
			&data.TitleKey,
			&data.ContentKey,
			&lenArgs,
		)
		if err2 != nil {
			ss_log.Error("push_temp表查询出错。TempNo：[%v],err:[%v]", data.TempNo, err2)
			continue
		}
		data.LenArgs = lenArgs.String

		datas = append(datas, data)
	}

	return datas, err
}

func (*PushDao) DelPushTemp(tx *sql.Tx, tempNo string) (err error) {
	sqlStr := "update push_temp set is_delete = '1' where temp_no = $1 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, tempNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*PushDao) ModifyPushConf(tx *sql.Tx, pusherNo, pusher, config, useStatus, conditionType, conditionValue string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	//condition_type
	//condition_value
	sqlStr := "update push_conf set pusher = $2, config = $3, use_status = $4, condition_type = $5, condition_value = $6 where pusher_no = $1 "
	errT := ss_sql.ExecTx(tx, sqlStr, pusherNo, pusher, config, useStatus, conditionType, conditionValue)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*PushDao) AddPushConf(tx *sql.Tx, pusher, config, useStatus, conditionType, conditionValue string) (pusherNoT string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	pusherNo := strext.NewUUID()
	sqlStr := "insert into push_conf(pusher_no, pusher, config, use_status, condition_type, condition_value, create_time) values($1,$2,$3,$4,$5,$6,current_timestamp)"
	errT := ss_sql.ExecTx(tx, sqlStr, pusherNo, pusher, config, useStatus, conditionType, conditionValue)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return pusherNo, nil
}

func (*PushDao) DelPushConf(tx *sql.Tx, pusherNo string) (err error) {
	sqlStr := "delete from push_conf where pusher_no = $1 "
	errT := ss_sql.ExecTx(tx, sqlStr, pusherNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*PushDao) ConfirmPusherIsExist(pusherNo string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var totalT sql.NullString
	sqlCnt := "SELECT count(1) FROM push_conf  where pusher_no = $1"
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, pusherNo)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return strext.ToInt(totalT.String) > 0
	}
	return strext.ToInt(totalT.String) > 0
}
