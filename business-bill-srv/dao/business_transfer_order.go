package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"fmt"
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
	FinishTime     string
	WrongReason    string
	AuthName       string
	ToAccount      string
	TransferType   string
	OutTransferNo  string
	AppId          string
}

var BusinessTransferOrderDaoInst BusinessTransferOrderDao

func (BusinessTransferOrderDao) InsertTx(tx *sql.Tx, d *BusinessTransferOrderDao) (logNo string, err error) {
	logNoT := strext.GetDailyId()

	if d.ToAccountNo != "" && d.ToBusinessNo == "" { //转账给个人的转账订单插入
		sqlStr := "INSERT INTO business_transfer_order(log_no, from_account_no, from_business_no, to_account_no, " +
			" amount, currency_type, rate, fee, payment_type," +
			" real_amount, order_status, remarks, transfer_type, out_transfer_no, create_time) " +
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP)"

		err := ss_sql.ExecTx(tx, sqlStr,
			logNoT, d.FromAccountNo, d.FromBusinessNo, d.ToAccountNo,
			d.Amount, d.CurrencyType, d.Rate, d.Fee, d.PaymentType,
			d.RealAmount, d.OrderStatus, d.Remarks, d.TransferType, d.OutTransferNo)
		if err != nil {
			return "", err
		}
	}

	return logNoT, nil
}

func (BusinessTransferOrderDao) GetOrderNoList(endTime, orderStatus string, limit int) ([]string, error) {
	dbhHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbhHandler)

	sqlStr := "SELECT bt.log_no " +
		"FROM business_transfer_order AS bt " +
		"LEFT JOIN enterprise_transfer_order AS et ON et.transfer_log_no = bt.log_no " +
		"WHERE bt.create_time <= $1 AND order_status = $2 " +
		"ORDER BY bt.create_time DESC Limit $3 "

	rows, stmt, err := ss_sql.Query(dbhHandler, sqlStr, endTime, orderStatus, limit)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var arr []string
	for rows.Next() {
		var transferNo sql.NullString
		if err := rows.Scan(&transferNo); err != nil {
			return nil, err
		}
		arr = append(arr, transferNo.String)
	}
	return arr, nil
}

func (BusinessTransferOrderDao) GetOrderByTransferNo(appId, transferNo, outTransferNo string) (*BusinessTransferOrderDao, error) {
	dbhHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbhHandler)

	var (
		i    = 2
		args []interface{}
	)
	whereStr := "WHERE et.app_id = $1 "
	args = append(args, appId)
	if transferNo != "" {
		whereStr += fmt.Sprintf("AND bt.log_no = $%v", i)
		args = append(args, transferNo)
		i++
	}

	if outTransferNo != "" {
		whereStr += fmt.Sprintf("AND bt.out_transfer_no = $%v", i)
		args = append(args, outTransferNo)
	}

	sqlStr := "SELECT bt.log_no, bt.out_transfer_no, bt.amount, bt.currency_type, bt.order_status, bt.finish_time, bt.wrong_reason " +
		"FROM business_transfer_order bt " +
		"LEFT JOIN enterprise_transfer_order et ON et.transfer_log_no = bt.log_no "

	sqlStr += whereStr
	var logNo, outTransferNoT, amount, currencyType, orderStatus, finishTime, wrongReason sql.NullString
	err := ss_sql.QueryRow(dbhHandler, sqlStr, []*sql.NullString{&logNo, &outTransferNoT, &amount, &currencyType,
		&orderStatus, &finishTime, &wrongReason}, args...)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessTransferOrderDao)
	obj.LogNo = logNo.String
	obj.OutTransferNo = outTransferNoT.String
	obj.Amount = amount.String
	obj.CurrencyType = currencyType.String
	obj.OrderStatus = orderStatus.String
	obj.FinishTime = finishTime.String
	obj.WrongReason = wrongReason.String
	obj.AppId = appId

	return obj, nil
}
