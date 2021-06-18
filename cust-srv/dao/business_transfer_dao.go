package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessTransferDao struct {
	LogNo               string
	FromBusinessAccount string
	FromBusinessName    string
	PayeeAccount        string
	ToBusinessName      string
	Amount              string
	CurrencyType        string
	Fee                 string
	PaymentType         string
	RealAmount          string
	OrderStatus         string
	BatchNo             string
	Remarks             string
	CreateTime          string
	WrongReason         string
	AuthName            string
	ToAccount           string //收款账号（批量转账产生的转账才有）
	TransferType        string
	OutTransferNo       string
}

var BusinessTransferDaoInst BusinessTransferDao

func (BusinessTransferDao) GetOrderList(whereStr string, whereArgs []interface{}) ([]*BusinessTransferDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bto.log_no, bto.amount, bto.fee, bto.real_amount, bto.currency_type, bto.create_time, " +
		"bto.order_status, bto.remarks, acc.account, bto.wrong_reason, bto.auth_name, bto.to_account, bto.transfer_type, bto.out_transfer_no " +
		"FROM business_transfer_order bto " +
		"LEFT JOIN account acc ON acc.uid=bto.to_account_no "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr+whereStr, whereArgs...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var list []*BusinessTransferDao
	for rows.Next() {
		var logNo, amount, fee, realAmount, currencyType, createTime, orderStatus, remarks, account,
			wrongReason, authName, toAccount, transferType, outTransferNo sql.NullString
		err := rows.Scan(&logNo, &amount, &fee, &realAmount, &currencyType, &createTime, &orderStatus, &remarks, &account,
			&wrongReason, &authName, &toAccount, &transferType, &outTransferNo)
		if err != nil {
			return nil, err
		}
		log := new(BusinessTransferDao)
		log.LogNo = logNo.String
		log.Amount = amount.String
		log.Fee = fee.String
		log.RealAmount = realAmount.String
		log.CurrencyType = currencyType.String
		log.CreateTime = createTime.String
		log.OrderStatus = orderStatus.String
		log.Remarks = remarks.String
		log.PayeeAccount = account.String
		log.WrongReason = wrongReason.String
		log.AuthName = authName.String
		log.ToAccount = toAccount.String
		log.TransferType = transferType.String
		log.OutTransferNo = outTransferNo.String
		list = append(list, log)
	}

	return list, nil
}

func (BusinessTransferDao) CountOrderNum(whereStr string, whereArgs []interface{}) (int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT COUNT(1) " +
		"FROM business_transfer_order bto " +
		"LEFT JOIN account acc ON acc.uid=bto.to_account_no "
	var num sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr+whereStr, []*sql.NullString{&num}, whereArgs...)
	if err != nil {
		return -1, err
	}

	return strext.ToInt32(num.String), nil
}

func (BusinessTransferDao) GetOrderDetail(logNo string) (*BusinessTransferDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT b1.full_name AS from_business_name, b2.full_name AS to_business_name," +
		"acc1.account AS from_account, acc2.account AS to_account," +
		"bt.log_no, bt.amount, bt.fee, bt.real_amount, bt.currency_type, bt.create_time, bt.order_status, bt.remarks, bt.batch_no," +
		"bt.transfer_type " +
		"FROM business_transfer_order bt " +
		"LEFT JOIN business b1 ON b1.account_no = bt.from_account_no " +
		"LEFT JOIN business b2 ON b2.account_no = bt.to_account_no " +
		"LEFT JOIN account acc1 ON acc1.uid = bt.from_account_no " +
		"LEFT JOIN account acc2 ON acc2.uid = bt.to_account_no " +
		"WHERE log_no=$1 "

	var (
		fromBusinessName, toBusinessName, fromAccount, toAccount, orderNo, amount, fee,
		realAmount, currencyType, createTime, orderStatus, remarks, batchNo, transferType sql.NullString
	)
	var err = ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{
		&fromBusinessName, &toBusinessName, &fromAccount, &toAccount, &orderNo, &amount, &fee,
		&realAmount, &currencyType, &createTime, &orderStatus, &remarks, &batchNo, &transferType},
		logNo,
	)
	if err != nil {
		return nil, err
	}

	order := new(BusinessTransferDao)
	order.LogNo = orderNo.String
	order.FromBusinessName = fromBusinessName.String
	order.ToBusinessName = toBusinessName.String
	order.FromBusinessAccount = fromAccount.String
	order.ToAccount = toAccount.String
	order.Amount = amount.String
	order.Fee = fee.String
	order.RealAmount = realAmount.String
	order.CurrencyType = currencyType.String
	order.CreateTime = createTime.String
	order.OrderStatus = orderStatus.String
	order.Remarks = remarks.String
	order.BatchNo = batchNo.String
	order.TransferType = transferType.String

	return order, nil
}

//管理后台获取
func (BusinessTransferDao) GetCnt(whereStr string, whereArgs []interface{}) (int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT COUNT(1) " +
		"FROM business_transfer_order bto " +
		"LEFT JOIN account acc ON acc.uid = bto.to_account_no " +
		"LEFT JOIN account acc2 ON acc2.uid = bto.from_account_no "

	var num sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr+whereStr, []*sql.NullString{&num}, whereArgs...)
	if err != nil {
		return -1, err
	}

	return strext.ToInt32(num.String), nil
}

func (BusinessTransferDao) GetTransferOrderList(whereStr string, whereArgs []interface{}) ([]*custProto.BusinessTransferOrderData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bto.log_no, bto.amount, bto.fee, bto.real_amount, bto.currency_type, bto.create_time, " +
		" bto.order_status, bto.remarks, acc.account, acc2.account, bto.rate," +
		" bto.batch_no, bto.batch_row_num, bto.wrong_reason, bto.transfer_type " +
		"FROM business_transfer_order bto " +
		"LEFT JOIN account acc ON acc.uid = bto.to_account_no " +
		"LEFT JOIN account acc2 ON acc2.uid = bto.from_account_no "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr+whereStr, whereArgs...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var datas []*custProto.BusinessTransferOrderData
	for rows.Next() {
		data := new(custProto.BusinessTransferOrderData)
		var toAccount, batchNo, batchRowNum, wrongReason, transferType sql.NullString
		err := rows.Scan(&data.LogNo, &data.Amount, &data.Fee, &data.RealAmount, &data.CurrencyType,
			&data.CreateTime, &data.OrderStatus, &data.Remarks, &toAccount, &data.FromAccount, &data.Rate,
			&batchNo, &batchRowNum, &wrongReason, &transferType,
		)
		if err != nil {
			return nil, err
		}
		data.ToAccount = toAccount.String
		data.BatchNo = batchNo.String
		data.BatchRowNum = batchRowNum.String
		data.WrongReason = wrongReason.String
		data.TransferType = transferType.String

		datas = append(datas, data)
	}

	return datas, nil
}
