package conf

const (
	Port = "9998"

	host = "127.0.0.1"
	//host = "10.41.1.241"

	PrepayUrl                  = "http://" + host + ":8888/api/prepay"
	PayUrl                     = "http://" + host + ":8888/api/pay"
	QueryPay                   = "http://" + host + ":8888/api/query"
	RefundUrl                  = "http://" + host + ":8888/api/refund"
	RefundQueryUrl             = "http://" + host + ":8888/api/query_refund"
	EnterpriseTransferUrl      = "http://" + host + ":8888/api/transfer/enterprise"
	EnterpriseTransferQueryUrl = "http://" + host + ":8888/api/transfer/query"

	// 异步通知地址
	NotifyUrl = "http://" + host + ":" + Port + "/order/notify"
	//NotifyUrl = "http://www.xxx.com/order/notify"

	// 同步跳转地址
	ReturnUrl = "http://" + host + ":" + Port + "/order/jump_back"
)

var (
	// 应用id
	AppId = ""

	// 应用名称
	AppName = ""

	// 商家的私钥
	SelfPrivateKey = ""
	// 平台的公钥
	PlatformPublicKey = ""
)

// 设置配置
func SetConfig(appId, appName, selfPrivateKey, platformPublicKey string) {
	AppId = appId
	AppName = appName
	SelfPrivateKey = selfPrivateKey
	PlatformPublicKey = platformPublicKey
}
