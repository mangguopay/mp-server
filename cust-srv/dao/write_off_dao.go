package dao

import (
	"a.a/cu/strext"
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type WriteoffDao struct {
}

type WriteOff struct {
	Code             string
	OrderNo          string
	OrderSource      string
	UseStatus        string
	OrderAmount      string
	RealAmount       string
	CurrencyType     string
	CreateTime       string
	FinishTime       string
	PayerAccount     string
	PayeeAccount     string
	PayeeVAccNo      string
	PayeeAccountType string
	DurationTime     string
}

var WriteoffDaoInst WriteoffDao

func (WriteoffDao) GetCode(logNo, orderType string) (code, err string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := ""
	switch orderType {
	case constants.VaReason_INCOME:
		sqlStr = "select code from writeoff where income_order_no = $1  "
	case constants.VaReason_OUTGO:
		sqlStr = "select code from writeoff where outgo_order_no = $1  "
	case constants.VaReason_TRANSFER:
		sqlStr = "select code from writeoff where transfer_order_no = $1  "
	default:
		ss_log.Error("orderType 类型错误 [%v]", orderType)
		return "", ss_err.ERR_PARAM
	}

	var codeT sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&codeT}, logNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", ss_err.ERR_PARAM
	}

	return codeT.String, ss_err.ERR_SUCCESS
}

//获取过期的核销码
func (WriteoffDao) GetExpiredCodeArr() (codes []string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select code from writeoff where duration_time <= current_timestamp and use_status = $1 "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, constants.WriteOffCodeWaitUse)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=%v", err)
		return nil
	}

	var codess []string
	for rows.Next() {
		var code sql.NullString
		err = rows.Scan(
			&code,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}

		if code.String != "" {
			codess = append(codess, code.String)
		}

	}

	return codess
}

//修改核销码状态
func (WriteoffDao) UpdateExpiredCodeStatus(code, status string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	sqlStr := "update writeoff set use_status = $1 "
	if status == constants.WriteOffCodeIsUse {
		sqlStr += ",finish_time=CURRENT_TIMESTAMP "
	}
	sqlStr += " where code = $2 "
	if err := ss_sql.Exec(dbHandler, sqlStr, status, code); err != nil {
		ss_log.Error("err=%v", err)
		return err
	}

	return nil
}

func (WriteoffDao) GetCodeList(whereStr string, args []interface{}) ([]*WriteOff, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT
		wo.code, wo.use_status, wo.create_time, wo.finish_time, wo.duration_time, wo.income_order_no, 
		wo.transfer_order_no, io.amount io_amount, io.real_amount io_real_amount, io.balance_type io_balance_type,
		tro.amount tro_amount, tro.real_amount tro_real_amount, tro.balance_type tro_balance_type,
		acc1.account payer_account, acc2.account payee_account, acc2.is_actived
	FROM writeoff wo 
	LEFT JOIN income_order io ON io.log_no = wo.income_order_no
	LEFT JOIN transfer_order tro ON tro.log_no = wo.transfer_order_no
	LEFT JOIN account acc1 ON acc1.uid = wo.send_acc_no 
	LEFT JOIN account acc2 ON acc2.uid = wo.recv_acc_no `

	sqlStr += whereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var dataList []*WriteOff
	for rows.Next() {
		var code, useStatus, createTime, finishTime, durationTime, inOrderNo, trOrderNo, inAmount, inRealAmount,
			inCurrencyType, trAmount, trRealAmount, trCurrencyType, payerAccount, payeeAccount, isActive sql.NullString
		err := rows.Scan(&code, &useStatus, &createTime, &finishTime, &durationTime, &inOrderNo, &trOrderNo,
			&inAmount, &inRealAmount, &inCurrencyType, &trAmount, &trRealAmount, &trCurrencyType,
			&payerAccount, &payeeAccount, &isActive)
		if err != nil {
			return nil, err
		}

		data := new(WriteOff)
		data.Code = code.String
		data.UseStatus = useStatus.String
		data.CreateTime = createTime.String
		data.FinishTime = finishTime.String
		data.PayerAccount = payerAccount.String
		data.PayeeAccount = payeeAccount.String
		data.PayeeAccountType = isActive.String
		data.DurationTime = durationTime.String
		data.OrderAmount = inAmount.String
		if data.OrderAmount == "" {
			data.OrderAmount = trAmount.String
		}
		data.RealAmount = inRealAmount.String
		if data.RealAmount == "" {
			data.RealAmount = trRealAmount.String
		}
		data.CurrencyType = inCurrencyType.String
		if data.CurrencyType == "" {
			data.CurrencyType = trCurrencyType.String
		}
		data.OrderNo = inOrderNo.String
		data.OrderSource = "pos端存款"
		if data.OrderNo == "" {
			data.OrderNo = trOrderNo.String
			data.OrderSource = "转账"
		}
		dataList = append(dataList, data)
	}

	return dataList, nil
}

func (WriteoffDao) CntCode(whereStr string, args []interface{}) (int32, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT COUNT(1)
		FROM writeoff wo 
		LEFT JOIN income_order io ON io.log_no = wo.income_order_no
		LEFT JOIN transfer_order tro ON tro.log_no = wo.transfer_order_no
		LEFT JOIN account acc1 ON acc1.uid = wo.send_acc_no 
		LEFT JOIN account acc2 ON acc2.uid = wo.recv_acc_no `
	sqlStr += whereStr

	var total sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&total}, args...)
	if err != nil {
		return 0, err
	}

	return strext.ToInt32(total.String), nil
}

func (WriteoffDao) GetCodeDetailByCode(code string) (*WriteOff, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `SELECT
		wo.code, wo.use_status, wo.duration_time,
		io.real_amount io_real_amount, tro.real_amount tro_real_amount, 
		io.balance_type io_balance_type, tro.balance_type tro_balance_type,
		io.recv_vacc io_payee_v_acc, tro.to_vaccount_no tro_payee_v_acc, 
		acc2.account payee_account, acc2.is_actived
	FROM writeoff wo 
	LEFT JOIN income_order io ON io.log_no = wo.income_order_no
	LEFT JOIN transfer_order tro ON tro.log_no = wo.transfer_order_no
	LEFT JOIN account acc1 ON acc1.uid = wo.send_acc_no 
	LEFT JOIN account acc2 ON acc2.uid = wo.recv_acc_no
	WHERE wo.code = $1 `

	var writeOffCode, useStatus, durationTime, ioRealAmount, ioCurrencyType, troRealAmount, troCurrencyType, ioVAccNo,
		troVAccNo, payeeAccount, isActive sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&writeOffCode, &useStatus, &durationTime, &ioRealAmount,
		&troRealAmount, &ioCurrencyType, &troCurrencyType, &ioVAccNo, &troVAccNo, &payeeAccount, &isActive}, code)
	if err != nil {
		return nil, err
	}

	payeeVAccNo := ioVAccNo.String
	if payeeVAccNo == "" {
		payeeVAccNo = troVAccNo.String
	}
	realAmount := ioRealAmount.String
	if realAmount == "" {
		realAmount = troRealAmount.String
	}
	currencyType := ioCurrencyType.String
	if currencyType == "" {
		currencyType = troCurrencyType.String
	}
	data := &WriteOff{
		Code:             writeOffCode.String,
		UseStatus:        useStatus.String,
		DurationTime:     durationTime.String,
		PayeeAccount:     payeeAccount.String,
		PayeeVAccNo:      payeeVAccNo,
		PayeeAccountType: isActive.String,
		RealAmount:       realAmount,
		CurrencyType:     currencyType,
	}
	return data, nil
}
