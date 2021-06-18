package adapter

import (
	"context"
	"testing"

	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	pushProto "a.a/mp-server/common/proto/push"
	_ "a.a/mp-server/push2-srv/test"
)

//发送短信验证码
func TestPushSms(t *testing.T) {
	num := util.RandomDigitStrOnlyNum(6)
	t.Logf("sms=[%v]", num)
	req := pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				Phone:       "13298690108",
				Lang:        "zh_CN",
				CountryCode: "86",
			},
		},
		TempNo: constants.Template_Reg,
		Args: []string{
			num,
		},
	}
	if err := Push(context.TODO(), &req); err != nil {
		t.Errorf("Push() error = %v", err)
	}

}

//发送邮件验证码
func TestPushEmail(t *testing.T) {
	num := util.RandomDigitStrOnlyNum(6)
	t.Logf("sms=[%v]", num)
	req := pushProto.PushReqest{
		Accounts: []*pushProto.PushAccout{
			{
				Email: "h13298690108@163.com",
				Lang:  "en",
			},
		},
		TempNo: constants.Template_VerifyEmail,
		Args: []string{
			"h13298690108@163.com",
			num,
		},
		TitleArgs: []string{
			num,
			"h13298690108@163.com",
		},
	}
	if err := Push(context.TODO(), &req); err != nil {
		t.Errorf("Push() error = %v", err)
	}

}
