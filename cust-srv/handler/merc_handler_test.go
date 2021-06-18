package handler

import (
	"a.a/cu/strext"
	custProto "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"testing"
)

func TestCustHandler_GetBusinessMessagesUnRead(t *testing.T) {
	req := &custProto.GetBusinessMessagesUnReadRequest{
		AccountNo:   "4c95eaa8-ee0c-4e21-aa96-2fa9b79a450e",
		AccountType: "8",
	}
	reply := &custProto.GetBusinessMessagesUnReadReply{}
	if err := CustHandlerInst.GetBusinessMessagesUnRead(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessMessagesUnRead() error = %v", err)
		return
	}

	t.Logf("未读消息数量：%v", reply.Number)
}

func TestCustHandler_IsEnabledScene(t *testing.T) {
	req := &custProto.IsEnabledSceneRequest{
		SceneNo:   "f4adfcd4-c490-4750-b5b9-80637dc1745c",
		IsEnabled: "0",
	}
	reply := &custProto.IsEnabledSceneReply{}
	if err := CustHandlerInst.IsEnabledScene(context.TODO(), req, reply); err != nil {
		t.Errorf("IsEnabledScene() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetBusinessAccountsProfit(t *testing.T) {
	req := &custProto.GetBusinessAccountsProfitRequest{
		Page:     0,
		PageSize: 10,
		//Account:  "h13298690108@163.com",
	}
	reply := &custProto.GetBusinessAccountsProfitReply{}
	if err := CustHandlerInst.GetBusinessAccountsProfit(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessAccountsProfit() error = %v", err)
		return
	}

	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetBusinessProfitList(t *testing.T) {
	req := &custProto.GetBusinessProfitListRequest{
		Reason:            "",
		CurrencyType:      "usd",
		BusinessAccountNo: "972617f3-c85b-465b-ae3a-8491647d869d",
		OpType:            "1",
	}
	reply := &custProto.GetBusinessProfitListReply{}

	if err := CustHandlerInst.GetBusinessProfitList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessProfitList() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetBusinessSceneDetail(t *testing.T) {
	req := &custProto.GetBusinessSceneDetailRequest{
		SceneNo: "246b72ad-2904-42b6-83b8-274dbf3a4927",
	}
	reply := &custProto.GetBusinessSceneDetailReply{}
	if err := CustHandlerInst.GetBusinessSceneDetail(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessSceneDetail() error = %v", err)
		return
	}
}

func TestCustHandler_GetBusinessSceneList(t *testing.T) {
	req := &custProto.GetBusinessSceneListRequest{
		Page:      "",
		PageSize:  "",
		StartTime: "",
		EndTime:   "",
		SceneName: "",
		IsDelete:  "",
		Lang:      "",
	}
	reply := &custProto.GetBusinessSceneListReply{}

	if err := CustHandlerInst.GetBusinessSceneList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetBusinessSceneList() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}
