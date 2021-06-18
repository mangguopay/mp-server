package p

import (
	"a.a/cu/encrypt"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/s1/handler"
	"io/ioutil"
)

const (
	salt    = `sa5d6g728ttg$%43JASHGFUIa72`
	UrlBase = `http://127.0.0.1:8080/mobile`
)

var (
	M = map[string]interface{}{}
)

func DoLogin() {
	str := util.RandomDigitStr(6)
	resp := handler.DoSend(UrlBase+"/auth/login", map[string]interface{}{
		"account":      "13800138000",
		"password":     doInitPassword3("1", str),
		"imei":         "123",
		"nonstr":       str,
		"lang":         "zh_CN",
		"pos_sn":       "pos-q",
		"account_type": 1,
		"app_version":  "",
	}, "")
	token := strext.ToStringNoPoint(resp.ExtMap["data"].(map[string]interface{})["userinfo"].(map[string]interface{})["token"])
	M["token"] = token
	ioutil.WriteFile("a.log", []byte(strext.ToJson(M)), 0777)
}

func doInitPassword1(password string) string {
	return encrypt.DoShaXXX(password, encrypt.HASHLENTYPE_SHA1)
}

func doInitPassword2(password string) string {
	return encrypt.DoMd5Salted(doInitPassword1(password), salt)
}

func doInitPassword3(password, nonstr string) string {
	return encrypt.DoMd5Salted(doInitPassword2(password), nonstr)
}

func Reload() {
	a, _ := ioutil.ReadFile("a.log")
	M = strext.Json2Map(a)
}
