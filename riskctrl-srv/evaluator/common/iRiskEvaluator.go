package common

import "a.a/mp-server/riskctrl-srv/m"

type iRiskEvaluator interface {
	// 评估风控处理方式
	Eva(req *m.RiskEvaluatorContext)
}
