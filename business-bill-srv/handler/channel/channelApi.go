package channel

import (
	"a.a/mp-server/business-bill-srv/model"
	"errors"
)

type PaymentChannels interface {
	//扫码支付
	ScanPay(request *model.ScanPayRequest) *model.ScanPayResponse
	//申请二维码
	PreCreate(request *model.PreCreateRequest) *model.PreCreateResponse
}

func GetChannelAPi(channelNo string) (PaymentChannels, error) {
	channelApi := GetChannel(channelNo)
	if channelApi == nil {
		return nil, errors.New("支付渠道不存在")
	}
	return channelApi, nil
}

func GetChannel(channelNo string) PaymentChannels {
	switch channelNo {
	case "01":
		return &ModernPay{}
	}

	return nil
}
