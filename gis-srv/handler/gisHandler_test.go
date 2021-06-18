package handler

import (
	"a.a/mp-server/common/proto/gis"
	"context"
	"testing"
)

func TestGisHandler_GetNearbyServicerList(t *testing.T) {
	type args struct {
		ctx   context.Context
		req   *go_micro_srv_gis.GetNearbyServicerListRequest
		reply *go_micro_srv_gis.GetNearbyServicerListReply
	}
	tests := []struct {
		name    string
		g       GisHandler
		args    args
		wantErr bool
	}{
		{
			args: args{
				ctx:   context.TODO(),
				req:   &go_micro_srv_gis.GetNearbyServicerListRequest{Lng: "32.333", Lat: "33.45"},
				reply: &go_micro_srv_gis.GetNearbyServicerListReply{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GisHandler{}
			if err := g.GetNearbyServicerList(tt.args.ctx, tt.args.req, tt.args.reply); (err != nil) != tt.wantErr {
				t.Errorf("GisHandler.GetNearbyServicerList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
