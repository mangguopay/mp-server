package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

var (
	AccDaoInstance AccDao
)

type AccDao struct {
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

func (*AccDao) GetAccByNickname(nickname string) *adminAuthProto.Account {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	account := adminAuthProto.Account{}
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

func (*AccDao) GetRoleAuthedUrlList(accountNo string) (errCode string, datas []*adminAuthProto.RoleSimpleData) {
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

	datas = []*adminAuthProto.RoleSimpleData{}
	for rows.Next() {
		data := adminAuthProto.RoleSimpleData{}
		err := rows.Scan(&data.RoleNo, &data.RoleName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		datas = append(datas, &data)
	}
	return ss_err.ERR_SUCCESS, datas
}

// 初始化密钥
func (AccDao) InitPassword(password string) string {
	k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	return encrypt.DoMd5Salted(password, passwordSalt)
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

func (*AccDao) GetCountryCodePhoneByUidTx(tx *sql.Tx, uid string) (string, string, error) {
	var countryCode, phone sql.NullString
	err := ss_sql.QueryRowTx(tx, `select country_code,phone from account where uid=$1 and is_delete = '0' `,
		[]*sql.NullString{&countryCode, &phone}, uid)

	return countryCode.String, phone.String, err
}
