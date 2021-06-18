package i

import (
	"a.a/mp-server/api-pay/handler"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_serv"
)

func InitSrv() {
	handler.BusinessBillHandlerInst.Client = businessBillProto.NewBusinessBillService(constants.ServerNameBusinessBill, ss_serv.GetRpcDefClient())
}
