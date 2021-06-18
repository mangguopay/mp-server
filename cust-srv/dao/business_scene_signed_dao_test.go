package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"context"
	"reflect"
	"testing"
)

func TestBusinessSceneSignedDao_SetStatusInvalidTx(t *testing.T) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	tx, errTx := dbHandler.BeginTx(context.TODO(), nil)
	if errTx != nil {
		ss_log.Error("开启事务失败,errTx=[%v]", errTx)
	}

	if err := BusinessSceneSignedDaoInst.SetStatusInvalidTx(tx, "a79123bd-0d06-418f-9c12-976a8643f82a", "246b72ad-2904-42b6-83b8-274dbf3a4927"); err != nil {
		t.Errorf("SetStatusInvalidTx() error = %v", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	t.Log("SetStatusInvalidTx() 成功")
}

func TestBusinessSceneSignedDao_GetBusinessSceneSignedDetail(t *testing.T) {
	type args struct {
		signedNo string
	}
	tests := []struct {
		name     string
		args     args
		wantData *BusinessSceneSignedData
		wantErr  bool
	}{
		{args: args{
			signedNo: "2020101917493899442243"}}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BusinessSceneSignedDao{}
			gotData, err := b.GetBusinessSceneSignedDetail(tt.args.signedNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBusinessSceneSignedDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("GetBusinessSceneSignedDetail() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
