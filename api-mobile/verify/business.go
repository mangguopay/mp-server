package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
)

func InsertOrUpdateBusinessCardVerify(req *custProto.InsertOrUpdateBusinessCardRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ChannelId, "!=", "", ss_err.ERR_PARAM, "ChannelId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Name, "!=", "", ss_err.ERR_PARAM, "Name")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNumber, "!=", "", ss_err.ERR_PARAM, "CardNumber")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateBusinessCard 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func AddBusinessToHeadVerify(req *custProto.AddIndividualBusinessToHeadRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PayPwd, "!=", "", ss_err.ERR_PARAM, "PayPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ImageId, "!=", "", ss_err.ERR_PARAM, "ImageId")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddIndividualBusinessToHead 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
