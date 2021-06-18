package router

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/business-settle-srv/service"
	"a.a/mp-server/common/constants"
	businessSettleProto "a.a/mp-server/common/proto/business-settle"
	"a.a/mp-server/common/ss_err"
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/v2/broker"
	"time"
)

//func sub(topic string) {
//	_, err := broker.Subscribe(topic, func(p broker.Event) error {
//		ss_log.Info("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
//
//		return nil
//	})
//	if err != nil {
//		ss_log.Error("err=[%v]", err)
//	}
//}

func R(na *broker.Broker) {
	_, err := (*na).Subscribe("go.micro.topic.settle", func(p broker.Event) error {
		r(p.Message().Header["m"], p.Message().Body)

		return nil
	})
	if err != nil {
		ss_log.Info("err=[%v]", err)
	}
}

func r(me string, body []byte) {
	switch me {
	case "settle":
		req := &businessSettleProto.BusinessSettleFeesRequest{}
		err := proto.Unmarshal(body, req)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		ss_log.Info("recv|method=[%v],body=[%v]", me, req)
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
		defer cancel()
		switch req.FeesType {
		case constants.FEES_TYPE_BILL: // 入金
			if errStr := service.BillService(ctx, req); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[处理入金手续费分成失败,err----->%s]", errStr)
			}
		case constants.FEES_TYPE_WITHDRAWAL: // 提现
			if errStr := service.WithdrawalService(ctx, req); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[处理出金手续费分成失败,err----->%s]", errStr)
			}
		default:
			ss_log.Error("err=[手续费类型不对,当前接受的类型为----->%s]", req.FeesType)
		}

	default:
		ss_log.Info("recv|not support|method=[%v],body=[%v]", me, body)
	}
}
