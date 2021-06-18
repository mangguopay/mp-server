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
	LogDaoInstance LogDao
)

type LogDao struct {
}

//此接口记APP、POS的操作日志
func (*LogDao) InsertAccountLog(description, accountUID, accountType string, typeInt int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, "insert into log_account(log_no,description,account_uid,type,account_type,log_time) values ($1,$2,$3,$4,$5,current_timestamp)",
		strext.GetDailyId(), description, accountUID, typeInt, accountType)
}

//此表只记WEB,即后台管理系统的关键操作日志
func (*LogDao) InsertWebAccountLog(description, accountUID string, typeStr string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, "insert into admin_log_account_web(log_no,description,account_uid,type,create_time) values ($1,$2,$3,$4,current_timestamp)",
		strext.GetDailyId(), description, accountUID, typeStr)
}

//此表只记WEB,即后台管理系统的关键操作日志
func (*LogDao) InsertWebAccountLogTx(tx *sql.Tx, description, accountUID string, typeStr string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.ExecTx(tx, "insert into admin_log_account_web(log_no,description,account_uid,type,create_time) values ($1,$2,$3,$4,current_timestamp)",
		strext.GetDailyId(), description, accountUID, typeStr)
}

func (*LogDao) GetLogAccountWebCnt(whereModelStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from admin_log_account_web web " +
		" left join admin_account acc on acc.uid = web.account_uid " + whereModelStr
	var totalT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereArgs...)

	return totalT.String, errT
}

func (*LogDao) GetLogAccountWebInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.LogAccountData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT web.log_no, web.description, web.account_uid, web.create_time, web.type, acc.account " +
		" FROM admin_log_account_web web " +
		" left join admin_account acc on acc.uid = web.account_uid " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.LogAccountData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.LogAccountData{}
		var account sql.NullString
		err2 := rows.Scan(
			&data.LogNo,
			&data.Description,
			&data.AccountUid,
			&data.CreateTime,
			&data.Type,

			&account,
		)
		if err2 != nil {
			ss_log.Error("admin_log_account_web表查询出错。LogNo：[%v],err:[%v]", data.LogNo, err2)
			continue
		}

		data.Account = account.String
		datas = append(datas, data)
	}

	return datas, err
}

func (*LogDao) GetLogAccountCnt(whereModelStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) " +
		" from log_account la " +
		" left join account acc on acc.uid = la.account_uid  " + whereModelStr
	var totalT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereArgs...)

	return totalT.String, errT
}

func (*LogDao) GetLogAccountInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.LogAccountData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT la.log_no, la.description, la.account_uid, la.log_time, la.type, la.account_type, acc.account " +
		" FROM log_account la " +
		" left join account acc on acc.uid = la.account_uid " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.LogAccountData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.LogAccountData{}
		var account, accountType sql.NullString
		err2 := rows.Scan(
			&data.LogNo,
			&data.Description,
			&data.AccountUid,
			&data.CreateTime,
			&data.Type,

			&accountType,
			&account,
		)
		if err2 != nil {
			ss_log.Error("log_account表查询出错。LogNo：[%v],err:[%v]", data.LogNo, err2)
			continue
		}
		data.Account = account.String
		data.AccountType = accountType.String
		datas = append(datas, data)
	}

	return datas, err
}
