package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	CashierDaoInstance CashierDao
)

type CashierDao struct{}

// 查询手机号和查询收银员ID
func (*CashierDao) QueryPhoneAndCID(uid, rechierType string) (string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 查询手机号,店员id
	var phone, cashierUID sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT a.phone,r.iden_no  FROM account a  LEFT JOIN rela_acc_iden r ON a.uid = r.account_no  "+
		" WHERE a.uid = $1 and r.account_type = $2 and a.is_delete = 0 LIMIT 1", []*sql.NullString{&phone, &cashierUID}, uid, rechierType); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", "", err
	}
	return phone.String, cashierUID.String, nil

}

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

func (*CashierDao) DeleteCashier(dbHandler *sql.DB, uid string) (errCode string) {
	err := ss_sql.Exec(dbHandler, `update cashier set is_delete='1' where uid=$1`, uid)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS

}

func (*CashierDao) GetCashierNoByAccNo(accNo string) (cashierNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select iden_no from rela_acc_iden where account_no = $1 and account_type = $2 "
	var cashierNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cashierNoT}, accNo, constants.AccountType_POS)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return cashierNoT.String
}

func (*CashierDao) GetSrvAccNoFromCaNo(caNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT s.account_no FROM rela_acc_iden rai LEFT JOIN cashier c ON rai.iden_no = c.uid LEFT JOIN servicer s ON " +
		"s.servicer_no = c.servicer_no WHERE s.is_delete = '0' AND rai.account_no = $1"
	var srvAccNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&srvAccNo}, caNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return srvAccNo.String
}
