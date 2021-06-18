package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/common/constants"
	"github.com/gin-gonic/gin"
)

type GenTraceNoMw struct {
}

var GenTraceNoMwInst GenTraceNoMw

/**
 * 生成跟踪号
 */
func (g *GenTraceNoMw) GenTraceNo() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := strext.GetDailyId()

		ss_log.Info("%v|start...", traceNo)

		c.Set(common.INNER_TRACE_NO, traceNo)

		// 设置默认语言类型到gin
		c.Set(common.InnerLanguage, constants.DefaultLang)
	}
}
