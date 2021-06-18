package load

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/etl-srv/m"
)

type PostgresInsertLoader struct {
}

var PostgresInsertLoaderInst PostgresInsertLoader

func (PostgresInsertLoader) Do(ctx *m.TaskContext) {
	//func (PostgresLoader) Do(dbname string, sqlStr string, args []interface{}) {
	dbHandler := db.GetDB(ctx.Load[ctx.LoadPc].Dbname)
	defer db.PutDB(ctx.Load[ctx.LoadPc].Dbname, dbHandler)

	l := []interface{}{}
	switch ctx.Load[ctx.LoadPc].ArgsType {
	case m.ArgsType_DataMap:
		for _, v := range ctx.Load[ctx.LoadPc].Keys {
			l = append(l, ctx.DataMap[v])
		}
	}

	if err := ss_sql.Exec(dbHandler, ctx.Load[ctx.LoadPc].SqlStr, l...); err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}

	return
}
