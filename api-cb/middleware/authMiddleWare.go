package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/api-cb/inner_util"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type AuthMiddleWare struct {
}

var AuthMiddleWareInst AuthMiddleWare

/**
 * 认证中间件
 */
func (*AuthMiddleWare) DoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ss_log.Info("AuthMiddleWare|DoAuth")
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================

		if !strext.ToBool(inner_util.S(c, common.INNER_SIGN_VERIFY)) {
			c.Set(common.RET_CODE, ss_err.ERR_SYS_SIGN)
			c.Set(common.INNER_IS_STOP, true)
			return
		}
	}
}

//
//func decodeData(in, orgNo, mercNo, encryptkey, action, platPriKey, timeStr, signType string) (map[string]interface{}, error) {
//	switch signType {
//	case common.SignType_Md5:
//		tmpMap := map[string]interface{}{}
//		json.Unmarshal([]byte(in), &tmpMap)
//		return tmpMap, nil
//	case common.SignType_Rsa:
//		// 正常
//		_encryptkey := strings.Replace(encryptkey, " ", "", -1)
//		ss_log.Info("key=%v\n", _encryptkey)
//
//		// key
//		_encryptkey2, err := encrypt.DoBase64(encrypt.HANDLE_DECRYPT, _encryptkey)
//		outStr, err := encrypt.DoRsa(encrypt.HANDLE_DECRYPT, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_STRING, encrypt.HASHLENTYPE_NONE,
//			_encryptkey2, platPriKey, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
//		if err != nil {
//			ss_log.Info("err=%v|decoded=%v\n", err, outStr)
//			return nil, err
//		}
//
//		crypted, err := hex.DecodeString(in)
//		dataOut, _ := encrypt.DoAes(encrypt.AES_BM_ECB, outStr.(string), encrypt.HANDLE_DECRYPT, []byte(crypted),
//			encrypt.PADDING_PKCS5, encrypt.RETFMT_STRING, encrypt.PRE_HANDLE_BYTE, nil)
//		ss_log.Info("dataOut=%v\n", dataOut)
//		tmpMap := map[string]interface{}{}
//		json.Unmarshal([]byte(dataOut.(string)), &tmpMap)
//
//		return tmpMap, nil
//	}
//
//	return nil, nil
//}
//
//func verifySign(inData, orgNo, mercNo, action, mercPubKey, mercMd5Key, timeStr, inSign, signType string) (bool, string) {
//	// key
//	ss_log.Info("sign=[%v]\npub=[%v]", inSign, mercPubKey)
//	// 全部小写
//	inSign = strings.ToLower(inSign)
//	var reqStrEnBefore string
//	switch signType {
//	case common.SignType_Rsa:
//		outStr, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA512,
//			map[string]interface{}{
//				"sign": inSign,
//				"data": ([]byte)(orgNo + mercNo + action + timeStr + inData),
//			}, mercPubKey, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP_BASE64)
//		if err != nil {
//			ss_log.Error("err=%v|decoded=%v\ndata=[%v]", err, outStr, orgNo+mercNo+action+timeStr+inData)
//			return false, ""
//		}
//	case common.SignType_Md5:
//		dataMap := make(map[string]interface{})
//		err := json.Unmarshal([]byte(inData), &dataMap)
//		if err != nil {
//			ss_log.Error("err=[%v]", err)
//			return false, ""
//		}
//
//		dataMap["org_no"] = orgNo
//		dataMap["merc_no"] = mercNo
//		dataMap["action"] = action
//		dataMap["time"] = timeStr
//
//		reqStrEnBefore = encrypt.Map2FormStr(dataMap, mercMd5Key, "&key=", encrypt.FIELD_ENCODED_NONE,
//			[]string{}, "sign", false)
//		// 全部小写
//		sign := strings.ToLower(encrypt.DoMd5(reqStrEnBefore))
//		ss_log.Info("before=[%v]", reqStrEnBefore)
//		if sign != inSign {
//			ss_log.Error("sign=[%v]|inSign=[%v]", sign, inSign)
//			return false, reqStrEnBefore
//		}
//	}
//
//	return true, reqStrEnBefore
//}
//
