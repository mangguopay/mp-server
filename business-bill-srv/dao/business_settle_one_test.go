package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"testing"
)

func TestBusinessSettleOneDao_InsertTx(t *testing.T) {
	//获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Error("获取数据库连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//开启事务
	tx, err := dbHandler.Begin()
	if err != nil {
		t.Errorf("开启事务失败, err:%v", err)
		return
	}

	id, err := BusinessSettleOneDaoInst.InsertTx(tx, nil)
	if err != nil {
		tx.Rollback()
		t.Errorf("InsertTx() error = %v", err)
		return
	}

	tx.Commit()
	t.Logf("结算日志id：%v", id)
}
