package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

type LogApiPayDao struct {
	TraceNo        string
	AppId          string
	ReqMethod      string
	ReqUri         string
	ReqBody        string
	RespDate       string
	ReqTime        int64
	TrafficStatus  string
	BusinessStatus string
}

var LogApiPayDaoInst LogApiPayDao

func (LogApiPayDao) Insert(d *LogApiPayDao) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "INSERT INTO log_api_pay (trace_no, req_method, req_uri, req_body, resp_data, traffic_status, business_status," +
		" app_id, req_time, create_time) " +
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP)"

	return ss_sql.Exec(dbHandler, sqlStr, d.TraceNo, d.ReqMethod, d.ReqUri, d.ReqBody, d.RespDate, d.TrafficStatus, d.BusinessStatus,
		d.AppId, d.ReqTime)
}
