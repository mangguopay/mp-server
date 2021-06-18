package dao

import (
	"database/sql"
	"errors"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var (
	CustDaoInstance CustDao
)

type CustDao struct {
}

func (*CustDao) QueryPwdFromOpAccNo(opAccNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pwdT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select payment_password from cust where cust_no=$1 limit 1", []*sql.NullString{&pwdT}, opAccNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return pwdT.String
}

func (*CustDao) QueryRateRoleFrom(custNo string) (string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var tradingAuthority, inAuthorizationT, outAuthorizationT, inTransferAuthorizationT, outTransferAuthorizationT sql.NullString
	err := ss_sql.QueryRow(dbHandler, "select trading_authority,in_authorization,out_authorization,in_transfer_authorization,out_transfer_authorization"+
		" from cust where cust_no=$1 limit 1", []*sql.NullString{&tradingAuthority, &inAuthorizationT, &outAuthorizationT, &inTransferAuthorizationT, &outTransferAuthorizationT}, custNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", ""
	}
	return tradingAuthority.String, inAuthorizationT.String, outAuthorizationT.String, inTransferAuthorizationT.String, outTransferAuthorizationT.String
}

type CustInfo struct {
	CustNo               string
	TradingAuthority     string
	IndividualAuthStatus string
}

func (*CustDao) QueryCustInfo(accountNo string) (*CustInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select cu.cust_no, cu.trading_authority, acc.individual_auth_status " +
		"from cust cu " +
		"left join account acc on acc.uid = cu.account_no " +
		"where cu.account_no=$1 limit 1"
	var custNo, tradingAuthority, individualAuthStatus sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&custNo, &tradingAuthority, &individualAuthStatus}, accountNo)
	if nil != err {
		return nil, err
	}

	custInfo := new(CustInfo)
	custInfo.CustNo = custNo.String
	custInfo.TradingAuthority = tradingAuthority.String
	custInfo.IndividualAuthStatus = individualAuthStatus.String
	return custInfo, nil
}

func (CustDao) UpdateTradingAuthority(custNo string, tradingAuthority int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	return ss_sql.Exec(dbHandler, `update cust set trading_authority=$1 where cust_no=$2`,
		tradingAuthority, custNo)
}
