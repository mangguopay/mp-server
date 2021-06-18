package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type WalletDao struct{}

var WalletDaoInst WalletDao

func (wa *WalletDao) ConfirmExistWallet(tx *sql.Tx, accNo, roleType string) (logNo, beforeAmount string) {
	logNo, beforeAmount = wa.GetLogNoFromAccNo(tx, accNo, roleType)
	if logNo == "" {
		logNo = wa.InitWallet(tx, accNo, roleType)
		return logNo, "0"
	}
	return logNo, beforeAmount
}

func (wa *WalletDao) GetLogNoFromAccNo(tx *sql.Tx, accNo, roleType string) (string, string) {
	var logNoT, amountT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select log_no,amount from business_wallet where acc_no=$1 and role_type = $2  limit 1`,
		[]*sql.NullString{&logNoT, &amountT}, accNo, roleType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return logNoT.String, amountT.String
}
func (wa *WalletDao) InitWallet(tx *sql.Tx, accNo, roleType string) (logNo string) {
	logNo = strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into business_wallet(log_no,amount,acc_no,create_time,frozen_amount,role_type) 
								values ($1,$2,$3 ,current_timestamp,$4,$5)`,
		logNo, 0, accNo, 0, roleType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNo
}

func (wa *WalletDao) ModifyAmountFromNo(tx *sql.Tx, amount, logNo, op string) (errCode string) {
	switch op {
	case "+":
		err := ss_sql.ExecTx(tx, `update business_wallet set amount=amount+$1,modify_time=current_timestamp where log_no=$2`, amount, logNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}
	case "-":
		err := ss_sql.ExecTx(tx, `update business_wallet set amount=amount+$1,modify_time=current_timestamp where log_no=$2`, amount, logNo)
		if nil != err {
			ss_log.Error("err=%v", err)
			return ss_err.ERR_PAY_AMT_NOT_ENOUGH
		}

	}

	return ss_err.ERR_SUCCESS
}
