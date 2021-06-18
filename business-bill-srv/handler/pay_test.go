package handler

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"context"
	"testing"
	"time"
)

func TestBusinessBillHandler_QrCodeFixedPay(t *testing.T) {
	req := &businessBillProto.QrCodeFixedPayRequest{
		OrderNo:         "2020111016115149149217",
		AccountNo:       "2cabe1a5-82f4-4c3f-b95e-e6f4b8559bc5",
		AccountType:     "4",
		PaymentPassword: "123456",
		NonStr:          "123456",
	}
	reply := &businessBillProto.QrCodeFixedPayReply{}
	if err := BusinessBillHandlerInst.QrCodeFixedPay(context.TODO(), req, reply); err != nil {
		t.Errorf("QrCodeFixedPay() error = %v", err)
		return
	}
	t.Logf("支付结果: %v", strext.ToJson(reply))
}

func TestBusinessBillHandler_QrCodeAmountPay(t *testing.T) {
	req := &businessBillProto.QrCodeAmountPayRequest{
		QrCodeId:        "5fe4762f34719450cf61639f8fcf67a0",
		AccountNo:       "2cabe1a5-82f4-4c3f-b95e-e6f4b8559bc5",
		PaymentPassword: "1231231",
		NonStr:          "1231231231",
		AccountType:     "4",
		PaymentMethod:   constants.PayMethodBalance,
		BankCardNo:      "",
	}
	reply := &businessBillProto.QrCodeAmountPayReply{}

	if err := BusinessBillHandlerInst.QrCodeAmountPay(context.TODO(), req, reply); err != nil {
		t.Errorf("QrCodeAmountPay() error = %v", err)
		return
	}

	t.Logf("支付结果: %v", strext.ToJson(reply))

}

func TestBusinessBillHandler_AppPay(t *testing.T) {
	//"app_pay_content":"{\"amount\":\"100\",\"currency_type\":\"USD\",\"order_no\":\"2020090114233176355866\",\"sign\":\"9cba3a67c4c3bdf5ac4b10d48abde03f\",\"subject\":\"APP支付\",\"timestamp\":\"1598941411\"}"
	req := &businessBillProto.AppPayRequest{
		AccountNo:     "972617f3-c85b-465b-ae3a-8491647d869d",
		AccountType:   "4",
		PaymentPwd:    "111",
		NonStr:        "111",
		AppPayContent: "{\"amount\":\"100\",\"app_name\":\"乐嘉超市\",\"currency_type\":\"USD\",\"order_no\":\"2020090415063726042730\",\"sign\":\"e90ff106b7ac6499effcbc11f9ba6ebf94218986d5d78abedaec6f56b279a042\",\"subject\":\"APP支付\",\"timestamp\":\"1599203197\"}",
	}
	reply := &businessBillProto.AppPayReply{}
	if err := BusinessBillHandlerInst.AppPay(context.TODO(), req, reply); err != nil {
		t.Errorf("AppPay() error = %v", err)
		return
	}

	t.Logf("支付结果：%v", strext.ToJson(reply))
}

func TestBusinessBillHandler_VerifySign(t *testing.T) {
	contentMap := map[string]interface{}{
		"amount":        "100",
		"app_name":      "乐嘉超市",
		"currency_type": "USD",
		"order_no":      "2020090311462931139429",
		"sign":          "87976316313b19540ca67b47a023364e76e3bb3c75fb09d2adfc7f8849d10f90",
		"subject":       "APP支付",
		"timestamp":     "1599104789",
	}
	reqSign := contentMap["sign"]
	verifySign := AppPayContentMakeSign(contentMap)
	if verifySign != reqSign {
		t.Errorf("验签失败，sign=[%v], verifySign=[%v]", reqSign, verifySign)
		return
	}
	t.Logf("验签成功")
}

func TestTime_Sql(t *testing.T) {
	d := time.Unix(0, 0)
	t.Logf("时间：%v,", d.String())
}
