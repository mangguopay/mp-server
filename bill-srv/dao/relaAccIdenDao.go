package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type RelaAccIdenDao struct {
}

var RelaAccIdenDaoInst RelaAccIdenDao

func (*RelaAccIdenDao) GetIdenNo(accountNo, accountType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var idenNOT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select iden_no from rela_acc_iden where account_no=$1 and account_type=$2 limit 1`,
		[]*sql.NullString{&idenNOT}, accountNo, accountType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return idenNOT.String
}
func (*RelaAccIdenDao) GetAccNo(idenNo, accountType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_no from rela_acc_iden where iden_no=$1 and account_type=$2 limit 1`,
		[]*sql.NullString{&accountNoT}, idenNo, accountType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return accountNoT.String
}
