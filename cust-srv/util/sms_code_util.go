package util

import (
	"fmt"
	"net/url"
	"strings"
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
