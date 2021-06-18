package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

type LogToHeadquartersDao struct {
}

var LogToHeadquartersDaoInst LogToHeadquartersDao

func (*LogToHeadquartersDao) InsertLogToHeadquarters(tx *sql.Tx, serviceNo, imageURL, carNo, collectType, amount, balanceType string, orderType int) string {
	logNoT := strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into log_to_headquarters(log_no,servicer_no,currency_type,amount,order_status,collection_type,card_no,order_type,image_id,create_time) 
				values ($1,$2,$3,$4,$5,$6,$7,$8,$9,current_timestamp)`,
		logNoT, serviceNo, balanceType, amount, "0", collectType, carNo, orderType, imageURL)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT
}
