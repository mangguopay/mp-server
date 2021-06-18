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

type TransferNotifyHandler struct {
}

var TransferNotifyH TransferNotifyHandler

func (t *TransferNotifyHandler) TransferSuccessNotify(logNo string) {
	order, err := dao.BusinessTransferOrderDaoInst.GetTransferOrderByLogNo(logNo)
	if err != nil || order == nil {
		ss_log.Error("查询订单信息失败, order_no=[%v], err=[%v]", logNo, err)
		return
	}

	//通知状态, 超时
	if order.NotifyStatus == constants.NotifyStatusTimeout {
		ss_log.Error("订单已超出通知时间, order_no=[%v]", logNo)
		return
	}

	if order.NotifyStatus == constants.NotifyStatusSuccess {
		ss_log.Error("订单已成功通知, order_no=[%v]", logNo)
		return
	}

	if order.NotifyUrl == "" {
		ss_log.Error("企业转账订单[%v]缺少回调地址，不再进行通知", logNo)
		//插入通知日志
		logData := &dao.BusinessNotifyLog{
			OrderNo:    logNo,
			OutOrderNo: order.OutTransferNo,
			Status:     common.NotifyStatusFail,
			Result:     fmt.Sprintf("订单[%v]缺少回调地址", logNo),
			OrderType:  constants.VaReason_BusinessTransferToBusiness,
		}
		_, err := dao.BusinessNotifyLogDaoInst.InsertLog(logData)
		if err != nil {
			ss_log.Error("插入订单异步通知日志失败,InputParam=[%v], err=[%v]", strext.ToString(logData), err)
			return
		}

		updateOrder := new(dao.UpdateTransferNotifyStatus)
		updateOrder.OrderNo = logNo
		updateOrder.NotifyStatus = constants.NotifyStatusTimeout
		err = dao.BusinessTransferOrderDaoInst.UpdateNotifyStatusByOrderNo(updateOrder)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", logNo, err)
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
	packReq := &PackTransferParamRequest{
		AppId:      order.AppId,
		SignType:   appInfo.SignMethod,
		Timestamp:  ss_time.Now(global.Tz).Unix(),
		NotifyType: common.NotifyTypeToTransfer,

		OutOrderNo:   order.OutTransferNo,
		OrderNo:      order.LogNo,
		Amount:       strext.ToString(order.Amount),
		CurrencyType: order.CurrencyType,
		OrderStatus:  order.OrderStatus,
		PayTime:      ss_time.ParseTimeFromPostgres(order.FinishTime, global.Tz).Unix(),

		PlatformPrivateKey: appInfo.PlatformPrivKey,
	}
	postData, err := t.packTransferParamTransfer(packReq)
	if err != nil {
		ss_log.Error("拼装数据失败, err=[%v]", err)
		return
	}

	ss_log.Info("转账订单第%v次通知, order_no=[%v], failed_times=[%v]", order.NotifyFailTimes+1, logNo, order.NotifyFailTimes)

	//插入通知日志
	var logId string
	logData := &dao.BusinessNotifyLog{
		OrderNo:    logNo,
		OutOrderNo: order.OutTransferNo,
		Status:     common.NotifyStatusPending,
		OrderType:  constants.VaReason_BusinessTransferToBusiness,
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
	updateDate := new(dao.UpdateTransferNotifyStatus)
	updateDate.OrderNo = logNo

	if resultStr != common.ReturnSuccess {
		//通知失败 日志状态
		status = common.NotifyStatusFail
		//下一次通知的等待时间
		waitTime := common.GetNotifyWaitTimeById(order.NotifyFailTimes + 1)
		nextTime := ss_time.Now(global.Tz).Add(time.Duration(waitTime) * time.Second).Format(ss_time.DateTimeDashFormat)

		updateDate.NextTime = nextTime
		updateDate.NotifyStatus = constants.NotifyStatusDoing
		updateDate.NotifyFailTimes = 1
		err := dao.BusinessTransferOrderDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", logNo, err)
		}

		if waitTime == -1 { //超时
			updateDate.NotifyStatus = constants.NotifyStatusTimeout
			updateDate.NextTime = ""
			updateDate.NotifyFailTimes = 0
			err := dao.BusinessTransferOrderDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
			if err != nil {
				ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", logNo, err)
			}

		} else {
			err := cache.RedisClient.SetNX(common.GetTransferNotifyExpireKey(logNo), order.NotifyFailTimes, time.Duration(waitTime)*time.Second).Err()
			if err != nil {
				ss_log.Error("redis插入数据失败, key=%v, err=%v", common.GetTransferNotifyExpireKey(logNo), err)
			}
		}
	} else {
		//通知成功
		updateDate.NotifyStatus = constants.NotifyStatusSuccess
		err := dao.BusinessTransferOrderDaoInst.UpdateNotifyStatusByOrderNo(updateDate)
		if err != nil {
			ss_log.Error("修改订单通知状态错误,order_no=[%v], err=[%v]", logNo, err)
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

	ss_log.Info("转账订单第%v次通知结束, ", order.NotifyFailTimes+1)

}

type PackTransferParamRequest struct {
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

func (t *TransferNotifyHandler) packTransferParamTransfer(req *PackTransferParamRequest) (string, error) {
	postData := map[string]interface{}{
		"app_id":          req.AppId,
		"notify_type":     req.NotifyType,
		"sign_type":       req.SignType,
		"timestamp":       fmt.Sprintf("%v", ss_time.Now(global.Tz).Unix()),
		"out_transfer_no": req.OutOrderNo,
		"transfer_no":     req.OrderNo,
		"transfer_status": req.OrderStatus,
		"amount":          req.Amount,
		"currency_type":   req.CurrencyType,
		"transfer_time":   req.PayTime,
		"subject":         req.Subject,
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
