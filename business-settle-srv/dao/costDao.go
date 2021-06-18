package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type CostDao struct {
}

var CostDaoInst CostDao

func (wa *CostDao) InsertCost(tx *sql.Tx, orderNo, orderType, amount, channelNo string) (errCode string) {
	err := ss_sql.ExecTx(tx, `insert into business_cost(order_no,order_type,amount,channel_no,create_time) 
								values ($1,$2,$3,$4,current_timestamp)`,
		orderNo, orderType, amount, channelNo)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}
