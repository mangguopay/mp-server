package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/ss_err"
)

func ModifyPWDReqVerify(req *go_micro_srv_auth.MobileModifyPwdRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OldPassword, "!=", "", ss_err.ERR_PARAM, "OldPassword")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NewPassword, "!=", "", ss_err.ERR_PARAM, "NewPassword")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPWDReq 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyPayPWDReqVerify(req *go_micro_srv_auth.MobileModifyPayPwdRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPayPWD 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

//
func ModifyPayPWDByMailCodeVerify(req *go_micro_srv_auth.ModifyPayPWDByMailCodeRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MailCode, "!=", "", ss_err.ERR_PARAM, "MailCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPayPWDByMailCode 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

//
func ModifyEmailVerify(req *go_micro_srv_auth.ModifyEmailRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MailCode, "!=", "", ss_err.ERR_PARAM, "MailCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Email, "!=", "", ss_err.ERR_PARAM, "Email")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Account, "!=", "", ss_err.ERR_PARAM, "Account")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IdenNo, "!=", "", ss_err.ERR_PARAM, "IdenNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CountryCode, "!=", "", ss_err.ERR_PARAM, "CountryCode")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyEmail 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

//
func BusinessModifyPhoneRequestVerify(req *go_micro_srv_auth.BusinessModifyPhoneRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CountryCode, "!=", "", ss_err.ERR_PARAM, "CountryCode")
	if !fildChk.IsOk {
		ss_log.Error("err=[BusinessModifyPhone 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

//
func ModifyPayPWDByOldPwdVerify(req *go_micro_srv_auth.ModifyPayPWDByOldPwdRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.IdenNo, "!=", "", ss_err.ERR_PARAM, "IdenNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OldPayPwd, "!=", "", ss_err.ERR_PARAM, "OldPayPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NewPayPwd, "!=", "", ss_err.ERR_PARAM, "NewPayPwd")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPayPWDByOldPwd 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func BusinessModifyPWDBySmsVerify(req *go_micro_srv_auth.BusinessModifyPWDBySmsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	if !fildChk.IsOk {
		ss_log.Error("err=[BusinessModifyPWDBySms 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func BusinessNotokenModifyPWDBySmsVerify(req *go_micro_srv_auth.BusinessModifyPWDBySmsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Account, "!=", "", ss_err.ERR_PARAM, "Account")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	if !fildChk.IsOk {
		ss_log.Error("err=[BusinessModifyPWDBySms 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyPWDByMailCodeVerify(req *go_micro_srv_auth.ModifyPWDByMailCodeRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MailCode, "!=", "", ss_err.ERR_PARAM, "MailCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Email, "!=", "", ss_err.ERR_PARAM, "Email")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPWDByMailCode 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func NoTokenModifyPWDByMailCodeVerify(req *go_micro_srv_auth.ModifyPWDByMailCodeRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MailCode, "!=", "", ss_err.ERR_PARAM, "MailCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Email, "!=", "", ss_err.ERR_PARAM, "Email")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Account, "!=", "", ss_err.ERR_PARAM, "Account")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPWDByMailCode 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InitBusinessPayPwdVerify(req *go_micro_srv_auth.InitBusinessPayPwdRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IdenNo, "!=", "", ss_err.ERR_PARAM, "IdenNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PayPwd, "!=", "", ss_err.ERR_PARAM, "PayPwd")
	if !fildChk.IsOk {
		ss_log.Error("err=[InitBusinessPayPwd 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func SaveBusinessAccountVerify(req *go_micro_srv_auth.SaveBusinessAccountRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Email, "!=", "", ss_err.ERR_PARAM, "Email")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PhoneISNull_FAILD, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CountryCode, "!=", "", ss_err.ERR_PARAM, "CountryCode")
	if !fildChk.IsOk {
		ss_log.Error("err=[SaveBusinessAccount 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
