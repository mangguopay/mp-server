package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_data"
	"database/sql"
	"errors"
)

type globalParamDao struct{}

var GlobalParamDaoInstance globalParamDao

// 全局配置项
func (*globalParamDao) GetGlobalParam(paramKey string) (redisKeyRet, value string, err error) {
	redisKey := ss_data.MkGlobalParamValue(paramKey)
	// ----------------第一版本的redis--------------------
	//tmp, err := ss_data.GetDataFromCache1st(paramKey, redisKey, cache.RedisCli, constants.DefPoolName, func(key string) (string, error) {

	// ----------------第二版本的redis--------------------
	tmp, err := cache.GetDataFromCache1st(paramKey, redisKey, cache.RedisClient, func(key string) (string, error) {
		dbHandler := db.GetDB(constants.DB_RISK)
		defer db.PutDB(constants.DB_RISK, dbHandler)
		selectSql := `SELECT r.rule FROM rela_api_event rav LEFT JOIN rela_event_rule rer ON rav.event_no = rer.event_no LEFT JOIN rule r ON rer.rule_no = r.rule_no WHERE rav.api_type = $1 and r.is_delete='0'`
		row := dbHandler.QueryRow(selectSql, key)
		tmp := sql.NullString{}
		row.Scan(&tmp)
		if tmp.Valid {
			return tmp.String, nil
		}
		ss_log.Error("not found param_value,key[%v]", redisKey)
		return "", errors.New("not found key")
	})
	return redisKey, strext.ToString(tmp), err
}
