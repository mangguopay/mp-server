package router

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/bill-srv/handler"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/v2/broker"
)

func R(na *broker.Broker) {
	_, err := (*na).Subscribe(constants.Nats_Topic_Listen_Exp_key, func(p broker.Event) error {
		r(p.Message().Header["m"], p.Message().Body)

		return nil
	})
	if err != nil {
		ss_log.Info("err=[%v]", err)
	}
}

func r(me string, body []byte) {
	switch me {
	case "SendListenExpKeyMsg":
		req := &go_micro_srv_push.SendListenExpKeyRequest{}
		err := proto.Unmarshal(body, req)
		if err != nil {
			ss_log.Info("err=[%v]", err)
		}
		ss_log.Info("recv|method=[%v],body=[%v]", me, req)
		// 调用取消确认接口
		// 调用取消的rpc流程.
		cancelReq := &go_micro_srv_bill.CancelWithdrawRequest{
			OrderNo:      req.OrderNo,
			CancelReason: "pos端超时确认",
		}
		cancelRepl := &go_micro_srv_bill.CancelWithdrawReply{}
		_ = handler.BillHandlerInst.CancelWithdraw(context.TODO(), cancelReq, cancelRepl)
		if cancelRepl.ResultCode != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[pos端超时确认,调用取消确认的rpc失败,订单号为----->%s]", req.OrderNo)
		}

	default:
		ss_log.Info("recv|not support|method=[%v],body=[%v]", me, body)
	}
}
