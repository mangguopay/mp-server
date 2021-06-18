package handler

import (
	"a.a/cu/strext"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"context"
	"testing"
)

func TestBusinessBillHandler_GetPersonalBusinessInfo(t *testing.T) {
	req := &businessBillProto.GetPersonalBusinessInfoRequest{
		FixedCode: "Fmp5MGDxPdFRufAL0CuUcS7NsypJmQWJ",
		Lang:      "",
	}
	reply := &businessBillProto.GetPersonalBusinessInfoReply{}
	if err := BusinessBillHandlerInst.GetPersonalBusinessInfo(context.TODO(), req, reply); err != nil {
		t.Errorf("GetPersonalBusinessInfo() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}
