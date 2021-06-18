package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type ChannelDao struct {
}

var ChannelDaoInstance ChannelDao

func (*ChannelDao) QeuryChannelNoFromNameAndType(name, channelType string) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "channel_name", Val: name, EqType: "="},
		{Key: "channel_type", Val: channelType, EqType: "in"},
		{Key: "is_delete", Val: "0", EqType: "="},
	})

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " LIMIT 1 ")
	var channelNo sql.NullString
	err := ss_sql.QueryRow(dbHandler, "SELECT channel_no FROM channel "+whereModel.WhereStr, []*sql.NullString{&channelNo}, whereModel.Args...)
	if err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return ""
	}
	return channelNo.String

}
