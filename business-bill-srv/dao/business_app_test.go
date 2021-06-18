package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	"testing"
)

func TestBusinessAppDao_GetAppInfoByFixedQrCode(t *testing.T) {
	fixedQrCode := "038dd0b8289cd06c0e08f9cb50aeda67"
	got, err := BusinessAppDaoInst.GetAppInfoByFixedQrCode(fixedQrCode)
	if err != nil {
		t.Errorf("GetAppInfoByFixedQrCode() error = %v", err)
		return
	}

	t.Logf("app信息：%v", strext.ToJson(got))
}

func TestBusinessAppDao_GetAppInfoByAppId(t *testing.T) {
	appId := ""
	app, err := BusinessAppDaoInst.GetAppInfoByAppId(appId)
	if err != nil {
		t.Errorf("GetAppInfoByAppId() error = %v", err)
		return
	}
	t.Logf("app: %v", strext.ToJson(app))
}
