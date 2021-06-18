package handler

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	custProto "a.a/mp-server/common/proto/cust"
	"context"
	"testing"
)

func TestCustHandler_AddAuthMaterialBusiness(t *testing.T) {
	req := &custProto.AddAuthMaterialBusinessRequest{
		//ImgId:      "2020110915442831226946",
		AuthName: "鸡大保",
		//AuthNumber: "123456789",
		AccountUid: "972617f3-c85b-465b-ae3a-8491647d869d",
		StartDate:  "",

		EndDate:      "",
		TermType:     constants.TermType_Long,
		IndustryNo:   "001001",
		SimplifyName: "小鸡岛牛杂店",
	}
	reply := &custProto.AddAuthMaterialBusinessReply{}

	if err := CustHandlerInst.AddAuthMaterialBusiness(context.TODO(), req, reply); err != nil {
		t.Errorf("AddAuthMaterialBusiness() error = %v", err)
		return
	}
	t.Logf("结果：%v", strext.ToJson(reply))
}

func TestCustHandler_GetAuthMaterialBusinessDetail(t *testing.T) {
	type fields struct {
		Client custProto.CustService
	}
	type args struct {
		ctx   context.Context
		req   *custProto.GetAuthMaterialBusinessDetailRequest
		reply *custProto.GetAuthMaterialBusinessDetailReply
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			args: args{
				ctx: context.TODO(),
				req: &custProto.GetAuthMaterialBusinessDetailRequest{
					AccountUid: "58fa37ce-24d7-4423-a5f3-5557f132ccc6",
				},
				reply: &custProto.GetAuthMaterialBusinessDetailReply{},
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cu := &CustHandler{
				Client: tt.fields.Client,
			}
			if err := cu.GetAuthMaterialBusinessDetail(tt.args.ctx, tt.args.req, tt.args.reply); (err != nil) != tt.wantErr {
				t.Errorf("GetAuthMaterialBusinessDetail() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("GetAuthMaterialBusinessDetail() ,tt.args.reply -----> %v  , wantErr %v", tt.args.reply, tt.wantErr)
		})
	}
}

func TestCustHandler_ModifyAuthMaterialBusinessStatus(t *testing.T) {
	type fields struct {
		Client custProto.CustService
	}
	type args struct {
		ctx   context.Context
		req   *custProto.ModifyAuthMaterialBusinessStatusRequest
		reply *custProto.ModifyAuthMaterialBusinessStatusReply
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			args: args{
				ctx: context.TODO(),
				req: &custProto.ModifyAuthMaterialBusinessStatusRequest{
					AuthMaterialNo: "739841cc-88f9-4905-a9ff-2757b8162055",
					Status:         constants.AuthMaterialStatus_Passed,
					LoginUid:       "3ef90251-88ca-41f8-8a3d-051b8a1077ca",
				},
				reply: &custProto.ModifyAuthMaterialBusinessStatusReply{},
			},
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cu := &CustHandler{
				Client: tt.fields.Client,
			}
			if err := cu.ModifyAuthMaterialBusinessStatus(tt.args.ctx, tt.args.req, tt.args.reply); (err != nil) != tt.wantErr {
				t.Errorf("ModifyAuthMaterialBusinessStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
