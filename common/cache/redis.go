package cache

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"errors"
	"fmt"
	goredis "github.com/go-redis/redis/v7"
	"net"
)

var RedisClient *goredis.Client

var RedisNil = errors.New("redis: nil")

const (
	PrePwdErrCountKey        = "err_pwd_"
	PrePaymentPwdErrCountKey = "err_payment_pwd_"
)

// 初始化redis实例
func InitRedis(host string, port string, Password string, db int) error {
	client, err := NewRedisClient(host, port, Password, db)
	if err != nil {
		return err
	}

	RedisClient = client

	return nil
}

// 获取redis客户端实例
func NewRedisClient(Addr string, port string, Password string, db int) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     net.JoinHostPort(Addr, port),
		Password: Password, // no password set
		DB:       db,       // use default DB
	})

	_, err := client.Ping().Result()
	return client, err
}

// 设置短信验证码进Redis
func SetSMSMsgToCache1stV2(key, value string, redis *goredis.Client) error {
	return redis.Set(key, value, constants.SmsKeySecV2).Err()
}

// 从Redis中获取验证码
func GetSMSMsgFromCache1stV2(key string, redis *goredis.Client) (string, error) {
	value, err := redis.Get(key).Result()
	if err == RedisNil {
		return "", nil
	}

	return value, err
}

func GetDataFromCache1st(key string, redisKey string, redis *goredis.Client,
	getDataFunc func(key string) (string, error)) (value string, err error) {
	ret, err := redis.Get(redisKey).Result()
	if ret == "" || err != nil {
		value, err := getDataFunc(key)
		if err != nil {
			return value, err
		}
		redis.Set(redisKey, value, constants.CacheKeySecV2)
		return value, nil
	}
	return strext.ToString(ret), nil
}

func GetPwdErrCountKey(pre, k string) string {
	return fmt.Sprintf("%s%s", pre, k)
}

//支付密码连续出错次数key
func GetPayPwdErrCountKey(pre, accType, idenNo string) string {
	return fmt.Sprintf("%s%s_%s", pre, accType, idenNo)
}
func GetPwdErrCountKey2(accType, pre, k string) string {
	return fmt.Sprintf("%s_%s_%s", accType, pre, k)
}

func GetDataFromCache1stI(keys []string, redisKey string,
	getDataFunc func(keysI []string) (interface{}, error)) (value interface{}, err error) {
	result, err := RedisClient.Get(redisKey).Result()
	if result == "" || err != nil {
		value, err := getDataFunc(keys)
		if err != nil {
			return value, err
		}
		RedisClient.Set(redisKey, strext.ToJsonNotChange(value), constants.CacheKeySecV2)
		return value, nil
	}
	return result, nil
}

func GetDataFromCache1stL(keys []string, redisKey string,
	getDataFunc func(keysI []string) (string, error)) (value string, err error) {
	result, err := RedisClient.Get(redisKey).Result()
	//ret, err := redisU.Get(poolname, redisKey)
	if result == "" || err != nil {
		value, err := getDataFunc(keys)
		if err != nil {
			return value, err
		}
		RedisClient.Set(redisKey, value, constants.CacheKeySecV2)
		return value, nil
	}
	return result, nil
}

func MkChannelSupplier(channelNo string) string {
	return fmt.Sprintf("%v_%v", constants.PRE_CHANNEL_SUPPLIER, channelNo)
}

//
func MkMercApiMode(mercNo string) string {
	return fmt.Sprintf("%v_%v", constants.PreMercApiMode, mercNo)
}

func MkChannelParam(channelNo string) string {
	return fmt.Sprintf("%v_%v", constants.PRE_CHANNEL_PARAM, channelNo)
}
