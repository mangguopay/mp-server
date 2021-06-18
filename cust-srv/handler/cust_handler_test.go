package handler

import (
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"testing"
)

func TestCustHandler_CustBills(t *testing.T) {
	req := &custProto.CustBillsRequest{
		AccountUid: "33b34e0f-e655-4a00-9e61-8ee026ddf264",
		//CurrencyType: "usd",
	}
	reply := &custProto.CustBillsReply{}

	if err := CustHandlerInst.CustBills(context.TODO(), req, reply); err != nil {
		t.Errorf("CustBills() error = %v", err)
		return
	}
	t.Logf(": %v", strext.ToJson(reply))
}
