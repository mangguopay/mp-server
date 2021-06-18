package dao

import (
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type riskResultDao struct{}

var RiskResultDaoInstance riskResultDao

func (*riskResultDao) InsertResult(riskNo, riskResult, riskThreshold, apiType, payerAccNo, actionTime, evaExecuteType, evaScore, moneyType, orderNo, productType string, score int32) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	err := ss_sql.Exec(dbHandler, `insert into risk_result(risk_no,risk_result,risk_threshold,create_time,api_type,payer_acc_no,action_time,eva_execute_type,eva_score,money_type,order_no,score,product_type)
				values($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		riskNo, riskResult, riskThreshold, apiType, payerAccNo, actionTime, evaExecuteType, evaScore, moneyType, orderNo, score, productType)
	if err != nil {
		ss_log.Error("err=%v", err)
		return ss_err.ERR_PARAM
	}

	return ss_err.ERR_SUCCESS
}

type InsertResultNewData struct {
	RiskNo     string
	Result     int
	Threshold  int
	Position   string
	ActionTime string
	Uid        string
	Score      int
	Params     string
	ItemResult string
}

func (*riskResultDao) InsertResultNew(data InsertResultNewData) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `INSERT INTO risk_result_new (risk_no, result, threshold, position, action_time, uid, score, params, item_result, create_time)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, current_timestamp)`

	err := ss_sql.Exec(dbHandler, sqlStr,
		data.RiskNo, data.Result, data.Threshold, data.Position, data.ActionTime,
		data.Uid, data.Score, data.Params, data.ItemResult,
	)

	return err
}

// 查找风控结果.
func (riskResultDao) GetRiskResult(riskNo string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var riskResultT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select risk_result from risk_result where risk_no =$1 limit 1`,
		[]*sql.NullString{&riskResultT}, riskNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return riskResultT.String
}
