package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
)

type BusinessTransferOrderDao struct {
	LogNo          string
	FromAccountNo  string
	FromBusinessNo string
	ToAccountNo    string
	ToBusinessNo   string
	Amount         string
	CurrencyType   string
	Rate           string
	Fee            string
	PaymentType    string
	RealAmount     string
	OrderStatus    string
	BatchNo        string
	BatchRowNum    string
	Remarks        string
	CreateTime     string
	WrongReason    string
	AuthName       string
	ToAccount      string
	TransferType   string

	TransferNo    string
	OutTransferNo string
}

var BusinessTransferOrderDaoInst BusinessTransferOrderDao

//此接口是商家单个转账到商家的转账订单插入
func (BusinessTransferOrderDao) Insert(d *BusinessTransferOrderDao) (logNoT string, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	logNo := strext.GetDailyId()
	ss_log.Info("data[%+v]", d)

	sqlStr := "INSERT INTO business_transfer_order(log_no, from_account_no, from_business_no, to_account_no, to_business_no," +
		" amount, currency_type, rate, fee, payment_type," +
		" real_amount, order_status, remarks, transfer_type, create_time) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP)"

	err := ss_sql.Exec(dbHandler, sqlStr,
		logNo, d.FromAccountNo, d.FromBusinessNo, d.ToAccountNo, d.ToBusinessNo,
		d.Amount, d.CurrencyType, d.Rate, d.Fee, d.PaymentType,
		d.RealAmount, d.OrderStatus, d.Remarks, d.TransferType)
	if err != nil {
		return "", err
	}

	return logNo, nil
}

//此接口是商家批量转账到商家的转账订单插入。
func (BusinessTransferOrderDao) Insert2(d *BusinessTransferOrderDao) (logNoT string, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	logNo := strext.GetDailyId()
	ss_log.Info("data[%+v]", d)

	if d.ToAccountNo == "" && d.ToBusinessNo == "" { //批量转账生成的，错误订单，商家账号可能不存在
		sqlStr := "INSERT INTO business_transfer_order(log_no, from_account_no, from_business_no," +
			" amount, currency_type, rate, fee, payment_type," +
			" real_amount, order_status, batch_no, batch_row_num, remarks," +
			" wrong_reason, auth_name, to_account, transfer_type, create_time) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, CURRENT_TIMESTAMP)"

		err := ss_sql.Exec(dbHandler, sqlStr,
			logNo, d.FromAccountNo, d.FromBusinessNo,
			d.Amount, d.CurrencyType, d.Rate, d.Fee, d.PaymentType,
			d.RealAmount, d.OrderStatus, d.BatchNo, d.BatchRowNum, d.Remarks,
			d.WrongReason, d.AuthName, d.ToAccount, d.TransferType)
		if err != nil {
			return "", err
		}
	} else if d.ToAccountNo != "" && d.ToBusinessNo == "" { //批量转账生成的，转账给用户订单，商家身份是没有的
		sqlStr := "INSERT INTO business_transfer_order(log_no, from_account_no, from_business_no, to_account_no, " +
			" amount, currency_type, rate, fee, payment_type," +
			" real_amount, order_status, batch_no, batch_row_num, remarks," +
			" wrong_reason, auth_name, to_account, transfer_type, create_time) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, CURRENT_TIMESTAMP)"

		err := ss_sql.Exec(dbHandler, sqlStr,
			logNo, d.FromAccountNo, d.FromBusinessNo, d.ToAccountNo,
			d.Amount, d.CurrencyType, d.Rate, d.Fee, d.PaymentType,
			d.RealAmount, d.OrderStatus, d.BatchNo, d.BatchRowNum, d.Remarks,
			d.WrongReason, d.AuthName, d.ToAccount, d.TransferType)
		if err != nil {
			return "", err
		}
	} else { //转账给商家的
		sqlStr := "INSERT INTO business_transfer_order(log_no, from_account_no, from_business_no, to_account_no, to_business_no," +
			" amount, currency_type, rate, fee, payment_type," +
			" real_amount, order_status, batch_no, batch_row_num, remarks," +
			" wrong_reason, auth_name, to_account, transfer_type, create_time) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, CURRENT_TIMESTAMP)"

		err := ss_sql.Exec(dbHandler, sqlStr,
			logNo, d.FromAccountNo, d.FromBusinessNo, d.ToAccountNo, d.ToBusinessNo,
			d.Amount, d.CurrencyType, d.Rate, d.Fee, d.PaymentType,
			d.RealAmount, d.OrderStatus, d.BatchNo, d.BatchRowNum, d.Remarks,
			d.WrongReason, d.AuthName, d.ToAccount, d.TransferType)
		if err != nil {
			return "", err
		}
	}

	return logNo, nil
}

//此接口是商家转账到个人的转账订单插入
func (BusinessTransferOrderDao) Add(d *BusinessTransferOrderDao) (logNoT string, errT error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	logNo := strext.GetDailyId()

	ss_log.Info("data[%+v]", d)
	if d.ToAccountNo != "" && d.ToBusinessNo == "" { //转账给个人的转账订单插入
		sqlStr := "INSERT INTO business_transfer_order(log_no, from_account_no, from_business_no, to_account_no, " +
			" amount, currency_type, rate, fee, payment_type," +
			" real_amount, order_status, remarks, transfer_type, out_transfer_no, create_time) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP)"

		err := ss_sql.Exec(dbHandler, sqlStr,
			logNo, d.FromAccountNo, d.FromBusinessNo, d.ToAccountNo,
			d.Amount, d.CurrencyType, d.Rate, d.Fee, d.PaymentType,
			d.RealAmount, d.OrderStatus, d.Remarks, d.TransferType, d.OutTransferNo)
		if err != nil {
			return "", err
		}
	}

	return logNo, nil
}
func (BusinessTransferOrderDao) AddTx(tx *sql.Tx, d *BusinessTransferOrderDao) (logNoT string, errT error) {
	if d.LogNo == "" {
		d.LogNo = strext.GetDailyId()
	}
	if d.ToAccountNo != "" && d.ToBusinessNo == "" { //转账给个人的转账订单插入
		sqlStr := "INSERT INTO business_transfer_order(log_no, from_account_no, from_business_no, to_account_no, " +
			" amount, currency_type, rate, fee, payment_type," +
			" real_amount, order_status, remarks, transfer_type, out_transfer_no, create_time) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP)"

		err := ss_sql.ExecTx(tx, sqlStr,
			d.LogNo, d.FromAccountNo, d.FromBusinessNo, d.ToAccountNo,
			d.Amount, d.CurrencyType, d.Rate, d.Fee, d.PaymentType,
			d.RealAmount, d.OrderStatus, d.Remarks, d.TransferType, d.OutTransferNo)
		if err != nil {
			return "", err
		}
	}

	return d.LogNo, nil
}

func (BusinessTransferOrderDao) UpdateOrderStatusByLogNoTx(tx *sql.Tx, logNo, orderStatus, wrongReason string) error {
	sqlStr := "UPDATE business_transfer_order SET order_status=$3, wrong_reason = $4, finish_time=CURRENT_TIMESTAMP " +
		"WHERE log_no = $1 AND order_status = $2 "
	return ss_sql.ExecTx(tx, sqlStr, logNo, constants.BusinessTransferOrderStatusPending, orderStatus, wrongReason)
}
func (BusinessTransferOrderDao) UpdateOrderStatusByLogNo(logNo, orderStatus, wrongReason string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "UPDATE business_transfer_order SET order_status=$3, wrong_reason = $4, finish_time=CURRENT_TIMESTAMP " +
		"WHERE log_no = $1 AND order_status = $2 "
	return ss_sql.Exec(dbHandler, sqlStr, logNo, constants.BusinessTransferOrderStatusPending, orderStatus, wrongReason)
}

//确认
func (BusinessTransferOrderDao) CheckBusinessTransferOrder(batchNo, batchRowNum string) (cnt string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from business_transfer_order where batch_no = $1 AND batch_row_num = $2 "
	var total sql.NullString
	if errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, batchNo, batchRowNum); errT != nil {
		return "", errT
	}
	return total.String, nil
}

//获取转账订单内的一个转账批次的订单数量
func (BusinessTransferOrderDao) GetBatchCnt(batchNo string) (cnt string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from business_transfer_order where batch_no = $1 "
	var total sql.NullString
	if errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, batchNo); errT != nil {
		ss_log.Error("err[%v]", errT)
		return ""
	}
	return total.String
}

//获取一个转账批次的订单内的转账订单信息
func (BusinessTransferOrderDao) GetTransferOrderList(whereList []*model.WhereSqlCond) (datas []*BusinessTransferOrderDao, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "select " +
		" log_no, from_account_no , to_account_no, currency_type, fee, real_amount  " +
		" from business_transfer_order " + whereModel.WhereStr
	row, stmt, errT := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if errT != nil {
		ss_log.Error("err[%v]", errT)
		return nil, errT
	}
	var logNo, fromAccountNo, toAccountNo, currencyType, fee, realAmount sql.NullString
	for row.Next() {
		errT = row.Scan(
			&logNo, &fromAccountNo, &toAccountNo, &currencyType, &fee, &realAmount,
		)
		if errT != nil {
			ss_log.Error("err[%v]", errT)
			return nil, errT
		}
		data := &BusinessTransferOrderDao{
			LogNo:         logNo.String,
			FromAccountNo: fromAccountNo.String,
			ToAccountNo:   toAccountNo.String,
			CurrencyType:  currencyType.String,
			Fee:           fee.String,
			RealAmount:    realAmount.String,
		}
		datas = append(datas, data)
	}

	return datas, nil
}

//获取一个批量转账批次的成功转账订单金额、手续费统计
func (BusinessTransferOrderDao) GetTransferOrderSum(batchNo, orderStatus string) (realAmountSum, feeSum, cnt string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select " +
		" SUM(real_amount), SUM(fee), count(1)" +
		" from business_transfer_order " +
		" where batch_no = $1 " +
		" and order_status = $2 "
	var realAmountSumT, feeSumT, cntT sql.NullString
	if errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&realAmountSumT, &feeSumT, &cntT}, batchNo, orderStatus); errT != nil {
		ss_log.Error("err[%v]", errT)
		return "0", "0", "0", errT
	}

	if realAmountSumT.String == "" {
		realAmountSumT.String = "0"
	}
	if feeSumT.String == "" {
		feeSumT.String = "0"
	}
	if cntT.String == "" {
		cntT.String = "0"
	}
	return realAmountSumT.String, feeSumT.String, cntT.String, nil
}

func (BusinessTransferOrderDao) BusinessTransferBillsDetail(whereStr string, whereArgs []interface{}) (*go_micro_srv_cust.BusinessTransferOrderData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bto.log_no, bto.amount, bto.fee, bto.real_amount, bto.currency_type, bto.create_time, " +
		" bto.order_status, bto.remarks, acc.account, acc2.account, bto.rate," +
		" bto.batch_no, bto.batch_row_num, bto.wrong_reason " +
		"FROM business_transfer_order bto " +
		"LEFT JOIN account acc ON acc.uid = bto.to_account_no " +
		"LEFT JOIN account acc2 ON acc2.uid = bto.from_account_no "
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr+whereStr, whereArgs...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}

	data := new(go_micro_srv_cust.BusinessTransferOrderData)
	var toAccount, batchNo, batchRowNum, wrongReason sql.NullString
	err = rows.Scan(&data.LogNo, &data.Amount, &data.Fee, &data.RealAmount, &data.CurrencyType,
		&data.CreateTime, &data.OrderStatus, &data.Remarks, &toAccount, &data.FromAccount, &data.Rate,
		&batchNo, &batchRowNum, &wrongReason,
	)
	if err != nil {
		return nil, err
	}
	data.ToAccount = toAccount.String
	data.BatchNo = batchNo.String
	data.BatchRowNum = batchRowNum.String
	data.WrongReason = wrongReason.String

	return data, nil
}

func (BusinessTransferOrderDao) GetTransferOrderByLogNo(logNo string) (*BusinessTransferOrderDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT amount, fee, currency_type, order_status, from_account_no, to_account, to_account_no, to_account_type " +
		"FROM business_transfer_order " +
		"WHERE log_no = $1 "

	var amount, fee, currencyType, orderStatus, fromAccountNo, toAccount, toAccountNo, toAccountType sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&amount, &fee, &currencyType, &orderStatus, &fromAccountNo,
		&toAccount, &toAccountNo, &toAccountType}, logNo)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessTransferOrderDao)
	obj.LogNo = logNo
	obj.Amount = amount.String
	obj.Fee = fee.String
	obj.CurrencyType = currencyType.String
	obj.OrderStatus = orderStatus.String
	obj.FromAccountNo = fromAccountNo.String
	obj.ToAccount = toAccount.String
	obj.ToAccountNo = toAccountNo.String
	return obj, nil
}
