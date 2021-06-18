package extract

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/etl-srv/m"
	"database/sql"
)

type PostgresExtractor struct {
}

var PostgresExtractorInst PostgresExtractor

func (PostgresExtractor) Do(ctx *m.TaskContext) {
	dbHandler := db.GetDB(ctx.Extract.Dbname)
	defer db.PutDB(ctx.Extract.Dbname, dbHandler)

	l := []*sql.NullString{}
	for i := 0; i < ctx.Extract.DataCnt; i++ {
		l = append(l, &sql.NullString{})
	}

	if ctx.Extract.Args == nil {
		if err := ss_sql.QueryRow(dbHandler, ctx.Extract.SqlStr, l); err != nil {
			ss_log.Error("err=[%v]", err)
			return
		}
	} else {
		if err := ss_sql.QueryRow(dbHandler, ctx.Extract.SqlStr, l, ctx.Extract.Args...); err != nil {
			ss_log.Error("err=[%v]", err)
			return
		}
	}

	l2 := map[string]interface{}{}
	for i, v := range l {
		l2[ctx.Extract.Keys[i]] = v.String
	}

	ctx.DataMap = l2
	return
}
