package handler

import (
	"a.a/cu/strext"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_GetHeadquartersProfitList(t *testing.T) {
	req := &go_micro_srv_cust.GetHeadquartersProfitListRequest{}
	reply := &go_micro_srv_cust.GetHeadquartersProfitListReply{}
	if err := CustHandlerInst.GetHeadquartersProfitList(context.TODO(), req, reply); err != nil {
		t.Errorf("GetHeadquartersProfitList() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(reply))
}
