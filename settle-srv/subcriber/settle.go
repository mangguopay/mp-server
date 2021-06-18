package subcriber

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_settle "a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/settle-srv/service"
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"time"
)

func InitSubcriber(s server.Server) error {
	return micro.RegisterSubscriber(constants.Settle_Type, s, processSettle, server.SubscriberQueue("queue.pubsub"))
}

func processSettle(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest) error {
	md, _ := metadata.FromContext(ctx)
	ss_log.Info("[ProcessPush] Received event %+v with metadata %+v\n", req, md)

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
		if errStr := service.WithdrawService(ctx, req); errStr != ss_err.ERR_SUCCESS {
			ss_log.Error("err=[处理取款利息失败,err----->%s]", errStr)
		}
	default:
		ss_log.Error("err=[手续费类型不对,当前接受的类型为----->%s]", req.FeesType)
	}

	return nil
}
