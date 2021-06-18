package dao

import (
	"database/sql"
	"errors"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

var PushConfigDaoInst PushConfigDao

type PushConfigDao struct {
}

type PushConfigInfo struct {
	Pusher         string
	Config         map[string]interface{}
	ConditionType  int
	ConditionValue string
}

func (PushConfigDao) GetPushConfigList() ([]*PushConfigInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	var list []*PushConfigInfo

	sqlStr := "SELECT pusher, config FROM push_conf WHERE use_status= '1' "

	rows, stmt, qErr := ss_sql.Query(dbHandler, sqlStr)
	if stmt != nil {
		defer stmt.Close()
		defer rows.Close()
	}

	if qErr != nil {
		if qErr.Error() != ss_sql.DB_NO_ROWS_MSG {
			return list, qErr
		}
		return list, nil
	}

	for rows.Next() {
		var pusher, config sql.NullString
		err := rows.Scan(&pusher, &config)
		if err != nil {
			ss_log.Error("err=%v", err)
			continue
		}
		pushConfigm := strext.Json2Map(config.String)
		list = append(list, &PushConfigInfo{
			Pusher: pusher.String,
			Config: pushConfigm,
		})
	}

	return list, nil
}

func (PushConfigDao) GetPushConfig(pusherNo string) (*PushConfigInfo, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var pusher, pushConfig, conditionType, conditionValue, useStatus sql.NullString
	sqlStr := "SELECT pusher,config,condition_type,condition_value,use_status FROM push_conf WHERE pusher_no=$1 limit 1 "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&pusher, &pushConfig, &conditionType, &conditionValue, &useStatus}, pusherNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	if useStatus.String != "1" {
		return nil, errors.New("NotUsable")
	}

	pushConfigm := strext.Json2Map(pushConfig.String)
	return &PushConfigInfo{Config: pushConfigm, Pusher: pusher.String, ConditionType: strext.ToInt(conditionType.String), ConditionValue: conditionValue.String}, nil
}
