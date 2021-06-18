package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/url"
)

type GetParamsMw struct {
}

var GetParamsMwInst GetParamsMw

/**
 * 获取post过来的json
 */
func (GetParamsMw) GetPostJsonBodyParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		if err != nil {
			ss_log.Error("请求包体为空|err=[%v]", err)
			c.Set(common.INNER_IS_STOP, true)
			c.Set(common.INNER_ERR_REASON, ss_err.ERR_SYS_EMPTY_BODY)
			c.Abort()
			return
		}

		p := strext.Json2Map(buf)
		if p == nil {
			ss_log.Error("请求包体不是json|buf=[%v]", string(buf))
			reqFormStr2, _ := url.QueryUnescape(string(buf))
			p = strext.FormString2Map(reqFormStr2)
			if p == nil {
				ss_log.Error("请求包体不是form|buf=[%v]", reqFormStr2)
				c.Set(common.INNER_IS_STOP, true)
				c.Set(common.INNER_ERR_REASON, ss_err.ERR_SYS_BODY_NOT_JSON)
				c.Abort()
				return
			}
		}

		ss_log.Info("----------------------------POST的参数")
		for k, v := range p {
			ss_log.Info("|[%v]=>[%v]", k, v)
		}
		ss_log.Info("----------------------------")

		c.Set(common.INNER_FMT, "json")
		ss_log.Info("recv[%v]=>[%v]", c.Request.Method, strext.ToStringNoPoint(c.Request.RequestURI))
		c.Set(common.INNER_PARAM_MAP, p)
		return
	}
}
