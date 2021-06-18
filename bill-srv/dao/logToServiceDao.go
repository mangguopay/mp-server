package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

type LogToServiceDao struct {
}

var (
	LogToServiceDaoInstance LogToServiceDao
)

/**
 * 增加服务商日志
 * @Param ServiceNo 商户号
 * @Param amount    金额
 * @Param collectType 收款类型
 * @Param orderType 订单类型
 * @Param currencyType 货币类型
 */
func (*LogToServiceDao) InsertLogToService(tx *sql.Tx, serviceNo, amount, collectType, cardNo string, orderType int, currencyType string) string {
	logNoT := strext.GetDailyId()

	err := ss_sql.ExecTx(tx, `insert into log_to_servicer(log_no,currency_type,servicer_no,collection_type,card_no,amount,order_type,order_status,create_time) 
				values ($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)`,
		logNoT, currencyType, serviceNo, collectType, cardNo, amount, orderType, "0")
	if err != nil {
		ss_log.Error("err=%v", err)
		return ""
	}
	return logNoT
}
