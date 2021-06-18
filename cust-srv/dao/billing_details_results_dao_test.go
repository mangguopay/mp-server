package dao

import (
	"a.a/mp-server/common/constants"
	_ "a.a/mp-server/cust-srv/test"
	"fmt"
	"testing"
)

func Test_BillingDetailsResultsDao_CountServicerByDate(t *testing.T) {
	date := "2020-04-20"

	total, err := BillingDetailsResultsDaoInst.CountServicerByDate(date)

	fmt.Printf("total:%+v,err:%v\n", total, err)
}

func Test_BillingDetailsResultsDao_GetServicerNoByDate(t *testing.T) {
	startTime := "2020-04-20"

	pageSize := 2 // 每次查询数量
	page := 1     // 总页数

	list, err := BillingDetailsResultsDaoInst.GetServicerNoByDate(startTime, page, pageSize)

	fmt.Println("list-len:", len(list))
	fmt.Println("list:", list)
	fmt.Println("err:", err)
}

func Test_BillingDetailsResultsDao_GetServicerStatis(t *testing.T) {
	startTime := "2020-04-20"

	servicerNo := "99f047f6-30d0-4e76-853d-57d124c76cc0"
	data, err := BillingDetailsResultsDaoInst.GetServicerStatis(servicerNo, constants.CURRENCY_USD, startTime)
	fmt.Printf("data:%+v,err:%v\n", data, err)
}
