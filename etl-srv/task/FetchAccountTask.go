package task

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/etl-srv/extract"
	"a.a/mp-server/etl-srv/load"
	"a.a/mp-server/etl-srv/m"
	"a.a/mp-server/etl-srv/transform"
)

type FetchAccountTask struct {
}

var (
	FetchAccountTaskInst FetchAccountTask
)

func (FetchAccountTask) Do() *m.TaskContext {
	taskParam := m.TaskParam{}
	taskParam.Task = extract.PostgresExtractorInst
	taskParam.Dbname = constants.DB_CRM
	taskParam.SqlStr = `select uid,account,create_time from account limit 1`
	taskParam.Args = nil
	taskParam.DataCnt = 3
	taskParam.Keys = []string{"account_no", "account", "create_time"}

	ctx := m.TaskContext{}
	ctx.Extract = &taskParam

	l := []*m.TaskParam{}
	taskParamT := m.TaskParam{}
	taskParamT.Task = transform.RemoveFieldsTransInst
	taskParamT.Keys = []string{"create_time"}
	l = append(l, &taskParamT)

	ctx.Transform = l

	l2 := []*m.TaskParam{}
	taskParamL := m.TaskParam{}
	taskParamL.Task = load.PostgresInsertLoaderInst
	taskParamL.Dbname = constants.DbStat
	taskParamL.SqlStr = `insert into account(uid,nickname)values($1,$2)`
	taskParamL.Keys = []string{
		"account_no", "account",
	}
	taskParamL.ArgsType = m.ArgsType_DataMap
	l2 = append(l2, &taskParamL)

	ctx.Load = l2
	ctx.LoadPc = 0
	ctx.TransformPc = 0
	return &ctx
}
