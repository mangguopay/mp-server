package common

import "fmt"

const (
	//redis key
	PayNotifyExpireKey      = "payNotifyFail"
	TransferNotifyExpireKey = "transferNotifyFail"
	RefundNotifyExpireKey   = "refundNotifyFail"
	ExpireKeyLock           = "ExpireKeyLock"
)

func GetPayNotifyExpireKey(orderNo string) string {
	return fmt.Sprintf("%s_%s", PayNotifyExpireKey, orderNo)
}

func GetTransferNotifyExpireKey(orderNo string) string {
	return fmt.Sprintf("%s_%s", TransferNotifyExpireKey, orderNo)
}

func GetRefundNotifyExpireKey(refundNo string) string {
	return fmt.Sprintf("%s_%s", RefundNotifyExpireKey, refundNo)
}

func GetLockKey(orderNo string) string {
	return fmt.Sprintf("%s_%s", ExpireKeyLock, orderNo)
}
