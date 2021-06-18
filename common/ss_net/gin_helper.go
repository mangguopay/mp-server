package ss_net

import (
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"log"
)

type SsNet struct {
}

var (
	NetU SsNet
)

/**
删除的常用封装
*/
func (*SsNet) DoDelete(c *gin.Context, handlerFunc func(interface{}) (string, error)) {
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
func (*SsNet) DoGetSingle(c *gin.Context, handlerFunc func(interface{}) (string, interface{}, error), key string) {
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
func (*SsNet) DoGetSingle2(c *gin.Context, handlerFunc func([]string) (string, interface{}, error), keys ...string) {
	var vals []string
	for _, v := range keys {
		vals = append(vals, c.Query(v))
	}
	resultCode, replyData, err := handlerFunc(vals)
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
func (*SsNet) DoGetV2(c *gin.Context, handlerFunc func(interface{}) (string, gin.H, error, []interface{}), key string) {
	var keyVal interface{}
	keyVal = c.Query(key)
	if "" == keyVal {
		params, _ := c.Get("params")
		keyVal = params
	}

	resultCode, retData, err, errMsg := handlerFunc(keyVal)
	if err != nil {
		log.Printf("DoGet|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		func(context *gin.Context, i ...interface{}) interface{} {
			return i[0]
		},
		RetNormalErrFunc,
		c, errMsg, retData)
}

/**
查询单个的常用封装
*/
func (*SsNet) DoGet(c *gin.Context, handlerFunc func(interface{}) (string, gin.H, error), key string) {
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
func (*SsNet) DoUpdate(c *gin.Context, handlerFunc func(interface{}) (string, interface{}, error)) {
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

func (*SsNet) DoUpdate3(c *gin.Context, handlerFunc func(interface{}) (string, interface{}, error)) {
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

func (*SsNet) DoUpdate4(c *gin.Context, handlerFunc func(interface{}) (string, string, interface{}, error)) {
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

func (*SsNet) DoUpdateWithMsg(c *gin.Context, handlerFunc func(interface{}) (string, string, interface{}, error)) {
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

/**
查询多个的常用封装
*/
func (*SsNet) DoGetList(c *gin.Context, handlerFunc func([]string) (string, interface{}, int32, error), keys ...string) {
	var vals []string
	for _, v := range keys {
		vals = append(vals, c.Query(v))
	}
	resultCode, replyDataList, replyTotal, err := handlerFunc(vals)
	if err != nil {
		log.Printf("DoGetSingle|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetListFunc,
		RetNormalErrFunc,
		c, nil, replyDataList, replyTotal)
}

func (*SsNet) DoGetList2(c *gin.Context, handlerFunc func([]string) (string, interface{}, int32, map[string]interface{}, error), keys ...string) {
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

/**
更新的常用封装
*/
func (*SsNet) DoUpdate2(c *gin.Context, handlerFunc func(interface{}) (string, interface{}, error)) {
	params, _ := c.Get("params")
	resultCode, replyUid, err := handlerFunc(params)
	if err != nil {
		log.Printf("DoUpdate|err=%v\n", err)
	}

	HandleRetMulti(resultCode, ss_err.ERR_SUCCESS, ss_err.GetErrMsgMulti,
		RetSingleFunc,
		RetNormalErrFunc,
		c, nil, replyUid)
}
