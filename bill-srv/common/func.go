package common

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 通过币种获取虚拟账号类型（用户）
func VirtualAccountTypeByMoneyType(moneyType string, isActived string) (int, error) {
	moneyType = strings.ToLower(moneyType)
	var vaType int
	switch moneyType {
	case constants.CURRENCY_USD:
		if isActived == "0" || isActived == "" {
			vaType = constants.VaType_FREEZE_USD_DEBIT
		} else {
			vaType = constants.VaType_USD_DEBIT
		}
	case constants.CURRENCY_KHR:
		if isActived == "0" || isActived == "" {
			vaType = constants.VaType_FREEZE_KHR_DEBIT
		} else {
			vaType = constants.VaType_KHR_DEBIT
		}
	default:
		return vaType, ss_err.ErrWrongeCurrecyType
	}

	return vaType, nil
}

// scene (1存款,2提现)
// 通过币种获取手续费类型
func FeesTypeByMoneyType(scene int, moneyType string) (int32, error) {
	feesType := int32(0)

	if scene == constants.Scene_Save { // 存款
		switch moneyType {
		case constants.CURRENCY_USD:
			feesType = constants.Fees_Type_USD_DEPOSIT_RATE
		case constants.CURRENCY_KHR:
			feesType = constants.Fees_Type_KHR_DEPOSIT_RATE
		}
	} else if scene == constants.Scene_Withdraw { // 提现
		switch moneyType {
		case constants.CURRENCY_USD:
			feesType = constants.Fees_Type_USD_PHONE_WITHDRAW_RATE
		case constants.CURRENCY_KHR:
			feesType = constants.Fees_Type_KHR_PHONE_WITHDRAW_RATE
		}
	} else if scene == constants.Scene_Transfer { // 转账
		switch moneyType {
		case constants.CURRENCY_USD:
			feesType = constants.Fees_Type_USD_TRANSFER_RATE
		case constants.CURRENCY_KHR:
			feesType = constants.Fees_Type_KHR_TRANSFER_RATE
		}
	}

	if feesType == 0 {
		return feesType, errors.New("币种类型错误或场景不存在")
	}

	return feesType, nil
}

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

// 将 2020-04-10T00:00:00Z 提取出 2020-04-10
func GetPostgresDate(dateStr string) string {
	if dateStr == "" {
		return dateStr
	}

	pos := strings.Index(dateStr, "T")
	if pos < 1 { // 没有找到“T”字符，或者就只有一个“T”字符
		return dateStr
	}

	return dateStr[0:pos]
}
