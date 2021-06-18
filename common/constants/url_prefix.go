package constants

import (
	"a.a/cu/encrypt"
	"fmt"
)

const (
	UrlQrCodeId              = "mp://pay/bizpay"   // 带金额临时二维码
	UrlBusinessFixedQrCodeId = "mp://pay/business" // 商家固定二维码
	UrlPersonalFixedQrCodeId = "mp://pay/personal" // 个人商家固定二维码
)

// 获取待金额的收款二维码
// 格式: mp://pay/bizpay?qr=821f8e60db7eefccf517916a0f6f345f
func GetQrCodeUrl(qrCodeId string) string {
	return fmt.Sprintf("%s?qr=%s", UrlQrCodeId, qrCodeId)
}

// 获取商家固定二维码
// 格式: mp://pay/business?qr=f8026630df48fb5e7830d46dd2e42bd3
func GetBusinessFixedQrCodeUrl(qrCodeId string) string {
	return fmt.Sprintf("%s?qr=%s", UrlBusinessFixedQrCodeId, qrCodeId)
}

// 获取个人商家固定二维码
// 格式: mp://pay/personal?qr=f8026630df48fb5e7830d46dd2e42bd3
func GetPersonalBusinessFixedQrCodeUrl(qrCodeId string) string {
	return fmt.Sprintf("%s?qr=%s", UrlPersonalFixedQrCodeId, qrCodeId)
}

func GetQrCodeId(param string) string {
	salt := "Y0zEyRTW4RIVMUHTvRsWVcvPoIv50dpL"
	return encrypt.DoMd5Salted(param, salt)
}
