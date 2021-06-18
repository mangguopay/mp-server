package handler

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"context"
	"testing"
)

//下单
func TestBusinessBillHandler_Prepay(t *testing.T) {
	req := &businessBillProto.PrepayRequest{
		Amount:       "100",
		Remark:       "接口测试商品",
		NotifyUrl:    "http://127.0.0.1:9998",
		AppId:        "2020090417003775361070",
		OutOrderNo:   strext.GetDailyId(),
		CurrencyType: "USD",
		Subject:      "支付",
		TradeType:    constants.TradeTypeModernpayFaceToFace, //TradeTypeModernpayAPP
		//PaymentCode:  "2020081715223098988137",
	}
	reply := &businessBillProto.PrepayReply{}

	if err := BusinessBillHandlerInst.Prepay(context.TODO(), req, reply); err != nil {
		t.Errorf("Prepay() error = %v", err)
	}

	t.Logf("下单结果: %v", strext.ToJson(reply))
}

func TestBusinessBillHandler_QrCodeFixedPrePay(t *testing.T) {
	req := &businessBillProto.QrCodeFixedPrePayRequest{
		QrCodeId:     "fe3b1753b0d7698bc19c5eeadc07340e",
		AccountNo:    "972617f3-c85b-465b-ae3a-8491647d869d",
		AccountType:  "4",
		Amount:       "10000",
		CurrencyType: "USD",
		Remark:       "广州酒家月饼",
		Subject:      "月饼23",
	}
	reply := &businessBillProto.QrCodeFixedPrePayReply{}
	if err := BusinessBillHandlerInst.QrCodeFixedPrePay(context.TODO(), req, reply); err != nil {
		t.Errorf("QrCodeFixedPrePay() error = %v", err)
		return
	}

	t.Logf("下单结果：%v", strext.ToJson(reply))
}

func TestBusinessBillHandler_PersonalBusinessPrepay1(t *testing.T) {
	req := &businessBillProto.PersonalBusinessPrepayRequest{
		AccountNo:    "972617f3-c85b-465b-ae3a-8491647d869d",
		Amount:       "200",
		CurrencyType: "USD",
		Subject:      "other",
		Remark:       "",
		Lang:         "",
	}
	reply := &businessBillProto.PersonalBusinessPrepayReply{}

	if err := BusinessBillHandlerInst.PersonalBusinessPrepay(context.TODO(), req, reply); err != nil {
		t.Errorf("PersonalBusinessPrepay() error = %v", err)
		return
	}
	//"order_no":"2020111015435664215982","qr_code_id":"mp://pay/bizpay?qr=9ba93b6f4f70ccfb41ee49e0e9ddd696"
	t.Logf("下单结果：%v", strext.ToJson(reply))
}

func TestBusinessBillHandler_PersonalBusinessCodeFixedPrePay(t *testing.T) {
	req := &businessBillProto.PersonalBusinessCodeFixedPrePayRequest{
		QrCodeId:     "Fmp5MGDxPdFRufAL0CuUcS7NsypJmQWJ",
		AccountNo:    "2cabe1a5-82f4-4c3f-b95e-e6f4b8559bc5",
		AccountType:  "4",
		Amount:       "100",
		CurrencyType: "USD",
		Remark:       "测试",
	}
	reply := &businessBillProto.PersonalBusinessCodeFixedPrePayReply{}
	if err := BusinessBillHandlerInst.PersonalBusinessCodeFixedPrePay(context.TODO(), req, reply); err != nil {
		t.Errorf("PersonalBusinessCodeFixedPrePay() error = %v", err)
		return
	}
}
