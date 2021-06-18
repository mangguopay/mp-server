package poly

import "a.a/mp-server/api-cb/m"

type Ipoly interface {
	// 支付回调
	Callback(req *m.PolyCallbackReq) *m.PolyCallbackResp
	// 代付回调
	TransferCallback(req *m.PolyTransferCallbackReq) *m.PolyTransferCallbackResp
}
