package m

//=============================================
type PolyCallbackReq struct {
	RecvMap map[string]interface{} // 接收到的报文
}

type PolyCallbackResp struct {
	RetCode      string
	RetMsg       string
	OrderStatus  string // 订单状态
	RetBody      string // 回复上游报文
	InnerOrderNo string // 平台订单号
	UpperOrderNo string // 上游订单号
	UpdateTime   string
	Amount       string // 发起金额
}

//=============================================
type PolyTransferCallbackReq struct {
	RecvMap map[string]interface{} // 接收到的报文
}

type PolyTransferCallbackResp struct {
	RetCode      string
	RetMsg       string
	OrderStatus  string // 订单状态
	RetBody      string // 回复上游报文
	InnerOrderNo string // 平台订单号
	UpperOrderNo string // 上游订单号
	UpdateTime   string
	Amount       string // 发起金额
}

//=============================================
