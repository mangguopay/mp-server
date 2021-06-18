package common

import (
	"fmt"

	"a.a/cu/strext"
	"a.a/mp-server/common/ss_struct"
)

var EncryptMap map[string]interface{}

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

// 国家码去前缀,如传入 855,返回 0855
func PreCountryCode(countryCode string) string {
	return fmt.Sprintf("%04d", strext.ToInt(countryCode))
}
