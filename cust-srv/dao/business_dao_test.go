package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestBusinessDao_GetBusinessProfit(t *testing.T) {
	whereList := []*model.WhereSqlCond{
		{Key: "acc.account", Val: "h13298690108@163.com", EqType: "like"},
	}

	opType := constants.VaOpType_Add
	reason := []string{constants.VaReason_Business_Settle, constants.VaReason_BusinessTransferToBusiness}
	ret, err := BusinessDaoInst.GetBusinessProfit(whereList, 0, 10, opType, reason)
	if err != nil {
		t.Errorf("GetBusinessProfit() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(ret))
}
