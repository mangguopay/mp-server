package handler

import (
	"a.a/cu/strext"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"context"
	"testing"
)

func TestBusinessBillHandler_ApiPayRefundCopy(t *testing.T) {
	req := &businessBillProto.ApiPayRefundRequest{
		AppId:        "2020090417003775361070",
		OrderNo:      "2020102811543245326039",
		OutOrderNo:   "",
		RefundAmount: "10000",
		OutRefundNo:  strext.GetDailyId(),
		RefundReason: "新结构测试",
		Lang:         "",
	}
	reply := &businessBillProto.ApiPayRefundReply{}

	if err := BusinessBillHandlerInst.ApiPayRefund(context.TODO(), req, reply); err != nil {
		t.Errorf("ApiPayRefundCopy() error = %v", err)
		return
	}

	t.Logf("退款结果：%v", reply)

}

func TestBusinessBillHandler_BusinessBillRefundCopy(t *testing.T) {
	req := &businessBillProto.BusinessBillRefundRequest{
		BusinessNo:    "558c1ae3-e0ca-4af1-96ae-9d5625795d2b",
		BusinessAccNo: "555d2d86-fef4-42d9-b2f0-a6adb8c3f325",
		AccountType:   "8",
		OrderNo:       "2020102814033604397776",
		RefundAmount:  "1000",
		PaymentPwd:    "",
		NonStr:        "",
		RefundReason:  "商家中心-新结构测试",
		Lang:          "",
	}
	reply := &businessBillProto.BusinessBillRefundReply{}

	if err := BusinessBillHandlerInst.BusinessBillRefund(context.TODO(), req, reply); err != nil {
		t.Errorf("BusinessBillRefundCopy() error = %v", err)
		return
	}
	t.Logf("退款结果：%v", strext.ToJson(reply))
}
