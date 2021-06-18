package global

import (
	"a.a/mp-server/common/constants"
	"strings"
)

/**
很多地方都要用币种去查虚拟账号类型:
	商家以币种分为USD、KHR，在币种的基础上又分已结算和未结算
	用户以币种分为USD、KHR，在币种的基础上又分已激活未激活
*/

func GetBusinessVAccType(currencyType string, isSettled bool) int {
	currencyType = strings.ToUpper(currencyType)
	var businessVAccType int
	if isSettled {
		switch currencyType {
		case constants.CURRENCY_UP_USD:
			businessVAccType = constants.VaType_USD_BUSINESS_SETTLED
		case constants.CURRENCY_UP_KHR:
			businessVAccType = constants.VaType_KHR_BUSINESS_SETTLED
		}
	} else {
		switch currencyType {
		case constants.CURRENCY_UP_USD:
			businessVAccType = constants.VaType_USD_BUSINESS_UNSETTLED
		case constants.CURRENCY_UP_KHR:
			businessVAccType = constants.VaType_KHR_BUSINESS_UNSETTLED
		}
	}
	return businessVAccType
}

func GetUserVAccType(currencyType string, isActivate bool) int {
	currencyType = strings.ToUpper(currencyType)
	var userVAccType int
	if isActivate {
		switch currencyType {
		case constants.CURRENCY_UP_USD:
			userVAccType = constants.VaType_USD_DEBIT
		case constants.CURRENCY_UP_KHR:
			userVAccType = constants.VaType_KHR_DEBIT
		}
	} else {
		switch currencyType {
		case constants.CURRENCY_UP_USD:
			userVAccType = constants.VaType_FREEZE_USD_DEBIT
		case constants.CURRENCY_UP_KHR:
			userVAccType = constants.VaType_FREEZE_KHR_DEBIT
		}
	}
	return userVAccType
}

func GetPlatFormVAccType(currencyType string) int {
	currencyType = strings.ToUpper(currencyType)
	var platFormVAccType int
	switch currencyType {
	case constants.CURRENCY_UP_USD:
		platFormVAccType = constants.VaType_USD_FEES
	case constants.CURRENCY_UP_KHR:
		platFormVAccType = constants.VaType_KHR_FEES
	default:
		platFormVAccType = 0
	}
	return platFormVAccType
}
