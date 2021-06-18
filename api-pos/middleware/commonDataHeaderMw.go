package middleware

import (
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/util"
	"a.a/mp-server/api-pos/common"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_struct"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	// md5签名时添加的盐
	CommonDataMd5Salt = "NbVcCHUJewP721WDXgyGlzKA50qtxdYL"

	// 请求的有效时间 (单位: 秒)
	RequestExpireTime int64 = 120
)

type ParseCommonDataHeaderMw struct {
}

var ParseCommonDataHeaderMwInst *ParseCommonDataHeaderMw

// 解析头部公共信息
func (p *ParseCommonDataHeaderMw) ParseCommonDataHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !p.parseCommonDataHeader(c) { // 验证失败
			c.Set(common.RET_CODE, ss_err.ERR_SYS_Common_Data)
			c.Set(common.INNER_IS_STOP, true)
			RespMwInst.PackInner(c)
			RsaMwInst.EncodeInner(c)
			RsaMwInst.SignInner(c)
			RespMwInst.RespInner(c)
			c.Abort()
			return
		}
	}
}

func (p *ParseCommonDataHeaderMw) parseCommonDataHeader(c *gin.Context) bool {
	logPrefix := c.GetString(common.INNER_TRACE_NO) + "|通用头信息:"
	dataSep := "&sign=" // 数据分隔标识

	cData := ss_struct.HeaderCommonData{}
	cData.Lang = constants.LangEnUS // 默认参数

	defer func() { // 在defer中设置参数，防止验证失败时没有默认参数
		// 设置公共数据
		c.Set(constants.Common_Data, cData)
	}()

	// 获取头信息  示例: `{"lang":"zh_CN","app_version":"1.0.0","timestamp":1588748802000,"platform":"android","app_name":"pos_client","utm_source":"official"}&sign=62e6d741646d8f7bec65a297d2321d5a`
	strData := strings.TrimSpace(c.GetHeader("CommonData"))
	if strData == "" {
		ss_log.Error(logPrefix + "CommonData字段为空")
		return false
	}

	// base64解码
	if byteData, dErr := base64.StdEncoding.DecodeString(strData); dErr != nil {
		ss_log.Error(logPrefix+"base64解码失败,err:%v, strData:%s", dErr, strData)
		return false
	} else {
		strData = string(byteData)
	}

	// 进行分割数据
	arr := strings.Split(strData, dataSep)
	if len(arr) != 2 {
		ss_log.Error(logPrefix+"数据格式错误,strData:%s", strData)
		return false
	}

	// 校验sign是否正确
	if encrypt.DoMd5Salted(arr[0], CommonDataMd5Salt) != arr[1] {
		ss_log.Error(logPrefix+"校验sign失败,数据:%s, sign:%s", arr[0], arr[1])
		return false
	}

	// 解析到结构体
	if jerr := json.Unmarshal([]byte(arr[0]), &cData); jerr != nil {
		ss_log.Error(logPrefix+"json解码失败,err:%v, data:%s", jerr, arr[0])
		return false
	}

	// 验证请求是否已经过期了
	if cData.Timestamp+RequestExpireTime < ss_time.NowTimestamp(global.Tz) {
		ss_log.Error(logPrefix+"请求地址已过期,请求时间:%d,当前时间:%d", cData.Timestamp, ss_time.NowTimestamp(global.Tz))
		return false
	}

	// 去掉空格
	cData.Lang = strings.TrimSpace(cData.Lang)
	cData.AppVersion = strings.TrimSpace(cData.AppVersion)

	// 如果未传则设置默认参数
	if cData.Lang == "" || !util.InSlice(cData.Lang, []string{constants.LangEnUS, constants.LangKmKH, constants.LangZhCN}) {
		cData.Lang = constants.LangEnUS
	}

	return true
}
