package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/tm-srv/common"
	"time"
)

// 事务日志时间记录表
type TmTimeLogDao struct {
	TxNo       string
	StartTime  string
	EndTime    string
	Duration   int64
	EndMode    string
	SqlList    string
	SqlNum     int
	TmServerId string
}

// 插入一条记录
func (t *TmTimeLogDao) Insert() error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sql := `INSERT INTO tm_time_log (tx_no, start_time, end_time, duration, end_mode, sql_list, sql_num, tm_server_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	return ss_sql.Exec(dbHandler, sql, t.TxNo, t.StartTime, t.EndTime, t.Duration, t.EndMode, t.SqlList, t.SqlNum, t.TmServerId)
}

// 添加事务时间记录日志
func AddTmTimeLog(txNo string, startTime time.Time, endTime time.Time, endType string, sqlList []string) {
	o := &TmTimeLogDao{}
	o.TxNo = txNo
	o.StartTime = ss_time.ForPostgres(startTime)
	o.EndTime = ss_time.ForPostgres(endTime)
	o.Duration = int64(endTime.Sub(startTime) / 1000 / 1000) // 单位: 毫秒
	o.EndMode = endType
	o.SqlList = strext.ToJson(sqlList)
	o.SqlNum = len(sqlList)
	o.TmServerId = common.TmServerFullid

	if err := o.Insert(); err != nil {
		ss_log.Error("%s:添加事务日志失败, err:%v, data:%+v", txNo, err, o)
	}
}
