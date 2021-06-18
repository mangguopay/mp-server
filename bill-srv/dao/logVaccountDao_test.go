package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/bill-srv/test"
	"a.a/mp-server/common/constants"
	"testing"
)

func TestLogVaccountDao_GetLogVAccountByBizLogNo(t *testing.T) {
	accountNo := ""
	bizLogNo := "2020073011115556133322"
	reason := constants.VaReason_TRANSFER
	got, err := LogVaccountDaoInst.GetLogVAccountByBizLogNo(accountNo, bizLogNo, reason)
	if err != nil {
		t.Errorf("GetLogVAccountByBizLogNo() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(got))
}

func TestLogVaccountDao_GetLogVAccountJoinWriteOff(t *testing.T) {
	accountNo := "b4203420-be32-4ee2-acbe-e9d444fac58a"
	LogNo := "2020101510334486150379"
	reason := constants.VaReason_PlatformFreeze
	ret, err := LogVaccountDaoInst.GetLogVAccountJoinWriteOff(accountNo, LogNo, reason)
	if err != nil {
		t.Errorf("GetLogVAccountJoinWriteOff() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(ret))
}
