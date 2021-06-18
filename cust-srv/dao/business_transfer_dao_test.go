package dao

import (
	"a.a/cu/strext"
	"testing"
)

func TestBusinessTransferDao_GetOrderDetail(t *testing.T) {
	logNo := "2020081014395392363596"
	got, err := BusinessTransferDaoInst.GetOrderDetail(logNo)
	if err != nil {
		t.Errorf("GetOrderDetail() error = %v", err)
		return
	}

	t.Logf("详情：%v", strext.ToJson(got))
}
