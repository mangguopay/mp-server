package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	billProto "a.a/mp-server/common/proto/bill"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
)

func BusinessWithdrawVerify(req *billProto.BusinessWithdrawRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IdenNo, "!=", "", ss_err.ERR_PARAM, "IdenNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.WithdrawType, "!=", "", ss_err.ERR_PARAM, "WithdrawType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PayPwd, "!=", "", ss_err.ERR_PARAM, "PayPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	if !fildChk.IsOk {
		ss_log.Error("err=[BusinessWithdraw 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func AddBusinessTransferVerify(req *billProto.AddBusinessTransferRequest) string {

	fildChk := ss_chk.DoFieldChk(nil, req.BusinessAccNo, "!=", "", ss_err.ERR_PARAM, "BusinessAccNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessNo, "!=", "", ss_err.ERR_PARAM, "BusinessNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PayeeNo, "!=", "", ss_err.ERR_PARAM, "PayeeNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PaymentPwd, "!=", "", ss_err.ERR_PARAM, "PaymentPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CurrencyType, "!=", "", ss_err.ERR_PARAM, "CurrencyType")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddBusinessTransfer 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func BusinessBillRefundVerify(req *businessBillProto.BusinessBillRefundRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.BusinessAccNo, "!=", "", ss_err.ERR_PARAM, "BusinessAccNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessNo, "!=", "", ss_err.ERR_PARAM, "BusinessNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RefundAmount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PaymentPwd, "!=", "", ss_err.ERR_PARAM, "PaymentPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddBusinessTransfer 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func BusinessBatchTransferConfirmVerify(req *billProto.BusinessBatchTransferConfirmRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.BusinessAccNo, "!=", "", ss_err.ERR_PARAM, "BusinessAccNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessNo, "!=", "", ss_err.ERR_PARAM, "BusinessNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PayPwd, "!=", "", ss_err.ERR_PARAM, "PayPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BatchNo, "!=", "", ss_err.ERR_PARAM, "BatchNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddBusinessTransfer 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func GetBatchAnalysisResult(req *billProto.GetBatchAnalysisResultRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessNo, "!=", "", ss_err.ERR_PARAM, "BusinessNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.FileId, "!=", "", ss_err.ERR_PARAM, "FileId")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetBatchAnalysisResult 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
