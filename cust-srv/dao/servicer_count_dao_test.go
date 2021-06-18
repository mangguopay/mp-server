package dao

import (
	"a.a/mp-server/common/constants"
	_ "a.a/mp-server/cust-srv/test"
	"fmt"
	"testing"
)

func Test_ServicerCountDao_UpdateCountData(t *testing.T) {
	var data ServicerCheckListStatis
	data.ServicerNo = "e357b129-d1a7-421c-9bf3-009e3834b94d"
	data.CurrencyType = constants.CURRENCY_USD
	data.Dates = "2020-04-13"
	data.InNum = 10
	data.InAmount = 100
	data.OutNum = 20
	data.OutAmount = 200
	data.ProfitNum = 30
	data.ProfitAmount = 300
	data.RechargeNum = 40
	data.RechargeAmount = 400
	data.WithdrawNum = 50
	data.WithdrawAmount = 500

	err := ServicerCountDaoInst.UpdateCountData(data)

	fmt.Println("err:", err)
}
