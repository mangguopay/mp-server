package subcriber

import (
	"context"

	go_micro_srv_gis "a.a/mp-server/common/proto/gis"

	"a.a/cu/ss_log"
	"a.a/mp-server/gis-srv/common"
	"a.a/mp-server/gis-srv/dao"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/log"
)

func InitSubcriber(s server.Server) error {
	err := micro.RegisterSubscriber("srv_gis", s, processEvent, server.SubscriberQueue("queue.pubsub"))
	return err
}

func processEvent(ctx context.Context, event *go_micro_srv_gis.ListenEvenRequest) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("[pubsub.1] Received event %+v with metadata %+v\n", event, md)

	//  同步数据
	if event.IsSync {
		coordinates, err := dao.ServiceDaoInst.GetSrvCoordinate()
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return err
		}
		common.SrvCoordinates = coordinates
		ss_log.Info("监听更新服务商成功|len=[%v]", len(coordinates))
	}
	return nil
}
