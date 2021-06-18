package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"context"
	"testing"
)

func TestBusinessFixedCodeDao_AddBusinessFixedCodeTx(t *testing.T) {

	type args struct {
		businessAccountNo string
		businessNo        string
		signedNo          string
	}
	tests := []struct {
		name           string
		args           args
		wantStaticCode string
		wantErr        bool
	}{
		{
			args: args{
				businessAccountNo: "6b669ba3-2496-41e2-90de-aaf547e2523a",
				businessNo:        "48b36f10-de17-4b6a-af6d-cac604220e8b",
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbHandler := db.GetDB(constants.DB_CRM)
			defer db.PutDB(constants.DB_CRM, dbHandler)

			tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
			if errTx != nil {
				ss_log.Error("开启事务失败,errTx=[%v]", errTx)
			}

			bu := BusinessFixedCodeDao{}
			gotStaticCode, err := bu.AddBusinessFixedCodeTx(tx, tt.args.businessAccountNo, tt.args.businessNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBusinessFixedCodeTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("AddBusinessFixedCodeTx() gotStaticCode = %v", gotStaticCode)

			tx.Commit()
		})
	}
}
