package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"strings"
)

type ApiDao struct {
}

var (
	ApiDaoInstance ApiDao
)

//func (*ApiDao) GetParamValue(channelNo, redisKey, key string) (value string, err error) {
//	tmp, err := cache.GetDataFromCache1st(key, redisKey, cache.RedisClient, func(key2 string) (string, error) {
//		dbHandler := db.GetDB(constants.DB_CRM)
//		defer db.PutDB(constants.DB_CRM, dbHandler)
//
//		var tmp sql.NullString
//		selectSql := `select "param_value" from interface_param where channel_no=$1 and "param_key"=$2 limit 1`
//		err := ss_sql.QueryRow(dbHandler, selectSql, []*sql.NullString{&tmp}, channelNo, key2)
//		if err != nil {
//			ss_log.Error("err=[%v]", err)
//		}
//		if tmp.Valid {
//			return tmp.String, nil
//		}
//		ss_log.Error("not found channelNo,key[%v|%v]", channelNo, key2)
//		return "", errors.New("not found key")
//	})
//	return strext.ToString(tmp), err
//}

//func (*ApiDao) GetParamValueWithBizType(channelNo string, bizType int32, key string) (redisKeyRet, value string, err error) {
//	redisKey := ss_data.MkParamCacheKey(key, bizType, channelNo)
//	tmp, err := ss_data.GetDataFromCache1st(key, redisKey, cache.RedisCli, "a", func(key2 string) (string, error) {
//		dbHandler := db.GetDB(constants.DB_CRM)
//		defer db.PutDB(constants.DB_CRM, dbHandler)
//
//		var tmp sql.NullString
//		selectSql := `select "param_value" from interface_param where channel_no=$1 and interface_biz=$2 and "param_key"=$3 limit 1`
//		err := ss_sql.QueryRow(dbHandler, selectSql, []*sql.NullString{&tmp}, channelNo, bizType, key2)
//		if err != nil {
//			ss_log.Error("err=[%v]", err)
//		}
//		if tmp.Valid {
//			return tmp.String, nil
//		}
//		ss_log.Error("not found channelNo,key[%v|%v]", channelNo, key2)
//		return "", errors.New("not found key")
//	})
//	return redisKey, strext.ToString(tmp), err
//}

//func (*ApiDao) IsOpenedApi(reqUrl, accNo string) (redisKeyRet string, value string, err error) {
//	redisKey := ss_data.MkRelaApi(reqUrl, accNo)
//	tmp, err := ss_data.GetDataFromCache1stL([]string{reqUrl, accNo}, redisKey, cache.RedisCli, "a", func(params []string) (string, error) {
//		dbHandler := db.GetDB(constants.DB_CRM)
//		defer db.PutDB(constants.DB_CRM, dbHandler)
//
//		var useStatus sql.NullInt64
//		selectSql := `select count(1) from api a left join rela_acc_api rela on rela.api_url=a.url where rela.acc_no=$2 and a.url=$1`
//		row, stmt, err := ss_sql.QueryRowN(dbHandler, selectSql, params[0], params[1])
//		if err != nil {
//			ss_log.Error("err=[%v]", err)
//			return "0", errors.New("not found key")
//		}
//		if stmt != nil {
//			defer stmt.Close()
//		}
//		err = row.Scan(&useStatus)
//		if err != nil {
//			ss_log.Error("err=[%v]", err)
//			return "0", errors.New("not found key")
//		}
//		ss_log.Error("err=[%v|%v]", err, useStatus.Int64)
//		if useStatus.Valid && useStatus.Int64 > 0 {
//			return "1", nil
//		}
//		return "0", errors.New("not found key")
//	})
//
//	return redisKey, tmp, err
//}

//
func (*ApiDao) GetMercSignInfo(accNo string) (redisKeyRet, signMethod, apiMode string, err error) {
	redisKey := cache.MkMercApiMode(accNo)
	tmp, err := cache.GetDataFromCache1st(accNo, redisKey, cache.RedisClient, func(key string) (string, error) {
		dbHandler := db.GetDB(constants.DB_CRM)
		defer db.PutDB(constants.DB_CRM, dbHandler)

		var signMethodT, apiModeT sql.NullString
		selectSql := `select sign_method,api_mode from business_acc_sign where acc_no=$1 limit 1`
		err := ss_sql.QueryRow(dbHandler, selectSql, []*sql.NullString{&signMethodT, &apiModeT}, key)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		if signMethodT.Valid {
			return strings.Join([]string{signMethodT.String, apiModeT.String}, ","), nil
		}
		ss_log.Error("not found sign_method,key[%v]", redisKey)
		return "", errors.New("not found key")
	})
	tmpS := strext.ToString(tmp)
	if tmpS != "" {
		tmpSList := strings.Split(tmpS, ",")
		return redisKey, tmpSList[0], tmpSList[1], nil
	}
	return redisKey, "", "", err
}

func (*ApiDao) GetMd5Key(accNo string) (md5Key string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var md5KeyT sql.NullString
	err := ss_sql.QueryRow(dbHandler, `select md5_key from business_acc_sign where acc_no=$1 limit 1`, []*sql.NullString{&md5KeyT}, accNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return ""
	}
	return md5KeyT.String
}
