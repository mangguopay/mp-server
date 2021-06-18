package middleware

import (
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-pos/common"
	"a.a/mp-server/common/ss_err"
	"encoding/base64"
	"github.com/gin-gonic/gin"
)

type RsaMw struct {
}

var RsaMwInst RsaMw

/**
 * 解密rsa
 */
func (RsaMw) DecodeRsa() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================

		isEncoded := c.GetBool(common.INNER_IS_ENCODED)
		// 需要解密才处理
		if !isEncoded {
			return
		}

		traceNo := c.GetString(common.INNER_TRACE_NO)
		p, _ := c.Get(common.INNER_PARAM_MAP)
		switch p2 := p.(type) {
		case map[string]interface{}:
			// rsa解密
			//ss_log.Info("pri--------------->%s", strext.ToStringNoPoint(i.EncryptMap["pri_key"]))
			dencr, derr := encrypt.DoRsa(encrypt.HANDLE_DECRYPT, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_STRING, encrypt.HASHLENTYPE_SHA256,
				strext.ToStringNoPoint(p2["data"]), strext.ToStringNoPoint(common.EncryptMap["pri_key"]), encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BASE64)
			if derr != nil {
				ss_log.Error("%v|============解密失败,err=[%v]", traceNo, derr)
				c.Set(common.INNER_IS_STOP, true)
				c.Set(common.RET_CODE, ss_err.ERR_SYS_DECODE)
				return
			}
			c.Set(common.INNER_POST_RSA_PARAM, dencr)

			// decodeBase64
			//str, err := base64.StdEncoding.DecodeString(dencr.(string)) // 明文byte
			ss_log.Info("dencr:-------------------->%s", strext.ToStringNoPoint(dencr))
			dencrBody, err := base64.StdEncoding.DecodeString(strext.ToStringNoPoint(dencr)) // 明文byte
			if err != nil {
				ss_log.Error("解析参数base64decode失败,err: %v|----------------------------", err.Error())
				c.Set(common.INNER_IS_STOP, true)
				c.Set(common.RET_CODE, ss_err.ERR_SYS_DECODE)
				return
			}
			p3 := strext.Json2Map(dencrBody)
			c.Set(common.INNER_PARAM_MAP, p3)
			ss_log.Info("%v|----------------------------POST解密后的参数", traceNo)
			for k, v := range p3 {
				ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v)
			}
			ss_log.Info("%v|----------------------------", traceNo)
			c.Set(common.INNER_SIGN, strext.ToStringNoPoint(p2["sign"]))
			return
		default:
			ss_log.Error("body格式错误")
			c.Set(common.INNER_IS_STOP, true)
			c.Set(common.RET_CODE, ss_err.ERR_SYS_BODY_NOT_JSON)
			return
		}
	}
}

/**
 * rsa验签
 */
func (RsaMw) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		isEncoded := c.GetBool(common.INNER_IS_ENCODED)
		// 需要解密才处理
		if !isEncoded {
			return
		}

		traceNo := c.GetString(common.INNER_TRACE_NO)
		param := c.GetString(common.INNER_POST_RSA_PARAM)
		sign := c.GetString(common.INNER_SIGN)

		m2 := encrypt.DoMd5(param)
		pubKey := c.GetString(common.Pub_Key)
		// 验证签名
		_, err := encrypt.DoRsa(encrypt.HANDLE_VERIFY, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_NONE, encrypt.HASHLENTYPE_SHA256,
			map[string]interface{}{
				"sign": sign,
				"data": []byte(m2),
			}, pubKey, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_SIGN_MAP_BASE64)
		if err != nil {
			ss_log.Error("%v|===========验证签名失败,err=[%v]", traceNo, err)
			c.Set(common.INNER_SIGN_VERIFY, false)
			return
		}

		c.Set(common.INNER_SIGN_VERIFY, true)
		return
	}
}

func (r *RsaMw) Encode() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		r.EncodeInner(c)
	}
}

func (r *RsaMw) Sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		r.SignInner(c)
		return
	}
}

func (*RsaMw) SignInner(c *gin.Context) {
	dataByteT := c.GetString(common.INNER_PRESEND_ENCODED)
	dataByteT2 := c.GetString(common.INNER_PRESEND_ENCODED2)
	//==================签名==========================
	dataByteMd5 := encrypt.DoMd5(dataByteT2)
	ss_log.Info("md5=[%v]", dataByteMd5)
	sign, err := encrypt.DoRsa(encrypt.HANDLE_SIGN, encrypt.KEYTYPE_PKCS8, encrypt.RETFMT_BASE64, encrypt.HASHLENTYPE_SHA256,
		[]byte(dataByteMd5), strext.ToStringNoPoint(common.EncryptMap["pri_key"]), encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
	if err != nil {
		ss_log.Error("============签名失败", err)
		return
	}
	//==================签名==========================
	resp2 := gin.H{
		"data": strext.ToStringNoPoint(dataByteT),
		"sign": strext.ToStringNoPoint(sign),
	}
	// ===============返回加密结束===============
	c.Set(common.RET_DATA_PRESEND, resp2)
	return
}

func (*RsaMw) EncodeInner(c *gin.Context) {
	traceNo := c.GetString(common.INNER_TRACE_NO)
	p := c.GetString(common.RET_DATA_PRESEND)
	// base64
	dataByteT := base64.StdEncoding.EncodeToString([]byte(p))
	pubKey := c.GetString(common.Pub_Key)
	ss_log.Info("%v|加密时使用的公钥=[%s]", traceNo, pubKey)
	//=====================加密=======================
	encr, err := encrypt.DoRsa(encrypt.HANDLE_ENCRYPT, encrypt.KEYTYPE_PKIX, encrypt.RETFMT_BASE64, encrypt.HASHLENTYPE_NONE,
		[]byte(dataByteT), pubKey, encrypt.KEYFMT_PEM, encrypt.PRE_HANDLE_BYTE)
	if err != nil {
		ss_log.Error("%v|============加密失败,err=[%v]", traceNo, err)
		return
	}
	//=====================加密=======================
	c.Set(common.INNER_PRESEND_ENCODED, encr)
	c.Set(common.INNER_PRESEND_ENCODED2, dataByteT)
	return
}
