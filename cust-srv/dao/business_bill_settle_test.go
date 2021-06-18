package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestBusinessBillSettle_GetSettleLogById(t *testing.T) {
	settleId := "2020082519265303881827"
	got, err := BusinessBillSettleDaoInst.GetSettleLogById(settleId)
	if err != nil {
		t.Errorf("GetSettleLogById() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(got))
}

func TestBusinessBillSettle_GetSingleSettleDetail(t *testing.T) {
	settleId := "2020090716171910268764"
	got, err := BusinessBillSettleDaoInst.GetSingleSettleDetail(settleId)
	if err != nil {
		t.Errorf("GetSettleLogById() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(got))
}
