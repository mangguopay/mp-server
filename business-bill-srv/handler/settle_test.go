package handler

import (
	"a.a/cu/strext"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
	"context"
	"testing"
)

func TestBusinessBillHandler_SingleOrderSettle(t *testing.T) {
	orderNo := ""
	settleId, err := BusinessBillHandlerInst.SingleOrderSettle(orderNo)
	if err != ss_err.Success {
		t.Errorf("SingleOrderSettle() error = %v", err)
		return
	}
	t.Logf("结算：%v", settleId)
}

func TestBusinessBillHandler_D(t *testing.T) {

	t.Logf("结算：%v", strext.GetDailyId())
}

func TestBusinessBillHandler_ManualSettle(t *testing.T) {
	orderNo := []string{
		"2020091415215044583020",
		"2020091415212338653817",
		"2020091117461516440244",
		"2020091117141846522927",
		"2020091116532836932466",
	}
	t.Logf("订单：%v", strext.ToJson(orderNo))
	req := &businessBillProto.ManualSettleRequest{
		OrderNos: orderNo,
	}
	reply := &businessBillProto.ManualSettleReply{}
	if err := BusinessBillHandlerInst.ManualSettle(context.TODO(), req, reply); err != nil {
		t.Errorf("ManualSettle() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}
