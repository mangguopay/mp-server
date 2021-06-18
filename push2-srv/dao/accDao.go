package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	AccDaoInst AccDao
)

type AccDao struct {
}

// 获取账号语言
func (r *AccDao) GetAccLang(accNo, accountType string) (langStr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select app_lang,pos_lang from account where uid=$1 limit 1"
	var appLang, posLang sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&appLang, &posLang}, accNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return constants.LangEnUS
	}

	if appLang.String == "" {
		appLang.String = constants.LangEnUS
	}
	if posLang.String == "" {
		posLang.String = constants.LangEnUS
	}

	switch accountType {
	case constants.AccountType_USER:
		langStr = appLang.String
	case constants.AccountType_POS:
		fallthrough
	case constants.AccountType_OPERATOR:
		langStr = posLang.String
	}

	return langStr
}

// 获取账号
func (r *AccDao) GetAcc(accNo string) (account string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select account from account where uid=$1 limit 1"
	var accountT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accountT}, accNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return accountT.String
}

// 获取电话
func (r *AccDao) GetPhone(accNo string) (countryCode, phone string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select phone,country_code from account where uid=$1 limit 1"
	var phoneT, countryCodeT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&phoneT, &countryCodeT}, accNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return countryCodeT.String, phoneT.String
	}

	return countryCodeT.String, phoneT.String
}

// 获取email
func (r *AccDao) GetEmail(accNo string) (emailAddr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select email from account where uid=$1 limit 1"
	var emailAddrT sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&emailAddrT}, accNo)
	if nil != err {
		ss_log.Error("err=[%v]", err)
		return emailAddrT.String
	}

	return emailAddrT.String
}
