package util

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func MkSmsCodeName(code string) string {
	return fmt.Sprintf("%s_%s", "sms", code)
}

func MkPickTokenName(code string) string {
	return fmt.Sprintf("%s_%s", "pic_token", code)
}

func AlipaySpecialUrlEncode(str string) string {
	tmp := url.QueryEscape(str)
	for k, v := range map[string]string{
		"+": "%20",
		"*": "%2A",
		"~": "%7E",
	} {
		tmp = strings.Replace(tmp, k, v, -1)
	}

	return tmp
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
