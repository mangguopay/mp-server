package dao

import (
	"database/sql"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_sql"
)

var (
	ServicerDaoInstance ServicerDao
)

type ServicerDao struct {
}

func (ServicerDao) InsertInitService(tx *sql.Tx, accountNo string) (string, error) {
	//创建运营商
	servicer_no := strext.NewUUID()
	sqlStr := "insert into servicer(servicer_no, account_no, create_time) " +
		" values ($1,$2,current_timestamp)"
	return servicer_no, ss_sql.ExecTx(tx, sqlStr, servicer_no, accountNo)
}
