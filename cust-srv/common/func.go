package common

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"fmt"
	"strconv"
)

// 修正各币种的金额
func NormalAmountByMoneyType(moneyType, amount string) string {
	amountRet := ""
	switch moneyType {
	case constants.CURRENCY_USD:
		amountRet = strext.ToStringNoPoint(strext.ToFloat64(amount) / 100)
		//todo 保留两位小数
		fInput, err := strconv.ParseFloat(amountRet, 64)
		if err != nil {
			ss_log.Error("string类型转float类型出错，amountRet=[%v],err=[%v]", amountRet, err)
			//提前返回，不再进行保留两位小数操作，如整数3则返回3，3.3则返回3.3。
			return amountRet
		}
		//返回两位小数的金额。
		amountRet = fmt.Sprintf("%.2f", fInput)
	case constants.CURRENCY_KHR:
		amountRet = strext.ToStringNoPoint(strext.ToFloat64(amount))
	}

	return amountRet
}
