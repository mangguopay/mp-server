package middleware

import (
	_ "a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/common"
	"a.a/mp-server/api-webadmin/inner_util"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RespMw struct {
}

var RespMwInst RespMw

func (*RespMw) Resp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		//ss_log.Info("RespMw|Resp")
		ss_log.Info("resp=[%v]", inner_util.S(c, common.RET_DATA_PRESEND))
		c.Header("Content-Type", "application/json; charset=utf-8")
		statusCode := http.StatusOK
		if inner_util.S(c, common.RET_CODE) == ss_err.ERR_SYS_NO_ROUTE {
			statusCode = http.StatusNotFound
		}
		p := inner_util.S(c, common.RET_DATA_PRESEND)
		c.String(statusCode, p)
		// 打印
		traceNo := c.GetString(common.INNER_TRACE_NO)
		p2 := strext.Json2Map(p)
		ss_log.Info("%v|----------------------------返回的参数", traceNo)
		for k, v := range p2 {
			ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v)
		}
		ss_log.Info("%v|----------------------------", traceNo)
		return
	}
}

func (*RespMw) Pack() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		//ss_log.Info("JsonPackerMw|Pack")
		var data2 gin.H
		data, _ := c.Get(common.RET_DATA)
		switch data.(type) {
		case gin.H:
			data2 = data.(gin.H)
		default:
			data2 = nil
		}
		inner_util.R(c, inner_util.S(c, common.RET_CODE), inner_util.S(c, common.RET_MSG), data2)
	}
}
