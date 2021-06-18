package handler

import (
	go_micro_srv_riskctrl "a.a/mp-server/common/proto/riskctrl"
	"context"
	"testing"
)

func TestRiskCtrlHandler_Login(t *testing.T) {
	req := &go_micro_srv_riskctrl.LoginRequest{
		DeviceId: "aaaaaaa",
		Ip:       "192.168.1.123",
		Uid:      "49c71695-e29a-4309-91a2-27ebe1547563",
	}
	reply := &go_micro_srv_riskctrl.LoginReply{}

	ri := RiskCtrlHandler{}
	if err := ri.Login(context.TODO(), req, reply); err != nil {
		t.Errorf("Login() error = %v,", err)
		return
	}

	t.Errorf("reply.ResultCode = %v,", reply.ResultCode)
}
