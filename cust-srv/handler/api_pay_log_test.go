package handler

import (
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
	"time"
)

func TestCustHandler_GetApiPayLogList(t *testing.T) {
	//time.Local =
	local, err := time.LoadLocation("Asia/Phnom_Penh")
	if err != nil {
		t.Errorf("time.LoadLocation() error = %v", err)
	}

	time.Local = local

	req := &custProto.GetApiPayLogListRequest{
		ReqStartTime: "2020/11/03 00:00:00",
		ReqEndTime:   "2020/11/03 18:00:00",
	}
	reply := &custProto.GetApiPayLogListReply{}

	if err := CustHandlerInst.GetApiPayLogList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetApiPayLogList() error = %v", err)
	}
	t.Logf("datas[%+v]", reply.Datas)

}
