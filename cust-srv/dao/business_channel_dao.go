package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"fmt"
)

type BusinessChannel struct {
	ChannelNo          string
	ChannelName        string
	ChannelType        string
	ChannelWeiChatRate string
	ChannelAliPayRate  string
	UpstreamNo         string
	CreateTime         string
}

var BusinessChannelDao BusinessChannel

func (BusinessChannel) InsertTx(tx *sql.Tx, d *BusinessChannel) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "INSERT INTO business_channel (business_channel_no, channel_name, channel_type, upstream_no, create_time) " +
		"VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) "

	return ss_sql.ExecTx(tx, sqlStr, d.ChannelNo, d.ChannelName, d.ChannelType, d.UpstreamNo)
}

func (BusinessChannel) UpdateTx(tx *sql.Tx, d *BusinessChannel) error {
	sqlStr, args, _, err := ss_sql.MkUpdateSql("business_channel", map[string]string{
		"channel_name": d.ChannelName,
		"channel_type": d.ChannelType,
		"upstream_no":  d.UpstreamNo,
		"update_time":  ss_time.NowForPostgres(global.Tz),
	}, fmt.Sprintf("where business_channel_no = '%v' ", d.ChannelNo))

	if err != nil {
		return err
	}
	return ss_sql.ExecTx(tx, sqlStr, args...)
}

func (BusinessChannel) CntChannel(whereStr string, args []interface{}) (int64, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT COUNT(1) FROM business_channel "

	sqlStr += whereStr
	var totalNum sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalNum}, args...)
	if err != nil {
		return 0, err
	}

	return strext.ToInt64(totalNum.String), nil
}

//查询全部渠道
func (BusinessChannel) GetAllChannel(whereStr string, args []interface{}) ([]*BusinessChannel, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT business_channel_no, channel_name, channel_type, upstream_no, create_time " +
		"FROM business_channel "

	sqlStr += whereStr
	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Scan()

	var list []*BusinessChannel
	for rows.Next() {
		var channelNo, channelName, channelType, upstreamNo, createTime sql.NullString
		err := rows.Scan(&channelNo, &channelName, &channelType, &upstreamNo, &createTime)
		if err != nil {
			return nil, err
		}
		channel := new(BusinessChannel)
		channel.ChannelNo = channelNo.String
		channel.ChannelName = channelName.String
		channel.ChannelType = channelType.String
		channel.UpstreamNo = upstreamNo.String
		channel.CreateTime = createTime.String

		list = append(list, channel)
	}

	return list, nil
}

//查询渠道名称和上游渠道费率(微信，支付宝)
func (BusinessChannel) GetChannelRate(appId, sceneNo string) (*BusinessChannel, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bc.channel_name, bi.up_wx_rate, bi.up_alipay_rate " +
		"FROM business_channel bc " +
		"LEFT JOIN business_scene bs ON bs.business_channel_no = bc.business_channel_no " +
		"LEFT JOIN business_signed bsi ON bsi.scene_no = bs.scene_no " +
		"LEFT JOIN business_industry bi ON bi.code = bsi.industry_no " +
		"LEFT JOIN business_app app ON app.app_id = bsi.app_id " +
		"WHERE app.app_id = $1  " +
		"AND bs.scene_no = $2 "

	var channelName, wxRate, aliRate sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&channelName, &wxRate, &aliRate},
		appId, sceneNo)
	if err != nil {
		return nil, err
	}

	obj := new(BusinessChannel)
	obj.ChannelName = channelName.String
	obj.ChannelWeiChatRate = wxRate.String
	obj.ChannelAliPayRate = aliRate.String
	return obj, nil
}

//查询渠道名称和上游渠道费率(根据订单交易类型)
func (BusinessChannel) GetOutChannelRateAndName(appId, sceneNo string) (name, rate string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT bc.channel_name, " +
		"(CASE bc.channel_name " +
		"WHEN $1 THEN bi.up_wx_rate " +
		"WHEN $2 THEN bi.up_alipay_rate " +
		"ELSE 0 " +
		"END " +
		") channel_rate " +
		"FROM business_channel bc " +
		"LEFT JOIN business_scene bs ON bs.business_channel_no = bc.business_channel_no " +
		"LEFT JOIN business_signed bsi ON bsi.scene_no = bs.scene_no " +
		"LEFT JOIN business_app app ON app.app_id = bsi.app_id " +
		"LEFT JOIN business_industry bi ON bi.code = bsi.industry_no " +
		"WHERE app.app_id = $3 AND bs.scene_no = $4 AND bc.channel_type = $5 "

	var channelName, channelRate sql.NullString
	err = ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&channelName, &channelRate},
		constants.WeChatChannel, constants.AliPayChannel, appId, sceneNo, constants.ChannelTypeOut)
	if err != nil {
		return "", "", err
	}
	return channelName.String, channelRate.String, nil
}

func (BusinessChannel) GetChannelName(businessChannelNo string) (channelName string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT channel_name " +
		" FROM business_channel" +
		" WHERE business_channel_no = $1  "
	var name sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&name}, businessChannelNo); err != nil {
		return "", err
	}

	return name.String, nil
}
