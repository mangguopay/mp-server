package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/model"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestBusinessBillDao_GetBusinessBills(t *testing.T) {
	whereList := []*model.WhereSqlCond{
		{Key: "bb.business_no", Val: "53be2111-b1cb-4041-8143-d4a2ccf7d995", EqType: "="},
	}
	page := 1
	pageSize := 10

	list, err := BusinessBillDaoInst.GetBusinessBills(whereList, page, pageSize)
	if err != nil {
		t.Errorf("GetBusinessBills() error = %v", err)
		return
	}
	t.Logf("list: %v", strext.ToJson(list))
}

func TestBusinessBillDao_GetBusinessBillDetail(t *testing.T) {
	orderNo := "2020100914293281633801"
	data, err := BusinessBillDaoInst.GetBusinessBillDetail(orderNo)
	if err != nil {
		t.Errorf("GetBusinessBillDetail() error = %v", err)
		return
	}
	t.Logf("订单详情：%v", strext.ToJson(data))
}

func TestBusinessBillDao_GetBusinessChannelBills(t *testing.T) {
	whereList := []*model.WhereSqlCond{
		//{Key: "bb.business_no", Val: "53be2111-b1cb-4041-8143-d4a2ccf7d995", EqType: "="},
	}
	page := 1
	pageSize := 10
	ret, err := BusinessBillDaoInst.GetBusinessChannelBills(whereList, page, pageSize)
	if err != nil {
		t.Errorf("GetBusinessChannelBills() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(ret))
}
