package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessDao struct{}

var BusinessDaoInst BusinessDao

func (*BusinessDao) GetAccNoByBusinessNo(businessNo string) (accNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select account_no from business where business_no = $1 "
	var accNoT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accNoT}, businessNo); err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return accNoT.String
}
func (*BusinessDao) GetBusinessByAccNo(accNo string) (businessNo string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select business_no from business where account_no = $1 "
	var businessNoT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&businessNoT}, accNo); err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}

	return businessNoT.String
}

type BusinessStatus struct {
	AccountNo           string
	AccountStatus       string
	BusinessNo          string
	BusinessStatus      string
	IncomeAuthorization string
	OutgoAuthorization  string
	FullName            string
	BusinessType        string
}

func (*BusinessDao) GetBusinessStatusInfo(account, accountNo string) (*BusinessStatus, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT acc.uid, acc.use_status, b.business_no, b.use_status, b.income_authorization," +
		" b.outgo_authorization, b.full_name, b.business_type " +
		"FROM business b " +
		"LEFT JOIN account acc ON acc.uid = b.account_no "
	whereStr := "WHERE acc.is_delete='0' AND b.is_delete='0' "
	var arg string
	if account != "" {
		whereStr += "AND acc.account=$1 "
		arg = account
	} else if accountNo != "" {
		whereStr += "AND b.account_no=$1 "
		arg = accountNo
	}

	var accNo, accStatus, businessNo, businessStatus, incomeAuth, outgoAuth, fullName, businessType sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr+whereStr,
		[]*sql.NullString{&accNo, &accStatus, &businessNo, &businessStatus, &incomeAuth, &outgoAuth, &fullName, &businessType}, arg,
	)
	if err != nil {
		return nil, err
	}

	info := new(BusinessStatus)
	info.AccountNo = accNo.String
	info.AccountStatus = accStatus.String
	info.BusinessNo = businessNo.String
	info.BusinessStatus = businessStatus.String
	info.IncomeAuthorization = incomeAuth.String
	info.OutgoAuthorization = outgoAuth.String
	info.FullName = fullName.String
	info.BusinessType = businessType.String

	return info, nil
}
