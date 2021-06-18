package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type WalletDetailDao struct{}

var WalletDetailDaoInst WalletDetailDao

func (wa *WalletDetailDao) InsertDetail(tx *sql.Tx, accNo, beforeAmount, fees, op, upperNo, orderType, orderNo, walletNo, description, channelNo, merchantNo, op1 string) (errCode string) {
	logNo := strext.GetDailyId()
	var afterAmount string
	switch op1 {
	case "+":
		afterAmount = ss_count.Add(beforeAmount, fees)
	case "-":

		afterAmount = ss_count.Add(beforeAmount, fees)
	}
	err := ss_sql.ExecTx(tx, `insert into business_wallet_detail(log_no,acc_no,before_amount,change_amount,after_amount,op,upper_no,order_type,order_no,wallet_no,description,channel_no,merchant_no,create_time) 
								values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,current_timestamp)`,
		logNo, accNo, beforeAmount, fees, afterAmount, op, upperNo, orderType, orderNo, walletNo, description, channelNo, merchantNo)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}
