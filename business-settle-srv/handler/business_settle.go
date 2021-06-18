package handler

import (
	businessSettleProto "a.a/mp-server/common/proto/business-settle"
	"context"
)

type BusinessSettle struct{}

func (b BusinessSettle) BusinessSettleFees(context.Context, *businessSettleProto.BusinessSettleFeesRequest, *businessSettleProto.BusinessSettleFeesReply) error {
	panic("implement me")
}
