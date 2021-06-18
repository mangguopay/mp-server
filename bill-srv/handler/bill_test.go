package handler

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/bill-srv/test"
	"a.a/mp-server/common/constants"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"context"
	"testing"
)

func TestCheckAmountIsMaxMinTransfer(t *testing.T) {
	type args struct {
		monType string
		amount  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"01", args{"usd", "10"}, false},
		{"02", args{"usd", "10"}, false},
		{"03", args{"khr", "10"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckAmountIsMaxMinTransfer(tt.args.monType, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("CheckAmountIsMaxMinTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckAmountIsMaxMinWithdraw(t *testing.T) {
	type args struct {
		monType string
		amount  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"01", args{"usd", "10"}, false},
		{"02", args{"khr", "100"}, false},
		{"03", args{"usd", "10"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckAmountIsMaxMinWithdraw(tt.args.monType, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("CheckAmountIsMaxMinWithdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_doFees(t *testing.T) {
	var feeType int32 = 1
	amount := "100"
	rate, feeAmount, err := doFees(feeType, amount)
	if err != nil {
		t.Errorf("doFees() error = %v", err)
		return
	}
	t.Logf("费率:%v，手续费:%v", rate, feeAmount)
}

func TestBillHandler_CustOrderBillDetail(t *testing.T) {
	req := &go_micro_srv_bill.CustOrderBillDetailRequest{
		OrderNo:   "2020102016392042118773",
		OrderType: constants.VaReason_BusinessTransferToBusiness,
		AccountNo: "972617f3-c85b-465b-ae3a-8491647d869d",
		LogNo:     "2020102016392042418247",
		//CurrencyType: "USD",
	}
	reply := &go_micro_srv_bill.CustOrderBillDetailReply{}

	b := &BillHandler{}
	if err := b.CustOrderBillDetail(context.TODO(), req, reply); err != nil {
		t.Errorf("CustOrderBillDetail() error = %v", err)
		return
	}

	t.Logf("查询结果：%v", strext.ToJson(reply))
}
