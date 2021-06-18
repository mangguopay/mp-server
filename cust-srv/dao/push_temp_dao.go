package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type PushTempDao struct{}

var PushTempDaoInst PushTempDao

func (*PushTempDao) GetPushTempInfoFromNo(tempNo string) (*go_micro_srv_cust.PushTempData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var tempNoT, pushNos, titleKey, contentKey, lenArgs sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT temp_no,push_nos,title_key,content_key,len_args FROM push_temp WHERE temp_no=$1 ",
		[]*sql.NullString{&tempNoT, &pushNos, &titleKey, &contentKey, &lenArgs}, tempNo)
	if err != nil {
		return nil, err
	}

	return &go_micro_srv_cust.PushTempData{
		TempNo:     tempNo,
		PushNos:    pushNos.String,
		TitleKey:   titleKey.String,
		ContentKey: contentKey.String,
		LenArgs:    lenArgs.String,
	}, nil

}

func (*PushTempDao) Insert(tempNo, pushNos, titleKey, contentKey string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "insert into push_temp(temp_no, push_nos, title_key, content_key,create_time) values($1,$2,$3,$4,current_timestamp)"
	err := ss_sql.Exec(dbHandler, sqlStr, tempNo, pushNos, titleKey, contentKey)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD
	}

	return ss_err.ERR_SUCCESS
}

func (*PushTempDao) ModifyPushTemp(tempNo, pushNos, titleKey, contentKey string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update push_temp set push_nos = $2, title_key = $3, content_key = $4,modify_time = current_timestamp where temp_no = $1 "
	err := ss_sql.Exec(dbHandler, sqlStr, tempNo, pushNos, titleKey, contentKey)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}

	return ss_err.ERR_SUCCESS
}

func (*PushTempDao) CheckPushTempNo(tempNo string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from push_temp where temp_no = $1 and is_delete = '0' "
	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&count}, tempNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	return count.String == "0"
}
