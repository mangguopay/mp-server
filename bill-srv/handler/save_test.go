package handler

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/bill-srv/test"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	"context"
	"testing"
)

func TestBillHandler_SaveMoney(t *testing.T) {
	req := &go_micro_srv_bill.SaveMoneyRequest{
		RecvPhone:       "13298690108",                          // 收款人手机号
		SendPhone:       "077778888",                            // 存款人手机号
		Amount:          "10000",                                // 金额
		MoneyType:       "usd",                                  // 币种
		AccountType:     "3",                                    // 操作
		Password:        "2b6c6b43d04c2b7d23e707cc32306a19",     // 支付密码
		OpAccNo:         "e84eefa8-e51e-41b5-99e6-2e60ef675618", // 操作员账号
		NonStr:          "gexZO7ngaBgaJOtE",
		SaveCountryCode: "855",
		RecvCountryCode: "855",
		AccountUid:      "230f2034-f526-4542-b460-853046d14fb9",
		Lang:            "zh_CN",
		Ip:              "10.41.6.132",
	}
	reply := &go_micro_srv_bill.SaveMoneyReply{}

	if err := BillHandlerInst.SaveMoney(context.TODO(), req, reply); err != nil {
		t.Errorf("SaveMoney() error = %v", err)
	}

}

func TestBillHandler_SaveMoneyDetail(t *testing.T) {
	req := go_micro_srv_bill.SaveMoneyDetailRequest{
		OrderNo: "2019123012400389755784",
	}
	reply := go_micro_srv_bill.SaveMoneyDetailReply{}
	if err := BillHandlerInst.SaveMoneyDetail(context.TODO(), &req, &reply); err != nil {
		t.Errorf("SaveMoneyDetail() error = %v", err)
	}

	t.Logf("存款详情：%v", strext.ToJson(reply.Data))

}

func TestBillHandler_QuerySaveReceipt(t *testing.T) {
	req := go_micro_srv_bill.QuerySaveReceiptRequest{
		OrderNo: "2019123012400389755784",
	}
	reply := go_micro_srv_bill.QuerySaveReceiptReply{}

	if err := BillHandlerInst.QuerySaveReceipt(context.TODO(), &req, &reply); err != nil {
		t.Errorf("QuerySaveReceipt() error = %v", err)
	}

	t.Logf("小票：%v", strext.ToJson(reply.Data))
}
