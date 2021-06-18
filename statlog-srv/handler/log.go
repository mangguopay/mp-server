package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/statlog"
	"a.a/mp-server/common/ss_sql"
)

type StatlogHandler struct{}

var StatlogHandlerInst StatlogHandler

func (StatlogHandler) PushApiLog(req *go_micro_srv_statlog.PushApiLogRequest) {
	if req.AccountNo == "" {
		req.AccountNo = ss_sql.UUID
	}
	if req.AccountType == "" {
		req.AccountType = "-1"
	}

	bp, err := db.SsInfluxDBInst.NewBatchPoints(constants.IDB_LOG, db.Precision_Ns)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}

	p, err := db.SsInfluxDBInst.NewPoint("api_log", nil, map[string]interface{}{
		"url":          req.Url,
		"method":       req.Method,
		"log_time":     req.LogTime,
		"during":       req.During,
		"status_code":  req.StatusCode,
		"account_no":   req.AccountNo,
		"account_type": req.AccountType,
		"trace_no":     req.TraceNo,
		"ua":           req.Ua,
	})
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}
	bp.AddPoint(p)
	err = db.SsInfluxDBInst.Write("t_log", bp)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}

	//
	dbHandler := db.GetDB("statlog")
	defer db.PutDB("mp_log", dbHandler)
	err = ss_sql.Exec(dbHandler, `insert into log_api(log_no,url,method,log_time,during,status_code,account_type,ip,account_no,trace_no,ua) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		strext.GetDailyId(), req.Url, req.Method, req.LogTime, req.During, req.StatusCode, req.AccountType, req.Ip, req.AccountNo, req.TraceNo, req.Ua)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}
	return
}
