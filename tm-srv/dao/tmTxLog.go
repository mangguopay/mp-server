package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/tm-srv/common"
)

// 事务日志记录表
type TmTxLogDao struct {
	Id         string
	UnFinishNo int
	CreateTime string
	TmServerId string
}

// 插入一条记录
func (t *TmTxLogDao) Insert() error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sql := `INSERT INTO tm_tx_log (id, un_finish_no, create_time, tm_server_id) VALUES ($1, $2, current_timestamp, $3)`
	return ss_sql.Exec(dbHandler, sql, t.Id, t.UnFinishNo, t.TmServerId)
}

// 添加事务记录日志
func AddTmTxLog(unFinishNo int) {
	o := &TmTxLogDao{}
	o.Id = strext.GetDailyId()
	o.UnFinishNo = unFinishNo
	o.TmServerId = common.TmServerFullid

	if err := o.Insert(); err != nil {
		ss_log.Error("添加事务记录日志, err:%v, data:%+v", err, o)
	}
}
