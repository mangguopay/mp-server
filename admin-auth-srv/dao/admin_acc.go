package dao

import (
	"a.a/cu/encrypt"
	"a.a/mp-server/common/cache"
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	AdminAccDaoInstance AdminAccDao
)

type AdminAccDao struct {
}

//检查账户是否已存在
func (*AdminAccDao) CheckAccount(tx *sql.Tx, account string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM admin_account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	return strext.ToInt(count.String), err
}

func (*AdminAccDao) GetAccountCnt(account string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM admin_account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return 0, err
	}
	return strext.ToInt(count.String), nil
}

func (*AdminAccDao) CheckAccountUpdate(tx *sql.Tx, accountUid string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM admin_account WHERE  uid = $1 and is_delete = '0' ", []*sql.NullString{&count}, accountUid)
	return strext.ToInt(count.String), err
}

func (r *AdminAccDao) AddAccount(tx *sql.Tx, nickname, _account, password, useStatus, masterAccount, phone, countryCode, utmSource string) (string, error) {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}

	ss_log.Info("masterAccount==[%v]", masterAccount)
	accountUid := strext.NewUUID()
	if masterAccount == "" {
		masterAccount = "00000000-0000-0000-0000-000000000000"
	}
	err := ss_sql.ExecTx(tx, "INSERT INTO admin_account(uid,nickname,account,password,use_status,"+
		"master_acc,phone,create_time,gen_key,is_actived,country_code,utm_source) VALUES ($1,$2,$3,$4,$5,$6,$7,current_timestamp,$8,'1',$9,$10)",
		accountUid, nickname, _account, pwdMD5, useStatus, masterAccount, phone, util.RandomDigitStrOnlyNum(10), countryCode, utmSource)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	ss_log.Info("SaveAccount | AdminAccDao | AddAccount accountUid=%v", accountUid)
	ss_log.Info("SaveAccount | AdminAccDao | AddAccount pwdMD5=%v", pwdMD5)

	//if masterAccount == ss_sql.UUID {
	//	// 主账号默认授权，子账号默认空
	//	err = authAccount(tx, accountType, accountUid)
	//	if nil != err {
	//		ss_log.Error("授权失败,没找到默认角色|err=[%v]", err)
	//		// 不管授权问题
	//		return accountUid, nil
	//	}
	//}
	return accountUid, nil
}

func (r *AdminAccDao) AddAdminAccount(tx *sql.Tx, _account, password, useStatus, phone string) (string, error) {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}

	accountUid := strext.NewUUID()

	err := ss_sql.ExecTx(tx, "INSERT INTO admin_account(uid,account,password,use_status,"+
		"phone,create_time) VALUES ($1,$2,$3,$4,$5,current_timestamp)",
		accountUid, _account, pwdMD5, useStatus, phone)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	ss_log.Info("SaveAccount | AdminAccDao | AddAccount accountUid=%v", accountUid)
	ss_log.Info("SaveAccount | AdminAccDao | AddAccount pwdMD5=%v", pwdMD5)

	return accountUid, nil
}

func (r *AdminAccDao) AuthAccountRetCode(tx *sql.Tx, accountType, accountUid string) string {
	var roleUid sql.NullString
	err := ss_sql.QueryRowTx(tx, "select role_no from admin_role where acc_type=$1 and def_type='1' and is_delete='0' LIMIT 1", []*sql.NullString{&roleUid}, accountType)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
	}
	err = ss_sql.ExecTx(tx, "insert into admin_rela_account_role(rela_uid,account_uid,role_uid) VALUES ($1,$2,$3)", strext.NewUUID(), accountUid, roleUid.String)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

// xxx
func (r *AdminAccDao) UpdateAccount(tx *sql.Tx, nickname, _account, password, useStatus, accountUid, phone, email string) error {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}
	var errUpdate error
	sqlUpdate, data, _, err := ss_sql.MkUpdateSql("admin_account", map[string]string{
		"modify_time": ss_time.NowForPostgres(global.Tz),
		"nickname":    nickname,
		"account":     _account,
		"password":    pwdMD5,
		"use_status":  useStatus,
		"phone":       phone,
		"email":       email,
	}, " WHERE uid='"+accountUid+"'")
	if err != nil {
		ss_log.Error("err=%v", err)
	}
	ss_log.Info("sql=[%v]\ndata=[%v]", sqlUpdate, data)
	errUpdate = ss_sql.ExecTx(tx, sqlUpdate, data...)
	if errUpdate != nil {
		ss_log.Error("errUpdate=%v", errUpdate)
	}
	return errUpdate
}

// xxx
func (r *AdminAccDao) UpdateAdminAccount(tx *sql.Tx, _account, password, useStatus, accountUid, phone, email string) error {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}
	var errUpdate error
	sqlUpdate, data, _, err := ss_sql.MkUpdateSql("admin_account", map[string]string{
		"modify_time": ss_time.NowForPostgres(global.Tz),
		"account":     _account,
		"password":    pwdMD5,
		"use_status":  useStatus,
		"phone":       phone,
		"email":       email,
	}, " WHERE uid='"+accountUid+"'")
	if err != nil {
		ss_log.Error("err=%v", err)
	}
	ss_log.Info("sql=[%v]\ndata=[%v]", sqlUpdate, data)
	errUpdate = ss_sql.ExecTx(tx, sqlUpdate, data...)
	if errUpdate != nil {
		ss_log.Error("errUpdate=%v", errUpdate)
	}
	return errUpdate
}

func (*AdminAccDao) GetAccountTypeFromAccNoAdminOrOp(accountNo string) (accountType string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_type from admin_rela_acc_iden where account_no=$1 and account_type in($2,$3) limit 1`,
		[]*sql.NullString{&accountTypeT}, accountNo, constants.AccountType_ADMIN, constants.AccountType_OPERATOR)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountTypeT.String
}

func (*AdminAccDao) GetIdenNoFromAcc(accountNo, accountType string) (idenNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var idenNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select iden_no from admin_rela_acc_iden where account_no=$1 and account_type=$2 limit 1`,
		[]*sql.NullString{&idenNoT}, accountNo, accountType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return idenNoT.String
}

func (*AdminAccDao) GetAdminAccountByUid(uid string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var account sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account from admin_account where uid=$1 and is_delete = '0' `, []*sql.NullString{&account}, uid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return account.String
}

func (*AdminAccDao) GetCountryCodePhoneByUidTx(tx *sql.Tx, uid string) (string, string, error) {
	var countryCode, phone sql.NullString
	err := ss_sql.QueryRowTx(tx, `select country_code,phone from admin_account where uid=$1 and is_delete = '0' `,
		[]*sql.NullString{&countryCode, &phone}, uid)

	return countryCode.String, phone.String, err
}

// 初始化密钥
func (AdminAccDao) InitPassword(password string) string {
	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	return encrypt.DoMd5Salted(password, passwordSalt)
}
