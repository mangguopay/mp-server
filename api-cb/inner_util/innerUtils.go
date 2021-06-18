package inner_util

import (
	"a.a/cu/ss_lang"
	"a.a/cu/strext"
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"a.a/net/consts"
	"github.com/gin-gonic/gin"
)

func S(c *gin.Context, key string) string {
	val, _ := c.Get(key)
	return strext.ToStringNoPoint(val)
}

func M(c *gin.Context, key string) string {
	p, _ := c.Get(common.INNER_PARAM_MAP)
	return strext.ToStringNoPoint(p.(map[string]interface{})[key])
}

// 返回值
func R(c *gin.Context, retCode, retMsg string, data gin.H) {
	if retCode == "" {
		retCode = ss_err.ERR_SYS_UNKNOWN
	}
	var msg string
	if retMsg == "" {
		// 获取客户端语种,并检查语种是否存在，不存在会使用默认语种
		lang := ss_lang.NormalLaguage(ss_net.GetCommonData(c).Lang)

		if retCode == ss_err.ERR_SUCCESS {
			msg = ss_err.GetErrMsgMulti(lang, retCode)
		} else {
			switch retCode[1:2] {
			case "A":
				msg = ss_err.GetErrMsgMulti(lang, retCode)
			case "B":
				msg = consts.GetErrMsg(retCode)
			}
		}
	} else {
		msg = retMsg
	}
	if retCode == ss_err.ERR_SUCCESS {
		c.Set(common.RET_DATA_PRESEND, strext.ToJsonNotChange(gin.H{
			"code": retCode,
			"msg":  msg,
			"data": data,
		}))
	} else {
		c.Set(common.RET_DATA_PRESEND, strext.ToJsonNotChange(gin.H{
			"code": retCode,
			"msg":  msg,
		}))
	}
}
