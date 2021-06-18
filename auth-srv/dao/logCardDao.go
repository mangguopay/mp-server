package dao

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type LogCardDao struct{}

var LogCardDaoInstance LogCardDao

func (*LogCardDao) InsertLogCard(tx *sql.Tx, cardNum, name, accountNo, channelNo, channelType, descript string, vaType int) string {
	logNOT := strext.GetDailyId()
	err := ss_sql.ExecTx(tx, `insert into log_card(log_no,card_num,name,account_no,va_type,`+
		`channel_no,channel_type,descript,create_time) values($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)`,
		logNOT, cardNum, name, accountNo, vaType, channelNo, channelType, descript)
	if nil != err {
		ss_log.Error("err=%v", err)
		return ""
	}
	return ss_err.ERR_SUCCESS
}
