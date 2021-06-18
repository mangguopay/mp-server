package common

import (
	"a.a/mp-server/common/constants"
)

func getTargetApi(evaluatorType string) iRiskEvaluator {
	switch evaluatorType {
	case constants.RiskEvaluator_Amount:
		return AmountEvaluator{}
	}
	return nil
}
