package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"context"
	"database/sql"
	"testing"
)

func TestLogAppFingerprintPayDao_AddTx(t *testing.T) {
	type args struct {
		TX   *sql.Tx
		data *LogAppFingerprintPayData
	}
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{args: args{TX: tx, data: &LogAppFingerprintPayData{
			AccountNo:    "58fa37ce-24d7-4423-a5f3-5557f132ccc6",
			DeviceUuid:   "482fa1bb91a30505",
			SignKey:      "18ed5b5a4c34a5910b415d7f0f9727035f7ec6c86c9e79c933f1ab23a5cf235b26d3fb4ee7cd53158b7054fc60660090a91ccc7b4b3ca447b0e1109eb75fb507ogtmdbzdsh",
			OrderNo:      "123111111",
			OrderType:    constants.VaReason_Exchange,
			Amount:       "10",
			CurrencyType: "usd",
		}}}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lo := LogAppFingerprintPayDao{}
			if err := lo.AddTx(tt.args.TX, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("AddTx() error = %v, wantErr %v", err, tt.wantErr)
			}
			tx.Commit()
		})
	}
}
