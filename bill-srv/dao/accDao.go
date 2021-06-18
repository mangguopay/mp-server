package dao

import (
	"database/sql"

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
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_sql"
)

var (
	AccDaoInstance AccDao
)

type AccDao struct {
	AccountNo   string
	Account     string
	AccountType string
	IdentityNo  string
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

func (r *AccDao) InsertAccount(account, nickname, phone, password string) (accountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	accountNoT := strext.NewUUID()
	pwdMD5 := r.InitPassword(password)
	err := ss_sql.Exec(dbHandler, `insert into account(uid,nickname,account,password,use_status,create_time,phone)values($1,$2,$3,$4,'1',current_timestamp,$5)`,
		accountNoT, nickname, account, pwdMD5, phone)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}

	return accountNoT
}

func (r *AccDao) InsertEmptyAccount(phone, countryCode string) (accountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	genKey := util.RandomDigitStrOnlyNum(10)
	account := ss_func.ComposeAccountByPhoneCountryCode(phone, countryCode)
	accountNoT := strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into account(uid,is_actived,create_time,phone,gen_key,account,country_code)values($1,$2,current_timestamp,$3,$4,$5,$6)`,
		accountNoT, 0, phone, genKey, account, countryCode)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}

	return accountNoT
}

func (*AccDao) InitializeWallet(tx *sql.Tx, accountUid string) error {
	walletNo := strext.NewUUID()
	sqlStr := "INSERT INTO wallet(wallet_no,remain,account_no,create_time) VALUES ($1,0,$2,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, walletNo, accountUid)
	return err
}

//检查账户是否已存在
func (*AccDao) CheckAccount(tx *sql.Tx, account string) (int, error) {
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
		ss_log.Error("err=%v", err)
		return 0, err
	}
	return strext.ToInt(count.String), nil
}

func (*AccDao) IsAccountExists(account string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	if err != nil {
		ss_log.Error("err=%v", err)
		return false
	}
	return strext.ToInt(count.String) > 0
}

func (*AccDao) CheckAccountUpdate(tx *sql.Tx, account, accountUid string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM account WHERE account = $1 and uid != $2", []*sql.NullString{&count}, account, accountUid)
	return strext.ToInt(count.String), err
}

func (r *AccDao) AddAccount(tx *sql.Tx, nickname, _account, password, useStatus, masterAccount, phone string) (string, error) {
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
		"master_acc,phone,create_time) VALUES ($1,$2,$3,$4,$5,$6,$7,current_timestamp)",
		accountUid, nickname, _account, pwdMD5, useStatus, masterAccount, phone)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", err
	}
	ss_log.Info("SaveAccount | AccDao | AddAccount accountUid=%v", accountUid)
	ss_log.Info("SaveAccount | AccDao | AddAccount pwdMD5=%v", pwdMD5)

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

func authAccount(tx *sql.Tx, accountType, accountUid string) error {
	var roleUid sql.NullString
	err := ss_sql.QueryRowTx(tx, "select role_no from role where acc_type=$1 and def_type='1' LIMIT 1", []*sql.NullString{&roleUid}, accountType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return err
	}
	err = ss_sql.ExecTx(tx, "insert into rela_account_role(rela_uid,account_uid,role_uid) VALUES ($1,$2,$3)", strext.NewUUID(), accountUid, roleUid.String)
	if nil != err {
		ss_log.Error("err=%v", err)
		return err
	}
	return nil
}

func (r *AccDao) QeuryNamePhoneFromGenKey(genKey string) (string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var uidT, phoneT, nickNameT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select uid,phone,nickname from account where gen_key=$1 and is_delete='0' LIMIT 1", []*sql.NullString{&uidT, &phoneT, &nickNameT}, genKey)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", ""
	}
	return uidT.String, phoneT.String, nickNameT.String
}

func authAccountRetCode(tx *sql.Tx, accountType, accountUid string) string {
	var roleUid sql.NullString
	err := ss_sql.QueryRowTx(tx, "select role_no from role where acc_type=$1 and def_type='1' and is_delete='0' LIMIT 1", []*sql.NullString{&roleUid}, accountType)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_ACCOUNT_ROLE_NOT_EXISTS
	}
	err = ss_sql.ExecTx(tx, "insert into rela_account_role(rela_uid,account_uid,role_uid) VALUES ($1,$2,$3)", strext.NewUUID(), accountUid, roleUid.String)
	if nil != err {
		ss_log.Error("err=%v", err)
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
		ss_log.Error("err=%v", err)
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
		ss_log.Error("err=%v", err)
		// 还有一种方式
		err := ss_sql.QueryRow(dbHandler, "select account_type,affiliation_uid from account where uid=$1 limit 1 ",
			[]*sql.NullString{&accountType, &mercNo}, loginAccNo)
		if nil != err {
			ss_log.Error("err=%v", err)
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
		ss_log.Error("err=%v", err)
		return accountNoT.String, err
	}
	return accountNoT.String, nil
}

func (*AccDao) GetAccNoFromAccount(account string) (accNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accNoT sql.NullString
	err = ss_sql.QueryRow(dbHandler, "select uid from account where account=$1 and is_delete = '0' limit 1", []*sql.NullString{&accNoT}, account)
	if nil != err {
		ss_log.Error("err=%v", err)
		return accNoT.String, err
	}
	return accNoT.String, nil
}

// 删除单个或多个账号
func (r *AccDao) DeleteAccountList(accNos interface{}) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	switch accNos.(type) {
	case string:
		err := r.DeleteAccount(dbHandler, accNos.(string))
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=%v", err)
			return
		}
		ss_log.Info("删除账户[%v]成功", accNos)
	case []string:
		for _, v := range accNos.([]string) {
			err := r.DeleteAccount(dbHandler, v)
			if err != ss_err.ERR_SUCCESS {
				ss_log.Error("err=%v", err)
				continue
			}
			ss_log.Info("删除账户[%v]成功", v)
		}
	}
}

// 删除单个账号
func (*AccDao) DeleteAccount(dbHandler *sql.DB, accNo string) (errCode string) {
	err := ss_sql.Exec(dbHandler, `update account set is_delete='1' where uid=$1`, accNo)
	if nil != err {
		ss_log.Error("err=%v", err)
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

	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_ADD, nil
	}

	datas = []*go_micro_srv_auth.RoleSimpleData{}
	for rows.Next() {
		data := go_micro_srv_auth.RoleSimpleData{}
		err := rows.Scan(&data.RoleNo, &data.RoleName)
		if err != nil {
			ss_log.Error("err=%v", err)
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

	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_ADD, nil
	}

	datas = []*go_micro_srv_auth.RoleSimpleData{}
	for rows.Next() {
		data := go_micro_srv_auth.RoleSimpleData{}
		err := rows.Scan(&data.RoleNo, &data.RoleName)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		datas = append(datas, &data)
	}
	return ss_err.ERR_SUCCESS, datas
}

func (*AccDao) GetAccountTypeFromAccNo(accountNo string) (accountType string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountTypeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account_type from rela_acc_iden where account_no=$1 limit 1`, []*sql.NullString{&accountTypeT}, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountTypeT.String
}

//func (*AccDao) CheckAccountExists(phone string) bool {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//
//	var cnt sql.NullString
//	err := ss_sql.QueryRow(dbHandler, `select count(1) from account where phone=$1 and is_delete='0'`, []*sql.NullString{&cnt}, phone)
//	if strext.ToInt(cnt.String) > 1 {
//		ss_log.Error("err=%v", err)
//		return true
//	}
//	err = ss_sql.QueryRow(dbHandler, `select count(1) from account where nickname=$1 and is_delete='0'`, []*sql.NullString{&cnt}, phone)
//	if strext.ToInt(cnt.String) > 1 {
//		ss_log.Error("err=%v", err)
//		return true
//	}
//	err = ss_sql.QueryRow(dbHandler, `select count(1) from account where account=$1 and is_delete='0'`, []*sql.NullString{&cnt}, phone)
//	if strext.ToInt(cnt.String) > 1 {
//		ss_log.Error("err=%v", err)
//		return true
//	}
//
//	return false
//}

// 初始化密钥
func (AccDao) InitPassword(password string) string {
	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	return encrypt.DoMd5Salted(password, passwordSalt)
}

//验证管理员密码
func (a AccDao) CheckAdminLoginPWD(uid, pwd, nonStr string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdDb sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT password FROM admin_account WHERE uid=$1 and is_delete = '0' LIMIT 1", []*sql.NullString{&pwdDb}, uid)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	// 数据库取出的密码加盐(加的是和前端传来的盐一样)
	pwdMD5FixedDB := encrypt.DoMd5Salted(pwdDb.String, nonStr)
	ss_log.Info("pwd[%v],pwdDb[%v],pwdMd5[%v]", pwd, pwdDb.String, pwdMD5FixedDB)

	return pwd == pwdMD5FixedDB
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

func (*AccDao) GetAccNoFromPhone(phone, countryCode string) (accNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select uid from account where phone=$1 and country_code = $2 and is_delete='0' limit 1`, []*sql.NullString{&accountNo}, phone, countryCode)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountNo.String
}

func (*AccDao) GetNameAndURLFromUID(tx *sql.Tx, uid string) (string, string) {

	var nicknameT, headURLT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select ac.nickname,di.image_url from account ac LEFT JOIN dict_images di 
	ON ac.head_portrait_img_no = di.image_id WHERE  
	ac.uid= $1 and ac.is_delete='0' limit 1`, []*sql.NullString{&nicknameT, &headURLT}, uid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return nicknameT.String, headURLT.String
}
func (*AccDao) GetIsActiveAccNoFromPhone(phone, countryCode string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo, isActived sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select uid,is_actived from account where phone=$1 and country_code = $2 and is_delete='0'  limit 1`, []*sql.NullString{&accountNo, &isActived}, phone, countryCode)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return accountNo.String, isActived.String
}
func (*AccDao) GetIsActiveFromPhone(phone, countryCode string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	account := ss_func.ComposeAccountByPhoneCountryCode(phone, countryCode)

	var isActive sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select is_actived from account where account = $1 and is_delete='0' limit 1`, []*sql.NullString{&isActive}, account)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return isActive.String
}

func (*AccDao) GetPhoneFromAccNo(accNo string) (phone string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var phoneT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select phone from account where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&phoneT}, accNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return phoneT.String
}
func (*AccDao) GetPhoneCountryCodeFromAccNo(accNo string) (string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var phoneT, countryCode sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select phone,country_code from account where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&phoneT, &countryCode}, accNo)
	return phoneT.String, countryCode.String, err
}

// 修改正式账号USD的余额

// 修改正式账号khr的余额 khr_balance
//func (*AccDao) UpdateKHRRemain(accountNo, amount, op string) string {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//	tx, _ := dbHandler.BeginTx(ctx, nil)
//	defer ss_sql.Rollback(tx)
//
//	switch op {
//	case "+":
//		err := ss_sql.ExecTx(tx, `update account set khr_balance=khr_balance+$1 where uid=$2 and is_delete='0'`, amount, accountNo)
//		if nil != err {
//			ss_log.Error("err=%v", err)
//			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//		}
//	case "-":
//		err := ss_sql.ExecTx(tx, `update account set khr_balance=khr_balance-$1 where uid=$2 and is_delete='0'`, amount, accountNo)
//		if nil != err {
//			ss_err.ss_log.Error("err=%v", err)(err)
//			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//		}
//	default:
//		return ss_err.ERR_PAY_VACCOUNT_OP_MISSING
//	}
//
//	descroption := "KHR金额变动" + op + amount
//	errCode := LogAccDaoInstance.InsertAccountLog(descroption, accountNo, constants.REMAINTYPE)
//	if errCode != ss_err.ERR_SUCCESS {
//		ss_log.Error("err=%v", errCode)
//		return errCode
//	}
//
//	var tmp sql.NullString
//	err := ss_sql.QueryRowTx(tx, `select khr_balance from account where uid=$1 and is_delete='0' limit 1`,
//		[]*sql.NullString{&tmp}, accountNo)
//	if nil != err {
//		ss_err.ss_log.Error("err=%v", err)(err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//	if strext.ToInt64(tmp.String) < 0 {
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	ss_sql.Commit(tx)
//	return ss_err.ERR_SUCCESS
//}

func (*AccDao) GetPhoneFromVAccNo(vaccNo string) (phone string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var phoneT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select a.phone FROM vaccount va LEFT JOIN account a ON va.account_no = a.uid WHERE va.vaccount_no = $1 and a.is_delete='0' limit 1`, []*sql.NullString{&phoneT}, vaccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return phoneT.String
}
func (*AccDao) GetPhoneCountryCodeFromVAccNo(vaccNo string) (string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var phoneT, countryCodeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select a.phone,a.country_code FROM vaccount va LEFT JOIN 
		account a ON va.account_no = a.uid WHERE va.vaccount_no = $1 and a.is_delete='0' limit 1`, []*sql.NullString{&phoneT, &countryCodeT}, vaccNo)

	return phoneT.String, countryCodeT.String, err
}

func (a *AccDao) ConfirmAccIsExit(phone, countryCode string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	account := ss_func.ComposeAccountByPhoneCountryCode(phone, countryCode)

	var accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT uid FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&accountNo}, account)
	if err != nil {
		ss_log.Error("err=%v", err)
	}

	if accountNo.String == "" {
		if ok := GlobalParamDaoInstance.QeuryParamValue(constants.TransferToUnRegisteredUser); ok != "true" {
			ss_log.Error("转账至未注册用户开关为[%v]，不允许进行转账", ok)
			return ""
		}

		// 插入空的账号
		if err := CountryCodePhoneDaoInst.Insert1(countryCode, phone); err != nil {
			ss_log.Error("新增手机号和国家码进唯一表失败,phone: %s,err: %s", phone, err.Error())
			return ""
		}
		accountNo.String = a.InsertEmptyAccount(phone, countryCode)

		//初始化钱包3、4
		if vaccNo := VaccountDaoInst.InitVaccountNo(accountNo.String, constants.CURRENCY_USD, constants.VaType_FREEZE_USD_DEBIT); vaccNo == "" {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v", accountNo.String, constants.VaType_FREEZE_USD_DEBIT)
			//return err
		}
		if vaccNo := VaccountDaoInst.InitVaccountNo(accountNo.String, constants.CURRENCY_KHR, constants.VaType_FREEZE_KHR_DEBIT); vaccNo == "" {
			ss_log.Error("初始化个人虚账失败，accountNo=%v, vaccType=%v", accountNo.String, constants.VaType_FREEZE_KHR_DEBIT)
			//return err
		}
	}
	return accountNo.String
}

// 获取账号通过的实名认证姓名
func (AccDao) GetAuthNameFromUid(uid string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var authName sql.NullString
	sqlStr := "select am.auth_name from account a LEFT JOIN auth_material am " +
		"ON a.individual_auth_material_no = am.auth_material_no where a.uid = $1 " +
		"and a.individual_auth_status = 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&authName}, uid)
	if err != nil {
		return "", err
	}
	return authName.String, nil
}

func (*AccDao) GetAccountFromAccNo(accountNo string) (account string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account from account where uid=$1 limit 1`, []*sql.NullString{&accountT}, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountT.String
}

func (*AccDao) GetAccountByAccNo(accountNo string) (*AccDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var account, accType, identityNo sql.NullString
	sqlStr := "select acc.account, ra.account_type, ra.iden_no " +
		"from account acc " +
		"left join rela_acc_iden ra on ra.account_no = acc.uid " +
		"where acc.uid=$1 and acc.use_status=$2 limit 1"
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&account, &accType, &identityNo},
		accountNo, constants.AccountUseStatusNormal)
	if nil != err {
		return nil, err
	}

	obj := new(AccDao)
	obj.AccountNo = accountNo
	obj.Account = account.String
	obj.AccountType = accType.String
	obj.IdentityNo = identityNo.String

	return obj, nil
}
