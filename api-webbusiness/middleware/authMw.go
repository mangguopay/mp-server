package middleware

import (
	"net/http"

	"a.a/cu/ss_lang"
	"a.a/cu/ss_log"
	"a.a/mp-server/api-webbusiness/common"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
)

type AuthMw struct {
}

var AuthMwInst AuthMw

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

			// 获取客户端语种,并检查语种是否存在，不存在会使用默认语种
			lang := ss_lang.NormalLaguage(ss_net.GetCommonData(c).Lang)

			c.JSON(http.StatusUnauthorized, gin.H{
				"retcode": ss_err.ERR_SYS_SIGN_JWT,
				"msg":     ss_err.GetErrMsgMulti(lang, ss_err.ERR_SYS_SIGN_JWT),
				"status":  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}
	}
}
