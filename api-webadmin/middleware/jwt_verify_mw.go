package middleware

import (
	"a.a/cu/container"
	"a.a/cu/jwt"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/common"
	"a.a/mp-server/api-webadmin/dao"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"strings"
)

type JwtVerifyMw struct {
}

var JwtVerifyMwInst JwtVerifyMw

/**
 * 认证中间件
 */
func (*JwtVerifyMw) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		traceNo := c.GetString(common.INNER_TRACE_NO)

		token := c.Request.Header.Get("Authorization")
		//ss_log.Info("toke=[%v]", token)
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

			xRealIp := c.Request.Header.Get("X-Real-IP")
			if xRealIp == "" {
				xRealIp = c.ClientIP()
			}
			c.Set("ip", xRealIp)
			//isOk := true
			ss_log.Error("%v|jwt token verify=[%v]", traceNo, isOk)
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
					xPath := c.GetString("xPath")
					xLoginToken := c.GetString("xLoginToken")
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
						c.Set(ss_err.ERR_ACCOUNT_NO_LOGIN, true)
						return // 未登录或认证已过期
					}
				}

				c.Set(common.INNER_JWT_DATA, decodedJwt)
				c.Set(common.INNER_IS_JWT_CHECKED, true)
				return
			}

			c.Set(ss_err.ERR_ACCOUNT_JWT_OUTDATED, true)
			c.Set(common.INNER_IS_JWT_CHECKED, false)
			return
		} else if c.Request.RequestURI == common.Login_URI {
			c.Set(common.INNER_IS_JWT_CHECKED, true)
			return
		}
	}
}
