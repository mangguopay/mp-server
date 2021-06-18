package dao

import (
	_ "a.a/mp-server/business-bill-srv/test"
	"context"
	"testing"
	"time"

	"a.a/cu/ss_time"

	"a.a/cu/db"

	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
)

func TestBusinessBillDao_OutOrderNoExist(t *testing.T) {
	businessNo := "d950f705-688a-4b35-b773-de4497fa7602"
	outOrderNo := "outOrderNo555555555555"

	exist, err := BusinessBillDaoInst.OutOrderNoExist(businessNo, outOrderNo)

	if err != nil {
		t.Errorf("OutOrderNoExist-err:%v, businessNo:%v, outOrderNo:%v", err, businessNo, outOrderNo)
		return
	}

	t.Logf("exist:%v", exist)
}

func TestBusinessBillDao_InsertOrder(t *testing.T) {
	order := BusinessBillDao{
		OrderNo:           strext.GetDailyId(),
		Fee:               "10",
		Amount:            "100",
		RealAmount:        "90",
		OrderStatus:       constants.BusinessOrderStatusPending,
		Remark:            "这是remark",
		NotifyUrl:         "http://127.0.0.1/notify",
		RreturnUrl:        "http://127.0.0.1/jump_back",
		OutOrderNo:        "13333333333333333333333",
		Rate:              "10",
		BusinessNo:        "d950f705-688a-4b35-b773-de4497fa7602",
		BusinessAccountNo: "789f5dbb-d6dc-478e-83a8-af10875ac0c6",
		AppId:             "2020063014584944251674",
		CurrencyType:      "USD",
		Subject:           "商品001",
		SceneNo:           "f163be38-3dab-4ee0-a05b-32deafbea51e",
		ExpireTime:        time.Now().Unix(),
		AccountNo:         "972617f3-c85b-465b-ae3a-8491647d869d",
	}

	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Error("获取数据连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Errorf("开启事物失败,err:%v", err)
		return
	}

	err = BusinessBillDaoInst.InsertOrderTx(tx, order)
	if err != nil {
		tx.Rollback()
		t.Errorf("OutOrderNoExist-err:%v, order:%v", err, strext.ToJson(order))
		return
	}

	tx.Commit()
	t.Logf("插入订单成功OrderNo:%v", order.OrderNo)
}

func TestBusinessBillDao_GetSettleData(t *testing.T) {
	startTime, _ := time.Parse(ss_time.DateTimeDashFormat, "2020-08-26 00:00:00")
	endTime, _ := time.Parse(ss_time.DateTimeDashFormat, "2020-08-26 23:59:59")

	data, err := BusinessBillDaoInst.GetSettleData(startTime.Unix(), endTime.Unix())
	if err != nil {
		t.Errorf("GetSettleData-err:%v", err)
		return
	}

	t.Logf("GetSettleData-ok, data:%v", strext.ToJson(data))
}

func TestBusinessBillDao_GetOrderInfoByOrderNo(t *testing.T) {
	orderNo := "2020090316163490748190"
	orderInfo, err := BusinessBillDaoInst.GetOrderInfoByOrderNo(orderNo)
	if err != nil {
		t.Errorf("GetOrderInfoByOrderNo() error = %v", err)
		return
	}

	t.Logf("%v订单信息: %v", orderNo, strext.ToJson(orderInfo))
}

func TestBusinessBillDao_GetCheckingData(t *testing.T) {
	isSettled := false //未结算
	businessNo := "d950f705-688a-4b35-b773-de4497fa7602"
	currencyType := "USD"
	data, err := BusinessBillDaoInst.GetBusinessTransData(isSettled, businessNo, currencyType)
	if err != nil {
		t.Errorf("GetCheckingData() error = %v", err)
		return
	}

	t.Logf("对账数据:%v", strext.ToJson(data))

}

func TestBusinessBillDao_UpdateOrderSettleId(t *testing.T) {
	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Error("获取数据连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Errorf("开启事物失败,err:%v", err)
		return
	}

	data := new(UpdateOrderSettleId)
	data.SettleId = "2020111207000021854264"
	data.BusinessNo = "8276897e-1ee3-471a-8563-f9a936678946"
	data.AppId = ""
	data.StartTime = time.Now().AddDate(0, 0, -1).Unix()
	data.EndTime = time.Now().Unix()
	if err := BusinessBillDaoInst.UpdateSettleIdTx(tx, data); err != nil {
		t.Errorf("UpdateOrderSettleId() error = %v", err)
		tx.Rollback()
		return
	}

	tx.Commit()

}

func TestBusinessBillDao_UpdateOrderPaid(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, err := dbHandler.Begin()
	if err != nil {
		t.Error("开启事务失败")
		return
	}
	paidData := UpdateOrderPaidData{
		OrderNo:            "2020082419100019534603",
		AccountNo:          "2cabe1a5-82f4-4c3f-b95e-e6f4b8559bc5",
		VaccountNo:         "40df1317-5388-4c5e-83b8-3730979ba92a",
		BusinessVaccountNo: "822d554e-fed2-4e76-948c-e7a26b14b2f3",
		PayTime:            "2020-08-24 18:10:13",
		Cycle:              "T+1",
		SettleDate:         time.Now().AddDate(0, 0, 1).Unix(),
	}

	if err := BusinessBillDaoInst.UpdateOrderPaid(tx, paidData); err != nil {
		tx.Rollback()
		t.Errorf("UpdateOrderPaid() error = %v", err)
		return
	}

	tx.Commit()
	t.Logf("修改订单为支付成功成功")

}

func TestBusinessBillDao_UpdateOrderOutTime(t *testing.T) {
	if err := BusinessBillDaoInst.UpdateOrderOutTime(); err != nil {
		t.Errorf("UpdateOrderOutTime() error = %v", err)
	}
	t.Logf("修改订单为超时成功")

}

func TestBusinessBillDao_OutOrderNoExist1(t *testing.T) {
	businessNo := ""
	outOrderNo := ""
	got, err := BusinessBillDaoInst.OutOrderNoExist(businessNo, outOrderNo)
	if err != nil {
		t.Errorf("OutOrderNoExist() error = %v", err)
		return
	}

	t.Logf("订单是否已存在：%v", got)
}

func TestBusinessBillDao_GetCustPendingPayOrder(t *testing.T) {
	accountNo := "f0347afa-a89d-411e-923e-d8f1ba47bd6c"
	orderNo := "2020082510571194132302"
	order, err := BusinessBillDaoInst.GetCustPendingPayOrder(accountNo, orderNo)
	if err != nil {
		t.Errorf("GetCustPendingPayOrder() error = %v", err)
		return
	}

	t.Logf("订单:%v", strext.ToJson(order))
}

func TestBusinessBillDao_AppQueryOrder(t *testing.T) {
	appId := "2020080510551083814320"
	orderNo := "2020083113424337056381"
	outOrderNo := "merchant2020083113424276549170"

	order, err := BusinessBillDaoInst.AppQueryOrder(appId, orderNo, outOrderNo)
	if err != nil {
		t.Errorf("GetCustPendingPayOrder() error = %v", err)
		return
	}

	t.Logf("订单:%v", strext.ToJson(order))
}

func TestBusinessBillDao_GetOrderInfo(t *testing.T) {
	orderNo := "2020090114233176355866"
	outOrderNo := ""
	order, err := BusinessBillDaoInst.GetOrderInfo(orderNo, outOrderNo)
	if err != nil {
		t.Errorf("GetOrderInfo() error = %v", err)
		return
	}
	t.Logf("查询结果：%v", strext.ToJson(order))
}

func TestBusinessBillDao_InsertOrder1(t *testing.T) {
	order := BusinessBillDao{
		OrderNo:      strext.GetDailyId(),
		OutOrderNo:   "",
		OrderStatus:  "1",
		Rate:         "100",
		Fee:          "1",
		Amount:       "100",
		RealAmount:   "99",
		CurrencyType: "USD",
	}

	if err := BusinessBillDaoInst.InsertOrder(order); err != nil {
		t.Errorf("InsertOrder() error = %v", err)
		return
	}
	t.Logf("插入成功")

}
