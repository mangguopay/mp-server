package handler

import (
	"a.a/cu/db"
	_ "a.a/mp-server/auth-srv/test"
	"a.a/mp-server/common/constants"
	"context"
	"testing"
)

func TestSyncFreezeVAccountBalance(t *testing.T) {
	accountNo := "f9b706be-e869-4f48-a009-93fae1983ed9"

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		t.Errorf("开启事务失败,errTx=%v", errTx)
		return
	}

	if err := VAccountHandlerInst.SyncFreezeVAccountBalance(tx, accountNo); err != nil {
		tx.Rollback()
		t.Errorf("SyncFreezeVAccountBalance() error = %v", err)
		return
	}

	tx.Commit()
	t.Logf("同步成功")
}

func TestVAccountHandler_InitVAccount(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		t.Errorf("开启事务失败,errTx=%v", errTx)
		return
	}

	accountNo := "b8c8414e-5061-4e3a-b6b0-78a0a53d559d"
	if err := VAccountHandlerInst.InitVAccount(tx, accountNo); err != nil {
		t.Errorf("InitVAccount() error = %v", err)
		tx.Rollback()
		return
	}

	tx.Commit()
	t.Logf("初始化成功")

}
