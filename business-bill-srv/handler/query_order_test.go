package handler

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"context"
	"regexp"
	"testing"
)

//查询交易信息
func TestBusinessBillHandler_QueryTransInfo(t *testing.T) {
	req := &businessBillProto.QueryTransInfoRequest{
		QrCodeId: "fca761585ef7bf3d08d8cf0cde1dcb6f",
	}
	resp := &businessBillProto.QueryTransInfoReply{}

	if err := BusinessBillHandlerInst.QueryTransInfo(context.TODO(), req, resp); err != nil {
		t.Errorf("QueryTransInfo() error = %v", err)
		return
	}
	t.Logf("交易信息查询结果resp: %v", strext.ToJson(resp))
}

//查询订单支付结果
func TestBusinessBillHandler_Query(t *testing.T) {

	req := &businessBillProto.QueryRequest{
		AppId:      "2020080510551083814320",
		OrderNo:    "2020081911374142042839",
		OutOrderNo: "merchant2020081911372744349170",
	}
	reply := &businessBillProto.QueryReply{}
	if err := BusinessBillHandlerInst.Query(context.TODO(), req, reply); err != nil {
		t.Errorf("Query() error = %v", err)
	}
	t.Logf("订单支付结果: %v", strext.ToJson(reply))

}

//付款码
func TestCustHandler_GetPaymentCode(t *testing.T) {
	req := &businessBillProto.GetPaymentCodeRequest{
		AccountNo: "0e8d24af-bec7-4f95-b038-c48045f51abf",
		//PaymentCode: "2020081711545564724699",
	}
	reply := &businessBillProto.GetPaymentCodeReply{}
	if err := BusinessBillHandlerInst.GetPaymentCode(context.TODO(), req, reply); err != nil {
		t.Errorf("GetPaymentCode() error = %v", err)
		return
	}
	t.Logf(strext.ToJson(reply))
}

func TestBusinessBillHandler_QueryPendingPayOrder(t *testing.T) {
	req := &businessBillProto.QueryPendingPayOrderRequest{
		AccountNo:   "0e8d24af-bec7-4f95-b038-c48045f51abf",
		PaymentCode: "2020081715223098988137",
	}
	reply := &businessBillProto.QueryPendingPayOrderReply{}
	if err := BusinessBillHandlerInst.QueryPendingPayOrder(context.TODO(), req, reply); err != nil {
		t.Errorf("QueryPendingPayOrder() error =%v", err)
		return
	}

	t.Logf("订单：%v", strext.ToJson(reply))

}

func TestBusinessBillHandler_QueryOrder(t *testing.T) {
	req := &businessBillProto.QueryOrderRequest{
		OrderNo: "2020090114233176355866",
	}

	reply := &businessBillProto.QueryOrderReply{}

	if err := BusinessBillHandlerInst.QueryOrder(context.TODO(), req, reply); err != nil {
		t.Errorf("QueryOrder() error = %v", err)
	}

	t.Logf("订单：%v", strext.ToJson(reply))

}

func TestRegexp_FindAll(t *testing.T) {
	str := "declare @shp Table(id int IDENTITY PRIMARY KEY,shpid int)  insert into @shp values(537) insert into @shp values(449) insert into @shp values(1027)"
	reg := regexp.MustCompile("\\d+")
	dataSlice := reg.FindAll([]byte(str), -1)
	for _, v := range dataSlice {
		t.Logf("结果：%v", string(v))
	}
}

func TestBusinessBillHandler_QueryTransfer(t *testing.T) {
	req := &businessBillProto.QueryTransferRequest{
		AppId:         "2020090416495834598604",
		TransferNo:    "2020103010464168182487",
		OutTransferNo: "",
	}
	reply := &businessBillProto.QueryTransferReply{}
	if err := BusinessBillHandlerInst.QueryTransfer(context.TODO(), req, reply); err != nil {
		t.Errorf("QueryTransfer() error = %v", err)
		return
	}

	t.Logf("查询结果：%v", strext.ToJson(reply))
}
