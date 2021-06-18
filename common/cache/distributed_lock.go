package cache

import (
	"fmt"
	"time"
)

// 分布式锁-获取
// 只有当key不存在的时候才能设置成功，并设置有过期时间,避免锁一直不能释放
// @auth xiaoyanchun 2020-04-15
func GetDistributedLock(key string, value string, expire time.Duration) bool {
	boolCmd := RedisClient.SetNX(key, value, expire)
	return boolCmd.Val()
}

// 分布式锁-释放
// 只有获取到key中的值和value相等的时候才能删除成功
// 避免自己设置的值被别人释放掉
// @auth xiaoyanchun 2020-04-15
func ReleaseDistributedLock(key string, value string) bool {
	script := `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`

	vals, err := RedisClient.Eval(script, []string{key}, value).Result()

	// 没有错误，并且结构等于1
	if err == nil && fmt.Sprintf("%v", vals) == "1" {
		return true
	}

	return false
}
