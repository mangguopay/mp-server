package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"context"
	"strings"
)

/**
 *平台盈利提现
 */
func (b *BillHandler) InsertHeadquartersProfitWithdraw(ctx context.Context, req *go_micro_srv_bill.InsertHeadquartersProfitWithdrawRequest, reply *go_micro_srv_bill.InsertHeadquartersProfitWithdrawReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	//查询平台盈利可提现余额是否满足此次提现
	amount, errGet1 := dao.HeadquartersProfitWithdrawDaoInst.GetProfitCashable(tx, req.CurrencyType)
	if errGet1 != ss_err.ERR_SUCCESS {
		ss_log.Error("查询平台盈利余额失败%s", errGet1)
		reply.ResultCode = errGet1
		return nil
	}

	// 判断金额是否包含小数点
	if req.CurrencyType == constants.CURRENCY_USD {
		if strings.Contains(req.Amount, ".") {
			reply.ResultCode = ss_err.ERR_PAY_AMOUNT_IS_NO_INTEGER
			return nil
		}
	}

	if strext.ToFloat64(amount) < strext.ToFloat64(req.Amount) {
		ss_log.Error("平台盈利可提现余额不足此次提现")
		reply.ResultCode = ss_err.ERR_PROFITC_ASHABLE_FAILD
		return nil
	}

	vaType := 0
	accPlatVaType := 0 //平台账户虚拟账户类型
	//更新提现账号的余额
	switch req.CurrencyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
		accPlatVaType = constants.VaType_USD_FEES
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
		accPlatVaType = constants.VaType_KHR_FEES
	default:
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//添加平台盈利提现日志
	logNo := strext.GetDailyId()
	sqlInsert := "insert into headquarters_profit_withdraw(order_no, currency_type, amount, note, create_time, account_no) " +
		" values($1,$2,$3,$4,current_timestamp,$5)"
	errAddLog := ss_sql.ExecTx(tx, sqlInsert, logNo, req.CurrencyType, req.Amount, req.Note, req.AccountNo)
	if errAddLog != nil {
		ss_log.Error("err=[%v]", errAddLog)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	//查询虚拟账户的vid
	vaccNo := InternalCallHandlerInst.ConfirmExistVAccount(req.AccountNo, req.CurrencyType, strext.ToInt32(vaType))
	//更新提现的账号虚账余额
	err := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, vaccNo, req.Amount, "+", logNo, constants.VaReason_PROFIT_OUTGO)
	if err != ss_err.ERR_SUCCESS {
		ss_log.Error("平台盈利提现 更新账户[%v]余额失败,err=[%v]", req.AccountNo, err)
		reply.ResultCode = err
		return nil
	}

	//更新平台账号的虚帐余额
	_, accPlat, _ := cache.ApiDaoInstance.GetGlobalParam("acc_plat")
	headVaccNo := InternalCallHandlerInst.ConfirmExistVAccount(accPlat, req.CurrencyType, strext.ToInt32(accPlatVaType))
	//更新提现的账号虚账余额
	errUp := dao.VaccountDaoInst.ModifyVaccRemainUpperZero(tx, headVaccNo, req.Amount, "-", logNo, constants.VaReason_PROFIT_OUTGO)
	if errUp != ss_err.ERR_SUCCESS {
		ss_log.Error("平台盈利提现 更新平台账户[%v]虚帐[%v]余额失败,err=[%v]", accPlat, headVaccNo, errUp)
		reply.ResultCode = err
		return nil
	}

	//更新平台可提现余额
	err2 := dao.HeadquartersProfitWithdrawDaoInst.ModifyProfitCashable(tx, req.Amount, req.CurrencyType)
	if err2 != ss_err.ERR_SUCCESS {
		ss_log.Error("更新平台可提现余额失败，err2=[%v]", err2)
		reply.ResultCode = err2
		return nil
	}

	//更新账号account表中的余额
	//if errStr := dao.VaccountDaoInst.SyncAccRemain(tx, req.AccountNo); errStr != ss_err.ERR_SUCCESS {
	//	ss_log.Error("同步账号[%v]余额失败。errStr=[%v]", req.AccountNo, errStr)
	//	reply.ResultCode = errStr
	//	return nil
	//}

	tx.Commit()
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
