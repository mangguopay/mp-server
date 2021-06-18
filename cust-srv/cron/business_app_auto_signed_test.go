package cron

import (
	"testing"

	_ "a.a/mp-server/cust-srv/test"
)

func TestBusinessSigned_AutoSigned(t *testing.T) {
	BusinessAppAutoSignedTask.AutoSigned()
}
