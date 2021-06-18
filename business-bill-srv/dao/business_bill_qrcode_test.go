package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"context"
	"testing"

	_ "a.a/mp-server/business-bill-srv/test"
)

//查询订单号
func TestBusinessBillQrCode_QueryOrderNoByQrCodeId(t *testing.T) {
	var qrCodeId = "a3316cdd1028d95f7b934ca4508b67a4"
	orderNo, err := BusinessBillQrCodeInst.QueryOrderNoByQrCodeId(qrCodeId)
	if err != nil {
		t.Errorf("查询订单号失败, err:%v", err)
	}
	t.Logf("qrCodeId:%v, orderNo: %v", qrCodeId, orderNo)
}

//插入记录
func TestBusinessBillQrCode_InsertOrderQrCode(t *testing.T) {
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

	var orderNo = ""
	var payQrCodeId = ""
	if err := BusinessBillQrCodeInst.InsertOrderQrCode(tx, orderNo, payQrCodeId); err != nil {
		t.Errorf("插入记录失败, err:%v", err)
	}
	t.Logf("插入成功")
}
