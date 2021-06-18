package common

import (
	"a.a/mp-server/push2-srv/common"
)

func getTargetApi(apiType string) IPushAdapter {
	switch apiType {
	case common.Pusher_Fcm:
		return &FcmAdapterInst
	case common.Pusher_ClSms:
		return &ClSmsAdapterInst
	case common.Pusher_ZtSms:
		return &ZtSmsAdapterInst
	case common.Pusher_GoogleEmail:
		return &GoogleEmailAdapterInst
	}
	return nil
}
