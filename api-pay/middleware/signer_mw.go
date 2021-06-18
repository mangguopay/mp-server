package middleware

import (
	"a.a/cu/strext"

	"a.a/cu/ss_log"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/api-pay/inner_util"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_rsa"
	"github.com/gin-gonic/gin"
)

var SignerMwInst SignerMw

type SignerMw struct {
}

// 验签
func (s *SignerMw) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 获取signType
		signType := inner_util.M(c, common.ParamsSignType)
		if signType == "" {
			ss_log.Error("%v|参数signType不存在或signType为空", traceNo)
			c.Set(common.RET_CODE, ss_err.ACErrMissParamSignType)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		// 检测app配置的前面方式是否和接口传入的一致
		if configSignType := inner_util.S(c, common.InnerAppSignMethod); signType != configSignType {
			ss_log.Error("%v|参数signType与配置的不一致|signType:%s|configSignType:%s", traceNo, signType, configSignType)
			c.Set(common.RET_CODE, ss_err.ACErrSignTypeNotMatch)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		retCode := ss_err.ACErrSignTypeErr
		verifyOk := false

		switch signType {
		case common.SignMethod_RSA2: // RSA2验签
			retCode, verifyOk = s.verifyRSA2(c, traceNo)
		case common.SignMethod_MD5: // MD5验签
			retCode, verifyOk = s.verifyMD5(c, traceNo)
		}

		if !verifyOk { // 验签失败了，停止后面的执行
			c.Set(common.RET_CODE, retCode)
			c.Abort()
			return
		}

		return
	}
}

// RSA2验签
func (s *SignerMw) verifyRSA2(c *gin.Context, traceNo string) (string, bool) {
	// 获取参数
	paramsMap := c.GetStringMap(common.INNER_PARAM_MAP)
	if len(paramsMap) == 0 {
		ss_log.Error("%v|参数为空", traceNo)
		return ss_err.ACErrParamEmpty, false
	}

	// 获取sign
	sign := inner_util.M(c, common.ParamsSign)
	if sign == "" {
		ss_log.Error("%v|参数sign不存在或sign为空", traceNo)
		return ss_err.ACErrMissParamSign, false
	}

	// 获取商家公钥
	pubKey := inner_util.S(c, common.InnerAppBusinessPublickKey)
	if pubKey == "" {
		ss_log.Error("%v|获取商家公钥失败", traceNo)
		return ss_err.ACErrSysErr, false
	}

	ss_log.Info("paramsMap:%v", paramsMap)

	// 将参数排序并拼接成字符串
	data := ss_func.ParamsMapToString(paramsMap, common.ParamsSign)

	// rsa2验签
	err := ss_rsa.RSA2Verify(data, sign, pubKey)
	if err != nil {
		ss_log.Error("%v|rsa2验签失败err:%v|data:%s|sign:%s|pubKey:%s", traceNo, err, data, sign, pubKey)
		return ss_err.ACErrVerifySignFailed, false
	}

	ss_log.Info("%v|RSA2验签成功", traceNo)
	return "", true
}

// MD5验签
func (s *SignerMw) verifyMD5(c *gin.Context, traceNo string) (string, bool) {
	// todo 暂时不支持MD5验签
	return ss_err.ACErrSignTypeErr, false
}

// 签名
func (s *SignerMw) Sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 有设置跳过返回签名步骤
		if c.GetString(common.InnerSkipSign) == common.InnerSkipSignValue {
			ss_log.Info("%v|跳过返回签名步骤", traceNo)
			return
		}

		// 签名失败, 返回固定的格式，固定的异常错误
		retDataMap := inner_util.GetExceptionRetData()

		// 获取signType
		signType := inner_util.M(c, common.ParamsSignType)

		switch signType {
		case common.SignMethod_RSA2: // RSA2签名方式
			if data, ok := s.signRSA2(c, traceNo); ok {
				retDataMap = data
			}
		case common.SignMethod_MD5: // MD5签名方式
			if data, ok := s.signMD5(c, traceNo); ok {
				retDataMap = data
			}
		default:
			ss_log.Error("%v|参数signType为空或取配置错误,signType:%v", traceNo, signType)
		}

		// 重新将数据设置回去
		c.Set(common.RET_DATA_PRESEND, retDataMap)
	}
}

// 使用RSA2对数据进行签名
func (s *SignerMw) signRSA2(c *gin.Context, traceNo string) (gin.H, bool) {
	// 获取返回数据
	retData, exists := c.Get(common.RET_DATA_PRESEND)
	if !exists {
		ss_log.Error("%v|返回|未设置返回数据", traceNo)
		return nil, false
	}

	// 类型断言
	retDataMap, ok := retData.(gin.H)
	if !ok { // 转换失败
		ss_log.Error("%v|返回|返回数据格式转换失败|retData:%s", traceNo, strext.ToJson(retData))
		return nil, false
	}

	if len(retDataMap) == 0 { // 返回数据为空
		ss_log.Error("%v|返回|获取返回数据为空", traceNo)
		return nil, false
	}

	// 获取平台的私钥
	privKey := inner_util.S(c, common.InnerAppPlatformPrivateKey)
	if privKey == "" {
		ss_log.Error("%v|返回|获取平台的私钥失败", traceNo)
		return nil, false
	}

	// 将map数据转换为字符串，以便给RSA2进行签名
	data := ss_func.ParamsMapToString(retDataMap, common.RetFieldSgin)

	// 使用平台私钥签名
	sign, err := ss_rsa.RSA2Sign(data, privKey)
	if err != nil {
		ss_log.Error("%v|返回|使用平台私钥签名失败|err:%v|data:%v|privKey:%v", traceNo, err, data, privKey)
		return nil, false
	}

	retDataMap[common.RetFieldSgin] = sign

	ss_log.Info("%v|返回|签名完成", traceNo)
	return retDataMap, true
}

// 使用MD5对数据进行签名
func (s *SignerMw) signMD5(c *gin.Context, traceNo string) (gin.H, bool) {
	// todo 暂时不支持MD5签名
	return nil, false
}
