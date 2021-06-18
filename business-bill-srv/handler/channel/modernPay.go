package channel

import (
	"a.a/mp-server/business-bill-srv/handler"
	"a.a/mp-server/business-bill-srv/model"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
	"context"
)

type ModernPay struct {
}

func (p *ModernPay) PreCreate(request *model.PreCreateRequest) *model.PreCreateResponse {
	prepayReply := &businessBillProto.PrepayReply{}
	err := handler.BusinessBillHandlerInst.Prepay(context.TODO(), &businessBillProto.PrepayRequest{
		OutOrderNo:   request.OutOrderNo,
		TradeType:    request.TradeType,
		Amount:       request.TradeAmount,
		CurrencyType: request.CurrencyType,
		NotifyUrl:    request.NotifyUrl,
	}, prepayReply)
	if err != nil {

	}

	resp := new(model.PreCreateResponse)
	resp.ReturnCode = ss_err.Success
	resp.ReturnMsg = ss_err.GetMsg(ss_err.Success, "")
	resp.InnerOrderNo = prepayReply.OrderNo
	resp.OutOrderNo = prepayReply.OutOrderNo
	resp.QrCode = prepayReply.QrCodeId
	return resp
}

func (*ModernPay) ScanPay(request *model.ScanPayRequest) *model.ScanPayResponse {
	return nil
}
