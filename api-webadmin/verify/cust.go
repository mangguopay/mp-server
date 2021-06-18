package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
)

func CheckUpdateLangRequestVerify(req *custProto.UpdateLangRequest) string {
	//fildChk := ss_chk.DoFieldChk(nil, req.Key, "!=", "", ss_err.ERR_PARAM, "Key")
	fildChk := ss_chk.DoFieldChk(nil, req.Type, "!=", "", ss_err.ERR_PARAM, "Type")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LangKm, "!=", "", ss_err.ERR_PARAM, "LangKm")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LangEn, "!=", "", ss_err.ERR_PARAM, "LangEn")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LangCh, "!=", "", ss_err.ERR_PARAM, "LangCh")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateLang 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UploadImageReqVerify(req *custProto.UploadImageRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ImageStr, "!=", "", ss_err.ERR_PARAM, "ImageStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UploadImage 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyServicerConfigReqVerify(req *custProto.ModifyServicerConfigRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ServicerNo, "!=", "", ss_err.ERR_PARAM, "ServicerNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IncomeAuthorization, "!=", "", ss_err.ERR_PARAM, "IncomeAuthorization")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OutgoAuthorization, "!=", "", ss_err.ERR_PARAM, "OutgoAuthorization")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CommissionSharing, "!=", "", ss_err.ERR_PARAM, "CommissionSharing")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IncomeSharing, "!=", "", ss_err.ERR_PARAM, "IncomeSharing")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UsdAuthCollectLimit, "!=", "", ss_err.ERR_PARAM, "UsdAuthCollectLimit")
	fildChk = ss_chk.DoFieldChk(fildChk, req.KhrAuthCollectLimit, "!=", "", ss_err.ERR_PARAM, "KhrAuthCollectLimit")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lat, "!=", "", ss_err.ERR_PARAM, "Lat")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lng, "!=", "", ss_err.ERR_PARAM, "Lng")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Scope, "!=", "", ss_err.ERR_PARAM, "Scope")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ScopeOff, "!=", "", ss_err.ERR_PARAM, "ScopeOff")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.ServicerName, "!=", "", ss_err.ERR_PARAM, "ServicerName")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessTime, "!=", "", ss_err.ERR_PARAM, "BusinessTime")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyServicerConfig 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UpdateFuncConfigReqVerify(req *custProto.UpdateFuncConfigRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.FuncName, "!=", "", ss_err.ERR_PARAM, "FuncName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.JumpUrl, "!=", "", ss_err.ERR_PARAM, "JumpUrl")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ImgBase64, "!=", "", ss_err.ERR_PARAM, "ImgBase64")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ApplicationType, "!=", "", ss_err.ERR_PARAM, "ApplicationType")

	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateFuncConfig 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InsertOrUpdateConsultationConfigVerify(req *custProto.InsertOrUpdateConsultationConfigRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Name, "!=", "", ss_err.ERR_PARAM, "Name")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Text, "!=", "", ss_err.ERR_PARAM, "Text")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UseStatus, "!=", "", ss_err.ERR_PARAM, "UseStatus")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LogoImgNo, "!=", "", ss_err.ERR_PARAM, "LogoImgNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateConsultationConfig 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InsertOrUpdateAgreementVerify(req *custProto.InsertOrUpdateAgreementRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Text, "!=", "", ss_err.ERR_PARAM, "Text")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Type, "!=", "", ss_err.ERR_PARAM, "Type")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UseStatus, "!=", "", ss_err.ERR_PARAM, "UseStatus")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")

	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateAgreement 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InsertOrUpdateChannelVerify(req *custProto.InsertOrUpdateChannelRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ChannelName, "!=", "", ss_err.ERR_PARAM, "ChannelName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LogoImgNo, "!=", "", ss_err.ERR_PARAM, "LogoImgNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LogoImgNoGrey, "!=", "", ss_err.ERR_PARAM, "LogoImgNoGrey")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ColorBegin, "!=", "", ss_err.ERR_PARAM, "ColorBegin")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ColorEnd, "!=", "", ss_err.ERR_PARAM, "ColorEnd")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateChannelRequest 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UpdateOrInsertCardReqVerify(req *custProto.UpdateOrInsertCardRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Name, "!=", "", ss_err.ERR_PARAM, "Name")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNumber, "!=", "", ss_err.ERR_PARAM, "CardNumber")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateOrInsertCard 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckInsertOrUpdateAppVersionVerify(req *custProto.InsertOrUpdateAppVersionRequest) string {
	//fildChk := ss_chk.DoFieldChk(nil, req.VId, "!=", "", ss_err.ERR_PARAM, "Addr")
	//fildChk := ss_chk.DoFieldChk(nil, req.VsType, "!=", "", ss_err.ERR_PARAM, "Nickname")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.System, "!=", "", ss_err.ERR_PARAM, "Phone")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.UpType, "!=", "", ss_err.ERR_PARAM, "Phone")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.Description, "!=", "", ss_err.ERR_PARAM, "Phone")
	//
	//fildChk = ss_chk.DoFieldChk(fildChk, req.Note, "!=", "", ss_err.ERR_PARAM, "Phone")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.VsCode, "!=", "", ss_err.ERR_PARAM, "Phone")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.IsForce, "!=", "", ss_err.ERR_PARAM, "Phone")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Phone")
	fildChk := ss_chk.DoFieldChk(nil, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateAppVersion 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckGetNewVersionVerify(req *custProto.GetNewVersionRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.System, "!=", "", ss_err.ERR_PARAM, "System")
	fildChk = ss_chk.DoFieldChk(fildChk, req.VsType, "!=", "", ss_err.ERR_PARAM, "VsType")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetNewVersion 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckDeleteCashierVerify(req *custProto.DeleteCashierRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.CashierNo, "!=", "", ss_err.ERR_PARAM, "CashierNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetNewVersion 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckModifyCashierVerify(req *custProto.ModifyCashierRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.CashierNo, "!=", "", ss_err.ERR_PARAM, "CashierNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PARAM, "Phone")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetNewVersion 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckGetCashiersVerify(req *custProto.GetCashiersRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ServicerNo, "!=", "", ss_err.ERR_PARAM, "ServicerNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetCashiers 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckInsertOrUpdateHelpVerify(req *custProto.InsertOrUpdateHelpRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Problem, "!=", "", ss_err.ERR_PARAM, "Problem")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Answer, "!=", "", ss_err.ERR_PARAM, "Answer")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lang, "!=", "", ss_err.ERR_PARAM, "Lang")
	fildChk = ss_chk.DoFieldChk(fildChk, req.VsType, "!=", "", ss_err.ERR_PARAM, "VsType")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateHelp 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckModifyCustInfoVerify(req *custProto.ModifyCustInfoRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.CustNo, "!=", "", ss_err.ERR_PARAM, "CustNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.InAuthorization, "!=", "", ss_err.ERR_PARAM, "InAuthorization")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OutAuthorization, "!=", "", ss_err.ERR_PARAM, "OutAuthorization")
	fildChk = ss_chk.DoFieldChk(fildChk, req.InTransferAuthorization, "!=", "", ss_err.ERR_PARAM, "InTransferAuthorization")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OutTransferAuthorization, "!=", "", ss_err.ERR_PARAM, "OutTransferAuthorization")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyCustInfo 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckModifyServicerStatusVerify(req *custProto.ModifyServicerStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ServicerNo, "!=", "", ss_err.ERR_PARAM, "ServicerNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UseStatus, "!=", "", ss_err.ERR_PARAM, "UseStatus")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyServicerStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InsertPosChannelVerify(req *custProto.InsertPosChannelRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ChannelNo, "!=", "", ss_err.ERR_PARAM, "ChannelNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IsRecom, "!=", "", ss_err.ERR_PARAM, "IsRecom")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CurrencyType, "!=", "", ss_err.ERR_PARAM, "CurrencyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertPosChannelVerify 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyAgreementStatusVerify(req *custProto.ModifyAgreementStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Id, "!=", "", ss_err.ERR_PARAM, "Id")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UseStatus, "!=", "", ss_err.ERR_PARAM, "UseStatus")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyAgreementStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func DeleteAgreementVerify(req *custProto.DeleteAgreementRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Id, "!=", "", ss_err.ERR_PARAM, "Id")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[DeleteAgreement 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyChannelPosStatusVerify(req *custProto.ModifyChannelPosStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Id, "!=", "", ss_err.ERR_PARAM, "Id")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UseStatus, "!=", "", ss_err.ERR_PARAM, "UseStatus")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyChannelPosStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func DeletePosChannelVerify(req *custProto.DeletePosChannelRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Id, "!=", "", ss_err.ERR_PARAM, "Id")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[DeletePosChannel 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyChannelPosIsRecomVerify(req *custProto.ModifyChannelPosIsRecomRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Id, "!=", "", ss_err.ERR_PARAM, "Id")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IsRecom, "!=", "", ss_err.ERR_PARAM, "IsRecom")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyChannelPosIsRecom 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyCollectStatusVerify(req *custProto.ModifyCollectStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.SetStatus, "!=", "", ss_err.ERR_PARAM, "SetStatus")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyCollectStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func DelectCardVerify(req *custProto.DelectCardRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[DelectCard 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InsertOrUpdatePushConfsVerify(req *custProto.InsertOrUpdatePushConfsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Pusher, "!=", "", ss_err.ERR_PARAM, "Pusher")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Config, "!=", "", ss_err.ERR_PARAM, "Config")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UseStatus, "!=", "", ss_err.ERR_PARAM, "UseStatus")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdatePushConfs 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func DeletePushConfVerify(req *custProto.DeletePushConfRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.PusherNo, "!=", "", ss_err.ERR_PARAM, "PusherNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[DeletePushConf 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InsertOrUpdatePushTempVerify(req *custProto.InsertOrUpdatePushTempRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.TempNo, "!=", "", ss_err.ERR_PARAM, "temp_no")
	fildChk = ss_chk.DoFieldChk(fildChk, req.PushNos, "!=", "", ss_err.ERR_PARAM, "push_nos")
	fildChk = ss_chk.DoFieldChk(fildChk, req.TitleKey, "!=", "", ss_err.ERR_PARAM, "title_key")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ContentKey, "!=", "", ss_err.ERR_PARAM, "content_key")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.LenArgs, "!=", "", ss_err.ERR_PARAM, "len_args")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdatePushTemp 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func InsertOrUpdateEventVerify(req *custProto.InsertOrUpdateEventRequest) string {
	//fildChk := ss_chk.DoFieldChk(nil, req.EventNo, "!=", "", ss_err.ERR_PARAM, "EventNo")
	fildChk := ss_chk.DoFieldChk(nil, req.EventName, "!=", "", ss_err.ERR_PARAM, "EventName")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateEvent 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func InsertOrUpdateEvaParamVerify(req *custProto.InsertOrUpdateEvaParamRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Key, "!=", "", ss_err.ERR_PARAM, "Key")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Val, "!=", "", ss_err.ERR_PARAM, "Val")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateEvaParam 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func InsertOrUpdateGlobalParamVerify(req *custProto.InsertOrUpdateGlobalParamRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ParamKey, "!=", "", ss_err.ERR_PARAM, "ParamKey")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ParamValue, "!=", "", ss_err.ERR_PARAM, "ParamValue")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateGlobalParam 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func InsertOrUpdateOpVerify(req *custProto.InsertOrUpdateOpRequest) string {
	//fildChk := ss_chk.DoFieldChk(nil, req.ParamKey, "!=", "", ss_err.ERR_PARAM, "ParamKey")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.ParamValue, "!=", "", ss_err.ERR_PARAM, "ParamValue")
	//if !fildChk.IsOk {
	//	ss_log.Error("err=[InsertOrUpdateGlobalParam 接口------> %s 为空]", fildChk.ErrReason)
	//	return fildChk.ErrCode
	//}
	return ""
}

func ModifyAuthMaterialBusinessStatusVerify(req *custProto.ModifyAuthMaterialBusinessStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AuthMaterialNo, "!=", "", ss_err.ERR_PARAM, "AuthMaterialNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyAuthMaterialBusinessStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func ModifyAuthMaterialStatusVerify(req *custProto.ModifyAuthMaterialStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AuthMaterialNo, "!=", "", ss_err.ERR_PARAM, "AuthMaterialNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyAuthMaterialBusinessStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyAuthMaterialEnterpriseStatusVerify(req *custProto.ModifyAuthMaterialEnterpriseStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AuthMaterialNo, "!=", "", ss_err.ERR_PARAM, "AuthMaterialNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AuthMaterialNo, "!=", "", ss_err.ERR_PARAM, "AuthMaterialNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyAuthMaterialEnterpriseStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func DeletePushTempVerify(req *custProto.DeletePushTempRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.TempNo, "!=", "", ss_err.ERR_PARAM, "TempNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[DeletePushTemp 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func ReStatisticVerify(req *custProto.ReStatisticRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Type, "!=", "", ss_err.ERR_PARAM, "Type")
	fildChk = ss_chk.DoFieldChk(fildChk, req.StartDate, "!=", "", ss_err.ERR_PARAM, "StartDate")
	fildChk = ss_chk.DoFieldChk(fildChk, req.EndDate, "!=", "", ss_err.ERR_PARAM, "EndDate")
	if !fildChk.IsOk {
		ss_log.Error("err=[ReStatistic 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func UpdateTransferSecurityConfigVerify(req *custProto.UpdateTransferSecurityConfigRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ContinuousErrPassword, "!=", "", ss_err.ERR_PARAM, "ContinuousErrPassword")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ErrPaymentPwdCount, "!=", "", ss_err.ERR_PARAM, "ErrPaymentPwdCount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateTransferSecurityConfig 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func UpdateWriteOffDurationDateConfigVerify(req *custProto.UpdateWriteOffDurationDateConfigRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.DurationDate, "!=", "", ss_err.ERR_PARAM, "DurationDate")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateWriteOffDurationDateConfig 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyBusinessStatusVerify(req *custProto.ModifyBusinessStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.BusinessNo, "!=", "", ss_err.ERR_PARAM, "BusinessNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	fildChk = ss_chk.DoFieldChk(fildChk, req.StatusType, "!=", "", ss_err.ERR_PARAM, "StatusType")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyBusinessStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func InsertOrUpdateBulletinVerify(req *custProto.InsertOrUpdateBulletinRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Title, "!=", "", ss_err.ERR_PARAM, "Title")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Content, "!=", "", ss_err.ERR_PARAM, "Content")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateBulletin 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func UpdateBulletinStatusVerify(req *custProto.UpdateBulletinStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.BulletinId, "!=", "", ss_err.ERR_PARAM, "BulletinId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	fildChk = ss_chk.DoFieldChk(fildChk, req.StatusType, "!=", "", ss_err.ERR_PARAM, "StatusType")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateBulletinStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func UpdateBusinessSignedStatusVerify(req *custProto.UpdateBusinessSignedStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.SignedId, "!=", "", ss_err.ERR_PARAM, "SignedId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateBusinessSignedCycle 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UpdateBusinessSignedInfoVerify(req *custProto.UpdateBusinessSignedInfoRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.SignedId, "!=", "", ss_err.ERR_PARAM, "SignedId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Cycle, "!=", "", ss_err.ERR_PARAM, "Cycle")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Rate, "!=", "", ss_err.ERR_PARAM, "Rate")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateBusinessSignedCycle 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModifyAuthMaterialBusinessUpdateStatusVerify(req *custProto.ModifyAuthMaterialBusinessUpdateStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Id, "!=", "", ss_err.ERR_PARAM, "Id")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModifyAuthMaterialBusinessUpdateStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func InsertOrUpdateBusinessIndustryRateCycleVerify(req *custProto.InsertOrUpdateBusinessIndustryRateCycleRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Code, "!=", "", ss_err.ERR_PARAM, "Code")
	fildChk = ss_chk.DoFieldChk(fildChk, req.BusinessChannelNo, "!=", "", ss_err.ERR_PARAM, "BusinessChannelNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Rate, "!=", "", ss_err.ERR_PARAM, "Rate")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Cycle, "!=", "", ss_err.ERR_PARAM, "Cycle")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateBusinessIndustryRateCycle 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func InsertOrUpdateBusinessIndustryVerify(req *custProto.InsertOrUpdateBusinessIndustryRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Code, "!=", "", ss_err.ERR_PARAM, "Code")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NameCh, "!=", "", ss_err.ERR_PARAM, "NameCh")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NameEn, "!=", "", ss_err.ERR_PARAM, "NameEn")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NameKm, "!=", "", ss_err.ERR_PARAM, "NameKm")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Level, "!=", "", ss_err.ERR_PARAM, "Level")
	if !fildChk.IsOk {
		ss_log.Error("err=[InsertOrUpdateBusinessIndustryRateCycle 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UpdateBusinessSceneSignedStatusVerify(req *custProto.UpdateBusinessSceneSignedStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.SignedId, "!=", "", ss_err.ERR_PARAM, "SignedId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Status, "!=", "", ss_err.ERR_PARAM, "Status")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateBusinessSignedCycle 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UpdateBusinessSceneSignedInfoVerify(req *custProto.UpdateBusinessSceneSignedInfoRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.SignedId, "!=", "", ss_err.ERR_PARAM, "SignedId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Cycle, "!=", "", ss_err.ERR_PARAM, "Cycle")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Rate, "!=", "", ss_err.ERR_PARAM, "Rate")
	fildChk = ss_chk.DoFieldChk(fildChk, req.LoginUid, "!=", "", ss_err.ERR_PARAM, "LoginUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UpdateBusinessSceneSignedInfo 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
