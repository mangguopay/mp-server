package m

type RiskEvaReq struct {
	// 发起金额
	Amount int64
	// 发起支付的账号
	PayerAccNo string
	// 收款人账号
	PayeeAccNo string
	// 创建时间
	ActionTime string
	PayType    string
}

type RiskEvaRsp struct {
	// 评价状态
	RiskExecuteType string
	Score           int64
}

type RiskEvaluatorContext struct {
	// 发起金额
	Amount int64
	// 发起支付的账号
	PayerAccNo string
	// 收款人账号
	PayeeAccNo string
	// 创建时间
	ActionTime string
	// 返回值
	RiskExecuteType string
	Score           int64
	PayType         string
}

type Op struct {
	Type     string `json:"type"`
	Rollback string `json:"rollback"`
	// pipeline
	Steps string `json:"steps"`
	// if
	Condition string `json:"condition"` // op_no
	Y         []*Op  `json:"y"`
	N         []*Op  `json:"n"`
	// switch
	// -- Condition string // op_no
	Ops map[string][]*Op `json:"ops"`
	// for
	IBegin string `json:"i_begin"`
	IStep  string `json:"i_step"`
	IEnd   string `json:"i_end"`
	Loop   []*Op  `json:"loop"`
}

type Setps struct {
	Steps []*Op
}
