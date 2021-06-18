package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	authProto "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_err"
)

func CheckAddCashierVerify(req *authProto.AddCashierRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ServicerNo, "!=", "", ss_err.ERR_PARAM, "ServicerNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CountryCode, "!=", "", ss_err.ERR_PARAM, "CountryCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddCashier 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckModifyUserStatusVerify(req *adminAuthProto.ModifyUserStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.SetStatus, "!=", "", ss_err.ERR_PARAM, "SetStatus")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyUserStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
