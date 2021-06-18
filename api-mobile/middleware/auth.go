package middleware

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/api-mobile/common"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type AuthMw struct {
}

var AuthMwInst AuthMw

/**
 * 认证中间件
 */
func (*AuthMw) DoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		// rsa验签
		if !c.GetBool(common.INNER_SIGN_VERIFY) {
			c.Set(common.RET_CODE, ss_err.ERR_SYS_SIGN)
			c.Set(common.INNER_IS_STOP, true)
			RespMwInst.PackInner(c)
			RsaMwInst.EncodeInner(c)
			RsaMwInst.SignInner(c)
			RespMwInst.RespInner(c)
			c.Abort()
			return
		}
	}
}

func (*AuthMw) DoAuthJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		// jwt验签
		if !c.GetBool(common.INNER_IS_JWT_CHECKED) {
			c.Set(common.RET_CODE, ss_err.ERR_SYS_SIGN_JWT)
			c.Set(common.INNER_IS_STOP, true)
			RespMwInst.PackInner(c)
			RsaMwInst.EncodeInner(c)
			RsaMwInst.SignInner(c)
			RespMwInst.RespInner(c)
			c.Abort()
			return
		}
	}
}
