package handler

import (
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_GetBusinessTransferOrderList(t *testing.T) {
	req := &custProto.GetBusinessTransferOrderListRequest{
		BusinessAccNo: "555d2d86-fef4-42d9-b2f0-a6adb8c3f325",
		CurrencyType:  "usd",
		TransferType:  "2",
	}
	reply := &custProto.GetBusinessTransferOrderListReply{}

	if err := CustHandlerInst.GetBusinessTransferOrderList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessTransferOrderList() error = %v", err)
		return
	}
	t.Logf("ret:%v", strext.ToJson(reply))
}

func TestCustHandler_GetBusinessTransferOrderDetail(t *testing.T) {
	req := &custProto.GetBusinessTransferOrderDetailRequest{
		LogNo: "2020081215154216913138",
	}
	reply := &custProto.GetBusinessTransferOrderDetailReply{}

	if err := CustHandlerInst.GetBusinessTransferOrderDetail(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessTransferOrderDetail() error = %v", err)
		return
	}

	t.Logf("商家转账订单详情：%v", strext.ToJson(reply))
}
