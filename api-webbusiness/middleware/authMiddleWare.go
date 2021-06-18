package middleware

import (
	"a.a/cu/container"
	"a.a/cu/encrypt"
	"a.a/cu/jwt"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webbusiness/common"
	"a.a/mp-server/api-webbusiness/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/ss_err"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/**
 * 认证中间件
 */
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)
		token := c.Request.Header.Get("Authorization")
		xPath := c.Request.Header.Get("x-path")
		c.Set("xPath", xPath)
		xSign := c.Request.Header.Get("x-sign")
		xRan := c.Request.Header.Get("x-ran")
		xLoginToken := c.Request.Header.Get("x-login-token")
		ss_log.Info("%v|xPath=[%v],xRan=[%v],xLoginToken=[%v],xSign=[%v]", traceNo, xPath, xRan, xLoginToken, xSign)

		k1, passwordSalt, err := cache.ApiDaoInstance.GetGlobalParam("password_salt")
		if err != nil {
			ss_log.Error("err=[%v],missing key=[%v]", err, k1)
		}
		signBefore := fmt.Sprintf("x-login-token=%s&x-path=%s&x-ran=%s&key=%s", xLoginToken, xPath, xRan, passwordSalt)
		sign := encrypt.DoMd5(signBefore)
		ss_log.Info("%v|signBefore=[%v],md5=[%v]", traceNo, signBefore, sign)
		// xxx 测试暂时屏蔽
		//if strings.ToLower(sign) != strings.ToLower(strext.ToStringNoPoint(xSign)) {
		//	c.JSON(http.StatusUnauthorized, gin.H{
		//		"error": "x-sign认证失败",
		//	})
		//	c.Abort()
		//	return
		//}

		xRealIp := c.Request.Header.Get("X-Real-IP")
		if xRealIp == "" {
			xRealIp = c.ClientIP()
		}
		c.Set("ip", xRealIp)

		if strings.HasPrefix(token, "Bearer ") {
			value := token[7:] // 截取前面的prefix
			//fmt.Println(value)
			// 验证jwt
			//ss_log.Error("token=[%v]\n", value)
			isOk := false
			k1, loginSignKey, err := cache.ApiDaoInstance.GetGlobalParam("login_sign_key")
			if err != nil {
				ss_log.Error("%v|err=[%v],missing key=[%v]", traceNo, err, k1)
			}
			k1, loginAesKey, err := cache.ApiDaoInstance.GetGlobalParam("login_aes_key")
			if err != nil {
				ss_log.Error("%v|err=[%v],missing key=[%v]", traceNo, err, k1)
			}
			if c.Request.RequestURI != "/auth/logout" {
				isOk = jwt.ValidateEncryptedJWT(value, loginAesKey, loginSignKey)
			} else {
				isOk = true
			}
			//isOk := true
			ss_log.Error("isOk=[%v]", isOk)
			if isOk {
				decodedJwt := jwt.DecodeEncryptedJWT(value, loginAesKey, loginSignKey)
				ss_log.Info("---------------------------jwt解密")
				for k, v := range decodedJwt {
					ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v)
				}
				ss_log.Info("---------------------------")
				accNo := container.GetValFromMapMaybe(decodedJwt, "account_uid").ToStringNoPoint()
				loginAccNo := container.GetValFromMapMaybe(decodedJwt, "login_account_no").ToStringNoPoint()
				c.Set("acc_no", accNo)
				if c.Request.RequestURI != "/auth/logout" {
					if loginAccNo != "" {
						accNo = loginAccNo
					}
					isOk2, isReset := dao.AccDaoInstance.GetLoginToken(accNo, strext.ToStringNoPoint(xLoginToken), xPath, xRealIp)
					//isReset = true
					if isReset {
						err := dao.AccDaoInstance.UpdateLoginTime(accNo)
						if err != ss_err.ERR_SUCCESS {
							ss_log.Error("%v|err=[%v]", traceNo, err)
						} else {
							jwt2 := jwt.GetNewEncryptedJWTToken(-1, jwt.MapClaims2Map(&decodedJwt), loginAesKey, loginSignKey)
							isOk2 = true
							c.Header("x-token", jwt2)
						}
					}
					if !isOk2 {
						c.JSON(http.StatusUnauthorized, gin.H{
							"error": "未登录或认证已过期",
						})
						c.Abort()
						return
					}
				}

				c.Set("decodedJwt", decodedJwt)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未登录或认证已过期",
		})
		c.Abort()
		return
	}
}
