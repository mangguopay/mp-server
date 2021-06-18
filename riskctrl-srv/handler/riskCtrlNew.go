package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/proto/riskctrl"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/riskctrl-srv/dao"
	"a.a/mp-server/riskctrl-srv/evaluator"
	"context"
	"strings"
)

// 登录位置风险评估
func (RiskCtrlHandler) Login(ctx context.Context, req *go_micro_srv_riskctrl.LoginRequest, reply *go_micro_srv_riskctrl.LoginReply) error {
	req.DeviceId = strings.TrimSpace(req.DeviceId)
	req.Ip = strings.TrimSpace(req.Ip)
	req.Uid = strings.TrimSpace(req.Uid)

	if req.DeviceId == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Ip == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	if req.Uid == "" {
		reply.ResultCode = ss_err.ERR_PARAM
		return nil
	}

	params := make(map[string]string)
	params["uid"] = req.Uid
	params["ip"] = req.Ip
	params["device_id"] = req.DeviceId

	// 日志记录
	riskLog := dao.InsertResultNewData{}
	riskLog.RiskNo = strext.GetDailyId()
	riskLog.Params = strext.ToJson(params)
	riskLog.ActionTime = ss_time.NowForPostgres(global.Tz)
	riskLog.Uid = req.Uid

	// 创建评估人，添加评估项， 进行评估
	em := evaluator.NewEvaluatorManager(evaluator.PositionLogin)
	em.AddItem(evaluator.NewDeviceItem(req.DeviceId))        // 评估设备
	em.AddItem(evaluator.NewIpItem(req.Ip))                  // 评估ip
	em.AddItem(evaluator.NewLoginPasswordWrongItem(req.Uid)) // 评估账号登录密码错误情况
	em.AddItem(evaluator.NewPayPasswordWrongItem(req.Uid))   // 评估账号支付密码错误情况

	// 计算总评分
	score, itemResult := em.CalculateScore()

	isStop := 0                                // 是否阻止当前行为
	if score > evaluator.LoginThresholdScore { // 超过阈值了
		isStop = 1
	}

	riskLog.Result = isStop
	riskLog.Position = string(em.GetPosition())
	riskLog.Threshold = evaluator.LoginThresholdScore
	riskLog.Score = score
	riskLog.ItemResult = strext.ToJson(itemResult)
	if err := dao.RiskResultDaoInstance.InsertResultNew(riskLog); err != nil {
		ss_log.Error("插入风控日志失败:%s", strext.ToJson(riskLog))
	}

	if isStop > 0 { // 风控不通过
		reply.OpResult = constants.Risk_Result_No_Pass_Str
		ss_log.Error("风控不通过,RiskNo:%s", riskLog.RiskNo)
	} else {
		reply.OpResult = constants.Risk_Result_Pass_Str
	}

	reply.RiskNo = riskLog.RiskNo
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
