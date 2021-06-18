package cron

import "testing"

func TestPayOrderTimeout_Doing(t *testing.T) {
	PayOrderTimeoutTask.Doing()
}

func TestRefundOrderTimeout_Doing(t *testing.T) {
	RefundOrderTimeoutTask.Doing()
}
