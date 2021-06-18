package handler

import (
	"a.a/cu/strext"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_count"
	"context"
	"testing"
)

func TestBusinessBillHandler_BusinessBillRefund(t *testing.T) {
	req := &businessBillProto.BusinessBillRefundRequest{
		BusinessNo:    "53be2111-b1cb-4041-8143-d4a2ccf7d995",
		BusinessAccNo: "e8a38dc2-7f59-4396-b00a-ec218a00d5bb",
		OrderNo:       "2020081911083983446497",
		RefundAmount:  "1000",
		//CurrencyType: "USD",
	}
	reply := &businessBillProto.BusinessBillRefundReply{}
	if err := BusinessBillHandlerInst.BusinessBillRefund(context.TODO(), req, reply); err != nil {
		t.Errorf("BusinessBillRefund() error = %v", err)
		return
	}

	t.Logf("退款结果：%v", strext.ToJson(reply))
}

func TestBusinessBillHandler_Count(t *testing.T) {
	ret := ss_count.Sub("10", "100")
	t.Logf("计算结果：%v", ret.String())
	if ret.Sign() < 0 {
		t.Logf("结果小于0，Sign()返回值：%v", ret.Sign())
	} else if ret.Sign() == 0 {
		t.Logf("结果等于0，Sign()返回值：%v", ret.Sign())
	} else {
		t.Logf("结果大于0，Sign()返回值：%v", ret.Sign())
	}

}
