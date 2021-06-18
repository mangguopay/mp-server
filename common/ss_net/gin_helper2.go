package ss_net

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"log"
)

type SsJet struct {
}

var (
	SsJetInst SsJet
)

//================================
func (*SsJet) transQueryToMap(c *gin.Context, mustNotNil []string) map[string]interface{} {
	urlValues := c.Request.URL.Query()
	vals := map[string]interface{}{}
	for k, v := range urlValues {
		// 查找是否需要过滤
		lenMust := len(mustNotNil)
		if lenMust > 0 {
			idx := container.GetKey(k, mustNotNil)
			if idx >= 0 && len(v) > 0 {
				if lenMust > idx+1 {
					mustNotNil = append(mustNotNil[:idx], mustNotNil[idx+1:]...)
				} else {
					mustNotNil = mustNotNil[:idx]
				}
			}
		}
		//
		if vals[k] == nil {
			if len(v) == 1 {
				vals[k] = v[0]
			} else {
				vals[k] = v
			}
		} else {
			tmp := vals[k]
			var tmp3 []string
			switch tmp2 := tmp.(type) {
			case []string:
				tmp3 = append(tmp2, v...)
			case string:
				tmp3 = append(tmp3, tmp2)
				tmp3 = append(tmp3, v...)
			}
			vals[k] = tmp3
		}
	}
	if len(mustNotNil) > 0 {
		ss_log.Error("少了入参，或者入参为空=[%v]", mustNotNil)
		HandleRetMulti(ss_err.ERR_PARAM, ss_err.ERR_PARAM, ss_err.GetErrMsgMulti,
			RetSingleFunc,
			RetNormalErrFunc,
			c, nil, nil)
		return nil
	}
	return vals
}

//================================
func (*SsJet) transPostToMap(c *gin.Context, mustNotNil []string) map[string]interface{} {
	params, _ := c.Get("params")
	var m map[string]interface{}
	switch params2 := params.(type) {
	case map[string]interface{}:
		m = params2
	case string:
		m = strext.Json2Map(params2)
	}

	for k, v := range m {
		// 查找是否需要过滤
		if len(mustNotNil) > 0 {
			idx := container.GetKey(k, mustNotNil)
			if idx >= 0 && v != nil {
				mustNotNil = append(mustNotNil[:idx], mustNotNil[idx:]...)
			}
		}
	}
	if len(mustNotNil) > 0 {
		ss_log.Error("少了入参，或者入参为空=[%v]", mustNotNil)
		HandleRetMulti(ss_err.ERR_PARAM, ss_err.ERR_PARAM, ss_err.GetErrMsgMulti,
			RetSingleFunc,
			RetNormalErrFunc,
			c, nil, nil)
		return nil
	}
	return m
}

//================================
/**
查询单个的常用封装
*/
func (r *SsJet) Get(c *gin.Context, handlerFunc func(params map[string]interface{}) (string, interface{}), mustNotNil []string) {
	vals := r.transQueryToMap(c, mustNotNil)
	if vals == nil {
		return
	}
	resultCode, replyData := handlerFunc(vals)
	if resultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("retcode=[%v]", resultCode)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetSingleFunc,
		RetNormalErrFunc,
		c, nil, replyData)
}

//================================
/**
查询多个的常用封装
*/
func (r *SsJet) GetList(c *gin.Context, handlerFunc func(map[string]interface{}) (string, interface{}, int32), mustNotNil []string) {
	vals := r.transQueryToMap(c, mustNotNil)
	if vals == nil {
		return
	}
	resultCode, replyDataList, replyTotal := handlerFunc(vals)
	if resultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("retcode=[%v]", resultCode)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetListFunc,
		RetNormalErrFunc,
		c, nil, replyDataList, replyTotal)
}

//================================
/**
更新的常用封装
*/
func (r *SsJet) Update(c *gin.Context, handlerFunc func(map[string]interface{}) string, mustNotNil []string) {
	resultCode := handlerFunc(r.transPostToMap(c, mustNotNil))
	if resultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("retcode=[%v]", resultCode)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetSingleFunc,
		RetNormalErrFunc,
		c, nil, 0)
}

//================================
/**
删除的常用封装
*/
func (r *SsJet) Delete(c *gin.Context, handlerFunc func(map[string]interface{}) string, mustNotNil []string) {
	resultCode := handlerFunc(r.transPostToMap(c, mustNotNil))
	if resultCode != ss_err.ERR_SUCCESS {
		ss_log.Error("retcode=[%v]", resultCode)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetMsgFunc,
		RetNormalErrFunc,
		c, nil)
}

//================================

/**
删除的常用封装
*/
func (*SsJet) DoDelete(c *gin.Context, handlerFunc func(interface{}) (string, error)) {
	params, _ := c.Get("params")
	resultCode, err := handlerFunc(params)
	if err != nil {
		log.Printf("DoDelete|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetMsgFunc,
		RetNormalErrFunc,
		c, nil)
}

/**
查询单个的常用封装
*/
func (*SsJet) DoGetSingle(c *gin.Context, handlerFunc func(interface{}) (string, interface{}, error), key string) {
	keyVal := c.Query(key)
	ChkEmptyAndAbort(c, keyVal, ss_err.ERR_ARGS, ss_err.GetErrMsgMulti, key)
	resultCode, replyData, err := handlerFunc(keyVal)
	if err != nil {
		log.Printf("DoGetSingle|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetSingleFunc,
		RetNormalErrFunc,
		c, nil, replyData)
}

/**
查询单个的常用封装
*/
func (*SsJet) DoGet(c *gin.Context, handlerFunc func(interface{}) (string, gin.H, error), key string) {
	var keyVal interface{}
	keyVal = c.Query(key)
	if "" == keyVal {
		params, _ := c.Get("params")
		keyVal = params
	}

	resultCode, retData, err := handlerFunc(keyVal)
	if err != nil {
		log.Printf("DoGet|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		func(context *gin.Context, i ...interface{}) interface{} {
			return i[0]
		},
		RetNormalErrFunc,
		c, nil, retData)
}

/**
更新的常用封装
*/
func (*SsJet) DoUpdate(c *gin.Context, handlerFunc func(interface{}) (string, interface{}, error)) {
	params, _ := c.Get("params")
	resultCode, replyUid, err := handlerFunc(params)
	if err != nil {
		log.Printf("DoUpdate|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetSingleUidFunc,
		RetNormalErrFunc,
		c, nil, replyUid)
}

func (*SsJet) DoUpdate3(c *gin.Context, handlerFunc func(interface{}) (string, interface{}, error)) {
	params, _ := c.Get("params")
	resultCode, reply, err := handlerFunc(params)
	if err != nil {
		log.Printf("DoUpdate|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetCommonFunc,
		RetNormalErrFunc,
		c, nil, reply)
}

func (*SsJet) DoUpdate4(c *gin.Context, handlerFunc func(interface{}) (string, string, interface{}, error)) {
	params, _ := c.Get("params")
	resultCode, retMsg, reply, err := handlerFunc(params)
	if err != nil {
		log.Printf("DoUpdate|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetSingleUidFunc,
		RetNormalErrFunc,
		c, []interface{}{retMsg}, reply)
}

func (*SsJet) DoUpdateWithMsg(c *gin.Context, handlerFunc func(interface{}) (string, string, interface{}, error)) {
	params, _ := c.Get("params")
	resultCode, retMsg, reply, err := handlerFunc(params)
	if err != nil {
		log.Printf("DoUpdate|err=%v\n", err)
	}

	tag := "0"
	if resultCode == ss_err.ERR_PARAM {
		tag = "1"
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetSingleUidFunc,
		RetNormalErrFunc,
		c, []interface{}{tag, retMsg}, reply)
}

func (*SsJet) DoGetList2(c *gin.Context, handlerFunc func([]string) (string, interface{}, int32, map[string]interface{}, error), keys ...string) {
	var vals []string
	for _, v := range keys {
		vals = append(vals, c.Query(v))
	}
	resultCode, replyDataList, replyTotal, extData, err := handlerFunc(vals)
	if err != nil {
		log.Printf("DoGetSingle|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetList2Func,
		RetNormalErrFunc,
		c, nil, replyDataList, replyTotal, extData)
}
