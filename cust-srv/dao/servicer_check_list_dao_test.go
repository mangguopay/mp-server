package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_ServicerCheckListDao_InsertByCurrency(t *testing.T) {
	var data ServicerCheckListStatis
	data.ServicerNo = "22222222-7141-43b0-9837-dd38e61b40ca"
	data.CurrencyType = constants.CURRENCY_USD
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
	data.Dates = "2020-05-11"

	err := ServicerCheckListDaoInst.InsertByCurrency(data)

	fmt.Println("err:", err)
}

func Test_ServicerCheckListDao_GetCheckListStatis(t *testing.T) {
	date := "2020-04-13"
	pageSize := 2 // 每次查询数量

	list, err := ServicerCheckListDaoInst.GetCheckListStatis(date, pageSize)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Printf("len:%+v\n", len(list))
	fmt.Printf("list:%+v\n", list)
}

func Test_ServicerCheckListDao_UpdateIsCountedTx(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//更新服务商的统计	servicer_count
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*120)
	defer cancel()

	tx, errTx := dbHandler.BeginTx(ctx, nil)
	if errTx != nil {

	}

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

	err := ServicerCheckListDaoInst.UpdateIsCountedTx(tx, data)

	tx.Commit()

	fmt.Printf("err:%v\n", err)
}
