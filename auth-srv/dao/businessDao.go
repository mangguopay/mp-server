package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	BusinessDaoInstance BusinessDao
)

type BusinessDao struct {
}

func (BusinessDao) AddBusinessTx(tx *sql.Tx, accountNo, fullName, businessType string) (err error, businessNo string) {
	//创建运营商
	businessNoT := strext.NewUUID()
	businessId := strext.GetDailyId()
	sqlStr := "insert into business(business_no, account_no ,is_delete ,use_status ,full_name ,business_id , business_type, create_time) " +
		" values ($1,$2,$3,$4,$5,$6,$7,current_timestamp)"
	err = ss_sql.ExecTx(tx, sqlStr, businessNoT, accountNo, "0", "1", fullName, businessId, businessType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	return err, businessNoT
}

//修改商家的支付密码
func (*BusinessDao) ModifyBusinessPayPwdByUid(businessNo, payPwd string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	pwdMD5 := AccDaoInstance.InitPassword(payPwd)

	err := ss_sql.Exec(dbHandler, `update business set payment_password = $1 , modify_time = current_timestamp where business_no = $2 and is_delete = '0' `, pwdMD5, businessNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//修改商家的支付密码
func (*BusinessDao) ModifyBusinessPayPwdByUidTx(tx *sql.Tx, businessNo, payPwd string) error {
	pwdMD5 := AccDaoInstance.InitPassword(payPwd)

	err := ss_sql.ExecTx(tx, `update business set payment_password = $1 , modify_time = current_timestamp where business_no = $2`, pwdMD5, businessNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (*BusinessDao) GetPayPwdFromIdenNo(idenNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var payPwdDb sql.NullString
	sqlStr := "select payment_password from business where business_no = $1 and is_delete = '0' "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&payPwdDb}, idenNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return payPwdDb.String
}

func (*BusinessDao) GetAccNoByBusiness(idenNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var payPwdDb sql.NullString
	sqlStr := "select payment_password from business where business_no = $1 and is_delete = '0' "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&payPwdDb}, idenNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return payPwdDb.String
}

//确认是否已初始化支付密码(0否，1是)
func (*BusinessDao) GetBusinessInitPayPwdStatus(businessNo string) (initPayPwdStatus string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select init_pay_pwd_status from business where business_no = $1 and is_delete = '0' "
	var initPayPwdStatusT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&initPayPwdStatusT}, businessNo); err != nil {
		ss_log.Error("err=[%v]", err)
		return "", err
	}

	return initPayPwdStatusT.String, nil
}

func (*BusinessDao) ModifyBusinessInitPayPwdStatus(tx *sql.Tx, businessNo string) (err error) {
	sqlStr := "update business set init_pay_pwd_status = $2 where business_no = $1 and is_delete = '0' and init_pay_pwd_status = $3 "
	if err := ss_sql.ExecTx(tx, sqlStr, businessNo, constants.InitPayPwdStatus_true, constants.InitPayPwdStatus_false); err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	return nil
}
