package dao

import "testing"

func TestAccDao_GetIsActiveFromPhone(t *testing.T) {
	phone := ""
	countryCode := ""
	isActive := AccDaoInstance.GetIsActiveFromPhone(phone, countryCode)
	t.Logf("查询结果：%v", isActive)
}
