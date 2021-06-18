package i

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/gis-srv/common"
	"a.a/mp-server/gis-srv/dao"
)

func InitGisInfo() {
	coordinates, err := dao.ServiceDaoInst.GetSrvCoordinate()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return
	}
	common.SrvCoordinates = coordinates
	ss_log.Info("加载文件成功|len=[%v]", len(coordinates))
}
