package dao

import (
	"testing"
)

func TestBusinessAppDao_GetSignInfo(t *testing.T) {
	appId := ""
	ret, err := BusinessAppDaoInst.GetSignInfo(appId)
	if err != nil {
		t.Errorf("GetSignInfo() error = %v", err)
		return
	}
	t.Logf("ret:%v", ret)
}
