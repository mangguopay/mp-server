package common

import (
	"a.a/mp-server/common/constants"
	"testing"
)

func TestFeesTypeByMoneyType(t *testing.T) {
	type args struct {
		scene     int
		moneyType string
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add test cases.
		{"A01", args{scene: 1, moneyType: constants.CURRENCY_USD}, 4, false},
		{"A02", args{scene: 1, moneyType: constants.CURRENCY_KHR}, 9, false},
		{"A03", args{scene: 2, moneyType: constants.CURRENCY_USD}, 5, false},
		{"A04", args{scene: 2, moneyType: constants.CURRENCY_KHR}, 10, false},
		{"A05", args{scene: 0, moneyType: constants.CURRENCY_USD}, 0, true},
		{"A06", args{scene: 3, moneyType: constants.CURRENCY_USD}, 2, false},
		{"A07", args{scene: 3, moneyType: constants.CURRENCY_KHR}, 7, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FeesTypeByMoneyType(tt.args.scene, tt.args.moneyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("FeesTypeByMoneyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FeesTypeByMoneyType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestVirtualAccountTypeByMoneyType(t *testing.T) {
	type args struct {
		moneyType string
		isActived string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{"A01", args{moneyType: constants.CURRENCY_USD, isActived: "0"}, constants.VaType_FREEZE_USD_DEBIT, false},
		{"A02", args{moneyType: constants.CURRENCY_KHR, isActived: "0"}, constants.VaType_FREEZE_KHR_DEBIT, false},
		{"A03", args{moneyType: constants.CURRENCY_USD, isActived: "1"}, constants.VaType_USD_DEBIT, false},
		{"A04", args{moneyType: constants.CURRENCY_KHR, isActived: "1"}, constants.VaType_KHR_DEBIT, false},
		{"A05", args{moneyType: "", isActived: "1"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VirtualAccountTypeByMoneyType(tt.args.moneyType, tt.args.isActived)
			if (err != nil) != tt.wantErr {
				t.Errorf("VirtualAccountTypeByMoneyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VirtualAccountTypeByMoneyType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPostgresDate(t *testing.T) {
	type args struct {
		dateStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"A01", args{dateStr: "2020-04-10T00:00:00Z"}, "2020-04-10"},
		{"A02", args{dateStr: ""}, ""},
		{"A03", args{dateStr: "T"}, "T"},
		{"A04", args{dateStr: "xxxx"}, "xxxx"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPostgresDate(tt.args.dateStr); got != tt.want {
				t.Errorf("GetPostgresDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
