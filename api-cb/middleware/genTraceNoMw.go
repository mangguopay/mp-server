package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-cb/common"
	"github.com/gin-gonic/gin"
)

type GenTraceNoMw struct {
}

var GenTraceNoMwInst GenTraceNoMw

/**
 * 生成跟踪号
 */
func (GenTraceNoMw) GenTraceNo() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := strext.GetDailyId()
		c.Set(common.INNER_TRACE_NO, traceNo)
		ss_log.Info("trace[%v]begin====================", traceNo)
	}
}
