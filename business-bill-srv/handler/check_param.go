package handler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"a.a/cu/ss_time"
	"a.a/mp-server/common/global"

	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
)

func CheckPrepayReq(req *businessBillProto.PrepayRequest) (string, error) {
	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	if req.AppId == "" {
		return ss_err.AppIdIsEmpty, errors.New("AppId参数为空")
	}

	if req.TradeType == "" {
		return ss_err.TradeTypeIsEmpty, errors.New("TradeType字段为空")
	}

	if req.OutOrderNo == "" {
		return ss_err.OutOrderNoIsEmpty, errors.New("OutOrderNo参数为空")
	}

	if req.Amount == "" {
		return ss_err.AmountIsEmpty, errors.New("Amount参数为空")
	}

	if req.CurrencyType == "" {
		return ss_err.CurrencyTypeIsEmpty, errors.New("CurrencyType参数为空")
	}

	if !util.InSlice(req.TradeType, []string{constants.TradeTypeModernpayAPP, constants.TradeTypeModernpayFaceToFace}) {
		return ss_err.TradeTypeValueIsIllegality, errors.New("TradeType字段值错误")
	}

	if !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_UP_USD, constants.CURRENCY_UP_KHR}) {
		return ss_err.CurrencyTypeValueIsIllegality, errors.New("币种错误")
	}

	if strext.ToInt64(req.Amount) < 1 {
		return ss_err.AmountNotLessThanOne, errors.New("金额不能小于1")
	}

	if req.Subject == "" {
		switch req.Lang {
		case constants.LangZhCN:
			req.Subject = "其它"
		case constants.LangEnUS:
			req.Subject = "other"
		case constants.LangKmKH:
			req.Subject = "ផ្សេងៗ"
		default:
			req.Subject = "other"
		}
	}

	expireTime := ss_time.Now(global.Tz).Add(constants.BusinessOrderExpireTime * time.Minute).Unix()
	if req.TimeExpire == "" {
		req.TimeExpire = strext.ToString(expireTime)
	} else {
		timestamp, err := strconv.ParseInt(req.TimeExpire, 10, 64)
		if err != nil {
			return ss_err.TimeFormatErr, errors.New("TimeExpire参数不是一个时间戳")
		}
		currentTime := ss_time.Now(global.Tz).Unix()
		if timestamp < currentTime || timestamp > expireTime {
			req.TimeExpire = strext.ToString(expireTime)
		}
	}

	return ss_err.Success, nil
}

func CheckQrCodeFixedPrePay(req *businessBillProto.QrCodeFixedPrePayRequest) (string, error) {
	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	if req.QrCodeId == "" {
		return ss_err.QrCodeIdIsEmpty, errors.New("QrCodeId参数为空")
	}
	if req.AccountNo == "" {
		return ss_err.AccountNoIsEmpty, errors.New("AccountNo参数为空")
	}
	if req.Amount == "" {
		return ss_err.AmountIsEmpty, errors.New("Amount参数为空")
	}
	if req.CurrencyType == "" {
		return ss_err.CurrencyTypeIsEmpty, errors.New("CurrencyType参数为空")
	}
	if req.AccountType != constants.AccountType_USER {
		return ss_err.AccountTypeNotUser, errors.New(fmt.Sprintf("AccountType=%v, 非用户(4)", req.AccountType))
	}

	if strext.ToInt64(req.Amount) < 1 {
		return ss_err.AmountNotLessThanOne, errors.New("金额不能小于1")
	}

	if !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_UP_USD, constants.CURRENCY_UP_KHR}) {
		return ss_err.CurrencyTypeValueIsIllegality, errors.New("币种错误")
	}
	return ss_err.Success, nil
}

func CheckQrCodeFixedPayReq(req *businessBillProto.QrCodeFixedPayRequest) (string, error) {
	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	if req.OrderNo == "" {
		return ss_err.OrderNoIsEmpty, errors.New("OrderNo参数为空")
	}
	if req.AccountNo == "" {
		return ss_err.AccountNoIsEmpty, errors.New("AccountNo参数为空")
	}
	if req.AccountType != constants.AccountType_USER {
		return ss_err.AccountTypeNotUser, errors.New(fmt.Sprintf("AccountType=%v, 非用户(4)", req.AccountType))
	}
	if req.PaymentPassword == "" {
		return ss_err.PaymentPwdIsEmpty, errors.New("PaymentPassword参数为空")
	}
	if req.NonStr == "" {
		return ss_err.NonStrIsEmpty, errors.New("NonStr参数为空")
	}

	if req.PaymentMethod == constants.PayMethodBankCard {
		if req.BankCardNo == "" {
			return ss_err.BankCardNumberIsEmpty, errors.New("BankCardNumber参数为空")
		}
	}

	return ss_err.Success, nil
}

func CheckOrderPayReq(req *businessBillProto.OrderPayRequest) (string, error) {
	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	if req.OrderNo == "" {
		return ss_err.OrderNoIsEmpty, errors.New("OrderNo参数为空")
	}
	if req.AccountNo == "" {
		return ss_err.AccountNoIsEmpty, errors.New("AccountNo参数为空")
	}
	if req.AccountType != constants.AccountType_USER {
		return ss_err.AccountTypeNotUser, errors.New(fmt.Sprintf("AccountType=%v, 非用户(4)", req.AccountType))
	}
	if req.PaymentPwd == "" {
		return ss_err.PaymentPwdIsEmpty, errors.New("PaymentPassword参数为空")
	}
	if req.NonStr == "" {
		return ss_err.NonStrIsEmpty, errors.New("NonStr参数为空")
	}

	if req.PaymentMethod == constants.PayMethodBankCard {
		if req.BankCardNo == "" {
			return ss_err.BankCardNumberIsEmpty, errors.New("BankCardNumber参数为空")
		}
	}

	return ss_err.Success, nil
}

func CheckAppPayReq(req *businessBillProto.AppPayRequest) (string, error) {
	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	if req.AccountNo == "" {
		return ss_err.AccountNoIsEmpty, errors.New("AccountNo参数为空")
	}
	if req.PaymentPwd == "" {
		return ss_err.PaymentPwdIsEmpty, errors.New("PaymentPassword参数为空")
	}
	if req.NonStr == "" {
		return ss_err.NonStrIsEmpty, errors.New("NonStr参数为空")
	}
	if req.AppPayContent == "" {
		return ss_err.AppPayContentIsEmpty, errors.New("AppPayContent参数为空")
	}

	if req.PaymentMethod == constants.PayMethodBankCard {
		if req.BankCardNo == "" {
			return ss_err.BankCardNumberIsEmpty, errors.New("BankCardNumber参数为空")
		}
	}

	return ss_err.Success, nil
}

func CheckPayReq(req PayRequest) (string, error) {
	if req.OrderNo == "" {
		return ss_err.OrderNoIsEmpty, errors.New("OrderNo参数为空")
	}
	if req.AccountNo == "" {
		return ss_err.AccountNoIsEmpty, errors.New("AccountNo参数为空")
	}
	if req.PaymentPassword == "" {
		return ss_err.PaymentPwdIsEmpty, errors.New("PaymentPassword参数为空")
	}
	if req.NonStr == "" {
		return ss_err.NonStrIsEmpty, errors.New("NonStr参数为空")
	}
	if req.AccountType != constants.AccountType_USER {
		return ss_err.AccountTypeNotUser, errors.New(fmt.Sprintf("AccountType=%v, 非用户(4)", req.AccountType))
	}

	if req.PaymentMethod == constants.PayMethodBankCard {
		if req.BankCardNo == "" {
			return ss_err.BankCardNumberIsEmpty, errors.New("BankCardNumber参数为空")
		}
	}

	return ss_err.Success, nil
}

func CheckEnterpriseTransferReq(req *businessBillProto.EnterpriseTransferRequest) (string, error) {
	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	if req.AppId == "" {
		return ss_err.AppIdIsEmpty, errors.New("AppId参数为空")
	}

	if req.OutTransferNo == "" {
		return ss_err.OutOrderNoIsEmpty, errors.New("OutTransferNo参数为空")
	}

	if req.Amount == "" {
		return ss_err.AmountIsEmpty, errors.New("Amount参数为空")
	}

	if req.CurrencyType == "" {
		return ss_err.CurrencyTypeIsEmpty, errors.New("CurrencyType参数为空")
	}

	if !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_UP_USD, constants.CURRENCY_UP_KHR}) {
		return ss_err.CurrencyTypeValueIsIllegality, errors.New("币种错误")
	}

	if strext.ToInt64(req.Amount) < 1 {
		return ss_err.AmountNotLessThanOne, errors.New("金额不能小于1")
	}

	if (req.PayeePhone != "" && req.PayeeEmail != "") || (req.PayeePhone == "" && req.PayeeEmail == "") {
		errStr := fmt.Sprintf("PayeePhone[%v]参数和PayeeEmail[%v]参数只能二选一", req.PayeePhone, req.PayeeEmail)
		return ss_err.ParamErr, errors.New(errStr)
	}

	if req.PayeePhone != "" {
		if req.PayeeCountryCode == "" {
			return ss_err.CountryCodeIsEmpty, errors.New("PayeeCountryCode参数为空")
		}
	}

	return ss_err.Success, nil
}

func CheckApiPayRefundReq(req *businessBillProto.ApiPayRefundRequest) (string, error) {
	if req.AppId == "" {
		return ss_err.AppIdIsEmpty, errors.New("AppId参数为空")
	}
	if req.OrderNo == "" && req.OutOrderNo == "" {
		return ss_err.OrderNoIsEmpty, errors.New("OrderNo参数为空")
	}
	if req.RefundAmount == "" {
		return ss_err.RefundAmountIsEmpty, errors.New("RefundAmount参数为空")
	}

	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	return ss_err.Success, nil
}

func CheckPersonalBusinessPrepayReq(req *businessBillProto.PersonalBusinessPrepayRequest) (string, error) {
	if req.AccountNo == "" {
		return ss_err.ParamErr, errors.New("AccountNo参数为空")
	}
	if req.Amount == "" {
		return ss_err.AmountIsEmpty, errors.New("amount参数为空")
	}
	if req.CurrencyType == "" {
		return ss_err.RefundAmountIsEmpty, errors.New("CurrencyType参数为空")
	}

	if !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_UP_USD, constants.CURRENCY_UP_KHR}) {
		return ss_err.CurrencyTypeValueIsIllegality, errors.New("币种错误")
	}

	if strext.ToInt64(req.Amount) < 1 {
		return ss_err.AmountNotLessThanOne, errors.New("金额不能小于1")
	}

	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	return ss_err.Success, nil
}

func CheckPersonalBusinessCodeFixedPrePay(req *businessBillProto.PersonalBusinessCodeFixedPrePayRequest) (string, error) {
	// 没有传入语言默认是英语
	if req.Lang == "" || !util.InSlice(req.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		req.Lang = constants.LangEnUS
	}

	if req.QrCodeId == "" {
		return ss_err.QrCodeIdIsEmpty, errors.New("QrCodeId参数为空")
	}
	if req.AccountNo == "" {
		return ss_err.AccountNoIsEmpty, errors.New("AccountNo参数为空")
	}
	if req.Amount == "" {
		return ss_err.AmountIsEmpty, errors.New("Amount参数为空")
	}
	if req.CurrencyType == "" {
		return ss_err.CurrencyTypeIsEmpty, errors.New("CurrencyType参数为空")
	}
	if req.AccountType != constants.AccountType_USER {
		return ss_err.AccountTypeNotUser, errors.New(fmt.Sprintf("AccountType=%v, 非用户(4)", req.AccountType))
	}

	if strext.ToInt64(req.Amount) < 1 {
		return ss_err.AmountNotLessThanOne, errors.New("金额不能小于1")
	}

	if !util.InSlice(req.CurrencyType, []string{constants.CURRENCY_UP_USD, constants.CURRENCY_UP_KHR}) {
		return ss_err.CurrencyTypeValueIsIllegality, errors.New("币种错误")
	}
	return ss_err.Success, nil
}
