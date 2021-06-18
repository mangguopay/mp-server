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
	ServicerDaoInstance ServicerDao
)

type ServicerDao struct {
}

func (ServicerDao) InsertInitService(tx *sql.Tx, accountNo string) (string, error) {
	//创建运营商
	servicer_no := strext.NewUUID()
	sqlStr := "insert into servicer(servicer_no, account_no, create_time) " +
		" values ($1,$2,current_timestamp)"
	return servicer_no, ss_sql.ExecTx(tx, sqlStr, servicer_no, accountNo)
}

// 修改支付密码
func (*ServicerDao) ModifyServicerPwdByUID(uid, pwd string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	pwdMD5 := AccDaoInstance.InitPassword(pwd)

	err := ss_sql.Exec(dbHandler, `update servicer set password = $1 , modify_time = current_timestamp where servicer_no = $2 and is_delete = 0`, pwdMD5, uid)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

// 修改支付密码
func (*ServicerDao) ModifyServicerPubKey(servicerNo, pubKey string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update servicer set pub_key = $1 , modify_time = current_timestamp where servicer_no = $2 and is_delete = 0`, pubKey, servicerNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}

func (*ServicerDao) GetAccountNoFromServiceNo(serviceNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_no from servicer where servicer_no=$1 and is_delete='0' limit 1`, []*sql.NullString{&accountNoT}, serviceNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}

	return accountNoT.String
}

func (*ServicerDao) GetServicerFromServiceNo(serviceNo string) (string, string, string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var openIdxT, contactPersonT, contactPhoneT, contactAddrT, addrT, incomeAuthorizationT, outgoAuthorizationT, createTimeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select open_idx,contact_person,contact_phone,contact_addr,addr,income_authorization,outgo_authorization,create_time
			from servicer where servicer_no=$1 and is_delete='0' limit 1`, []*sql.NullString{&openIdxT, &contactPersonT, &contactPhoneT, &contactAddrT, &addrT, &incomeAuthorizationT, &outgoAuthorizationT, &createTimeT}, serviceNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", "", "", "", ""
	}

	return openIdxT.String, contactPersonT.String, contactPhoneT.String, contactAddrT.String, addrT.String, incomeAuthorizationT.String, outgoAuthorizationT.String, createTimeT.String
}

func (*ServicerDao) GetServicerNoByAccNo(accNo string) (servicerNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select servicer_no from servicer where account_no = $1 and is_delete = '0' "
	var servicerNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&servicerNoT}, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return servicerNoT.String
}

func (ServicerDao) GetServicerPWDFromIdenNo(idenNo string) (payPwd string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select password from servicer where servicer_no =$1 and is_delete='0' limit 1`, []*sql.NullString{&pwdT}, idenNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return pwdT.String
}

//根据服务商id获取账号
func (ServicerDao) GetAccountBySerNo(servicerNo string) (account string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select acc.account " +
		" from servicer ser " +
		" left join account acc on ser.account_no = acc.uid and acc.is_delete ='0' " +
		" where ser.servicer_no = $1 and ser.is_delete = '0' "
	var accountT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accountT}, servicerNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return accountT.String, nil
}
