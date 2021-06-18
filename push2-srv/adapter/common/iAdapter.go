package common

import (
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/push2-srv/m"
)

type IPushAdapter interface {
	// 获取配置
	GetAccToken(accReq *go_micro_srv_push.PushAccout) (string, string)
	// 发送
	Send(m.SendReq) (string, string)
}
