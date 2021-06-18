package dao

import "testing"

func TestAccountDao_GetAccountById(t *testing.T) {
	account, err := AccountDaoInst.GetAccountById("555d2d86-fef4-42d9-b2f0-a6adb8c3f325")
	if err != nil {
		t.Errorf("GetAccountById() error = %v", err)
		return
	}
	t.Logf("账号：%v", account)
}
