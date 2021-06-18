package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/gis"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/gis-srv/common"
	"context"
	"github.com/json-iterator/go"
)

var GisHandlerInst GisHandler

type GisHandler struct{}

func (GisHandler) GetNearbyServicerList(ctx context.Context, req *go_micro_srv_gis.GetNearbyServicerListRequest, reply *go_micro_srv_gis.GetNearbyServicerListReply) error {
	if req.Page < 1 || req.PageSize < 1 {
		ss_log.Error("参数错误:Page:%d,PageSize:%d", req.Page, req.PageSize)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	redisKey := "nearby_servicer_" + req.Lat + "_" + req.Lng
	lat := strext.ToFloat64(req.Lat)
	lng := strext.ToFloat64(req.Lng)
	//返回拍好序的
	serviceList := []*go_micro_srv_gis.NearbyServicerData{}
	var total int32
	if value, _ := cache.RedisClient.Get(redisKey).Result(); value == "" { // 查询数据库并设置进redis
		serviceList = countSrvCoordinate(lat, lng)
		// ==========如果 lat 和 lng 为0 则不排序==========
		if lat == 0 && lng == 0 {
			serviceList = common.SrvCoordinates
			ss_log.Info("app定位不了定位,返回未排序的所有的服务商列表,serviceList: %+v", serviceList)
		}
		// ==========如果 lat 和 lng 为0 则不排序==========

		if len(serviceList) == 0 {
			reply.ResultCode = ss_err.ERR_SUCCESS
			reply.Datas = serviceList
			reply.Total = 0
			return nil
		}

		// 设置到redis中去，请求第二页时候就从redis中拿数据。
		datasString := strext.ToJson(serviceList)
		if err := cache.RedisClient.Set(redisKey, datasString, constants.CacheKeySecV2).Err(); err != nil {
			ss_log.Error("附近服务商存进redis失败,err=[%v],datas--->%s", err, serviceList)
		}
		total = strext.ToInt32(len(serviceList))
	} else {
		jsoniter.Unmarshal([]byte(value), &serviceList)
		total = strext.ToInt32(len(serviceList))
	}

	startInt := strext.ToInt((req.Page - 1) * req.PageSize)
	if startInt > len(serviceList) {
		ss_log.Error("分页条数大于当前切片长度,请求为第 %v 页,页面大小为 %v, 切片长度为: %v", req.Page, req.PageSize, len(serviceList))
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}
	endInt := strext.ToInt(req.Page * req.PageSize)
	if len(serviceList) < endInt {
		serviceList = serviceList[startInt:]
	} else {
		serviceList = serviceList[startInt:endInt]
	}

	reply.Total = total
	reply.Datas = serviceList
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 快速排序
func quickSort(values []*go_micro_srv_gis.NearbyServicerData, left int, right int) {
	if left < right {
		// 设置分水岭
		temp := values[left]
		// 设置哨兵
		i, j := left, right
		for {
			// 从右向左找，找到第一个比分水岭小的数
			for values[j].Distance >= temp.Distance && i < j {
				j--
			}
			// 从左向右找，找到第一个比分水岭大的数
			for values[i].Distance <= temp.Distance && i < j {
				i++
			}
			// 如果哨兵相遇，则退出循环
			if i >= j {
				break
			}
			// 交换左右两侧的值
			values[i], values[j] = values[j], values[i]
		}
		// 将分水岭移到哨兵相遇点
		values[left] = values[i]
		values[i] = temp
		// 递归，左右两侧分别排序
		quickSort(values, left, i-1)
		quickSort(values, i+1, right)
	}
}

// 计算和排序周围的服务商坐标
func countSrvCoordinate(lat, lng float64) []*go_micro_srv_gis.NearbyServicerData {
	serviceList := make([]*go_micro_srv_gis.NearbyServicerData, len(common.SrvCoordinates))
	copy(serviceList, common.SrvCoordinates)

	if len(serviceList) == 0 {
		ss_log.Error("CountSrvCoordinate, 服务商坐标列表为空")
		return serviceList
	}

	for _, data := range serviceList {
		//与查询的经纬度的大圆距离
		distance := ss_count.CountCircleDistance(strext.ToFloat64(data.Lat), strext.ToFloat64(data.Lng), strext.ToFloat64(lat), strext.ToFloat64(lng))
		data.Distance = distance
	}
	//bf, _ := jsoniter.MarshalToString(serviceList)
	//fmt.Println("排序前--------->", bf)
	//sort.Sort()
	quickSort(serviceList, 0, len(serviceList)-1)
	//af, _ := jsoniter.MarshalToString(serviceList)
	//fmt.Println("排序后--------->", af)
	return serviceList
}
