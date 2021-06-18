package i

import (
	"a.a/mp-server/api-mobile/handler"
	"a.a/mp-server/common/constants"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	go_micro_srv_bill "a.a/mp-server/common/proto/bill"
	go_micro_srv_business_bill "a.a/mp-server/common/proto/business-bill"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_serv"
)

func InitSrv() {
	handler.AuthHandlerInst.Client = go_micro_srv_auth.NewAuthService(constants.ServerNameAuth, ss_serv.GetRpcDefClient())
	handler.CustHandlerInst.Client = go_micro_srv_cust.NewCustService(constants.ServerNameCust, ss_serv.GetRpcDefClient())
	handler.BillHandlerInst.Client = go_micro_srv_bill.NewBillService(constants.ServerNameBill, ss_serv.GetRpcDefClient())
	handler.BusinessBillHandlerInst.Client = go_micro_srv_business_bill.NewBusinessBillService(constants.ServerNameBusinessBill, ss_serv.GetRpcDefClient())
}
