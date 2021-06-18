package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type TransferDao struct {
}

var TransferDaoInst TransferDao

func (TransferDao) UpdateTransferOrderIsCount(orderNo string, isCount int) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update business_transfer_order set is_count=$1 where order_no=$2`,
		isCount, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

func (*TransferDao) GetTransferInfoFromNo(tx *sql.Tx, orderNo string) (string, string, string, string, string, string, string, string) {
	var roleTypeT, settlementMethodT, feeT, accNoT, channelNoT, rateT, amountT, countFeeT sql.NullString
	err := ss_sql.QueryRowTx(tx, `SELECT tr.role_type,tr.settlement_method,tr.fee,tr.acc_no,tr.channel_no,c.rate,tr.amount,c.count_fee 
				from business_transfer_order tr LEFT JOIN business_channel c ON tr.channel_no = c.channel_no 
				where tr.order_no=$1 and tr.is_count='1' limit 1`,
		[]*sql.NullString{&roleTypeT, &settlementMethodT, &feeT, &accNoT, &channelNoT, &rateT, &amountT, &countFeeT}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", "", "", "", ""
	}
	return roleTypeT.String, settlementMethodT.String, feeT.String, accNoT.String, channelNoT.String, rateT.String, amountT.String, countFeeT.String
}

func (*TransferDao) ConfirmIsNoCount(orderNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_transfer_order set is_count = '1' where order_no=$1 and is_count='0'"
	result, err := ss_sql.ExecWithResult(dbHandler, sqlStr, orderNo)
	if err != nil {
		return ""
	} else {
		affected, _ := result.RowsAffected()
		if affected == 0 { // 没更新成功
			return ""
		}
	}
	return ss_err.ERR_SUCCESS
}
