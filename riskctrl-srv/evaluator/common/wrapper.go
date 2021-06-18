package common

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/riskctrl-srv/m"
)

var (
	RiskEvaluatorWrapperInst RiskEvaluatorWrapper
)

type RiskEvaluatorWrapper struct {
}

// 评估风控处理方式
func (*RiskEvaluatorWrapper) Eva(evaluatorType string, req *m.RiskEvaluatorContext) {
	targetApi := getTargetApi(evaluatorType)
	if targetApi == nil {
		ss_log.Error("[%v]is nil", evaluatorType)
		return
	}
	targetApi.Eva(req)
	return
}
