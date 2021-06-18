package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type CashierDao struct {
}

var CashierDaoInst CashierDao

func (CashierDao) GetServicerNoFromOpAccNo(opAccNo string) (servierNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var serviceNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select servicer_no from cashier where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&serviceNoT}, opAccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return serviceNoT.String
}

// 获取收银员密码
func (CashierDao) GetCashierPwdFromOpAccNo(opAccNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var opPWD sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select op_password from cashier where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&opPWD}, opAccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return opPWD.String
}

func (CashierDao) GetSerAccNoByCashierAccNo(dbHandler *sql.DB, cashierAccNo string) (serAccNo string) {
	sqlStr := "select ser.account_no " +
		" from rela_acc_iden rai " +
		" LEFT JOIN cashier ca ON ca.uid = rai.iden_no " +
		" LEFT JOIN servicer ser ON ser.servicer_no = ca.servicer_no " +
		" where rai.account_no = $1 and rai.account_type = $2 "
	var serAccNoT sql.NullString

	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&serAccNoT}, cashierAccNo, constants.AccountType_POS)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return serAccNoT.String
}
