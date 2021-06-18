package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"testing"
)

func TestTransferDao_InsertTransfer(t *testing.T) {
	log := new(TransferDao)
	log.FromVacc = "589370e5-0784-42ca-9dba-0efa2a3f5822"
	log.ToVacc = "4facbf3f-399f-4c1d-8c0d-edbe6fa71401"
	log.Amount = "1000"
	log.ExchangeType = "2"
	log.Fees = "30"
	log.MoneyType = "usd"
	log.FeeRate = "300"
	log.RealAmount = "1000"
	log.Lat = "23.08592800"
	log.Lng = "113.34624500"
	log.Ip = "10.41.6.251"

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	tx, _ := dbHandler.Begin()

	orderNo, err := TransferDaoInst.InsertTransfer(tx, log)
	if err != nil {
		ss_sql.Rollback(tx)
		t.Errorf("InsertTransfer() error = %v", err)
		return
	}

	ss_sql.Commit(tx)
	t.Logf("订单号：%v", orderNo)
}
