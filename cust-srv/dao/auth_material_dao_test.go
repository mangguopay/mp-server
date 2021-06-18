package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	_ "a.a/mp-server/cust-srv/test"
	"context"
	"testing"
)

func TestAuthMaterialDao_ModifyAuthMaterialBusinessStatus(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		return
	}

	authMaterialNo := ""
	status := ""
	oldStatus := ""
	au := AuthMaterialDao{}
	if err := au.ModifyAuthMaterialBusinessStatus(tx, authMaterialNo, status, oldStatus); err != nil {
		t.Errorf("ModifyAuthMaterialBusinessStatus() error = %v", err)
		tx.Rollback()
		return
	}
	tx.Commit()
}

func TestAuthMaterialDao_ModifyAuthMaterialEnterpriseStatus(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		return
	}

	authMaterialNo := "1f8f478f-8a9d-4dfa-a6de-612579d22615"
	status := "1"
	oldStatus := "1"
	au := AuthMaterialDao{}
	if err := au.ModifyAuthMaterialEnterpriseStatus(tx, authMaterialNo, status, oldStatus); err != nil {
		t.Errorf("ModifyAuthMaterialEnterpriseStatus() error = %v", err)
		tx.Rollback()
		return
	}

	tx.Commit()
}
