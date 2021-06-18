package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/api-pay/inner_util"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"github.com/gin-gonic/gin"
)

type JsonPackerMw struct {
}

var JsonPackerMwInst JsonPackerMw

func (*JsonPackerMw) Pack() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 获取返回code
		retCode := inner_util.S(c, common.RET_CODE)
		if retCode == "" { // 没有设置返回code
			ss_log.Info("%v|返回|未设置返回code", traceNo)
			retCode = ss_err.ACErrSysErr
		}

		// 获取返回的data数据
		data := gin.H{}
		if retData, exists := c.Get(common.RET_DATA); exists { // 有设置返回数据
			if d, ok := retData.(gin.H); ok { // 返回数据的格式正确
				data = d
			} else { // 返回数据的格式错误
				ss_log.Error("%v|返回|返回数据格式不正确|retData:%s", traceNo, strext.ToJson(retData))
			}
		}

		// 获取语言类型
		lang := ss_func.NormalizeLang(c.GetString(common.InnerLanguage))
		ss_log.Info("%v|返回|retCode=%v|lang=%v", traceNo, retCode, lang)

		// 将code和msg字段合并到返回数据中
		data[common.RetFieldCode] = retCode
		data[common.RetFieldMsg] = ss_err.GetPayApiErrMsg(retCode, lang)

		// 设置用户返回的数据到gin,以便下一步返回的处理
		c.Set(common.RET_DATA_PRESEND, data)

		ss_log.Info("%v|返回|打包完成", traceNo)
	}
}
