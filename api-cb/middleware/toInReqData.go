package middleware

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/api-cb/dao"
	"a.a/mp-server/api-cb/inner_util"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

type ToInReqDataMw struct {
}

var ToInReqDataMwInst ToInReqDataMw

/**
 *
 */
func (ToInReqDataMw) DoTrans() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ss_log.Info("DoTrans")
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================

		accNo := inner_util.M(c, common.INNER_DATA_ACCNO)
		if accNo == "" {
			ss_log.Info("DoTrans|no account")
			c.Set(common.RET_CODE, ss_err.ERR_SYS_NO_ACCNO)
			c.Set(common.INNER_IS_STOP, true)
			return
		}

		// 提前拿商户配置
		rKey, signMethod, apiMode, err := dao.ApiDaoInstance.GetMercSignInfo(accNo)
		if err != nil {
			ss_log.Error("no rKey=[%v]|err=[%v]", rKey, err)
			c.Set(common.RET_CODE, ss_err.ERR_SYS_API_SIGN_INFO)
			c.Set(common.INNER_IS_STOP, true)
			return
		}
		c.Set(common.INNER_SIGN_METHOD, signMethod)
		c.Set(common.INNER_API_MODE, apiMode)
		return
	}
}
