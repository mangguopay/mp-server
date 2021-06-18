package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
)

func ModifyAuthMaterialStatusVerify(req *go_micro_srv_cust.ModifyAuthMaterialStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AuthMaterialNo, "!=", "", ss_err.ERR_PARAM, "AuthMaterialNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AuthMaterialNo, "!=", "", ss_err.ERR_PARAM, "AuthMaterialNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyAuthMaterialStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func AddAuthMaterialEnterpriseVerify(req *go_micro_srv_cust.AddAuthMaterialEnterpriseRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ImgBase64, "!=", "", ss_err.ERR_PARAM, "ImgBase64")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AuthName, "!=", "", ss_err.ERR_PARAM, "AuthName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AuthNumber, "!=", "", ss_err.ERR_PARAM, "AuthNumber")
	fildChk = ss_chk.DoFieldChk(fildChk, req.TermType, "!=", "", ss_err.ERR_PARAM, "TermType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.SimplifyName, "!=", "", ss_err.ERR_PARAM, "SimplifyName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddAuthMaterialEnterprise 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func UpdateBusinessAuthMaterialInfoVerify(req *go_micro_srv_cust.UpdateBusinessAuthMaterialInfoRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.SimplifyName, "!=", "", ss_err.ERR_PARAM, "SimplifyName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateBusinessAuthMaterialInfo 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func UpdateBusinessBaseInfoVerify(req *go_micro_srv_cust.UpdateBusinessBaseInfoRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MainIndustry, "!=", "", ss_err.ERR_PARAM, "MainIndustry")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MainBusiness, "!=", "", ss_err.ERR_PARAM, "MainBusiness")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ContactPerson, "!=", "", ss_err.ERR_PARAM, "ContactPerson")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ContactPhone, "!=", "", ss_err.ERR_PARAM, "ContactPhone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CountryCode, "!=", "", ss_err.ERR_PARAM, "CountryCode")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateBusinessBaseInfo 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckSmsReqVerify(req *go_micro_srv_cust.CheckSmsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Sms, "!=", "", ss_err.ERR_PARAM, "Sms")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Function, "!=", "", ss_err.ERR_PARAM, "Function")
	if !fildChk.IsOk {
		ss_log.Error("err=[CheckSms 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func SmsReqVerify(req *go_micro_srv_cust.RegSmsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Function, "!=", "", ss_err.ERR_PARAM, "Function")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CountryCode, "!=", "", ss_err.ERR_PARAM, "country_code")
	if !fildChk.IsOk {
		ss_log.Error("err=[RegSms 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func InsertOrUpdateBusinessCardVerify(req *go_micro_srv_cust.InsertOrUpdateBusinessCardRequest) string {
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
func AddBusinessToHeadVerify(req *go_micro_srv_cust.AddBusinessToHeadRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PayPwd, "!=", "", ss_err.ERR_PARAM, "PayPwd")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ImageBase64, "!=", "", ss_err.ERR_PARAM, "ImageBase64")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddBusinessToHead 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func AddBusinessAppVerify(req *go_micro_srv_cust.InsertOrUpdateBusinessAppRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.IdenNo, "!=", "", ss_err.ERR_PARAM, "IdenNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ApplyType, "!=", "", ss_err.ERR_PARAM, "ApplyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AppName, "!=", "", ss_err.ERR_PARAM, "AppName")

	fildChk = ss_chk.DoFieldChk(fildChk, req.SmallImgNo, "!=", "", ss_err.ERR_PARAM, "SmallImgNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BigImgNo, "!=", "", ss_err.ERR_PARAM, "BigImgNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddBusinessApp 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func GetHeadquartersCardsReqVerify(req *go_micro_srv_cust.GetHeadquartersCardsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetHeadquartersCards 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func SendMailReqVerify(req *go_micro_srv_cust.SendMailRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Function, "!=", "", ss_err.ERR_PARAM, "Function")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Email, "!=", "", ss_err.ERR_PARAM, "Email")
	if !fildChk.IsOk {
		ss_log.Error("err=[SendMail 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func IdentityVerifyVerify(req *go_micro_srv_cust.IdentityVerifyRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Account, "!=", "", ss_err.ERR_PARAM, "Account")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Verifyid, "!=", "", ss_err.ERR_PARAM, "Verifyid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Verifynum, "!=", "", ss_err.ERR_PARAM, "Verifynum")
	if !fildChk.IsOk {
		ss_log.Error("err=[IdentityVerify 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckMailCodeVerify(req *go_micro_srv_cust.CheckMailCodeRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Mail, "!=", "", ss_err.ERR_PARAM, "Mail")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MailCode, "!=", "", ss_err.ERR_PARAM, "MailCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Function, "!=", "", ss_err.ERR_PARAM, "Function")
	if !fildChk.IsOk {
		ss_log.Error("err=[CheckMailCode 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func AddBusinessSignedVerify(req *go_micro_srv_cust.AddBusinessSignedRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccUid, "!=", "", ss_err.ERR_PARAM, "AccUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessNo, "!=", "", ss_err.ERR_PARAM, "BusinessNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AppId, "!=", "", ss_err.ERR_PARAM, "AppId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.SceneNo, "!=", "", ss_err.ERR_PARAM, "SceneNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IndustryNo, "!=", "", ss_err.ERR_PARAM, "IndustryNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddBusinessSigned 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func AddBusinessSceneSignedVerify(req *go_micro_srv_cust.AddBusinessSceneSignedRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccUid, "!=", "", ss_err.ERR_PARAM, "AccUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessNo, "!=", "", ss_err.ERR_PARAM, "BusinessNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.SceneNo, "!=", "", ss_err.ERR_PARAM, "SceneNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IndustryNo, "!=", "", ss_err.ERR_PARAM, "IndustryNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddBusinessSceneSigned 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
