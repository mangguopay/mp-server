package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	ApiPayLogDaoInstance ApiPayLogDao
)

//
type ApiPayLogDao struct {
	Page      string
	PageSize  string
	StartTime string
	EndTime   string
	ReqMethod string //请求方法

	ReqUri         string //请求uri
	ReqBody        string //请求body
	RespData       string //返回数据
	TrafficStatus  string //通信状态(0失败，1成功)
	BusinessStatus string //业务处理装(0失败，1成功)

	AppId        string //应用id
	ReqStartTime string //请求时间开始（要转成时间戳）
	ReqEndTime   string //请求时间结束（要转成时间戳）
}

//返回数据封装
type ApiPayLogData struct {
	TraceNo   string //
	ReqMethod string //请求方法
	ReqUri    string //请求uri
	ReqBody   string //请求body
	RespData  string //返回数据

	TrafficStatus  string //通信状态(0失败，1成功)
	BusinessStatus string //业务处理装(0失败，1成功)
	CreateTime     string //创建时间
	AppId          string //应用id
	ReqTime        string //请求时间(时间戳)
}

func (ApiPayLogDao) GetList(reqData ApiPayLogDao) (string, []*ApiPayLogData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "app_id", Val: reqData.AppId, EqType: "="},
		{Key: "req_method", Val: reqData.ReqMethod, EqType: "="},
		{Key: "traffic_status", Val: reqData.TrafficStatus, EqType: "="},
		{Key: "business_status", Val: reqData.BusinessStatus, EqType: "="},

		{Key: "create_time", Val: reqData.StartTime, EqType: ">="},
		{Key: "create_time", Val: reqData.EndTime, EqType: "<="},
		{Key: "req_time", Val: reqData.ReqStartTime, EqType: ">="},
		{Key: "req_time", Val: reqData.ReqEndTime, EqType: "<="},

		{Key: "req_uri", Val: reqData.ReqUri, EqType: "begin like"},
		{Key: "req_body", Val: reqData.ReqBody, EqType: "like"},
		{Key: "resp_data", Val: reqData.RespData, EqType: "like"},
	})

	var total sql.NullString
	sqlCnt := "SELECT COUNT(1) FROM log_api_pay " + whereModel.WhereStr
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&total}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err)
		return "0", nil, err
	}

	if total.String == "" || total.String == "0" {
		return "0", nil, nil
	}

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY create_time DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(reqData.PageSize), strext.ToInt(reqData.Page))
	sqlStr := "SELECT trace_no, req_method, req_uri, req_body, resp_data, traffic_status, business_status, create_time, app_id, req_time " +
		" FROM log_api_pay "
	sqlStr += whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		return total.String, nil, err
	}

	var dataList []*ApiPayLogData
	for rows.Next() {
		var traceNo, reqMethod, reqUri, reqBody, respData,
			trafficStatus, businessStatus, createTime, appId,
			reqTime sql.NullString
		if err := rows.Scan(
			&traceNo, &reqMethod, &reqUri, &reqBody, &respData,
			&trafficStatus, &businessStatus, &createTime, &appId, &reqTime,
		); err != nil {
			ss_log.Error("err=[%v] traceNo=[%v]", err, traceNo.String)
			return total.String, nil, err
		}
		data := new(ApiPayLogData)
		data.TraceNo = traceNo.String
		data.ReqMethod = reqMethod.String
		data.ReqUri = reqUri.String
		data.ReqBody = reqBody.String
		data.RespData = respData.String

		data.TrafficStatus = trafficStatus.String
		data.BusinessStatus = businessStatus.String
		data.CreateTime = createTime.String
		data.AppId = appId.String
		data.ReqTime = reqTime.String

		dataList = append(dataList, data)
	}

	return total.String, dataList, nil
}
