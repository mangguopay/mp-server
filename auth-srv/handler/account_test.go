package handler

import (
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"context"
	"testing"
)

func TestAuth_MobileModifyPwd(t *testing.T) {
	type args struct {
		ctx   context.Context
		req   *go_micro_srv_auth.MobileModifyPwdRequest
		reply *go_micro_srv_auth.MobileModifyPwdReply
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{args: args{
			req: &go_micro_srv_auth.MobileModifyPwdRequest{
				Uid:         "c930475e-2e76-4e28-91bb-2b3d478778e5",
				OldPassword: "1f82c942befda29b6ed487a51da199f78fce7f05",
				NewPassword: "b0a5c59d95469ad94c1391e2575ca734a8b740eb",
			},
			reply: &go_micro_srv_auth.MobileModifyPwdReply{},
			ctx:   context.TODO(),
		},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Auth{}
			if err := r.MobileModifyPwd(tt.args.ctx, tt.args.req, tt.args.reply); (err != nil) != tt.wantErr {
				t.Errorf("MobileModifyPwd() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Errorf("MobileModifyPwd() tt.args.reply = %v, wantErr %v", tt.args.reply, tt.wantErr)

		})
	}
}
