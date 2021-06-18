package handler

import (
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/global"
	custProto "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"testing"
	"time"
)

func TestCustHandler_GetRefundBills(t *testing.T) {
	req := &custProto.GetRefundBillsRequest{
		Page:         "1",
		PageSize:     "10",
		StartTime:    "",
		EndTime:      "",
		BusinessNo:   "53be2111-b1cb-4041-8143-d4a2ccf7d995",
		CurrencyType: "USD",
		RefundNo:     "",
		TransOrderNo: "",
		OrderStatus:  "",
	}
	reply := &custProto.GetRefundBillsReply{}
	if err := CustHandlerInst.GetRefundBills(context.TODO(), req, reply); err != nil {
		t.Errorf("GetRefundBills() error = %v", err)
		return
	}

	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetRefundDetail(t *testing.T) {
	req := &custProto.GetRefundDetailRequest{
		RefundNo: "2020091517302956088025",
	}
	reply := &custProto.GetRefundDetailReply{}
	if err := CustHandlerInst.GetRefundDetail(context.TODO(), req, reply); err != nil {
		t.Errorf("GetRefundDetail() error = %v", err)
		return
	}

	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetBusinessBills(t *testing.T) {
	req := &custProto.GetBusinessBillsRequest{
		Page:       1,
		PageSize:   10,
		BusinessNo: "8276897e-1ee3-471a-8563-f9a936678946",
		//BusinessName: "hex",
		IsSettled:    "1",
		OrderStatus:  "",
		CurrencyType: "USD",
		//ChannelType:  constants.ChannelTypeOut,
	}
	reply := &custProto.GetBusinessBillsReply{}
	if err := CustHandlerInst.GetBusinessBills(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessBills() error = %v", err)
		return
	}

	t.Logf("结果: %v", strext.ToJson(reply))
}

func TestCustHandler_GetBusinessBillDetail(t *testing.T) {
	i := time.Now().Unix()
	ss_time.Unixtime2Time(strext.ToString(i), global.Tz).Format(ss_time.DateFormat)

	req := &custProto.GetBusinessBillDetailRequest{
		OrderNo: "2020102115193184151501",
	}
	reply := &custProto.GetBusinessBillDetailReply{}
	if err := CustHandlerInst.GetBusinessBillDetail(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessBillDetail() error = %v", err)
		return
	}
	t.Logf("订单详情：%v", strext.ToJson(reply))
}

func TestCustHandler_CreateBillFile(t *testing.T) {

	req := &custProto.CreateBillFileRequest{
		Page:         1,
		PageSize:     1000,
		Uid:          "555d2d86-fef4-42d9-b2f0-a6adb8c3f325",
		IdenNo:       "558c1ae3-e0ca-4af1-96ae-9d5625795d2b",
		IsSettled:    "1",
		CurrencyType: "USD",
		StartTime:    "2020-08-01 00:00:00",
		EndTime:      "2020-09-03 17:00:00",
	}
	reply := &custProto.CreateBillFileReply{}

	if err := CustHandlerInst.CreateBillFile(context.TODO(), req, reply); err != nil {
		t.Errorf("CreateBillFile() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))

}

func TestCustHandler_CreateRefundFile(t *testing.T) {
	req := &custProto.CreateRefundFileRequest{
		Page:         1,
		PageSize:     1000,
		StartTime:    "2020-08-01 00:00:00",
		EndTime:      "2020-08-28 17:00:00",
		BusinessNo:   "558c1ae3-e0ca-4af1-96ae-9d5625795d2b",
		AccountNo:    "555d2d86-fef4-42d9-b2f0-a6adb8c3f325",
		CurrencyType: "USD",
	}
	reply := &custProto.CreateRefundFileReply{}

	if err := CustHandlerInst.CreateRefundFile(context.TODO(), req, reply); err != nil {
		t.Errorf("CreateRefundFile() error = %v", err)
		return
	}

	t.Logf("结果: %v", strext.ToJson(reply))

}

func TestCustHandler_GetChannelBills(t *testing.T) {
	req := &custProto.GetChannelBillsRequest{
		Page:      1,
		PageSize:  10,
		StartTime: "",
		EndTime:   "",
		//Account:      "h13298690108@163.com",
		OrderStatus:  "",
		SceneNo:      "",
		CurrencyType: "",
		IsSettled:    "1",
	}
	reply := &custProto.GetChannelBillsReply{}
	if err := CustHandlerInst.GetChannelBills(context.TODO(), req, reply); err != nil {
		t.Errorf("GetChannelBills() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetChannelRefundBills(t *testing.T) {
	req := &custProto.GetChannelRefundBillsRequest{}
	reply := &custProto.GetChannelRefundBillsReply{}
	if err := CustHandlerInst.GetChannelRefundBills(context.TODO(), req, reply); err != nil {
		t.Errorf("GetChannelRefundBills() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}
