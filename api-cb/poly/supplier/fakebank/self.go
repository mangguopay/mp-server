package fakebank

import (
	"a.a/cu/strext"
	"a.a/mp-server/api-cb/m"
	"a.a/mp-server/common/ss_err"
)

var PolyFakebankInst PolyFakebank

type PolyFakebank struct {
}

func (s PolyFakebank) Callback(req *m.PolyCallbackReq) *m.PolyCallbackResp {
	return &m.PolyCallbackResp{
		Amount:       strext.ToStringNoPoint(req.RecvMap["amount"]),
		InnerOrderNo: strext.ToStringNoPoint(req.RecvMap["order_no"]),
		UpperOrderNo: strext.ToStringNoPoint(req.RecvMap["inner_order_no"]),
		UpdateTime:   strext.ToStringNoPoint(req.RecvMap["update_time"]),
		RetCode:      ss_err.ERR_SUCCESS,
		RetMsg:       strext.ToStringNoPoint(req.RecvMap["msg"]),
		OrderStatus:  strext.ToStringNoPoint(req.RecvMap["order_status"]),
		RetBody:      "success",
	}
}

func (s PolyFakebank) TransferCallback(req *m.PolyTransferCallbackReq) *m.PolyTransferCallbackResp {
	return &m.PolyTransferCallbackResp{
		Amount:       strext.ToStringNoPoint(req.RecvMap["amount"]),
		InnerOrderNo: strext.ToStringNoPoint(req.RecvMap["order_no"]),
		UpperOrderNo: strext.ToStringNoPoint(req.RecvMap["inner_order_no"]),
		UpdateTime:   strext.ToStringNoPoint(req.RecvMap["update_time"]),
		RetCode:      ss_err.ERR_SUCCESS,
		RetMsg:       strext.ToStringNoPoint(req.RecvMap["msg"]),
		OrderStatus:  strext.ToStringNoPoint(req.RecvMap["order_status"]),
		RetBody:      "success",
	}
}
