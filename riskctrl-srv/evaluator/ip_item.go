package evaluator

func NewIpItem(ip string) Item {
	return &IpItem{ip: ip}
}

// ip评估项
type IpItem struct {
	ip string
}

// 评估项名称
func (i *IpItem) Name() string {
	return ItemNameIp
}

// 对ip进行评估
func (i *IpItem) Evaluate() (int, error) {
	// todo
	return 0, nil
}
