package dao

import (
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestUserExchangeDao_GetStatisticData(t *testing.T) {
	startDate := "2020-05-17"
	endDate := "2020-05-23"

	list, err := StatisticUserExchangeDaoInst.GetStatisticData(startDate, endDate)
	if err != nil {
		t.Errorf("GetStatisticData-err:%v", err)
		return
	}

	t.Logf("GetStatisticData-list:%v", list)
}

func TestStatisticUserExchangeDao_GetStatisticDataList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticUserExchangeListRequest{
		StartDate: "2020-05-17",
		EndDate:   "2020-05-23",
		//PageSize:     5,
		//Page:         3,
	}

	list, total, err := StatisticUserExchangeDaoInst.GetStatisticDataList(req)
	if err != nil {
		t.Errorf("GetStatisticDataList-err:%v", err)
		return
	}
	t.Logf("GetStatisticDataList-result,total:%d,list:%v", total, list)
}
