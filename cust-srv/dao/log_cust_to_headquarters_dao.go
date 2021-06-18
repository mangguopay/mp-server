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
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type LogCustToHeadquartersDao struct{}

var LogCustToHeadquartersDaoInst LogCustToHeadquartersDao

func (LogCustToHeadquartersDao) Insert(tx *sql.Tx, custNo, currencyType, amount, orderStatus, collectionType, cardNo, orderType, imageId, arriveAmount, fees, lat, lng, ip string) string {
	logNo := strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into log_cust_to_headquarters(log_no,cust_no,currency_type,amount,order_status,collection_type,card_no,create_time,order_type,image_id,arrive_amount,fees,lat,lng,ip,payment_type) 
				values ($1,$2,$3,$4,$5,$6,$7,current_timestamp,$8,$9,$10,$11,$12,$13,$14,$15)`,
		logNo, custNo, currencyType, amount, orderStatus, collectionType, cardNo, orderType, imageId, arriveAmount, fees, lat, lng, ip, constants.ORDER_PAYMENT_BANK_TRANSFER)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNo
}

func (*LogCustToHeadquartersDao) QueryOrderStatusFromLogNo(orderNo string) (string, string, string, string, string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var statusT, currencyType, custNo, arriveAmount, fees, amount sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select order_status,currency_type,cust_no,arrive_amount,fees,amount from log_cust_to_headquarters where log_no=$1  limit 1`, []*sql.NullString{&statusT, &currencyType, &custNo, &arriveAmount, &fees, &amount}, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", "", "", "", "", ""
	}
	return statusT.String, currencyType.String, custNo.String, arriveAmount.String, fees.String, amount.String
}

func (*LogCustToHeadquartersDao) UpdateStatusFromLogNo(orderNo string, status int32) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `update log_cust_to_headquarters set order_status= $1 where log_no=$2`, status, orderNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_OPERATE_FAILD
	}
	return ss_err.ERR_SUCCESS
}

func (*LogCustToHeadquartersDao) CustToHeadquartersDetail(logNo string) (data *go_micro_srv_bill.CustToHeadquartersDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT log_no, amount, order_status, create_time, finish_time, currency_type,  fees,arrive_amount,payment_type  FROM  " +
		"log_cust_to_headquarters  WHERE  log_no = $1"
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
			&data.ArriveAmount,
			&data.PaymentType,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}
		data.OrderType = constants.VaReason_Cust_Save
	}

	return data, ss_err.ERR_SUCCESS
}

func (*LogCustToHeadquartersDao) CustToHeadquartersList(whereStr string, whereArgs []interface{}) (datas []*go_micro_srv_cust.CustToHeadquartersData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT lcth.log_no, lcth.currency_type, lcth.collection_type" +
		", lcth.amount, lcth.create_time, lcth.order_type, lcth.order_status, lcth.finish_time" +
		",lcth.lat, lcth.lng, lcth.fees, lcth.ip, lcth.image_id " +
		", acc.account, ca.name, ca.card_number, ch.channel_name " +
		" FROM log_cust_to_headquarters lcth " +
		" LEFT JOIN cust cu ON cu.cust_no = lcth.cust_no " +
		" LEFT JOIN account acc ON acc.uid = cu.account_no " +
		" LEFT JOIN card_head ca ON ca.card_no = lcth.card_no " +
		" LEFT JOIN channel_cust_config ccc ON ccc.id = ca.channel_cust_config_id " +
		" LEFT JOIN channel ch ON ch.channel_no = ccc.channel_no "
	sqlStr += whereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	var datasT []*go_micro_srv_cust.CustToHeadquartersData
	for rows.Next() {
		var data go_micro_srv_cust.CustToHeadquartersData
		var finishTime, imageId, channelName sql.NullString
		if err = rows.Scan(
			&data.LogNo,
			&data.CurrencyType,
			&data.CollectionType,
			&data.Amount,
			&data.CreateTime,

			&data.OrderType,
			&data.OrderStatus,
			&finishTime,
			&data.Lat,
			&data.Lng,

			&data.Fees,
			&data.Ip,
			&imageId,
			&data.Account,
			&data.Name,

			&data.CardNumber,
			&channelName,
		); err != nil {
			ss_log.Error("err=[%v]", err)
			continue
		}
		data.FinishTime = finishTime.String
		data.ImageId = imageId.String
		data.ChannelName = channelName.String

		datasT = append(datasT, &data)
	}

	return datasT, nil

}

// GetCustToHeadquartersCountByDate 统计
func (*LogCustToHeadquartersDao) GetCustToHeadquartersCountByDate(beforeDay, currencyType string) (*DataCount, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	startTime := beforeDay
	endTime, tErr := ss_time.TimeAfter(beforeDay, ss_time.DateFormat, time.Hour*24) // 下一天的日期
	if tErr != nil {
		return nil, tErr
	}

	sqlStr := `select count(1) as totalCount,sum(amount) as totalAmount,sum(fees) as totalFees from log_cust_to_headquarters WHERE create_time >= $1 
		and create_time < $2 AND currency_type = $3 and order_status = 1`

	var num, amount, fee sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&num, &amount, &fee}, startTime, endTime, currencyType)

	if errT != nil {
		return nil, errT
	}
	data := &DataCount{
		Num:    strext.ToInt64(num.String),
		Amount: strext.ToInt64(amount.String),
		Fee:    strext.ToInt64(fee.String),
		Type:   1, // 向总部充值
		Day:    startTime,
		CType:  currencyType,
	}
	return data, nil
}
