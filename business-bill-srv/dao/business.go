package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var BusinessDaoInst BusinessDao

type BusinessDao struct {
	BusinessNo          string
	BusinessAccNo       string
	BusinessAcc         string
	IsEnabled           bool
	AccountNo           string
	SimplifyName        string
	IncomeAuthorization int
	OutgoAuthorization  int
}

// 查询交易配置
func (b *BusinessDao) GetTransConfig(businessNo string) (*BusinessDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var businessNoT, simplifyName, useStatus, businessAccNo, incomeAuth, outgoAuth, account sql.NullString
	sqlStr := "SELECT b.business_no, b.simplify_name, b.use_status, b.account_no, b.income_authorization, " +
		"b.outgo_authorization, acc.account " +
		"FROM business b " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"WHERE b.business_no=$1 AND acc.use_status=$2 AND b.is_delete = 0 LIMIT 1"
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&businessNoT, &simplifyName, &useStatus, &businessAccNo, &incomeAuth,
		&outgoAuth, &account}, businessNo, constants.AccountUseStatusNormal)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessDao)
	obj.BusinessNo = businessNoT.String
	obj.BusinessAccNo = businessAccNo.String
	obj.SimplifyName = simplifyName.String
	obj.BusinessAcc = account.String
	obj.IsEnabled = strext.ToInt(useStatus.String) == constants.BusinessUseStatusEnable
	obj.IncomeAuthorization = strext.ToInt(incomeAuth.String)
	obj.OutgoAuthorization = strext.ToInt(outgoAuth.String)

	return obj, err
}

func (b *BusinessDao) GetTransConfigByAccNo(accNo, businessType string) (*BusinessDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var businessNoT, simplifyName, useStatus, incomeAuth, outgoAuth, account sql.NullString
	sqlStr := "SELECT b.business_no, b.simplify_name, b.use_status, b.income_authorization, " +
		"b.outgo_authorization, acc.account " +
		"FROM business b " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"WHERE b.account_no=$1 AND b.business_type=$2 AND acc.use_status=$3 AND b.is_delete = 0 LIMIT 1"
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&businessNoT, &simplifyName, &useStatus, &incomeAuth,
		&outgoAuth, &account}, accNo, businessType, constants.AccountUseStatusNormal)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessDao)
	obj.BusinessNo = businessNoT.String
	obj.BusinessAccNo = accNo
	obj.SimplifyName = simplifyName.String
	obj.BusinessAcc = account.String
	obj.IsEnabled = strext.ToInt(useStatus.String) == constants.BusinessUseStatusEnable
	obj.IncomeAuthorization = strext.ToInt(incomeAuth.String)
	obj.OutgoAuthorization = strext.ToInt(outgoAuth.String)

	return obj, err
}

func (BusinessDao) QueryAccByBusinessNo(businessNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT acc.account " +
		"FROM business b " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"WHERE b.business_no=$1 and acc.use_status=$2 and b.is_delete=0  LIMIT 1 "
	var account sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&account}, businessNo, constants.AccountUseStatusNormal)
	return account.String, err
}

func (BusinessDao) QueryAccNoByBusinessNo(businessNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accountNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT account_no FROM business WHERE business_no=$1 and is_delete = 0 LIMIT 1",
		[]*sql.NullString{&accountNo}, businessNo)
	return accountNo.String, err
}

func (BusinessDao) QueryNameByBusinessNo(businessNo string) (string, string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var fullName, simplifyName sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT full_name, simplify_name FROM business WHERE business_no=$1 and is_delete = 0 LIMIT 1",
		[]*sql.NullString{&fullName, &simplifyName}, businessNo)
	return fullName.String, simplifyName.String, err
}
