package dao

import (
	"a.a/cu/strext"
	"testing"
)

func TestBusinessDao_GetBusinessStatusInfo(t *testing.T) {
	//account := "h13298690108@163.com"
	accountNo := "555d2d86-fef4-42d9-b2f0-a6adb8c3f325"
	got, err := BusinessDaoInst.GetBusinessStatusInfo("", accountNo)
	if err != nil {
		t.Errorf("GetBusinessStatusInfo() error = %v", err)
		return
	}
	t.Logf(strext.ToJson(got))
}
