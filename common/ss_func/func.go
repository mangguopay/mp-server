package ss_func

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/ss_err"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
)

// 获取服务的完整名称
// serverName : service.Server.Options().Name
// serverId : service.Server.Options().Id
func GetServerFullId(serverName, serverId string) string {
	return serverName + "-" + serverId
}

// 国家码加前缀,如传入 855,返回 0855
func PreCountryCode(countryCode string) string {
	return fmt.Sprintf("%04d", strext.ToInt(countryCode))
}

// 处理前端传来的手机号前缀，注意如果是jwt或数据库获取的手机号则不需要调用本方法
func PrePhone(countryCode, phone string) string {
	//柬埔寨的去掉手机号第1位0，如果有多个0，也只去掉第1个,如传入 001002,返回 01002
	if countryCode == constants.CountryCode_KH {
		phone = strings.TrimPrefix(phone, "0")
	}

	return phone
}

// 账号组合成国家码+手机号
func ComposeAccountByPhoneCountryCode(phone, countryCode string) string {
	return fmt.Sprintf("%s%s", PreCountryCode(countryCode), phone)
}

// 正常化语言类型
func NormalizeLang(lang string) string {
	// 如果为空或不是允许的语言类型,则默认为英语
	if lang == "" || !util.InSlice(lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		lang = constants.DefaultLang
	}

	return lang
}

// 判断金额是正数？是负数？还是0
//
// 参数amount支持“整数字符串”和“浮点字符串”
//
// 如果amount>0 返回 1
// 如果amount=0 返回 0
// 如果amount<0 返回 -1
func JudgeAmountPositiveOrNegative(amount string) (int, error) {
	ret := 0

	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return ret, err
	}

	if f > 0 {
		ret = 1
	} else if f < 0 {
		ret = -1
	}

	return ret, nil
}

// 将参数排序并拼接成字符串
func ParamsMapToString(params map[string]interface{}, signKey string) string {
	var pList = make([]string, 0)

	for key, value := range params {
		if strings.TrimSpace(key) == signKey { // 忽略验签字段
			continue
		}
		// 将interface转换为字符串
		val := strings.TrimSpace(strext.ToStringNoPoint(value))
		if len(val) > 0 { // 忽略空值
			pList = append(pList, key+"="+val)
		}
	}

	// 按键排序
	sort.Strings(pList)

	// 使用&符号拼接
	return strings.Join(pList, "&")
}

//校验国家码是否合法
func CheckCountryCode(countryCode string) (errStr string) {
	switch countryCode {
	case constants.CountryCode_ZHCN: //中国86
	case constants.CountryCode_KH: //柬埔寨855
	default: //其他情况不合法
		return ss_err.ERR_CountryCode_FAILD
	}
	return ss_err.ERR_SUCCESS
}

//校验商家前面方式是否合法
func CheckSignMethod(signMethod string) (errStr string) {
	switch signMethod {
	case constants.SignMethod_RSA2: //RSA2
	default: //其他情况不合法
		return ss_err.ERR_SignMethod_FAILD
	}
	return ss_err.ERR_SUCCESS
}

/**
手机号脱敏
规则：
手机号9位以下保留前后两位，如 12345678 ---> 12****78
手机号9位和9位以上保留前后三位, 如123456789 --->123****789
*/
func GetDesensitizationPhone(phone string) (phone2 string) {
	if phone != "" {
		length := len(phone)
		if length > 9 {
			phone2 = phone[:3] + "****" + phone[(length-3):]
		} else if length < 9 && length >= 5 {
			phone2 = phone[:2] + "****" + phone[(length-2):]
		} else {
			ss_log.Error("手机号码[%v]位数异常(长度小于5),", phone)
		}
	}
	return phone2
}

/**
根据国家码给手机号脱敏
规则：
区号为0855的只显示首尾2位数，区号为0086的显示首尾3位
例如“付款方手机号”为0855-88889999，脱敏后前端显示为“0855-88****99”
例如“付款方手机号”为0086-15939519160，脱敏后前端显示为“0086-159****160”
*/
func GetDesensitizationPhoneByCountryCode(countryCode, phone string) (string, error) {
	var phoneDes string
	var err error
	if countryCode != "" && phone != "" {
		length := len(phone)
		if countryCode == constants.CountryCode_ZHCN {
			phoneDes = phone[:3] + "****" + phone[(length-3):]
		} else if countryCode == constants.CountryCode_KH {
			phoneDes = phone[:2] + "****" + phone[(length-2):]
		} else {
			err = errors.New(fmt.Sprintf("国家码[%v]或手机号码[%v]异常", countryCode, phone))
		}
	} else {
		err = errors.New(fmt.Sprintf("国家码[%v]或手机号码[%v]错误", countryCode, phone))
	}
	phoneDes = fmt.Sprintf("%v-%v", countryCode, phoneDes)
	return phoneDes, err
}

/**
邮箱脱敏
规则：
@前面位数<=3,保留两位，后面用****代替。如a22@163.com ----> a2****@163.com
大于3位则保留三位，后面用*代替。如 445522214@qq.com ---->  445****@qq.com
*/
func GetDesensitizationEmail(email string) (email2 string) {
	if email != "" {
		arr := strings.Split(email, "@")
		preEmail := arr[0]
		lenInt := len(arr[0])
		if lenInt > 3 {
			email2 = preEmail[:3] + "****@" + arr[1]
		} else if lenInt <= 3 && lenInt > 2 {
			email2 = preEmail[:2] + "****@" + arr[1]
		} else {
			ss_log.Error("邮箱@前面[%v]位数异常(长度<=2),", preEmail)
			email2 = email
		}
	}

	return email2
}

/**
账号脱敏
规则：
区号为855的只显示首尾2位数，区号为86的显示首尾3位
例如“付款方账号”为085588889999，脱敏后返回“085588****99”
例如“付款方账号”为008615939519160，脱敏后返回“0086159****160”
*/
func GetDesensitizationAccount(account string) (accountT string) {
	if strings.Contains(account, "@") { //商家账号处理
		accountT = GetDesensitizationEmail(account)
	} else { //用户账号处理
		if len(account) >= 5 {
			countryCode := account[:4]
			account2 := account[4:]

			switch countryCode {
			case PreCountryCode(constants.CountryCode_ZHCN): //中国0086
				account2 = account2[:3] + "****" + account2[len(account2)-3:]
			case PreCountryCode(constants.CountryCode_KH): //柬埔寨0855
				account2 = account2[:2] + "****" + account2[len(account2)-2:]
			default:
				account2 = account2[:1] + "****"
			}

			accountT = countryCode + account2
		}
	}

	return accountT
}

// 处理前端传来的第三方渠道卡号（第三方的渠道卡号都是柬埔寨的手机号，所以要删除用户可能输入的0前缀）
// 去掉第1位0，如果有多个0，也只去掉第1个,如传入 001002,返回 01002
func PreThirdpartyCardNumber(cardNumber string) string {
	cardNumber = strings.TrimPrefix(cardNumber, "0")

	return cardNumber
}
