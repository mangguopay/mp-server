package common

import (
	"time"

	"a.a/cu/encrypt"
	"a.a/cu/jwt"
	"a.a/cu/ss_log"
	"a.a/mp-server/auth-srv/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_struct"
)

var EncryptMap map[string]interface{}

func CheckPayPWD(accountType, idenNo, nonStr, reqPWD string) string {
	var payPwdDb string // 数据库密码
	switch accountType {
	case constants.AccountType_SERVICER: //服务商
		payPwdDb = dao.ServicerDaoInstance.GetServicerPWDFromIdenNo(idenNo)
	case constants.AccountType_POS: // 收银员
		payPwdDb = dao.CashierDaoInstance.GetCashierPwdFromIdenNo(idenNo)
	case constants.AccountType_USER:
		payPwdDb = dao.CustDaoInstance.QueryPwdFromIdenNo(idenNo)
	case constants.AccountType_PersonalBusiness:
		fallthrough
	case constants.AccountType_EnterpriseBusiness:
		payPwdDb = dao.BusinessDaoInstance.GetPayPwdFromIdenNo(idenNo)
	}
	if payPwdDb == "" {
		ss_log.Error("idenNo[%v],accountType[%v]数据库支付密码为空，无法开始验证支付密码", idenNo, accountType)
		return ss_err.ERR_PAY_PWD_IS_NULL
	}
	// 数据库取出的支付密码加盐(加的是和前端传来的盐一样)
	pwdMD5FixedDB := encrypt.DoMd5Salted(payPwdDb, nonStr)
	if reqPWD != pwdMD5FixedDB {
		ss_log.Error("reqPWD[%v]---pwdMD5FixedDB[%v]密码校验失败，密码错误", reqPWD, pwdMD5FixedDB)
		return ss_err.ERR_DB_PWD
	}
	return ss_err.ERR_SUCCESS
}

func CreateWebBusinessJWT(data ss_struct.JwtDataWebBusiness) (jwtStr string) {
	retMap := JwtStructToMapWebBusiness(data)

	k1, loginSignKey, err := cache.ApiDaoInstance.GetGlobalParam("login_sign_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k1)
	}
	k2, loginAesKey, err := cache.ApiDaoInstance.GetGlobalParam("login_aes_key")
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err, k2)
	}
	jwt2 := jwt.GetNewEncryptedJWTToken(time.Hour*4, retMap, loginAesKey, loginSignKey)
	return jwt2
}

// app端jwt结构体转map
func JwtStructToMapApp(data ss_struct.JwtDataApp) map[string]string {
	return map[string]string{
		"account":      data.Account,
		"account_uid":  data.AccountUid,
		"iden_no":      data.IdenNo,
		"account_type": data.AccountType,
		//"account_name":     data.AccountName,
		"login_account_no": data.LoginAccountNo,
		"pub_key":          data.PubKey,
		"jump_iden_no":     data.JumpIdenNo,
		"jump_iden_type":   data.JumpIdenType,
		"master_acc_no":    data.MasterAccNo,
		"is_master_acc":    data.IsMasterAcc,
		"pos_sn":           data.PosSn,
	}
}

// webAdmin端jwt结构体转map
func JwtStructToMapWebAdmin(data ss_struct.JwtDataWebAdmin) map[string]string {
	return map[string]string{
		"account":     data.Account,
		"account_uid": data.AccountUid,
		//"merchant_uid": data.MerchantUid,
		"iden_no":      data.IdenNo,
		"account_type": data.AccountType,
		//"account_name":     getAcc.Nickname,
		"login_account_no": data.LoginAccountNo,
		"jump_iden_no":     data.JumpIdenNo,
		"jump_iden_type":   data.JumpIdenType,
		"master_acc_no":    data.MasterAccNo,
		"is_master_acc":    data.IsMasterAcc,
	}

}

// webBusiness端jwt结构体转map
func JwtStructToMapWebBusiness(data ss_struct.JwtDataWebBusiness) map[string]string {
	return map[string]string{
		"account":          data.Account,
		"account_uid":      data.AccountUid,
		"iden_no":          data.IdenNo,
		"account_type":     data.AccountType,
		"email":            data.Email,
		"phone":            data.Phone,
		"country_code":     data.CountryCode,
		"login_account_no": data.LoginAccountNo,
		"jump_iden_no":     data.JumpIdenNo,
		"jump_iden_type":   data.JumpIdenType,
		"master_acc_no":    data.MasterAccNo,
		"is_master_acc":    data.IsMasterAcc,
	}

}
