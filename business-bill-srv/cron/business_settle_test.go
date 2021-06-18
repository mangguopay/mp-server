package cron

import (
	_ "a.a/mp-server/business-bill-srv/test"
	"testing"
	"time"
)

func TestBusinessSettle_Doing(t *testing.T) {
	t.Logf("时间戳：%v", time.Now().Add(30*time.Minute).Unix())
	BusinessSettleTask.Doing()
}

func TestBusinessSettle_CheckUnSettledBalance(t *testing.T) {
	req := &CheckUnSettledBalanceReq{
		SettledId:     "",
		BusinessNo:    "",
		BusinessAccNo: "",
		CurrencyType:  "",
	}
	BusinessSettleTask.checkUnSettledBalance(req)

}
