package util

import (
	"fmt"
)

func MkSmsCodeName(code string) string {
	return fmt.Sprintf("%s_%s", "sms", code)
}
