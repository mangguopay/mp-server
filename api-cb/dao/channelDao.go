package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-cb/m"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"encoding/json"
)

type ChannelDao struct {
}

var ChannelDaoInst ChannelDao

// 获取通道参数
func (ChannelDao) GetChannelParam(channelNo string) (channelParam *m.ChannelParam) {
	redisKey := cache.MkChannelParam(channelNo)
	tmp, err := cache.GetDataFromCache1stI([]string{channelNo}, redisKey, func(params []string) (interface{}, error) {
		dbHandler := db.GetDB(constants.DB_CRM)
		defer db.PutDB(constants.DB_CRM, dbHandler)

		var orgNo, mercNo, gateway, signKey, key1, key2, key3, key4, key5, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10 sql.NullString
		err := ss_sql.QueryRow(dbHandler, `select org_no,merc_no,gateway,sign_key,key1,key2,key3,key4,key5,p1,p2,p3,p4,p5,p6,p7,p8,p9,p10 from business_channel_param where channel_no=$1 limit 1`,
			[]*sql.NullString{&orgNo, &mercNo, &gateway, &signKey, &key1, &key2, &key3, &key4, &key5, &p1, &p2, &p3, &p4, &p5, &p6, &p7, &p8, &p9, &p10}, channelNo)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return "", err
		}

		m2 := m.ChannelParam{
			OrgNo:   orgNo.String,
			MercNo:  mercNo.String,
			Gateway: gateway.String,
			SignKey: signKey.String,
			Key1:    key1.String,
			Key2:    key2.String,
			Key3:    key3.String,
			Key4:    key4.String,
			Key5:    key5.String,
			P1:      p1.String,
			P2:      p2.String,
			P3:      p3.String,
			P4:      p4.String,
			P5:      p5.String,
			P6:      p6.String,
			P7:      p7.String,
			P8:      p8.String,
			P9:      p9.String,
			P10:     p10.String,
		}

		return m2, nil
	})
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil
	}

	ss_log.Info("tmp=[%v]", tmp)

	switch v := tmp.(type) {
	case string:
		ch := m.ChannelParam{}
		err = json.Unmarshal([]byte(strext.ToStringNoPoint(tmp)), &ch)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return nil
		}

		return &ch
	case m.ChannelParam:
		return &v
	}
	return nil
}
