package middleware

import (
	"a.a/cu/ss_log"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"runtime/debug"
)

type RecoveryMiddleWare struct {
}

var RecoveryMiddleWareInst RecoveryMiddleWare

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func (RecoveryMiddleWare) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ss_log.Info("RecoveryMiddleWare|Recovery")
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				ss_log.Error("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, "")
				c.AbortWithStatus(500)
			}
		}()
	}
}
