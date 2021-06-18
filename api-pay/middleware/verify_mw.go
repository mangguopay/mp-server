package middleware

import (
	"strconv"
	"strings"

	"a.a/cu/ss_time"
	"a.a/mp-server/common/global"

	"a.a/cu/ss_log"

	"a.a/mp-server/common/ss_err"

	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/api-pay/inner_util"
	"github.com/gin-gonic/gin"
)

var (
	VerifyMwInst VerifyMw
)

type VerifyMw struct {
}

// IP白名单验证
func (v *VerifyMw) IpWhiteList() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 获取app的ip白名单
		ipWhiteList := inner_util.S(c, common.InnerAppIpWhiteList)

		if ipWhiteList == "" { // 没有设置IP白名单限制
			return
		}

		xRealIp := c.ClientIP()

		flagIsPass := false
		split := strings.Split(ipWhiteList, "\n")
		for _, v := range split {
			if strings.TrimSpace(v) == xRealIp {
				flagIsPass = true
				break
			}
		}

		if !flagIsPass {
			ss_log.Error("%v|IP白名单限制|ipWhiteList:[%s]|xRealIp:[%s]", traceNo, ipWhiteList, xRealIp)
			c.Set(common.RET_CODE, ss_err.ACErrIpWhiteListForbid)
			c.Abort()
			return
		}
		return
	}
}

// 校验真是ip
func (v *VerifyMw) RequestTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		// 获取请求的时间戳
		timestamp := inner_util.M(c, common.ParamsTimestamp)
		if timestamp == "" {
			ss_log.Error("%v|参数timestamp不存在或timestamp为空", traceNo)
			c.Set(common.RET_CODE, ss_err.ACErrMissParamTimestamp)
			c.Abort()
			return
		}

		// 解析timestamp为int64类型
		reqTime, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			ss_log.Error("%v|参数timestamp不存在或timestamp为空|timestamp=%s|err:%v", traceNo, timestamp, err)
			c.Set(common.RET_CODE, ss_err.ACErrParamTimestampErr)
			c.Abort()
			return
		}

		// 验证请求是否已经过期了
		if nowTime := ss_time.NowTimestamp(global.Tz); reqTime+common.RequestExpireTime < nowTime {
			ss_log.Error("%v|请求已过期|nowTime=%v|reqTime=%v|expire:%v", traceNo, nowTime, reqTime, common.RequestExpireTime)
			c.Set(common.RET_CODE, ss_err.ACErrRequestExpired)
			c.Abort()
			return
		}

		return
	}
}
