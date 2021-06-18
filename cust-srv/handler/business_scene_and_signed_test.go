package handler

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_GetSceneSignedList(t *testing.T) {
	req := &custProto.GetSceneSignedListRequest{
		SceneName: "当面付",
		Lang:      constants.LangEnUS,
	}
	reply := &custProto.GetSceneSignedListReply{}
	if err := CustHandlerInst.GetSceneSignedList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetSceneSignedList() error = %v", err)
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}
