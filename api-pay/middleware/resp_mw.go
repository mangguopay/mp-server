package middleware

import (
	"net/http"

	"a.a/mp-server/api-pay/dao"
	"a.a/mp-server/api-pay/inner_util"
	"a.a/mp-server/common/ss_err"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-pay/common"
	"github.com/gin-gonic/gin"
)

type RespMw struct {
}

var RespMwInst RespMw

func (*RespMw) Resp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 获取返回数据
		retData, exists := c.Get(common.RET_DATA_PRESEND)
		if !exists {
			ss_log.Error("%v|返回|返回数据未设置,使用默认异常返回数据", traceNo)
			retData = inner_util.GetExceptionRetData()
		}

		ss_log.Info("%v|返回|end:%s", traceNo, strext.ToJson(retData))

		go recordLog(c, traceNo)

		c.Header("Content-Type", "application/json; charset=utf-8")

		c.String(http.StatusOK, strext.ToStringNoPoint(retData))
		return
	}
}

func recordLog(c *gin.Context, traceNo string) {
	var trafficStatus, businessStatus, reqTime string
	reqBody, exists := c.Get(common.REQUEST_BODY)
	if exists {
		reqTime = strext.ToString(reqBody.(map[string]interface{})["timestamp"])
	}

	retData, exists := c.Get(common.RET_DATA_PRESEND)
	if exists {
		if strext.ToString(retData.(gin.H)["code"]) != ss_err.ACErrSuccess {
			trafficStatus = "0"
		} else {
			trafficStatus = "1"
		}

		if strext.ToString(retData.(gin.H)["sub_code"]) != ss_err.Success {
			businessStatus = "0"
		} else {
			businessStatus = "1"
		}
	}

	log := new(dao.LogApiPayDao)
	log.TraceNo = traceNo
	log.AppId = inner_util.M(c, "app_id")
	log.ReqUri = c.Request.RequestURI
	log.ReqMethod = c.Request.Method
	log.ReqBody = strext.ToJson(reqBody)
	log.RespDate = strext.ToJson(retData)
	log.ReqTime = strext.ToInt64(reqTime)
	log.TrafficStatus = trafficStatus
	log.BusinessStatus = businessStatus
	if err := dao.LogApiPayDaoInst.Insert(log); err != nil {
		ss_log.Info("%v|插入日志失败|%v, err=%v", traceNo, strext.ToJson(log), err)
	}
}
