package handler

import (
	"a.a/cu/strext"
	billProto "a.a/mp-server/common/proto/bill"
	"context"
	"testing"
)

func TestBillHandler_AddBusinessTransfer(t *testing.T) {
	req := &billProto.AddBusinessTransferRequest{
		BusinessNo:    "515a3299-5af2-4af0-98b2-28d76b6c2223",
		BusinessAccNo: "49307c6f-9e03-4535-b3ef-8aaa581f91bd",
		Amount:        "100",
		CurrencyType:  "usd",
		PayeeNo:       "h13298690108@163.com",
	}
	reply := &billProto.AddBusinessTransferReply{}
	if err := BillHandlerInst.AddBusinessTransfer(context.TODO(), req, reply); err != nil {
		t.Errorf("AddBusinessTransfer() error = %v", err)
		return
	}
	t.Logf("转账结果：%v", strext.ToJson(reply))
}

func TestBillHandler_EnterpriseTransferToUser(t *testing.T) {
	req := &billProto.EnterpriseTransferToUserRequest{
		TransferNo:   "2020102916265953054177",
		TransferType: "2",
	}
	reply := &billProto.EnterpriseTransferToUserReply{}
	if err := BillHandlerInst.EnterpriseTransferToUser(context.TODO(), req, reply); err != nil {
		t.Errorf("EnterpriseTransferToUser() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}
