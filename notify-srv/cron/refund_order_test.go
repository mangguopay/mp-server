package cron

import (
	"testing"

	_ "a.a/mp-server/notify-srv/test"
)

func TestRefundNotify_Doing(t *testing.T) {
	RefundNotifyTask.Run()
}
