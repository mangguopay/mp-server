package middleware

import (
	"database/sql"

	"a.a/mp-server/common/constants"

	"a.a/cu/ss_log"
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/api-pay/dao"
	"a.a/mp-server/api-pay/inner_util"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type VerifyAppMw struct {
}

var VerifyAppMwInst VerifyAppMw

func (VerifyAppMw) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceNo := c.GetString(common.INNER_TRACE_NO)

		appId := inner_util.M(c, common.ParamsAppId)
		if appId == "" {
			ss_log.Error("%v|缺少app_id参数", traceNo)
			c.Set(common.RET_CODE, ss_err.ACErrMissParamAppId)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		// 获取应用的配置信息
		app, err := dao.AppDaoInstance.GetSignInfo(appId)
		if err != nil {
			if err == sql.ErrNoRows { // 应用记录不存在
				ss_log.Error("%v|应用记录不存在|appId:%s", traceNo, appId)
				c.Set(common.RET_CODE, ss_err.ACErrAppNotExists)
			} else {
				ss_log.Error("%v|查询应用记录失败|appId:%s, err:%v", traceNo, appId, err)
				c.Set(common.RET_CODE, ss_err.ACErrSysErr)
			}
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		if app.Status != constants.BusinessAppStatus_Up { // 应用状态不是上架状态
			ss_log.Error("%v|应用状态不是上架状态|appId:%s, status:%v", traceNo, appId, app.Status)
			c.Set(common.RET_CODE, ss_err.ACErrAppNotPutOn)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		if app.SignMethod != common.SignMethod_RSA2 { // 签名方式配置错误
			ss_log.Error("%v|签名方式配置错误|appId:%s, SignMethod:%v", traceNo, appId, app.SignMethod)
			c.Set(common.RET_CODE, ss_err.ACErrAppConfigErr)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		if app.BusinessPubKey == "" { // 商家的公钥为空
			ss_log.Error("%v|商家的公钥为空|appId:%s", traceNo, appId)
			c.Set(common.RET_CODE, ss_err.ACErrAppConfigErr)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		if app.PlatformPrivKey == "" { // 平台的私钥为空
			ss_log.Error("%v|平台的私钥为空|appId:%s", traceNo, appId)
			c.Set(common.RET_CODE, ss_err.ACErrAppConfigErr)
			c.Set(common.InnerSkipSign, common.InnerSkipSignValue) // 跳过返回签名步骤
			c.Abort()
			return
		}

		ss_log.Info("%v|应用各项配置正确|appId:%s", traceNo, appId)

		c.Set(common.InnerAppSignMethod, app.SignMethod)
		c.Set(common.InnerAppIpWhiteList, app.IpWhiteList)
		c.Set(common.InnerAppBusinessPublickKey, app.BusinessPubKey)
		c.Set(common.InnerAppPlatformPrivateKey, app.PlatformPrivKey)
		return
	}
}
