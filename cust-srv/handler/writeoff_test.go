package handler

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_GetWriteOffList(t *testing.T) {
	req := &custProto.GetWriteOffListRequest{
		Page:         1,
		PageSize:     10,
		StartTime:    "",
		EndTime:      "",
		Code:         "",
		PayerAccount: "",
		PayeeAccount: "",
		UseStatus:    "",
	}
	reply := &custProto.GetWriteOffListReply{}
	if err := CustHandlerInst.GetWriteOffList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetWriteOffList() error = %v", err)
		return
	}

	t.Logf("结果：%v", strext.ToJson(reply))
}

func TestCustHandler_DisposeWriteOffCode(t *testing.T) {
	req := &custProto.DisposeWriteOffCodeRequest{
		Code:     "1354408099",
		OpType:   constants.WrittenOffCodeOpFreeze,
		LoginUid: "f5f9fe55-ce7c-4cad-80eb-682f8b41e87d",
	}
	reply := &custProto.DisposeWriteOffCodeReply{}
	if err := CustHandlerInst.DisposeWriteOffCode(context.TODO(), req, reply); err != nil {
		t.Errorf("DisposeWriteOffCode() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}
