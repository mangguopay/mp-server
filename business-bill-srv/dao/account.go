package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type AccountDao struct {
	AccountNo           string
	Account             string
	AccountType         string
	UseStatus           string
	Phone               string
	Email               string
	IncomeAuthorization int
}

var AccountDaoInst AccountDao

func (*AccountDao) GetAccountById(uid string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT account FROM account WHERE uid=$1 "
	var account sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&account}, uid); nil != err {
		return "", err
	}
	return account.String, nil
}

func (*AccountDao) GetInfoByAccount(account string) (*AccountDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT acc.uid, acc.phone, acc.email, acc.use_status, rel.account_type, cu.in_transfer_authorization " +
		"FROM account acc " +
		"LEFT JOIN rela_acc_iden rel ON rel.account_no = acc.uid " +
		"LEFT JOIN cust cu ON cu.account_no = acc.uid " +
		"WHERE acc.account=$1 AND acc.use_status=$2 "

	var accountNo, phone, email, useStatus, accountType, incomeAuthorization sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accountNo, &phone, &email, &useStatus, &accountType,
		&incomeAuthorization},
		account, constants.AccountUseStatusNormal)
	if nil != err {
		return nil, err
	}

	obj := new(AccountDao)
	obj.AccountNo = accountNo.String
	obj.Account = account
	obj.AccountType = accountType.String
	obj.Phone = phone.String
	obj.Email = email.String
	obj.UseStatus = useStatus.String
	obj.IncomeAuthorization = strext.ToInt(incomeAuthorization.String)

	return obj, nil
}

//获取用户的语言
func (*AccountDao) QueryAccountLang(accountNo string) (string, string) {
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
