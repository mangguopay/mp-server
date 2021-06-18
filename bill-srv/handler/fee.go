package handler

import (
	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"context"
	"errors"
)

func (b *BillHandler) QeuryFees(ctx context.Context, req *go_micro_srv_bill.QeuryRateRequest, reply *go_micro_srv_bill.QeuryRateReply) error {
	// 判断金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[查询手续费, 金额为0或者为空,传入的金额为----->%s]", req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}
	rate, fees, err := doFees(req.Type, req.Amount)
	if err != nil {
		ss_log.Error("计算手续费失败,err:--->%s", err.Error())
		reply.ResultCode = ss_err.ERR_RATE_FAILD
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Data = &go_micro_srv_bill.RateData{
		Fees: fees,
		Rate: rate,
	}
	return nil
}

func doFees(feeType int32, amount string) (rate, feeAmount string, retErr error) {
	// 计算费率
	if amount == "" || amount == "0" {
		retErr = errors.New("计算手续费金额为0")
		return
	}

	var paramKey string
	var minDefaultFees string
	switch feeType {
	//case 1:
	//	paramKey = constants.SCORE_SETTING_BUY_RATE
	//	return nil
	case constants.Fees_Type_USD_TRANSFER_RATE: // 判断转出权限
		paramKey = constants.USD_TRANSFER_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_USD_MIN_TRANSFER_FEE)
	case constants.Fees_Type_USD_RECV_RATE: // 判断转入权限
		paramKey = constants.USD_RECV_RATE
		// 收款属于转账类
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_USD_MIN_TRANSFER_FEE)
	case constants.Fees_Type_USD_DEPOSIT_RATE: // 判断充值权限
		paramKey = constants.USD_DEPOSIT_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_USD_MIN_DEPOSIT_FEE)
	case constants.Fees_Type_USD_PHONE_WITHDRAW_RATE: // 判断提现权限 usd手机号取款
		paramKey = constants.USD_PHONE_WITHDRAW_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_USD_PHONE_MIN_WITHDRAW_FEE)
	case constants.Fees_Type_USD_FACE_WITHDRAW_RATE: // 判断提现权限 usd 面对面取款
		paramKey = constants.USD_FACE_WITHDRAW_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_USD_FACE_MIN_WITHDRAW_FEE)
	case constants.Fees_Type_KHR_TRANSFER_RATE: // 判断转出权限
		paramKey = constants.KHR_TRANSFER_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_KHR_MIN_TRANSFER_FEE)
	case constants.Fees_Type_KHR_RECV_RATE: // 判断转入权限
		paramKey = constants.KHR_RECV_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_KHR_MIN_TRANSFER_FEE)
	case constants.Fees_Type_KHR_DEPOSIT_RATE: // 判断充值权限
		paramKey = constants.KHR_DEPOSIT_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_KHR_MIN_DEPOSIT_FEE)
	case constants.Fees_Type_KHR_PHONE_WITHDRAW_RATE: // 判断提现权限 khr手机号提现费率
		paramKey = constants.KHR_PHONE_WITHDRAW_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_KHR_PHONE_MIN_WITHDRAW_FEE)
	case constants.Fees_Type_KHR_FACE_WITHDRAW_RATE: // 判断提现权限 khr手机号提现费率
		paramKey = constants.KHR_FACE_WITHDRAW_RATE
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_KHR_FACE_MIN_WITHDRAW_FEE)
	case constants.Fees_Type_Usd_To_Khr_Count_Fee: // usd-->khr 单笔手续费
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_USD_TO_KHR_FEE)
		rate = "0"
		feeAmount = minDefaultFees
		return
	case constants.Fees_Type_Khr_To_Usd_Count_Fee: //khr-->usd 单笔手续费
		_, minDefaultFees, _ = cache.ApiDaoInstance.GetGlobalParam(constants.KEY_KHR_TO_USD_FEE)
		rate = "0"
		feeAmount = minDefaultFees
		return
	default:
		retErr = errors.New("输入的取款类型不正确")
		return
	}

	_, paramValue, _ := cache.ApiDaoInstance.GetGlobalParam(paramKey)
	if paramValue == "" {
		retErr = errors.New("查询globalParam中paramValue为空")
		return
	}
	if paramValue == "0" { // 没有手续费
		rate = paramValue
		feeAmount = "0"
		return
	}

	fees := ss_count.CountFees(amount, paramValue, minDefaultFees)
	// 取整数
	feeAmount = ss_big.SsBigInst.ToRound(fees, 0, ss_big.RoundingMode_HALF_EVEN).String()
	rate = paramValue
	return
}
