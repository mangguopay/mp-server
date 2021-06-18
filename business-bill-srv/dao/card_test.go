package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/business-bill-srv/test"
	"testing"
)

func TestCard_GetCardBaseInfo(t *testing.T) {
	accountNo := "972617f3-c85b-465b-ae3a-8491647d869d"
	bankCardNo := "b505d10c-0455-4b1d-8923-63516c78d6e1"
	info, err := CardDao.GetCardBaseInfo(accountNo, bankCardNo)
	if err != nil {
		t.Errorf("GetCardBaseInfo() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(info))
}
