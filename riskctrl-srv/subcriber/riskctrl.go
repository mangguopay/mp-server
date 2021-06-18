package subcriber

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	"a.a/mp-server/riskctrl-srv/handler"
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
)

func InitSubcriber(s server.Server) error {
	return micro.RegisterSubscriber(constants.Nats_Broker_Header_Risk, s, processRiskOffline, server.SubscriberQueue("queue.pubsub"))
}

func processRiskOffline(ctx context.Context, req *go_micro_srv_riskctrl.RiskOfflineRequest) {
	md, _ := metadata.FromContext(ctx)
	ss_log.Info("[ProcessPush] Received event %+v with metadata %+v\n", req, md)

	// 处理消息
	handler.RecvMQMsg(req)
}
