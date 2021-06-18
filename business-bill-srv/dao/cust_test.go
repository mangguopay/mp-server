package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	"testing"
)

func TestCustDao_QueryCustInfo(t *testing.T) {
	accountNo := "972617f3-c85b-465b-ae3a-8491647d869d"
	got, err := CustDaoInst.QueryCustInfo("", accountNo)
	if err != nil {
		t.Errorf("QueryCustInfo() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(got))
}
