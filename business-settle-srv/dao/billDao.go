package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type BillDao struct {
}

var BillDaoInst BillDao

func (*BillDao) ConfirmIsNoCountTx(tx *sql.Tx, orderNo string) string {
	var feesT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select fees from business_bill where order_no=$1 and is_count='1'  limit 1`, []*sql.NullString{&feesT}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return feesT.String
}

func (*BillDao) ConfirmIsNoCount(orderNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//var feesT sql.NullString
	sqlStr := "update business_bill set is_count = '1' where order_no=$1 and is_count='0'"
	result, err := ss_sql.ExecWithResult(dbHandler, sqlStr, orderNo)
	if err != nil {
		return ""
	} else {
		affected, _ := result.RowsAffected()
		if affected == 0 { // 没更新成功
			return ""
		}
	}

	//err := ss_sql.QueryRowTx(tx, `select fees from bill where order_no=$1 and is_count='1'  limit 1`, []*sql.NullString{&feesT}, orderNo)
	//if nil != err {
	//	ss_log.Error("err=%v", err)
	//	return ""
	//}
	//return feesT.String
	return ss_err.ERR_SUCCESS
}

func (BillDao) UpdateExchangeOrderStatus(tx *sql.Tx, logNo, orderStatus, errReason string) string {
	err := ss_sql.ExecTx(tx, `update business_bill set order_status=$2, err_reason=$3 where log_no=$1`,
		logNo, orderStatus, errReason)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (BillDao) UpdateIsCountFromLogNo(logNo string, countStatus, isWallted int) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update business_bill set is_count=$2,is_wallted=$3  where order_no=$1`, logNo, countStatus, isWallted)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}
func (BillDao) UpdateIsCountFromLogNoTx(tx *sql.Tx, logNo string, countStatus, isWallted int) string {
	err := ss_sql.ExecTx(tx, `update business_bill set is_count=$2,is_wallted=$3,is_settled = $4  where order_no=$1`, logNo, countStatus, isWallted, 1)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_SYS_DB_OP
	}
	return ss_err.ERR_SUCCESS
}

func (*BillDao) QueryRateInfo(tx *sql.Tx, orderNo string) (string, string, string, string, string, string, string) {
	var rateNoT, agencyNoT, accNoT, channelRateT, channelNoT, rateT, countFeeT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select b.rate_no,m.agency_no,b.acc_no,c.rate,b.channel_no,b.rate,c.count_fee  from business_bill b  
							LEFT JOIN business_merchant m ON m.acc_no  = b.acc_no 
							LEFT JOIN business_channel c ON c.channel_no = b.channel_no  and c.is_delete = '0'
							WHERE b.order_no = $1 and b.is_count='1' limit 1`,
		[]*sql.NullString{&rateNoT, &agencyNoT, &accNoT, &channelRateT, &channelNoT, &rateT, &countFeeT}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", "", "", ""
	}
	return rateNoT.String, agencyNoT.String, accNoT.String, channelRateT.String, channelNoT.String, rateT.String, countFeeT.String
}
func (*BillDao) QueryBillInfo(tx *sql.Tx, orderNo string) (string, string) {
	var amountT, feesT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select amount,fee from business_bill where order_no = $1 limit 1`, []*sql.NullString{&amountT, &feesT}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", ""
	}
	return amountT.String, feesT.String
}
