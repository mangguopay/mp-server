package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"errors"
)

type BusinessBatchTransferOrderDao struct {
	BatchNo          string
	TotalNumber      string
	TotalAmount      string
	SuccessfulNumber string
	SuccessfulAmount string
	FailNumber       string
	FailAmount       string
	ProcessingNumber string
	Status           string
	PaymentType      string
	Remarks          string
	BusinessNo       string
	CreateTime       string
	FinishTime       string
	UpdateTime       string
	FileContent      string
	GenerateAll      string //转账订单是否全部生成
	CurrencyType     string
	RealAmount       string
	Fee              string
}

var BusinessBatchTransferOrderDaoInst BusinessBatchTransferOrderDao

func (BusinessBatchTransferOrderDao) GetBatchOrderDetail(batchNo string) (*BusinessBatchTransferOrderDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT 
				batch_no, business_no, status, file_content, currency_type,
				real_amount, generate_all, total_number, fee 
			 FROM business_transfer_batch_order 
			 WHERE batch_no = $1 `
	row, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, batchNo)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}

	var batchNoT, businessNo, status, fileContent, currencyType, realAmount, generateAll, totalNumber, fee sql.NullString
	err = row.Scan(
		&batchNoT,
		&businessNo,
		&status,
		&fileContent,
		&currencyType,
		&realAmount,
		&generateAll,
		&totalNumber,
		&fee,
	)
	if err != nil {
		return nil, err
	}

	return &BusinessBatchTransferOrderDao{
		BatchNo:      batchNoT.String,
		BusinessNo:   businessNo.String,
		Status:       status.String,
		FileContent:  fileContent.String,
		CurrencyType: currencyType.String,
		RealAmount:   realAmount.String,
		GenerateAll:  generateAll.String,
		TotalNumber:  totalNumber.String,
		Fee:          fee.String,
	}, err

}

func (BusinessBatchTransferOrderDao) GetBatchOrderList(whereList []*model.WhereSqlCond) (datas []*BusinessBatchTransferOrderDao, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := `SELECT 
					batch_no, business_no, status, file_content, currency_type,
					real_amount, generate_all, total_number 
				FROM business_transfer_batch_order  `
	rows, stmt, errT := ss_sql.Query(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if errT != nil {
		return nil, errT
	}
	for rows.Next() {
		var batchNo, businessNo, status, fileContent, currencyType, realAmount, generateAll, totalNumber sql.NullString
		errT = rows.Scan(
			&batchNo,
			&businessNo,
			&status,
			&fileContent,
			&currencyType,
			&realAmount,
			&generateAll,
			&totalNumber,
		)
		if errT != nil {
			return nil, errT
		}
		data := &BusinessBatchTransferOrderDao{
			BatchNo:      batchNo.String,
			BusinessNo:   businessNo.String,
			Status:       status.String,
			FileContent:  fileContent.String,
			CurrencyType: currencyType.String,
			RealAmount:   realAmount.String,
			GenerateAll:  generateAll.String,
			TotalNumber:  totalNumber.String,
		}
		datas = append(datas, data)
	}
	return datas, errT

}

func (BusinessBatchTransferOrderDao) UpdateOrderStatusPaySuccess(tx *sql.Tx, batchNo string) error {
	sqlStr := "UPDATE business_transfer_batch_order SET status= $3, payment_type = $4, processing_number = total_number WHERE batch_no=$1 AND status=$2 "
	return ss_sql.ExecTx(tx, sqlStr, batchNo, constants.BusinessBatchTransferOrderStatusPending,
		constants.BusinessBatchTransferOrderStatusPaySuccess, constants.ORDER_PAYMENT_BALANCE)
}

func (BusinessBatchTransferOrderDao) UpdateOrderStatusSuccess(tx *sql.Tx, batchNo, successCnt, successSum, failCnt, failSum string) error {
	sqlStr := "UPDATE business_transfer_batch_order" +
		" SET " +
		" status= $3," +
		" processing_number = processing_number - $4," +
		" successful_number = successful_number + $5," +
		" successful_amount = successful_amount + $6, " +
		" fail_number = fail_number + $7," +
		" fail_amount = fail_amount + $8," +
		" finish_time = CURRENT_TIMESTAMP" +
		" WHERE" +
		"  batch_no = $1" +
		"  AND status = $2 "
	return ss_sql.ExecTx(tx, sqlStr, batchNo, constants.BusinessBatchTransferOrderStatusPaySuccess,
		constants.BusinessBatchTransferOrderStatusSuccess,
		ss_count.Add(successCnt, failCnt),
		successCnt,
		successSum,
		failCnt,
		failSum,
	)
}

func (BusinessBatchTransferOrderDao) UpdateGenerateAllByBatchNo(batchNo, generateAll string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "UPDATE business_transfer_batch_order SET generate_all= $2 WHERE batch_no = $1  "

	return ss_sql.Exec(dbHandler, sqlStr, batchNo, generateAll)
}

func (BusinessBatchTransferOrderDao) UpdateProcessingNumber(batchNo, processingNumber string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if processingNumber == "" {
		return errors.New("参数错误")
	}
	sqlStr := "UPDATE business_transfer_batch_order " +
		" SET processing_number = processing_number + $2," +
		" update_time = CURRENT_TIMESTAMP " +
		" WHERE batch_no = $1  "

	return ss_sql.Exec(dbHandler, sqlStr, batchNo, processingNumber)
}

func (BusinessBatchTransferOrderDao) UpdateFailAmountTx(tx *sql.Tx, batchNo, amount string) error {

	if amount == "" {
		return errors.New("参数错误")
	}
	sqlStr := "UPDATE business_transfer_batch_order " +
		" SET processing_number = processing_number - 1, " +
		" fail_number = fail_number + 1," +
		" fail_amount = fail_amount + $2," +
		" update_time = CURRENT_TIMESTAMP " +
		" WHERE batch_no = $1  "

	return ss_sql.ExecTx(tx, sqlStr, batchNo, amount)
}

func (BusinessBatchTransferOrderDao) UpdateFailAmount(batchNo, cnt, amount string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if amount == "" || batchNo == "" {
		return errors.New("参数错误")
	}
	sqlStr := "UPDATE business_transfer_batch_order " +
		" SET " +
		" fail_number = fail_number + $2," +
		" fail_amount = fail_amount + $3," +
		" update_time = CURRENT_TIMESTAMP " +
		" WHERE batch_no = $1  "

	return ss_sql.Exec(dbHandler, sqlStr, batchNo, cnt, amount)
}

func (BusinessBatchTransferOrderDao) UpdateSuccessAmount(batchNo, cnt, amount string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	if amount == "" {
		return errors.New("参数错误")
	}
	sqlStr := "UPDATE business_transfer_batch_order " +
		" SET " +
		" successful_number = successful_number + $2," +
		" successful_amount = successful_amount + $3, " +
		" update_time = CURRENT_TIMESTAMP " +
		" WHERE batch_no = $1  "

	return ss_sql.Exec(dbHandler, sqlStr, batchNo, cnt, amount)
}
func (BusinessBatchTransferOrderDao) InsertOrder(data BusinessBatchTransferOrderDao) (batchNo string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	id := strext.GetDailyId()
	sqlStr := "INSERT INTO business_transfer_batch_order(batch_no, total_number, total_amount, status, " +
		" remarks, business_no, file_content, generate_all, currency_type," +
		" real_amount, fee, create_time)" +
		" VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,CURRENT_TIMESTAMP) "
	errT := ss_sql.Exec(dbHandler, sqlStr,
		id, data.TotalNumber, data.TotalAmount, constants.BusinessBatchTransferOrderStatusPending,
		data.Remarks, data.BusinessNo, data.FileContent, data.GenerateAll, data.CurrencyType,
		data.RealAmount, data.Fee)
	if errT != nil {
		return "", errT
	}

	return id, nil
}
