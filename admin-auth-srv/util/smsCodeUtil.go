package util

import (
	"fmt"
	"strings"
	"time"

	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
)

func MkSmsCodeName(code string) string {
	return fmt.Sprintf("%s_%s", "sms", code)
}

func DoChkSmsCode(code, phone string) bool {
	//---------------第一版本的redis-----------------
	//ret, err := cache.RedisCli.Get("a", MkSmsCodeName(code))

	ret, err := cache.RedisClient.Get(MkSmsCodeName(code)).Result()
	if ret == "" || err != nil {
		return false
	}

	l := strings.Split(ret, ",")
	nowStr := time.Now().Unix()
	if nowStr > strext.ToInt64(l[0]) || l[1] != phone {
		return false
	}
	return true
}
