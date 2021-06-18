package dao

import (
	"a.a/cu/ss_time"
	_ "a.a/mp-server/bill-srv/test"
	"a.a/mp-server/common/global"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_ServiceDao_GetServicerCheckListTotal(t *testing.T) {
	startTime := ""
	endTime := ""

	servicerNo := "ea3091a5-7141-43b0-9837-dd38e61b40ca"

	returnTotals, returnErr := ServiceDaoInst.GetServicerCheckListTotal(startTime, endTime, servicerNo)

	fmt.Printf("returnTotals:%v, returnErr:%v \n", returnTotals, returnErr)
}

func Test_ServiceDao_GetServicerCheckList(t *testing.T) {
	startTime := "2020-04-10"
	endTime := "2020-05-10"
	servicerNo := "60f15170-c1db-41b0-bb3d-14185ab43d28"
	page := int32(1)
	pageSize := int32(10)

	datas, returnTotals, returnErr := ServiceDaoInst.GetServicerCheckList(startTime, endTime, servicerNo, page, pageSize)

	aa, _ := json.Marshal(datas)

	fmt.Println("json-datas--->", string(aa))

	fmt.Printf("datas:%+v, returnTotals:%v, returnErr:%v \n", datas, returnTotals, returnErr)
}

func Test_ServiceDao_GetServicerCheckListByDate(t *testing.T) {
	date := "2020-04-11"
	servicerNo := "ea3091a5-7141-43b0-9837-dd38e61b40ca"

	datas, returnErr := ServiceDaoInst.GetServicerCheckListByDate(date, servicerNo)

	aa, _ := json.Marshal(datas)

	fmt.Println("json-datas-->", string(aa))

	fmt.Printf("datas:%+v,  returnErr:%v \n", datas, returnErr)
}

func TestCreateTime(t *testing.T) {
	fmt.Println(ss_time.PostgresTimeToTime(time.Now().String(), global.Tz))
	//timeT, err := IncomeOrderDaoInst.QueryCreateTime("99f047f6-30d0-4e76-853d-57d124c76cc0")
	//if err != nil {
	//	ss_log.Error("err: %s", err.Error())
	//	return
	//}
	//fmt.Println(timeT, time.Now())
	fmt.Println(time.Now(), time.Now().Add(-5))
}
