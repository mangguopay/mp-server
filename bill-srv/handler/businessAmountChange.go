package handler

import (
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/model"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	"context"
	"fmt"
	"strings"

	"a.a/cu/db"
	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/bill-srv/i"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	authProto "a.a/mp-server/common/proto/auth"
	billProto "a.a/mp-server/common/proto/bill"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

// 商家中心-商家提现
func (b *BillHandler) BusinessWithdraw(ctx context.Context, req *billProto.BusinessWithdrawRequest, reply *billProto.BusinessWithdrawReply) error {
	//验证支付密码
	replyCheckPayPwd, errCheckPayPwd := i.AuthHandlerInst.Client.CheckPayPWD(ctx, &authProto.CheckPayPWDRequest{
		AccountUid:  req.AccountUid,
		AccountType: req.AccountType,
		Password:    req.PayPwd,
		NonStr:      req.NonStr,
		IdenNo:      req.IdenNo,
	})
	if errCheckPayPwd != nil {
		ss_log.Error("调用验证支付密码接口出错,err[%v]", errCheckPayPwd)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	if replyCheckPayPwd.ResultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("支付密码错误。req[%+v]", req)
		reply.ResultCode = ss_err.ERR_BusinessPayPwd_FAILD
		return nil
	}

	// 判断提现金额
	if strext.ToFloat64(req.Amount) <= 0 || req.Amount == "" {
		ss_log.Error("err=[商家[%v]提现, 金额为0或者为空,传入的金额为----->%s]", req.AccountUid, req.Amount)
		reply.ResultCode = ss_err.ERR_WALLET_AMOUNT_NULL
		return nil
	}

	// 判断用户的卡是否存在
	_, _, balanceType, channelNo := dao.CardBusinessDaoInst.QueryNameAndNumFromNo(req.CardNo)
	if balanceType != req.MoneyType {
		ss_log.Error("银行卡的币种和取款的币种不一致,数据库的币种为: %s,取款的币种为: %s", balanceType, req.MoneyType)
		reply.ResultCode = ss_err.ERR_REC_CARD_NUM_FAILD
		return nil
	}

	// 获取渠道手续费率
	// 提现手续费率      单笔提现最大金额      提现单笔手续费      提现计算手续费类型(1-按比例收取手续费，2按单笔手续费收取)
	withdrawRate, withdrawMaxAmount, withdrawSingleMinFee, withdrawChargeType := dao.ChannelBusinessDaoInst.QueryChannelBusinessWithdrawInfoFromNo(channelNo, balanceType)
	if withdrawChargeType == "" {
		ss_log.Error("根据 channelNo 查询 渠道信息失败,channelNo: %s", channelNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 判断限额
	if withdrawMaxAmount != "" {
		if strext.ToFloat64(req.Amount) > strext.ToFloat64(withdrawMaxAmount) {
			ss_log.Error("用户向总部提现,交易金额超出单笔最大金额,交易金额为: %s,单笔提现最大金额为: %s", req.Amount, withdrawMaxAmount)
			reply.ResultCode = ss_err.ERR_LOCAL_RULE_EXCEED_AMOUNT
			return nil
		}
	}

	//获取手续费
	var fees string
	switch withdrawChargeType {
	case constants.Charge_Type_Rate: // 手续费比例收取
		feesDeci := ss_count.CountFees(req.Amount, withdrawRate, "0")
		// 取整
		fees = ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
	case constants.Charge_Type_Count: // 按单笔手续费收取
		fees = withdrawSingleMinFee
	}

	// 根据币种获取虚账类型
	var vaType int32
	switch req.MoneyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_BUSINESS_SETTLED
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_BUSINESS_SETTLED
	default:
		ss_log.Error("商家提现,币种错误,MoneyType: %s", req.MoneyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 确保虚拟账号存在
	recvVaccNo := InternalCallHandlerInst.ConfirmExistVAccount(req.AccountUid, req.MoneyType, strext.ToInt32(vaType))

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	balance, _ := dao.VaccountDaoInst.GetBalance(recvVaccNo)
	// 判断是普通提现还是全部提现,获取金额,计算手续费
	var withdrawAmount string
	switch req.WithdrawType {
	case constants.WITHDRAWAL_TYPE_ORDINARY:
		/**普通提现
		申请金额 + 手续费 <= 账户余额，申请金额 = 到账金额
		*/
		f, _ := ss_count.Sub(balance, ss_count.Add(req.Amount, fees)).Float64()
		if f < 0 {
			ss_log.Error("普通提现,商家账户[%v]余额不足以提现,提交的金额为:%s,账户余额为:%s,手续费为:%s", req.AccountUid, req.Amount, balance, fees)
			reply.ResultCode = ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
			return nil
		}
		withdrawAmount = req.Amount
	case constants.WITHDRAWAL_TYPE_ALL:
		/**全部提现
		申请金额 + 手续费 = 账户余额， 到账金额 = 申请金额 - 手续费
		*/
		if req.Amount != balance {
			ss_log.Error("全部提现,提交的金额和商家账户余额对应不上,提交的金额为: %s,账户余额为: %s", req.Amount, balance)
			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_FAILD
			return nil
		}
		withdrawAmountDeci := ss_count.Sub(balance, fees)
		f, _ := withdrawAmountDeci.Float64()
		if f <= 0 {
			ss_log.Error("err=[全部提现手续费不够扣,当前余额为----->%s,应扣手续费为---->%s]", balance, fees)
			reply.ResultCode = ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH
			return nil
		}
		withdrawAmount = withdrawAmountDeci.String()
	default:
		ss_log.Error("商家提现,提现类型错误,WithdrawType: %s", req.WithdrawType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_log.Info("商家[%v]提现,手续费类型为:%s,提现金额为:%s,计算手续费率为:%s,手续费为:%s", req.AccountUid, withdrawChargeType, req.Amount, withdrawRate, fees)

	logNo := dao.LogToBusinessDaoInst.Insert(tx, req.MoneyType, req.IdenNo, constants.COLLECTION_TYPE_BANK_TRANSFER,
		req.CardNo, req.Amount, withdrawAmount, fees)
	if logNo == "" {
		reply.ResultCode = ss_err.ERR_LOG_TO_BUSINESS_FAILD
		return nil
	}

	//减少商户金额，增加冻结金额
	// 判断输入的金额是否超额（如果余额小于0则会报错）
	if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", recvVaccNo, withdrawAmount, logNo, constants.VaReason_Business_Withdraw, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
		if errStr == ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH {
			reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
			return nil
		}
		reply.ResultCode = errStr
		return nil
	}
	if fees != "0" && fees != "" {
		// 修改手续费
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "-", recvVaccNo, fees, logNo, constants.VaReason_FEES, constants.VaOpType_Freeze); errStr != ss_err.ERR_SUCCESS {
			if errStr == ss_err.ERR_WITHDRAW_AMT_NOT_ENOUGH {
				reply.ResultCode = ss_err.ERR_PAY_AMT_NOT_ENOUGH
				return nil
			}
			reply.ResultCode = errStr
			return nil
		}
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.LogNo = logNo
	return nil
}

// 管理后台-商家提现订单审核
func (b *BillHandler) UpdateToBusinessStatus(ctx context.Context, req *billProto.UpdateToBusinessStatusRequest, reply *billProto.UpdateToBusinessStatusReply) error {
	if req.LogNo == "" {
		ss_log.Error("订单号为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 判断订单是否是初始化状态
	logData, errGet := dao.LogToBusinessDaoInst.GetToBusinessDetail(req.LogNo)
	if errGet != nil {
		ss_log.Error("查询商家提现订单[%v]信息出错,err[%v]", req.LogNo, errGet)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 根据币种获取虚账类型
	var vaType, plantVaType int32
	switch logData.CurrencyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_BUSINESS_SETTLED
		plantVaType = constants.VaType_USD_FEES
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_BUSINESS_SETTLED
		plantVaType = constants.VaType_KHR_FEES
	default:
		ss_log.Error("用户向总部提现,币种错误,MoneyType: %s", logData.CurrencyType)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil

	}
	accNo := dao.BusinessDaoInst.GetAccNoByBusinessNo(logData.BusinessNo)
	if accNo == "" {
		ss_log.Error("查询不到商家[%v]的账号uid", logData.BusinessNo)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	// 确保虚拟账号存在
	recvVaccNo := InternalCallHandlerInst.ConfirmExistVAccount(accNo, logData.CurrencyType, strext.ToInt32(vaType))

	imageId := ""
	description := "" //关键操作日志描述

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	switch strext.ToStringNoPoint(req.OrderStatus) {
	case constants.WithdrawalOrderStatusPassed: //审核通过，已受理
		description = fmt.Sprintf("审核商家提现订单[%v],操作[%v] ", req.LogNo, "审核通过")

		if logData.OrderStatus != constants.WithdrawalOrderStatusPending {
			ss_log.Error(" 订单不是在待审核状态,不能审核商家提现订单,OrderNo: %s,订单状态为: %s", req.LogNo, logData.OrderStatus)
			reply.ResultCode = ss_err.ERR_ORDER_STATUS_NO_INIT
			return nil
		}

	case constants.WithdrawalOrderStatusSuccess: //完成
		description = fmt.Sprintf("完成商家提现订单[%v],操作[%v] ", req.LogNo, "完成")

		if logData.OrderStatus != constants.WithdrawalOrderStatusPassed {
			ss_log.Error(" 订单不是审核通过，已受理状态,不能完成商家提现订单,OrderNo: %s,订单状态为: %s", req.LogNo, logData.OrderStatus)
			reply.ResultCode = ss_err.ERR_ORDER_STATUS_NO_INIT
			return nil
		}

		//上传凭证
		if req.ImgStr == "" {
			ss_log.Error("未上传凭证图片req.ImageBase64为空")
			reply.ResultCode = ss_err.ERR_PARAM
			return nil
		}
		upReply, errU := i.CustHandleInstance.Client.UploadImage(ctx, &custProto.UploadImageRequest{
			ImageStr:     req.ImgStr,
			AccountUid:   req.AccountUid,
			Type:         constants.UploadImage_Auth,
			AddWatermark: constants.AddWatermark_True,
		}, global.RequestTimeoutOptions)
		if errU != nil {
			ss_log.Error("errU=[%v]", errU)
			reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
			return nil
		}

		if upReply.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("errU=[%v]", errU)
			reply.ResultCode = ss_err.ERR_SAVE_IMAGE_FAILD
			return nil
		}

		imageId = upReply.ImageId

		// 修改冻结余额
		//amountB := ss_count.Add(amount,fees)
		if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero(tx, recvVaccNo, logData.RealAmount, req.LogNo, constants.VaReason_Business_Withdraw, logData.Fee); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		if logData.Fee != "" && logData.Fee != "0" {
			//继续减手续费
			if errStr := dao.VaccountDaoInst.ModifyVaccFrozenUpperZero(tx, recvVaccNo, logData.Fee, req.LogNo, constants.VaReason_FEES, ""); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}

			// ============平台收益====================
			// 查询总部的账号
			_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
			// 确保虚拟账号存在
			headVacc := InternalCallHandlerInst.ConfirmExistVAccount(headAcc, logData.CurrencyType, plantVaType)

			// 修改总部的临时虚账余额
			if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, logData.Fee, "+", req.LogNo, constants.VaReason_Business_Withdraw); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}

			// 插入利润表
			d := &dao.HeadquartersProfit{
				GeneralLedgerNo: req.LogNo,
				Amount:          logData.Fee,
				OrderStatus:     constants.OrderStatus_Paid,
				BalanceType:     strings.ToLower(logData.CurrencyType),
				ProfitSource:    constants.ProfitSource_ToBusinessFee,
				OpType:          constants.PlatformProfitAdd,
			}
			if errStr := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}

			// 修改收益 总部虚账的余额是等于收益表中的可提现余额
			if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, logData.Fee, logData.CurrencyType); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
		}
	case constants.WithdrawalOrderStatusDeny: //审核不通过
		description = fmt.Sprintf("审核商家提现订单[%v],操作[%v]", req.LogNo, "审核不通过")
		fallthrough
	case constants.WithdrawalOrderStatusFail: //提现失败
		if description == "" {
			description = fmt.Sprintf("审核商家提现订单[%v],操作[%v]", req.LogNo, "未完成")
		}

		// 恢复余额
		if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "+", recvVaccNo, logData.RealAmount, req.LogNo, constants.VaReason_Business_Cancel_Withdraw, constants.VaOpType_Defreeze_Add); errStr != ss_err.ERR_SUCCESS {
			reply.ResultCode = errStr
			return nil
		}
		if logData.Fee != "" && logData.Fee != "0" {
			if errStr := dao.VaccountDaoInst.ModifyVaccRemainAndFrozenUpperZero(tx, "+", recvVaccNo, logData.Fee, req.LogNo, constants.VaReason_FEES, constants.VaOpType_Defreeze_Add); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
		}
	default:
		ss_log.Error("需要修改的订单状态有误,req.OrderStatus: %v", req.OrderStatus)
		reply.ResultCode = ss_err.ERR_OPERATE_FAILD
		return nil
	}

	// 修改状态
	if err := dao.LogToBusinessDaoInst.UpdateStatusFromLogNo(tx, req.LogNo, req.OrderStatus, req.Notes, imageId); err != nil {
		ss_log.Error("修改订单状态失败,err=[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 添加关键操作日志
	errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, req.AccountUid, constants.LogAccountWebType_Trading_Order)
	if errAddLog != nil {
		ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
		reply.ResultCode = ss_err.ERR_SYS_DB_OP
		return nil
	}

	ss_sql.Commit(tx)
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

//管理后台-审核商家充值
func (b *BillHandler) UpdateBusinessToHeadStatus(ctx context.Context, req *billProto.UpdateBusinessToHeadStatusRequest, reply *billProto.UpdateBusinessToHeadStatusReply) error {
	whereList := []*model.WhereSqlCond{
		{Key: "bth.log_no", Val: req.LogNo, EqType: "="},
	}
	logData, err := dao.LogBusinessToHeadquartersDaoInst.GetBusinessToHeadDetail(whereList)
	if err != nil {
		ss_log.Error("查询订单[%v]信息出错,err[%v]", req.LogNo, err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if logData.CurrencyType == "" || logData.ArriveAmount == "" || logData.BusinessAccNo == "" {
		ss_log.Error("根据订单号查询出的必要参数CurrencyType、Amount、BusinessAccNo其中有为空的,logData[%+v]", logData)
		reply.ResultCode = ss_err.ERR_PAY_QUERY_ORDER_ERROR
		return nil
	}

	// 检查状态
	if logData.OrderStatus != constants.AuditOrderStatus_Pending {
		ss_log.Error("订单[%v]不是待审核状态", req.LogNo)
		reply.ResultCode = ss_err.ERR_PAY_ORDER_STATUS_MISTAKE
		return nil
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
		reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
		return nil
	}
	defer ss_sql.Rollback(tx)

	//修改订单状态
	if err2 := dao.LogBusinessToHeadquartersDaoInst.ModifyStatus(tx, req.LogNo, req.Status); err2 != nil {
		ss_log.Error("修改订单[%v]订单状态为[%v]出错。err=[%v]", req.LogNo, req.Status, err2)
		reply.ResultCode = ss_err.ERR_Audit_FAILD
		return nil
	}

	switch req.Status {
	case constants.AuditOrderStatus_Passed: //通过
		// 调用ps rpc,进行加商家余额
		quotaRepl, quotaErr := i.QuotaHandleInstance.Client.ModifyQuota(context.TODO(), &go_micro_srv_quota.ModifyQuotaRequest{
			CurrencyType: logData.CurrencyType,
			Amount:       logData.ArriveAmount,
			AccountNo:    logData.BusinessAccNo,
			OpType:       constants.QuotaOp_BusinessSave,
			LogNo:        req.LogNo,
		})
		if quotaErr != nil {
			ss_log.Error("调用远程ModifyQuota rpc失败,err: %v", quotaErr)
			reply.ResultCode = ss_err.ERR_SYS_REMOTE_API_ERR
			return nil
		}

		if quotaRepl.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[--------------->%s]", "调用服务失败,操作为商家充值")
			reply.ResultCode = quotaRepl.ResultCode
			return nil
		}

		// ============平台收益====================
		if logData.Fee != "" && strext.ToInt64(logData.Fee) > 0 {
			var plantVaType int32
			switch logData.CurrencyType {
			case constants.CURRENCY_USD:
				plantVaType = constants.VaType_USD_FEES
			case constants.CURRENCY_KHR:
				plantVaType = constants.VaType_KHR_FEES
			default:
				ss_log.Error("商家充值订单[%v],币种错误,MoneyType: %s", logData.LogNo, logData.CurrencyType)
				reply.ResultCode = ss_err.ERR_PARAM
				return nil
			}

			// 查询总部的账号
			_, headAcc, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
			// 确保虚拟账号存在
			headVacc := InternalCallHandlerInst.ConfirmExistVAccount(headAcc, logData.CurrencyType, plantVaType)

			// 修改总部的临时虚账余额
			if errStr := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVacc, logData.Fee, "+", req.LogNo, constants.VaReason_Business_Save); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}

			// 插入利润表
			d := &dao.HeadquartersProfit{
				GeneralLedgerNo: req.LogNo,
				Amount:          logData.Fee,
				OrderStatus:     constants.OrderStatus_Paid,
				BalanceType:     strings.ToLower(logData.CurrencyType),
				ProfitSource:    constants.ProfitSource_BusinessToFee,
				OpType:          constants.PlatformProfitAdd,
			}
			if errStr := dao.HeadquartersProfitDaoInstance.InsertHeadquartersProfit(tx, d); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}

			// 修改收益 总部虚账的余额是等于收益表中的可提现余额
			if errStr := dao.HeadquartersProfitCashableDaoInstance.SyncAccProfit(tx, headVacc, logData.Fee, logData.CurrencyType); errStr != ss_err.ERR_SUCCESS {
				reply.ResultCode = errStr
				return nil
			}
		}
	case constants.AuditOrderStatus_Deny:
	default:
		ss_log.Error("审核操作异常,status[%v]", req.Status)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
