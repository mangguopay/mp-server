package common

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/riskctrl-srv/m"
)

type PayTypeEvaluator struct {
}

var PayTypeEvaluatorInst PayTypeEvaluator

func (r PayTypeEvaluator) Eva(ctx *m.RiskEvaluatorContext) {
	if ctx.PayType == constants.Risk_Pay_Type_Transfer || ctx.PayType == constants.Risk_Pay_Type_Exchange || ctx.PayType == constants.Risk_Pay_Type_Save_Money || ctx.PayType == constants.Risk_Pay_Type_Mobile_Num_Withdrawal || ctx.PayType == constants.Risk_Pay_Type_Collection {
		ctx.RiskExecuteType = constants.RiskExecuteType_Online
	}
}
