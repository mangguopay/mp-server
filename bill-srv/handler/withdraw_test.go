package handler

import (
	"a.a/cu/strext"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"context"
	"testing"
)

func TestBillHandler_Withdrawal(t *testing.T) {
	req := go_micro_srv_bill.WithdrawalRequest{}
	reply := go_micro_srv_bill.WithdrawalReply{}

	if err := BillHandlerInst.Withdrawal(context.TODO(), &req, &reply); err != nil {
		t.Errorf("Withdrawal() error = %v", err)
	}
	t.Logf("取款结果：%v", strext.ToJson(reply))
}
