package i

import (
	"a.a/mp-server/common/constants"
	authProto "a.a/mp-server/common/proto/auth"
	billProto "a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/ss_serv"
)

func InitSrv() {
	AuthHandlerInst.Client = authProto.NewAuthService(constants.ServerNameAuth, ss_serv.GetRpcDefClient())
	BillHandlerInst.Client = billProto.NewBillService(constants.ServerNameBill, ss_serv.GetRpcDefClient())
}

type AuthHandler struct {
	Client authProto.AuthService
}

var AuthHandlerInst AuthHandler

type BillHandler struct {
	Client billProto.BillService
}

var BillHandlerInst BillHandler
