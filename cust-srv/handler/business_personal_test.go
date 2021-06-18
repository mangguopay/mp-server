package handler

import (
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_GetPersonalBusinessBalance(t *testing.T) {
	req := &custProto.GetPersonalBusinessBalanceRequest{
		AccountNo: "972617f3-c85b-465b-ae3a-8491647d869d",
	}
	reply := &custProto.GetPersonalBusinessBalanceReply{}
	if err := CustHandlerInst.GetPersonalBusinessBalance(context.TODO(), req, reply); err != nil {
		t.Errorf("GetPersonalBusinessBalance() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetPersonalBusinessBills(t *testing.T) {
	req := &custProto.GetPersonalBusinessBillsRequest{
		Page:         0,
		PageSize:     0,
		StartTime:    "",
		EndTime:      "",
		AccountNo:    "972617f3-c85b-465b-ae3a-8491647d869d",
		CurrencyType: "",
	}
	reply := &custProto.GetPersonalBusinessBillsReply{}
	if err := CustHandlerInst.GetPersonalBusinessBills(context.TODO(), req, reply); err != nil {
		t.Errorf("GetPersonalBusinessBills() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetPersonalBusinessFixedCode(t *testing.T) {
	req := &custProto.GetPersonalBusinessFixedCodeRequest{
		AccountNo:   "972617f3-c85b-465b-ae3a-8491647d869d",
		AccountType: "4",
	}
	reply := &custProto.GetPersonalBusinessFixedCodeReply{}
	if err := CustHandlerInst.GetPersonalBusinessFixedCode(context.TODO(), req, reply); err != nil {
		t.Errorf("GetPersonalBusinessFixedCode() error = %v", err)
		return
	}

	//mp://pay/personal?qr=Fmp5MGDxPdFRufAL0CuUcS7NsypJmQWJ
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetPersonalBusinessBillDetail(t *testing.T) {
	req := &custProto.GetPersonalBusinessBillDetailRequest{
		AccountNo: "972617f3-c85b-465b-ae3a-8491647d869d",
		OrderNo:   "2020111015552397120612",
	}
	reply := &custProto.GetPersonalBusinessBillDetailReply{}

	if err := CustHandlerInst.GetPersonalBusinessBillDetail(context.TODO(), req, reply); err != nil {
		t.Errorf("GetPersonalBusinessBillDetail() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetPersonalBusinessInfo(t *testing.T) {
	req := &custProto.GetPersonalBusinessInfoRequest{
		AccountNo: "972617f3-c85b-465b-ae3a-8491647d869d",
	}
	reply := &custProto.GetPersonalBusinessInfoReply{}
	if err := CustHandlerInst.GetPersonalBusinessInfo(context.TODO(), req, reply); err != nil {
		t.Errorf("GetPersonalBusinessInfo() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}
