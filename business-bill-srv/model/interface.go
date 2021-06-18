package model

//预创建
type PreCreateRequest struct {
	ReqTime int64 //发起请求时间

	//订单信息
	OutOrderNo   string //外部订单号
	InnerOrderNo string //平台内部订单号
	TradeType    string //交易类型
	TradeAmount  string //支付金额
	CurrencyType string //币种
	Subject      string //商品名称
	NotifyUrl    string //回调地址

	//商户信息
	MerchantNo string                   //商家编号
	SignKey    string                   //签名key
	VerifyKey  string                   //验签key
	ExtendData []map[string]interface{} //扩展数据
}

type PreCreateResponse struct {
	ReturnCode   string //通信结果码
	ReturnMsg    string //通信结果描述
	SubCode      string //业务结果码
	SubMsg       string //业务结果描述
	OutOrderNo   string //外部订单号
	InnerOrderNo string //平台内部订单号
	TradeAmount  string //支付金额
	CurrencyType string //币种
	QrCode       string //二维码
	UpstreamData string //上游返回数据json字符串
}

//==============================================

//支付
type ScanPayRequest struct {
	ReqTime int64 //发起请求时间

	//订单信息
	OutOrderNo   string //外部订单号
	InnerOrderNo string //平台内部订单号
	TradeType    string //交易类型
	TradeAmount  string //交易金额
	CurrencyType string //币种
	Subject      string //商品名称
	NotifyUrl    string //回调地址

	//商户信息
	MerchantNo string                   //商家编号
	SignKey    string                   //签名key
	VerifyKey  string                   //验签key
	ExtendData []map[string]interface{} //扩展数据

}

type ScanPayResponse struct {
	ReturnCode   string //通信结果码
	ReturnMsg    string //通信结果描述
	SubCode      string //业务结果码
	SubMsg       string //业务结果描述
	OutOrderNo   string //外部订单号
	InnerOrderNo string //平台内部订单号
	TradeAmount  string //支付金额
	CurrencyType string //币种
	OrderStatus  string //订单状态
	UpstreamData string //上游返回数据json字符串
}

//================================================
