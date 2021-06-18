package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"context"
	"database/sql"
	"testing"
)

func TestAppFingerprintDao_SetUseStatusDisableByAccount(t *testing.T) {
	type args struct {
		accountNo string
		devicUuid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{args: args{accountNo: "58fa37ce-24d7-4423-a5f3-5557f132ccc6", devicUuid: "482fa1bb91a30505"}}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := &AppFingerprintDao{}
			if err := ap.SetUseStatusDisableByAccount(tt.args.accountNo, tt.args.devicUuid); (err != nil) != tt.wantErr {
				t.Errorf("SetUseStatusDisableByAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppFingerprintDao_AddTx(t *testing.T) {
	type args struct {
		tx        *sql.Tx
		accountNo string
		devicUuid string
	}

	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
	}

	tests := []struct {
		name        string
		args        args
		wantSignKey string
		wantErr     bool
	}{
		{args: args{accountNo: "58fa37ce-24d7-4423-a5f3-5557f132ccc6", devicUuid: "482fa1bb91a30505", tx: tx}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := &AppFingerprintDao{}
			gotSignKey, err := ap.AddTx(tt.args.tx, tt.args.accountNo, tt.args.devicUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSignKey != tt.wantSignKey {
				t.Errorf("AddTx() gotSignKey = %v, want %v", gotSignKey, tt.wantSignKey)
			}
			tx.Commit()
		})
	}
}
