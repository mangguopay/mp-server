package cache

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_data"
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func InitDb() {
	db.DoDBInitPostgres(constants.DB_CRM, "10.41.1.241", "5432", "postgres", "123", "mp_crm")
}
func InitRedisTest() {
	err := InitRedis("10.41.1.241", "6379", "123456a")
	fmt.Println("InitRedis-err:", err)
}

func TestInitRedis(t *testing.T) {
	InitRedisTest()

	fmt.Printf("RedisClient:%+v \n", RedisClient)
}

func TestCheckSMS(t *testing.T) {
	InitRedisTest()

	function := "xxxxxx"
	phone := "123456"
	msg := "dddsd"

	ok := CheckSMS(function, phone, msg)

	fmt.Println("ok:", ok)
}

func TestSetSMSMsgToCache1stV2(t *testing.T) {
	InitRedisTest()
	key := "kkkkk"
	val := "vvvvv"
	err := SetSMSMsgToCache1stV2(key, val, RedisClient)
	fmt.Println("err:", err)
}

func TestGetSMSMsgFromCache1stV2(t *testing.T) {
	InitRedisTest()
	key := "kkkkk"

	value, err := GetSMSMsgFromCache1stV2(key, RedisClient)
	fmt.Println("err:", err)
	fmt.Println("value:", value)
}

func TestGetDataFromCache1st(t *testing.T) {
	InitRedisTest()
	InitDb()
	paramKey := "usd_recv_rate"
	redisKey := ss_data.MkGlobalParamValue(paramKey)
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

	fmt.Println(redisKey, strext.ToString(tmp), err)
}
