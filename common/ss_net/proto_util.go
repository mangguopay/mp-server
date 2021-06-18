package ss_net

import (
	"a.a/cu/ss_lang"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_struct"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	RET_CODE       = "_ret_code"
	RET_MSG        = "_ret_msg"
	RET_DATA       = "_ret_data"
	RET_CUSTOM_MSG = "_ret_custom_msg"
)

func HandleRet(retcode, successCode string, errCodeFunc func(string, ...interface{}) string, okFunc func(*gin.Context, ...interface{}) interface{},
	errFunc func(string, *gin.Context, ...interface{}) interface{},
	c *gin.Context, errMsg []interface{}, reply ...interface{}) {
	c.Set("status", http.StatusOK)
	if retcode == successCode {
		c.Set(RET_CODE, retcode)
		c.Set(RET_DATA, okFunc(c, reply...))
	} else {
		var msg string
		if errMsg == nil {
			msg = errCodeFunc(retcode)
		} else {
			if len(errMsg) >= 2 {
				if errMsg[0] == "1" {
					msg = errCodeFunc(retcode, errMsg...)
				} else {
					msg = strext.ToStringNoPoint(errMsg[1])
				}
			} else {
				msg = errCodeFunc(retcode, errMsg...)
			}
		}

		c.Set(RET_CODE, retcode)
		c.Set(RET_MSG, msg)
		c.Set(RET_DATA, errFunc(retcode, c, reply...))
	}
	return
}

// 返回retcode经过语言转换
func HandleRetMulti(retcode, successCode string, errCodeFunc func(ss_lang.SsLang, string, ...interface{}) string, okFunc func(*gin.Context, ...interface{}) interface{},
	errFunc func(string, *gin.Context, ...interface{}) interface{},
	c *gin.Context, errMsg []interface{}, reply ...interface{}) {
	c.Set("status", http.StatusOK)
	if retcode == successCode {
		c.Set(RET_CODE, retcode)
		c.Set(RET_DATA, okFunc(c, reply...))

	} else {
		// 获取客户端语种,并检查语种是否存在，不存在会使用默认语种
		lang := ss_lang.NormalLaguage(GetCommonData(c).Lang)

		var msg string
		if errMsg == nil {
			msg = errCodeFunc(lang, retcode)
		} else {
			if len(errMsg) >= 2 {
				if errMsg[0] == "1" {
					msg = errCodeFunc(lang, retcode, errMsg...)
				} else {
					msg = strext.ToStringNoPoint(errMsg[1])
				}
			} else {
				msg = errCodeFunc(lang, retcode, errMsg...)
			}
		}

		c.Set(RET_CODE, retcode)
		c.Set(RET_MSG, msg)
		c.Set(RET_DATA, errFunc(retcode, c, reply...))
	}
	return
}

func RetListFunc(c2 *gin.Context, r ...interface{}) interface{} {
	return gin.H{
		"data":  r[0],
		"total": r[1],
	}
}

func RetSingleFunc(c2 *gin.Context, r ...interface{}) interface{} {
	return gin.H{
		"data": r[0],
	}
}

func RetSingleUidFunc(c2 *gin.Context, r ...interface{}) interface{} {
	return gin.H{
		"uid": r[0],
	}
}

func RetMsgFunc(c2 *gin.Context, r ...interface{}) interface{} {
	return gin.H{}
}

func RetCommonFunc(c2 *gin.Context, r ...interface{}) interface{} {
	return r
}

func RetNormalErrFunc(retCode string, c2 *gin.Context, r ...interface{}) interface{} {
	return gin.H{}
}

func RetList2Func(c2 *gin.Context, r ...interface{}) interface{} {
	l := gin.H{
		"data":  r[0],
		"total": r[1],
	}

	for k, v := range r[2].(map[string]interface{}) {
		l[k] = v
	}

	return l
}

//===============================
//func ChkEmptyAndAbort(c *gin.Context, data interface{}, errCode string, errCodeFunc func(string, ...interface{}) string, inArgs ...interface{}) {
//	var msg string
//	if len(inArgs) < 1 {
//		msg = errCodeFunc(errCode)
//	} else {
//		msg = errCodeFunc(errCode, inArgs...)
//	}
//
//	if data == nil || data == "" {
//		c.Set("static", http.StatusForbidden)
//		c.Set("resp", gin.H{
//			"status":  400,
//			"retcode": errCode,
//			"msg":     msg,
//			"errmsg":  msg,
//		})
//		c.Abort()
//	}
//}

// modified by xiaoyanchun 2020-04-10
// 修改errCodeFunc参数的类型,支持多语言环境
// 修改前: errCodeFunc func(string, ...interface{}) string
// 修改后: errCodeFunc func(ss_lang.SsLang,string, ...interface{}) string
func ChkEmptyAndAbort(c *gin.Context, data interface{}, errCode string, errCodeFunc func(ss_lang.SsLang, string, ...interface{}) string, inArgs ...interface{}) {
	var msg string

	// 获取客户端语种,并检查语种是否存在，不存在会使用默认语种
	lang := ss_lang.NormalLaguage(GetCommonData(c).Lang)

	if len(inArgs) < 1 {
		msg = errCodeFunc(lang, errCode)
	} else {
		msg = errCodeFunc(lang, errCode, inArgs...)
	}

	if data == nil || data == "" {
		c.Set("static", http.StatusForbidden)
		c.Set("resp", gin.H{
			"status":  400,
			"retcode": errCode,
			"msg":     msg,
			"errmsg":  msg,
		})
		c.Abort()
	}
}

func HandleRetWebCli(retcode, successCode string, errCodeFunc func(string, ...interface{}) string, okFunc func(*gin.Context, ...interface{}) interface{},
	errFunc func(string, *gin.Context, ...interface{}) interface{},
	c *gin.Context, errMsg []interface{}, reply ...interface{}) {
	c.Set("status", http.StatusOK)
	if retcode == successCode {
		c.Set("resp", gin.H{
			"code": retcode,
			"msg":  "成功",
			"data": okFunc(c, reply...),
		})
	} else {
		var msg string
		if errMsg == nil {
			msg = errCodeFunc(retcode)
		} else {
			msg = errCodeFunc(retcode, errMsg...)
		}

		c.Set("resp", gin.H{
			"code": retcode,
			"msg":  msg,
			"data": errFunc(retcode, c, reply...),
		})
	}
	c.Next()
}

func RetList3Func(c2 *gin.Context, r ...interface{}) interface{} {
	return r[0]
}

func RetSingle3Func(c2 *gin.Context, r ...interface{}) interface{} {
	return r[0]
}

// 设置jwt的头信息
func SetJwtAuthentication(c *gin.Context, jwtStr string) {
	c.Header("Authentication", "Bearer "+jwtStr)
}

// 获取公共数据
func GetCommonData(c *gin.Context) ss_struct.HeaderCommonData {
	if data, exist := c.Get(constants.Common_Data); exist {
		if a, ok := data.(ss_struct.HeaderCommonData); ok {
			return a
		}
	}

	return ss_struct.HeaderCommonData{}
}
