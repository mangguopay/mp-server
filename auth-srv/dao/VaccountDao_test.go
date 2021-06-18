package dao

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"context"
	"testing"
)

func TestVaccountDao_SyncAccRemain(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		t.Errorf("开启事务失败,errTx=%v", errTx)
		return
	}

	accNo := "f9b706be-e869-4f48-a009-93fae1983ed9"
	if err := VaccountDaoInst.SyncAccRemain(tx, accNo); err != nil {
		t.Errorf("SyncAccRemain() error = %v", err)
	}

}

func TestVaccountDao_GetVAccNoByAccountNo(t *testing.T) {
	accountNo := "0e8d24af-bec7-4f95-b038-c48045f51abf"
	got, err := VaccountDaoInst.GetVAccNoByAccountNo(accountNo)
	if err != nil {
		t.Errorf("GetVAccNoByAccountNo() error = %v", err)
		return
	}
	t.Logf("虚账:%v", strext.ToJson(got))
}
