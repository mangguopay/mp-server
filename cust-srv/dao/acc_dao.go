package dao

import (
	"database/sql"
	"time"

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
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	AccDaoInstance AccDao
)

type AccDao struct {
}

func (*AccDao) UpdateAccountStatusByAccount(account string, status int) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update account set use_status= $3 where account=$1 and use_status = $2 `, account, constants.AccountUseStatusTemporaryDisabled, status)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
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

func (r *AccDao) UpdateAccountPwd(phone, pwd string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	pwdMD5 := r.InitPassword(pwd)
	err := ss_sql.Exec(dbHandler, `update account set password = $1 , modify_time = current_timestamp where phone = $2 and is_delete = 0`, pwdMD5, phone)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
func (r *AccDao) UpdateAccountPwdByUID(uid, oldPassword, newPassword string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	oldPwdMD5 := r.InitPassword(oldPassword)
	newPwdMD5 := r.InitPassword(newPassword)
	var dbPWD sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT password FROM account WHERE uid=$1 and is_delete = 0 LIMIT 1", []*sql.NullString{&dbPWD}, uid)
	if err != nil {
		ss_log.Error("err=[%v]", err.Error())
	}

	if dbPWD.String != oldPwdMD5 {
		ss_log.Error("err=[%v]", "原密码不正确")
		return ss_err.ERR_ACCOUNT_OLD_PWD_FAILD
	}

	// 修改密码
	if err := ss_sql.Exec(dbHandler, `update account set password = $1 , modify_time = current_timestamp,is_first_login = '1' where uid = $2 and is_delete = 0`, newPwdMD5, uid); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return ss_err.ERR_MODIFY_ACCOUNT_PWD_FAILD
	}

	return ss_err.ERR_SUCCESS
}

func (r *AccDao) QueryUIDByPhone(phone string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var uid sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT uid FROM account WHERE phone=$1 and is_delete = 0 LIMIT 1", []*sql.NullString{&uid}, phone)
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

//检查账户是否已存在,并返回uid
func (*AccDao) GetUidByAccountTx(tx *sql.Tx, account string) (uid string, err error) {
	var uidT sql.NullString
	errT := ss_sql.QueryRowTx(tx, "SELECT uid FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&uidT}, account)

	return uidT.String, errT
}

//检查账户是否已存在,并返回uid
func (*AccDao) GetUidByAccount(account string) (uid string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var uidT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, "SELECT uid FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&uidT}, account)

	return uidT.String, errT
}

//检查账户是否已存在
func (*AccDao) CheckAccountTx(tx *sql.Tx, account string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	return strext.ToInt(count.String), err
}

//检查账户是否已存在
func (*AccDao) CheckAccount(account string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	return strext.ToInt(count.String), err
}

func (*AccDao) GetAccountCnt(account string) (int, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var count sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT COUNT(1) FROM account WHERE account = $1 and is_delete='0'", []*sql.NullString{&count}, account)
	if err != nil {
		ss_log.Error("err=[%v]", err)
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
		ss_log.Error("err=[%v]", err)
		return false
	}
	return strext.ToInt(count.String) > 0
}

func (*AccDao) CheckAccountUpdate(tx *sql.Tx, account, accountUid string) (int, error) {
	var count sql.NullString
	err := ss_sql.QueryRowTx(tx, "SELECT COUNT(1) FROM account WHERE account = $1 and uid = $2 and is_delete = '0' ", []*sql.NullString{&count}, account, accountUid)
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
		"master_acc,phone,create_time,gen_key,is_actived) VALUES ($1,$2,$3,$4,$5,$6,$7,current_timestamp,$8,'1')",
		accountUid, nickname, _account, pwdMD5, useStatus, masterAccount, phone, util.RandomDigitStrOnlyNum(10))
	if nil != err {
		ss_log.Error("err=[%v]", err)
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

func (r *AccDao) AddCashier(tx *sql.Tx, name, servicer_no, op_password string) (string, error) {
	var pwdMD5 string
	if op_password != "" {
		pwdMD5 = r.InitPassword(op_password)
	}

	ss_log.Info("AccDao | AddCashier ")
	uid := strext.NewUUID()
	insertSql := "INSERT INTO cashier(uid, name, servicer_no, is_delete, create_time, op_password)" +
		" VALUES ($1,$2,$3,$4,current_timestamp,$5)"
	err := ss_sql.ExecTx(tx, insertSql, uid, name, servicer_no, "0", pwdMD5)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	ss_log.Info(" AccDao | AddCashier Uid=%v", uid)
	ss_log.Info(" AccDao | AddCashier pwdMD5=%v", pwdMD5)

	return uid, nil
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
		err := r.DeleteAccount(dbHandler, accNos.(string))
		if err != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[%v]", err)
			return
		}
		ss_log.Info("删除账户[%v]成功", accNos)
	case []string:
		for _, v := range accNos.([]string) {
			err := r.DeleteAccount(dbHandler, v)
			if err != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[%v]", err)
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

// 初始化密钥
func (AccDao) InitPassword(password string) string {
	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	return encrypt.DoMd5Salted(password, passwordSalt)
}

func (*AccDao) GetAccountByUid(dbHandler *sql.DB, uid string) string {
	var account sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account from account where uid=$1 `, []*sql.NullString{&account}, uid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return account.String
}

func (*AccDao) GetAccByUid(uid string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var account sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select account from account where uid=$1 `, []*sql.NullString{&account}, uid)
	if nil != err {
		return "", err
	}
	return account.String, nil
}

func (*AccDao) GetAccNoFromPhone(phone string) (accNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select uid from account where phone=$1 and is_delete='0' limit 1`, []*sql.NullString{&accountNo}, phone)
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

func (AccDao) UpdatePhone(accNo, phone string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update account set phone = $2,account = $2 where uid = $1 and is_delete = '0' "
	err := ss_sql.Exec(dbHandler, sqlStr, accNo, phone)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_SYS_DB_UPDATE
	}
	return ss_err.ERR_SUCCESS
}

func (AccDao) DeleteAccountTx(tx *sql.Tx, accNo string) string {
	sqlStr2 := "update account set is_delete = '1' where uid = $1 and is_delete = '0' "
	if err2 := ss_sql.ExecTx(tx, sqlStr2, accNo); err2 != nil {
		ss_log.Error("删除店员账号失败，err2=[%v]", err2)
		return ss_err.ERR_SYS_DB_DELETE
	}
	return ss_err.ERR_SUCCESS
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
func (*AccDao) GetBusinessPhoneFromAccNo(accNo string) (businessPhone string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var businessPhoneT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select business_phone from account where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&businessPhoneT}, accNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return businessPhoneT.String
}
func (*AccDao) GetEmailFromAccNo(accNo string) (email string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var emailT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select email from account where uid=$1 and is_delete='0' limit 1`, []*sql.NullString{&emailT}, accNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return emailT.String
}

// 获取实名认证的姓名
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

//获取账号的个人实名认证状态
func (AccDao) GetAuthStatusFromUid(uid string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var authStatus sql.NullString
	sqlStr := "select individual_auth_status " +
		" from account " +
		" where uid = $1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&authStatus}, uid)
	if err != nil {
		return "", err
	}
	return authStatus.String, nil
}

func (*AccDao) GetRegCountByDate(beforeDay string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return &DataCount{}, tErr
	}

	sqlStr := `select count(1) as totalCount from account 
	WHERE create_time >= $1 and create_time  < $2 and is_actived = 1 and  is_delete = 0 `
	var num sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num}, startTime, endTime)
	if errT != nil {
		return &DataCount{}, errT
	}
	return &DataCount{
		RegNum: strext.ToInt64(num.String),
		Day:    startTime,
	}, nil
}

// 用户余额、冻结金额 统计
func (*AccDao) GetUserMoneyCount() *UserMoneyDataCount {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select sum(balance),sum(frozen_balance)  from vaccount where va_type = $1 `
	var usdBalance, usdFrozenBalance sql.NullString
	if errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&usdBalance, &usdFrozenBalance}, constants.VaType_USD_DEBIT); errT != nil {
		usdBalance.String = "0"
		usdFrozenBalance.String = "0"
	}

	var khrBalance, khrFrozenBalance sql.NullString
	if errT2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&khrBalance, &khrFrozenBalance}, constants.VaType_KHR_DEBIT); errT2 != nil {
		khrBalance.String = "0"
		khrFrozenBalance.String = "0"
	}

	//未激活的账号也有钱在平台
	var usdBalance2, usdFrozenBalance2 sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&usdBalance2, &usdFrozenBalance2}, constants.VaType_FREEZE_USD_DEBIT)
	if errT != nil {
		usdBalance2.String = "0"
		usdFrozenBalance2.String = "0"
	}

	var khrBalance2, khrFrozenBalance2 sql.NullString
	if errT2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&khrBalance2, &khrFrozenBalance2}, constants.VaType_FREEZE_KHR_DEBIT); errT2 != nil {
		khrBalance2.String = "0"
		khrFrozenBalance2.String = "0"
	}

	data := &UserMoneyDataCount{
		UserUsdBalance:       strext.ToInt64(ss_count.Add(usdBalance.String, usdBalance2.String)),
		UserUsdFrozenBalance: strext.ToInt64(ss_count.Add(usdFrozenBalance.String, usdFrozenBalance2.String)),
		UserKhrBalance:       strext.ToInt64(ss_count.Add(khrBalance.String, khrBalance2.String)),
		UserKhrFrozenBalance: strext.ToInt64(ss_count.Add(khrFrozenBalance.String, khrFrozenBalance2.String)),
	}
	return data
}

func (*AccDao) GetAccountByUidTX(tx *sql.Tx, uid string) string {
	var account sql.NullString
	err := ss_sql.QueryRowTx(tx, `select account from account where uid=$1 `, []*sql.NullString{&account}, uid)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return account.String
}

type AccountAuthInfo struct {
	Account        string
	AuthMaterialNo string
	AuthName       string
	AuthNumber     string
	AuthStatus     string
}

type BusinessAuthInfo struct {
	Account      string
	BusinessName string
	SimplifyName string
	AuthStatus   string
	StartDate    string
	EndDate      string
	AuthNumber   string
	LicenseImgNo string
}

//查询商家认证状态
func (*AccDao) GetAuthInfoByAccountNo(accountNo, accountType string) (*AccountAuthInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var sqlStr string
	if accountType == constants.AccountType_PersonalBusiness {
		sqlStr = `SELECT status, auth_material_no, auth_name, auth_number
		FROM auth_material_business
		WHERE account_uid = $1 
		ORDER BY create_time DESC 
		LIMIT 1 `
	}

	if accountType == constants.AccountType_EnterpriseBusiness {
		sqlStr = `SELECT status, auth_material_no, auth_name, auth_number
		FROM auth_material_enterprise
		WHERE account_uid = $1 
		ORDER BY create_time DESC 
		LIMIT 1 `
	}

	var authStatus, materialNo, authName, authNumber sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&authStatus, &materialNo, &authName, &authNumber}, accountNo)
	if err != nil {
		return nil, err
	}

	info := new(AccountAuthInfo)
	info.AuthMaterialNo = materialNo.String
	info.AuthName = authName.String
	info.AuthNumber = authNumber.String
	info.AuthStatus = authStatus.String

	return info, nil
}

//查询普通用户账号基本信息
func (*AccDao) GetUserAccountInfo(accountNo string) (*AccountAuthInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT acc.account, au.auth_name, au.status " +
		"FROM account acc " +
		"LEFT JOIN auth_material au ON au.account_uid = acc.uid " +
		"WHERE acc.uid = $1 "

	var account, realName, authStatus sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&account, &realName, &authStatus}, accountNo); err != nil {
		return nil, err
	}

	obj := new(AccountAuthInfo)
	obj.Account = account.String
	obj.AuthName = realName.String
	obj.AuthStatus = authStatus.String

	return obj, nil
}

//查询个人商家账号基本信息
func (*AccDao) GetBusinessAccountInfo(accountNo string) (*BusinessAuthInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT acc.account, au.auth_name, au.simplify_name, au.status, au.start_date, au.end_date, " +
		"au.license_img_no, au.auth_number " +
		"FROM account acc " +
		"LEFT JOIN auth_material_business au ON au.account_uid = acc.uid " +
		"WHERE acc.uid = $1 "

	var account, fullName, simplifyName, authStatus, startDate, endDate, licenseImgNo, authNumber sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&account, &fullName, &simplifyName, &authStatus,
		&startDate, &endDate, &licenseImgNo, &authNumber}, accountNo)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessAuthInfo)
	obj.Account = account.String
	obj.BusinessName = fullName.String
	obj.SimplifyName = simplifyName.String
	obj.AuthStatus = authStatus.String
	obj.StartDate = startDate.String
	obj.EndDate = endDate.String
	obj.LicenseImgNo = licenseImgNo.String
	obj.AuthNumber = authNumber.String

	return obj, nil
}
