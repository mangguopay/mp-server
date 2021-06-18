package inner_util

import (
	"a.a/cu/strext"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

func S(c *gin.Context, key string) string {
	val, _ := c.Get(key)
	return strext.ToStringNoPoint(val)
}

func M(c *gin.Context, key string) string {
	p, exists := c.Get(common.INNER_PARAM_MAP)
	if !exists {
		return ""
	}
	return strext.ToStringNoPoint(p.(map[string]interface{})[key])
}

// 接口组织数据返回异常时的默认返回数据
func GetExceptionRetData() gin.H {
	return gin.H{
		common.RetFieldCode: ss_err.ACErrSysRetAbnormal,
		common.RetFieldMsg:  ss_err.GetPayApiErrMsg(ss_err.ACErrSysRetAbnormal, constants.DefaultLang),
		common.RetFieldSgin: "",
	}
}
