package constants

const (
	// 线上-商户扫码
	ProductNo_ScanCommon = "10001"
	// 线上-个码扫码
	ProductNo_ScanPersonCode = "10002"
	// 线上-随机金额扫码
	ProductNo_ScanRandom = "10003"
	// 线上-商户条码
	ProductNo_CodePayCommon = "10004"
	// 线上-静态二维码支付
	ProductNo_ScanStatic = "10005"
	// 线上-商户池码
	ProductNo_ScanUpPool = "10006"

	// 线上-商户代付
	ProductNo_AgentpayCommon = "20001"
	// 线下-线下出款
	ProductNo_AgentpayOffline = "20002"

	// 线上-商户退款
	ProductNo_RefundCommon = "30001"

	// 线上-快捷支付
	ProductNo_QuickpayCommon = "60001"

	// 线上-代还
	ProductNo_PaybackCommon = "70001"
	ProductNo_PaybackBack   = "70002"

	// 绑卡短信-代还
	ProductNo_BindCardSmsPayback = "80001"
)

const (
	// 扫码
	Scene_Scan = "1"
	// 条码
	Scene_Codepay = "2"
	// 代付
	Scene_Agentpay = "3"
	// 退款
	Scene_Refund = "4"
	// 快捷鉴权短信
	Scene_QuickpaySms = "5"
	// 快捷支付
	Scene_Quickpay = "6"
	// 代还
	Scene_Payback = "7"
	// 代还绑卡短信
	Scene_BindCardSms = "8"
)

const (
	SceneOp_Do          = "1"
	SceneOp_QuerySingle = "2"
	SceneOp_Callback    = "3"
)

const (
	ProductTypePoly     = "poly"     // 聚合
	ProductTypeTransfer = "transfer" // 代付
	ProductTypeWithdraw = "withdraw" // 银行卡
)

//是否可手动签约：0-否，1-是
const (
	ProductIsManualSigned_False = "0" //不可以
	ProductIsManualSigned_True  = "1" //可以
)
