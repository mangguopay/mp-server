package middleware

import (
	colloection "a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/common"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RespMiddleWare struct {
}

func (*RespMiddleWare) DoResp(filter []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)
		// 开始时间
		//start := time.Now()
		path := c.Request.URL.Path

		method := c.Request.Method
		//var params string
		ss_log.Info("%v|path=[%v]", traceNo, path)

		if http.MethodGet == method {
			queryForm, _ := url.ParseQuery(c.Request.URL.RawQuery)
			p := map[string]interface{}{}
			ss_log.Info("%v|----------------------------GET的参数", traceNo)
			for k, v := range queryForm {
				ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v[0])
				p[k] = v[0]
			}
			ss_log.Info("%v|----------------------------", traceNo)
			c.Set("params", p)
		} else if colloection.GetKey(path, filter) > -1 {
			ss_log.Info("%v|Index----------------------------filter=%v,path=%v", traceNo, filter, path)
		} else {
			buf, _ := ioutil.ReadAll(c.Request.Body)
			defer c.Request.Body.Close()
			//params = string(buf)
			var p map[string]interface{}
			json.Unmarshal(buf, &p)
			c.Set("params", p)
			ss_log.Info("%v|----------------------------POST的参数", traceNo)
			for k, v := range p {
				ss_log.Info("%v|[%v]=>[%v]", traceNo, k, strext.ForShort(strext.ToStringNoPoint(v), 32))
			}
			ss_log.Info("%v|----------------------------", traceNo)
		}

		// Process request
		c.Next()

		_, isStop := c.Get("is_stop")
		if isStop {
			return
		}

		statusCode, exStatus := c.Get("status")
		if exStatus {
			c.Writer.WriteHeader(int(statusCode.(int)))
		} else if c.Writer.Status() > 0 {
			c.Writer.WriteHeader(c.Writer.Status())
			statusCode = c.Writer.Status()
		}

		resp, exStatus := c.Get("resp")
		//ss_log.Info("resp=[%v]", resp)
		ss_log.Info("%v|----------------------------响应的数据", traceNo)
		switch resp.(type) {
		case gin.H:
			for k, v := range resp.(gin.H) {
				switch v2 := v.(type) {
				case gin.H:
					for k3, v3 := range v2 {
						ss_log.Info("%v|[data|%v]=>[%v]", traceNo, k3, v3)
					}
				case map[string]interface{}:
					for k3, v3 := range v2 {
						ss_log.Info("%v|[data|%v]=>[%v]", traceNo, k3, v3)
					}
				default:
					ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v)
				}
			}
		default:
			ss_log.Info("%v|resp=[%v]", traceNo, resp)
		}
		ss_log.Info("%v|----------------------------", traceNo)

		if exStatus {
			c.JSON(statusCode.(int), resp)
			(resp.(gin.H))["data"] = gin.H{}
		} else {
			resp = ""
		}

		// Stop timer
		//end := time.Now()
		//latency := end.Sub(start)

		//clientIP := c.ClientIP()
		//xPath, _ := c.Get("xPath")
		//accNo, _ := c.Get("acc_no")
		//
		////comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		//during := int32(latency.Nanoseconds() / 1000)
		//
		//statusCode2 := strconv.Itoa(c.Writer.Status())
		//if colloection.GetKey(path, filter) > -1 {
		//	rpc_cli.RpcCliMemberInst.InsertAdminLog(&member.InsertAdminLogRequest{
		//		CreateTime: start.Format("2006-01-02 15:04:05.999999-07:00"),
		//		During:     during, Url: path, Param: strext.ToStringNoPoint(xPath), Op: method, OpType: 1, OpAccUid: strext.ToStringNoPoint(accNo),
		//		Ip: clientIP, StatusCode: statusCode2, Response: "",
		//	})
		//} else {
		//	rpc_cli.RpcCliMemberInst.InsertAdminLog(&member.InsertAdminLogRequest{CreateTime: start.Format("2006-01-02 15:04:05.999999-07:00"),
		//		During: during, Url: path, Param: strext.ToStringNoPoint(xPath), Op: method, OpType: 1, OpAccUid: strext.ToStringNoPoint(accNo),
		//		Ip: clientIP, StatusCode: statusCode2, Response: "",
		//	})
		//}
		return
	}
}
