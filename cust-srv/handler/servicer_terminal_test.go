package handler

import (
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"testing"
)

func TestCustHandler_AddTerminal(t *testing.T) {
	req := &custProto.AddTerminalRequest{
		ServicerAccount: "085513888888888",
		TerminalNumber:  "98211909306922",
		PosSn:           "98211909306922",
		UseStatus:       "1",
		LoginAccount:    "",
	}
	reply := &custProto.AddTerminalReply{}
	if err := CustHandlerInst.AddTerminal(context.TODO(), req, reply); err != nil {
		t.Errorf("AddServicerPOS() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))

}

func TestCustHandler_GetTerminalList(t *testing.T) {
	req := &custProto.GetTerminalListRequest{}
	reply := &custProto.GetTerminalListReply{}
	if err := CustHandlerInst.GetTerminalList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetTerminalList() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_UpdateTerminal(t *testing.T) {
	req := &custProto.UpdateTerminalRequest{
		TerminalNo:   "2020092813582210165478",
		UseStatus:    "1",
		LoginAccount: "22222222-2222-2222-2222-222222222222",
	}
	reply := &custProto.UpdateTerminalReply{}
	if err := CustHandlerInst.UpdateTerminal(context.TODO(), req, reply); err != nil {
		t.Errorf("UpdateTerminal() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}
