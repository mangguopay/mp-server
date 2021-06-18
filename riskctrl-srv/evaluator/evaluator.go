package evaluator

import (
	"a.a/cu/ss_log"
)

type RiskPosition string

const (
	// 评估项名称
	ItemNameDivice             = "DeviceItem"             // 设备评估项名称
	ItemNameIp                 = "IpItem"                 // ip评估项名称
	ItemNameLoginPasswordWrong = "LoginPasswordWrongItem" // 登录密码错误评估项名称
	ItemNamePayPasswordWrong   = "PayPasswordWrongItem"   // 支付密码错误评估项名称

	// 评估位置的风险阈值
	LoginThresholdScore = 50 // 登录位置风险阈值

	// 风控位置(风控点)
	PositionLogin RiskPosition = "login"
)

// 评估项接口
type Item interface {
	Name() string
	Evaluate() (int, error)
}

type EvaluatorManager struct {
	items    []Item       // 需要计算的评估项
	position RiskPosition // 风控位置
}

func NewEvaluatorManager(position RiskPosition) *EvaluatorManager {
	return &EvaluatorManager{items: make([]Item, 0), position: position}
}

// 添加评估项
func (e *EvaluatorManager) AddItem(item Item) {
	e.items = append(e.items, item)
}

// 获取风控的位置点
func (e *EvaluatorManager) GetPosition() RiskPosition {
	return e.position
}

// 计算评估项的总分数
func (e *EvaluatorManager) CalculateScore() (int, map[string]int) {
	total := 0
	itemResult := make(map[string]int)

	for _, i := range e.items {
		score, err := i.Evaluate()
		if err != nil {
			ss_log.Error("RiskPosition:%s,Item:%s, err:%v", e.GetPosition(), i.Name(), err)
			continue
		}
		itemResult[i.Name()] = score
		total += score
	}

	return total, itemResult
}
