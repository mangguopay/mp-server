package handler

import (
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_GetFingerprintList(t *testing.T) {
	req := &custProto.GetFingerprintListRequest{
		Page:      0,
		PageSize:  0,
		StartTime: "",
		EndTime:   "",
		Account:   "0855789456123",
		DeviceNo:  "",
		UseStatus: "",
	}
	reply := &custProto.GetFingerprintListReply{}

	err := CustHandlerInst.GetFingerprintList(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetFingerprintList() error = %v", err)
		return
	}

	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_CloseFingerprintFunction(t *testing.T) {
	req := &custProto.CloseFingerprintFunctionRequest{
		LoginUid: "829b7c29-4bed-4d58-8010-716d2e14f1b8",
		IsOpen:   "true",
	}
	reply := &custProto.CloseFingerprintFunctionReply{}
	if err := CustHandlerInst.CloseFingerprintFunction(context.TODO(), req, reply); err != nil {
		t.Errorf("CloseFingerprintFunction() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_CleanFingerprintData(t *testing.T) {
	req := &custProto.CleanFingerprintDataRequest{
		LoginUid: "829b7c29-4bed-4d58-8010-716d2e14f1b8",
		OpType:   "all",
		Id:       "2020102614374367718404",
	}
	reply := &custProto.CleanFingerprintDataReply{}
	if err := CustHandlerInst.CleanFingerprintData(context.TODO(), req, reply); err != nil {
		t.Errorf("CleanFingerprintData() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetAppFingerprintOn(t *testing.T) {
	req := &custProto.GetAppFingerprintOnRequest{}
	reply := &custProto.GetAppFingerprintOnReply{}

	if err := CustHandlerInst.GetAppFingerprintOn(context.TODO(), req, reply); err != nil {
		t.Errorf("GetAppFingerprintOn() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}
