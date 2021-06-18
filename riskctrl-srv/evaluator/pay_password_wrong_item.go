package evaluator

func NewPayPasswordWrongItem(uid string) Item {
	return &PayPasswordWrongItem{uid: uid}
}

// 支付密码错评估项
type PayPasswordWrongItem struct {
	uid string
}

// 评估项名称
func (p *PayPasswordWrongItem) Name() string {
	return ItemNamePayPasswordWrong
}

// 对登录密码错误进行评估
func (p *PayPasswordWrongItem) Evaluate() (int, error) {
	// todo
	return 0, nil
}
