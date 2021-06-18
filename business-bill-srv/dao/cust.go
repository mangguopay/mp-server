package dao

import (
	"database/sql"
	"fmt"

	"a.a/cu/strext"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type CustDao struct {
}

var CustDaoInst CustDao

type CustInfo struct {
	CustNo               string
	AccountNo            string
	PaymentPassword      string
	IndividualAuthStatus string
	TradingAuthority     int
	IncomeAuthorization  int
	OutgoAuthorization   int
}

func (*CustDao) QueryCustInfo(accountNo, account string) (*CustInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return nil, GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereStr := "where 1=1 "
	var i = 1
	if accountNo != "" {
		whereStr += fmt.Sprintf("and cu.account_no=$%v", i)
		i++
	}
	if account != "" {
		whereStr += fmt.Sprintf("and acc.account=$%v", i)
	}

	sqlStr := "select cu.cust_no, cu.trading_authority, cu.payment_password, acc.individual_auth_status, " +
		"cu.in_transfer_authorization, cu.out_transfer_authorization, acc.uid " +
		"from cust cu " +
		"left join account acc on acc.uid = cu.account_no "
	sqlStr += whereStr

	var custNo, tradingAuthority, paymentPwd, individualAuthStatus, inAuthorization, outAuthorization, accNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&custNo, &tradingAuthority, &paymentPwd, &individualAuthStatus, &inAuthorization,
			&outAuthorization, &accNo}, accountNo)
	if nil != err {
		return nil, err
	}

	obj := new(CustInfo)
	obj.CustNo = custNo.String
	obj.AccountNo = accNo.String
	obj.TradingAuthority = strext.ToInt(tradingAuthority.String)
	obj.PaymentPassword = paymentPwd.String
	obj.IndividualAuthStatus = individualAuthStatus.String
	obj.IncomeAuthorization = strext.ToInt(inAuthorization.String)
	obj.OutgoAuthorization = strext.ToInt(outAuthorization.String)

	return obj, nil
}

func (*CustDao) UpdateTradingAuthority(custNo string, tradingAuthority int) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "update cust set trading_authority=$1 where cust_no=$2"
	return ss_sql.Exec(dbHandler, sqlStr, tradingAuthority, custNo)
}

func (*CustDao) QueryCustNo(accountNo string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		return "", GetDBConnectFailedErr
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var custNo sql.NullString
	sqlStr := "select cust_no from cust where account_no=$1 limit 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&custNo}, accountNo)
	if nil != err {
		return "", err
	}
	return custNo.String, nil
}
