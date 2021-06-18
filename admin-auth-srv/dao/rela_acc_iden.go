package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type RelaAccIdenDao struct {
}

var RelaAccIdenDaoInst RelaAccIdenDao
var errPG = `pq: duplicate key value violates unique constraint "rela_acc_iden_pkey"`

func (RelaAccIdenDao) InsertRelaAccIden(tx *sql.Tx, accountNo, idenNo, accountType string) (retCode string) {
	err := ss_sql.ExecTx(tx, "insert into rela_acc_iden (account_no,account_type,iden_no) VALUES ($1,$2,$3)", accountNo, accountType, idenNo)
	if err != nil {
		if err.Error() == errPG {
			return ss_err.ERR_ACCOUNT_IS_RELA
		}
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
