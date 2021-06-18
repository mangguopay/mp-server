package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-pos/common"
	"a.a/mp-server/api-pos/dao"
	"a.a/mp-server/api-pos/inner_util"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type LoginOutMw struct {
}

var LoginOutMwInst LoginOutMw

/**
 * 认证中间件
 */
func (*LoginOutMw) DoLoginOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//=========================================
		// 获取 token
		xLoginToken := c.Request.Header.Get("LoginToken")
		if xLoginToken == "" {
			ss_log.Error("1")
			c.Set(common.RET_CODE, ss_err.ERR_ACCOUNT_NO_LOGIN)
			c.Set(common.INNER_IS_STOP, true)
			RespMwInst.PackInner(c)
			RsaMwInst.EncodeInner(c)
			RsaMwInst.SignInner(c)
			RespMwInst.RespInner(c)
			c.Abort()
			return
		}
		// 从 jwt 里获取用户id
		uid := inner_util.GetJwtDataString(c, "account_uid")
		isOk := dao.AccDaoInstance.GetPosLoginToken(uid, strext.ToStringNoPoint(xLoginToken), c.ClientIP())
		if !isOk {
			ss_log.Error("2")
			c.Set(common.RET_CODE, ss_err.ERR_ACCOUNT_NO_LOGIN)
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
