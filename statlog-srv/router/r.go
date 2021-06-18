package router

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_statlog "a.a/mp-server/common/proto/statlog"
	handler2 "a.a/mp-server/statlog-srv/handler"
	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/v2/broker"
)

func R(na *broker.Broker) {
	_, err := (*na).Subscribe(constants.Nats_Topic_Statlog, func(p broker.Event) error {
		r(p.Message().Header["m"], p.Message().Body)

		return nil
	})
	if err != nil {
		ss_log.Info("err=[%v]", err)
	}
}

func r(me string, body []byte) {
	switch me {
	case "PushApiLog":
		req := &go_micro_srv_statlog.PushApiLogRequest{}
		err := proto.Unmarshal(body, req)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		ss_log.Info("recv|method=[%v],body=[%v]", me, req)

		// 处理消息
		handler2.StatlogHandlerInst.PushApiLog(req)
	default:
		ss_log.Info("recv|not support|method=[%v],body=[%v]", me, body)
	}
}
