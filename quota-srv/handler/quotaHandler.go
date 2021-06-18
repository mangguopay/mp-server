package handler

import (
	"context"
	"database/sql"

	"a.a/cu/db"
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	ss_sql2 "a.a/cu/ss_sql"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	go_micro_srv_quota "a.a/mp-server/common/proto/quota"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/common/ss_struct"
	"a.a/mp-server/quota-srv/common"
	"a.a/mp-server/quota-srv/dao"
)

type QuotaHandler struct{}

var QuotaHandlerInst QuotaHandler

//
func (b QuotaHandler) ModifyQuota(ctx context.Context, req *go_micro_srv_quota.ModifyQuotaRequest, resp *go_micro_srv_quota.ModifyQuotaReply) error {
	// 参数校验
	fildChk := ss_chk.DoFieldChk(nil, req.CurrencyType, "!=", "", ss_err.ERR_PARAM, "CurrencyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpType, "!=", "", ss_err.ERR_PARAM, "OpType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LogNo, "!=", "", ss_err.ERR_PARAM, "LogNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyQuota 接口------> 参数 %s 为空]", fildChk.ErrReason)
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	switch req.OpType {
	case constants.QuotaOp_PreSave:
		//errCode := b.preSave(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		errCode := b.preSaveV2(req.TxNo, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_Save:
		//errCode := b.save(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		errCode := b.saveV2(req.TxNo, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_Withdraw:
		errCode := b.withdraw(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_SvrPreSave:
	case constants.QuotaOp_SvrSave:
		errCode := b.svrSave(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_SvrWithdraw:
		errCode := b.svrWithdraw(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_SvrPreWithdraw:
		errCode := b.svrPreWithdraw(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_SvrWithdraw_Cancel:
		errCode := b.svrWithdrawCancel(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_CustPreSave:
		errCode := b.custPreSave(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_CustSave:
		errCode := b.custSave(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_CustSave_Cancel:
		errCode := b.custSaveCancel(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_Rollback:
		errCode := b.rollback(req.AccountNo, req.CurrencyType, req.Amount, req.LogNo, req.Reason)
		//errCode := b.rollbackV2(req.TxNo, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo, req.Reason)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_SvrCashRecharge:
		errCode := b.svrCashRecharge(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_BusinessSave:
		errCode := b.businessSave(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_ChangeCustBalanceAdd:
		errCode := b.changeCustBalance(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo, constants.VaOpType_Add)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_ChangeCustBalanceMinus:
		errCode := b.changeCustBalance(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo, constants.VaOpType_Minus)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_ChangeSvrBalanceAdd:
		errCode := b.changeSvrBalance(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo, constants.VaOpType_Add)
		resp.ResultCode = errCode
		return nil
	case constants.QuotaOp_ChangeSvrBalanceMinus:
		errCode := b.changeSvrBalance(ctx, req.AccountNo, req.CurrencyType, req.Amount, req.LogNo, constants.VaOpType_Minus)
		resp.ResultCode = errCode
		return nil
	default:
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	return nil
}

//
func (b QuotaHandler) ModifyDefaultQuota(ctx context.Context, req *go_micro_srv_quota.ModifyDefaultQuotaRequest, resp *go_micro_srv_quota.ModifyDefaultQuotaReply) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	switch req.OpType {
	case constants.VaOpType_Add:
		fallthrough
	case constants.VaOpType_Minus:
		fallthrough
	case constants.VaOpType_Freeze:
		fallthrough
	case constants.VaOpType_Defreeze_Add:
		var vaType int32 = 0
		switch req.CurrencyType {
		case "usd":
			vaType = constants.VaType_QUOTA_USD
		case "khr":
			vaType = constants.VaType_QUOTA_KHR
		}
		recvVaccNo := b.confirmExistVaccountTx(tx, req.AccountNo, vaType)
		errCode := dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, req.Amount, req.OpType, req.LogNo, "reason") // todo reason 需要确定
		if errCode != ss_err.ERR_SUCCESS {
			ss_log.Error("errCode=[%v]", errCode)
			resp.ResultCode = errCode
			return nil
		}
	default:
		resp.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	ss_sql.Commit(tx)
	resp.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (QuotaHandler) confirmExistVaccountTx(tx *sql.Tx, accountNo string, vaType int32) (vaccountNo string) {
	vaccountNo = dao.VaccountDaoInst.GetVaccountNoTx(tx, accountNo, vaType)
	if vaccountNo == "" {
		vaccountNo = dao.VaccountDaoInst.InitVaccountNoTx(tx, accountNo, vaType)
		return vaccountNo
	}
	return vaccountNo
}

func (QuotaHandler) confirmExistVaccountTxV2(tmProxy *ss_struct.TmServerProxy, accountNo string, vaType int32) (string, error) {
	vaccountNo, err := dao.VaccountDaoInst.GetVaccountNoTxV2(tmProxy, accountNo, vaType)
	if err != nil {
		return "", err
	}
	if vaccountNo == "" {
		vaccountNo, err = dao.VaccountDaoInst.InitVaccountNoTxV2(tmProxy, accountNo, vaType)
		if err != nil {
			return "", err
		}
		return vaccountNo, nil
	}
	return vaccountNo, nil
}

func (b QuotaHandler) preSave(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	dbHandler.SetConnMaxLifetime(-1)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}

	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)
	if recvVaccNo == "" {
		return ss_err.ERR_PARAM
	}
	// 修改实时额度
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Defreeze_But_Minus, logNo, constants.VaReason_INCOME)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	balance, fbalance := dao.VaccountDaoInst.GetBalanceTx(tx, recvVaccNo)
	// 规定额度
	var vaTypeFixed int32 = 0
	switch currencyType {
	case "usd":
		vaTypeFixed = constants.VaType_QUOTA_USD
	case "khr":
		vaTypeFixed = constants.VaType_QUOTA_KHR
	}
	recvVaccNoFixed := b.confirmExistVaccountTx(tx, accountNo, vaTypeFixed)
	balanceFixed, _ := dao.VaccountDaoInst.GetBalanceTx(tx, recvVaccNoFixed)

	useAmount := ss_count.Sub(balanceFixed, balance).String()    // 固定-实时=可用额度
	resultAmount, _ := ss_count.Sub(useAmount, amount).Float64() // 可用额度-存款金额<0,余额不足
	if resultAmount < 0 {
		ss_log.Error("err=[服务商的可用额度为----->%s,固定额度为----->%s, 实时额度为----->%s]", useAmount, balanceFixed, balance)
		return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	}
	if strext.ToInt64(balance)+strext.ToInt64(fbalance) > strext.ToInt64(balanceFixed) {
		return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) preSaveV2(txNo, accountNo, currencyType, amount, logNo string) (errCode string) {

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}

	tmProxy, err := ss_struct.GetTmServerProxyFromTxNo(common.QuotaServerFullid, txNo)
	if err != nil {
		ss_log.Error("GetTmServerProxyFromTxNo err: %s", err.Error())
		return ss_err.ERR_PARAM
	}

	recvVaccNo, err := b.confirmExistVaccountTxV2(tmProxy, accountNo, vaType)
	if err != nil {
		ss_log.Error("confirmExistVaccountTxV2 err: %s", err.Error())
		return ss_err.ERR_PARAM
	}
	if recvVaccNo == "" {
		ss_log.Error("confirmExistVaccountTxV2 recvVaccNo为空")
		return ss_err.ERR_PARAM
	}
	// 修改实时额度
	err2 := dao.VaccountDaoInst.ModifyVaccV2(tmProxy, recvVaccNo, amount, constants.VaOpType_Defreeze_But_Minus, logNo, constants.VaReason_INCOME)
	if err2 != nil {
		ss_log.Error("ModifyVaccV2 err=[%v]", err2)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	balance, fbalance, err3 := dao.VaccountDaoInst.GetBalanceV2(tmProxy, recvVaccNo)
	if err3 != nil {
		ss_log.Error("GetBalanceV2 err=[%v]", err2)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}
	// 规定额度
	var vaTypeFixed int32 = 0
	switch currencyType {
	case "usd":
		vaTypeFixed = constants.VaType_QUOTA_USD
	case "khr":
		vaTypeFixed = constants.VaType_QUOTA_KHR
	}
	recvVaccNoFixed, err4 := b.confirmExistVaccountTxV2(tmProxy, accountNo, vaTypeFixed)
	if err4 != nil {
		ss_log.Error("confirmExistVaccountTxV2 err=[%v]", err2)
		return ss_err.ERR_PARAM
	}
	balanceFixed, _, err5 := dao.VaccountDaoInst.GetBalanceV2(tmProxy, recvVaccNoFixed)
	if err5 != nil {
		ss_log.Error("GetBalanceV2 err=[%v]", err2)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	useAmount := ss_count.Sub(balanceFixed, balance).String()    // 固定-实时=可用额度
	resultAmount, _ := ss_count.Sub(useAmount, amount).Float64() // 可用额度-存款金额<0,余额不足
	if resultAmount < 0 {
		ss_log.Error("err=[服务商的可用额度为----->%s,固定额度为----->%s, 实时额度为----->%s]", useAmount, balanceFixed, balance)
		return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	}
	if strext.ToInt64(balance)+strext.ToInt64(fbalance) > strext.ToInt64(balanceFixed) {
		ss_log.Error("err=[金额加冻结金额大于固定额度,balance(%s)+fbalance(%s)= %v, balanceFixed=%v, recvVaccNo:%v]",
			balance, fbalance,
			strext.ToInt64(balance)+strext.ToInt64(fbalance),
			strext.ToInt64(balanceFixed),
			recvVaccNo,
		)
		return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) saveV2(txNo, accountNo, currencyType, amount, logNo string) (errCode string) {
	tmProxy, err := ss_struct.GetTmServerProxyFromTxNo(common.QuotaServerFullid, txNo)
	if err != nil {
		ss_log.Error("GetTmServerProxyFromTxNo err: %s", err.Error())
		return ss_err.ERR_PARAM
	}

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo, err1 := b.confirmExistVaccountTxV2(tmProxy, accountNo, vaType)
	if err1 != nil {
		ss_log.Error("confirmExistVaccountTxV2 err=[%v]", err1)
		return ss_err.ERR_PARAM
	}
	// 冻结金额可能为负数,当客户向总部充值时,冻结金额就是负数
	//_, fbalance, err2 := dao.VaccountDaoInst.GetBalanceV2(tmProxy, recvVaccNo)
	//if err2 != nil {
	//	ss_log.Error("GetBalanceV2 err=[%v]", err2)
	//	return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	//}
	//if strext.ToInt64(fbalance) < strext.ToInt64(amount) {
	//	ss_log.Error("客户存款,冻结金额小于金额,strext.ToInt64(fbalance) < strext.ToInt64(amount),fbalance: %s,amount: %s", fbalance, amount)
	//	return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	//}

	err3 := dao.VaccountDaoInst.ModifyVaccV2(tmProxy, recvVaccNo, amount, constants.VaOpType_Defreeze_Add, logNo, constants.VaReason_INCOME)
	if err3 != nil {
		ss_log.Error("ModifyVaccV2 err=[%v]", err3)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}
func (b QuotaHandler) save(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)
	_, fbalance := dao.VaccountDaoInst.GetBalanceTx(tx, recvVaccNo)
	if strext.ToInt64(fbalance) < strext.ToInt64(amount) {
		return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	}

	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Defreeze_Add, logNo, constants.VaReason_INCOME)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) withdraw(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)
	//_, fbalance := dao.VaccountDaoInst.GetBalanceTx(tx, recvVaccNo)
	//if strext.ToInt64(fbalance) < strext.ToInt64(amount) {
	//	return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	//}

	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Minus, logNo, constants.VaReason_OUTGO)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) custPreSave(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Defreeze_But_Minus, logNo, constants.VaReason_Cust_Save)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) svrSave(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	famount := ss_count.Sub("0", amount).String()
	//errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, famount, constants.VaOpType_Defreeze_Add, logNo, constants.VaReason_Srv_Save)
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, famount, constants.VaOpType_Add, logNo, constants.VaReason_Srv_Save)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}
func (b QuotaHandler) custSave(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	//famount := ss_count.Sub("0", amount).String()
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Defreeze_Add, logNo, constants.VaReason_Cust_Save)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}
func (b QuotaHandler) custSaveCancel(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_USD_DEBIT
	case "khr":
		vaType = constants.VaType_KHR_DEBIT
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	//famount := ss_count.Sub("0", amount).String()
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Defreeze_Minus, logNo, constants.VaReason_Cust_Save)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) svrWithdraw(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)
	_, fbalance := dao.VaccountDaoInst.GetBalanceTx(tx, recvVaccNo)
	if strext.ToInt64(fbalance) < strext.ToInt64(amount) {
		return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	}

	//冻结-，实时不变
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Defreeze, logNo, constants.VaReason_Srv_Withdraw)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) rollbackV2(txNo, accountNo, currencyType, amount, logNo, reason string) (errCode string) {
	tmProxy, err := ss_struct.GetTmServerProxyFromTxNo(common.QuotaServerFullid, txNo)
	if err != nil {
		ss_log.Error("GetTmServerProxyFromTxNo err: %s", err.Error())
		return ss_err.ERR_PARAM
	}

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo, err1 := b.confirmExistVaccountTxV2(tmProxy, accountNo, vaType)
	if err1 != nil {
		ss_log.Error("confirmExistVaccountTxV2 err: %s", err1.Error())
		return ss_err.ERR_PARAM
	}

	err2 := dao.VaccountDaoInst.ModifyVaccV2(tmProxy, recvVaccNo, amount, constants.VaOpType_Defreeze_Minus, logNo, reason)
	if err2 != nil {
		ss_log.Error("ModifyVaccV2 err=[%v]", err2)
		return ss_err.ERR_PAY_AMT_NOT_ENOUGH
	}

	return ss_err.ERR_SUCCESS
}
func (b QuotaHandler) rollback(accountNo, currencyType, amount, logNo, reason string) (errCode string) {
	dbHandler := ss_sql2.NewTxInst(constants.DB_CRM)
	defer dbHandler.Rollback()

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo := b.confirmExistVaccountTx(dbHandler.Tx, accountNo, vaType)

	errCode = dao.VaccountDaoInst.ModifyVacc(dbHandler.Tx, recvVaccNo, amount, constants.VaOpType_Defreeze_Minus, logNo, reason)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	dbHandler.Commit()
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) svrWithdrawCancel(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)
	//_, fbalance := dao.VaccountDaoInst.GetBalanceTx(tx, recvVaccNo)
	//if strext.ToInt64(fbalance) < strext.ToInt64(amount) {
	//	return ss_err.ERR_PAY_QUOTA_NOT_ENOUGH
	//}

	//冻结-，实时+
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Balance_Defreeze_Add, logNo, constants.VaReason_Srv_Withdraw)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) svrPreWithdraw(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case "usd":
		vaType = constants.VaType_QUOTA_USD_REAL
	case "khr":
		vaType = constants.VaType_QUOTA_KHR_REAL
	}
	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	//冻结+
	//fAmount := ss_count.Sub("0", amount).String()
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Balance_Frozen_Add, logNo, constants.VaReason_Srv_Withdraw)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}

	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) svrCashRecharge(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_QUOTA_USD_REAL
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_QUOTA_KHR_REAL
	}

	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	amountB := ss_count.Sub("0", amount).String()

	//余额增加
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amountB, constants.VaOpType_Add, logNo, constants.VaReason_Srv_CashRecharge)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}
	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) businessSave(ctx context.Context, accountNo, currencyType, amount, logNo string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_BUSINESS_SETTLED
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_BUSINESS_SETTLED
	}

	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	//余额增加
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, constants.VaOpType_Add, logNo, constants.VaReason_Business_Save)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}
	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) changeCustBalance(ctx context.Context, accountNo, currencyType, amount, logNo, op string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_USD_DEBIT
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_KHR_DEBIT
	}

	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)

	//余额增加
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amount, op, logNo, constants.VaReason_ChangeCustBalance)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}
	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}

func (b QuotaHandler) changeSvrBalance(ctx context.Context, accountNo, currencyType, amount, logNo, op string) (errCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.BeginTx(ctx, nil)
	defer ss_sql.Rollback(tx)

	var vaType int32 = 0
	switch currencyType {
	case constants.CURRENCY_USD:
		vaType = constants.VaType_QUOTA_USD_REAL
	case constants.CURRENCY_KHR:
		vaType = constants.VaType_QUOTA_KHR_REAL
	}

	recvVaccNo := b.confirmExistVaccountTx(tx, accountNo, vaType)
	amountB := ss_count.Sub("0", amount).String()

	//余额增加
	errCode = dao.VaccountDaoInst.ModifyVacc(tx, recvVaccNo, amountB, op, logNo, constants.VaReason_ChangeSrvBalance)
	if errCode != ss_err.ERR_SUCCESS {
		ss_log.Error("errCode=[%v]", errCode)
		return errCode
	}
	ss_sql.Commit(tx)
	return ss_err.ERR_SUCCESS
}
