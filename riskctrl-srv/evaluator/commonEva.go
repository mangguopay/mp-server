package evaluator

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/riskctrl-srv/evaluator/common"
	"a.a/mp-server/riskctrl-srv/m"
)

/**
 ** 评估风控的执行方式
 *  如果是犹豫中，那么需要返回评分，最后返回的时候由评分来决定如果处理
 *  如果是非犹豫状态，那么就直接返回结果
 */
func Eva(req *m.RiskEvaReq) *m.RiskEvaRsp {
	l := []string{
		constants.RiskEvaluator_Amount,
	}
	ctx := m.RiskEvaluatorContext{
		Amount:          req.Amount,
		PayerAccNo:      req.PayerAccNo,
		PayeeAccNo:      req.PayeeAccNo,
		ActionTime:      req.ActionTime,
		RiskExecuteType: constants.RiskExecuteType_Hesitate,
		Score:           0,
		PayType:         req.PayType,
	}
	for _, v := range l {
		common.RiskEvaluatorWrapperInst.Eva(v, &ctx)
		if ctx.Score >= 10000 {
			ctx.RiskExecuteType = constants.RiskExecuteType_Online
		}
		if ctx.RiskExecuteType != constants.RiskExecuteType_Hesitate {
			ss_log.Info("终于决定了!=>[%v],score=[%v]", ctx.RiskExecuteType, ctx.Score)
			break
		}
	}

	if ctx.Score >= 10000 {
		ctx.RiskExecuteType = constants.RiskExecuteType_Online
	} else if ctx.Score >= 5000 {
		ctx.RiskExecuteType = constants.RiskExecuteType_Half
	} else {
		ctx.RiskExecuteType = constants.RiskExecuteType_Offline
	}

	common.PayTypeEvaluatorInst.Eva(&ctx)
	//// fixme 测试
	//ctx.RiskExecuteType = constants.RiskExecuteType_Online
	return &m.RiskEvaRsp{
		RiskExecuteType: ctx.RiskExecuteType,
		Score:           ctx.Score,
	}
}
