package dao

import (
	"a.a/cu/strext"
	"testing"

	_ "a.a/mp-server/business-bill-srv/test"
)

func TestBusinessCheckingDao_Insert(t *testing.T) {
	data := BusinessChecking{
		CheckingId:         strext.GetDailyId(),
		BusinessNo:         "d950f705-688a-4b35-b773-de4497fa7602",
		BusinessAccountNo:  "789f5dbb-d6dc-478e-83a8-af10875ac0c6",
		CurrencyType:       "USD",
		BusinessBillAmount: 0,
		AccountBalance:     0,
		SettledId:          strext.GetDailyId(),
	}
	if err := BusinessCheckingDaoInst.Insert(data); err != nil {
		t.Errorf("Insert() error = %v", err)
	}
	t.Logf("插入数据成功")
}
