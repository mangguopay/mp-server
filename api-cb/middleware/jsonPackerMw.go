package middleware

import (
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/api-cb/inner_util"
	"github.com/gin-gonic/gin"
)

type JsonPackerMw struct {
}

var JsonPackerMwInst JsonPackerMw

func (*JsonPackerMw) Pack() gin.HandlerFunc {
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
