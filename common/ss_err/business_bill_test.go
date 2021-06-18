package ss_err

import (
	"testing"

	"a.a/mp-server/common/constants"
)

func TestGetMsg(t *testing.T) {
	type args struct {
		code string
		lang string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "A01", args: args{code: Success, lang: constants.LangZhCN}, want: "成功"},
		{name: "A02", args: args{code: Success, lang: constants.LangEnUS}, want: "success"},
		{name: "A03", args: args{code: Success, lang: constants.LangKmKH}, want: "ជោគជ័យ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMsg(tt.args.code, tt.args.lang); got != tt.want {
				t.Errorf("GetMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
