package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessFixedCodeDao struct {
}

var BusinessFixedCodeDaoInst BusinessFixedCodeDao

//添加个人商家固定二维码
func (BusinessFixedCodeDao) AddBusinessFixedCodeTx(tx *sql.Tx, businessAccountNo, businessNo string) (staticCode string, err error) {
	//创建运营商
	code := util.RandomDigitStr(32)
	sqlStr := "insert into business_fixed_code(fixed_code, business_account_no, business_no, create_time) " +
		" values ($1,$2,$3,CURRENT_TIMESTAMP)"
	err = ss_sql.ExecTx(tx, sqlStr, code, businessAccountNo, businessNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	return code, nil
}

func (BusinessFixedCodeDao) GetFixedCodeByBusinessNo(businessNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT fixed_code FROM business_fixed_code WHERE business_no = $1 "

	var fixedCode sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&fixedCode}, businessNo); err != nil {
		return "", err
	}

	return fixedCode.String, nil
}
