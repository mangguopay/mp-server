package dao

import (
	"testing"

	"a.a/cu/strext"
)

func TestApp_Insert(t *testing.T) {
	InitDB()

	order := &App{
		AppId:              "1234567890",
		AppName:            "测试插入",
		MerchantPrivateKey: "merchantPrivateKey",
		MerchantPublicKey:  "merchantPublicKey",
		PlatformPublicKey:  "platformPublicKey",
		MerchantKeyType:    "PKCS1",
	}

	if err := AppInstance.Insert(order); err != nil {
		t.Errorf("插入应用失败, err:%v, order:%+v", err, order)
		return
	}

	t.Logf("插入应用成功,AppId:%s", order.AppId)
}

func TestApp_GetList(t *testing.T) {
	InitDB()

	page := 2
	pageSize := 5

	list, err := AppInstance.GetList(page, pageSize)

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("list:%s", strext.ToJson(list))
}

func TestApp_GetUsingRow(t *testing.T) {
	InitDB()

	row, err := AppInstance.GetUsingRow()

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("row:%s", strext.ToJson(row))
}

func TestApp_SetAppUsing(t *testing.T) {
	InitDB()

	appId := "2020081115122507673612"
	err := AppInstance.SetAppUsing(appId)

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("ok")
}

func TestApp_GetOneByAppId(t *testing.T) {
	InitDB()

	appId := "2020090416495834598604"

	app, err := AppInstance.GetOneByAppId(appId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("appInfo:%s", strext.ToJson(app))
}
