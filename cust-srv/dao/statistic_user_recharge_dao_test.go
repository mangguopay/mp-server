package dao

import (
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestUserRechargeDao_GetStatisticData(t *testing.T) {
	startDate := "2020-06-02"
	endDate := "2020-06-05"
	currencyType := "usd"

	list, err := StatisticUserRechargeDaoInst.GetStatisticData(startDate, endDate, currencyType)

	if err != nil {
		t.Errorf("GetStatisticData-err:%v", err)
		return
	}

	t.Logf("GetStatisticData-list:%v", list)
}

func TestStatisticUserRechargeDao_GetStatisticDataList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserRechargeListRequest{
		StartDate:    "2020-06-02",
		EndDate:      "2020-06-05",
		CurrencyType: "usd",
		RechargeType: "1",
		//PageSize:     5,
		//Page:         3,
	}

	list, total, err := StatisticUserRechargeDaoInst.GetStatisticDataList(req)
	if err != nil {
		t.Errorf("GetStatisticDataList-err:%v", err)
		return
	}
	t.Logf("GetStatisticDataList-result,total:%d,list:%v", total, list)
}
