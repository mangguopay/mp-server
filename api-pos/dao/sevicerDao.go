package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	ServicerDaoInstance ServicerDao
)

type ServicerDao struct {
}

func (*ServicerDao) GetAccountNoFromServiceNo(serviceNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_no from servicer where servicer_no=$1 and is_delete='0' limit 1`, []*sql.NullString{&accountNoT}, serviceNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return accountNoT.String
}

func (*ServicerDao) GetServicerNoByAccNo(dbHandler *sql.DB, accNo string) (servicerNo string) {
	sqlStr := "select servicer_no from servicer where account_no = $1 and is_delete = '0' "
	var servicerNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&servicerNoT}, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return servicerNoT.String
}

func (*ServicerDao) GetServicerPubKeyFromNo(servicerNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select pub_key from servicer where servicer_no = $1 and is_delete = '0' "
	var pubKeyT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&pubKeyT}, servicerNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return pubKeyT.String
}
