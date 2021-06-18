package handler

import (
	"a.a/cu/ss_big"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_count"
	"errors"
	"fmt"
)

type InternalCallHandler struct {
}

var InternalCallHandlerInst InternalCallHandler

/**
查询账号的虚拟账号，如果不存在则初始化一个新的虚拟账号
*/
func (InternalCallHandler) ConfirmExistVAccount(accountNo, balanceType string, vaType int32) (vAccountNo string) {
	vAccountNo = dao.VaccountDaoInst.GetVaccountNo(accountNo, vaType)
	if vAccountNo == "" {
		vAccountNo = dao.VaccountDaoInst.InitVaccountNo(accountNo, balanceType, vaType)
	}
	return vAccountNo
}

/**
美金兑换成瑞尔手续费计算
*/
func (InternalCallHandler) ExchangeUsdToKhr(amount, rate string) string {
	a := ss_count.Multiply(amount, rate)
	result := ss_count.Div(a.String(), "100") // 结果需要处于100然后取整
	// 取整数
	return ss_big.SsBigInst.ToRound(result, 0, ss_big.RoundingMode_HALF_EVEN).String()
}

/**
瑞尔兑换成美金手续费计算
*/
func (InternalCallHandler) ExchangeKhrToUsd(amount, rate string) (string, error) {
	reqAmount := ss_count.Multiply(amount, "100") // 先把khr *100后再比较,因为 usd 的单位是分

	// 判断是否能足够兑换1美元
	_, exchangeRateAmount, _ := cache.ApiDaoInstance.GetGlobalParam(constants.Exchange_Khr_To_Usd)
	reqAmountF := strext.StringToFloat64(amount)
	exchangeRateF := strext.StringToFloat64(exchangeRateAmount) / 100
	if reqAmountF < exchangeRateF { // 最低兑换金额为0.01USD
		return "", errors.New(fmt.Sprintf("khr_to_usd, khr 的余额不足,khr的余额为---> %v,最低需要为---> %v", reqAmountF, exchangeRateF))
	}

	fromString := ss_count.Div(reqAmount.String(), rate)
	// 美元  四舍六入五成双
	return ss_big.SsBigInst.ToRound(fromString, 0, ss_big.RoundingMode_HALF_EVEN).String(), nil
}
