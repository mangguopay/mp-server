package middleware

import (
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/common"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_struct"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
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
		signBefore := fmt.Sprintf("x-login-token=%s&x-lang=%s&x-path=%s&x-ran=%s&key=%s", xLoginToken, xLang, xPath, xRan, passwordSalt)
		sign := encrypt.DoMd5(signBefore)
		ss_log.Info("%v|signBefore=[%v],md5=[%v]", traceNo, signBefore, sign)
		if isChkSign && strings.ToLower(sign) != strings.ToLower(strext.ToStringNoPoint(xSign)) {
			c.Set(ss_err.ERR_ACCOUNT_X_SIGN_FAILD, true)
			return // x-sign认证失败
		}

		cData := ss_struct.HeaderCommonData{}
		switch xLang {
		case constants.LangEnUS:
			fallthrough
		case constants.LangKmKH:
			fallthrough
		case constants.LangZhCN:
			cData.Lang = xLang
		default:
			cData.Lang = constants.LangEnUS // 默认参数
		}
		c.Set(constants.Common_Data, cData)
	}
}
