package handler

import (
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
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type PayNotifyHandler struct {
}

var PayNotifyH PayNotifyHandler

func (p *PayNotifyHandler) PaySuccessNotify(orderNo string) {
	orderInfo, err := dao.BusinessBillDaoInst.QueryOrderInfoByOrderNo(orderNo)
	if err != nil || orderInfo == nil {
		ss_log.Error("查询订单信息失败, order_no=[%v], err=[%v]", orderNo, err)
		return
	}

	//通知状态, 超时
	if orderInfo.NotifyStatus == constants.NotifyStatusTimeout {
		ss_log.Error("订单已超出通知时间, order_no=[%v]", orderNo)
		return
	}

	if orderInfo.NotifyStatus == constants.NotifyStatusSuccess {
		ss_log.Error("订单已成功通知, order_no=[%v]", orderNo)
		return
	}

	if orderInfo.NotifyUrl == "" {
		ss_log.Error("普通交易订单[%v]缺少回调地址，不再进行通知", orderNo)
		//插入通知日志
		logData := &dao.BusinessNotifyLog{
			OrderNo:    orderNo,
			OutOrderNo: orderInfo.OutOrderNo,
			Status:     common.NotifyStatusFail,
			Result:     fmt.Sprintf("订单[%v]缺少回调地址", orderNo),
			OrderType:  constants.VaReason_Cust_Pay_Order,
		}
		_, err := dao.BusinessNotifyLogDaoInst.InsertLog(logData)
		if err != nil {
			ss_log.Error("插入订单异步通知日志失败,InputParam=[%v], err=[%v]", strext.ToString(logData), err)
			return
		}

		updateOrder := new(dao.UpdateBillNotifyStatus)
		updateOrder.OrderNo = orderNo
		updateOrder.NotifyStatus = constants.NotifyStatusTimeout
		err = dao.BusinessBillDaoInst.UpdateNotifyStatusByOrderNo(updateOrder)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", orderNo, err)
		}
		return
	}

	appInfo, err := dao.BusinessAppDaoInst.GetSignInfo(orderInfo.AppId)
	if err != nil {
		if err == sql.ErrNoRows {
			ss_log.Error("应用记录不存在|appId:%s", orderInfo.AppId)
			return
		} else {
			ss_log.Error("查询应用记录失败|appId:%s, err:%v", orderInfo.AppId, err)
			return
		}
	}

	//打包参数
	packReq := &PackPayParamRequest{
		AppId:      orderInfo.AppId,
		SignType:   appInfo.SignMethod,
		Timestamp:  ss_time.Now(global.Tz).Unix(),
		NotifyType: common.NotifyTypeToPayment,

		OutOrderNo:   orderInfo.OutOrderNo,
		OrderNo:      orderInfo.OrderNo,
		Amount:       orderInfo.Amount,
		CurrencyType: orderInfo.CurrencyType,
		OrderStatus:  orderInfo.OrderStatus,
		PayTime:      ss_time.ParseTimeFromPostgres(orderInfo.PayTime, global.Tz).Unix(),
		Subject:      orderInfo.Subject,

		PlatformPrivateKey: appInfo.PlatformPrivKey,
	}
	postData, err := p.packPayParam(packReq)
	if err != nil {
		ss_log.Error("拼装数据失败, err=[%v]", err)
		return
	}

	ss_log.Info("支付成功订单第%v次通知, order_no=[%v], failed_times=[%v]", orderInfo.NotifyFailTimes+1, orderNo, orderInfo.NotifyFailTimes)

	//插入通知日志
	var logId string
	logData := &dao.BusinessNotifyLog{
		OrderNo:    orderNo,
		OutOrderNo: orderInfo.OutOrderNo,
		Status:     common.NotifyStatusPending,
		OrderType:  constants.VaReason_Cust_Pay_Order,
	}
	logId, err = dao.BusinessNotifyLogDaoInst.InsertLog(logData)
	if err != nil {
		ss_log.Error("插入订单异步通知日志失败,InputParam=[%v], err=[%v]", strext.ToString(logData), err)
		return
	}

	//发送通知,接收通知结果
	result, err := common.PostSend(orderInfo.NotifyUrl, postData)
	if err != nil {
		ss_log.Error("PostSend错误, err=[%v]", err)
		result = []byte(err.Error())
	}
	resultStr := strext.ToString(result)
	resultMap := map[string]string{
		"url": orderInfo.NotifyUrl,
		"res": resultStr,
	}

	var status = common.NotifyStatusSuccess
	updateDate := new(dao.UpdateBillNotifyStatus)
	updateDate.OrderNo = orderNo

	if resultStr != common.ReturnSuccess {
		//通知失败 日志状态
		status = common.NotifyStatusFail
		//下一次通知的等待时间
		waitTime := common.GetNotifyWaitTimeById(orderInfo.NotifyFailTimes + 1)
		nextTime := ss_time.Now(global.Tz).Add(time.Duration(waitTime) * time.Second).Format(ss_time.DateTimeDashFormat)

		updateDate.NextTime = nextTime
		updateDate.NotifyStatus = constants.NotifyStatusDoing
		updateDate.NotifyFailTimes = 1

		err := dao.BusinessBillDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", orderNo, err)
		}

		if waitTime == -1 { //超时
			updateDate.NotifyStatus = constants.NotifyStatusTimeout
			updateDate.NextTime = ""
			updateDate.NotifyFailTimes = 0
			err := dao.BusinessBillDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
			if err != nil {
				ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", orderNo, err)
			}

		} else {
			err := cache.RedisClient.SetNX(common.GetPayNotifyExpireKey(orderNo), orderInfo.NotifyFailTimes, time.Duration(waitTime)*time.Second).Err()
			if err != nil {
				ss_log.Error("redis插入数据失败, key=%v, err=%v", common.GetPayNotifyExpireKey(orderNo), err)
			}
		}
	} else {
		//通知成功
		updateDate.NotifyStatus = constants.NotifyStatusSuccess
		err := dao.BusinessBillDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", orderNo, err)
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

	ss_log.Info("支付成功订单第%v次通知结束, ", orderInfo.NotifyFailTimes+1)

}

type PackPayParamRequest struct {
	AppId      string //AppId
	SignType   string //数据签名方式
	Timestamp  int64  //时间戳
	NotifyType string //通知类型

	OutOrderNo   string //外部订单号
	OrderNo      string //平台订单号
	Amount       string //金额
	CurrencyType string //币种
	OrderStatus  string //订单状态
	PayTime      int64  //支付时间(时间戳)
	Subject      string //商品名称

	PlatformPrivateKey string //平台私钥
}

func (PayNotifyHandler) packPayParam(req *PackPayParamRequest) (string, error) {
	postData := map[string]interface{}{
		"app_id":      req.AppId,
		"sign_type":   req.SignType,
		"timestamp":   fmt.Sprintf("%v", req.Timestamp),
		"notify_type": req.NotifyType,

		"out_order_no":  req.OutOrderNo,
		"order_no":      req.OrderNo,
		"amount":        req.Amount,
		"currency_type": req.CurrencyType,
		"order_status":  req.OrderStatus,
		"pay_time":      fmt.Sprintf("%v", req.PayTime),
		"subject":       req.Subject,
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
