package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type LogCustToHeadquartersDao struct{}

var LogCustToHeadquartersDaoInst LogCustToHeadquartersDao

func (LogCustToHeadquartersDao) Insert(tx *sql.Tx, custNo, currencyType, amount, orderStatus, collectionType, cardNo, orderType, imageId, arriveAmount, fees, lat, lng, ip string) string {
	logNo := strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into log_cust_to_headquarters(log_no,cust_no,currency_type,amount,order_status,collection_type,card_no,create_time,order_type,image_id,arrive_amount,fees,lat,lng,ip) 
				values ($1,$2,$3,$4,$5,$6,$7,current_timestamp,$8,$9,$10,$11,$12,$13,$14)`,
		logNo, custNo, currencyType, amount, orderStatus, collectionType, cardNo, orderType, imageId, arriveAmount, fees, lat, lng, ip)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNo
}

func (*LogCustToHeadquartersDao) CustToHeadquartersDetail(logNo string) (data *go_micro_srv_bill.CustToHeadquartersDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT lcth.log_no, lcth.amount, lcth.order_status, lcth.create_time" +
		", lcth.finish_time, lcth.currency_type, lcth.fees, lcth.arrive_amount, lcth.payment_type" +
		", ca.card_number,ca.name, ch.channel_name " +
		" FROM log_cust_to_headquarters lcth " +
		" left join card_head ca on ca.card_no = lcth.card_no " +
		" left join channel_cust_config ccc on ccc.id = ca.channel_cust_config_id " +
		" left join channel ch on ch.channel_no = ccc.channel_no " +
		" WHERE  lcth.log_no = $1"
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
			&data.CardNumber,
			&data.Name,
			&data.ChannelName,
		)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		if finishTime.String != "" {
			data.FinishTime = finishTime.String
		}
		data.OrderType = constants.VaReason_Cust_Save
		//以下直接等于+是错误的，
		//data.OpType = "1"
	}

	return data, ss_err.ERR_SUCCESS
}
