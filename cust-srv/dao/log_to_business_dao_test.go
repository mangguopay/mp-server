package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestLogToBusinessDao_GetToBusinessListDetail(t *testing.T) {
	logNo := "2020081216165503533361"
	got, err := LogToBusinessDaoInst.GetToBusinessDetail(logNo)
	if err != nil {
		t.Errorf("GetToBusinessListDetail() error = %v", err)
		return
	}

	t.Logf("详情：%v", strext.ToJson(got))

}
