package dao

import (
	"context"
	"testing"

	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"

	_ "a.a/mp-server/business-bill-srv/test"
)

func TestBusinessBillSettleDao_InsertTx(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Errorf("获取数据库连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		t.Errorf("开启事物失败-err:%v", errTx)
		return
	}

	data := new(BusinessBillSettleDao)
	data.SettleId = strext.GetDailyId()
	data.BusinessNo = "d950f705-688a-4b35-b773-de4497fa7602"
	data.StartTime = "2020-07-20 00:00:00"
	data.EndTime = "2020-07-20 23:59:59"
	data.TotalAmount = 8100
	data.TotalRealAmount = 8019
	data.TotalFees = 81
	data.TotalOrder = 3

	err := BusinessBillSettDaoInst.InsertTx(tx, data)
	if err != nil {
		t.Errorf("BusinessBillSettDaoInst-InsertTx-err:%v, data:%+v", err, data)
		return
	}

	tx.Commit()

	t.Logf("BusinessBillSettDaoInst-InsertTx-ok,settleId:%v", data.SettleId)
}
