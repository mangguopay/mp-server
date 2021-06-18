package handler

import (
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_GetAllPaymentChannel(t *testing.T) {
	req := &custProto.GetAllPaymentChannelRequest{}
	reply := &custProto.GetAllPaymentChannelReply{}
	if err := CustHandlerInst.GetAllPaymentChannel(context.TODO(), req, reply); err != nil {
		t.Errorf("GetAllPaymentChannel() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_AddPaymentChannel(t *testing.T) {
	req := &custProto.AddPaymentChannelRequest{
		ChannelName: "测试5",
		ChannelType: "OUT",
		UpstreamNo:  "",
	}
	reply := &custProto.AddPaymentChannelReply{}
	if err := CustHandlerInst.AddPaymentChannel(context.TODO(), req, reply); err != nil {
		t.Errorf("AddPaymentChannel() error = %v", err)
		return
	}
	t.Logf("添加成功")
}
