package dao

import (
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestUserTransferDao_GetStatisticData(t *testing.T) {
	startDate := "2020-05-14"
	endDate := "2020-05-19"
	currencyType := "usd"

	list, err := StatisticUserTransferDaoInst.GetStatisticData(startDate, endDate, currencyType)

	if err != nil {
		t.Errorf("GetStatisticData-err:%v", err)
		return
	}

	t.Logf("GetStatisticData-list:%v", list)
}

func TestStatisticUserTransferDao_GetStatisticDataList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserTransferListRequest{
		StartDate:    "2020-05-14",
		EndDate:      "2020-05-19",
		CurrencyType: "usd",
		//PageSize:     5,
		//Page:         3,
	}

	list, total, err := StatisticUserTransferDaoInst.GetStatisticDataList(req)
	if err != nil {
		t.Errorf("GetStatisticDataList-err:%v", err)
		return
	}
	t.Logf("GetStatisticDataList-result,total:%d,list:%v", total, list)
}
