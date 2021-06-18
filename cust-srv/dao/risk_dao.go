package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	RiskDaoInstance RiskDao
)

type RiskDao struct {
}

/**********************************************  event  ************************************************************/
func (RiskDao) GetEventCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM event " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetEventInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.EventData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT event_no, event_name, create_time " +
		" FROM event  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.EventData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.EventData{}
		err2 := rows.Scan(
			&data.EventNo,
			&data.EventName,
			&data.CreateTime,
		)
		if err2 != nil {
			ss_log.Error("event表查询出错。EventNo：[%v],err:[%v]", data.EventNo, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) GetEventDetail(whereModelStr string, whereArgs []interface{}) (dataT *go_micro_srv_cust.EventData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT event_no, event_name, create_time " +
		" FROM event  " + whereModelStr
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	data := &go_micro_srv_cust.EventData{}
	err2 := rows.Scan(
		&data.EventNo,
		&data.EventName,
		&data.CreateTime,
	)
	if err2 != nil {
		ss_log.Error("push_conf表查询出错。EventNo：[%v],err:[%v]", data.EventNo, err2)
	}

	return data, err
}

func (*RiskDao) ModifyEvent(tx *sql.Tx, eventNo, eventName string) (err error) {

	sqlStr := "update event set event_name = $2 where event_no = $1 "
	errT := ss_sql.ExecTx(tx, sqlStr, eventNo, eventName)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) AddEvent(tx *sql.Tx, eventName string) (eventNoT string, err error) {

	eventNo := strext.NewUUID()
	sqlStr := "insert into event(event_no, event_name, create_time) values($1,$2,current_timestamp)"
	errT := ss_sql.ExecTx(tx, sqlStr, eventNo, eventName)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return eventNo, nil
}

func (*RiskDao) DelEvent(tx *sql.Tx, eventNo string) (err error) {
	//sqlStr := "delete from event where event_no = $1 "
	sqlStr := "update event set is_delete = '1' where event_no = $1 and  is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, eventNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

/**********************************************  eva_param  ************************************************************/
func (RiskDao) GetEvaParamsCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM eva_param " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetEvaParamsInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.EvaParamsData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT key,val " +
		" FROM eva_param  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.EvaParamsData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.EvaParamsData{}
		err2 := rows.Scan(
			&data.Key,
			&data.Val,
		)
		if err2 != nil {
			ss_log.Error("g表查询出错。key：[%v],val:[%v],err:[%v]", data.Key, data.Val, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) AddOrUpdateEvaParam(tx *sql.Tx, key, val string) (err error) {

	sqlStr := "insert into eva_param(key,val)values($1,$2) on conflict (key) do update set val=$2"
	errT := ss_sql.ExecTx(tx, sqlStr, key, val)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) DelEvaParam(tx *sql.Tx, key string) (err error) {
	sqlStr := "delete from eva_param where key = $1 "
	errT := ss_sql.ExecTx(tx, sqlStr, key)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//global_param
func (RiskDao) GetGlobalParamCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM global_param " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetGlobalParamInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.GlobalParamData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT param_key,param_value,remark " +
		" FROM global_param  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.GlobalParamData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.GlobalParamData{}
		err2 := rows.Scan(
			&data.ParamKey,
			&data.ParamValue,
			&data.Remark,
		)
		if err2 != nil {
			ss_log.Error("查询出错。ParamKey：[%v],err:[%v]", data.ParamKey, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) AddOrUpdateGlobalParam(tx *sql.Tx, paramKey, paramValue, remark string) (err error) {
	sqlStr := "insert into global_param(param_key,param_value,remark)values($1,$2,$3) on conflict (param_key) do update set param_value=$2,remark=$3"
	errT := ss_sql.ExecTx(tx, sqlStr, paramKey, paramValue, remark)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) DelGlobalParam(tx *sql.Tx, paramKey string) (err error) {
	sqlStr := "delete from global_param where param_key = $1 "
	errT := ss_sql.ExecTx(tx, sqlStr, paramKey)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//log_result
func (RiskDao) GetLogResultsCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM log_result " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetLogResultsInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.LogResultsData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT log_no,event_no,rule_no,op_no,result " +
		" FROM log_result  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.LogResultsData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.LogResultsData{}
		err2 := rows.Scan(
			&data.LogNo,
			&data.EventNo,
			&data.RuleNo,
			&data.OpNo,
			&data.Result,
		)
		if err2 != nil {
			ss_log.Error("查询出错。LogNo：[%v],err:[%v]", data.LogNo, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

//op
func (RiskDao) GetOpsCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM op " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetOpsInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.OpData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT op_no,op_name,script_name,param,score " +
		" FROM op  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.OpData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.OpData{}
		var opName, scriptName, param, score sql.NullString
		err2 := rows.Scan(
			&data.OpNo,
			&opName,
			&scriptName,
			&param,
			&score,
		)
		if err2 != nil {
			ss_log.Error("查询出错。OpNo：[%v],err:[%v]", data.OpNo, err2)
			continue
		}

		if opName.String != "" {
			data.OpName = opName.String
		}
		if scriptName.String != "" {
			data.ScriptName = scriptName.String
		}
		if param.String != "" {
			data.Param = param.String
		}
		if score.String != "" {
			data.Score = score.String
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) GetOpDetail(whereModelStr string, whereArgs []interface{}) (dataT *go_micro_srv_cust.OpData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT op_no,op_name,script_name,param,score " +
		" FROM op  " + whereModelStr
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	data := &go_micro_srv_cust.OpData{}
	var opName, scriptName, param, score sql.NullString
	err2 := rows.Scan(
		&data.OpNo,
		&opName,
		&scriptName,
		&param,
		&score,
	)
	if err2 != nil {
		ss_log.Error("查询出错。OpNo：[%v],err:[%v]", data.OpNo, err2)
	}

	if opName.String != "" {
		data.OpName = opName.String
	}
	if scriptName.String != "" {
		data.ScriptName = scriptName.String
	}
	if param.String != "" {
		data.Param = param.String
	}
	if score.String != "" {
		data.Score = score.String
	}

	return data, err
}

func (*RiskDao) ModifyOp(tx *sql.Tx, opNo, opName, scriptName, param, score string) (err error) {
	sqlStr := "update op set op_name = $2, script_name = $3, param = $4, score = $5 where op_no = $1 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, opNo, opName, scriptName, param, score)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) AddOp(tx *sql.Tx, opName, scriptName, param, score string) (eventNoT string, err error) {

	opNo := strext.NewUUID()
	sqlStr := "insert into op(op_no, op_name, script_name, param, score) values($1,$2,$3,$4,$5)"
	errT := ss_sql.ExecTx(tx, sqlStr, opNo, opName, scriptName, param, score)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return opNo, nil
}

func (*RiskDao) DelOp(tx *sql.Tx, opNo string) (err error) {
	sqlStr := "update op set is_delete = '1' where op_no = $1 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, opNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//rela_api_event
func (RiskDao) GetGetRelaApiEventsCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM rela_api_event " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetGetRelaApiEventsInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.RelaApiEventData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT api_type, event_no, create_time" +
		" FROM rela_api_event  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.RelaApiEventData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.RelaApiEventData{}
		err2 := rows.Scan(
			&data.ApiType,
			&data.EventNo,
			&data.CreateTime,
		)
		if err2 != nil {
			ss_log.Error("查询出错。ApiType：[%v],err:[%v]", data.ApiType, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) AddOrUpdateRelaApiEvent(tx *sql.Tx, apiType, eventNo string) (err error) {
	sqlStr := "insert into rela_api_event(api_type,event_no,create_time)values($1,$2,current_timestamp) on conflict (api_type) do update set event_no=$2 "
	errT := ss_sql.ExecTx(tx, sqlStr, apiType, eventNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) DelRelaApiEvent(tx *sql.Tx, key string) (err error) {
	sqlStr := "delete from rela_api_event where api_type = $1 "
	errT := ss_sql.ExecTx(tx, sqlStr, key)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//rela_event_rule
func (RiskDao) GetRelaEventRulesCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM rela_event_rule " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetRelaEventRulesInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.RelaEventRuleData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT rela_no, event_no, rule_no, create_time " +
		" FROM rela_event_rule  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.RelaEventRuleData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.RelaEventRuleData{}
		err2 := rows.Scan(
			&data.RelaNo,
			&data.EventNo,
			&data.RuleNo,
			&data.CreateTime,
		)
		if err2 != nil {
			ss_log.Error("查询出错。RelaNo：[%v],err:[%v]", data.RelaNo, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) ModifyRelaEventRule(tx *sql.Tx, relaNo, eventNo, ruleNo string) (err error) {
	sqlStr := "update rela_event_rule set event_no = $2, rule_no = $3 where rela_no = $1 "
	errT := ss_sql.ExecTx(tx, sqlStr, relaNo, eventNo, ruleNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) AddRelaEventRule(tx *sql.Tx, eventNo, ruleNo string) (relaNoT string, err error) {
	relaNo := strext.NewUUID()
	sqlStr := "insert into rela_event_rule(rela_no, event_no, rule_no, create_time) values($1,$2,$3,current_timestamp)"
	errT := ss_sql.ExecTx(tx, sqlStr, relaNo, eventNo, ruleNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return relaNo, nil
}

func (*RiskDao) DelRelaEventRule(tx *sql.Tx, relaNo string) (err error) {
	sqlStr := "delete from rela_event_rule where rela_no = $1  "
	errT := ss_sql.ExecTx(tx, sqlStr, relaNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//risk_threshold
func (RiskDao) GetRiskThresholdsCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM risk_threshold " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetRiskThresholdsInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.RiskThresholdData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT rule_no, risk_threshold, create_time " +
		" FROM risk_threshold  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.RiskThresholdData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.RiskThresholdData{}
		err2 := rows.Scan(
			&data.RuleNo,
			&data.RiskThreshold,
			&data.CreateTime,
		)
		if err2 != nil {
			ss_log.Error("查询出错。RuleNo：[%v],err:[%v]", data.RuleNo, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) AddOrUpdateRiskThreshold(tx *sql.Tx, ruleNo, riskThreshold string) (err error) {
	sqlStr := "insert into risk_threshold(rule_no, risk_threshold, create_time)values($1,$2,current_timestamp) on conflict (rule_no) do update set risk_threshold=$2 "
	errT := ss_sql.ExecTx(tx, sqlStr, ruleNo, riskThreshold)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) DelRiskThreshold(tx *sql.Tx, ruleNo string) (err error) {
	sqlStr := "update risk_threshold set is_delete = '1' where rule_no = $1 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, ruleNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//rule
func (RiskDao) GetRulesCnt(whereStr string, whereArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)
	//全部数量统计
	var totalT sql.NullString
	sqlCnt := "SELECT count(1) " +
		" FROM rule " + whereStr
	cnterr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereArgs...)
	if cnterr != nil {
		ss_log.Error("err=[%v]", cnterr)
		return "0", cnterr
	}
	return totalT.String, nil
}

func (*RiskDao) GetRulesInfos(whereModelStr string, whereArgs []interface{}) (datasT []*go_micro_srv_cust.RuleData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT rule_no, rule_name, rule,	create_time	" +
		" FROM rule  " + whereModelStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datas := []*go_micro_srv_cust.RuleData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	for rows.Next() {
		data := &go_micro_srv_cust.RuleData{}
		err2 := rows.Scan(
			&data.RuleNo,
			&data.RuleName,
			&data.Rule,
			&data.CreateTime,
		)
		if err2 != nil {
			ss_log.Error("查询出错。RuleNo：[%v],err:[%v]", data.RuleNo, err2)
			continue
		}

		datas = append(datas, data)
	}

	return datas, err
}

func (*RiskDao) GetRuleDetail(whereModelStr string, whereArgs []interface{}) (dataT *go_micro_srv_cust.RuleData, err error) {
	dbHandler := db.GetDB(constants.DB_RISK)
	defer db.PutDB(constants.DB_RISK, dbHandler)

	sqlStr := "SELECT rule_no, rule_name, rule,	create_time	" +
		" FROM rule  " + whereModelStr
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	data := &go_micro_srv_cust.RuleData{}
	err2 := rows.Scan(
		&data.RuleNo,
		&data.RuleName,
		&data.Rule,
		&data.CreateTime,
	)
	if err2 != nil {
		ss_log.Error("rule表查询出错。RuleNo：[%v],err:[%v]", data.RuleNo, err2)
	}

	return data, err
}

func (*RiskDao) ModifyRule(tx *sql.Tx, ruleNo, ruleName, rule string) (err error) {
	sqlStr := "update rule set rule_name = $2, rule = $3 where rule_no = $1 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, ruleNo, ruleName, rule)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*RiskDao) AddRule(tx *sql.Tx, ruleName, rule string) (ruleNoT string, err error) {
	ruleNo := strext.NewUUID()
	sqlStr := "insert into rule(rule_no, rule_name, rule, create_time) values($1,$2,$3,current_timestamp)"
	errT := ss_sql.ExecTx(tx, sqlStr, ruleNo, ruleName, rule)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return ruleNo, nil
}

func (*RiskDao) DelRule(tx *sql.Tx, ruleNo string) (err error) {
	sqlStr := "update rule set is_delete = '1' where rule_no = $1 and is_delete = '0' "
	errT := ss_sql.ExecTx(tx, sqlStr, ruleNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//
