package middleware

import (
	"a.a/cu/encrypt"
	"a.a/cu/ss_log"
	"a.a/mp-server/api-cb/common"
	"a.a/mp-server/api-cb/dao"
	"a.a/mp-server/api-cb/inner_util"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"strings"
)

type Md5SignerMw struct {
}

var Md5SignerMwInst Md5SignerMw

func (*Md5SignerMw) Sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		if inner_util.S(c, common.INNER_SIGN_METHOD) != common.SignMethod_Md5 {
			return
		}

		accNo := inner_util.M(c, common.INNER_DATA_ACCNO)
		if accNo == "" {
			c.Set(common.RET_CODE, ss_err.ERR_SYS_NO_ACCNO)
			c.Set(common.INNER_SIGN_VERIFY, false)
			return
		}
		inSign := inner_util.M(c, common.INNER_DATA_SIGN)
		if accNo == "" {
			c.Set(common.RET_CODE, ss_err.ERR_SYS_SIGN)
			c.Set(common.INNER_SIGN_VERIFY, false)
			return
		}

		paramMap := c.GetStringMap(common.INNER_PARAM_MAP)
		md5Key := dao.ApiDaoInstance.GetMd5Key(accNo)
		// 字典序 md5(k=v&k=v&...&key=md5key)
		reqStrEnBefore := encrypt.Map2FormStr(paramMap, md5Key, "&key=", encrypt.FIELD_ENCODED_NONE,
			[]string{}, "sign", false)
		// 全部小写
		sign := strings.ToLower(encrypt.DoMd5(reqStrEnBefore))
		ss_log.Info("before=[%v]", reqStrEnBefore)
		if sign != inSign {
			ss_log.Error("sign=[%v]|inSign=[%v]", sign, inSign)
			c.Set(common.INNER_SIGN_VERIFY, false)
			c.Set(common.RET_CODE, ss_err.ERR_SYS_SIGN)
			return
		}
		c.Set(common.INNER_SIGN_VERIFY, true)
		return
	}
}

func (*Md5SignerMw) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ss_log.Info("Md5SignerMw|Verify")
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================

		if inner_util.S(c, common.INNER_SIGN_METHOD) != common.SignMethod_Md5 {
			return
		}

		accNo := inner_util.M(c, common.INNER_DATA_ACCNO)
		if accNo == "" {
			c.Set(common.RET_CODE, ss_err.ERR_SYS_NO_ACCNO)
			c.Set(common.INNER_SIGN_VERIFY, false)
			return
		}
		inSign := inner_util.M(c, common.INNER_DATA_SIGN)
		if accNo == "" {
			c.Set(common.RET_CODE, ss_err.ERR_SYS_SIGN)
			c.Set(common.INNER_SIGN_VERIFY, false)
			return
		}

		paramMap := c.GetStringMap(common.INNER_PARAM_MAP)
		md5Key := dao.ApiDaoInstance.GetMd5Key(accNo)
		// 字典序 md5(k=v&k=v&...&key=md5key)
		reqStrEnBefore := encrypt.Map2FormStr(paramMap, md5Key, "&key=", encrypt.FIELD_ENCODED_NONE,
			[]string{}, "sign", false)
		// 全部小写
		sign := strings.ToLower(encrypt.DoMd5(reqStrEnBefore))
		ss_log.Biz("sign", "=============\nurl=[%v]\nparamMap=[%v]\nbefore=[%v]\nsign=[%v]\ninSign=[%v]",
			c.Request.URL.String(), paramMap, reqStrEnBefore, sign, inSign)
		if sign != inSign {
			//ss_log.Error("sign=[%v]|inSign=[%v]", sign, inSign)
			c.Set(common.INNER_SIGN_VERIFY, false)
			c.Set(common.RET_CODE, ss_err.ERR_SYS_SIGN)
			return
		}
		c.Set(common.INNER_SIGN_VERIFY, true)
		return
	}
}
