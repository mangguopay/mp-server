package handler

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"context"
	"testing"
)

func TestBusinessBillHandler_Transfers(t *testing.T) {
	req := &businessBillProto.EnterpriseTransferRequest{
		AppId:            "2020090417003775361070",
		Amount:           "1000",
		CurrencyType:     "USD",
		PayeePhone:       "789456123",
		PayeeCountryCode: "855",
		Remark:           "测试",
		Lang:             "",
		OutTransferNo:    "interface001",
	}
	reply := &businessBillProto.EnterpriseTransferReply{}
	if err := BusinessBillHandlerInst.EnterpriseTransfer(context.TODO(), req, reply); err != nil {
		t.Errorf("Transfers() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}
