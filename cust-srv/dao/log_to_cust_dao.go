package dao

import (
	"database/sql"
	"time"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type LogToCustDao struct {
}

var LogToCustDaoInst LogToCustDao

//func (LogToCustDao) Insert(tx *sql.Tx, currencyType, custNo, collectionType, cardNo, amount, orderType, orderStatus, lat, lng, fees, ip string) string {
//	logNoT := strext.GetDailyId()
//	err := ss_sql.ExecTx(tx, `insert into log_to_cust(log_no,currency_type,cust_no,collection_type,card_no,amount,create_time,order_type,order_status,lat,lng,fees,ip)
//				values ($1,$2,$3,$4,$5,$6,current_timestamp,$7,$8,$9,$10,$11,$12)`,
//		logNoT, currencyType, custNo, collectionType, cardNo, amount, orderType, orderStatus, lat, lng, fees, ip)
//	if err != nil {
//		ss_log.Error("err=%v", err)
//		return ""
//	}
//	return logNoT
//}

func (*LogToCustDao) QueryOrderStatusFromLogNo(orderNo string) (string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var statusT, currencyType, custNo, amount, fees sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select order_status,currency_type,cust_no,amount,fees from log_to_cust where log_no=$1  limit 1`, []*sql.NullString{&statusT, &currencyType, &custNo, &amount, &fees}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", ""
	}
	return statusT.String, currencyType.String, custNo.String, amount.String, fees.String
}
func (*LogToCustDao) UpdateStatusFromLogNo(tx *sql.Tx, orderNo, notes, imageId string, status int32) (errCode string) {

	err := ss_sql.ExecTx(tx, `update log_to_cust set order_status = $2, notes = $3, image_id = $4, finish_time = current_timestamp where log_no=$1`, orderNo, status, notes, imageId)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_OPERATE_FAILD
	}
	return ss_err.ERR_SUCCESS
}

func (*LogToCustDao) LogToCustDetail(logNo string) (data *go_micro_srv_bill.CustToHeadquartersDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT log_no, amount, order_status, create_time, finish_time, currency_type,fees,payment_type FROM  " +
		"log_to_cust  WHERE  log_no = $1"
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, logNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.CustToHeadquartersDetailData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, ss_err.ERR_SYS_DB_GET
	}
	for rows.Next() {
		var finishTime sql.NullString
		err = rows.Scan(
			&data.LogNo,
			&data.Amount,
			&data.OrderStatus,
			&data.CreateTime,
			&finishTime,
			&data.BalanceType,
			&data.Fees,
			&data.PaymentType,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}
		data.OrderType = constants.VaReason_Cust_Withdraw
	}

	return data, ss_err.ERR_SUCCESS
}

func (*LogToCustDao) GetToCustInfoByLogNo(logNo string) (accountUid, currencyTypeStr, amountStr string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var accUid, currencyType, amount sql.NullString
	sqlStr := " select acc.uid " +
		" from log_to_cust ltc " +
		" left join cust cu on cu.cust_no = ltc.cust_no " +
		" left join account acc on acc.uid = cu.account_no " +
		" where ltc.log_no = $1 "
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&accUid, &currencyType, &amount}, logNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", "", "", errT
	}
	return accUid.String, currencyType.String, amount.String, nil

}

// GetToCustInfoCount 统计
func (*LogToCustDao) GetToCustInfoCountByDate(beforeDay, currencyType string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := `select count(1) as totalCount,sum(amount) as totalAmount,sum(fees) as totalFees  from
	log_to_cust WHERE create_time >= $1 AND create_time < $2 and currency_type = $3 and order_status = 1 `

	//====================================
	var num, amount, fee sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount, &fee}, startTime, endTime, currencyType)

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}
	data := &DataCount{
		Num:    strext.ToInt64(num.String),
		Amount: strext.ToInt64(amount.String),
		Fee:    strext.ToInt64(fee.String),
		Type:   1, // 银行卡提现
		Day:    startTime,
		CType:  currencyType,
	}
	return data, nil
}

// DataCount 提现统计
type DataCount struct {
	Type   int64
	CType  string
	Num    int64
	Amount int64
	Fee    int64
	Day    string

	Usd2khrNum    int64
	Usd2khrAmount int64
	Usd2khrFee    int64

	Khr2usdNum    int64
	Khr2usdAmount int64
	Khr2usdFee    int64

	RegNum    int64 // 新增注册数量
	ServerNum int64 // 服务商数量
}
