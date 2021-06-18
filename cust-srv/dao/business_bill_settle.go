package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessBillSettle struct {
	SettleId        string
	StartTime       string
	EndTime         string
	TotalAmount     string
	TotalRealAmount string
	TotalFees       string
	TotalOrderNum   string
	CurrencyType    string
	CreateTime      string
	BusinessName    string
	Account         string
	AppName         string
}

var BusinessBillSettleDaoInst BusinessBillSettle

func (*BusinessBillSettle) GetSettleLogById(settleId string) (*BusinessBillSettle, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT bbs.settle_id, bbs.start_time, bbs.end_time, bbs.total_amount, bbs.total_real_amount, bbs.total_fees, 
		bbs.total_order, bbs.currency_type,bbs.create_time, b.full_name, acc.account, app.app_name
		FROM business_bill_settle bbs
		LEFT JOIN business b ON b.business_no=bbs.business_no
		LEFT JOIN account acc ON acc.uid=b.account_no
		LEFT JOIN business_app app ON app.app_id=bbs.app_id 
		WHERE bbs.settle_id=$1`

	var settleIdT, startTime, endTime, totalAmount, totalRealAmount, totalFees, totalOrderNum, currencyType,
		createTime, businessName, account, appName sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&settleIdT, &startTime, &endTime, &totalAmount, &totalRealAmount,
		&totalFees, &totalOrderNum, &currencyType, &createTime, &businessName, &account, &appName}, settleId)
	if err != nil {
		return nil, err
	}

	data := new(BusinessBillSettle)
	data.SettleId = settleIdT.String
	data.StartTime = startTime.String
	data.EndTime = endTime.String
	data.TotalAmount = totalAmount.String
	data.TotalRealAmount = totalRealAmount.String
	data.TotalFees = totalFees.String
	data.TotalOrderNum = totalOrderNum.String
	data.CurrencyType = currencyType.String
	data.CreateTime = createTime.String
	data.BusinessName = businessName.String
	data.Account = account.String
	data.AppName = appName.String

	return data, nil
}

func (*BusinessBillSettle) GetSingleSettleDetail(settleId string) (*BusinessBillSettle, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bo.settle_id, bo.create_time, bb.amount, bb.real_amount, bb.currency_type, bb.fee," +
		"b.full_name, acc.account, app.app_name " +
		"FROM business_bill_settle_one  bo " +
		"LEFT JOIN business_bill bb ON bb.settle_id = bo.settle_id " +
		"LEFT JOIN business b ON b.business_no = bb.business_no " +
		"LEFT JOIN account acc ON acc.uid = b.account_no " +
		"LEFT JOIN business_app app ON app.app_id = bb.app_id " +
		"WHERE bb.settle_id=$1 "

	var settleIdT, createTime, amount, realAmount, currencyType, fee, businessName, account, appName sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{
		&settleIdT, &createTime, &amount, &realAmount, &currencyType, &fee, &businessName, &account, &appName,
	}, settleId)
	if err != nil {
		return nil, err
	}

	data := new(BusinessBillSettle)
	data.SettleId = settleIdT.String
	data.CreateTime = createTime.String
	data.TotalAmount = amount.String
	data.TotalRealAmount = realAmount.String
	data.TotalFees = fee.String
	data.TotalOrderNum = "1"
	data.CurrencyType = currencyType.String
	data.CreateTime = createTime.String
	data.BusinessName = businessName.String
	data.Account = account.String
	data.AppName = appName.String

	return data, nil
}
