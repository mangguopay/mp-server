package i

import (
	"a.a/mp-server/api-webadmin/handler"
	"a.a/mp-server/common/constants"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	authProto "a.a/mp-server/common/proto/auth"
	billProto "a.a/mp-server/common/proto/bill"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_serv"
)

func InitSrv() {
	handler.AdminAuthHandlerInst.Client = adminAuthProto.NewAdminAuthService(constants.ServerNameAdminAuth, ss_serv.GetRpcDefClient())
	handler.AuthHandlerInst.Client = authProto.NewAuthService(constants.ServerNameAuth, ss_serv.GetRpcDefClient())
	handler.CustHandlerInst.Client = custProto.NewCustService(constants.ServerNameCust, ss_serv.GetRpcDefClient())
	handler.BillHandlerInst.Client = billProto.NewBillService(constants.ServerNameBill, ss_serv.GetRpcDefClient())
	handler.BusinessBillHandlerInst.Client = businessBillProto.NewBusinessBillService(constants.ServerNameBusinessBill, ss_serv.GetRpcDefClient())
}
