package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/riskctrl"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/riskctrl-srv/common"
	"a.a/mp-server/riskctrl-srv/dao"
	"a.a/mp-server/riskctrl-srv/evaluator"
	"a.a/mp-server/riskctrl-srv/m"
	"context"
)

type RiskCtrlHandler struct{}

func (RiskCtrlHandler) RiskOffline(context.Context, *go_micro_srv_riskctrl.RiskOfflineRequest, *go_micro_srv_riskctrl.RiskOfflineReply) error {
	panic("implement me")
}

var RiskCtrlHandlerInst RiskCtrlHandler

// 下单
func (RiskCtrlHandler) GetRiskCtrlReuslt(ctx context.Context, req *go_micro_srv_riskctrl.GetRiskCtrlResultRequest, reply *go_micro_srv_riskctrl.GetRiskCtrlResultReply) error {
	if req.ApiType == "" {
		ss_log.Error("err=[---->%s]", "风控结果接口,参数为空")
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	// 估算初步风险值
	resp := evaluator.Eva(&m.RiskEvaReq{
		// 发起金额
		Amount: strext.ToInt64(req.Amount),
		// 发起支付的账号
		PayerAccNo: req.PayerAccNo,
		// 收款人账号
		PayeeAccNo: req.PayeeAccNo,
		// 创建时间
		ActionTime: req.ActionTime,
		PayType:    req.PayType,
	})
	riskNo := strext.GetDailyId()
	var opResult, result, errStr string
	switch resp.RiskExecuteType {
	case constants.RiskExecuteType_Online:
		result, _, errStr = riskResult(req.ApiType, req.PayerAccNo, req.ActionTime, resp.RiskExecuteType, strext.ToStringNoPoint(resp.Score), req.MoneyType, req.OrderNo, riskNo, req.ProductType)
		if errStr != ss_err.ERR_SUCCESS {
			opResult = constants.Risk_Result_No_Pass_Str
		} else {
			if result == constants.Risk_Result_No_Pass {
				opResult = constants.Risk_Result_No_Pass_Str
			} else {
				opResult = constants.Risk_Result_Pass_Str
			}
		}
	case constants.RiskExecuteType_Half:
		if util.Random(0, 10000) > 5000 {
			result, _, errStr = riskResult(req.ApiType, req.PayerAccNo, req.ActionTime, resp.RiskExecuteType, strext.ToStringNoPoint(resp.Score), req.MoneyType, req.OrderNo, riskNo, req.ProductType)
			if errStr != ss_err.ERR_SUCCESS {
				opResult = constants.Risk_Result_No_Pass_Str
			} else {
				if result == constants.Risk_Result_No_Pass {
					opResult = constants.Risk_Result_No_Pass_Str
				} else {
					opResult = constants.Risk_Result_Pass_Str
				}
			}
		} else {
			offline(req.ApiType, req.PayerAccNo, req.ActionTime, resp.RiskExecuteType, strext.ToStringNoPoint(resp.Score), req.MoneyType, req.OrderNo, riskNo, req.ProductType)
		}
	case constants.RiskExecuteType_Offline:
		opResult = constants.Risk_Result_Pending_Str
		offline(req.ApiType, req.PayerAccNo, req.ActionTime, resp.RiskExecuteType, strext.ToStringNoPoint(resp.Score), req.MoneyType, req.OrderNo, riskNo, req.ProductType)
	}
	// 离线时是pending的状态.
	reply.OpResult = opResult
	reply.RiskNo = riskNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

func (RiskCtrlHandler) GetRiskCtrlReuslt2(ctx context.Context, req *go_micro_srv_riskctrl.GetRiskCtrlResult2Request, reply *go_micro_srv_riskctrl.GetRiskCtrlResult2Reply) error {
	riskResult := dao.RiskResultDaoInstance.GetRiskResult(req.RiskNo)
	if riskResult == "" || riskResult == constants.Risk_Result_No_Pass { // 被风控了
		ss_log.Error("err=[被风控了,riskNo为-----%s]", req.RiskNo)
		//reply.ResultCode = ss_err.ERR_RISK_IS_RISK
		reply.ResultCode = ss_err.ERR_SUCCESS
		return nil
	}
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}

// 推送消息进队列
func offline(apiType, payerAccNo, actionTime, evaExecuteType, evaScore, moneyType, orderNo, riskNo, ProductType string) {

	// 消息推送
	ev := &go_micro_srv_riskctrl.RiskOfflineRequest{
		ApiType:        apiType,
		PayerAccNo:     payerAccNo,
		ActionTime:     actionTime,
		RiskNo:         riskNo,
		ProductType:    ProductType,
		EvaExecuteType: evaExecuteType,
		EvaScord:       evaScore,
		MoneyType:      moneyType,
		OrderNo:        orderNo,
	}
	ss_log.Info("publishing %+v\n", ev)
	// publish an event
	if err := common.RiskEvent.Publish(context.TODO(), ev); err != nil {
		ss_log.Error("err=[风控接口,请求参数推送到MQ失败,err----->%s]", err.Error())
	}
}

// 队列消费端
func RecvMQMsg(msg *go_micro_srv_riskctrl.RiskOfflineRequest) {
	riskResult(msg.ApiType, msg.PayerAccNo, msg.ActionTime, msg.EvaExecuteType, msg.EvaScord, msg.MoneyType, msg.OrderNo, msg.RiskNo, msg.ProductType)
}

func riskResult(apiType, payerAccNo, actionTime, evaExecuteType, evaScore, moneyType, orderNo, riskNo, productType string) (result, threshold, errStr string) {
	// 进redis
	_, rule, _ := dao.GlobalParamDaoInstance.GetGlobalParam(apiType)

	if rule == "" {
		errStr = ss_err.ERR_PARAM
		return
	}

	opArr, errCode := Init(rule)
	if errCode != ss_err.ERR_SUCCESS {
		errStr = errCode
		return
	}
	//  获取 rule_No
	rulefNo := dao.RuleDaoInstance.GetRuleNoFromApiType(apiType)
	if rulefNo == "" {
		errStr = ss_err.ERR_PARAM
		return
	}

	ruleParser := DoInitRuleParser(rulefNo)
	err := ruleParser.ParsRule(opArr)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}

	result = strext.ToStringNoPoint(ruleParser.Result)
	score := ruleParser.Ctx.Score

	threshold = strext.ToStringNoPoint(ruleParser.ThresHold)

	// 生成uuid
	//riskNo = strext.GetDailyId()
	if errStr = dao.RiskResultDaoInstance.InsertResult(riskNo, result, threshold, apiType, payerAccNo, actionTime, evaExecuteType, evaScore, moneyType, orderNo, productType, score); errStr != ss_err.ERR_SUCCESS {
		ss_log.Error("err=[----------->%s]", "插入结果失败")
		return
	}
	errStr = ss_err.ERR_SUCCESS
	return
}
