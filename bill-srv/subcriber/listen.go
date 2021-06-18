package subcriber

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/bill-srv/handler"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"

	"context"

	go_micro_srv_push "a.a/mp-server/common/proto/push"
	"a.a/mp-server/common/ss_err"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
)

func InitSubcriber(s server.Server) error {
	err := micro.RegisterSubscriber(constants.Nats_Listen_Exp_key, s, processSettle, server.SubscriberQueue("queue.pubsub"))
	return err
}

func processSettle(ctx context.Context, req *go_micro_srv_push.SendListenExpKeyRequest) error {
	md, _ := metadata.FromContext(ctx)
	ss_log.Info("pos端超时确认 [ProcessPush] Received event %+v with metadata %+v\n", req, md)

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
	return nil
}
