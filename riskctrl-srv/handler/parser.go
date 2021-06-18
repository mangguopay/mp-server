package handler

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/riskctrl-srv/dao"
	"a.a/mp-server/riskctrl-srv/m"
	"errors"
	"github.com/json-iterator/go"
)

type RuleParser struct {
	Ctx       *ExecContext
	ThresHold int32 // [0,10000]
	Result    bool
}

func DoInitRuleParser(ruleNo string) *RuleParser {
	i := new(RuleParser)
	i.Ctx = &ExecContext{
		Score: 0,
		Param: make(map[string]string),
	}
	threshold := dao.RuleDaoInstance.GetRuleThreshold(ruleNo)
	i.ThresHold = threshold
	return i
}

// ===================
func (r *RuleParser) ParsRule(op []*m.Op) error {
	for _, v := range op {
		switch v.Type {
		case "pipeline": // 顺序
			err := r.pipelineRoute(v)
			if err != nil {
				break
			}
		case "if":
			err := r.ifRoute(v)
			if err != nil {
				break
			}
		case "switch": // switch
			err := r.switchRoute(v)
			if err != nil {
				break
			}
		}

		r.Result = r.Ctx.Score >= r.ThresHold
		if r.Result {
			return errors.New("not passed")
		}
	}

	return nil
}

//===========================================

// 操作
func (r *RuleParser) doOp(op *m.Op) error {
	r.Ctx.LastResult = ""
	_, scriptName, param, score := dao.OpDaoInstance.GetOpFromNo(op.Steps) // 读取操作的id对应的数据
	var m map[string]string
	_ = jsoniter.Unmarshal([]byte(param), &m)
	err := r.doAction(scriptName, m)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		if op.Rollback != "" {
			_, scriptName, param, score = dao.OpDaoInstance.GetOpFromNo(op.Rollback) // 读取操作的id对应的数据
			err2 := r.doAction(scriptName, m)
			if err2 != nil {
				ss_log.Error("err=[%v]", err2)
			}
			// 由回滚脚本控制他的分数
			r.Ctx.Score += score
		}
		return err
	}
	// 处理完成,分数调整
	r.Ctx.Score += score
	return nil
}

// 解析
func unmarshalRuleStr(rule string) []*m.Op {
	st := &m.Setps{
		Steps: []*m.Op{},
	}
	if err := jsoniter.Unmarshal([]byte(rule), &st); err != nil {
		ss_log.Error("------->%s", err.Error())
		return nil
	}
	return st.Steps
}

func (r *RuleParser) pipelineRoute(op *m.Op) error {
	// pipeline没有控制参数
	return r.doOp(op)
}

func (r *RuleParser) ifRoute(op *m.Op) error {
	err := r.doOp(op)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	// 判断
	if strext.ToBool(r.Ctx.LastResult) {
		if len(op.Y) > 0 {
			return r.ParsRule(op.Y)
		}
	} else {
		if len(op.N) > 0 {
			return r.ParsRule(op.N)
		}
	}
	return nil
}

func (r *RuleParser) switchRoute(op *m.Op) error {
	err := r.doOp(op)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	if op.Ops[r.Ctx.LastResult] == nil {
		if op.Ops["default"] == nil {
			return nil
		}
		return r.ParsRule(op.Ops["default"])
	}
	return r.ParsRule(op.Ops[r.Ctx.LastResult])
}
