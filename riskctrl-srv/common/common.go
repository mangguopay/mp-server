package common

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
)

// 获取时间类型
func ParsingTimeType(timeType string) string {
	switch timeType {
	case "s": // 秒
		return constants.TimeType_Second
	case "m": // 分
		return constants.TimeType_Minute
	case "h": // 小时
		return constants.TimeType_Hour
	case "d": // 天
		return constants.TimeType_Day
	case "w": // 周
		return constants.TimeType_Week
	case "Mon": // 月
		return constants.TimeType_Month
	case "y": // 年
		return constants.TimeType_Year
	default:
		ss_log.Error("err=[格式化时间,参数错误----->%s]", timeType)
		return ""
	}
	return ""
}

const (
	ExeParam_LoginTimes = "loginTimes"
)
