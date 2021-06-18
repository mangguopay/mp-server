package dao

import (
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestUserWithdrawDao_GetStatisticData(t *testing.T) {
	startDate := "2020-05-20"
	endDate := "2020-05-22"
	currencyType := "usd"

	list, err := StatisticUserWithdrawDaoInst.GetStatisticData(startDate, endDate, currencyType)

	if err != nil {
		t.Errorf("GetStatisticData-err:%v", err)
		return
	}

	t.Logf("GetStatisticData-list:%v", list)
}

func TestStatisticUserWithdrawDao_GetStatisticDataList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserWithdrawListRequest{
		StartDate:    "2020-05-20",
		EndDate:      "2020-05-22",
		CurrencyType: "usd",
		WithdrawType: "1",
		PageSize:     5,
		Page:         3,
	}

	list, total, err := StatisticUserWithdrawDaoInst.GetStatisticDataList(req)
	if err != nil {
		t.Errorf("GetStatisticDataList-err:%v", err)
		return
	}
	t.Logf("GetStatisticDataList-result,total:%d,list:%v", total, list)
}
