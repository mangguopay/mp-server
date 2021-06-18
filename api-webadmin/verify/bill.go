package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	billProto "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
)

func CheckInsertHeadquartersProfitWithdrawVerify(req *billProto.InsertHeadquartersProfitWithdrawRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CurrencyType, "!=", "", ss_err.ERR_PARAM, "CurrencyType")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertHeadquartersProfitWithdraw 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UpdateToBusinessStatusVerify(req *billProto.UpdateToBusinessStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LogNo, "!=", "", ss_err.ERR_PARAM, "LogNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OrderStatus, "!=", "", ss_err.ERR_PARAM, "OrderStatus")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateToBusinessStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func AddChangeBalanceOrderVerify(req *billProto.AddChangeBalanceOrderRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccUid, "!=", "", ss_err.ERR_PARAM, "AccUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CurrencyType, "!=", "", ss_err.ERR_PARAM, "CurrencyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpType, "!=", "", ss_err.ERR_PARAM, "OpType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ChangeReason, "!=", "", ss_err.ERR_PARAM, "ChangeReason")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginPwd, "!=", "", ss_err.ERR_PARAM, "LoginPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddChangeBalanceOrder 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
