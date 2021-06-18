package dao

import (
	_ "a.a/mp-server/business-bill-srv/test"
	"context"
	"testing"

	_ "a.a/mp-server/business-bill-srv/test"

	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
)

func TestVaccountDao_FreezeBalance(t *testing.T) {
	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Errorf("获取数据连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Errorf("开启事物失败,err:%v", err)
		return
	}

	vaccountNo := "ab28f001-0d9c-4c90-8131-4fb385233df2"
	amount := "30"

	balance, frozenBalance, aerr := VaccountDaoInst.FreezeBalance(tx, vaccountNo, amount)
	if aerr != nil {
		t.Errorf("FreezeBalance-error:%v, vaccountNo:%v, amount:%v", aerr, vaccountNo, amount)
		return
	}

	ss_sql.Commit(tx)

	t.Logf("balance:%v, frozenBalance:%v", balance, frozenBalance)
}

func TestVaccountDao_UnfreezeBalance(t *testing.T) {
	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Errorf("获取数据连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Errorf("开启事物失败,err:%v", err)
		return
	}

	vaccountNo := "ab28f001-0d9c-4c90-8131-4fb385233df2"
	amount := "30"

	balance, frozenBalance, aerr := VaccountDaoInst.UnfreezeBalance(tx, vaccountNo, amount)
	if aerr != nil {
		t.Errorf("UnfreezeBalance-error:%v, vaccountNo:%v, amount:%v", aerr, vaccountNo, amount)
		return
	}

	ss_sql.Commit(tx)

	t.Logf("balance:%v, frozenBalance:%v", balance, frozenBalance)

}

func TestVaccountDao_PlusBalance(t *testing.T) {
	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Errorf("获取数据连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Errorf("开启事物失败,err:%v", err)
		return
	}

	vaccountNo := "ab28f001-0d9c-4c90-8131-4fb385233df2"
	amount := "30"

	balance, frozenBalance, aerr := VaccountDaoInst.PlusBalance(tx, vaccountNo, amount)
	if aerr != nil {
		t.Errorf("PlusBalance-error:%v, vaccountNo:%v, amount:%v", aerr, vaccountNo, amount)
		return
	}

	ss_sql.Commit(tx)

	t.Logf("balance:%v, frozenBalance:%v", balance, frozenBalance)
}

func TestVaccountDao_MinusBalance(t *testing.T) {
	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Errorf("获取数据连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Errorf("开启事物失败,err:%v", err)
		return
	}

	vaccountNo := "ab28f001-0d9c-4c90-8131-4fb385233df2"
	amount := "30"

	balance, frozenBalance, aerr := VaccountDaoInst.MinusBalance(tx, vaccountNo, amount)
	if aerr != nil {
		t.Errorf("MinusBalance-error:%v, vaccountNo:%v, amount:%v", aerr, vaccountNo, amount)
		return
	}

	ss_sql.Commit(tx)

	t.Logf("balance:%v, frozenBalance:%v", balance, frozenBalance)
}

func TestVaccountDao_MinusFrozenBalance(t *testing.T) {
	// 获取数据库连接
	dbHandler := db.GetDB(constants.DB_CRM)
	if dbHandler == nil {
		t.Errorf("获取数据连接失败")
		return
	}
	defer db.PutDB(constants.DB_CRM, dbHandler)

	// 开启事物
	tx, err := dbHandler.BeginTx(context.TODO(), nil)
	if err != nil {
		t.Errorf("开启事物失败,err:%v", err)
		return
	}

	vaccountNo := "ab28f001-0d9c-4c90-8131-4fb385233df2"
	amount := "30"

	balance, frozenBalance, aerr := VaccountDaoInst.MinusFrozenBalance(tx, vaccountNo, amount)
	if aerr != nil {
		t.Errorf("MinusFrozenBalance-error:%v, vaccountNo:%v, amount:%v", aerr, vaccountNo, amount)
		return
	}

	ss_sql.Commit(tx)

	t.Logf("balance:%v, frozenBalance:%v", balance, frozenBalance)
}

func TestVaccountDao_GetBalanceByVAccNo(t *testing.T) {
	vAccountNo := "997cabc4-76b7-4572-b535-eda8f179f0e2"
	balance, err := VaccountDaoInst.GetBalanceByVAccNo(vAccountNo)
	if err != nil {
		t.Errorf("GetBalanceByVAccNo() error = %v", err)
		return
	}

	t.Logf("账户(%v)余额: %v", vAccountNo, balance)
}
