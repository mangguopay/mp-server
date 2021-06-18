package common

import (
	"a.a/cu/ss_log"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/push2-srv/m"
	"fmt"
)

var (
	PushAdapterWrapperInst PushAdapterWrapper
)

type PushAdapterWrapper struct {
}

// 评估风控处理方式
func (*PushAdapterWrapper) Send(adpterType string, req m.SendReq) (string, string) {
	targetApi := getTargetApi(adpterType)
	if targetApi == nil {
		ss_log.Error("[%v]is nil", adpterType)
		return fmt.Sprintf("%v is nil", adpterType), ss_err.ERR_PARAM
	}
	return targetApi.Send(req)
}

// 评估风控处理方式
func (*PushAdapterWrapper) GetAccToken(adpterType string, accReq *go_micro_srv_push.PushAccout) (string, string) {
	targetApi := getTargetApi(adpterType)
	if targetApi == nil {
		ss_log.Error("[%v]is nil", adpterType)
		return "", "没有对应适配器"
	}
	return targetApi.GetAccToken(accReq)
}
