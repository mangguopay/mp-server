package service

import (
	"a.a/cu/db"
	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/business-settle-srv/dao"
	"a.a/mp-server/common/constants"
	businessSettleProto "a.a/mp-server/common/proto/business-settle"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const DB_NO_ROWS_MSG = "sql: no rows in result set"

// 入金手续费分成
func BillService(ctx context.Context, req *businessSettleProto.BusinessSettleFeesRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)
	// 确定手续费没有统计
	if errStr := dao.BillDaoInst.ConfirmIsNoCount(req.OrderNo); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[BillService 该手续费已统计,订单号为----->%s]", req.OrderNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" {
		ss_log.Error("err=[手续费为0 ,req.Fees为----->%s ,订单号为--------->%s]", req.Fees, req.OrderNo)
		return ss_err.ERR_PARAM
	}

	//// 根据订单id查询 rate 和upperRateNo
	_, agencyNo, merchantNo, channelRate, channelNo, rate, countFee := dao.BillDaoInst.QueryRateInfo(tx, req.OrderNo)
	if merchantNo == "" {
		ss_log.Error("err=[根据订单id查询 rate 和 upperRateNo 失败,订单号为--------->%s]", req.OrderNo)
		return ss_err.ERR_PARAM
	}
	// 对手续费进行清分
	if errStr := doAgencyFees1(tx, req.OrderNo, agencyNo, channelNo, merchantNo, req.Amount, req.Fees, rate, channelRate, req.Fees, agencyNo, countFee); errStr != ss_err.ERR_SUCCESS {
		// 修改清分失败.is_count,is_wallet为3
		if errStr := dao.BillDaoInst.UpdateIsCountFromLogNo(req.OrderNo, constants.FEES_FAILE_COUNT, constants.FEES_FAILE_COUNT); errStr != ss_err.ERR_SUCCESS {
			return errStr
		}
		return errStr
	}

	// 修改手续费已经统计
	if errStr := dao.BillDaoInst.UpdateIsCountFromLogNoTx(tx, req.OrderNo, constants.FEES_IS_COUNT, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

/**
orderNo:订单号
upperAgencyNo:代理id
channelNo:通道id
amount:交易金额
totalFees:手续费总额
rate:与上级代理相减前的费率
*/
//func doAgencyFees1(tx *sql.Tx, orderNo, upperAgencyNo, channelNo, merchantNo, amount, fees, rate, channelRate, totalFees, merchantAgencyNo, countFee string) string {
//	if upperAgencyNo == constants.INIT_UUID { // 目前最上级的代理了,上面就是平台了
//		// 计算上游的手续费
//		channelDeci := ss_count.CountFees(amount, channelRate)
//		upperFees := ss_big.SsBigInst.ToRound(channelDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
//		// 处理倒扣情况
//		if strext.ToInt(countFee) > 0 {
//			upperFees = ss_count.Add(upperFees, countFee)
//		}
//		plantFormFees := ss_count.Sub(fees, upperFees).String()
//		var op string
//		var decription string
//		if strext.ToInt(plantFormFees) < 0 { // 倒扣
//			op = "-"
//			decription = fmt.Sprintf("%s", "入金操作,因手续费不够给上游,导致负扣手续费")
//			// 关闭该通道
//			if errStr := dao.ChannelDaoInst.ForbidChannelFromNo(tx, channelNo, constants.CHANNEL_STATUS_FORBID); errStr != ss_err.ERR_SUCCESS {
//				return errStr
//			}
//		} else {
//			op = "+"
//			decription = fmt.Sprintf("%s", "入金操作,分的手续费")
//		}
//		// 确定平台账号是否存在
//		logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, upperAgencyNo, constants.RoleType_Admin)
//		if logNo == "" {
//			ss_log.Error("获取钱包id失败,accNo为----->%s", upperAgencyNo)
//			return ss_err.ERR_PARAM
//		}
//		if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, plantFormFees, logNo, op); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//		// 记录平台手续费分成钱包明细
//		if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, upperAgencyNo, beforeAmount, plantFormFees, "1",
//			upperAgencyNo, constants.FEES_TYPE_BILL, orderNo, logNo, decription, channelNo, merchantNo, op); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//
//		// 商户变动金额
//		merchantChangeAmount := ss_big.SsBigInst.ToRound(ss_count.Sub(amount, totalFees), 0, ss_big.RoundingMode_HALF_EVEN).String()
//
//		// 记录商户钱包
//		merchantLogNo, merchantBeforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, merchantNo, constants.RoleType_Merc)
//		if merchantLogNo == "" {
//			ss_log.Error("获取商户钱包id失败,accNo为----->%s", merchantNo)
//			return ss_err.ERR_PARAM
//		}
//		if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, merchantChangeAmount, merchantLogNo, "+"); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//		// 记录商户钱包明细
//		merChantDecription := "商户入金操作"
//		if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, merchantNo, merchantBeforeAmount, merchantChangeAmount, "1",
//			merchantAgencyNo, constants.FEES_TYPE_BILL, orderNo, merchantLogNo, merChantDecription, channelNo, merchantNo, "+"); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//
//		return ss_err.ERR_SUCCESS
//	}
//	// 查找代理的费率
//	upperRate := dao.RateDaoInst.QueryRateFromAccNo(tx, upperAgencyNo, constants.RoleType_Agency)
//	if upperRate == "" {
//		return ss_err.ERR_PARAM
//	}
//	// 计算费率差
//	upperRateDeci := ss_count.Sub(rate, upperRate)
//	upperFeesDeci := ss_count.CountFees(amount, upperRateDeci.String())
//	// 取整
//	upperFees := ss_big.SsBigInst.ToRound(upperFeesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
//
//	// 记录日志进钱包表
//	logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, upperAgencyNo, constants.RoleType_Agency)
//	if logNo == "" {
//		ss_log.Error("获取钱包id失败,accNo为----->%s", upperAgencyNo)
//		return ss_err.ERR_PARAM
//	}
//	if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, upperFees, logNo, "+"); errStr != ss_err.ERR_SUCCESS {
//		return errStr
//	}
//
//	// 根据代理找到上级代理accNo
//	nextUpperNo := dao.AgencyDaoInst.QueryUpper(tx, upperAgencyNo)
//	if nextUpperNo == "" {
//		ss_log.Error("查询当前代理的上级代理失败,当前代理accNo为----->%s", upperAgencyNo)
//		return ss_err.ERR_PARAM
//	}
//	// 记录平台手续费分成钱包明细
//	decription := "入金操作,分得手续费"
//	if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, upperAgencyNo, beforeAmount, upperFees, "1",
//		nextUpperNo, constants.FEES_TYPE_BILL, orderNo, logNo, decription, channelNo, merchantNo, "+"); errStr != ss_err.ERR_SUCCESS {
//		return errStr
//	}
//	newFees := ss_count.Sub(fees, upperFees).String()
//
//	if errStr := doAgencyFees1(tx, orderNo, nextUpperNo, channelNo, merchantNo, amount, newFees, upperRate,
//		channelRate, totalFees, merchantAgencyNo, countFee); errStr != ss_err.ERR_SUCCESS {
//		return errStr
//	}
//
//	return ss_err.ERR_SUCCESS
//}

/**
orderNo:订单号
upperAgencyNo:代理id
channelNo:通道id
amount:交易金额
totalFees:手续费总额
rate:与上级代理相减前的费率
*/
func doAgencyFees1(tx *sql.Tx, orderNo, upperAgencyNo, channelNo, merchantNo, amount, fees, rate, channelRate, totalFees, merchantAgencyNo, countFee string) string {
	if upperAgencyNo == constants.INIT_UUID { // 目前最上级的代理了,上面就是平台了
		// 计算上游的手续费
		channelDeci := ss_count.CountFees(amount, channelRate, "0")
		upperFees := ss_big.SsBigInst.ToRound(channelDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
		// 处理倒扣情况
		if strext.ToInt(countFee) > 0 {
			upperFees = ss_count.Add(upperFees, countFee)
		}
		plantFormFees := ss_count.Sub(fees, upperFees).String()
		var op string
		var dbOp string
		var decription string
		if strext.ToInt(plantFormFees) < 0 { // 倒扣
			op = "-"
			decription = fmt.Sprintf("%s", "入金操作,因手续费不够给上游,导致负扣手续费")
			dbOp = "2"
			// 关闭该通道
			if errStr := dao.ChannelDaoInst.ForbidChannelFromNo(tx, channelNo, constants.CHANNEL_STATUS_FORBID); errStr != ss_err.ERR_SUCCESS {
				return errStr
			}
		} else {
			op = "+"
			dbOp = "1"
			decription = fmt.Sprintf("%s", "入金操作,分的手续费")
		}
		// 确定平台账号是否存在
		logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, upperAgencyNo, constants.RoleType_Admin)
		if logNo == "" {
			ss_log.Error("获取钱包id失败,accNo为----->%s", upperAgencyNo)
			return ss_err.ERR_PARAM
		}
		ss_log.Info("平台 %s 清分-------------->%s", upperAgencyNo, plantFormFees)
		if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, plantFormFees, logNo, op); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("errStr=[%v]", errStr)
			return errStr
		}
		// 记录平台手续费分成钱包明细
		if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, upperAgencyNo, beforeAmount, plantFormFees, dbOp,
			upperAgencyNo, constants.FEES_TYPE_BILL, orderNo, logNo, decription, channelNo, merchantNo, op); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("errStr=[%v]", errStr)
			return errStr
		}

		// 商户变动金额
		merchantChangeAmount := ss_big.SsBigInst.ToRound(ss_count.Sub(amount, totalFees), 0, ss_big.RoundingMode_HALF_EVEN).String()
		ss_log.Info("商户 %s,清分------------->%s", merchantNo, merchantChangeAmount)
		// 记录商户钱包
		merchantLogNo, merchantBeforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, merchantNo, constants.RoleType_Merc)
		if merchantLogNo == "" {
			ss_log.Error("获取商户钱包id失败,accNo为----->%s", merchantNo)
			return ss_err.ERR_PARAM
		}
		if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, merchantChangeAmount, merchantLogNo, "+"); errStr != ss_err.ERR_SUCCESS {
			return errStr
		}
		// 记录商户钱包明细
		merChantDecription := "商户入金操作"
		if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, merchantNo, merchantBeforeAmount, merchantChangeAmount, "1",
			merchantAgencyNo, constants.FEES_TYPE_BILL, orderNo, merchantLogNo, merChantDecription, channelNo, merchantNo, "+"); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("errStr=[%v]", errStr)
			return errStr
		}

		// 记录成本
		if errStr := dao.CostDaoInst.InsertCost(tx, orderNo, constants.FEES_TYPE_BILL, upperFees, channelNo); errStr != ss_err.ERR_SUCCESS {
			return errStr
		}
		// todo 修改bill_result表,submit_amount_sum,real_amount_sum,enter_fees,enter_amount_sum,
		if errStr := dao.BillResultDaoInst.ModifyBillEnterResult(tx, amount, amount, totalFees, merchantChangeAmount, merchantNo, constants.RoleType_Merc); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("修改bill_result表 失败,订单号为-----> %s", orderNo)
			return errStr
		}
		// todo 修改统计表,成功交易总金额(amount),平台利润,商户收入总金额
		if errStr := dao.CumulativeCountDaoInst.ModifyCumulative1(tx, amount, plantFormFees, merchantChangeAmount); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("修改统计利润表失败,订单号为-----> %s", orderNo)
			return errStr
		}

		return ss_err.ERR_SUCCESS
	}

	// 根据代理找到上级代理accNo
	nextUpperNo := dao.AgencyDaoInst.QueryUpper(tx, upperAgencyNo)
	if nextUpperNo == "" {
		ss_log.Error("查询当前代理的上级代理失败,当前代理accNo为----->%s", upperAgencyNo)
		return ss_err.ERR_PARAM
	}
	upperFees := "0"
	newFees := fees
	// 查找代理的费率
	upperRate, err := dao.RateDaoInst.QueryRateFromAccNo(tx, upperAgencyNo, constants.RoleType_Agency, channelNo)
	if err != nil && !strings.Contains(err.Error(), DB_NO_ROWS_MSG) {
		ss_log.Error("err=[%v]", err)
		return ss_err.ERR_PARAM
	}
	if upperRate != "" {
		// 计算费率差
		upperRateDeci := ss_count.Sub(rate, upperRate)
		// 判断费率差是否会出现倒扣的情况
		if f, _ := upperRateDeci.Float64(); f < 0 { // 倒扣
			// 这个上级代理不分手续费,同时关闭商户的通道
			if errStr := dao.ChannelDaoInst.ForbidChannelFromNo(tx, channelNo, constants.CHANNEL_STATUS_FORBID); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[%v]", err)
				return errStr
			}
		} else {
			upperFeesDeci := ss_count.CountFees(amount, upperRateDeci.String(), "0")
			// 取整
			upperFees = ss_big.SsBigInst.ToRound(upperFeesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()

			// 记录日志进钱包表
			logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, upperAgencyNo, constants.RoleType_Agency)
			if logNo == "" {
				ss_log.Error("获取钱包id失败,accNo为----->%s", upperAgencyNo)
				return ss_err.ERR_PARAM
			}
			if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, upperFees, logNo, "+"); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("errStr=[%v]", errStr)
				return errStr
			}

			// 记录平台手续费分成钱包明细
			decription := "入金操作,分得手续费"
			if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, upperAgencyNo, beforeAmount, upperFees, "1",
				nextUpperNo, constants.FEES_TYPE_BILL, orderNo, logNo, decription, channelNo, merchantNo, "+"); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("errStr=[%v]", errStr)
				return errStr
			}

			ss_log.Info("代理 %s: 费率: %s,上级费率为: %s 清分手续费为-------->%s", upperAgencyNo, rate, upperRate, upperFees)

			newFees = ss_count.Sub(fees, upperFees).String()
			// todo 修改统计表,代理收入
			if errStr := dao.CumulativeCountDaoInst.ModifyCumulative2(tx, upperFees); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("修改统计利润表 代理收入 失败,订单号为-----> %s", orderNo)
				return errStr
			}
		}

	}

	//newFees := ss_count.Sub(fees, upperFees).String()

	if errStr := doAgencyFees1(tx, orderNo, nextUpperNo, channelNo, merchantNo, amount, newFees, upperRate,
		channelRate, totalFees, merchantAgencyNo, countFee); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("errStr=[%v]", errStr)
		return errStr
	}

	return ss_err.ERR_SUCCESS
}

// 提现(代付)
func WithdrawalService(ctx context.Context, req *businessSettleProto.BusinessSettleFeesRequest) string {
	if errStr := dao.TransferDaoInst.ConfirmIsNoCount(req.OrderNo); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[WithdrawalService 该手续费已统计,订单号为----->%s]", req.OrderNo)
		return ss_err.ERR_PARAM
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 获取订单信息
	roleType, settlementMethod, fees, accNo, channelNo, channelRate, amount, countFee := dao.TransferDaoInst.GetTransferInfoFromNo(tx, req.OrderNo)
	if roleType == "" || settlementMethod == "" {
		ss_log.Error("err=[获取角色类型失败,订单号为----->%s]", req.OrderNo)
		return ss_err.ERR_PARAM
	}
	switch settlementMethod {
	case constants.SettlementMethodWithdraw: // 提现
		if errStr := doWithdrawalFees(tx, req.OrderNo, fees, channelNo, accNo, roleType); errStr != ss_err.ERR_SUCCESS {
			// 修改统计失败
			if errStr := dao.TransferDaoInst.UpdateTransferOrderIsCount(req.OrderNo, constants.FEES_FAILE_COUNT); errStr != ss_err.ERR_SUCCESS {
				return errStr
			}
			return errStr
		}
	case constants.SettlementMethodTransfer: // 代付
		if roleType == constants.RoleType_Agency {
			ss_log.Error("代理没有代付操作, 订单号为--->%s", req.OrderNo)
			return ss_err.ERR_ACCOUNT_NO_PERMISSION
		}
		if errStr := doTransferFees(tx, req.OrderNo, amount, fees, channelNo, accNo, roleType, channelRate, countFee); errStr != ss_err.ERR_SUCCESS {
			// 修改统计失败
			if errStr := dao.TransferDaoInst.UpdateTransferOrderIsCount(req.OrderNo, constants.FEES_FAILE_COUNT); errStr != ss_err.ERR_SUCCESS {
				return errStr
			}
			return errStr

		}
	default:
		ss_log.Error("RoleType_Agency 结算方式不正确 ----->%s", settlementMethod)
		return ss_err.ERR_PARAM
	}

	// todo 更新商户出金(bill_result),out_go_amount, out_go_fee, out_go_real_amount
	if roleType == constants.RoleType_Merc {
		if errStr := dao.BillResultDaoInst.ModifyBillOutGoResult(tx, amount, fees, ss_count.Add(amount, fees), accNo, roleType); errStr != ss_err.ERR_SUCCESS {
			return errStr
		}
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

// 提现操作
func doWithdrawalFees(tx *sql.Tx, orderNo, fees, channelNo, accNo, roleType string) string {
	// 手续费直接进平台账号,确定平台账号是否存在
	logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, constants.INIT_UUID, constants.RoleType_Admin)
	if logNo == "" {
		ss_log.Error("获取钱包id失败,accNo为----->%s", constants.INIT_UUID)
		return ss_err.ERR_PARAM
	}
	if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, fees, logNo, "+"); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 记录平台手续费分成钱包明细
	var decription string
	switch roleType {
	case constants.RoleType_Agency:
		decription = fmt.Sprintf("%s", "代理 提现操作,分得手续费")
	case constants.RoleType_Merc:
		decription = fmt.Sprintf("%s", "商户 提现操作,分得手续费")
	}
	if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, constants.INIT_UUID, beforeAmount, fees, "1",
		constants.INIT_UUID, constants.FEES_TYPE_WITHDRAWAL, orderNo, logNo, decription, channelNo, accNo, "+"); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// todo 修改平台利润表(cumulative_count) 中的平台利润 headquarters_profit 字段
	if errStr := dao.CumulativeCountDaoInst.ModifyHeadquartersProfit(tx, fees); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("%s", "修改平台利润数据失败")
		return errStr
	}

	return ss_err.ERR_SUCCESS
}

// 代付
func doTransferFees(tx *sql.Tx, orderNo, amount, fees, channelNo, accNo, roleType, channelRate, countFee string) string {
	// 计算上游的手续费
	channelDeci := ss_count.CountFees(amount, channelRate, "0")
	upperFees := ss_big.SsBigInst.ToRound(channelDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
	// 如果存在单笔费用,则累加
	if strext.ToInt(countFee) > 0 {
		upperFees = ss_count.Add(upperFees, countFee)
	}
	plantFormFees := ss_count.Sub(fees, upperFees).String()

	var op string
	switch strext.ToInt(plantFormFees) > 0 {
	case true:
		op = "+"
	case false:
		op = "-"
	}

	// 确定平台账号是否存在
	logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, constants.INIT_UUID, constants.RoleType_Admin)
	if logNo == "" {
		ss_log.Error("获取钱包id失败,accNo为----->%s", constants.INIT_UUID)
		return ss_err.ERR_PARAM
	}
	if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, plantFormFees, logNo, op); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 记录平台手续费分成钱包明细
	var decription string
	switch roleType {
	case constants.RoleType_Agency:
		switch op {
		case "+":
			decription = fmt.Sprintf("%s", "代理 代付操作,分得手续费")
		case "-":
			decription = fmt.Sprintf("%s", "代理 代付操作,因手续费不够给上游,导致负扣手续费")
		}

	case constants.RoleType_Merc:
		switch op {
		case "+":
			decription = fmt.Sprintf("%s", "商户 代付操作,分得手续费")
		case "-":
			decription = fmt.Sprintf("%s", "商户 代付操作,因手续费不够给上游,导致负扣手续费")
		}

	}
	if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, constants.INIT_UUID, beforeAmount, plantFormFees, "1",
		constants.INIT_UUID, constants.FEES_TYPE_WITHDRAWAL, orderNo, logNo, decription, channelNo, accNo, op); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	if errStr := dao.CostDaoInst.InsertCost(tx, orderNo, constants.FEES_TYPE_WITHDRAWAL, upperFees, channelNo); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// todo 修改平台利润表(cumulative_count) 中的平台利润 headquarters_profit 字段
	if errStr := dao.CumulativeCountDaoInst.ModifyHeadquartersProfit(tx, plantFormFees); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("%s", "修改平台利润数据失败")
		return errStr
	}

	return ss_err.ERR_SUCCESS
}

/**
upperAccNo: 上级代理的
*/
//func doAgencyFees(tx *sql.Tx, upperAccNo, nextUpperAccNo, rateNo, upperRateNo, amount, fees, orderNo, channelRate, channelNo, merchantNo string) string {
//	if upperRateNo == constants.INIT_UUID { // 目前最上级的代理了,上面就是平台了
//		// 计算上游的手续费
//		channelDeci := ss_count.CountFees(amount, channelRate)
//		upperFees := ss_big.SsBigInst.ToRound(channelDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
//		plantFormFees := ss_count.Sub(fees, upperFees).String()
//		// 记录日志进钱包表
//		logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, upperAccNo)
//		if logNo == "" {
//			ss_log.Error("获取钱包id失败,accNo为----->%s", upperAccNo)
//			return ss_err.ERR_PARAM
//		}
//		if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, plantFormFees, logNo); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//		// 记录平台手续费分成钱包明细
//		decription := "入金操作,分得手续费"
//		if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, upperAccNo, beforeAmount, plantFormFees, "1",
//			nextUpperAccNo, constants.FEES_TYPE_BILL, orderNo, logNo, decription, channelNo, merchantNo); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//
//		amountT, feesT := dao.BillDaoInst.QueryBillInfo(tx, orderNo)
//		if amountT == "" || feesT == "" {
//			ss_log.Error("根据订单号查询金额和手续费失败,订单号为----->%s", orderNo)
//			return ss_err.ERR_PARAM
//		}
//		// 商户变动金额
//		merchantChangeAmount := ss_big.SsBigInst.ToRound(ss_count.Sub(amountT, feesT), 0, ss_big.RoundingMode_HALF_EVEN).String()
//
//		// 记录商户钱包
//		merchantLogNo, merchantBeforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, merchantNo)
//		if merchantLogNo == "" {
//			ss_log.Error("获取商户钱包id失败,accNo为----->%s", merchantNo)
//			return ss_err.ERR_PARAM
//		}
//		if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, merchantChangeAmount, merchantLogNo); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//		merchantAgenNo := dao.MerchantDaoInst.QueryAgencyNo(tx, merchantNo)
//		if merchantAgenNo == "" {
//			ss_log.Error("查询商户的代理失败,商户号为----->%s", merchantNo)
//			return ss_err.ERR_PARAM
//		}
//		// 记录商户钱包明细
//		merChantDecription := "商户入金操作"
//		if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, merchantNo, merchantBeforeAmount, merchantChangeAmount, "1",
//			merchantAgenNo, constants.FEES_TYPE_BILL, orderNo, merchantLogNo, merChantDecription, channelNo, merchantNo); errStr != ss_err.ERR_SUCCESS {
//			return errStr
//		}
//		return ss_err.ERR_SUCCESS
//	}
//	rate := dao.RateDaoInst.GetRateFromNo(tx, rateNo)
//	// 计算自身和上级的费率
//	upperRate := dao.RateDaoInst.GetRateFromNo(tx, upperRateNo)
//	// 计算二者之间的费率差从而获得上级的手续费
//	upperRateDeci := ss_count.Sub(rate, upperRate)
//	upperFeesDeci := ss_count.CountFees(amount, upperRateDeci.String())
//	// 取整
//	upperFees := ss_big.SsBigInst.ToRound(upperFeesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
//
//	// 记录日志进钱包表
//	logNo, beforeAmount := dao.WalletDaoInst.ConfirmExistWallet(tx, upperAccNo)
//	if logNo == "" {
//		ss_log.Error("获取钱包id失败,accNo为----->%s", upperAccNo)
//		return ss_err.ERR_PARAM
//	}
//	if errStr := dao.WalletDaoInst.ModifyAmountFromNo(tx, upperFees, logNo); errStr != ss_err.ERR_SUCCESS {
//		return errStr
//	}
//	// 记录平台手续费分成钱包明细
//	decription := "入金操作,分得手续费"
//	if errStr := dao.WalletDetailDaoInst.InsertDetail(tx, upperAccNo, beforeAmount, upperFees, "1",
//		nextUpperAccNo, constants.FEES_TYPE_BILL, orderNo, logNo, decription, channelNo, merchantNo); errStr != ss_err.ERR_SUCCESS {
//		return errStr
//	}
//
//	nextUpperRateNo := dao.RateDaoInst.QueryUpperRateNoFromRateNo(tx, upperRateNo)
//	nextUpperNo := dao.AgencyDaoInst.QueryUpper(tx, nextUpperAccNo) // nextUpperNo 可能为空,最顶级的代理是没有上级的了
//	if nextUpperNo == "" {
//		nextUpperNo = nextUpperAccNo
//	}
//	if nextUpperRateNo == "" {
//		return ss_err.ERR_PARAM
//	}
//	newFees := ss_count.Sub(fees, upperFees).String()
//	if errStr := doAgencyFees(tx, nextUpperAccNo, nextUpperNo, upperRateNo, nextUpperRateNo, amount, newFees, orderNo, channelRate, channelNo, merchantNo); errStr != ss_err.ERR_SUCCESS {
//		return errStr
//	}
//	return ss_err.ERR_SUCCESS
//}
