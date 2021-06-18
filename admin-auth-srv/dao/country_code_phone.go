package dao

import (
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type CountryCodePhoneDao struct{}

var CountryCodePhoneDaoInst CountryCodePhoneDao

func (CountryCodePhoneDao) Insert(tx *sql.Tx, countryCode, phone string) error {
	return ss_sql.ExecTx(tx, `insert into country_code_phone(country_code,phone)VALUES ($1,$2)`, countryCode, phone)
}

func (CountryCodePhoneDao) Delete(tx *sql.Tx, countryCode, phone string) error {
	return ss_sql.ExecTx(tx, `delete from country_code_phone where country_code = $1 and phone = $2`, countryCode, phone)
}
