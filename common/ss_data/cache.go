package ss_data

import (
	"a.a/mp-server/common/constants"
	"errors"
	"fmt"
	"strings"
	"time"
)

// 全局配置
func MkGlobalParamValue(paramKey string) string {
	return fmt.Sprintf("%v_%v", constants.PreGlobalParam, paramKey)
}

func GetSMSKey(business, mobile string) string {
	return fmt.Sprintf("%s_%s_%s", constants.PreSms, business, mobile)
}

func MkRelaApi(reqUrl, accNo string) string {
	reqUrl = strings.Replace(reqUrl, "/", "_", -1)
	return fmt.Sprintf("%v_%v_%v", constants.PreRelaApi, reqUrl, accNo)
}

func GetMailKey(business, mobile string) string {
	return fmt.Sprintf("%s_%s_%s", constants.PreMail, business, mobile)
}

func FormatDateTostring(template string, date interface{}) (string, error) {
	switch date.(type) {
	case string:
		tempTime, err := time.Parse(template, date.(string))
		if err != nil {
			return "", err
		}
		return tempTime.String(), nil
	case time.Time:
		t := date.(time.Time)
		return t.Format(template), nil
	}

	return "", errors.New("Parameter types do not match")
}
