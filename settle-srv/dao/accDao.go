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
	"a.a/mp-server/common/ss_sql"
)

var (
	AccDaoInstance AccDao
)

type AccDao struct {
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

func (r *AccDao) InsertEmptyAccount(phone string) (accountNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	genKey := util.RandomDigitStrOnlyNum(10)
	accountNoT := strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into account(uid,is_actived,create_time,phone,gen_key,account)values($1,$2,current_timestamp,$3,$4,$5)`,
		accountNoT, 0, phone, genKey, phone)
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
	ss_log.Error("err=%v", err)
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
	ss_log.Error("err=%v", err)
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
	err := ss_sql.QueryRow(dbHandler, `select khr_balance,usd_balance from account where uid=$1 limit 1`, []*sql.NullString{&khrT, &usdT}, accountNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "0", "0"
	}
	return khrT.String, usdT.String
}

func (*AccDao) GetAccNoFromPhone(phone string) (accNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select uid from account where phone=$1 and is_delete='0' limit 1`, []*sql.NullString{&accountNo}, phone)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return accountNo.String
}
func (*AccDao) GetIsActiveAccNoFromPhone(phone string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo, isActived sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select uid,is_actived from account where phone=$1 and is_delete='0'  limit 1`, []*sql.NullString{&accountNo, &isActived}, phone)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return accountNo.String, isActived.String
}
func (*AccDao) GetIsActiveFromPhone(phone string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var isActive sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select is_actived from account where phone=$1 and is_delete='0' limit 1`, []*sql.NullString{&isActive}, phone)
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

// 修改正式账号USD的余额
//func (*AccDao) UpdateUSDRemain(accountNo, amount, op string) string {
//	dbHandler := db.GetDB(constants.DB_CRM)
//	defer db.PutDB(constants.DB_CRM, dbHandler)
//	tx, _ := dbHandler.BeginTx(ctx, nil)
//	defer ss_sql.Rollback(tx)
//
//	switch op {
//	case "+":
//		err := ss_sql.ExecTx(tx, `update account set usd_balance=usd_balance+$1 where uid=$2 and is_delete='0'`, amount, accountNo)
//		if nil != err {
//			ss_log.Error("err=%v", err)
//			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//		}
//	case "-":
//		err := ss_sql.ExecTx(tx, `update account set usd_balance=usd_balance-$1 where uid=$2 and is_delete='0'`, amount, accountNo)
//		if nil != err {
//			ss_log.Error("err=%v", err)
//			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//		}
//	default:
//		return ss_err.ERR_PAY_VACCOUNT_OP_MISSING
//	}
//
//	descroption := "USD金额变动" + op + amount
//	errCode := LogAccDaoInstance.InsertAccountLog(descroption, accountNo, constants.REMAINTYPE)
//	if errCode != ss_err.ERR_SUCCESS {
//		ss_log.Error("err=%v", errCode)
//		return errCode
//	}
//
//	var tmp sql.NullString
//	err := ss_sql.QueryRowTx(tx, `select usd_balance from account where uid=$1 and is_delete='0' limit 1`,
//		[]*sql.NullString{&tmp}, accountNo)
//	if nil != err {
//		ss_log.Error("err=%v", err)
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//	if strext.ToInt64(tmp.String) < 0 {
//		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
//	}
//
//	ss_sql.Commit(tx)
//	return ss_err.ERR_SUCCESS
//}

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
//			ss_log.Error("err=%v", err)
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
//		ss_log.Error("err=%v", err)
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
