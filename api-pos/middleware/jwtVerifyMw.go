package middleware

import (
	"a.a/cu/jwt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/api-pos/common"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
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
			if c.Request.RequestURI != "/auth/logout" {
				isOk = jwt.ValidateEncryptedJWT(value, constants.MobileKeyJwtAes, constants.MobileKeyJwtSign)
			} else {
				isOk = true
			}
			//isOk := true
			ss_log.Error("%v|jwt token verify=[%v]", traceNo, isOk)
			if isOk {
				decodedJwt := jwt.DecodeEncryptedJWT(value, constants.MobileKeyJwtAes, constants.MobileKeyJwtSign)
				ss_log.Info("---------------------------jwt解密")
				for k, v := range decodedJwt {
					if k == "iat" || k == "exp" {
						ss_log.Info("%v|[%v]=>[%v]", traceNo, k, ss_time.Unixtime2Time(strext.ToStringNoPoint(v), global.Tz))
					} else {
						ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v)
					}
				}
				ss_log.Info("---------------------------")

				c.Set(common.INNER_JWT_DATA, decodedJwt)
				c.Set(common.INNER_IS_JWT_CHECKED, true)
				return
			}
		} else if c.Request.RequestURI == common.Login_URI {
			c.Set(common.INNER_IS_JWT_CHECKED, true)
			return
		}
	}
}
