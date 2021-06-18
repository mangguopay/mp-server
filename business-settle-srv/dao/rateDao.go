package dao

import (
	"database/sql"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_sql"
)

type RateDao struct {
}

var (
	RateDaoInst RateDao
)

// 查询费率
func (*RateDao) QueryRateFromAccNo(tx *sql.Tx, accNo, roleTye, channelNo string) (string, error) {
	var rateT sql.NullString
	err := ss_sql.QueryRowTx(tx, `select rate from business_rate where acc_no = $1 and role_type = $2 and channel_no = $3 and use_status = $4 limit 1`, []*sql.NullString{&rateT}, accNo, roleTye, channelNo, 1)
	if nil != err {
		ss_log.Error("err=%v", err)
		return "", err
	}
	return rateT.String, nil
}
