package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestServicerTerminal_GetTerminalByNumber(t *testing.T) {
	terminalNumber := "98212004010003"
	useStatus := constants.Status_Enable
	list, err := ServicerTerminalDao.GetTerminalByNumber(terminalNumber, useStatus)
	if err != nil {
		t.Errorf("GetTerminalByNumber() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(list))
}
