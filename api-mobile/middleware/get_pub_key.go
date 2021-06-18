package middleware

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/common"
	"a.a/mp-server/api-mobile/inner_util"
	"github.com/gin-gonic/gin"
	"strings"
)

type GetPriPubKeyMw struct {
}

var GetPriPubKeyMwInst GetPriPubKeyMw

func (*GetPriPubKeyMw) GetPubKeyMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		var pubKey string
		if strings.HasPrefix(uri, "/mobile/auth") { // 获取默认的pub
			pubKey = strext.ToStringNoPoint(common.EncryptMap["default_pub_key"])
		} else {
			// 从jwt里获取
			pubKey = inner_util.GetJwtDataString(c, "pub_key")
		}
		traceNo := c.GetString(common.INNER_TRACE_NO)
		ss_log.Info("%v|对方公钥=[%s]", traceNo, pubKey)
		c.Set(common.Pub_Key, pubKey)
	}
}

func (*GetPriPubKeyMw) GetPubKeyRetMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		params, _ := c.Get("params")
		pubKey := container.GetValFromMapMaybe(params, "pub_key").ToStringNoPoint()
		c.Set(common.Pub_Key, pubKey)
	}
}

// 未登陆换公钥
func (*GetPriPubKeyMw) UpdatePubkeyWhenLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		if strings.HasPrefix(uri, "/mobile/auth") { // 获取默认的pub
			params, _ := c.Get("params")
			pubKey := container.GetValFromMapMaybe(params, "pub_key").ToStringNoPoint()
			c.Set(common.Pub_Key, pubKey)
		} else {
			return
		}
	}
}
