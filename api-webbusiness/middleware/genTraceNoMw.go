package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/api-webbusiness/common"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_struct"
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

		xLang := c.Request.Header.Get("x-lang")

		ss_log.Info("xLang[%v]====================", xLang)

		// 如果未传则设置默认参数
		if xLang == "" || !util.InSlice(xLang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
			xLang = constants.LangEnUS
		}

		cData := ss_struct.HeaderCommonData{}
		cData.Lang = xLang
		c.Set(constants.Common_Data, cData)
	}
}
