package poly

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/api-cb/m"
	"a.a/mp-server/common/ss_err"
)

type PolyWrapper struct {
}

var (
	PolyWrapperInst PolyWrapper
)

func (*PolyWrapper) Callback(supplierCode string, req *m.PolyCallbackReq) (resp *m.PolyCallbackResp) {
	targetApi := getTargetApi(supplierCode)
	if targetApi == nil {
		resp = &m.PolyCallbackResp{}
		resp.RetCode = ss_err.ERR_UNKNOW_ERR
		return resp
	}
	// do calling
	ss_log.Info("calling to upstream...[%v]", req)
	resp = targetApi.Callback(req)
	ss_log.Info("resp from upstream...[%v]", resp)
	if resp == nil {
		resp = &m.PolyCallbackResp{}
		resp.RetCode = ss_err.ERR_UNKNOW_ERR
		return resp
	}
	return resp
}

func (*PolyWrapper) TransferCallback(supplierCode string, req *m.PolyTransferCallbackReq) (resp *m.PolyTransferCallbackResp) {
	targetApi := getTargetApi(supplierCode)
	if targetApi == nil {
		resp = &m.PolyTransferCallbackResp{}
		resp.RetCode = ss_err.ERR_UNKNOW_ERR
		return resp
	}
	// do calling
	resp = targetApi.TransferCallback(req)
	if resp == nil {
		resp = &m.PolyTransferCallbackResp{}
		resp.RetCode = ss_err.ERR_UNKNOW_ERR
		return resp
	}
	return resp
}
