package common

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/riskctrl-srv/m"
)

type AmountEvaluator struct {
}

func (r AmountEvaluator) Eva(ctx *m.RiskEvaluatorContext) {
	if ctx.Amount > 100000 /*100.00*/ {
		ctx.RiskExecuteType = constants.RiskExecuteType_Online
		return
	}

	// 大概先定这么个大小吧
	ctx.Score += 10
}
