package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	"testing"
)

func TestBusinessRefundOrderDao_GetRefundByPayOrderNo(t *testing.T) {
	payOrderNo := "2020081911113290077480"
	ret, err := BusinessRefundOrderDaoInst.GetRefundByPayOrderNo(payOrderNo)
	if err != nil {
		t.Errorf("GetRefundByPayOrderNo() error = %v", err)
		return
	}
	t.Logf("查询结果: %v", strext.ToJson(ret))

}

func TestBusinessRefundOrderDao_GetRefundByRefundNo(t *testing.T) {
	appId := "2020081211021112759835"
	refundNo := "2020082719372727728603"
	outRefundNo := ""
	ret, err := BusinessRefundOrderDaoInst.GetRefundByRefundNo(appId, refundNo, outRefundNo)
	if err != nil {
		t.Errorf("GetRefundByRefundNo() error = %v", err)
		return
	}
	t.Logf("查询结果: %v", strext.ToJson(ret))
}

func TestBusinessRefundOrderDao_UpdateTimeOutOrderStatus(t *testing.T) {
	list, err := BusinessRefundOrderDaoInst.UpdateTimeOutOrderStatus()
	if err != nil {
		t.Errorf("UpdateTimeOutOrderStatus() error = %v", err)
		return
	}
	t.Logf("修改订单列表：%v", strext.ToJson(list))
}
