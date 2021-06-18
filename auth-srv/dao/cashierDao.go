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
	CashierDaoInstance CashierDao
)

type CashierDao struct{}

// 修改支付密码
func (*CashierDao) ModifyCashierPwdByUID(uid, pwd string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	pwdMD5 := AccDaoInstance.InitPassword(pwd)

	err := ss_sql.Exec(dbHandler, `update cashier set op_password = $1 , modify_time = current_timestamp where uid = $2 and is_delete = 0`, pwdMD5, uid)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

// 获取收银员密码
func (CashierDao) GetCashierPwdFromIdenNo(idenNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var opPWD sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select op_password from cashier where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&opPWD}, idenNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return opPWD.String
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

// 查询手机号和查询 账号绑定的类型的身份UID
func (*CashierDao) QueryPhoneAndCID(uid, rechierType string) (string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 查询手机号,关系id
	var phone, idenUID sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT a.phone,r.iden_no  FROM account a  LEFT JOIN rela_acc_iden r ON a.uid = r.account_no  "+
		" WHERE a.uid = $1 and r.account_type = $2 and a.is_delete = 0 LIMIT 1", []*sql.NullString{&phone, &idenUID}, uid, rechierType); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", "", err
	}
	return phone.String, idenUID.String, nil

}

// 查询手机号和查询 账号绑定的类型的身份UID
func (*CashierDao) QueryBusinessPhoneAndCID(uid, rechierType string) (string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 查询手机号,关系id
	var businessPhone, idenUID sql.NullString
	if err := ss_sql.QueryRow(dbHandler, "SELECT a.business_phone,r.iden_no  FROM account a  LEFT JOIN rela_acc_iden r ON a.uid = r.account_no  "+
		" WHERE a.uid = $1 and r.account_type = $2 and a.is_delete = 0 LIMIT 1", []*sql.NullString{&businessPhone, &idenUID}, uid, rechierType); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", "", err
	}
	return businessPhone.String, idenUID.String, nil

}

func (*CashierDao) DeleteCashier(dbHandler *sql.DB, uid string) (errCode string) {
	err := ss_sql.Exec(dbHandler, `update cashier set is_delete='1' where uid=$1`, uid)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS

}

//验证服务商支付密码是否正确
func (*CashierDao) CheckCashierPayPWD(dbHandler *sql.DB, cashierNo, password string) (errR string) {
	pwdMD5 := AccDaoInstance.InitPassword(password)
	var passwordDB sql.NullString
	sqlStr := "select op_password from cashier where uid = $1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&passwordDB}, cashierNo)
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

//根据操作员id查询账号uid
func (*CashierDao) GetAccNoByCaNo(caNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT account_no FROM rela_acc_iden WHERE iden_no = $1 and account_type = $2 "
	var accNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accNo}, caNo, constants.AccountType_POS)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return accNo.String

}

func (r *CashierDao) AddCashier(tx *sql.Tx, servicerNo, accountNo string) (string, error) {
	uid := strext.NewUUID()
	insertSql := "INSERT INTO cashier(uid, servicer_no, account_no, create_time)" +
		" VALUES ($1,$2,$3,current_timestamp)"
	err := ss_sql.ExecTx(tx, insertSql, uid, servicerNo, accountNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	return uid, nil
}

//确认一个服务商没有添加过该店员
func (*CashierDao) CheckServicerNoCashier(cashierAccNo, servicerNo string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from cashier where account_no = $1 and servicer_no = $2 and is_delete = '0' "
	var cnt sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, cashierAccNo, servicerNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return true
	}

	return cnt.String != "0" //返回true则说明是添加过该店员
}
