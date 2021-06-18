package util

import (
	"a.a/cu/ss_big"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/bill-srv/dao"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"database/sql"
	"github.com/shopspring/decimal"
	"strings"
)

type FileContentJsonStruct struct {
	Row          string `json:"Row"`
	ToAccount    string `json:"ToAccount"`
	Amount       string `json:"Amount"`
	CurrencyType string `json:"CurrencyType"`
	Name         string `json:"Name"`
	Remarks      string `json:"Remarks"`
}

//校验转账金额
//注意,这里的amount金额是乘以100后的数(如USD1.5则传的是150,而KHR100还是100)
func CheckTransferAmount(amount, currencyType string, transferConf *dao.BusinessTransferParamValue) (wrongReason string) {
	if strext.ToInt(amount) <= 0 {
		wrongReason = ss_err.ERR_WALLET_AMOUNT_NULL // "金额小于或等于0"
		return wrongReason
	}

	switch currencyType {
	case constants.CURRENCY_UP_USD:
		if strext.ToInt64(amount) < transferConf.USDMinAmount && strext.ToInt64(amount) > transferConf.USDMaxAmount {
			ss_log.Error("转账金额超出限制,USDMinAmount[%v],USDMaxAmount[%v],Amount[%v]", transferConf.USDMinAmount, transferConf.USDMaxAmount, amount)
			wrongReason = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
		}

	case constants.CURRENCY_UP_KHR:
		if strext.ToInt64(amount) < transferConf.KHRMinAmount && strext.ToInt64(amount) > transferConf.KHRMaxAmount {
			ss_log.Error("转账金额超出限制,KHRMinAmount[%v],KHRMaxAmount[%v],Amount[%v]", transferConf.KHRMinAmount, transferConf.KHRMaxAmount, amount)
			wrongReason = ss_err.ERR_PAY_AMOUNT_IS_LIMIT
		}
	case "":
		wrongReason = ss_err.ERR_UnFilledCurrencyType_FAILD // "未填写币种"
	default:
		wrongReason = ss_err.ERR_MONEY_TYPE_FAILD
	}

	return wrongReason
}

//计算手续费
func QueryTransferFeeAndRate(amount, currencyType string, transferConf *dao.BusinessTransferParamValue) (fee, rate, wrongReason string) {
	var feesDeci decimal.Decimal
	switch currencyType {
	case constants.CURRENCY_UP_USD:
		rate = strext.ToString(transferConf.USDRate)
		feesDeci = ss_count.CountFees(amount, rate, strext.ToStringNoPoint(transferConf.USDMinFee))
	case constants.CURRENCY_UP_KHR:
		rate = strext.ToString(transferConf.KHRRate)
		feesDeci = ss_count.CountFees(amount, rate, strext.ToStringNoPoint(transferConf.KHRMinFee))
	}
	// 取整
	fee = ss_big.SsBigInst.ToRound(feesDeci, 0, ss_big.RoundingMode_HALF_EVEN).String()
	if fee == "" {
		fee = "0"
	}
	if rate == "" {
		rate = "0"
	}

	return fee, rate, wrongReason
}

//确认商家信息
func checkBusinessInfo(toAccount, authName, fromAccountNo string) (toAccUid, toBusinessNo, wrongReason string) {
	//根据商家账号查询商家uid、账号uid
	toAccUid, err := dao.AccDaoInstance.GetAccNoFromAccount(toAccount)
	if err != nil {
		ss_log.Error("err[%v]", err)
		wrongReason = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS //账号不存在
		return toAccUid, toBusinessNo, wrongReason
	}

	if toAccUid == "" {
		ss_log.Error("账号[%v]不存在", toAccount)
		wrongReason = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS //账号不存在
		return toAccUid, toBusinessNo, wrongReason
	} else if toAccUid == fromAccountNo {
		ss_log.Error("不能付款给自己，toAccUid[%v],FromAccountNo[%v]", toAccUid, fromAccountNo)
		wrongReason = ss_err.ERR_ACCOUNT_TRANSFER_TO_SELF // "自己不能给自己转账"
		return toAccUid, toBusinessNo, wrongReason
	} else {
		//查询收款商家的状态和收款权限
		inBusiness, err := dao.BusinessDaoInst.GetBusinessStatusInfo("", toAccUid)
		if err != nil {
			if err == sql.ErrNoRows {
				ss_log.Error("收款商家不存在，toAccUid=%v, err=%v", toAccUid, err)
				wrongReason = ss_err.ERR_PayeeNotExist
				return toAccUid, toBusinessNo, wrongReason
			}
			ss_log.Error("查询收款商家状态失败，toAccUid=%v, err=%v", toAccUid, err)
			wrongReason = ss_err.ERR_SYSTEM
			return toAccUid, toBusinessNo, wrongReason
		}

		toBusinessNo = inBusiness.BusinessNo

		if strext.ToInt(inBusiness.AccountStatus) != constants.AccountUseStatusNormal {
			ss_log.Error("收款商家账号已被禁用，toAccUid=%v, err=%v", toAccUid, err)
			wrongReason = ss_err.ERR_MERC_NO_USE
			return toAccUid, toBusinessNo, wrongReason
		}
		if strext.ToInt(inBusiness.BusinessStatus) == constants.BusinessUseStatusDisabled {
			ss_log.Error("收款商家已被禁用，toAccUid=%v, err=%v", toAccUid, err)
			wrongReason = ss_err.ERR_MERC_NO_USE
			return toAccUid, toBusinessNo, wrongReason
		}
		if strext.ToInt(inBusiness.IncomeAuthorization) == constants.BusinessIncomeAuthDisabled {
			ss_log.Error("收款商家没有收款权限，toAccUid=%v, err=%v", toAccUid, err)
			wrongReason = ss_err.ERR_MERC_NO_USE
			return toAccUid, toBusinessNo, wrongReason
		}

		if authName == "" { //等于空说明没有认证，或查询认证信息失败
			ss_log.Error("账号未认证商家名称")
			wrongReason = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH
			return toAccUid, toBusinessNo, wrongReason
		}

		if inBusiness.FullName != authName {
			ss_log.Error("商家名称错误,FullName[%v],authName[%v]", inBusiness.FullName, authName)
			wrongReason = ss_err.ERR_AuthName_FAILD
			return toAccUid, toBusinessNo, wrongReason
		}
	}

	return toAccUid, toBusinessNo, wrongReason
}

func checkUserInfo(toAccount, authName, fromAccountNo string) (toAccUid, wrongReason string) {
	accountArr := strings.Split(toAccount, "-")

	if len(accountArr) != 2 {
		wrongReason = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS //"账号不存在"
		return "", wrongReason
	}

	//处理国家码将其变成前缀无0的格式
	countryCode := strext.ToString(strext.ToInt(accountArr[0]))
	account := ss_func.PrePhone(countryCode, accountArr[1])

	//验证国家码是否是正确合法的
	if errStr := ss_func.CheckCountryCode(countryCode); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("err[%v]", errStr)
		wrongReason = ss_err.ERR_CountryCode_FAILD //"国家码错误"
		return "", wrongReason
	}

	//将国家码变成0086、0855的格式，组成账号
	toAccount = ss_func.PreCountryCode(countryCode) + account
	toAccUid, err := dao.AccDaoInstance.GetAccNoFromAccount(toAccount)
	if err != nil {
		ss_log.Error("err[%v]", err)
		wrongReason = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
		return "", wrongReason
	}

	if toAccUid == "" {
		ss_log.Error("账号[%v]不存在", toAccount)
		wrongReason = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
		return toAccUid, wrongReason
	} else if toAccUid == fromAccountNo {
		ss_log.Error("不能付款给自己，toAccUid[%v],FromAccountNo[%v]")
		wrongReason = ss_err.ERR_ACCOUNT_TRANSFER_TO_SELF
		return toAccUid, wrongReason
	} else {

		idenNo := dao.RelaAccIdenDaoInst.GetIdenNo(toAccUid, constants.AccountType_USER)
		if idenNo == "" {
			ss_log.Error("转账的账号不是用户，uid[%v]", toAccUid)
			wrongReason = ss_err.ERR_MSG_ACCOUNT_NOT_EXISTS
			return toAccUid, wrongReason
		}

		//确认用户打开收款权限
		_, _, _, inTransferAuthorization, _ := dao.CustDaoInstance.QueryRateRoleFrom(idenNo)
		if strext.ToInt(inTransferAuthorization) == constants.CustInTransferAuthorizationDisabled {
			ss_log.Error("收款人没有转入权限,accountNo: %s,custNo: %s", toAccount, idenNo)
			wrongReason = ss_err.ERR_MERC_NO_USE //"收款人没有转入权限"
			return toAccUid, wrongReason
		}

		// 获取实名制姓名,报错或者为空说明未实名制
		authNameDb, err := dao.AccDaoInstance.GetAuthNameFromUid(toAccUid)
		if err != nil {
			ss_log.Error(" 查询用户实名认证的姓名失败,toAccUid为: %s,err: %s", toAccUid, err)
			wrongReason = ss_err.ERR_ACCOUNT_NOT_REAL_AUTH //"账号未实名认证"
			return toAccUid, wrongReason
		}

		if authName == "" {
			ss_log.Error("未填写认证名称")
			wrongReason = ss_err.ERR_UnFilledAuthName_FAILD //"未填写认证名称"
			return toAccUid, wrongReason
		}

		// 填写认证名称要和实名认证的名称一致
		if authName != authNameDb {
			ss_log.Error("认证名称错误，authName[%v],authNameDb[%v]", authName, authNameDb)
			wrongReason = ss_err.ERR_AuthName_FAILD //"认证名称错误"
			return toAccUid, wrongReason
		}
	}

	return toAccUid, wrongReason
}

//确认收款方账号存不存在
func CheckToAccountAndAuthName(toAccount, authName, fromAccountNo string) (toAccUid, toBusinessNo, wrongReason string) {

	if strings.Contains(toAccount, "@") { //只有企业商家的账号有@
		toAccUid, toBusinessNo, wrongReason = checkBusinessInfo(toAccount, authName, fromAccountNo)
	} else { //todo 其他情况视为转账给个人
		toAccUid, wrongReason = checkUserInfo(toAccount, authName, fromAccountNo)
	}

	return toAccUid, toBusinessNo, wrongReason
}
