package dao

import (
	"a.a/mp-server/common/ss_func"
	"database/sql"
	"fmt"

	"a.a/cu/db"
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	AccDaoInstance AccDao
)

type AccDao struct {
}

func (r *AccDao) InsertAccount(account, nickname, phone, password, genKey, isActived string) (accountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	accountNoT := strext.NewUUID()
	pwdMD5 := r.InitPassword(password)
	err := ss_sql.Exec(dbHandler, `insert into account(uid,nickname,account,password,use_status,create_time,phone,gen_key,is_actived)values($1,$2,$3,$4,'1',current_timestamp,$5,$6,$7)`,
		accountNoT, nickname, account, pwdMD5, phone, genKey, isActived)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return accountNoT
}

func (*AccDao) QueryGenKeyFromAccNo(accNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var genKey sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT gen_key FROM account WHERE uid=$1 and is_delete = '0' and is_actived='1' LIMIT 1",
		[]*sql.NullString{&genKey}, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return genKey.String
}

func (r *AccDao) UpdateAccountPwd(phone, pwd, CountryCode string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	pwdMD5 := r.InitPassword(pwd)
	err := ss_sql.Exec(dbHandler, `update account set password = $1 , modify_time = current_timestamp where 
		phone = $2 and country_code = $3 and is_delete = 0`, pwdMD5, phone, CountryCode)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
func (r *AccDao) UpdateAccountPwdTx(tx *sql.Tx, phone, pwd, CountryCode, nickNane string) error {
	pwdMD5 := r.InitPassword(pwd)
	err := ss_sql.ExecTx(tx, `update account set password = $1 ,nickname=$2, modify_time = current_timestamp where 
		phone = $3 and country_code = $4 and is_delete = 0`, pwdMD5, nickNane, phone, CountryCode)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
func (r *AccDao) UpdateAccountPwdByUID(uid, oldPassword, newPassword, nonStr string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var dbPWD sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT password FROM account WHERE uid=$1 and is_delete = 0 LIMIT 1", []*sql.NullString{&dbPWD}, uid)
	if err != nil {
		ss_log.Error("err=[%v]", err.Error())
	}

	// 数据库取出的支付密码加盐(加的是和前端传来的盐一样)
	pwdMD5FixedDB := encrypt.DoMd5Salted(dbPWD.String, nonStr)

	if pwdMD5FixedDB != oldPassword {
		ss_log.Error("err=[%v],pwdMD5FixedDB[%v],oldPassword[%v]", "原密码不正确", pwdMD5FixedDB, oldPassword)
		//return ss_err.ERR_ACCOUNT_OLD_PWD_FAILD
		return ss_err.ERR_Business_OLD_PWD_FAILD
	}

	// 修改新登录密码
	newPwdMD5 := r.InitPassword(newPassword)
	if err := ss_sql.Exec(dbHandler, `update account set password = $1 , modify_time = current_timestamp where uid = $2 and is_delete = 0`, newPwdMD5, uid); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return ss_err.ERR_MODIFY_ACCOUNT_PWD_FAILD
	}

	return ss_err.ERR_SUCCESS
}

func (r *AccDao) QueryUIDByPhone(phone, CountryCode string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var uid sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT uid FROM account WHERE phone=$1 and country_code = $2 and "+
		"is_delete = 0 LIMIT 1", []*sql.NullString{&uid}, phone, CountryCode)
	if err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "", err
	}
	return uid.String, nil
}

func (*AccDao) InitializeWallet(tx *sql.Tx, accountUid string) error {
	walletNo := strext.NewUUID()
	sqlStr := "INSERT INTO wallet(wallet_no,remain,account_no,create_time) VALUES ($1,0,$2,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, walletNo, accountUid)
	return err
}

//检查账户是否已存在
func (*AccDao) CheckAccountTx(tx *sql.Tx, account string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	return strext.ToInt(count.String), err
}

func (*AccDao) GetAccountCnt(account string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return -1, err
	}
	return strext.ToInt(count.String), nil
}

func (*AccDao) IsAccountExists(account string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}
	return strext.ToInt(count.String) > 0
}

func (*AccDao) CheckAccountUpdate(tx *sql.Tx, accountUid string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM account WHERE  uid = $1 and is_delete = '0' ", []*sql.NullString{&count}, accountUid)
	return strext.ToInt(count.String), err
}

func (r *AccDao) AddAccount(tx *sql.Tx, nickname, _account, password, useStatus, masterAccount, phone, countryCode, utmSource, email string) (string, error) {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}

	ss_log.Info("masterAccount==[%v]", masterAccount)
	accountUid := strext.NewUUID()
	if masterAccount == "" {
		masterAccount = "00000000-0000-0000-0000-000000000000"
	}
	err := ss_sql.ExecTx(tx, "INSERT INTO account(uid,nickname,account,password,use_status,"+
		"master_acc,phone,create_time,gen_key,is_actived,country_code,utm_source,email) VALUES ($1,$2,$3,$4,$5,$6,$7,current_timestamp,$8,'1',$9,$10,$11)",
		accountUid, nickname, _account, pwdMD5, useStatus, masterAccount, phone, util.RandomDigitStrOnlyNum(10), countryCode, utmSource, email)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	ss_log.Info("SaveAccount | AccDao | AddAccount accountUid=%v", accountUid)
	ss_log.Info("SaveAccount | AccDao | AddAccount pwdMD5=%v", pwdMD5)

	return accountUid, nil
}

//商家的手机号是存在business_phone
func (r *AccDao) AddBusinessAccount(tx *sql.Tx, nickname, account, password, useStatus, masterAccount, phone, countryCode, utmSource, email string) (string, error) {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}

	ss_log.Info("masterAccount==[%v]", masterAccount)
	accountUid := strext.NewUUID()
	if masterAccount == "" {
		masterAccount = "00000000-0000-0000-0000-000000000000"
	}
	err := ss_sql.ExecTx(tx, "INSERT INTO account(uid,nickname,account,password,use_status,"+
		"master_acc,business_phone,create_time,gen_key,is_actived,country_code,utm_source,email) VALUES ($1,$2,$3,$4,$5,$6,$7,current_timestamp,$8,'1',$9,$10,$11)",
		accountUid, nickname, account, pwdMD5, useStatus, masterAccount, phone, util.RandomDigitStrOnlyNum(10), countryCode, utmSource, email)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return "", err
	}

	return accountUid, nil
}

func authAccount(tx *sql.Tx, accountType, accountUid string) error {
	var roleUid sql.NullString
	err := ss_sql.QueryRowTx(tx, "select role_no from role where acc_type=$1 and def_type='1' LIMIT 1", []*sql.NullString{&roleUid}, accountType)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	err = ss_sql.ExecTx(tx, "insert into rela_account_role(rela_uid,account_uid,role_uid) VALUES ($1,$2,$3)", strext.NewUUID(), accountUid, roleUid.String)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (r *AccDao) AuthAccountRetCode(tx *sql.Tx, accountType, accountUid string) string {
	var roleUid sql.NullString
	err := ss_sql.QueryRowTx(tx, "select role_no from role where acc_type=$1 and def_type='1' and is_delete='0' LIMIT 1", []*sql.NullString{&roleUid}, accountType)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
	}
	err = ss_sql.ExecTx(tx, "insert into rela_account_role(rela_uid,account_uid,role_uid) VALUES ($1,$2,$3)", strext.NewUUID(), accountUid, roleUid.String)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (r *AccDao) AddAccountWithAccNo(tx *sql.Tx, nickname, accNo, _account, password, useStatus string) string {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}
	err := ss_sql.ExecTx(tx, "INSERT INTO account(uid,nickname,account,password,use_status,create_time) VALUES ($1,$2,$3,$4,$5,current_timestamp)",
		accNo, nickname, _account, pwdMD5, useStatus)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_ACCOUNT_ALREADY_EXISTS
	}

	//errCode := authAccountRetCode(tx, accountType, accNo)
	return ss_err.ERR_SUCCESS
}

func (*AccDao) CheckUpdateAccount(tx *sql.Tx, _account, accountUid string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM account WHERE account = $1 and uid != $2", []*sql.NullString{&count}, _account, accountUid)
	return strext.ToInt(count), err
}

// xxx
func (r *AccDao) UpdateAccount(tx *sql.Tx, nickname, _account, password, useStatus, accountUid, phone, email string) error {
	var pwdMD5 string
	if password != "" {
		pwdMD5 = r.InitPassword(password)
	}
	var errUpdate error
	sqlUpdate, data, _, err := ss_sql.MkUpdateSql("account", map[string]string{
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

func (*AccDao) HasJumpRela(loginAccNo, idenNo, idenType string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count, accountType, mercNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select count(1) from rela_keyman_acc_iden where acc_no=$1 and iden_no=$2 and iden_type=$3 ",
		[]*sql.NullString{&count}, loginAccNo, idenNo, idenType)
	if nil != err || strext.ToInt(count.String) <= 0 {
		ss_log.Error("err=[%v]", err)
		// 还有一种方式
		err := ss_sql.QueryRow(dbHandler, "select account_type,affiliation_uid from account where uid=$1 limit 1 ",
			[]*sql.NullString{&accountType, &mercNo}, loginAccNo)
		if nil != err {
			ss_log.Error("err=[%v]", err)
			return false
		}

		switch accountType.String {
		case constants.AccountType_ADMIN:
			fallthrough
		case constants.AccountType_OPERATOR:
			return true
		case constants.AccountType_SERVICER:
			return false
		default:
			return false
		}
	}
	return strext.ToInt(count.String) > 0
}

/**
 * 身份id获取账号id
 */
func (*AccDao) GetAccNoFromIden(idenNo, idenType string) (accountNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNoT sql.NullString
	switch idenType {
	case constants.AccountType_OPERATOR:
		err = ss_sql.QueryRow(dbHandler, "select account_no from agency where agency_no=$1 limit 1", []*sql.NullString{&accountNoT}, idenNo)
	case constants.AccountType_SERVICER:
		err = ss_sql.QueryRow(dbHandler, "select account_no from merchant where merchant_no=$1 limit 1", []*sql.NullString{&accountNoT}, idenNo)
	}
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return accountNoT.String, err
	}
	return accountNoT.String, nil
}

// 删除单个或多个账号
func (r *AccDao) DeleteAccountList(accNos interface{}) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	switch accNos.(type) {
	case string:
		err := r.DeleteAccount(accNos.(string))
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return
		}
		ss_log.Info("删除账户[%v]成功", accNos)
	case []string:
		for _, v := range accNos.([]string) {
			err := r.DeleteAccount(v)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				continue
			}
			ss_log.Info("删除账户[%v]成功", v)
		}
	}
}

// 删除单个账号
func (*AccDao) DeleteAccount(accNo string) (errr error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update account set is_delete='1' where uid=$1`, accNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
func (*AccDao) UpdateAccountStatusByAccount(account string, status int) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update account set use_status= $2 where account=$1`, account, status)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}

func (*AccDao) GetAccByNickname(nickname string) *go_micro_srv_auth.Account {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	account := go_micro_srv_auth.Account{}
	var uid, accountT, useStatus, createTime, modifyTime sql.NullString
	rowAccount := ss_sql.QueryRow(dbHandler, "SELECT uid,account,use_status,create_time,modify_time FROM account WHERE nickname = $1 and is_delete='0' LIMIT 1",
		[]*sql.NullString{&uid, &accountT, &useStatus, &createTime, &modifyTime}, nickname)
	if nil == rowAccount {
		return nil
	}

	account.Uid = uid.String
	account.Nickname = nickname
	account.Account = accountT.String
	account.UseStatus = useStatus.String
	account.CreateTime = createTime.String
	account.ModifyTime = modifyTime.String
	return &account
}

func (*AccDao) GetRoleAuthedUrlList(accountNo string) (errCode string, datas []*go_micro_srv_auth.RoleSimpleData) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT role_uid,role_name FROM role WHERE role_uid in (SELECT role_uid FROM rela_account_role WHERE account_uid = $1) and is_delete='0'",
		accountNo)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}
	ss_log.Error("err=[%v]", err)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD, nil
	}

	datas = []*go_micro_srv_auth.RoleSimpleData{}
	for rows.Next() {
		data := go_micro_srv_auth.RoleSimpleData{}
		err := rows.Scan(&data.RoleNo, &data.RoleName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}
	return ss_err.ERR_SUCCESS, datas
}

func (*AccDao) GetRoleAllUrlList(accountNo string) (errCode string, datas []*go_micro_srv_auth.RoleSimpleData) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	rows, stmt, err := ss_sql.Query(dbHandler, "SELECT role_uid,role_name FROM role where is_delete='0'",
		accountNo)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}
	ss_log.Error("err=[%v]", err)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_ADD, nil
	}

	datas = []*go_micro_srv_auth.RoleSimpleData{}
	for rows.Next() {
		data := go_micro_srv_auth.RoleSimpleData{}
		err := rows.Scan(&data.RoleNo, &data.RoleName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}
	return ss_err.ERR_SUCCESS, datas
}

func (*AccDao) GetAccountTypeFromAccNoAdminOrOp(accountNo string) (accountType string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_type from rela_acc_iden where account_no=$1 and account_type in($2,$3) limit 1`,
		[]*sql.NullString{&accountTypeT}, accountNo, constants.AccountType_ADMIN, constants.AccountType_OPERATOR)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountTypeT.String
}

func (*AccDao) GetAccountTypeFromAccNoBusiness(accountNo string) (accountType string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_type from rela_acc_iden where account_no = $1 and account_type in($2,$3)  limit 1`,
		[]*sql.NullString{&accountTypeT}, accountNo, constants.AccountType_PersonalBusiness, constants.AccountType_EnterpriseBusiness)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountTypeT.String
}

func (*AccDao) GetAccountTypeByAccountNo(accountNo string) (accountTypes string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountTypeT []string
	sqlStr := " select account_type from rela_acc_iden where account_no= $1 "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, accountNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var accountType sql.NullString
			err = rows.Scan(
				&accountType,
			)
			if err != nil {
				ss_log.Error("err=%v", err)
				continue
			}
			if accountType.String != "" {
				accountTypeT = append(accountTypeT, accountType.String)
			}
		}
	} else {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	for _, v := range accountTypeT {
		if accountTypes != "" {
			accountTypes = accountTypes + "," + v
		} else {
			accountTypes = v
		}
	}

	return accountTypes
}

func (*AccDao) GetIdenNoFromAcc(accountNo, accountType string) (idenNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var idenNoT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select iden_no from rela_acc_iden where account_no=$1 and account_type=$2 limit 1`,
		[]*sql.NullString{&idenNoT}, accountNo, accountType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return idenNoT.String
}

func (*AccDao) GetAccountIsActived(preCountryCode, countryCode, phone string) (accountUid, isActived string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountUidT, isActivedT sql.NullString
	//err := ss_sql.QueryRow(dbHandler, `select uid, is_actived from account where phone=$1 and country_code = $2 and is_delete='0'`, []*sql.NullString{&accountUidT, &isActivedT}, phone, countryCode)
	//if accountUidT.String != "" {
	//	ss_log.Error("err=[%v]", err)
	//	return accountUidT.String, isActivedT.String
	//}
	//err = ss_sql.QueryRow(dbHandler, `select uid, is_actived from account where nickname=$1 and is_delete='0'`, []*sql.NullString{&accountUidT, &isActivedT}, phone)
	//if accountUidT.String != "" {
	//	ss_log.Error("err=[%v]", err)
	//	return accountUidT.String, isActivedT.String
	//}
	err := ss_sql.QueryRow(dbHandler, `select uid, is_actived from account where account=$1 and is_delete='0'`,
		[]*sql.NullString{&accountUidT, &isActivedT}, fmt.Sprintf("%s%s", preCountryCode, phone))
	if accountUidT.String != "" {
		ss_log.Error("err=[%v]", err)
		return accountUidT.String, isActivedT.String
	}

	return accountUidT.String, isActivedT.String
}

//将未激活的账号修改为激活
func (*AccDao) UpdateAccountIsActived(tx *sql.Tx, accUid string) (err error) {
	if errT := ss_sql.ExecTx(tx, `update account set is_actived ='1' where uid = $1 and is_actived = '0' and is_delete = '0' `, accUid); errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*AccDao) CheckAccountExists(preCountryCode, countryCode, phone string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var cnt sql.NullString
	//err := ss_sql.QueryRow(dbHandler, `select count(1) from account where phone=$1 and country_code = $2 and is_delete='0'`, []*sql.NullString{&cnt}, phone, countryCode)
	//if strext.ToInt(cnt.String) >= 1 {
	//	ss_log.Error("err=[%v]", err)
	//	return true
	//}
	//err = ss_sql.QueryRow(dbHandler, `select count(1) from account where nickname=$1 and is_delete='0'`, []*sql.NullString{&cnt}, phone)
	//if strext.ToInt(cnt.String) >= 1 {
	//	ss_log.Error("err=[%v]", err)
	//	return true
	//}
	err := ss_sql.QueryRow(dbHandler, `select count(1) from account where account=$1 and is_delete='0'`,
		[]*sql.NullString{&cnt}, fmt.Sprintf("%s%s", preCountryCode, phone))
	if strext.ToInt(cnt.String) >= 1 {
		ss_log.Error("err=[%v]", err)
		return true
	}

	return false
}

// 初始化密钥
func (AccDao) InitPassword(password string) string {
	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	return encrypt.DoMd5Salted(password, passwordSalt)
}

func (*AccDao) GetRemain(accountNo string) (khr, usd string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var khrT, usdT sql.NullString
	sqlStr := "select balance from vaccount where account_no = $1 and va_type = $2 and balance_type = $3 limit 1"
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&usdT}, accountNo, constants.VaType_USD_DEBIT, constants.CURRENCY_USD); nil != err {
		ss_log.Error("err=%v", err)
		//return "0", "0"
	}

	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&khrT}, accountNo, constants.VaType_KHR_DEBIT, constants.CURRENCY_KHR); nil != err {
		ss_log.Error("err=%v", err)
		//return "0", "0"
	}

	if khrT.String == "" {
		khrT.String = "0"
	}

	if usdT.String == "" {
		usdT.String = "0"
	}
	return khrT.String, usdT.String
}

func (*AccDao) GetAccountByUid(uid string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var account sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account from account where uid=$1 and is_delete = '0' `, []*sql.NullString{&account}, uid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return account.String
}

func (*AccDao) GetAccountMailByUid(uid string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var email sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select email from account where uid=$1 and is_delete = '0' `, []*sql.NullString{&email}, uid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return email.String
}

func (*AccDao) GetCountryCodePhoneByUid(uid string) (countryCodeT string, phoneT string, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var countryCode, phone sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select country_code,phone from account where uid=$1 and is_delete = '0' `,
		[]*sql.NullString{&countryCode, &phone}, uid)

	return countryCode.String, phone.String, err
}

func (*AccDao) GetCountryCodePhoneByUidTx(tx *sql.Tx, uid string) (countryCodeT string, phoneT string, errT error) {
	var countryCode, phone sql.NullString
	err := ss_sql.QueryRowTx(tx, `select country_code,phone from account where uid=$1 and is_delete = '0' `,
		[]*sql.NullString{&countryCode, &phone}, uid)

	return countryCode.String, phone.String, err
}

func (*AccDao) GetAccNoFromPhone(phone, countryCode string) (accNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	account := ss_func.ComposeAccountByPhoneCountryCode(phone, countryCode)
	var accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select uid from account where account = $1 and is_delete='0' limit 1`, []*sql.NullString{&accountNo}, account)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return accountNo.String
}

func (*AccDao) GetAccNoFromAccount(account string) (accNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select uid from account where account=$1 and is_delete='0' limit 1`, []*sql.NullString{&accountNo}, account)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return accountNo.String
}

//检查该账号是否有该条最近转账人信息
func (*AccDao) CheckAccountCollect(accountNo, toPhone string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	sqlCnt := "SELECT COUNT(1) FROM account_collect WHERE account_no = $1 and collect_phone = $2 and is_delete='0' "

	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&count}, accountNo, toPhone)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	return strext.ToInt(count.String), err
}

// 检查该账号有多少条记录
func (*AccDao) GetAccountCollectCount(accountNo string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	sqlCnt := "SELECT COUNT(1) FROM account_collect WHERE account_no = $1 and is_delete='0' "
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&count}, accountNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	return strext.ToInt(count.String), err
}

func (*AccDao) GetDefPayNo(custNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var genKey sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select def_pay_no from cust where cust_no=$1 limit 1",
		[]*sql.NullString{&genKey}, custNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return genKey.String
}

func (*AccDao) DelFriend(accountNo, toPhone string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlCnt := "update account_collect set is_delete='1' WHERE account_no = $1 and collect_phone = $2 and is_delete='0' "
	err := ss_sql.Exec(dbHandler, sqlCnt, accountNo, toPhone)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
func (*AccDao) UpdateLang(accountNo, lang string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlCnt := "update account set app_lang= $2 WHERE uid = $1  and is_delete='0' "
	err := ss_sql.Exec(dbHandler, sqlCnt, accountNo, lang)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

func (r *AccDao) QueryAccountLang(accountNo string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var appLang, posLang sql.NullString
	sqlCnt := "select app_lang,pos_lang from account where uid=$1 and is_delete='0'  limit 1 "
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&appLang, &posLang}, accountNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", ""
	}
	return appLang.String, posLang.String
}

func (r *AccDao) DelLastFriend(accountNo string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var phone sql.NullString
	sqlCnt := "select collect_phone from account_collect where account_no=$1 and is_delete='0' order by modify_time asc limit 1 "
	err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&phone}, accountNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	if phone.String == "" {
		// 没有可以删的那么就不删
		return nil
	}

	return r.DelFriend(accountNo, phone.String)
}

//用户、个人商家修改手机号
func (AccDao) ModifyPhone(tx *sql.Tx, uid, account, countryCode, phone string) error {
	return ss_sql.ExecTx(tx, `update account set phone = $1 ,account = $2,country_code = $3, modify_time = current_timestamp 
		where uid = $4 and is_delete = 0 `, phone, account, countryCode, uid)
}

//企业商家修改手机号。
func (AccDao) BusinessModifyBusinessPhoneTx(tx *sql.Tx, uid, businessPhone, countryCode string) error {
	return ss_sql.ExecTx(tx, `update account set business_phone = $2, country_code = $3, modify_time = current_timestamp 
		where uid = $1 and is_delete = 0 `, uid, businessPhone, countryCode)
}

func (AccDao) QueryPWDByAccountUid(uid string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdMD5 sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT password FROM account WHERE uid=$1 and is_delete = '0' LIMIT 1", []*sql.NullString{&pwdMD5}, uid)
	if err != nil {
		return ""
	}

	return pwdMD5.String
}

func (AccDao) ModifyAccountEmailTx(tx *sql.Tx, uid, email string) error {
	sqlStr := "update account set email = $2 where uid = $1 and is_delete = '0' "

	err := ss_sql.ExecTx(tx, sqlStr, uid, email)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//检查账户是否已存在
func (*AccDao) CheckAccount(account string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	return strext.ToInt(count.String), err
}

//修改账号account
func (*AccDao) ModifyAccountByUidTx(tx *sql.Tx, uid, account string) error {
	sqlStr := "update account set account = $2 where uid = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, uid, account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

func (r *AccDao) UpdateLoginPwdByUID(uid, newPassword string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	newPwdMD5 := r.InitPassword(newPassword)

	// 修改密码
	if err := ss_sql.Exec(dbHandler, `update account set password = $1 , modify_time = current_timestamp where uid = $2 and is_delete = 0`, newPwdMD5, uid); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return err
	}

	return nil
}

//确认企业商家的手机号与国家码组成是唯一的
func (r *AccDao) CheckBusinessPhoneAndCountryCodeUnique(businessPhone, countryCode string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from account where business_phone = $1 and country_code = $2 "
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, businessPhone, countryCode); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return false
	}

	return cnt.String == "0"
}
