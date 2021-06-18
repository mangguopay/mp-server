package dao

import (
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"reflect"
	"testing"
)

func TestCardDao_GetCustPaymentCard(t *testing.T) {
	type args struct {
		accountNo string
	}
	tests := []struct {
		name       string
		args       args
		wantDatasR []*go_micro_srv_auth.CardData
		wantErrR   string
	}{
		{args: args{accountNo: "0e8d24af-bec7-4f95-b038-c48045f51abf"}}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := &CardDao{}
			gotDatasR, gotErrR := ca.GetCustPaymentCard(tt.args.accountNo)
			if !reflect.DeepEqual(gotDatasR, tt.wantDatasR) {
				t.Errorf("GetCustPaymentCard() gotDatasR = %v, want %v", gotDatasR, tt.wantDatasR)
			}
			if gotErrR != tt.wantErrR {
				t.Errorf("GetCustPaymentCard() gotErrR = %v, want %v", gotErrR, tt.wantErrR)
			}
		})
	}
}
