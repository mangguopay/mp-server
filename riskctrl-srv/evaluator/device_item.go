package evaluator

func NewDeviceItem(deviceId string) Item {
	return &DeviceItem{deviceId: deviceId}
}

// 设备评估项
type DeviceItem struct {
	deviceId string
}

// 评估项名称
func (d *DeviceItem) Name() string {
	return ItemNameDivice
}

// 对设备进行评估
func (d *DeviceItem) Evaluate() (int, error) {
	// todo
	return 0, nil
}
