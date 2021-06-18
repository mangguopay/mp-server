package router

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/settle-srv/service"
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
	_, err := (*na).Subscribe(constants.Nats_Topic_Settle, func(p broker.Event) error {
		r(p.Message().Header["m"], p.Message().Body)

		return nil
	})
	if err != nil {
		ss_log.Info("err=[%v]", err)
	}
}

func r(me string, body []byte) {
	switch me {
	case constants.Settle_Type:
		req := &go_micro_srv_settle.SettleTransferRequest{}
		err := proto.Unmarshal(body, req)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
		ss_log.Info("recv|method=[%v],body=[%v]", me, req)

		switch req.FeesType {
		case constants.FEES_TYPE_EXCHANGE: // 兑换
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
			defer cancel()
			if errStr := service.ExchangeService(ctx, req); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[处理兑换利息失败,err----->%s]", errStr)
			}
		case constants.FEES_TYPE_TRANSFER: // 转账
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
			defer cancel()
			if errStr := service.TransferService(ctx, req); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[处理转账利息失败,err----->%s]", errStr)
			}
		case constants.FEES_TYPE_COLLECTION: // 收款
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
			defer cancel()
			if errStr := service.CollectionService(ctx, req); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[处理收款利息失败,err----->%s]", errStr)
			}
		case constants.FEES_TYPE_SAVEMONEY: // 存款
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
			defer cancel()
			if errStr := service.SavemoneyService(ctx, req); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[处理存款利息失败,err----->%s]", errStr)
			}
		case constants.FEES_TYPE_WITHDRAW: // 取款
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
			defer cancel()
			if errStr := service.MobileNumWithdrawService(ctx, req); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("err=[处理手机号取款利息失败,err----->%s]", errStr)
			}
		//case constants.FEES_TYPE_SWEEP_WITHDRAW: // 扫码取款
		//	if errStr := service.SweepWithdrawService(req); errStr != ss_err.ERR_SUCCESS {
		//		ss_log.Error("err=[处理扫码取款利息失败,err----->%s]", errStr)
		//	}
		default:
			ss_log.Error("err=[手续费类型不对,当前接受的类型为----->%s]", req.FeesType)
		}

	default:
		ss_log.Info("recv|not support|method=[%v],body=[%v]", me, body)
	}
}
