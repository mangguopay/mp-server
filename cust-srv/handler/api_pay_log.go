package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/global"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/cust-srv/dao"
	"context"
	"time"
)

func (*CustHandler) GetApiPayLogList(ctx context.Context, req *custProto.GetApiPayLogListRequest, reply *custProto.GetApiPayLogListReply) error {

	reqStartTimeStr := ""
	reqEndTimeStr := ""

	if req.ReqStartTime != "" {
		time2, err := time.ParseInLocation(ss_time.DateTimeSlashFormat, req.ReqStartTime, global.Tz)
		if err != nil {
			ss_log.Error("转换时间[%v]为时间戳出错,err[%v]", req.ReqStartTime, err)
			reply.ResultCode = ss_err.ERR_TIMEFORMAT
			return nil
		}
		reqStartTimeStr = strext.ToStringNoPoint(time2.Unix())
	}

	if req.ReqEndTime != "" {
		time2, err := time.ParseInLocation(ss_time.DateTimeSlashFormat, req.ReqEndTime, global.Tz)
		if err != nil {
			ss_log.Error("转换时间[%v]为时间戳出错,err[%v]", req.ReqEndTime, err)
			reply.ResultCode = ss_err.ERR_TIMEFORMAT
			return nil
		}
		reqEndTimeStr = strext.ToStringNoPoint(time2.Unix())
	}

	total, datas, err := dao.ApiPayLogDaoInstance.GetList(dao.ApiPayLogDao{
		Page:      req.Page,
		PageSize:  req.PageSize,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		ReqMethod: req.ReqMethod, //请求方法

		ReqUri:         req.ReqUri,         //请求uri
		ReqBody:        req.ReqBody,        //请求body
		RespData:       req.RespData,       //返回数据
		TrafficStatus:  req.TrafficStatus,  //通信状态(0失败，1成功)
		BusinessStatus: req.BusinessStatus, //业务处理装(0失败，1成功)

		AppId:        req.AppId,       //应用id
		ReqStartTime: reqStartTimeStr, //请求时间开始（要转成时间戳）
		ReqEndTime:   reqEndTimeStr,   //请求时间结束（要转成时间戳）
	})

	if err != nil {
		ss_log.Error("err[%v]", err)
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	var dataList []*custProto.ApiPayLogData
	for _, data := range datas {
		temp := &custProto.ApiPayLogData{
			TraceNo:        data.TraceNo,
			ReqMethod:      data.ReqMethod,
			ReqUri:         data.ReqUri,
			ReqBody:        data.ReqBody,
			RespData:       data.RespData,
			TrafficStatus:  data.TrafficStatus,
			BusinessStatus: data.BusinessStatus,
			CreateTime:     data.CreateTime,
			AppId:          data.AppId,
			ReqTime:        data.ReqTime,
		}

		//将时间戳转换成指定格式的时间字符串
		if temp.ReqTime == "0" {
			temp.ReqTime = ""
		} else {
			temp.ReqTime = time.Unix(strext.ToInt64(temp.ReqTime), 0).Format(ss_time.DateTimeSlashFormat)
		}

		dataList = append(dataList, temp)
	}

	reply.ResultCode = ss_err.ERR_SUCCESS
	reply.Total = strext.ToInt32(total)
	reply.Datas = dataList
	return nil
}
