package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_rsa"
	"a.a/mp-server/notify-srv/common"
	"a.a/mp-server/notify-srv/dao"
)

type RefundNotifyHandler struct {
}

var RefundNotifyH RefundNotifyHandler

func (r *RefundNotifyHandler) RefundNotify(refundNo string) {
	order, err := dao.BusinessRefundOrderDaoInst.GetRefundDetail(refundNo)
	if err != nil || order == nil {
		ss_log.Error("查询退款订单信息失败, refund_no=[%v], err=[%v]", refundNo, err)
		return
	}

	//通知状态, 超时
	if order.NotifyStatus == constants.NotifyStatusTimeout {
		ss_log.Error("订单已超出通知时间, refund_no=[%v]", refundNo)
		return
	}

	if order.NotifyStatus == constants.NotifyStatusSuccess {
		ss_log.Error("订单已成功通知, refund_no=[%v]", refundNo)
		return
	}

	if order.NotifyUrl == "" {
		ss_log.Error("退款订单[%v]缺少回调地址，不再进行通知", refundNo)
		//插入通知日志
		logData := &dao.BusinessNotifyLog{
			OrderNo:    refundNo,
			OutOrderNo: order.OutRefundNo,
			Status:     common.NotifyStatusFail,
			Result:     fmt.Sprintf("订单[%v]缺少回调地址", refundNo),
			OrderType:  constants.VaReason_BusinessRefund,
		}
		_, err := dao.BusinessNotifyLogDaoInst.InsertLog(logData)
		if err != nil {
			ss_log.Error("插入退款订单异步通知日志失败,InputParam=[%v], err=[%v]", strext.ToString(logData), err)
			return
		}

		updateOrder := new(dao.UpdateRefundNotifyStatus)
		updateOrder.OrderNo = refundNo
		updateOrder.NotifyStatus = constants.NotifyStatusTimeout
		err = dao.BusinessRefundOrderDaoInst.UpdateNotifyStatusByOrderNo(updateOrder)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,refund_no=[%v], err=[%v]", refundNo, err)
		}
		return
	}

	appInfo, err := dao.BusinessAppDaoInst.GetSignInfo(order.AppId)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("应用记录不存在|appId:%s", order.AppId)
			return
		} else {
			ss_log.Error("查询应用记录失败|appId:%s, err:%v", order.AppId, err)
			return
		}
	}

	//打包参数
	packReq := &PackRefundParamRequest{
		AppId:      order.AppId,
		SignType:   appInfo.SignMethod,
		Timestamp:  ss_time.Now(global.Tz).Unix(),
		NotifyType: common.NotifyTypeToRefund,

		OrderNo:            order.OrderNo,
		OutOrderNo:         order.OutOrderNo,
		TransAmount:        order.TransAmount,
		RefundNo:           order.RefundNo,
		OutRefundNo:        order.OutRefundNo,
		RefundAmount:       order.RefundAmount,
		CurrencyType:       order.CurrencyType,
		RefundStatus:       order.RefundStatus,
		RequestTime:        ss_time.ParseTimeFromPostgres(order.CreateTime, global.Tz).Unix(),
		FinishTime:         ss_time.ParseTimeFromPostgres(order.FinishTime, global.Tz).Unix(),
		PlatformPrivateKey: appInfo.PlatformPrivKey,
	}
	postData, err := r.packRefundParam(packReq)
	if err != nil {
		ss_log.Error("拼装数据失败, err=[%v]", err)
		return
	}

	ss_log.Info("退款订单第%v次通知, refund_no=[%v], failed_times=[%v]", order.NotifyFailTimes+1, refundNo, order.NotifyFailTimes)

	//插入通知日志
	var logId string
	logData := &dao.BusinessNotifyLog{
		OrderNo:    refundNo,
		OutOrderNo: order.OutOrderNo,
		Status:     common.NotifyStatusPending,
		OrderType:  constants.VaReason_BusinessRefund,
	}
	logId, err = dao.BusinessNotifyLogDaoInst.InsertLog(logData)
	if err != nil {
		ss_log.Error("插入订单异步通知日志失败,InputParam=[%v], err=[%v]", strext.ToString(logData), err)
		return
	}

	//发送通知,接收通知结果
	result, err := common.PostSend(order.NotifyUrl, postData)
	if err != nil {
		ss_log.Error("PostSend错误, err=[%v]", err)
		result = []byte(err.Error())
	}
	resultStr := strext.ToString(result)
	resultMap := map[string]string{
		"url": order.NotifyUrl,
		"res": resultStr,
	}

	var status = common.NotifyStatusSuccess
	updateDate := new(dao.UpdateRefundNotifyStatus)
	updateDate.OrderNo = refundNo

	if resultStr != common.ReturnSuccess {
		//通知失败 日志状态
		status = common.NotifyStatusFail
		//下一次通知的等待时间
		waitTime := common.GetNotifyWaitTimeById(order.NotifyFailTimes + 1)
		nextTime := ss_time.Now(global.Tz).Add(time.Duration(waitTime) * time.Second).Format(ss_time.DateTimeDashFormat)

		updateDate.NextTime = nextTime
		updateDate.NotifyStatus = constants.NotifyStatusDoing
		updateDate.NotifyFailTimes = 1
		err := dao.BusinessRefundOrderDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,refund_no=[%v], err=[%v]", refundNo, err)
		}

		if waitTime == -1 { //超时
			updateDate.NotifyStatus = constants.NotifyStatusTimeout
			updateDate.NextTime = ""
			updateDate.NotifyFailTimes = 0
			err := dao.BusinessRefundOrderDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
			if err != nil {
				ss_log.Error("修改订单通知状态错误,refund_no=[%v], err=[%v]", refundNo, err)
			}

		} else {
			err := cache.RedisClient.SetNX(common.GetRefundNotifyExpireKey(refundNo), order.NotifyFailTimes,
				time.Duration(waitTime)*time.Second).Err()
			if err != nil {
				ss_log.Error("redis插入数据失败, key=%v, err=%v", common.GetTransferNotifyExpireKey(refundNo), err)
			}
		}
	} else {
		//通知成功
		updateDate.NotifyStatus = constants.NotifyStatusSuccess
		err := dao.BusinessRefundOrderDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,refund_no=[%v], err=[%v]", refundNo, err)
		}
	}

	//修改日志的通知结果
	logData = &dao.BusinessNotifyLog{
		LogId:  logId,
		Result: strext.ToJson(resultMap),
		Status: status,
	}
	if updateLogErr := dao.BusinessNotifyLogDaoInst.UpdateNotifyResultById(logData); updateLogErr != nil {
		ss_log.Error("修改订单通知状态错误,InputParam=[%v], err=[%v]", strext.ToString(logData), updateLogErr)
	}

	ss_log.Info("退款订单第%v次通知结束, ", order.NotifyFailTimes+1)

}

type PackRefundParamRequest struct {
	AppId      string //AppId
	SignType   string //数据签名方式
	Timestamp  int64  //时间戳
	NotifyType string //通知类型

	OrderNo      string //平台交易单号
	OutOrderNo   string //外部交易单号
	TransAmount  string //交易金额
	RefundNo     string //平台退款订单号
	OutRefundNo  string //外部退款订单号
	RefundAmount string //退款金额
	CurrencyType string //币种
	RefundStatus string //订单状态
	RequestTime  int64  //退款创建时间(时间戳)
	FinishTime   int64  //退款完成时间(时间戳)

	PlatformPrivateKey string //平台私钥
}

func (r *RefundNotifyHandler) packRefundParam(req *PackRefundParamRequest) (string, error) {
	postData := map[string]interface{}{
		"app_id":      req.AppId,
		"sign_type":   req.SignType,
		"timestamp":   fmt.Sprintf("%v", req.Timestamp),
		"notify_type": req.NotifyType,

		"order_no":      req.OrderNo,
		"out_order_no":  req.OutOrderNo,
		"trans_amount":  req.TransAmount,
		"refund_no":     req.RefundNo,
		"out_refund_no": req.OutRefundNo,
		"refund_amount": req.RefundAmount,
		"currency_type": req.CurrencyType,
		"refund_status": req.RefundStatus,
		"refund_time":   fmt.Sprintf("%v", req.FinishTime),
	}
	//签名
	switch req.SignType {
	case "RSA2":
		bodyStr := ss_func.ParamsMapToString(postData, common.NotifySignField)
		signRet, err := ss_rsa.RSA2Sign(bodyStr, req.PlatformPrivateKey)
		if err != nil {
			ss_log.Error("签名失败|err:%v|data:%v|privateKey:%v", err, postData, req.PlatformPrivateKey)
			return "", err
		}
		// 将签名添加到尾部
		bodyStr += fmt.Sprintf("&%s=%s", common.NotifySignField, signRet)
		return bodyStr, nil
	case "MD5":
		ss_log.Info("暂不支持MD5签名")
		return "", errors.New("暂不支持MD5签名")
	}
	return "", errors.New("签名类型不存在")
}
