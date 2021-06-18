package cache

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_data"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type ApiDao struct {
}

var (
	ApiDaoInstance ApiDao
)

// 全局配置项
func (*ApiDao) GetGlobalParam(paramKey string) (redisKeyRet, value string, err error) {
	redisKey := ss_data.MkGlobalParamValue(paramKey)
	// ----------------第二版本的redis--------------------
	tmp, err := GetDataFromCache1st(paramKey, redisKey, RedisClient, func(key string) (string, error) {
		dbHandler := db.GetDB(constants.DB_CRM)
		defer db.PutDB(constants.DB_CRM, dbHandler)
		selectSql := `select param_value from global_param where param_key=$1 limit 1`
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
