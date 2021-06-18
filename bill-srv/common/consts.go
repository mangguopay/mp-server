package common

import (
	"fmt"
	"time"

	"a.a/mp-server/common/constants"
)

const (
	//扫一扫取款超时
	SweepWithdrawalExpireTime = time.Second * 300
)

func GetExpEey(key string) string {
	return fmt.Sprintf("%s%s", constants.EXP_GEN_CODE_KEY, key)
}
