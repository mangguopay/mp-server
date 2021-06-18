package dao

import (
	"testing"

	_ "a.a/mp-server/bill-srv/test"
)

func TestBillingDetailsResultsDao_GetRealAmountByOutOrderLogNo(t *testing.T) {
	logNo := "20200713154210351110231"
	realAmount, err := BillingDetailsResultsDaoInstance.GetRealAmountByOutOrderLogNo(logNo)
	if err != nil {
		t.Errorf("GetRealAmountByLogNo-err:%v", err)
		return
	}

	t.Logf("realAmount:%v", realAmount)
}
