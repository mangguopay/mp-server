package middleware

import (
	_ "a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/api-cb/inner_util"
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
		c.String(http.StatusOK, inner_util.S(c, common.RET_DATA_PRESEND))
		return
	}
}
