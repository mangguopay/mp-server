package common

const (
	ReturnSuccess = "success"
)

const (
	//通知类型
	NotifyTypeToPayment  = "PAYMENT"
	NotifyTypeToRefund   = "REFUND"
	NotifyTypeToTransfer = "TRANSFER"
)

const (
	//数据签名参数key
	NotifySignField = "sign"

	//通知状态
	NotifyStatusPending = 0 //进行中
	NotifyStatusSuccess = 1 //成功
	NotifyStatusFail    = 2 //失败
)
