package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type VaccountDao struct {
}

var VaccountDaoInst VaccountDao

func (VaccountDao) InitVaccountNo(accountNo, balanceType string, vaType int32) (vaccountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	vaccountNo = strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into vaccount(vaccount_no,account_no,va_type,balance,create_time,is_delete,use_status,frozen_balance,balance_type) values ($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8)`,
		vaccountNo, accountNo, vaType, "0", "0", "1", "0", balanceType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return vaccountNo
}
