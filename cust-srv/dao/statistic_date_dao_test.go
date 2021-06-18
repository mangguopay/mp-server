package dao

import (
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestDateDao_GetStatisticData(t *testing.T) {
	startDate := "2020-05-01"
	endDate := "2020-05-07"

	list, err := StatisticDateDaoInst.GetStatisticData(startDate, endDate)

	if err != nil {
		t.Errorf("GetStatisticData-err:%v", err)
		return
	}

	t.Logf("GetStatisticData-list:%v", list)
}

func TestStatisticDateDao_GetStatisticDataList(t *testing.T) {
	req := &go_micro_srv_cust.GetStatisticDateListRequest{
		StartDate: "2020-05-01",
		EndDate:   "2020-05-07",
		//PageSize:  5,
		//Page:      3,
	}

	list, total, err := StatisticDateDaoInst.GetStatisticDataList(req)
	if err != nil {
		t.Errorf("GetStatisticDataList-err:%v", err)
		return
	}
	t.Logf("GetStatisticDataList-result,total:%d,list:%v", total, list)
}
