package handler

import (
	"a.a/cu/strext"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"testing"
)

func TestCustHandler_GetStatisticUserRecharge(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserRechargeRequest{
		StartDate: "2020-06-02",
		EndDate:   "2020-06-06",
	}
	reply := &go_micro_srv_cust.GetStatisticUserRechargeReply{}

	cu := &CustHandler{}
	err := cu.GetStatisticUserRecharge(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticUserRecharge-err:%v", err)
		return
	}
	t.Logf("GetStatisticUserRecharge-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_GetStatisticUserRechargeList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserRechargeListRequest{
		StartDate: "2020-06-02",
		EndDate:   "2020-06-06",
	}
	reply := &go_micro_srv_cust.GetStatisticUserRechargeListReply{}

	cu := &CustHandler{}

	err := cu.GetStatisticUserRechargeList(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticUserRechargeList-err:%v", err)
		return
	}
	t.Logf("GetStatisticUserRechargeList-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_GetStatisticUserTransfer(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserTransferRequest{
		StartDate: "2020-05-14",
		EndDate:   "2020-05-19",
	}
	reply := &go_micro_srv_cust.GetStatisticUserTransferReply{}

	cu := &CustHandler{}
	err := cu.GetStatisticUserTransfer(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticUserTransfer-err:%v", err)
		return
	}
	t.Logf("GetStatisticUserTransfer-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_GetStatisticUserTransferList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserTransferListRequest{
		StartDate: "2020-05-14",
		EndDate:   "2020-05-19",
	}
	reply := &go_micro_srv_cust.GetStatisticUserTransferListReply{}

	cu := &CustHandler{}
	err := cu.GetStatisticUserTransferList(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticUserTransferList-err:%v", err)
		return
	}
	t.Logf("GetStatisticUserTransferList-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_GetStatisticUserExchange(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserExchangeRequest{
		StartDate: "2020-05-17",
		EndDate:   "2020-05-23",
	}
	reply := &go_micro_srv_cust.GetStatisticUserExchangeReply{}

	cu := &CustHandler{}
	err := cu.GetStatisticUserExchange(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticUserExchange-err:%v", err)
		return
	}
	t.Logf("GetStatisticUserExchange-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_GetStatisticUserExchangeList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserExchangeListRequest{
		StartDate: "2020-05-17",
		EndDate:   "2020-05-23",
	}
	reply := &go_micro_srv_cust.GetStatisticUserExchangeListReply{}

	cu := &CustHandler{}
	err := cu.GetStatisticUserExchangeList(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticUserExchangeList-err:%v", err)
		return
	}
	t.Logf("GetStatisticUserExchangeList-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_GetStatisticDate(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticDateRequest{
		StartDate: "2020-05-01",
		EndDate:   "2020-05-07",
	}
	reply := &go_micro_srv_cust.GetStatisticDateReply{}

	cu := &CustHandler{}
	err := cu.GetStatisticDate(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticDate-err:%v", err)
		return
	}
	t.Logf("GetStatisticDate-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_GetStatisticDateList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticDateListRequest{
		StartDate: "2020-05-01",
		EndDate:   "2020-05-07",
	}
	reply := &go_micro_srv_cust.GetStatisticDateListReply{}

	cu := &CustHandler{}
	err := cu.GetStatisticDateList(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticDateList-err:%v", err)
		return
	}
	t.Logf("GetStatisticDateList-reply-json:%v", strext.ToJson(reply))
}

func TestCustHandler_ReStatistic(t *testing.T) {
	req := &go_micro_srv_cust.ReStatisticRequest{
		Type:      ReStatisticTypeDate,
		StartDate: "2020-06-04",
		EndDate:   "2020-06-04",
	}
	reply := &go_micro_srv_cust.ReStatisticReply{}

	cu := &CustHandler{}
	err := cu.ReStatistic(context.TODO(), req, reply)
	if err != nil {
		t.Errorf("GetStatisticDate-err:%v", err)
		return
	}
	t.Logf("GetStatisticDate-reply-json:%v", strext.ToJson(reply))

	select {}
}
