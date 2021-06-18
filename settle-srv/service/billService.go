package service

import (
	"context"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	go_micro_srv_settle "a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/settle-srv/dao"
)

// 兑换手续费处理
func ExchangeService(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 确定手续费没有统计
	if logNoT := dao.ExchangeOrderDaoInst.ConfirmIsNoCount(tx, req.BillNo); logNoT != "" {
		ss_log.Error("err=[ExchangeService 该手续费已统计,订单号为----->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[兑换手续费或者MoneyType为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}
	var vaType int32
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_FEES
	case "khr":
		vaType = constants.VaType_KHR_FEES
	default:
		return ss_err.ERR_PARAM
	}

	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, req.MoneyType, vaType)
	// 修改总部的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, req.Fees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入利润表
	d := &dao.HeadquartersProfit{
		OrderNo:      req.BillNo,
		Amount:       req.Fees,
		OrderStatus:  constants.OrderStatus_Paid,
		BalanceType:  strings.ToLower(req.MoneyType),
		ProfitSource: constants.ProfitSource_Exchange,
		OpType:       constants.PlatformProfitAdd,
	}
	_, err := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d)
	if err != nil {
		ss_log.Error("插入手续费盈利失败，data=%v, err=%v", strext.ToJson(d), err)
		return ss_err.ERR_SYSTEM
	}

	// 修改收益 总部虚账的余额是等于收益表中的可提现余额
	if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, req.Fees, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 修改手续费已经统计
	if errStr := dao.ExchangeOrderDaoInst.UpdateIsCountFromLogNo(tx, req.BillNo, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

// 转账手续费处理
func TransferService(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 确定手续费没有统计
	if logNoT := dao.TransferDaoInst.ConfirmIsNoCount(tx, req.BillNo); logNoT != "" {
		ss_log.Error("err=[TransferService 该手续费已统计,订单号为----->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[转账手续费或者MoneyType为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}
	var vaType int32
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_FEES
	case "khr":
		vaType = constants.VaType_KHR_FEES
	default:
		return ss_err.ERR_PARAM
	}

	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, req.MoneyType, vaType)
	// 修改总部的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, req.Fees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入利润表
	d := &dao.HeadquartersProfit{
		OrderNo:      req.BillNo,
		Amount:       req.Fees,
		OrderStatus:  constants.OrderStatus_Paid,
		BalanceType:  strings.ToLower(req.MoneyType),
		ProfitSource: constants.ProfitSource_TRANSFERFee,
		OpType:       constants.PlatformProfitAdd,
	}
	_, err := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d)
	if err != nil {
		ss_log.Error("插入手续费盈利失败，data=%v, err=%v", strext.ToJson(d), err)
		return ss_err.ERR_SYSTEM
	}

	// 修改收益 总部虚账的余额是等于收益表中的可提现余额
	if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, req.Fees, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 修改手续费已经统计
	if errStr := dao.TransferDaoInst.UpdateIsCountFromLogNo(tx, req.BillNo, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

// 收款手续费处理
func CollectionService(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 确定手续费没有统计
	if logNoT := dao.CollectionDaoInst.ConfirmIsNoCount(tx, req.BillNo); logNoT != "" {
		ss_log.Error("err=[CollectionService 该手续费已统计,订单号为----->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[收款手续费或者 MoneyType 为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	var vaType int32
	switch req.MoneyType {
	case "usd":
		vaType = constants.VaType_USD_FEES
	case "khr":
		vaType = constants.VaType_KHR_FEES
	default:
		return ss_err.ERR_PARAM
	}

	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, req.MoneyType, vaType)
	// 修改总部的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, req.Fees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入利润表
	d := &dao.HeadquartersProfit{
		OrderNo:      req.BillNo,
		Amount:       req.Fees,
		OrderStatus:  constants.OrderStatus_Paid,
		BalanceType:  strings.ToLower(req.MoneyType),
		ProfitSource: constants.ProfitSource_COLLECTION,
		OpType:       constants.PlatformProfitAdd,
	}
	_, err := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d)
	if err != nil {
		ss_log.Error("插入手续费盈利失败，data=%v, err=%v", strext.ToJson(d), err)
		return ss_err.ERR_SYSTEM
	}

	// 修改收益 总部虚账的余额是等于收益表中的可提现余额
	if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, req.Fees, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 修改手续费已经统计
	if errStr := dao.CollectionDaoInst.UpdateIsCountFromLogNo(tx, req.BillNo, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

// 存钱手续费处理
func SavemoneyService(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 确定手续费没有统计
	if logNoT := dao.IncomeOrderDaoInst.ConfirmIsNoCount(tx, req.BillNo); logNoT != "" {
		ss_log.Error("err=[CollectionService 该手续费已统计,订单号为----->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[存钱手续费或者MoneyType为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 通过logNo找到服务商id,找到对应的服务商虚账id
	srvAccNo, servicerNo, opAccNo, opAccType := dao.IncomeOrderDaoInst.QuerySrvAccNoFromLogNo(req.BillNo)
	if srvAccNo == "" {
		ss_log.Error("err=[存钱手续费或者MoneyType为空,服务商的 account id为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 获取分成比例
	incomeSharing, _ := dao.ServiceDaoInst.GetSharingFromSrvNo(srvAccNo)
	// 计算总部,服务商能分到多少手续费
	headFees, srvFees := ss_count.CountSharing(req.Fees, incomeSharing)
	ss_log.Info("客户通过pos机存款,订单号为----->%s,一共手续费为----->%s,服务商分成比例为----->%s,总部得到的手续费为----->%s,服务商得到的手续费为----->%s",
		req.BillNo, req.Fees, incomeSharing, headFees, srvFees)

	var vaType int32
	var servRealType, servQuotaType int32
	switch req.MoneyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_FEES
		servRealType = constants.VaType_QUOTA_USD_REAL
		servQuotaType = constants.VaType_QUOTA_USD
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_FEES
		servRealType = constants.VaType_QUOTA_KHR_REAL
		servQuotaType = constants.VaType_QUOTA_KHR
	default:
		return ss_err.ERR_PARAM
	}
	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, req.MoneyType, vaType)
	srvVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, vaType)
	srvRealVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servRealType)
	srvQuotaVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servQuotaType)
	// 修改总部的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, headFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, headAcc); errStr != ss_err.ERR_SUCCESS {
	//	return errStr
	//}

	// 修改服务商的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, srvVaccNo, srvFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入平台利润表
	d := &dao.HeadquartersProfit{
		OrderNo:      req.BillNo,
		Amount:       req.Fees,
		OrderStatus:  constants.OrderStatus_Paid,
		BalanceType:  strings.ToLower(req.MoneyType),
		ProfitSource: constants.ProfitSource_INCOME,
		OpType:       constants.PlatformProfitAdd,
	}
	_, err := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d)
	if err != nil {
		ss_log.Error("插入手续费盈利失败，data=%v, err=%v", strext.ToJson(d), err)
		return ss_err.ERR_SYSTEM
	}

	// 修改收益 总部虚账的余额是等于收益表中的可提现余额
	if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, headFees, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	amount, _, _, _, _, _, _, _ := dao.IncomeOrderDaoInst.QueryIncomeOrder(req.BillNo, constants.OrderStatus_Paid)
	// 插入服务商利润表
	if errStr := dao.ServicerprofitledgerDaoInstance.InsertServicerProfitLedger(tx, req.BillNo, amount, req.Fees, incomeSharing, srvFees, servicerNo, req.MoneyType, constants.Order_Type_In_Come); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 更新手续费分成到服务商实时额度,可用余额是做成跟信用卡一样的,相减
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpper(tx, srvRealVaccNo, srvFees, "-", req.BillNo, constants.VaReason_FEES, srvQuotaVaccNo); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入服务商交易明细
	var accountType string
	switch strext.ToInt(opAccType) {
	case constants.OpAccType_Servicer: // 服务商
		accountType = constants.AccountType_SERVICER
	case constants.OpAccType_Pos: // 店员
		accountType = constants.AccountType_POS
	}
	if errStr := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, srvFees, req.MoneyType, srvAccNo, accountType, req.BillNo, "0", constants.OrderStatus_Paid, servicerNo, opAccNo, constants.BillDetailTypeProfit, "0", srvFees); errStr == "" {
		return ss_err.ERR_PARAM
	}
	// 修改手续费已经统计
	if errStr := dao.IncomeOrderDaoInst.UpdateIsCountFromLogNo(tx, req.BillNo, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

// 取款手续费处理
func WithdrawService(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 确定手续费没有统计
	if logNoT := dao.OutgoOrderDaoInst.ConfirmIsNoCount(tx, req.BillNo); logNoT != "" {
		ss_log.Error("err=[WithdrawService 该手续费已统计,订单号为----->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[取款手续费或者MoneyType为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 通过logNo找到服务商id,找到对应的服务商虚账id
	srvAccNo, servicerNo, opAccNo, opAccType := dao.OutgoOrderDaoInst.QuerySrvAccNoFromLogNo(req.BillNo)
	if srvAccNo == "" {
		ss_log.Error("err=[取款,服务商的 account id为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 获取分成比例
	_, outGoSharing := dao.ServiceDaoInst.GetSharingFromSrvNo(srvAccNo)
	// 计算总部,服务商能分到多少手续费
	headFees, srvFees := ss_count.CountSharing(req.Fees, outGoSharing)
	ss_log.Info("客户取款,订单号为----->%s,一共手续费为----->%s,服务商分成比例为----->%s,总部得到的手续费为----->%s,服务商得到的手续费为----->%s",
		req.BillNo, req.Fees, outGoSharing, headFees, srvFees)

	var vaType int32
	var servRealType, servQuotaType int32

	switch req.MoneyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_FEES
		servRealType = constants.VaType_QUOTA_USD_REAL
		servQuotaType = constants.VaType_QUOTA_USD
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_FEES
		servRealType = constants.VaType_QUOTA_KHR_REAL
		servQuotaType = constants.VaType_QUOTA_KHR
	default:
		return ss_err.ERR_PARAM
	}
	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, req.MoneyType, vaType)
	srvVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, vaType)
	srvRealVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servRealType)
	srvQuotaVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servQuotaType)

	// 修改总部的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, headFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 修改服务商的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, srvVaccNo, srvFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入平台利润表
	d := &dao.HeadquartersProfit{
		OrderNo:      req.BillNo,
		Amount:       headFees,
		OrderStatus:  constants.OrderStatus_Paid,
		BalanceType:  strings.ToLower(req.MoneyType),
		ProfitSource: constants.ProfitSource_WithdrawFee,
		OpType:       constants.PlatformProfitAdd,
	}
	_, err := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d)
	if err != nil {
		ss_log.Error("插入手续费盈利失败，data=%v, err=%v", strext.ToJson(d), err)
		return ss_err.ERR_SYSTEM
	}

	// 修改收益 总部虚账的余额是等于收益表中的可提现余额
	if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, headFees, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	amount, _, _, _, _, _, withdrawType, _ := dao.OutgoOrderDaoInst.QueryOutGoOrderFromLogNo(req.BillNo, constants.OrderStatus_Paid)

	ss_log.Info("amount[%v]----withdrawType[%v]", amount, withdrawType)

	orderTypeWithdraw := ""
	switch withdrawType {
	case constants.OutgoOrderPaymentType_MobileNum:
		orderTypeWithdraw = constants.Order_Type_Mobile_Withdraw
	case constants.OutgoOrderPaymentType_Sweep:
		fallthrough
	case constants.OutgoOrderPaymentType_SweepAll:
		orderTypeWithdraw = constants.Order_Type_Sweep_Withdraw
	default:
		ss_log.Error("未知取款订单类型[%v]", withdrawType)

	}

	// 插入服务商利润表
	if errStr := dao.ServicerprofitledgerDaoInstance.InsertServicerProfitLedger(tx, req.BillNo, amount, req.Fees, outGoSharing, srvFees, servicerNo, req.MoneyType, orderTypeWithdraw); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 更新手续费分成到服务商实时额度,可用余额是做成跟信用卡一样的,相减
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpper(tx, srvRealVaccNo, srvFees, "-", req.BillNo, constants.VaReason_FEES, srvQuotaVaccNo); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入服务商交易明细
	var accountType string
	switch strext.ToInt(opAccType) {
	case constants.OpAccType_Servicer: // 服务商
		accountType = constants.AccountType_SERVICER
	case constants.OpAccType_Pos: // 店员
		accountType = constants.AccountType_POS
	}
	if errStr := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, srvFees, req.MoneyType, srvAccNo, accountType, req.BillNo, "0", constants.OrderStatus_Paid, servicerNo, opAccNo, constants.BillDetailTypeProfit, "0", srvFees); errStr == "" {
		return ss_err.ERR_PARAM
	}

	// 修改手续费已经统计
	if errStr := dao.OutgoOrderDaoInst.UpdateIsCountFromLogNo(tx, req.BillNo, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

/*
// 手机号取款手续费处理
func MobileNumWithdrawService(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	// 确定手续费没有统计
	if logNoT := dao.OutgoOrderDaoInst.ConfirmIsNoCount(tx, req.BillNo); logNoT != "" {
		ss_log.Error("err=[MobileNumWithdrawService 该手续费已统计,订单号为----->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[手机号取款手续费或者MoneyType为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 通过logNo找到服务商id,找到对应的服务商虚账id
	srvAccNo, servicerNo, opAccNo, opAccType := dao.OutgoOrderDaoInst.QuerySrvAccNoFromLogNo(req.BillNo)
	if srvAccNo == "" {
		ss_log.Error("err=[手机号取款,服务商的 account id为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 获取分成比例
	_, outGoSharing := dao.ServiceDaoInst.GetSharingFromSrvNo(srvAccNo)
	// 计算总部,服务商能分到多少手续费
	headFees, srvFees := ss_count.CountSharing(req.Fees, outGoSharing)
	ss_log.Info("客户通过pos机取款,订单号为----->%s,一共手续费为----->%s,服务商分成比例为----->%s,总部得到的手续费为----->%s,服务商得到的手续费为----->%s",
		req.BillNo, req.Fees, outGoSharing, headFees, srvFees)

	var vaType int32
	var servRealType, servQuotaType int32

	switch req.MoneyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_FEES
		servRealType = constants.VaType_QUOTA_USD_REAL
		servQuotaType = constants.VaType_QUOTA_USD
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_FEES
		servRealType = constants.VaType_QUOTA_KHR_REAL
		servQuotaType = constants.VaType_QUOTA_KHR
	default:
		return ss_err.ERR_PARAM
	}
	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, req.MoneyType, vaType)
	srvVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, vaType)
	srvRealVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servRealType)
	srvQuotaVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servQuotaType)

	// 修改总部的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, headFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 修改服务商的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, srvVaccNo, srvFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 同步平台账户的余额
	if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, headAcc); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 同步服务商账户的余额
	if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, srvAccNo); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入平台利润表
	if errStr := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, req.BillNo, headFees, constants.OrderStatus_Paid, req.MoneyType, "1"); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 修改收益 总部虚账的余额是等于收益表中的可提现余额
	if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, headFees, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	amount, _, _, _, _, _, _ := dao.OutgoOrderDaoInst.QueryOutGoOrderFromLogNo(req.BillNo, constants.OrderStatus_Paid)
	// 插入服务商利润表
	if errStr := dao.ServicerprofitledgerDaoInstance.InsertServicerProfitLedger(tx, req.BillNo, amount, req.Fees, outGoSharing, srvFees, servicerNo, req.MoneyType, constants.Order_Type_Mobile_Withdraw); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 更新手续费分成到服务商实时额度,可用余额是做成跟信用卡一样的,相减
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpper(tx, srvRealVaccNo, srvFees, "-", req.BillNo, constants.VaReason_FEES, srvQuotaVaccNo); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入服务商交易明细
	var accountType string
	switch strext.ToInt(opAccType) {
	case constants.OpAccType_Servicer: // 服务商
		accountType = constants.AccountType_SERVICER
	case constants.OpAccType_Pos: // 店员
		accountType = constants.AccountType_POS
	}
	if errStr := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, srvFees, req.MoneyType, srvAccNo, accountType, req.BillNo, "0", constants.OrderStatus_Paid, servicerNo, opAccNo, constants.BillDetailTypeProfit, "0", srvFees); errStr == "" {
		return ss_err.ERR_PARAM
	}

	// 修改手续费已经统计
	if errStr := dao.OutgoOrderDaoInst.UpdateIsCountFromLogNo(tx, req.BillNo, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

// 扫码取款手续费处理
func SweepWithdrawService(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) string {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[扫码取款手续费或者MoneyType为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 确定手续费没有统计
	if logNoT := dao.OutgoOrderDaoInst.ConfirmIsNoCount(tx, req.BillNo); logNoT != "" {
		ss_log.Error("err=[SweepWithdrawService 该手续费已统计,订单号为----->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	if req.Fees == "0" || req.MoneyType == "" {
		ss_log.Error("err=[扫一扫取款手续费或者 MoneyType 为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}

	// 通过logNo找到服务商id,找到对应的服务商虚账id
	srvAccNo, servicerNo, opAccNo, opAccType := dao.OutgoOrderDaoInst.QuerySrvAccNoFromLogNo(req.BillNo)
	if srvAccNo == "" {
		ss_log.Error("err=[扫一扫取款,服务商的 account id为空,订单号为--------->%s]", req.BillNo)
		return ss_err.ERR_PARAM
	}
	// 获取分成比例
	_, outGoSharing := dao.ServiceDaoInst.GetSharingFromSrvNo(srvAccNo)
	// 计算总部,服务商能分到多少手续费
	headFees, srvFees := ss_count.CountSharing(req.Fees, outGoSharing)
	ss_log.Info("客户通过pos机扫码取款,订单号为----->%s,一共手续费为----->%s,服务商分成比例为----->%s,总部得到的手续费为----->%s,服务商得到的手续费为----->%s",
		req.BillNo, req.Fees, outGoSharing, headFees, srvFees)

	var vaType int32
	var servRealType, servQuotaType int32
	switch req.MoneyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_FEES
		servRealType = constants.VaType_QUOTA_USD_REAL
		servQuotaType = constants.VaType_QUOTA_USD
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_FEES
		servRealType = constants.VaType_QUOTA_KHR_REAL
		servQuotaType = constants.VaType_QUOTA_KHR
	default:
		return ss_err.ERR_PARAM
	}
	// 查询总部的账号
	_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVacc := dao.VaccountDaoInst.ConfirmExistVaccount(headAcc, req.MoneyType, vaType)
	srvVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, vaType)
	srvRealVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servRealType)
	srvQuotaVaccNo := dao.VaccountDaoInst.ConfirmExistVaccount(srvAccNo, req.MoneyType, servQuotaType)

	// 修改总部的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, headFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 修改服务商的临时虚账余额
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, srvVaccNo, srvFees, "+", req.BillNo, constants.VaReason_FEES); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 同步平台账户的余额
	if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, headAcc); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 同步服务商账户的余额
	if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, srvAccNo); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入平台利润表
	if errStr := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, req.BillNo, headFees, constants.OrderStatus_Paid, req.MoneyType, "1"); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 修改收益 总部虚账的余额是等于收益表中的可提现余额
	if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, headFees, req.MoneyType); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	amount, _, _, _, _, _, _ := dao.OutgoOrderDaoInst.QueryOutGoOrderFromLogNo(req.BillNo, constants.OrderStatus_Paid)
	// 插入服务商利润表
	if errStr := dao.ServicerprofitledgerDaoInstance.InsertServicerProfitLedger(tx, req.BillNo, amount, req.Fees, outGoSharing, srvFees, servicerNo, req.MoneyType, constants.Order_Type_Sweep_Withdraw); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}
	// 更新手续费分成到服务商实时额度,可用余额是做成跟信用卡一样的,相减
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpper(tx, srvRealVaccNo, srvFees, "-", req.BillNo, constants.VaReason_FEES, srvQuotaVaccNo); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	// 插入服务商交易明细
	var accountType string
	switch strext.ToInt(opAccType) {
	case constants.OpAccType_Servicer: // 服务商
		accountType = constants.AccountType_SERVICER
	case constants.OpAccType_Pos: // 店员
		accountType = constants.AccountType_POS
	}
	if errStr := dao.BillingDetailsResultsDaoInstance.InsertResult(tx, srvFees, req.MoneyType, srvAccNo, accountType, req.BillNo, "0", constants.OrderStatus_Paid, servicerNo, opAccNo, constants.BillDetailTypeProfit, "0", srvFees); errStr == "" {
		return ss_err.ERR_PARAM
	}

	// 修改手续费已经统计
	if errStr := dao.OutgoOrderDaoInst.UpdateIsCountFromLogNo(tx, req.BillNo, constants.FEES_IS_COUNT); errStr != ss_err.ERR_SUCCESS {
		return errStr
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}
*/
