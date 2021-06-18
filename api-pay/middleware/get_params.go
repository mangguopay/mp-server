package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"a.a/cu/util"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"github.com/gin-gonic/gin"
)

type GetParamsMw struct {
}

var GetParamsMwInst GetParamsMw

// 检测路由信息是否匹配
func (g *GetParamsMw) CheckRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取跟踪号
		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 因gin框架判断路由是否匹配是在中间件执行完之后才判断的
		// 固在此进行提前判断
		routeList := []string{
			"/api/prepay",
			"/api/query",
			"/api/refund",
			"/api/query_refund",
			"/api/transfer/enterprise",
			"/api/transfer/query",
		}

		if !util.InSlice(c.Request.URL.Path, routeList) {
			ss_log.Error("%v|路由不匹配,path=[%v]", traceNo, c.Request.URL.Path)
			c.Set(common.RET_CODE, ss_err.ACErrSysRouteNotFound)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}
	}
}

// 读取get的参数
func (g *GetParamsMw) GetParamsFromGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	只处理get
		if http.MethodGet != c.Request.Method {
			return
		}

		// 读取跟踪号
		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 读取get的参数
		queryForm, _ := url.ParseQuery(c.Request.URL.RawQuery)
		p := map[string]interface{}{}
		for k, v := range queryForm {
			p[k] = v[0]
		}

		ss_log.Info("%v|uri=[%v]|GET的参数:%v", traceNo, strext.ToStringNoPoint(c.Request.RequestURI), strext.ToJson(p))

		// 设置传入参数到gin
		c.Set(common.INNER_PARAM_MAP, p)

		// 获取语言类型
		lang := fmt.Sprintf("%s", p[common.ParamsLang])

		// 设置语言类型到gin
		c.Set(common.InnerLanguage, ss_func.NormalizeLang(lang))
		return
	}
}

// 获取post过来的json
func (g *GetParamsMw) GetParamsFromPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		//只处理post
		if http.MethodPost != c.Request.Method {
			return
		}

		traceNo := c.GetString(common.INNER_TRACE_NO)

		body, err := ioutil.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()

		c.Set(common.REQUEST_BODY, strext.Json2Map(body))

		if err != nil {
			ss_log.Error("%v|读取body失败|err:%v", traceNo, err)
			c.Set(common.RET_CODE, ss_err.ACErrPostReadBodyErr)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		if len(body) == 0 {
			ss_log.Error("%v|请求body为空", traceNo)
			c.Set(common.RET_CODE, ss_err.ACErrPostEmptyBody)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		p := strext.Json2Map(body)
		if p == nil {
			ss_log.Error("%v|请求包体不是json|buf=[%v]", traceNo, string(body))
			c.Set(common.RET_CODE, ss_err.ACErrPostBodyNotJson)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		ss_log.Info("%v|uri=[%v]|POST的参数:%v", traceNo, strext.ToStringNoPoint(c.Request.RequestURI), strext.ToJson(p))

		// 设置传入参数到gin
		c.Set(common.INNER_PARAM_MAP, p)

		// 获取语言类型
		lang := fmt.Sprintf("%s", p[common.ParamsLang])

		// 设置语言类型到gin
		c.Set(common.InnerLanguage, ss_func.NormalizeLang(lang))
		return
	}
}

// 检查参数
func (g *GetParamsMw) CheckRequestMethod() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 暂时只支持post请求
		if !util.InSlice(c.Request.Method, []string{http.MethodPost}) {
			c.Set(common.RET_CODE, ss_err.ACErrSysMethodErr)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		return
	}
}
