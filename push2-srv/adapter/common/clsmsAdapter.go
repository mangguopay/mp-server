package common

import (
	"errors"
	"fmt"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/push2-srv/dao"
	"a.a/mp-server/push2-srv/m"
	"a.a/net/consts"
	"a.a/net/module"
	"a.a/net/proc"
)

var (
	ClSmsAdapterInst ClSmsAdapter
)

type ClSmsAdapter struct {
}

func redisErrorIsNil(redisErr error) bool {
	return redisErr.Error() == "redis: nil"
}

func (f ClSmsAdapter) GetAccToken(accReq *go_micro_srv_push.PushAccout) (string, string) {
	ss_log.Info("GetAccToken -------------->%v", accReq)
	var phone, country string
	if accReq.Phone != "" && accReq.CountryCode != "" {
		country, phone = accReq.CountryCode, accReq.Phone
	} else {
		phone, country = dao.AccDaoInst.GetPhone(accReq.AccountNo)
	}
	return fmt.Sprintf("%s%s", country, phone), ""
}

func (f ClSmsAdapter) Send(req m.SendReq) (string, string) {
	account := strext.ToStringNoPoint(req.Config["account"])
	secret := strext.ToStringNoPoint(req.Config["secret"])
	gateway := strext.ToStringNoPoint(req.Config["gateway"])
	//mobile := req.AccToken

	return doSms(account, secret, gateway, req.AccToken, req.Content)
}

func doSms(account, secret, gateway, accountToken, msg string) (string, string) {
	// 1.获取发送短信的配置信息
	switch "" {
	case account:
		ss_log.Error("err=[%v]", "获取短信账号失败")
	case secret:
		ss_log.Error("err=[%v]", "获取短信密码失败")
	case gateway:
		ss_log.Error("err=[%v]", "获取短信链接失败")
	}

	// 执行发送
	message, err := clSendSms(gateway, account, secret, accountToken, msg)
	if err != nil {
		ss_log.Error("err=[%v],missing key=[%v]", err.Error(), "发送消息的日志存放进数据库失败")
		return message, ss_err.ERR_PARAM
	}
	return message, ss_err.ERR_SUCCESS
}

// 创蓝发送短信
func clSendSms(smsURL, account, secret, mobile, msg string) (string, error) {
	fmt.Println("--------------", smsURL, account, secret, mobile, msg)

	reqInner := &module.CommonSendReq{}
	reqInner.IsUseMap = false
	reqInner.UrlFull = smsURL
	reqInner.InData = map[string]interface{}{
		"account":  account,
		"password": secret,
		"msg":      msg,
		"mobile":   mobile,
	}
	reqInner.Fields = []string{"account", "password", "msg", "mobile"}
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
	respCode := strext.ToStringNoPoint(respInner.ExtMap["code"])
	respErr := strext.ToStringNoPoint(respInner.ExtMap["error"])
	if respCode != "0" {
		return respErr, errors.New("发送失败,code: " + respCode)
	}
	return respErr, nil
}

func doPackBodyStr(execContext *module.ExecContext) {
	reqMap := make(map[string]interface{})
	for _, v := range execContext.Fields {
		reqMap[v] = strext.ToStringNoPoint((*execContext.InData)[v])
	}
	//(*execContext.InData)[consts.ExecInDataField_bodyPack1] = reqMap
	execContext.PostData = strext.ToJson(reqMap)
}

func doRouteResult(execContext *module.ExecContext) int {
	// 上游返回错误
	if execContext.RetMap == nil || (*execContext.RetMap)["code"] == nil {
		return 2
	}

	// 上游通信异常
	if strext.ToStringNoPoint((*execContext.RetMap)["code"]) != "0" {
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
