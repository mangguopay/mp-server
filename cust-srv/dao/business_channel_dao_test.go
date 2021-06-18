package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestBusinessChannel_GetAllChannel(t *testing.T) {
	channels, err := BusinessChannelDao.GetAllChannel("", nil)
	if err != nil {
		t.Errorf("GetAllChannel() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(channels))

}

func TestBusinessChannel_Insert(t *testing.T) {
	d := &BusinessChannel{
		ChannelNo:   strext.GetDailyId(),
		ChannelName: "测试2",
		ChannelType: constants.ChannelTypeOut,
		UpstreamNo:  "",
	}
	if err := BusinessChannelDao.InsertTx(nil, d); err != nil {
		t.Errorf("Insert() error = %v", err)
		return
	}
	t.Logf("添加成功")
}

func TestBusinessChannel_GetChannelAndRate(t *testing.T) {
	appId := "2020090417003775361070"
	sceneNo := "f4adfcd4-c490-4750-b5b9-80637dc1745c"
	channel, err := BusinessChannelDao.GetChannelRate(appId, sceneNo)
	if err != nil {
		t.Errorf("GetChannelAndRate() error = %v", err)
		return
	}
	t.Logf("渠道：%v, ", strext.ToJson(channel))
}

func TestBusinessChannel_GetOutChannelRateAndName(t *testing.T) {
	appId := "2020090417003775361070"
	sceneNo := "f4adfcd4-c490-4750-b5b9-80637dc1745c"
	channelName, channelRate, err := BusinessChannelDao.GetOutChannelRateAndName(appId, sceneNo)
	if err != nil {
		t.Errorf("GetRate() error = %v", err)
		return
	}
	t.Logf("渠道：%v, 费率：%v", channelName, channelRate)
}

func TestBusinessChannel_Update(t *testing.T) {
	d := new(BusinessChannel)
	d.ChannelNo = "2020091610181255429942"
	d.ChannelName = "测试4"
	d.ChannelType = ""
	d.UpstreamNo = ""
	if err := BusinessChannelDao.UpdateTx(nil, d); err != nil {
		t.Errorf("Update() error = %v", err)
		return
	}
	t.Logf("修改成功")
}
