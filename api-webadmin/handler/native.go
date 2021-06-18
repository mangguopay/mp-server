package handler

import (
	"a.a/mp-server/api-webadmin/dao"
	"a.a/mp-server/api-webadmin/inner_util"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
)

//
func (a *AdminAuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNo := inner_util.GetJwtDataString(c, "account_uid")
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			dao.AccDaoInstance.DeleteLoginToken(accountNo)
			return ss_err.ERR_SUCCESS, "", nil
		})
	}
}
