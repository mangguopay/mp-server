package handler

import (
	"errors"
	"fmt"
	"time"

	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/business-bill-srv/common"
	"a.a/mp-server/business-bill-srv/dao"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
)

//检查商户的签约是否能交易
func CheckBusinessSigned(signed *dao.BusinessSceneSignedDao) (string, error) {
	if signed.SignedNo == "" {
		return ss_err.ProductUnsigned, errors.New("signedNo为空")
	}

	if signed.Status == constants.SignedStatusInvalid {
		return ss_err.SignedExpired, errors.New(fmt.Sprintf("商家(%v)签约(%v)已过期", signed.BusinessNo, signed.SignedNo))
	}

	//当前时间距离签约过期时间只剩10s就不给商户下单，目前程序设计是不会出现这种情况
	expireTime := ss_time.ParseTimeFromPostgres(signed.EndTime, global.Tz).Add(-10 * time.Second)
	timeDifference := expireTime.Sub(ss_time.Now(global.Tz)).Seconds()
	ss_log.Info("签约过期时间与当前时间差：%v", timeDifference)
	if timeDifference <= 0 {
		return ss_err.SignedExpired, errors.New(fmt.Sprintf("商家(%v)签约(%v)即将过期", signed.BusinessNo, signed.SignedNo))
	}

	if signed.Status != constants.SignedStatusPassed {
		return ss_err.ProductUnsigned, errors.New(fmt.Sprintf("商家(%v)签约(%v)不可用", signed.BusinessNo, signed.SignedNo))
	}

	return ss_err.Success, nil
}

//检查APP所属商户是否有收款权限
func CheckBusinessIncomeAuth(b *dao.BusinessDao) (string, error) {
	//商户账号
	if b.BusinessAccNo == "" {
		return ss_err.BusinessNotExist, errors.New(fmt.Sprintf("商户[%v]账号为空", b.BusinessNo))
	}

	// 商户状态不可用
	if !b.IsEnabled {
		return ss_err.BusinessNotAvailable, errors.New(fmt.Sprintf("商户[%v]不可用", b.BusinessNo))
	}

	//商户是否有收款权限
	if b.IncomeAuthorization == constants.BusinessIncomeAuthDisabled {
		return ss_err.AccountNoNotTradeForbid, errors.New(fmt.Sprintf("商户[%v]没有收款权限, incomeAuthorization:%v", b.BusinessNo, b.IncomeAuthorization))
	}

	return ss_err.Success, nil
}

//组织APP支付参数
func GetAppPayContent(order dao.BusinessBillDao) string {
	m := make(map[string]interface{})

	m["order_no"] = order.OrderNo
	//m["app_name"] = order.AppName //以前是显示应用名称，现在是商家简称
	m["app_name"] = order.SimplifyName
	m["subject"] = order.Subject
	m["amount"] = order.Amount
	m["currency_type"] = order.CurrencyType
	m["timestamp"] = fmt.Sprintf("%v", ss_time.Now(global.Tz).Unix())
	m[common.SignField] = AppPayContentMakeSign(m)

	return strext.ToJson(m)
}

//APP支付参数签名
func AppPayContentMakeSign(data map[string]interface{}) string {
	sortStr := ss_func.ParamsMapToString(data, common.SignField)
	ss_log.Info("APP支付参数签名字符串：%v", sortStr)
	return encrypt.DoSha256(sortStr, constants.AppPaySignKey)
}
