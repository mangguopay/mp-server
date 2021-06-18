package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	CustDaoInstance CustDao
)

type CustDao struct {
}

func (*CustDao) QueryPwdFromOpAccNo(opAccNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select payment_password from cust where cust_no=$1 limit 1", []*sql.NullString{&pwdT}, opAccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return pwdT.String
}

func (*CustDao) GetCustNoByAccNo(dbHandler *sql.DB, accNo string) (custNo string) {
	sqlStr := "select cust_no from cust where account_no = $1 and is_delete = '0' "
	var custNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&custNoT}, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return custNoT.String
}

func (*CustDao) GetCustPubKeyFromNo(custNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select pub_key from cust where cust_no = $1 and is_delete = '0' "
	var pubKeyT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&pubKeyT}, custNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return pubKeyT.String
}

func (*CustDao) UpdateAccountPub(accNo, pub string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	err := ss_sql.Exec(dbHandler, `update cust set pub_key=$2 where account_no=$1 and is_delete = 0`, accNo, pub)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}
