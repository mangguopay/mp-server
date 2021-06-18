package evaluator

import (
	"a.a/cu/ss_log"
	"a.a/cu/util"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/riskctrl-srv/m"
	"testing"
)

func TestEva(t *testing.T) {
	ss_log.InitLog(".")
	resp := Eva(&m.RiskEvaReq{
		// 发起金额
		Amount: 10000,
		// 发起支付的账号
		PayerAccNo: ss_sql.UUID,
		// 收款人账号
		PayeeAccNo: ss_sql.UUID,
		// 创建时间
		ActionTime: util.NowWithFmt("2006-01-02 15:04:05"),
	})
	ss_log.Info("resp=[%v]", resp)
}
