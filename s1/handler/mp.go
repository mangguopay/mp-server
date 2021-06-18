package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/net/consts"
	"a.a/net/module"
	"a.a/net/proc"
)

// 创蓝发送短信
func DoSend(url string, req map[string]interface{}, signKey string) *module.CommonSendResp {
	ss_log.InitLog(".")
	reqInner := &module.CommonSendReq{}
	reqInner.IsUseMap = false
	reqInner.UrlFull = url
	reqInner.InData = req
	reqInner.SignKey = signKey
	//
	reqInner.PreSendSeq = map[string][]module.ExecCmd{
		consts.ExecTagPay: []module.ExecCmd{
			{
				ExecType: consts.EXEC_CallFunc,
				ExecParam: []interface{}{
					doPackBodyStr,
				},
			},
		},
	}
	reqInner.PostSendSeq = map[string][]module.ExecCmd{
		consts.ExecTagPay: []module.ExecCmd{
			{
				ExecType: consts.EXEC_CallFunc,
				ExecParam: []interface{}{
					doRouteResult,
					doRetOk,
					doRetErrSend,
					doRetErrUpTech,
				},
			},
		},
	}
	reqInner.ContentType = consts.HTTP_CONTENT_JSON
	reqInner.RetContentType = consts.HTTP_CONTENT_JSON
	reqInner.SenderType = consts.HTTP_METHOD_POST_BODY
	reqInner.LogicParter = consts.ExecTagPay
	reqInner.IsChkHttpsName = false
	reqInner.IsInitCertPair = false
	respInner := proc.CommonSend(reqInner)
	ss_log.Info("respInner=[%v]", respInner)
	return respInner
}

func doPackBodyStr(execContext *module.ExecContext) {
	if execContext.SignKey != "" {
		execContext.HeaderMap = &map[string]interface{}{
			"Authorization": "Bearer " + execContext.SignKey,
		}
	}
	execContext.PostData = strext.ToJson(execContext.InData)
}

func doRouteResult(execContext *module.ExecContext) int {
	// 上游返回错误
	if execContext.RetMap == nil || (*execContext.RetMap)["retcode"] == nil {
		return 2
	}

	// 上游通信异常
	if strext.ToStringNoPoint((*execContext.RetMap)["retcode"]) != "0" {
		return 3
	}

	return 1
}

func doRetOk(execContext *module.ExecContext) {
	execContext.RetCode = consts.RETCODE_SUCCESS
	execContext.RetMsg = consts.GetErrMsg(consts.RETCODE_SUCCESS)
	execContext.OutRetCode = strext.ToStringNoPoint((*execContext.RetMap)["code"])
	execContext.OutRetMsg = strext.ToStringNoPoint((*execContext.RetMap)["error"])
}

func doRetErrSend(execContext *module.ExecContext) {
	execContext.RetCode = consts.RETCODE_NETWORK
	execContext.RetMsg = consts.GetErrMsg(consts.RETCODE_NETWORK)
}

func doRetErrUpTech(execContext *module.ExecContext) {
	execContext.RetCode = consts.RETCODE_BIZ_ERR
	execContext.RetMsg = consts.GetErrMsg(consts.RETCODE_BIZ_ERR)
	execContext.OutRetCode = strext.ToStringNoPoint((*execContext.RetMap)["code"])
	execContext.OutRetMsg = strext.ToStringNoPoint((*execContext.RetMap)["error"])

}
