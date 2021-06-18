package verify

import (
	"a.a/cu/ss_chk"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_err"
)

func ExchangeReqVerify(req *go_micro_srv_bill.ExchangeRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.InType, "!=", "", ss_err.ERR_PARAM, "InType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OutType, "!=", "", ss_err.ERR_PARAM, "OutType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.TransFrom, "!=", "", ss_err.ERR_PARAM, "TransFrom")
	if !fildChk.IsOk {
		ss_log.Error("err=[Exchange 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func CollectionReqVerify(req *go_micro_srv_bill.CollectionRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.RecAccountUid, "!=", "", ss_err.ERR_PARAM, "RecAccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.GenCode, "!=", "", ss_err.ERR_PARAM, "GenCode")
	if !fildChk.IsOk {
		ss_log.Error("err=[Collection 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func TransferReqVerify(req *go_micro_srv_bill.TransferRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ToPhone, "!=", "", ss_err.ERR_PARAM, "ToPhone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ExchangeType, "!=", "", ss_err.ERR_PARAM, "ExchangeType")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	if !fildChk.IsOk {
		ss_log.Error("err=[Transfer 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func SaveMoneyReqVerify(req *go_micro_srv_bill.SaveMoneyRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.RecvPhone, "!=", "", ss_err.ERR_PARAM, "RecvPhone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.SendPhone, "!=", "", ss_err.ERR_PARAM, "SendPhone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	if !fildChk.IsOk {
		ss_log.Error("err=[SaveMoney 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

//手机号取款
func MobileNumWithdrawalReqVerify(req *go_micro_srv_bill.WithdrawalRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.RecvPhone, "!=", "", ss_err.ERR_PARAM, "RecvPhone")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.SendPhone, "!=", "", ss_err.ERR_PARAM, "SendPhone")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_WALLET_AMOUNT_FAILD, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpAccNo, "!=", "", ss_err.ERR_PARAM, "OpAccNo")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.SaveCode, "!=", "", ss_err.ERR_PARAM, "SaveCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	if !fildChk.IsOk {
		ss_log.Error("err=[MobileNumWithdrawal 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

//扫一扫取款参数校验
func SweepWithdrawalReqVerify(req *go_micro_srv_bill.SweepWithdrawRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.GenCode, "!=", "", ss_err.ERR_PARAM, "GenCode")
	if !fildChk.IsOk {
		ss_log.Error("err=[SweepWithdrawal 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

//pos机确认扫码操作参数验证
func ConfirmpWithdrawalReqVerify(req *go_micro_srv_bill.ConfirmWithdrawRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.UseAccountUid, "!=", "", ss_err.ERR_PARAM, "UseAccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.GenCode, "!=", "", ss_err.ERR_PARAM, "GenCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OutOrderNo, "!=", "", ss_err.ERR_PARAM, "OutOrderNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[ConfirmpWithdrawal 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ModityGenCodeStatusReqVerify(req *go_micro_srv_bill.ModifyGenCodeStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.GenKey, "!=", "", ss_err.ERR_PARAM, "GenKey")
	if !fildChk.IsOk {
		ss_log.Error("err=[ModityGenCodeStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func CancelWithdrawalReqVerify(req *go_micro_srv_bill.CancelWithdrawRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CancelReason, "!=", "", ss_err.ERR_PARAM, "CancelReason")
	if !fildChk.IsOk {
		ss_log.Error("err=[CancelWithdrawal 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func WithdrawReceiptReqVerify(req *go_micro_srv_bill.WithdrawReceiptRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[WithdrawReceipt 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func QuerySaveReceiptReqVerify(req *go_micro_srv_bill.QuerySaveReceiptRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[QuerySaveReceipt 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func TransferToHeadquartersReqVerify(req *go_micro_srv_bill.TransferToHeadquartersRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ImageId, "!=", "", ss_err.ERR_PARAM, "ImageId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecName, "!=", "", ss_err.ERR_PARAM, "RecName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecCarNum, "!=", "", ss_err.ERR_PARAM, "RecCarNum")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpAccNo, "!=", "", ss_err.ERR_PARAM, "OpAccNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[TransferToHeadquarters 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func CustTransferToHeadquartersReqVerify(req *go_micro_srv_bill.CustTransferToHeadquartersRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ImageId, "!=", "", ss_err.ERR_PARAM, "ImageId")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecName, "!=", "", ss_err.ERR_PARAM, "RecName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecCarNum, "!=", "", ss_err.ERR_PARAM, "RecCarNum")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	//fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpAccNo, "!=", "", ss_err.ERR_PARAM, "OpAccNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.CardNo, "!=", "", ss_err.ERR_PARAM, "CardNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[TransferToHeadquarters 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}

func ApplyMoneyReqVerify(req *go_micro_srv_bill.ApplyMoneyRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.ChannelName, "!=", "", ss_err.ERR_PARAM, "ChannelName")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecCarNum, "!=", "", ss_err.ERR_PARAM, "RecCarNum")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OpAccNo, "!=", "", ss_err.ERR_PARAM, "OpAccNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[ApplyMoney 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func CustWithdrawReqVerify(req *go_micro_srv_bill.CustWithdrawRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.RecCarNo, "!=", "", ss_err.ERR_PARAM, "RecCarNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Password, "!=", "", ss_err.ERR_PARAM, "Password")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.NonStr, "!=", "", ss_err.ERR_PARAM, "NonStr")
	fildChk = ss_chk.DoFieldChk(fildChk, req.WithdrawType, "!=", "", ss_err.ERR_PARAM, "WithdrawType")
	if !fildChk.IsOk {
		ss_log.Error("err=[CustWithdraw 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func GenRecvCodeReqVerify(req *go_micro_srv_bill.GenRecvCodeRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	if !fildChk.IsOk {
		ss_log.Error("err=[GenRecvCode 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func QueryRateReqVerify(req *go_micro_srv_bill.QeuryRateRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountUid, "!=", "", ss_err.ERR_PARAM, "AccountUid")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.IdenNo, "!=", "", ss_err.ERR_PARAM, "IdenNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[QueryRate 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func CustQueryRateReqVerify(req *go_micro_srv_bill.CustQeuryRateRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.AccountType, "!=", "", ss_err.ERR_PARAM, "AccountType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.MoneyType, "!=", "", ss_err.ERR_PARAM, "MoneyType")
	fildChk = ss_chk.DoFieldChk(fildChk, req.Amount, "!=", "", ss_err.ERR_PARAM, "Amount")
	fildChk = ss_chk.DoFieldChk(fildChk, req.ChannelNo, "!=", "", ss_err.ERR_PARAM, "ChannelNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[CustQueryRate 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func QuerySweepCodeStatusReqVerify(req *go_micro_srv_bill.QuerySweepCodeStatusRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.GenCode, "!=", "", ss_err.ERR_PARAM, "GenCode")
	fildChk = ss_chk.DoFieldChk(fildChk, req.AccountNo, "!=", "", ss_err.ERR_PARAM, "AccountNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[QuerySweepCodeStatus 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func CustOrderBillDetailReqVerify(req *go_micro_srv_bill.CustOrderBillDetailRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	fildChk = ss_chk.DoFieldChk(fildChk, req.OrderType, "!=", "", ss_err.ERR_PARAM, "OrderType")
	if !fildChk.IsOk {
		ss_log.Error("err=[CustOrderBillDetail 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func SaveWithdrawDetailReqVerify(req *go_micro_srv_bill.SaveDetailRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[SaveWithdrawDetail 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func SweepWithdrawDetailReqVerify(req *go_micro_srv_bill.SweepWithdrawDetailRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[SweepWithdrawDetail 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
func SaveMoneyDetailReqVerify(req *go_micro_srv_bill.SaveMoneyDetailRequest) string {
	fildChk := ss_chk.DoFieldChk(nil, req.OrderNo, "!=", "", ss_err.ERR_PARAM, "OrderNo")
	if !fildChk.IsOk {
		ss_log.Error("err=[SaveMoneyDetail 接口------> %s 为空]", fildChk.ErrReason)
		return fildChk.ErrCode
	}
	return ""
}
