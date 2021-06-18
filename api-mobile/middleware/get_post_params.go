package middleware

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/common"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GetParamsMw struct {
}

var GetParamsMwInst GetParamsMw

/**
 * 获取post过来的json
 */
func (GetParamsMw) FetchPostJsonBodyParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		//	只处理post
		if http.MethodPost != c.Request.Method && http.MethodDelete != c.Request.Method && http.MethodPatch != c.Request.Method {
			ss_log.Info("not post")
			return
		}
		ss_log.Info("post")

		traceNo := c.GetString(common.INNER_TRACE_NO)
		c.Set(common.INNER_IS_ENCODED, true)
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
			ss_log.Error("%v|请求包体不是json|buf=[%v]", traceNo, buf)
			c.Set(common.INNER_IS_STOP, true)
			c.Set(common.INNER_ERR_REASON, ss_err.ERR_SYS_BODY_NOT_JSON)
			c.Abort()
			return
		}

		ss_log.Info("%v|----------------------------POST的参数", traceNo)
		for k, v := range p {
			ss_log.Info("%v|[%v]=>[%v]", k, v, traceNo)
		}
		ss_log.Info("%v|----------------------------", traceNo)

		c.Set(common.INNER_FMT, "json")
		c.Set(common.INNER_PARAM_MAP, p)
		ss_log.Info("%v|recv[%v]=>[%v]", traceNo, c.Request.Method, strext.ToStringNoPoint(c.Request.RequestURI))
		return
	}
}

/**
 * 读取get的参数
 */
func (GetParamsMw) FetchGetParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		//	只处理get
		if http.MethodGet != c.Request.Method {
			return
		}
		// 读取跟踪号
		traceNo := c.GetString(common.INNER_TRACE_NO)
		// 读取get的参数
		queryForm, _ := url.ParseQuery(c.Request.URL.RawQuery)
		p := map[string]interface{}{}
		ss_log.Info("%v|----------------------------GET的参数", traceNo)
		for k, v := range queryForm {
			ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v[0])
			p[k] = v[0]
		}
		ss_log.Info("%v|----------------------------", traceNo)
		c.Set(common.INNER_PARAM_MAP, p)
		c.Set(common.INNER_IS_ENCODED, false)
		c.Set(common.INNER_SIGN_VERIFY, true)
		ss_log.Info("recv[%v]=>[%v]", c.Request.Method, strext.ToStringNoPoint(c.Request.RequestURI))
		return
	}
}
