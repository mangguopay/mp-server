package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"a.a/cu/ss_lang"
	"a.a/mp-server/common/ss_net"

	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/mp-server/api-webbusiness/common"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type XPathVerifyMw struct {
}

var XPathVerifyMwInst XPathVerifyMw

/**
 * 认证中间件
 */
func (*XPathVerifyMw) Verify(isChkSign bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		traceNo := c.GetString(common.INNER_TRACE_NO)

		xPath := c.Request.Header.Get("x-path")
		c.Set("xPath", xPath)
		xSign := c.Request.Header.Get("x-sign")
		xRan := c.Request.Header.Get("x-ran")
		xLoginToken := c.Request.Header.Get("x-login-token")
		xLang := c.Request.Header.Get("x-lang")

		c.Set("xLoginToken", xLoginToken)
		ss_log.Info("%v|xPath=[%v],xRan=[%v],xLoginToken=[%v],xLang=[%v],xSign=[%v]", traceNo, xPath, xRan, xLoginToken, xLang, xSign)

		k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		}
		signBefore := fmt.Sprintf("x-login-token=%s&x-path=%s&x-ran=%s&key=%s&x-lang=%s", xLoginToken, xPath, xRan, passwordSalt, xLang)
		sign := encrypt.DoMd5(signBefore)
		ss_log.Info("%v|signBefore=[%v],md5Sign=[%v]", traceNo, signBefore, sign)
		verifySign := strings.ToLower(sign) != strings.ToLower(xSign)
		if isChkSign && verifySign {
			ss_log.Error("x-sign认证失败, isChkSign=%v, verifySign=%v", isChkSign, !verifySign)
			c.Set(common.RET_CODE, ss_err.ERR_ACCOUNT_X_SIGN_FAILD)
			//c.Set(common.INNER_IS_STOP, true)

			// 获取客户端语种,并检查语种是否存在，不存在会使用默认语种
			lang := ss_lang.NormalLaguage(ss_net.GetCommonData(c).Lang)

			c.JSON(http.StatusUnauthorized, gin.H{
				"retcode": ss_err.ERR_ACCOUNT_X_SIGN_FAILD,
				"msg":     ss_err.GetErrMsgMulti(lang, ss_err.ERR_ACCOUNT_X_SIGN_FAILD),
				"status":  http.StatusUnauthorized,
			})
			c.Abort()
			return // x-sign认证失败
		}

	}
}
