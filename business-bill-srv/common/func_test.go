package common

import (
	"a.a/mp-server/common/constants"
	"testing"
)

func TestNormalAmountByMoneyType(t *testing.T) {
	type args struct {
		moneyType string
		amount    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"A01", args{moneyType: constants.CURRENCY_USD, amount: "100"}, "1"},
		{"A02", args{moneyType: constants.CURRENCY_KHR, amount: "200"}, "200"},
		{"A03", args{moneyType: "", amount: "100"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalAmountByMoneyType(tt.args.moneyType, tt.args.amount); got != tt.want {
				t.Errorf("NormalAmountByMoneyType() = %v, want %v", got, tt.want)
			}
		})
	}
}
