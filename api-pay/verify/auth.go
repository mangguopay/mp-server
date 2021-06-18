package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
)

func MobileLoginReqVerify(req *go_micro_srv_auth.MobileLoginRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Account, "!=", "", ss_err.ERR_PARAM, "Account")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.Imei, "!=", "", ss_err.ERR_PARAM, "Imei")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Nonstr, "!=", "", ss_err.ERR_PARAM, "Nonstr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	if !fildChk.IsOk {
		ss_log.Error("err=[MobileLogin 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func MobileRegReqVerify(req *go_micro_srv_auth.MobileRegRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	if !fildChk.IsOk {
		ss_log.Error("err=[MobileReg 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func SmsReqVerify(req *go_micro_srv_cust.RegSmsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Function, "!=", "", ss_err.ERR_PARAM, "Function")
	if !fildChk.IsOk {
		ss_log.Error("err=[RegSms 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func BackPWDReqVerify(req *go_micro_srv_auth.MobileBackPwdRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	if !fildChk.IsOk {
		ss_log.Error("err=[BackPWD 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func RegLoginReqVerify(req *go_micro_srv_auth.RegLoginRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	if !fildChk.IsOk {
		ss_log.Error("err=[RegLogin 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func ModifyPhoneReqVerify(req *go_micro_srv_auth.MobileModifyPhonedRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPhone 接口------> %s 为空]", fildChk.ErrReason)
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
func ModifyNicknameReqVerify(req *go_micro_srv_auth.ModifyNicknameRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Nickname, "!=", "", ss_err.ERR_PARAM, "Nickname")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyNickname 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func ModifyPWDReqVerify(req *go_micro_srv_auth.MobileModifyPwdRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Uid, "!=", "", ss_err.ERR_PARAM, "Uid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OldPassword, "!=", "", ss_err.ERR_PARAM, "OldPassword")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NewPassword, "!=", "", ss_err.ERR_PARAM, "NewPassword")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyPWD 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func ModifyDefaultCardReqVerify(req *go_micro_srv_auth.ModifyDefaultCardRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpAccNo, "!=", "", ss_err.ERR_PARAM, "OpAccNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyDefaultCard 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func AddCardReqVerify(req *go_micro_srv_auth.AddCardRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ChannelName, "!=", "", ss_err.ERR_PARAM, "ChannelName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecName, "!=", "", ss_err.ERR_PARAM, "RecName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpAccNo, "!=", "", ss_err.ERR_PARAM, "OpAccNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecCarNum, "!=", "", ss_err.ERR_PARAM, "RecCarNum")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddCardReqVerify 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func DeleteBindCardReqVerify(req *go_micro_srv_auth.DeleteBindCardRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ChannelName, "!=", "", ss_err.ERR_PARAM, "ChannelName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CarNum, "!=", "", ss_err.ERR_PARAM, "CarNum")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpAccNo, "!=", "", ss_err.ERR_PARAM, "OpAccNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[DeleteBindCard 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func PerfectingInfoReqVerify(req *go_micro_srv_auth.PerfectingInfoRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.Nickname, "!=", "", ss_err.ERR_PARAM, "Nickname")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.ImageStr, "!=", "", ss_err.ERR_PARAM, "ImageStr")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[PerfectingInfo 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func CheckPayPWDReqVerify(req *go_micro_srv_auth.CheckPayPWDRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IdenNo, "!=", "", ss_err.ERR_PARAM, "IdenNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	if !fildChk.IsOk {
		ss_log.Error("err=[CheckPayPWD 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func GetPosRemainReqVerify(req *go_micro_srv_auth.GetPosRemainRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetPosRemain 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func VersionInfoReqVerify(req *go_micro_srv_auth.GetVersinInfoRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AppVersion, "!=", "", ss_err.ERR_PARAM, "AppVersion")
	if !fildChk.IsOk {
		ss_log.Error("err=[VersionInfo 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
