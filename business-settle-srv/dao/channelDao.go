package dao

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ChannelDao struct{}

var ChannelDaoInst ChannelDao

func (*ChannelDao) ForbidChannelFromNo(tx *sql.Tx, channelNo, channelStatus string) string {

	sqlStr := "update business_channel set interface_status = $1 where channel_no = $2 and is_delete='0 ' "
	err := ss_sql.ExecTx(tx, sqlStr, channelStatus, channelNo)
	if err != nil {
		ss_log.Error("%s", err.Error())
		return ss_err.ERR_PARAM
	}
	return ss_err.ERR_SUCCESS
}
