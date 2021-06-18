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

type LogToCustDao struct{}

var LogToCustDaoInst LogToCustDao

func (LogToCustDao) Insert(tx *sql.Tx, currencyType, custNo, collectionType, cardNo, amount, orderType, orderStatus, lat, lng, fees, ip string) string {
	logNoT := strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into log_to_cust(log_no,currency_type,cust_no,collection_type,card_no,amount,create_time,order_type,order_status,lat,lng,fees,ip,payment_type) 
				values ($1,$2,$3,$4,$5,$6,current_timestamp,$7,$8,$9,$10,$11,$12,$13)`,
		logNoT, currencyType, custNo, collectionType, cardNo, amount, orderType, orderStatus, lat, lng, fees, ip, constants.ORDER_PAYMENT_BANK_WITHDRAW)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT
}

func (*LogToCustDao) LogToCustDetail(logNo string) (data *go_micro_srv_bill.ToCustDetailData, returnErr string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT ltc.log_no, ltc.amount, ltc.order_status, ltc.create_time, ltc.finish_time, ltc.currency_type, ltc.fees, ltc.payment_type " +
		", ca.card_number, ca.name, ch.channel_name " +
		" FROM log_to_cust ltc " +
		"left join card ca on ca.card_no = ltc.card_no " +
		"left join channel ch on ch.channel_no = ca.channel_no " +
		" WHERE  ltc.log_no = $1 "
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, logNo)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	data = &go_micro_srv_bill.ToCustDetailData{}
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
		data.OrderType = constants.VaReason_Cust_Withdraw

		//以下直接等于2的错误的
		//data.OpType = "2"
	}

	return data, ss_err.ERR_SUCCESS
}
