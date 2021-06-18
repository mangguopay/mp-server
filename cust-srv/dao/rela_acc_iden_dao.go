package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type RelaAccIdenDao struct {
}

var RelaAccIdenDaoInst RelaAccIdenDao

func (RelaAccIdenDao) InsertRelaAccIden(tx *sql.Tx, accountNo, idenNo, accountType string) (retCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.ExecTx(tx, "insert into rela_acc_iden(account_no,account_type,iden_no) VALUES ($1,$2,$3)", accountNo, accountType, idenNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_ADD
	}
	return ss_err.ERR_SUCCESS
}

func (RelaAccIdenDao) GetIdenFromAcc(accNo, accountType string) (idenNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var idenNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT iden_no FROM rela_acc_iden WHERE account_no=$1 and account_type = $2 LIMIT 1",
		[]*sql.NullString{&idenNoT}, accNo, accountType)
	if err != nil || idenNoT.String == "" {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return idenNoT.String
}

// 收银员还是服务商?
func (RelaAccIdenDao) IsCashierOrServicer(accNo string) (accType string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT account_type FROM rela_acc_iden WHERE account_no=$1 and account_type in ($2,$3)",
		[]*sql.NullString{&accTypeT}, accNo, constants.AccountType_SERVICER, constants.AccountType_POS)
	if err != nil || accTypeT.String == "" {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return accTypeT.String
}
func (RelaAccIdenDao) GetAccNo(idenNo, accType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT account_no FROM rela_acc_iden WHERE iden_no=$1 and account_type=$2",
		[]*sql.NullString{&accNoT}, idenNo, accType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return accNoT.String
}

func (RelaAccIdenDao) DeleteRelaAccIden(tx *sql.Tx, idenNo, accountType string) (err error) {
	sqlStr := " delete from rela_acc_iden where iden_no = $1 and account_type = $2 "
	errT := ss_sql.ExecTx(tx, sqlStr, idenNo, accountType)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}
