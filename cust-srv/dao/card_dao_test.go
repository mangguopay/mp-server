package dao

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"reflect"
	"testing"
)

func TestCardDao_GetCustCards(t *testing.T) {
	type args struct {
		whereList []*model.WhereSqlCond
	}
	tests := []struct {
		name            string
		args            args
		wantReturnDatas []*go_micro_srv_cust.UserCardsData
		wantReturnTotal string
		wantReturnErr   string
	}{
		{
			args: args{
				whereList: []*model.WhereSqlCond{
					{Key: "ca.account_no", Val: "0e8d24af-bec7-4f95-b038-c48045f51abf", EqType: "="},
					{Key: "ca.collect_status", Val: "1", EqType: "="},
					{Key: "ca.is_delete", Val: "0", EqType: "="},
					{Key: "ca.balance_type", Val: "", EqType: "="},
					{Key: "ca.account_type", Val: constants.AccountType_USER, EqType: "="},
				},
			}, // TODO: Add test cases.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := CardDao{}
			gotReturnDatas, gotReturnTotal, gotReturnErr := ca.GetCustCards(tt.args.whereList)
			if !reflect.DeepEqual(gotReturnDatas, tt.wantReturnDatas) {
				t.Errorf("GetCustCards() gotReturnDatas = %v, want %v", gotReturnDatas, tt.wantReturnDatas)
			}
			if gotReturnTotal != tt.wantReturnTotal {
				t.Errorf("GetCustCards() gotReturnTotal = %v, want %v", gotReturnTotal, tt.wantReturnTotal)
			}
			if gotReturnErr != tt.wantReturnErr {
				t.Errorf("GetCustCards() gotReturnErr = %v, want %v", gotReturnErr, tt.wantReturnErr)
			}
		})
	}
}
