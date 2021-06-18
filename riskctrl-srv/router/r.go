package router

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/riskctrl"
	"a.a/mp-server/riskctrl-srv/handler"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/v2/broker"
)

func R(na *broker.Broker) {
	_, err := (*na).Subscribe(constants.Nats_Topic_Risk, func(p broker.Event) error {
		r(p.Message().Header["m"], p.Message().Body)

		return nil
	})
	if err != nil {
		ss_log.Info("err=[%v]", err)
	}
}

func r(me string, body []byte) {
	switch me {
	case "risk":
		req := &go_micro_srv_riskctrl.RiskOfflineRequest{}
		err := proto.Unmarshal(body, req)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		ss_log.Info("recv|method=[%v],body=[%v]", me, req)

		// 处理消息
		handler.RecvMQMsg(req)

	default:
		ss_log.Info("recv|not support|method=[%v],body=[%v]", me, body)
	}
}
