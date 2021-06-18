package dao

import (
	"testing"

	_ "a.a/mp-server/cust-srv/test"
)

func TestVaccountDao_GetBusinessRecordedAmount(t *testing.T) {
	vAccountNo := "d246571a-b950-4c26-87b4-e58dc81dff56"
	startTime := "2020-08-04 00:00:00"
	endTime := "2020-08-04 23:59:59"

	got, err := LogVaccountDaoInst.GetBusinessRecordedAmount(vAccountNo, startTime, endTime)
	if err != nil {
		t.Errorf("GetBusiness() error = %v", err)
		return
	}

	t.Logf("总金额：%v", got)
}

func TestVaccountDao_GetBusinessExpenditureAmount(t *testing.T) {
	vAccountNo := "d246571a-b950-4c26-87b4-e58dc81dff56"
	startTime := "2020-08-04 00:00:00"
	endTime := "2020-08-04 23:59:59"

	got, err := LogVaccountDaoInst.GetBusinessExpenditureAmount(vAccountNo, startTime, endTime)
	if err != nil {
		t.Errorf("GetBusiness() error = %v", err)
		return
	}

	t.Logf("总金额：%v", got)
}
