package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"a.a/cu/strext"

	"a.a/cu/db"
	"a.a/mp-server/common/ss_sql"
)

var TransferInstance Transfer

type Transfer struct {
	OutTransferNo string
	TransferNo    string
	CurrencyType  string
	Amount        int64
	AppId         string
	CreateTime    string
	CountryCode   string
	PayeePhone    string
	PayeeEmail    string
	Remark        string
	Status        int64
	StatusStr     string
}

const (
	TransferStatusPending = 1 // 转账中
	TransferStatusSuccess = 2 // 成功
	TransferStatusFail    = 3 // 失败

)

func GetTransferStatusString(transferStatus int64) string {
	switch transferStatus {
	case TransferStatusPending:
		return "转账中"
	case TransferStatusSuccess:
		return "成功"
	case TransferStatusFail:
		return "失败"
	}

	return fmt.Sprintf("%d", transferStatus)
}

// 插入一条记录
func (o *Transfer) Insert(transfer *Transfer) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	transfer.OutTransferNo = "transfer" + strext.GetDailyId()

	sqlStr := "INSERT INTO transfer (create_time, out_transfer_no, amount, currency_type, app_id, country_code, payee_phone, payee_email, remark, status)"
	sqlStr += " VALUES (current_timestamp, $1, $2, $3, $4, $5, $6, $7, $8, $9)"

	execErr := ss_sql.Exec(dbHandler, sqlStr,
		transfer.OutTransferNo, transfer.Amount,
		transfer.CurrencyType, transfer.AppId, transfer.CountryCode,
		transfer.PayeePhone, transfer.PayeeEmail, transfer.Remark, TransferStatusPending)

	return execErr
}

func (o *Transfer) GetTransferList(page, pageSize int) ([]Transfer, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	offset := (page - 1) * pageSize

	sqlStr := "SELECT out_transfer_no, transfer_no, amount, currency_type, app_id, country_code, payee_phone, payee_email, remark, create_time, status FROM transfer  "
	sqlStr += " ORDER BY create_time DESC LIMIT $1 OFFSET $2 "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr, pageSize, offset)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		return nil, qErr
	}

	list := []Transfer{}

	for rows.Next() {
		var outTransferNo, transferNo, amount, currencyType, appId, countryCode, payeePhone, payeeEmail, remark, createTime, status sql.NullString

		err := rows.Scan(&outTransferNo, &transferNo, &amount, &currencyType,
			&appId, &countryCode, &payeePhone, &payeeEmail, &remark, &createTime, &status)
		if err != nil {
			return nil, err
		}
		list = append(list, Transfer{
			OutTransferNo: outTransferNo.String,
			TransferNo:    transferNo.String,
			Amount:        strext.ToInt64(amount.String),
			CurrencyType:  currencyType.String,
			AppId:         appId.String,
			CountryCode:   countryCode.String,
			PayeePhone:    payeePhone.String,
			PayeeEmail:    payeeEmail.String,
			Remark:        remark.String,
			CreateTime:    createTime.String,

			Status:    strext.ToInt64(status.String),
			StatusStr: GetTransferStatusString(strext.ToInt64(status.String)),
		})
	}

	return list, nil
}

// 通过订单号获取一条记录
func (o *Transfer) GetOneByOutTransferNo(outTransferNo string) (*Transfer, error) {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := "SELECT out_transfer_no, transfer_no, amount, currency_type, app_id, create_time, country_code, payee_phone, payee_email, remark, status "
	sqlStr += " FROM transfer WHERE out_transfer_no=$1 "

	var out_transfer_no, transferNo, amount, currencyType, appId, createTime, countryCode, payeePhone, payeeEmail, remark, status sql.NullString

	qErr := ss_sql.QueryRow(dbHandler, sqlStr,
		[]*sql.NullString{&out_transfer_no, &transferNo, &amount, &currencyType, &appId, &createTime, &countryCode, &payeePhone, &payeeEmail, &remark, &status},
		outTransferNo,
	)

	if qErr != nil {
		return nil, qErr
	}

	order := &Transfer{
		OutTransferNo: out_transfer_no.String,
		TransferNo:    transferNo.String,
		CurrencyType:  currencyType.String,
		AppId:         appId.String,
		Amount:        strext.ToInt64(amount.String),
		Status:        strext.ToInt64(status.String),
		StatusStr:     GetTransferStatusString(strext.ToInt64(status.String)),
		CreateTime:    createTime.String,
		CountryCode:   countryCode.String,
		PayeePhone:    payeePhone.String,
		PayeeEmail:    payeeEmail.String,
		Remark:        remark.String,
	}

	return order, nil
}

func (o *Transfer) UpdateTransferNo(transferNo string, outTransferNo string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE transfer SET transfer_no=$1 WHERE out_transfer_no=$2`

	err := ss_sql.Exec(dbHandler, sqlStr, transferNo, outTransferNo)
	if nil != err {
		return err
	}

	return nil
}

// 更新订单转账成功
func (o *Transfer) UpdateTransferSuccess(outTransferNo string, transferTime string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE transfer SET status=$1, transfer_time=$2 WHERE out_transfer_no=$3`

	err := ss_sql.Exec(dbHandler, sqlStr, TransferStatusSuccess, transferTime, outTransferNo)
	if nil != err {
		return err
	}

	return nil
}

// 更新订单转账失败
func (o *Transfer) UpdateTransferFail(outTransferNo string) error {
	dbHandler := db.GetDB(DbMerchantMock)
	if dbHandler == nil {
		return errors.New("获取数据库连接失败")
	}
	defer db.PutDB(DbMerchantMock, dbHandler)

	sqlStr := `UPDATE transfer SET status=$1 WHERE out_transfer_no=$2`

	err := ss_sql.Exec(dbHandler, sqlStr, TransferStatusFail, outTransferNo)
	if nil != err {
		return err
	}

	return nil
}
