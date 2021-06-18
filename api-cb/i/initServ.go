package i

import (
	"a.a/mp-server/api-cb/handler"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_serv"
)

func InitSrv() {
	handler.BusinessBillHandlerInst.Client = businessBillProto.NewBusinessBillService("go.micro.srv.business_bill", ss_serv.GetRpcDefClient())
}
