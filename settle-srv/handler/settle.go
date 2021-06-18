package handler

import (
	go_micro_srv_settle "a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_err"
	"context"
)

type Settle struct{}

func (r *Settle) SettleTransfer(ctx context.Context, req *go_micro_srv_settle.SettleTransferRequest, reply *go_micro_srv_settle.SettleTransferReply) error {
	reply.ResultCode = ss_err.ERR_SUCCESS
	return nil
}
