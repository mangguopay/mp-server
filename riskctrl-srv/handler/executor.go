package handler

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/riskctrl-srv/common"
	"a.a/mp-server/riskctrl-srv/dao"
	"a.a/mp-server/riskctrl-srv/m"
	"errors"
	"fmt"
)

type Executor interface {
	Execute(op *m.Op) *ExecContext
}

type PipelineExecuter struct{}
type IfExecuter struct{}
type ForExecuter struct{}
type SwitchExecuter struct{}

var (
	PipelineExecuterInst PipelineExecuter
	IfExecuterInst       IfExecuter
	ForExecuterInst      ForExecuter
	SwitchExecuterInst   SwitchExecuter
)

// PipelineExecuter
//func (*PipelineExecuter) Execute(op *m.Op) *ExecContext {
//	var scriptName string
//	var param string
//	_, scriptName, param = common.GetOpMsgFromOpNo(op.Steps) // 读取操作的id对应的数据
//	var m map[string]string
//	_ = jsoniter.Unmarshal([]byte(param), &m)
//	return &ExecContext{
//		Cmd:   scriptName,
//		Param: m,
//	}
//}

// ForExecuter
//func (*ForExecuter) Execute(op *m.Op) *ExecContext {
//	var scriptName string
//	var param string
//	opNo := op.Condition
//	_, scriptName, param = common.GetOpMsgFromOpNo(opNo) // 读取操作的id对应的数据
//	var m map[string]string
//	_ = jsoniter.Unmarshal([]byte(param), &m)
//	return &ExecContext{
//		Cmd:   scriptName,
//		Param: m,
//	}
//}

// SwitchExecuter
//func (*SwitchExecuter) Execute(op *Op) *ExecContext {
//	var scriptName string
//	var param string
//	opNo := op.Condition
//	_, scriptName, param = common.GetOpMsgFromOpNo(opNo) // 读取操作的id对应的数据
//	var m map[string]string
//	_ = jsoniter.Unmarshal([]byte(param), &m)
//	return &ExecContext{
//		Cmd:   scriptName,
//		Param: m,
//	}
//}

func Init(rule string) ([]*m.Op, string) {
	opArr := unmarshalRuleStr(rule)
	if len(opArr) == 0 {
		ss_log.Error("err=[---->%s]", "风控结果接口,解析rule失败")
		return nil, ss_err.ERR_PARAM
	}
	return opArr, ss_err.ERR_SUCCESS
}

type ExecContext struct {
	Param      map[string]string
	Score      int32  // [0,10000]
	LastResult string // 流程控制参数
}

func (r *RuleParser) doAction(cmd string, param map[string]string) error {
	switch cmd {
	case "getLoginTimeInNSec":
		//xxxx 调用该函数,执行相对应的代码返回结果.
		r.Ctx.Param[common.ExeParam_LoginTimes] = r.getLoginTimesInNTime(param)
		return nil
	case constants.Risk_Ctrl_Save_Money:
		return nil
	case constants.Risk_Ctrl_Mobile_Num_Withdrawal:
		return nil
	case constants.Risk_Ctrl_Exchange:
		return nil
	case constants.Risk_Ctrl_Transfer:
		return nil
	case constants.Risk_Ctrl_Collection:
		return nil
	case constants.Risk_Ctrl_Sweep_Withdrawal:
		return nil
	default:
		ss_log.Error("no cmd[%v]", cmd)
	}
	return errors.New(fmt.Sprintf("no cmd[%v]", cmd))
}

func (r *RuleParser) getLoginTimesInNTime(param map[string]string) string {
	timeType := common.ParsingTimeType(param["type"])
	value := param["time"]
	loginCount := dao.LogLoginDaoInstance.GetCountFromTime(fmt.Sprintf("-%s %s", value, timeType))
	// 查询阈值
	//threshold := dao.RiskThresholdDaoInstance.GetThresholdFromOpNo(opNo)
	//ss_log.Info("%s %s 时间内,登录的条数为----->%s,阈值为----->%s", value, timeType, loginCount, threshold)
	return loginCount
}
