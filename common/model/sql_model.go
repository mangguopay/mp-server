package model

type WhereSql struct {
	WhereStr string
	I        int
	Args     []interface{}
}

type WhereSqlCond struct {
	EqType string
	Key    string
	Val    string
}

type WhereSqlOrderCond struct {
	Key   string
	IsAsc bool
}
