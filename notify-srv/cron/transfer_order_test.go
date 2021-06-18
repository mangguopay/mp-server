package cron

import (
	_ "a.a/mp-server/notify-srv/test"
	"testing"

	_ "a.a/mp-server/notify-srv/test"
)

func TestTransferNotify_QueryNotifyOmissionOrder(t *testing.T) {
	TransferNotifyTask.QueryNotifyOmissionOrder()
}

func TestTransferNotify_QueryNotifyBreakOrder(t *testing.T) {
	TransferNotifyTask.QueryNotifyBreakOrder()
}
