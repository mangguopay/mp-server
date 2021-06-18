package handler

import (
	"errors"
	"strings"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/common/constants"
	billProto "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
)

type CheckAmountIsMaxMinParam struct {
	UsdMax string
	UsdMin string
	KhrMax string
	KhrMin string
}

// 校验最大金额最下金额是否超出限制-转账
func CheckAmountIsMaxMinTransfer(monType, amount string) error {
	p := CheckAmountIsMaxMinParam{
		UsdMax: "usd_transfer_single_max",
		UsdMin: "usd_transfer_single_min",
		KhrMax: "khr_transfer_single_max",
		KhrMin: "khr_transfer_single_min",
	}
	return checkAmountIsMaxMin(monType, amount, p)
}

// 校验最大金额最下金额是否超出限制-提现
func CheckAmountIsMaxMinWithdraw(monType, amount string) error {
	p := CheckAmountIsMaxMinParam{
		UsdMax: "usd_phone_single_max",
		UsdMin: "usd_phone_single_min",
		KhrMax: "khr_phone_single_max",
		KhrMin: "khr_phone_single_min",
	}
	return checkAmountIsMaxMin(monType, amount, p)
}

// 校验最大金额最下金额是否超出限制-存款
func CheckAmountIsMaxMinSave(monType, amount string) error {
	p := CheckAmountIsMaxMinParam{
		UsdMax: "usd_deposit_single_max",
		UsdMin: "usd_deposit_single_min",
		KhrMax: "khr_deposit_single_max",
		KhrMin: "khr_deposit_single_min",
	}
	return checkAmountIsMaxMin(monType, amount, p)
}

// 检验最大最小金额是否超出限制金额
func checkAmountIsMaxMin(monType, amount string, p CheckAmountIsMaxMinParam) error {
	switch monType {
	case constants.CURRENCY_USD:
		// 判断金额是否包含小数点
		if strings.Contains(amount, ".") {
			return errors.New("usd 取款金额应为整数")
		}
		// 最大限额
		usdTransferSingleMax := dao.GlobalParamDaoInstance.QeuryParamValue(p.UsdMax)
		// 最小限额
		usdTransferSingleMin := dao.GlobalParamDaoInstance.QeuryParamValue(p.UsdMin)

		ss_log.Info("usd 最大金额:%s=%s,最小金额:%s=%s\n", p.UsdMax, usdTransferSingleMax, p.UsdMin, usdTransferSingleMin)

		if strext.ToFloat64(amount) < strext.ToFloat64(usdTransferSingleMin) || strext.ToFloat64(amount) > strext.ToFloat64(usdTransferSingleMax) {
			//ss_log.Error("扫码取款,币种 美元,超出金额限制,当前金额为--->%s,最大金额限制为--->%s,最小金额限制为--->%s", amount, usdTransferSingleMax, usdTransferSingleMin)
			//reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
			return errors.New("usd 操作金额超出金额限制")
		}
	case constants.CURRENCY_KHR:
		// 最大限额
		khrTransferSingleMax := dao.GlobalParamDaoInstance.QeuryParamValue(p.KhrMax)
		// 最小限额
		khrTransferSingleMin := dao.GlobalParamDaoInstance.QeuryParamValue(p.KhrMin)
		ss_log.Info("khr 最大金额:%s=%s,最小金额:%s=%s\n", p.KhrMax, khrTransferSingleMax, p.KhrMin, khrTransferSingleMin)
		if strext.ToFloat64(amount) < strext.ToFloat64(khrTransferSingleMin) || strext.ToFloat64(amount) > strext.ToFloat64(khrTransferSingleMax) {
			//ss_log.Error("扫码取款,币种 瑞尔,超出金额限制,当前金额为--->%s,最大金额限制为--->%s,最小金额限制为--->%s", amount, khrTransferSingleMax, khrTransferSingleMin)
			//reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
			return errors.New("khr 操作金额超出金额限制")
		}
	default:
		return errors.New("币种类型不存在")
	}
	return nil
}

func CheckEnterpriseTransferToUserParam(req *billProto.EnterpriseTransferToUserRequest) (errCode string, err error) {
	if req.TransferNo == "" {
		return ss_err.ERR_PARAM, errors.New("TransferNo参数为空")
	}

	if req.TransferType == "" {
		return ss_err.ERR_PARAM, errors.New("TransferType参数为空")
	}

	return ss_err.ERR_SUCCESS, nil
}
