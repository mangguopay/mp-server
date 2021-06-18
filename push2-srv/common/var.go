package common

import (
	"github.com/micro/go-micro/v2/broker"
)

var (
	MqPushMsg *broker.Broker
	MqTopic   string
)

const (
	Pusher_Jpush       = "jpush"
	Pusher_Fcm         = "fcmpush"
	Pusher_ClSms       = "cl_sms"
	Pusher_ZtSms       = "zt_sms"
	Pusher_GoogleEmail = "google_email"
)
