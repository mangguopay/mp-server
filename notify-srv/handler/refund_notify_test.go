package handler

import (
	"testing"

	_ "a.a/mp-server/notify-srv/test"
)

func TestRefundNotifyHandler_RefundNotify(t *testing.T) {
	refundNo := "2020102814251557394466"
	RefundNotifyH.RefundNotify(refundNo)
}
