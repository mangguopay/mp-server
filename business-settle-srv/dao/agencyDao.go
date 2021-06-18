package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/mp-server/common/ss_sql"
)

type AgencyDao struct {
}

var AgencyDaoInst AgencyDao

// 查询代理费率
func (*AgencyDao) QueryUpper(tx *sql.Tx, accNo string) string {
	var upperT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select upper from business_agency where acc_no=$1 and is_delete='0' limit 1`,
		[]*sql.NullString{&upperT}, accNo)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return upperT.String
}
