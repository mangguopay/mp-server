package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessBatchTransferDao struct {
	BatchNo          string
	TotalNumber      string
	TotalAmount      string
	SuccessfulNumber string
	SuccessfulAmount string
	FailNumber       string
	FailAmount       string
	ProcessingNumber string
	//Status           string
	PaymentType  string
	Remarks      string
	BusinessNo   string
	CreateTime   string
	FinishTime   string
	UpdateTime   string
	FileContent  string
	GenerateAll  string
	CurrencyType string
	RealAmount   string
}

var BusinessBatchTransferDaoInst BusinessBatchTransferDao

func (BusinessBatchTransferDao) GetOrderList(whereStr string, whereArgs []interface{}) ([]*go_micro_srv_cust.BusinessTransferBatchData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT btb.batch_no, btb.total_number, btb.total_amount, btb.successful_number, btb.successful_amount," +
		" btb.fail_number, btb.fail_amount, btb.processing_number, btb.status, btb.payment_type," +
		" btb.remarks, btb.business_no, btb.create_time , btb.finish_time, btb.update_time," +
		" btb.currency_type, acc.account, btb.real_amount " +
		" FROM business_transfer_batch_order btb" +
		" LEFT JOIN business bu ON bu.business_no = btb.business_no " +
		" LEFT JOIN account acc ON acc.uid = bu.account_no  "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr+whereStr, whereArgs...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var datas []*go_micro_srv_cust.BusinessTransferBatchData
	for rows.Next() {
		var (
			remarks, finishTime, updateTime sql.NullString
		)
		data := new(go_micro_srv_cust.BusinessTransferBatchData)

		err := rows.Scan(
			&data.BatchNo,
			&data.TotalNumber,
			&data.TotalAmount,
			&data.SuccessfulNumber,
			&data.SuccessfulAmount,

			&data.FailNumber,
			&data.FailAmount,
			&data.ProcessingNumber,
			&data.Status,
			&data.PaymentType,

			&remarks,
			&data.BusinessNo,
			&data.CreateTime,
			&finishTime,
			&updateTime,

			&data.CurrencyType,
			&data.Account,
			&data.RealAmount,
		)
		if err != nil {
			return nil, err
		}
		data.Remarks = remarks.String
		data.FinishTime = finishTime.String
		data.UpdateTime = updateTime.String
		datas = append(datas, data)
	}

	return datas, nil
}

func (BusinessBatchTransferDao) CountOrderNum(whereStr string, whereArgs []interface{}) (int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT COUNT(1) " +
		" FROM business_transfer_batch_order btb " +
		" LEFT JOIN business bu ON bu.business_no = btb.business_no " +
		" LEFT JOIN account acc ON acc.uid = bu.account_no  "
	var num sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr+whereStr, []*sql.NullString{&num}, whereArgs...)
	if err != nil {
		return -1, err
	}

	return strext.ToInt32(num.String), nil
}

func (BusinessBatchTransferDao) InsertOrder(data BusinessBatchTransferDao) (batchNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	id := strext.GetDailyId()
	sqlStr := "INSERT INTO business_transfer_batch_order(batch_no, total_number, total_amount, status, " +
		" remarks, business_no, file_content, generate_all, currency_type," +
		" real_amount, create_time)" +
		" VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,CURRENT_TIMESTAMP) "
	errT := ss_sql.Exec(dbHandler, sqlStr,
		id, data.TotalNumber, data.TotalAmount, constants.BusinessBatchTransferOrderStatusPending,
		data.Remarks, data.BusinessNo, data.FileContent, data.GenerateAll, data.CurrencyType, data.RealAmount)
	if errT != nil {
		return "", errT
	}

	return id, nil
}

func (BusinessBatchTransferDao) GetBatchOrderDetail(whereStr string, whereArgs []interface{}) (*go_micro_srv_cust.BusinessTransferBatchData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT btb.batch_no, btb.total_number, btb.total_amount, btb.successful_number, btb.successful_amount," +
		" btb.fail_number, btb.fail_amount, btb.processing_number, btb.status, btb.payment_type," +
		" btb.remarks, btb.business_no, btb.create_time , btb.finish_time, btb.update_time," +
		" acc.account, btb.real_amount, btb.currency_type, btb.fee  " +
		" FROM business_transfer_batch_order btb" +
		" LEFT JOIN business bu ON bu.business_no = btb.business_no " +
		" LEFT JOIN account acc ON acc.uid = bu.account_no  "
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr+whereStr, whereArgs...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}

	var (
		batchNo, totalNumber, totalAmount, successfulNumber, successfulAmount,
		failNumber, failAmount, processingNumber, status, paymentType,
		remarks, businessNo, createTime, finishTime, updateTime,
		account, realAmount, currencyType, fee sql.NullString
	)
	err = rows.Scan(
		&batchNo,
		&totalNumber,
		&totalAmount,
		&successfulNumber,
		&successfulAmount,

		&failNumber,
		&failAmount,
		&processingNumber,
		&status,
		&paymentType,

		&remarks,
		&businessNo,
		&createTime,
		&finishTime,
		&updateTime,
		&account,
		&realAmount,
		&currencyType,
		&fee,
	)
	if err != nil {
		return nil, err
	}
	data := new(go_micro_srv_cust.BusinessTransferBatchData)
	data.BatchNo = batchNo.String
	data.TotalNumber = totalNumber.String
	data.TotalAmount = totalAmount.String
	data.SuccessfulNumber = successfulNumber.String
	data.SuccessfulAmount = successfulAmount.String
	data.FailNumber = failNumber.String
	data.FailAmount = failAmount.String
	data.ProcessingNumber = processingNumber.String
	data.Status = status.String
	data.PaymentType = paymentType.String
	data.Remarks = remarks.String
	data.BusinessNo = businessNo.String
	data.CreateTime = createTime.String
	data.FinishTime = finishTime.String
	data.UpdateTime = updateTime.String
	data.Account = account.String
	data.RealAmount = realAmount.String
	data.CurrencyType = currencyType.String
	data.Fee = fee.String

	return data, nil
}
