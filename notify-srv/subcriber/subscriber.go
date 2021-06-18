package subcriber

import (
	"context"
	"time"

	"a.a/mp-server/common/cache"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/notify-srv/common"
	"a.a/mp-server/notify-srv/handler"

	notifyProto "a.a/mp-server/common/proto/notify"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
)

func InitSubscriber(s server.Server) error {
	return micro.RegisterSubscriber(constants.PaySystemResultNotify, s, processNotify, server.SubscriberQueue("queue.pubsub"))
}

func processNotify(ctx context.Context, req *notifyProto.PaySystemResultNotify) error {
	lockKey := common.GetLockKey(req.OrderNo)
	lockValue := strext.NewUUID()
	// 获取分布式锁
	if !cache.GetDistributedLock(lockKey, lockValue, 30*time.Second) {
		return nil
	}

	md, _ := metadata.FromContext(ctx)
	ss_log.Info("支付系统消息异步通知 [ProcessPush] Received event %+v with metadata %+v\n", req, md)

	if req.OrderNo == "" {
		ss_log.Error("订单号为空")
		return nil
	}
	switch req.OrderType {
	case constants.VaReason_Cust_Pay_Order:
		handler.PayNotifyH.PaySuccessNotify(req.OrderNo)
	case constants.VaReason_BusinessTransferToBusiness:
		handler.TransferNotifyH.TransferSuccessNotify(req.OrderNo)
	case constants.VaReason_BusinessRefund:
		handler.RefundNotifyH.RefundNotify(req.OrderNo)
	default:
		ss_log.Error("订单类型错误")
	}

	// 释放分布式锁
	cache.ReleaseDistributedLock(lockKey, lockValue)

	return nil
}
