package dao

import (
	"reflect"
	"testing"

	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
)

func TestBusinessSceneDao_GetBusinessSceneDetail(t *testing.T) {
	type args struct {
		sceneNo string
	}
	tests := []struct {
		name     string
		args     args
		wantData *go_micro_srv_cust.BusinessSceneData
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args: args{sceneNo: "d714c8f3-b0cc-4409-8fee-b2936fd1386d"},
			//wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bu := BusinessSceneDao{}
			gotData, err := bu.GetBusinessSceneDetail(tt.args.sceneNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBusinessSceneDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("GetBusinessSceneDetail() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}
