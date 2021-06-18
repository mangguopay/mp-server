package ss_sql

/***
 *    Todo 组装后的sql其实是可以塞到缓存里的，加个key作为索引就是了，但是现在还没有到那步，就先不管
 ***/

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/model"
	"database/sql"
	"encoding/json"
	"fmt"
)

type SsSqlFactory struct {
}

var (
	SsSqlFactoryInst SsSqlFactory
)

func (*SsSqlFactory) InitWhere() *model.WhereSql {
	m := &model.WhereSql{
		WhereStr: " where 1=1 ",
		I:        1,
		Args:     []interface{}{},
	}
	return m
}

func (r *SsSqlFactory) InitWhereList(c []*model.WhereSqlCond) *model.WhereSql {
	m := r.InitWhere()
	for _, v := range c {
		r.AppendWhere(m, v.Key, v.Val, v.EqType)
	}
	return m
}

func (r *SsSqlFactory) InitWhereListLimit(c []*model.WhereSqlCond, pageSize, pages int) *model.WhereSql {
	m := r.InitWhere()
	for _, v := range c {
		r.AppendWhere(m, v.Key, v.Val, v.EqType)
	}
	r.AppendWhereLimit(m, pageSize, pages)
	return m
}

func (r *SsSqlFactory) AppendWhereListNormal(m *model.WhereSql, c []*model.WhereSqlCond) {
	for _, v := range c {
		r.AppendWhere(m, v.Key, v.Val, v.EqType)
	}
}

func (r *SsSqlFactory) AppendWhereList(m *model.WhereSql, c []*model.WhereSqlCond, pageSize, pages int) {
	for _, v := range c {
		r.AppendWhere(m, v.Key, v.Val, v.EqType)
	}
}

func (r *SsSqlFactory) AppendWhere(m *model.WhereSql, key, val string, equalType string) {
	if val != "" {
		if val == "__empty_string" {
			val = ""
		}
		switch equalType {
		case "like":
			m.WhereStr = fmt.Sprintf("%s and %s like $%d ", m.WhereStr, key, m.I)
			m.Args = append(m.Args, "%"+val+"%")
		case "begin like":
			m.WhereStr = fmt.Sprintf("%s and %s like $%d ", m.WhereStr, key, m.I)
			m.Args = append(m.Args, val+"%")
		case "ending like":
			m.WhereStr = fmt.Sprintf("%s and %s like $%d ", m.WhereStr, key, m.I)
			m.Args = append(m.Args, "%"+val)
		case "in":
			m.WhereStr = fmt.Sprintf("%s and %s in %s ", m.WhereStr, key, val)
		case "not in":
			m.WhereStr = fmt.Sprintf("%s and %s not in %s ", m.WhereStr, key, val)
		case "or_group":
			data := []model.WhereSqlCond{}
			err := json.Unmarshal([]byte(val), &data)
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}
			m.WhereStr += " and ( 1="
			w := "0"
			for _, v := range data {
				w = r.AppendWhereOrGroup(m, v.Key, v.Val, v.EqType, w)
			}
			if w == "" {
				m.WhereStr += " ) "
			} else {
				m.WhereStr += "1 ) "
			}
		default:
			m.WhereStr = fmt.Sprintf("%s and %s %s $%d ", m.WhereStr, key, equalType, m.I)
			m.Args = append(m.Args, val)
		}
		ss_log.Info("key=[%v],equalType=[%v],val=[%v],i=[%v]", key, equalType, val, m.I)
		if equalType != "not in" && equalType != "or_group" && equalType != "in" {
			m.I = m.I + 1
		}
	}
}

func (*SsSqlFactory) AppendWhereOrGroup(m *model.WhereSql, key, val string, equalType, ext string) (newExt string) {
	if val != "" {
		if ext != "" {
			m.WhereStr += ext
		}
		switch equalType {
		case "like":
			m.WhereStr = fmt.Sprintf("%s or %s like $%d ", m.WhereStr, key, m.I)
			m.Args = append(m.Args, "%"+val+"%")
		default:
			m.WhereStr = fmt.Sprintf("%s or %s %s $%d ", m.WhereStr, key, equalType, m.I)
			m.Args = append(m.Args, val)
		}
		ss_log.Info("key=[%v],equalType=[%v],val=[%v],i=[%v]", key, equalType, val, m.I)
		m.I = m.I + 1
		ext = ""
		return ext
	}
	return ext
}

func (*SsSqlFactory) AppendWhereLimit(m *model.WhereSql, pageSize, pages int) {
	if pages == 0 {
		pages = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	m.Args = append(m.Args, pageSize, (pages-1)*pageSize)
	m.WhereStr = fmt.Sprintf("%s limit $%d offset $%d ", m.WhereStr, m.I, m.I+1)
	m.I += 2
}

func (*SsSqlFactory) AppendWhereLimitI32(m *model.WhereSql, pageSize, pages int32) {
	if pages == 0 {
		pages = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	m.Args = append(m.Args, pageSize, (pages-1)*pageSize)
	m.WhereStr = fmt.Sprintf("%s limit $%d offset $%d ", m.WhereStr, m.I, m.I+1)
	m.I += 2
}

func (*SsSqlFactory) AppendWhereExtra(m *model.WhereSql, extraStr string) {
	m.WhereStr = fmt.Sprintf("%s %s ", m.WhereStr, extraStr)
}

func (*SsSqlFactory) AppendWhereOrderBy(m *model.WhereSql, orderKey string, isAsc bool) {
	if isAsc {
		m.WhereStr = fmt.Sprintf("%s order by %s asc ", m.WhereStr, orderKey)
	} else {
		m.WhereStr = fmt.Sprintf("%s order by %s desc ", m.WhereStr, orderKey)
	}
}

func (r *SsSqlFactory) AppendWhereOrderByList(m *model.WhereSql, orders []model.WhereSqlOrderCond) {
	m.WhereStr = fmt.Sprintf("%s order by ", m.WhereStr)
	for _, v := range orders {
		m.WhereStr = fmt.Sprintf("%s %s %s,", m.WhereStr, v.Key, r.isAsc(v.IsAsc))
	}
	m.WhereStr = m.WhereStr[:len(m.WhereStr)-1] + " "
}

func (*SsSqlFactory) isAsc(isAsc bool) string {
	if isAsc {
		return "asc"
	} else {
		return "desc"
	}
}

func (*SsSqlFactory) Query(dbHandler *sql.DB, sqlStr string, m *model.WhereSql) (*sql.Rows, *sql.Stmt, error) {
	return db.Query(dbHandler, sqlStr+m.WhereStr, ss_log.Error, m.Args...)
}
func (*SsSqlFactory) QueryRow(dbHandler *sql.DB, sqlStr string, retStr []*sql.NullString, m *model.WhereSql) error {
	return db.QueryRow(dbHandler, sqlStr+m.WhereStr, ss_log.Error, retStr, m.Args...)
}

// 获取个数
func (*SsSqlFactory) GetCnt(dbHandler *sql.DB, tabname string, m *model.WhereSql) int64 {
	var cnt sql.NullString
	err := db.QueryRow(dbHandler, fmt.Sprintf(`SELECT count(1) FROM %s %s`, tabname, m.WhereStr), ss_log.Error, []*sql.NullString{&cnt}, m.Args...)
	if err != nil {
		return 0
	}
	return strext.ToInt64(cnt.String)
}

func (r *SsSqlFactory) DeepClone(m *model.WhereSql) *model.WhereSql {
	x := r.InitWhere()
	x.Args = m.Args
	x.I = m.I
	x.WhereStr = m.WhereStr
	return x
}
