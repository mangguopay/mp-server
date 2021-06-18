package cache

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_data"
)

//var (
//	RedisCli = new(cache.RedisU)
//)

func CheckSMS(function, phone, msg string) bool {
	// 从redis获取短信验证码
	//value, err := ss_data.GetSMSMsgFromCache1st(ss_data.GetSMSKey(function, phone), RedisCli, constants.DefPoolName)
	key := ss_data.GetSMSKey(function, phone)
	value, err := RedisClient.Get(key).Result()
	if value == "" || err != nil {
		ss_log.Error("获取短信验证码失败 key=[%s],err=[%v]", key, err)
		return false
	}
	return msg == value
}

func CheckMailCode(function, mail, msg string) bool {
	// 从redis获取邮箱验证码
	key := ss_data.GetMailKey(function, mail)
	value, err := RedisClient.Get(key).Result()
	if value == "" || err != nil {
		ss_log.Error("获取邮箱验证码失败 key=[%s],err=[%v]", key, err)
		return false
	}
	return msg == value
}
