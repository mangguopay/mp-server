package channel

import (
	"a.a/mp-server/business-bill-srv/model"
	"testing"
)

func TestGetChannelAPi(t *testing.T) {
	channelNo := "01"
	channel, err := GetChannelAPi(channelNo)
	if err != nil {
		t.Logf("获取API失败，err=%v", err)
		return
	}
	ret := channel.ScanPay(&model.ScanPayRequest{})
	t.Logf("测试结果：%v", ret)

}
