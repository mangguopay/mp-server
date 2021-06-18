package dao

import "testing"

func TestServiceDao_GetServiceNoByAccount(t *testing.T) {
	account := "085513888888888"
	servicerNo, err := ServiceDaoInst.GetServiceNoByAccount(account)
	if err != nil {
		t.Errorf("GetServiceNoByAccount() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", servicerNo)
}
