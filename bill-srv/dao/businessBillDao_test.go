package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/bill-srv/test"
	"testing"
)

func Test_businessBillDao_GetBusinessBillDetail(t *testing.T) {
	gotDataT, err := BusinessBillDaoInst.GetBusinessBillDetail("2020091417401430956415")
	if err != nil {
		t.Errorf("GetBusinessBillDetail() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(gotDataT))
}
