package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	CustDaoInstance CustDao
)

type CustDao struct {
}

func (CustDao) AddCustTx(tx *sql.Tx, accountNo, gender string) (err error, custNo string) {
	//创建运营商
	custNo = strext.NewUUID()
	sqlStr := "insert into cust(cust_no, account_no, payment_password, gender) " +
		" values ($1,$2,$3,$4)"
	err = ss_sql.ExecTx(tx, sqlStr, custNo, accountNo, "", gender)
	if err != nil {
		ss_log.Error("CustDao |AddCust err=[%v]", err)
	}
	return err, custNo
}
func (*CustDao) ModifyCustPwdByUID(uid, pwd string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	pwdMD5 := AccDaoInstance.InitPassword(pwd)

	err := ss_sql.Exec(dbHandler, `update cust set payment_password = $1 , modify_time = current_timestamp where cust_no = $2`, pwdMD5, uid)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (*CustDao) QueryPwdFromIdenNo(idenNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select payment_password from cust where cust_no=$1 limit 1", []*sql.NullString{&pwdT}, idenNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return pwdT.String
}

func (*CustDao) GetCustNoByAccNo(accNo string) (custNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select cust_no from cust where account_no = $1 and is_delete = '0' "
	var custNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&custNoT}, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return custNoT.String
}

//验证用户支付密码是否正确
func (*CustDao) CheckCustPayPWD(dbHandler *sql.DB, custNo, password string) (errR string) {
	pwdMD5 := AccDaoInstance.InitPassword(password)
	var passwordDB sql.NullString
	sqlStr := "select payment_password from cust where cust_no = $1 and is_delete = '0' "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&passwordDB}, custNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_GET
	}

	ss_log.Info("pwdMD5=[%v],passwordDB=[%v]", pwdMD5, passwordDB.String)
	if pwdMD5 != passwordDB.String {
		return ss_err.ERR_CHECK_PAY_PWD_FAILD
	}

	return ss_err.ERR_SUCCESS
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

func (*CustDao) GetAccNoByCustNo(custNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `selete account_no from cust where cust_no = $1 and is_delete = '0'' `, []*sql.NullString{&accNo}, custNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accNo.String
}

//查询支付密码是否为空
func (*CustDao) GetPaymentPwdIsNull(custNo string) (bool, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var paymentPwd sql.NullString
	sqlStr := "SELECT payment_password FROM cust WHERE cust_no=$1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&paymentPwd}, custNo)
	if err != nil {
		return false, err
	}

	if paymentPwd.String == "" {
		return false, nil
	}

	return true, nil
}

func (CustDao) UpdateTradingAuthority(custNo string, tradingAuthority int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	return ss_sql.Exec(dbHandler, `update cust set trading_authority=$1 where cust_no=$2`,
		tradingAuthority, custNo)
}
