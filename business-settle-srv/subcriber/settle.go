package subcriber

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/business-settle-srv/handler"
	"a.a/mp-server/business-settle-srv/service"
	"a.a/mp-server/common/constants"
	businessSettleProto "a.a/mp-server/common/proto/business-settle"
	"a.a/mp-server/common/ss_err"
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"time"
)

func InitSubcriber(s server.Server) error {
	err := micro.RegisterSubscriber(constants.BusinessSettle, s, processSettle, server.SubscriberQueue("queue.pubsub"))
	if err != nil {
		fmt.Println("err: --->", err.Error())
	}
	return businessSettleProto.RegisterBusinessSettleHandler(s, new(handler.BusinessSettle))
}

func processSettle(ctx context.Context, req *businessSettleProto.BusinessSettleFeesRequest) error {
	fmt.Println("接受推送--------------------------->")
	md, _ := metadata.FromContext(ctx)
	ss_log.Info("[ProcessPush] Received event %+v with metadata %+v\n", req, md)

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
	return nil
}
