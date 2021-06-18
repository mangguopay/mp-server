package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
)

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
func PosFenceReqVerify(req *go_micro_srv_cust.PosFenceRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.PosSn, "!=", "", ss_err.ERR_PARAM, "pos_no")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lat, "!=", "", ss_err.ERR_PARAM, "lat")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Lng, "!=", "", ss_err.ERR_PARAM, "lng")
	if !fildChk.IsOk {
		ss_log.Error("err=[PosFence 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func UploadImageReqVerify(req *go_micro_srv_cust.UploadImageRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ImageStr, "!=", "", ss_err.ERR_PARAM, "ImageStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	if !fildChk.IsOk {
		ss_log.Error("err=[UploadImage 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func UnAuthDownloadImageReqVerify(req *go_micro_srv_cust.UnAuthDownloadImageRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ImageId, "!=", "", ss_err.ERR_PARAM, "ImageId")
	if !fildChk.IsOk {
		ss_log.Error("err=[UnAuthDownloadImage 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

/*
func DownloadImageReqVerify(req *go_micro_srv_cust.DownloadImageRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ImageId, "!=", "", ss_err.ERR_PARAM, "ImageId")
	if !fildChk.IsOk {
		ss_log.Error("err=[DownloadImage 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
*/

func GetHeadquartersCardsReqVerify(req *go_micro_srv_cust.GetHeadquartersCardsRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetHeadquartersCards 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func CheckModifyCashierVerify(req *go_micro_srv_cust.ModifyCashierRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.CashierNo, "!=", "", ss_err.ERR_PARAM, "System")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Phone, "!=", "", ss_err.ERR_PARAM, "System")
	if !fildChk.IsOk {
		ss_log.Error("err=[GetNewVersion 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func AddAuthMaterialBusinessVerify(req *go_micro_srv_cust.AddAuthMaterialBusinessRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AuthName, "!=", "", ss_err.ERR_PARAM, "AuthName")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.AuthNumber, "!=", "", ss_err.ERR_PARAM, "AuthNumber")
	fildChk = ss_chk.DoFieldChk(fildChk, req.TermType, "!=", "", ss_err.ERR_PARAM, "TermType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.SimplifyName, "!=", "", ss_err.ERR_PARAM, "SimplifyName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IndustryNo, "!=", "", ss_err.ERR_PARAM, "IndustryNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[AddAuthMaterialBusiness 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
