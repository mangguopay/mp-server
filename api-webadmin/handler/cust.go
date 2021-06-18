package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"a.a/mp-server/common/constants"

	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/common"
	"a.a/mp-server/api-webadmin/dao"
	"a.a/mp-server/api-webadmin/inner_util"
	"a.a/mp-server/api-webadmin/util"
	"a.a/mp-server/api-webadmin/verify"
	"a.a/mp-server/common/aws_s3"
	"a.a/mp-server/common/global"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/util/file"
)

type CustHandler struct {
	Client custProto.CustService
}

var (
	CustHandlerInst CustHandler
)

/**
获取会员信息
*/
func (*CustHandler) GetCustList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCustList(context.TODO(), &custProto.GetCustListRequest{
				Page:          strext.ToInt32(params[0]),
				PageSize:      strext.ToInt32(params[1]),
				Uid:           strext.ToString(params[2]),
				StartTime:     strext.ToString(params[3]),
				EndTime:       strext.ToString(params[4]),
				QueryNickname: strext.ToString(params[5]),
				QueryPhone:    strext.ToString(params[6]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "uid", "start_time", "end_time", "query_nickname", "query_phone")
	}
}

func (*CustHandler) GetCustInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			uid := strext.ToString(params[0])

			if uid == "" {
				ss_log.Error("uid参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetCustInfo(context.TODO(), &custProto.GetCustRequest{
				Uid: uid,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, gin.H{
				"cust_info":  reply.CustData,
				"card_total": reply.CardTotal,
				"card_datas": reply.CardDatas,
			}, 0, err
		}, "uid")
	}
}

func (*CustHandler) ModifyCustInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyCustInfoRequest{
				CustNo:                   container.GetValFromMapMaybe(params, "cust_no").ToString(),
				InAuthorization:          container.GetValFromMapMaybe(params, "in_authorization").ToString(),
				OutAuthorization:         container.GetValFromMapMaybe(params, "out_authorization").ToString(),
				InTransferAuthorization:  container.GetValFromMapMaybe(params, "in_transfer_authorization").ToString(),
				OutTransferAuthorization: container.GetValFromMapMaybe(params, "out_transfer_authorization").ToString(),
				LoginUid:                 inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.CheckModifyCustInfoVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := CustHandlerInst.Client.ModifyCustInfo(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

//用户账单
func (*CustHandler) GetCustBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetCustBillsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				StartTime:   strext.ToString(params[2]),
				EndTime:     strext.ToString(params[3]),
				Reason:      strext.ToString(params[4]),
				Uid:         strext.ToString(params[5]),
				BizLogNo:    strext.ToString(params[6]),
				BalanceType: strext.ToString(params[7]),
			}

			if req.Uid == "" {
				ss_log.Error("参数Uid为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetCustBills(context.TODO(), req)

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "reason", "uid", "log_no", "balance_type")
	}

}

//获取ModernPay服务商列表
func (*CustHandler) GetServicerList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServicerList(context.TODO(), &custProto.GetServicerListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToString(params[2]),
				EndTime:      strext.ToString(params[3]),
				ServicerNo:   strext.ToString(params[4]),
				QueryName:    strext.ToString(params[5]),
				QueryPhone:   strext.ToString(params[6]),
				Account:      strext.ToString(params[7]),
				ServicerName: strext.ToString(params[8]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "servicer_no", "query_name", "query_phone", "account", "servicer_name")
	}
}

func (*CustHandler) GetServicerPhoneList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServicerPhoneList(context.TODO(), &custProto.GetServicerPhoneListRequest{
				QueryPhone: strext.ToString(params[0]),
			})
			return reply.ResultCode, reply.DataList, 0, err
		}, "query_phone")
	}
}

//获取ModernPay指定商户信息
func (*CustHandler) GetServicerInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			servicerNo := strext.ToString(params[0])
			if servicerNo == "" {
				ss_log.Error("参数ServicerNo为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetServicerInfo(context.TODO(), &custProto.GetServicerInfoRequest{
				ServicerNo: servicerNo,
			})
			return reply.ResultCode, gin.H{
				"info":                   reply.Data,
				"card_pack_data":         reply.CardList,
				"servicer_img_data":      reply.ServicerImgData,
				"servicer_terminal_data": reply.ServicerTerminalDataList,
			}, 0, err
		}, "servicer_no")
	}
}

//获取服务商交易查询信息列表
func (*CustHandler) GetServiceTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServiceTransactions(context.TODO(), &custProto.GetServiceTransactionsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToString(params[2]),
				EndTime:      strext.ToString(params[3]),
				Nickname:     strext.ToString(params[4]),
				Phone:        strext.ToString(params[5]),
				CurrencyType: strext.ToString(params[6]),
				OrderNo:      strext.ToString(params[7]),
				BillType:     strext.ToString(params[8]),
				Account:      strext.ToString(params[9]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "nickname", "phone", "currency_type", "order_no", "bill_type", "account")
	}
}

//获取服务商收益明细查询列表
func (*CustHandler) GetServicerProfitLedgerList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServicerProfitLedgerList(context.TODO(), &custProto.GetServicerProfitLedgerListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				Nickname:     strext.ToString(params[2]),
				Phone:        strext.ToString(params[3]),
				CurrencyType: strext.ToString(params[4]),
				LogNo:        strext.ToString(params[5]),
				StartTime:    strext.ToString(params[6]),
				EndTime:      strext.ToString(params[7]),
				OrderType:    strext.ToString(params[8]),
				Account:      strext.ToString(params[9]),
			})
			return reply.ResultCode, gin.H{
				"datas":      reply.DataList,
				"count_data": reply.CountData,
			}, reply.Total, err
		}, "page", "page_size", "nickname", "phone", "currency_type", "log_no", "start_time", "end_time", "order_type", "account")
	}
}

//修改ModernPay商户状态(停用)
func (*CustHandler) ModifyServicerStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyServicerStatusRequest{
				ServicerNo: container.GetValFromMapMaybe(params, "servicer_no").ToString(),
				UseStatus:  container.GetValFromMapMaybe(params, "service_status").ToString(),
				LoginUid:   inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.CheckModifyServicerStatusVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.ModifyServicerStatus(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

//修改ModernPay指定商户配置
func (*CustHandler) ModifyServicerConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyServicerConfigRequest{
				ServicerNo:          container.GetValFromMapMaybe(params, "servicer_no").ToString(),
				IncomeAuthorization: container.GetValFromMapMaybe(params, "income_authorization").ToString(),
				OutgoAuthorization:  container.GetValFromMapMaybe(params, "outgo_authorization").ToString(),
				CommissionSharing:   container.GetValFromMapMaybe(params, "commission_sharing").ToString(),
				IncomeSharing:       container.GetValFromMapMaybe(params, "income_sharing").ToString(),
				UsdAuthCollectLimit: container.GetValFromMapMaybe(params, "usd_auth_collect_limit").ToString(),
				KhrAuthCollectLimit: container.GetValFromMapMaybe(params, "khr_auth_collect_limit").ToString(),
				Lat:                 container.GetValFromMapMaybe(params, "lat").ToString(),
				Lng:                 container.GetValFromMapMaybe(params, "lng").ToString(),
				Scope:               container.GetValFromMapMaybe(params, "scope").ToString(),
				ScopeOff:            container.GetValFromMapMaybe(params, "scope_off").ToString(),
				ServicerName:        container.GetValFromMapMaybe(params, "servicer_name").ToString(),
				BusinessTime:        container.GetValFromMapMaybe(params, "business_time").ToString(),
				Addr:                container.GetValFromMapMaybe(params, "addr").ToString(),
				LoginUid:            inner_util.GetJwtDataString(c, "account_uid"),
				//TerminalNumberSn:    container.GetValFromMapMaybe(params, "terminal_number_sn").ToString(),
			}

			if errStr := verify.ModifyServicerConfigReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.ModifyServicerConfig(context.TODO(), req)
			return reply.ResultCode, "", err
		})
	}
}

//获取服务商收取款统计信息
func (*CustHandler) GetServicerOrderCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServicerOrderCount(context.TODO(), &custProto.GetServicerOrderCountRequest{
				ServicerNo: strext.ToString(params[0]),
			})
			return reply.ResultCode, gin.H{
				"usdData": reply.UsdData,
				"khrData": reply.KhrData,
			}, 0, err
		}, "servicer_no")
	}
}

//获取每天服务商的对账统计
func (*CustHandler) GetFinancialServicerCheckList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetFinancialServicerCheckList(context.TODO(), &custProto.GetFinancialServicerCheckListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToString(params[2]),
				EndTime:      strext.ToString(params[3]),
				ServicerNo:   strext.ToString(params[4]),
				Phone:        strext.ToString(params[5]),
				Account:      strext.ToString(params[6]),
				Id:           strext.ToString(params[7]),
				CurrencyType: strext.ToString(params[8]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "servicer_no", "phone", "account", "id", "currency_type")
	}
}

//查看指定服务商的某天账单明细
func (*CustHandler) GetBillingDetailsResultsList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBillingDetailsResultsList(context.TODO(), &custProto.GetBillingDetailsResultsListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				OrderNo:      strext.ToString(params[2]),
				BillType:     strext.ToString(params[3]),
				StartTime:    strext.ToString(params[4]),
				EndTime:      strext.ToString(params[5]),
				ServicerNo:   strext.ToString(params[6]),
				CurrencyType: strext.ToString(params[7]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "order_no", "bill_type", "start_time", "end_time", "servicer_no", "currency_type")
	}
}

func (*CustHandler) ModifyCollectStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyCollectStatusRequest{
				SetStatus: container.GetValFromMapMaybe(params, "set_status").ToString(),
				CardNo:    container.GetValFromMapMaybe(params, "card_no").ToString(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.ModifyCollectStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			var reply, err = CustHandlerInst.Client.ModifyCollectStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

//收款账户管理
func (*CustHandler) CollectionManagementList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.CollectionManagementList(context.TODO(), &custProto.CollectionManagementListRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				AccountType: strext.ToStringNoPoint(params[2]), //3服务商4用户 7个人商家，8企业商家
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "account_type")
	}
}

func (*CustHandler) DelectCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoDelete(c, func(params interface{}) (string, error) {
			req := &custProto.DelectCardRequest{
				CardNo:   container.GetValFromMapMaybe(params, "card_no").ToString(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.DelectCardVerify(req); errStr != "" {
				return errStr, nil
			}

			reply, err := CustHandlerInst.Client.DelectCard(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, err
			}
			return strext.ToString(reply.ResultCode), err
		})
	}
}

func (*CustHandler) GetCardlInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCardInfo(context.TODO(), &custProto.GetCardInfoRequest{
				CardNo: strext.ToString(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "card_no")
	}
}

func (*CustHandler) UpdateOrInsertCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateOrInsertCardRequest{
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToString(),
				ChannelNo:   container.GetValFromMapMaybe(params, "channel_no").ToString(),
				Name:        container.GetValFromMapMaybe(params, "name").ToString(),
				CardNumber:  container.GetValFromMapMaybe(params, "card_number").ToString(),
				Note:        container.GetValFromMapMaybe(params, "note").ToString(),
				BalanceType: container.GetValFromMapMaybe(params, "balance_type").ToString(),
				IsDefalut:   container.GetValFromMapMaybe(params, "is_defalut").ToString(),
				AccountType: container.GetValFromMapMaybe(params, "account_type").ToString(),
				ChannelId:   container.GetValFromMapMaybe(params, "channel_id").ToString(),
				LoginUid:    inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.UpdateOrInsertCardReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			var reply, err = CustHandlerInst.Client.UpdateOrInsertCard(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetHeadquartersProfitList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetHeadquartersProfitList(context.TODO(), &custProto.GetHeadquartersProfitListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				OrderNo:      strext.ToString(params[2]),
				StartTime:    strext.ToString(params[3]),
				EndTime:      strext.ToString(params[4]),
				ProfitSource: strext.ToString(params[5]),
				BalanceType:  strext.ToString(params[6]),
			})
			return reply.ResultCode, gin.H{
				"datas":     reply.DataList,
				"countData": reply.CountData,
			}, reply.Total, err
		}, "page", "page_size", "order_no", "start_time", "end_time", "profit_source", "balance_type")
	}
}

func (*CustHandler) GetToHeadquartersList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetToHeadquartersList(context.TODO(), &custProto.GetToHeadquartersListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				LogNo:        strext.ToString(params[2]),
				StartTime:    strext.ToString(params[3]),
				EndTime:      strext.ToString(params[4]),
				Account:      strext.ToString(params[5]),
				OrderStatus:  strext.ToString(params[6]),
				OrderType:    strext.ToString(params[7]),
				CurrencyType: strext.ToString(params[8]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "account", "order_status", "order_type", "currency_type")
	}
}

func (*CustHandler) UpdateToHeadquarters() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			var reply, err = CustHandlerInst.Client.UpdateToHeadquarters(context.TODO(), &custProto.UpdateToHeadquartersRequest{
				LogNo:       container.GetValFromMapMaybe(params, "log_no").ToString(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToString(),
				LoginUid:    inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetToServicerList() gin.HandlerFunc {

	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetToServicerList(context.TODO(), &custProto.GetToServicerListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				LogNo:        strext.ToString(params[2]),
				StartTime:    strext.ToString(params[3]),
				EndTime:      strext.ToString(params[4]),
				Account:      strext.ToString(params[5]),
				OrderStatus:  strext.ToString(params[6]),
				OrderType:    strext.ToString(params[7]),
				CurrencyType: strext.ToString(params[8]),
				Nickname:     strext.ToString(params[9]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "account", "order_status", "order_type", "currency_type", "nickname")
	}
}

func (*CustHandler) AddToServicer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.AddToServicer(context.TODO(), &custProto.AddToServicerRequest{
				CurrencyType:   container.GetValFromMapMaybe(params, "currency_type").ToString(),
				ServicerNo:     container.GetValFromMapMaybe(params, "servicer_no").ToString(),
				CollectionType: container.GetValFromMapMaybe(params, "collection_type").ToString(),
				CardNo:         container.GetValFromMapMaybe(params, "card_no").ToString(),
				Amount:         container.GetValFromMapMaybe(params, "amount").ToString(),
				OrderType:      container.GetValFromMapMaybe(params, "order_type").ToString(),
				OrderStatus:    container.GetValFromMapMaybe(params, "order_status").ToString(),
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetTransferOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetTransferOrderList(context.TODO(), &custProto.GetTransferOrderListRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				LogNo:       strext.ToString(params[2]),
				StartTime:   strext.ToString(params[3]),
				EndTime:     strext.ToString(params[4]),
				FromAccount: strext.ToString(params[5]),
				ToAccount:   strext.ToString(params[6]),
				OrderStatus: strext.ToString(params[7]),
				BalanceType: strext.ToString(params[8]),
				WriteOff:    strext.ToString(params[9]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "from_account", "to_account", "order_status", "balance_type", "write_off")
	}
}

func (*CustHandler) GetOutgoOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetOutgoOrderList(context.TODO(), &custProto.GetOutgoOrderListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				LogNo:        strext.ToString(params[2]),
				StartTime:    strext.ToString(params[3]),
				EndTime:      strext.ToString(params[4]),
				Phone:        strext.ToString(params[5]),
				OrderStatus:  strext.ToString(params[6]),
				BalanceType:  strext.ToString(params[7]),
				Account:      strext.ToString(params[8]),
				WriteOff:     strext.ToString(params[9]),
				OutAccount:   strext.ToString(params[10]),
				WithdrawType: strext.ToString(params[11]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "phone", "order_status", "balance_type", "account", "write_off", "out_account", "withdraw_type")
	}
}

func (*CustHandler) GetIncomeOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetIncomeOrderList(context.TODO(), &custProto.GetIncomeOrderListRequest{
				Page:          strext.ToInt32(params[0]),
				PageSize:      strext.ToInt32(params[1]),
				LogNo:         strext.ToString(params[2]),
				StartTime:     strext.ToString(params[3]),
				EndTime:       strext.ToString(params[4]),
				OrderStatus:   strext.ToString(params[5]),
				BalanceType:   strext.ToString(params[6]),
				IncomePhone:   strext.ToString(params[7]),
				Account:       strext.ToString(params[8]),
				RecvPhone:     strext.ToString(params[9]),
				WriteOff:      strext.ToString(params[10]),
				IncomeAccount: strext.ToString(params[11]),
				RecvAccount:   strext.ToString(params[12]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "order_status", "balance_type", "income_phone", "account", "recv_phone", "write_off", "income_account", "recv_account")
	}
}

func (*CustHandler) GetExchangeOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetExchangeOrderList(context.TODO(), &custProto.GetExchangeOrderListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				LogNo:        strext.ToString(params[2]),
				StartTime:    strext.ToString(params[3]),
				EndTime:      strext.ToString(params[4]),
				Phone:        strext.ToString(params[5]),
				OrderStatus:  strext.ToString(params[6]),
				ExchangeType: strext.ToString(params[7]),
				Account:      strext.ToString(params[8]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "phone", "order_status", "exchange_type", "account")
	}
}

func (*CustHandler) GetCollectionOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCollectionOrders(context.TODO(), &custProto.GetCollectionOrdersRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				LogNo:       strext.ToString(params[2]),
				StartTime:   strext.ToString(params[3]),
				EndTime:     strext.ToString(params[4]),
				FromAccount: strext.ToString(params[5]),
				ToAccount:   strext.ToString(params[6]),
				OrderStatus: strext.ToString(params[7]),
				BalanceType: strext.ToString(params[8]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "from_account", "to_account", "order_status", "balance_type")
	}
}

func (*CustHandler) GetChannels() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetChannels(context.TODO(), &custProto.GetChannelsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				ChannelName: strext.ToStringNoPoint(params[2]),
				ChannelType: strext.ToStringNoPoint(params[3]),
				UseStatus:   strext.ToStringNoPoint(params[4]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "channel_name", "channel_type", "use_status")
	}
}

func (*CustHandler) HeadquartersProfitWithdraws() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.HeadquartersProfitWithdraws(context.TODO(), &custProto.HeadquartersProfitWithdrawsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				OrderNo:      strext.ToString(params[2]),
				StartTime:    strext.ToString(params[3]),
				EndTime:      strext.ToString(params[4]),
				CurrencyType: strext.ToString(params[5]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "order_no", "start_time", "end_time", "currency_type")
	}
}

func (*CustHandler) GetLogVaccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetLogVaccounts(context.TODO(), &custProto.GetLogVaccountsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				LogNo:       strext.ToString(params[2]),
				StartTime:   strext.ToString(params[3]),
				EndTime:     strext.ToString(params[4]),
				OpType:      strext.ToString(params[5]),
				BalanceType: strext.ToString(params[6]),
				BizLogNo:    strext.ToString(params[7]),
				Reason:      strext.ToString(params[8]),
				Account:     strext.ToString(params[9]),
				VaType:      strext.ToString(params[10]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "op_type", "balance_type", "biz_log_no", "reason", "account", "va_type")
	}
}

func (*CustHandler) UpdateServiceWithdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			var reply, err = CustHandlerInst.Client.UpdateServiceWithdraw(context.TODO(), &custProto.UpdateServiceWithdrawRequest{
				OrderNo:  container.GetValFromMapMaybe(params, "log_no").ToString(),
				Status:   container.GetValFromMapMaybe(params, "order_status").ToInt32(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, 0, err
		})
	}
}

// 修改用户向总部提现的订单状态
func (*CustHandler) UpdateCustWithdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			var reply, err = CustHandlerInst.Client.UpdateCustWithdraw(context.TODO(), &custProto.UpdateCustWithdrawRequest{
				OrderNo:     container.GetValFromMapMaybe(params, "log_no").ToString(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToInt32(),
				Notes:       container.GetValFromMapMaybe(params, "notes").ToString(),
				ImageBase64: container.GetValFromMapMaybe(params, "base64_img").ToString(),
				LoginUid:    inner_util.GetJwtDataString(c, "account_uid"),
			})

			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

// 修改用户向总部存款的订单状态
func (*CustHandler) UpdateCustSave() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			var reply, err = CustHandlerInst.Client.UpdateCustSave(context.TODO(), &custProto.UpdateCustSaveRequest{
				OrderNo:     container.GetValFromMapMaybe(params, "log_no").ToString(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToInt32(),
				LoginUid:    inner_util.GetJwtDataString(c, "account_uid"),
			}, global.RequestTimeoutOptions)

			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetServicerAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServicerAccounts(context.TODO(), &custProto.GetServicerAccountsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Account:  strext.ToString(params[2]),
				SortType: strext.ToString(params[3]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "account", "sort_type")
	}
}

//用户账单
func (*CustHandler) GetSerBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetServicerBillsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				StartTime:   strext.ToString(params[2]),
				EndTime:     strext.ToString(params[3]),
				Reason:      strext.ToString(params[4]),
				Uid:         strext.ToString(params[5]),
				BizLogNo:    strext.ToString(params[6]),
				BalanceType: strext.ToString(params[7]),
			}

			if req.Uid == "" {
				ss_log.Error("参数Uid为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetServicerBills(context.TODO(), req)

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "reason", "uid", "log_no", "balance_type")
	}

}

func (s *CustHandler) GetCommonHelps() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCommonHelps(context.TODO(), &custProto.GetCommonHelpsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Problem:  strext.ToStringNoPoint(params[2]),
				Lang:     strext.ToStringNoPoint(params[3]),
				VsType:   strext.ToStringNoPoint(params[4]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "problem", "lang", "vs_type")

	}
}

func (s *CustHandler) GetCommonHelpDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCommonHelpDetail(context.TODO(), &custProto.GetCommonHelpDetailRequest{
				HelpNo: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "help_no")

	}
}

func (s *CustHandler) InsertOrUpdateHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateHelpRequest{
				HelpNo:  container.GetValFromMapMaybe(params, "help_no").ToStringNoPoint(),
				Problem: container.GetValFromMapMaybe(params, "problem").ToStringNoPoint(),
				Answer:  container.GetValFromMapMaybe(params, "answer").ToStringNoPoint(),
				Lang:    container.GetValFromMapMaybe(params, "lang").ToStringNoPoint(),
				VsType:  container.GetValFromMapMaybe(params, "vs_type").ToStringNoPoint(),
			}
			if errStr := verify.CheckInsertOrUpdateHelpVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.InsertOrUpdateHelp(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteHelp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			var reply, err = CustHandlerInst.Client.DeleteHelp(context.TODO(), &custProto.DeleteHelpRequest{
				HelpNo: container.GetValFromMapMaybe(params, "help_no").ToStringNoPoint(),
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) ModifyHelpStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			var reply, err = CustHandlerInst.Client.ModifyHelpStatus(context.TODO(), &custProto.ModifyHelpStatusRequest{
				HelpNo:    container.GetValFromMapMaybe(params, "help_no").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
			})
			return reply.ResultCode, 0, err
		})

	}
}

func (s *CustHandler) SwapHelpIdx() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, err := CustHandlerInst.Client.SwapHelpIdx(context.TODO(), &custProto.SwapHelpIdxRequest{
				Idx:      container.GetValFromMapMaybe(params, "idx").ToStringNoPoint(),
				HelpNo:   container.GetValFromMapMaybe(params, "help_no").ToStringNoPoint(),
				Lang:     container.GetValFromMapMaybe(params, "lang").ToStringNoPoint(),
				SwapType: container.GetValFromMapMaybe(params, "swap_type").ToStringNoPoint(),
				VsType:   container.GetValFromMapMaybe(params, "vs_type").ToStringNoPoint(),
			})
			if err != nil {
				ss_log.Error("SwapHelpIdx err=[%v]", err)
				return response.ResultCode, "", nil
			}
			return response.ResultCode, "", err
		})
	}
}

func (s *CustHandler) GetConsultationConfigs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetConsultationConfigs(context.TODO(), &custProto.GetConsultationConfigsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				UseStatus: strext.ToStringNoPoint(params[2]),
				Lang:      strext.ToStringNoPoint(params[3]),
				Name:      strext.ToStringNoPoint(params[4]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "use_status", "lang", "name")

	}
}

func (s *CustHandler) GetConsultationConfigDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetConsultationConfigDetail(context.TODO(), &custProto.GetConsultationConfigDetailRequest{
				Id: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "id")

	}
}

func (s *CustHandler) InsertOrUpdateConsultationConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			//上传图片
			imageStr := container.GetValFromMapMaybe(params, "logo_img_base64").ToStringNoPoint()
			if len(imageStr) > constants.UploadImgBase64LengthMax {
				return ss_err.ERR_ACCOUNT_IMAGE_BIG, nil, nil
			}
			ss_log.Info("imageStr:%v", imageStr)
			req := &custProto.UploadImageRequest{
				ImageStr:   imageStr,
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
				Type:       constants.UploadImage_UnAuth, //类型是不需要授权的图片
			}

			if errStr := verify.UploadImageReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.UploadImage(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用保存图片服务失败")
				return ss_err.ERR_SAVE_IMAGE_FAILD, "", nil
			}

			req2 := &custProto.InsertOrUpdateConsultationConfigRequest{
				Id:        container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				Name:      container.GetValFromMapMaybe(params, "name").ToStringNoPoint(),
				Text:      container.GetValFromMapMaybe(params, "text").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				Lang:      container.GetValFromMapMaybe(params, "lang").ToStringNoPoint(),
				LogoImgNo: reply.ImageId,
			}

			if errStr := verify.InsertOrUpdateConsultationConfigVerify(req2); errStr != "" {
				return errStr, nil, nil
			}

			reply2, err2 := CustHandlerInst.Client.InsertOrUpdateConsultationConfig(context.TODO(), req2)
			return reply2.ResultCode, 0, err2
		})
	}
}

func (s *CustHandler) DeleteConsultationConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			id := container.GetValFromMapMaybe(params, "id").ToStringNoPoint()
			loginUid := inner_util.GetJwtDataString(c, "account_uid")
			if id == "" {
				ss_log.Error("参数id为空")
				return ss_err.ERR_PARAM, 0, nil
			}
			if loginUid == "" {
				ss_log.Error("参数loginUid为空")
				return ss_err.ERR_PARAM, 0, nil
			}

			var reply, err = CustHandlerInst.Client.DeleteConsultationConfig(context.TODO(), &custProto.DeleteConsultationConfigRequest{
				Id:       id,
				LoginUid: loginUid,
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) ModifyConsultationConfigStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			var reply, err = CustHandlerInst.Client.ModifyConsultationConfigStatus(context.TODO(), &custProto.ModifyConsultationConfigStatusRequest{
				Id:        container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, 0, err
		})

	}
}

//
func (s *CustHandler) GetAgreements() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAgreements(context.TODO(), &custProto.GetAgreementsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Type:     strext.ToStringNoPoint(params[2]),
				Lang:     strext.ToStringNoPoint(params[3]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "type", "lang")
	}
}

func (s *CustHandler) GetAgreementDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			id := strext.ToStringNoPoint(params[0])
			if id == "" {
				ss_log.Error("参数id为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}
			reply, err := CustHandlerInst.Client.GetAgreementDetail(context.TODO(), &custProto.GetAgreementDetailRequest{
				Id: id,
			})

			return reply.ResultCode, reply.Data, 0, err
		}, "id")

	}
}

func (s *CustHandler) InsertOrUpdateAgreement() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {

			req := &custProto.InsertOrUpdateAgreementRequest{
				Id:        container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				Text:      container.GetValFromMapMaybe(params, "text").ToStringNoPoint(),
				Lang:      container.GetValFromMapMaybe(params, "lang").ToStringNoPoint(),
				Type:      container.GetValFromMapMaybe(params, "type").ToStringNoPoint(), //0 用户协议  1隐私协议 2实名认证协议
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.InsertOrUpdateAgreementVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			var reply, err = CustHandlerInst.Client.InsertOrUpdateAgreement(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteAgreement() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteAgreementRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.DeleteAgreementVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			var reply, err = CustHandlerInst.Client.DeleteAgreement(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) ModifyAgreementStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyAgreementStatusRequest{
				Id:        container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.ModifyAgreementStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			if req.UseStatus == "0" { //不可修改使用状态为不使用，已预防使用中的变为不使用
				ss_log.Error("不可修改使用状态为不使用,UseStatus[%v]", req.UseStatus)
				return ss_err.ERR_PARAM, 0, nil
			}

			var reply, err = CustHandlerInst.Client.ModifyAgreementStatus(context.TODO(), req)
			return reply.ResultCode, 0, err
		})

	}
}

func (*CustHandler) CreateXlsxFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			page := strext.ToInt32(1)
			pageSize := strext.ToInt32(10000)

			exchangeReq := &custProto.GetExchangeOrderListRequest{
				Page:         page,
				PageSize:     pageSize,
				LogNo:        container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:    container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:      container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				Phone:        container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				OrderStatus:  container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				ExchangeType: container.GetValFromMapMaybe(params, "exchange_type").ToStringNoPoint(),
			}
			incomeReq := &custProto.GetIncomeOrderListRequest{
				Page:        page,
				PageSize:    pageSize,
				LogNo:       container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:   container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:     container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				BalanceType: container.GetValFromMapMaybe(params, "balance_type").ToStringNoPoint(),
				IncomePhone: container.GetValFromMapMaybe(params, "income_phone").ToStringNoPoint(),
				Account:     container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				RecvPhone:   container.GetValFromMapMaybe(params, "recv_phone").ToStringNoPoint(),
			}
			outgoReq := &custProto.GetOutgoOrderListRequest{
				Page:        page,
				PageSize:    pageSize,
				LogNo:       container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:   container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:     container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				Phone:       container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				BalanceType: container.GetValFromMapMaybe(params, "balance_type").ToStringNoPoint(),
				Account:     container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
			}
			transferReq := &custProto.GetTransferOrderListRequest{
				Page:        page,
				PageSize:    pageSize,
				LogNo:       container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:   container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:     container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				FromAccount: container.GetValFromMapMaybe(params, "from_account").ToStringNoPoint(),
				ToAccount:   container.GetValFromMapMaybe(params, "to_account").ToStringNoPoint(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				BalanceType: container.GetValFromMapMaybe(params, "balance_type").ToStringNoPoint(),
			}
			collectionReq := &custProto.GetCollectionOrdersRequest{
				Page:        page,
				PageSize:    pageSize,
				LogNo:       container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:   container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:     container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				FromAccount: container.GetValFromMapMaybe(params, "from_account").ToStringNoPoint(),
				ToAccount:   container.GetValFromMapMaybe(params, "to_account").ToStringNoPoint(),
				OrderStatus: container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				BalanceType: container.GetValFromMapMaybe(params, "balance_type").ToStringNoPoint(),
			}
			toHeadReq := &custProto.GetToHeadquartersListRequest{
				Page:         page,
				PageSize:     pageSize,
				LogNo:        container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:    container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:      container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				Account:      container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				OrderStatus:  container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				OrderType:    container.GetValFromMapMaybe(params, "order_type").ToStringNoPoint(),
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
			}
			toSerReq := &custProto.GetToServicerListRequest{
				Page:         page,
				PageSize:     pageSize,
				LogNo:        container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:    container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:      container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				Account:      container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				OrderStatus:  container.GetValFromMapMaybe(params, "order_status").ToStringNoPoint(),
				OrderType:    container.GetValFromMapMaybe(params, "order_type").ToStringNoPoint(),
				CurrencyType: container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				Nickname:     container.GetValFromMapMaybe(params, "nickname").ToStringNoPoint(),
			}
			logVaccReq := &custProto.GetLogVaccountsRequest{
				Page:        page,
				PageSize:    pageSize,
				LogNo:       container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
				StartTime:   container.GetValFromMapMaybe(params, "start_time").ToStringNoPoint(),
				EndTime:     container.GetValFromMapMaybe(params, "end_time").ToStringNoPoint(),
				OpType:      container.GetValFromMapMaybe(params, "op_type").ToStringNoPoint(),
				BalanceType: container.GetValFromMapMaybe(params, "balance_type").ToStringNoPoint(),
				BizLogNo:    container.GetValFromMapMaybe(params, "biz_log_no").ToStringNoPoint(),
				Reason:      container.GetValFromMapMaybe(params, "reason").ToStringNoPoint(),
			}

			req := &custProto.GetCreateXlsxFileContentRequest{
				//订单类型（1兑 2存 3取 4转 5收 6服务商充值 7服务商请款 8虚拟账户日志流水）
				BillFileType:  container.GetValFromMapMaybe(params, "bill_file_type").ToStringNoPoint(),
				AccountNo:     inner_util.GetJwtDataString(c, "account_uid"),  //
				RoleType:      inner_util.GetJwtDataString(c, "account_type"), //
				ExchangeReq:   exchangeReq,
				IncomeReq:     incomeReq,
				OutgoReq:      outgoReq,
				TransferReq:   transferReq,
				CollectionReq: collectionReq,
				ToHeadReq:     toHeadReq,
				ToSerReq:      toSerReq,
				LogVaccReq:    logVaccReq,
			}

			//取出处理过的数据
			reply, err := CustHandlerInst.Client.GetCreateXlsxFileContent(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("查询数据失败,ResultCode=[%v]", reply.ResultCode)
				return reply.ResultCode, "", nil
			}

			xlsxTaskNo := strext.GetDailyId()
			pathStr := os.TempDir()

			fileName := fmt.Sprintf("%v/%v.xlsx", pathStr, xlsxTaskNo)
			ss_log.Info("fileName:[%v]", fileName)
			errC := util.CreateXlsxFile(fileName, reply.BillType, reply.QueryStr, reply.Datas)
			if errC != nil {
				ss_log.Error("创建临时文件出错,err[%v]", errC)
				return ss_err.ERR_PARAM, nil, nil
			}

			//读取并返回数据给前端
			file, err := os.Open(fileName)
			defer file.Close()
			if err != nil {
				ss_log.Error("err=[%v]\n", err)
				c.Set(ss_net.RET_CODE, ss_err.ERR_PARAM)
				return ss_err.ERR_PARAM, nil, nil
			} else {
				buff, err := ioutil.ReadAll(file)
				if err != nil {
					ss_log.Error("err=[%v]\n", err)
					c.Set(ss_net.RET_CODE, ss_err.ERR_PARAM)
					return ss_err.ERR_PARAM, nil, nil
				}
				_, err2 := c.Writer.Write(buff)
				if err2 != nil {
					ss_log.Error("err=[%v]\n", err2)
					c.Set(ss_net.RET_CODE, ss_err.ERR_PARAM)
					return ss_err.ERR_PARAM, nil, nil
				}

				ss_log.Info("成功返回文件数据")

				file.Close()
				//删除临时文件
				errDel := os.Remove(fileName)
				if errDel != nil {
					// 删除失败
					ss_log.Error("删除临时文件失败,err=[%v],fileName=[%v]", errDel, fileName)
				} else {
					ss_log.Info("已删除临时文件fileName[%v]", fileName)
				}

				c.Set(ss_net.RET_CODE, ss_err.ERR_SUCCESS)
				return ss_err.ERR_SUCCESS, nil, nil
			}

			//return ss_err.ERR_SUCCESS, gin.H{
			//	"xlsx_task_no": reply.XlsxTaskNo,
			//}, nil
		})

	}
}

func (*CustHandler) ModifyServicerInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyServicerInfoRequest{
				//ServicerBase64Img: container.GetValFromMapMaybe(params, "servicer_base64_img").ToString(), //营业执照图片
				//ImgType:           container.GetValFromMapMaybe(params, "img_type").ToString(),            //1营业执照图片,234营业场所
				AccountNo:    container.GetValFromMapMaybe(params, "account_no").ToString(),    //
				ServicerNo:   container.GetValFromMapMaybe(params, "servicer_no").ToString(),   //
				ServicerImg1: container.GetValFromMapMaybe(params, "servicer_img1").ToString(), //
				ServicerImg2: container.GetValFromMapMaybe(params, "servicer_img2").ToString(), //
				ServicerImg3: container.GetValFromMapMaybe(params, "servicer_img3").ToString(), //
				ServicerImg4: container.GetValFromMapMaybe(params, "servicer_img4").ToString(), //
				Nickname:     container.GetValFromMapMaybe(params, "nickname").ToString(),      //
				Phone:        container.GetValFromMapMaybe(params, "phone").ToString(),         //
				Addr:         container.GetValFromMapMaybe(params, "addr").ToString(),          //
				CountryCode:  "855",                                                            //
				LoginUid:     inner_util.GetJwtDataString(c, "account_uid"),                    //
			}
			//
			if req.LoginUid == "" {
				ss_log.Error("登陆的账号uid获取不到")
				return ss_err.ERR_ACCOUNT_JWT_OUTDATED, nil, nil
			}

			reply, err := CustHandlerInst.Client.ModifyServicerInfo(context.TODO(), req, global.RequestTimeoutOptions)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetAppVersions() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAppVersions(context.TODO(), &custProto.GetAppVersionsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				VsType:   strext.ToStringNoPoint(params[2]),
				System:   strext.ToStringNoPoint(params[3]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "vs_type", "system")
	}
}

func (s *CustHandler) GetAppVersionsCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAppVersionsCount(context.TODO(), &custProto.GetAppVersionsCountRequest{})
			//return reply.ResultCode, gin.H{
			//	"ios_app_count_data":     reply.IosAppCountData,
			//	"ios_pos_count_data":     reply.IosPosCountData,
			//	"android_app_count_data": reply.AndroidAppCountData,
			//	"android_pos_count_data": reply.AndroidPosCountData,
			//}, 0, err
			return reply.ResultCode, reply.CountData, 0, err
		}, "")
	}
}

func (s *CustHandler) GetAppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAppVersion(context.TODO(), &custProto.GetAppVersionRequest{
				VId: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "v_id")
	}
}

func (s *CustHandler) ModifyAppVersionStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.ModifyAppVersionStatus(context.TODO(), &custProto.ModifyAppVersionStatusRequest{
				VId:      container.GetValFromMapMaybe(params, "v_id").ToStringNoPoint(),
				Status:   container.GetValFromMapMaybe(params, "status").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) ModifyAppVersionIsForce() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.ModifyAppVersionIsForce(context.TODO(), &custProto.ModifyAppVersionIsForceRequest{
				VId:      container.GetValFromMapMaybe(params, "v_id").ToStringNoPoint(),
				IsForce:  container.GetValFromMapMaybe(params, "is_force").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) InsertOrUpdateAppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		params, _ := c.Get("params")
		ss_log.Info("%v", params)

		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateAppVersionRequest{
				VId:         container.GetValFromMapMaybe(params, "v_id").ToStringNoPoint(),
				VsType:      container.GetValFromMapMaybe(params, "vs_type").ToStringNoPoint(),
				System:      container.GetValFromMapMaybe(params, "system").ToStringNoPoint(),
				UpType:      container.GetValFromMapMaybe(params, "up_type").ToStringNoPoint(),
				Description: container.GetValFromMapMaybe(params, "description").ToStringNoPoint(),

				Note: container.GetValFromMapMaybe(params, "note").ToStringNoPoint(),
				//VsCode:   container.GetValFromMapMaybe(params, "vs_code").ToStringNoPoint(),
				FileId: container.GetValFromMapMaybe(params, "file_id").ToStringNoPoint(),
				//IsForce:   container.GetValFromMapMaybe(params, "is_force").ToStringNoPoint(),
				Status:                  container.GetValFromMapMaybe(params, "status").ToStringNoPoint(),
				ConsecutiveLitersNumber: container.GetValFromMapMaybe(params, "consecutive_liters_number").ToStringNoPointReg(`^[\d]+$`),
				AccountNo:               inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.CheckInsertOrUpdateAppVersionVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateAppVersion(context.TODO(), req, global.RequestTimeoutOptions)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetNewVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &custProto.GetNewVersionRequest{
				System: container.GetValFromMapMaybe(params, "system").ToStringNoPoint(),
				VsType: container.GetValFromMapMaybe(params, "vs_type").ToStringNoPoint(),
			}
			if errStr := verify.CheckGetNewVersionVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := CustHandlerInst.Client.GetNewVersion(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"new_version": reply.NewVersion,
				"update_time": reply.UpdateTime,
			}, err
		}, "params")
	}
}

func (s *CustHandler) UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			accNo := inner_util.GetJwtDataString(c, "account_uid")
			accType := inner_util.GetJwtDataString(c, "account_type")
			//fileId := container.GetValFromMapMaybe(params, "file_id").ToStringNoPoint()
			baseName := container.GetValFromMapMaybe(params, "filename").ToStringNoPoint()
			filenameWithSuffix := strings.ToLower(container.GetValFromMapMaybe(params, "filename_with_suffix").ToStringNoPoint())
			uploadPath := container.GetValFromMapMaybe(params, "upload_path").ToStringNoPoint()

			ss_log.Error("获取到的数据 accNo=[%v],baseName=[%v],filenameWithSuffix=[%v]", accNo, baseName, filenameWithSuffix)
			if accNo == "" || baseName == "" || filenameWithSuffix == "" {
				ss_log.Error("获取到的数据为空accNo=[%v],baseName=[%v],filenameWithSuffix=[%v]", accNo, baseName, filenameWithSuffix)
				return ss_err.ERR_PARAM, nil, nil
			}

			//确认一遍是否文件存在
			boolean, _ := file.Exists(uploadPath)
			if !boolean {
				ss_log.Error("上传的文件不存在,[%v]", uploadPath)
				return ss_err.ERR_FILE_OP_FAILD, nil, nil
			}

			fileName := ""
			fileType := ""
			//转存到该存的地方
			switch filenameWithSuffix {
			case ".ipa":
				fileType = constants.UploadFileType_IPA
				fileName = aws_s3.App_Dir + "/" + baseName
			case ".apk":
				fileType = constants.UploadFileType_APK
				fileName = aws_s3.App_Dir + "/" + baseName
			default:
				ss_log.Error("文件后缀名[%v]错误", filenameWithSuffix)
				return ss_err.ERR_PARAM, nil, nil
			}

			var errAwsS3 error
			if filenameWithSuffix == ".apk" { // apk文件单独上传
				_, errAwsS3 = common.UploadS3.UploadAPKFile(uploadPath, fileName, true)
			} else {
				_, errAwsS3 = common.UploadS3.UploadFile(uploadPath, fileName, true)
			}

			if errAwsS3 != nil {
				ss_log.Error("上传到AwsS3失败，errAwsS3:[%v]", errAwsS3)
				return ss_err.ERR_UPLOAD, nil, nil
			}

			//添加上传app版本日志
			fileLogId, err := dao.UploadFileLogDaoInstance.AddUploadFileLog(accNo, accType, fileName, fileType)
			if err != ss_err.ERR_SUCCESS {
				ss_log.Error("添加上传文件日志失败")
				return ss_err.ERR_PARAM, nil, nil
			}

			return ss_err.ERR_SUCCESS, gin.H{
				"filename": fileLogId,
			}, nil
		}, "params")

	}
}

func (s *CustHandler) GetCashiers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetCashiersRequest{
				Page:       strext.ToInt32(params[0]),
				PageSize:   strext.ToInt32(params[1]),
				ServicerNo: strext.ToStringNoPoint(params[2]),
			}

			if errStr := verify.CheckGetCashiersVerify(req); errStr != "" {
				return errStr, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetCashiers(context.TODO(), req)
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "servicer_no")
	}
}

func (s *CustHandler) GetCashierDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &custProto.GetCashierDetailRequest{
				CashierNo: container.GetValFromMapMaybe(params, "cashier_no").ToStringNoPoint(),
			}
			if req.CashierNo == "" {
				ss_log.Error("CashierNo参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}
			reply, err := CustHandlerInst.Client.GetCashierDetail(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Datas,
			}, err
		}, "params")
	}
}

func (s *CustHandler) DeleteCashier() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteCashierRequest{
				CashierNo: container.GetValFromMapMaybe(params, "cashier_no").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.CheckDeleteCashierVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteCashier(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) ModifyCashier() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			ss_log.Error("此接口已废除，后台只允许帮服务商创建和删除店员，不允许更改店员的手机号")
			return ss_err.ERR_SYS_NO_API_AUTH, 0, nil
			//
			//req := &custProto.ModifyCashierRequest{
			//	CashierNo: container.GetValFromMapMaybe(params, "cashier_no").ToStringNoPoint(),
			//	Phone:     container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
			//}
			//if errStr := verify.CheckModifyCashierVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}
			//
			//reply, err := CustHandlerInst.Client.ModifyCashier(context.TODO(), req)
			//return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetCommonHelpCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCommonHelpCount(context.TODO(), &custProto.GetCommonHelpCountRequest{})
			return reply.ResultCode, reply.Datas, 0, err
		}, "")
	}
}

func (s *CustHandler) SwapConsultationIdx() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, err := CustHandlerInst.Client.SwapConsultationIdx(context.TODO(), &custProto.SwapConsultationIdxRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				Idx:      container.GetValFromMapMaybe(params, "idx").ToStringNoPoint(),
				Lang:     container.GetValFromMapMaybe(params, "lang").ToStringNoPoint(),
				SwapType: container.GetValFromMapMaybe(params, "swap_type").ToStringNoPoint(),
			})
			if err != nil {
				ss_log.Error("SwapConsultationIdx err=[%v]", err)
				return response.ResultCode, "", nil
			}
			return response.ResultCode, "", err
		})
	}
}

func (*CustHandler) GetPosChannels() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetPosChannels(context.TODO(), &custProto.GetPosChannelsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				ChannelName: strext.ToStringNoPoint(params[2]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "channel_name")
	}
}

func (*CustHandler) GetUseChannels() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetUseChannels(context.TODO(), &custProto.GetUseChannelsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				ChannelName:  strext.ToStringNoPoint(params[2]),
				CurrencyType: strext.ToStringNoPoint(params[3]),
				ChannelType:  strext.ToStringNoPoint(params[4]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "channel_name", "currency_type", "channel_type")
	}
}

func (s *CustHandler) InsertUseChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertUseChannelRequest{
				Id:               container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				ChannelNo:        container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				CurrencyType:     container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				SupportType:      container.GetValFromMapMaybe(params, "support_type").ToStringNoPoint(),
				SaveRate:         container.GetValFromMapMaybe(params, "save_rate").ToStringNoPoint(),
				SaveSingleMinFee: container.GetValFromMapMaybe(params, "save_single_min_fee").ToStringNoPoint(),

				SaveMaxAmount:        container.GetValFromMapMaybe(params, "save_max_amount").ToStringNoPoint(),
				SaveChargeType:       container.GetValFromMapMaybe(params, "save_charge_type").ToStringNoPoint(),
				WithdrawRate:         container.GetValFromMapMaybe(params, "withdraw_rate").ToStringNoPoint(),
				WithdrawSingleMinFee: container.GetValFromMapMaybe(params, "withdraw_single_min_fee").ToStringNoPoint(),

				WithdrawMaxAmount:  container.GetValFromMapMaybe(params, "withdraw_max_amount").ToStringNoPoint(),
				WithdrawChargeType: container.GetValFromMapMaybe(params, "withdraw_charge_type").ToStringNoPoint(),

				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}
			//if errStr := verify.InsertPosChannelVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}

			//开始添加或更新渠道
			response, _ := CustHandlerInst.Client.InsertUseChannel(context.TODO(), req)
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) DeleteUseChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.DeleteUseChannel(context.TODO(), &custProto.DeleteUseChannelRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) ModifyUseChannelStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.ModifyUseChannelStatus(context.TODO(), &custProto.ModifyUseChannelStatusRequest{
				Id:        container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) GetUseChannelDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &custProto.GetUseChannelDetailRequest{
				Id: container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
			}
			if req.Id == "" {
				ss_log.Error("Id参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}
			reply, err := CustHandlerInst.Client.GetUseChannelDetail(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "params")
	}
}

func (*CustHandler) GetToCustList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetToCustList(context.TODO(), &custProto.GetToCustListRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				LogNo:       strext.ToStringNoPoint(params[4]),
				OrderStatus: strext.ToStringNoPoint(params[5]),
				Account:     strext.ToStringNoPoint(params[6]),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "log_no", "order_status", "account")
	}
}

func (*CustHandler) GetCustToHeadquartersList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCustToHeadquartersList(context.TODO(), &custProto.GetCustToHeadquartersListRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				LogNo:       strext.ToStringNoPoint(params[4]),
				OrderStatus: strext.ToStringNoPoint(params[5]),
				Account:     strext.ToStringNoPoint(params[6]),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "log_no", "order_status", "account")
	}
}

func (s *CustHandler) GetImgUrl() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			imageId := container.GetValFromMapMaybe(params, "image_id").ToStringNoPoint()
			imgType := container.GetValFromMapMaybe(params, "img_type").ToStringNoPoint()
			if imageId == "" {
				ss_log.Error("参数ImageId为空")
				return ss_err.ERR_PARAM, nil, nil
			}
			switch imgType {
			case "1": // 需要授权

				req := &custProto.AuthDownloadImageRequest{ImageId: imageId}
				reply, err := CustHandlerInst.Client.AuthDownloadImage(context.TODO(), req, global.RequestTimeoutOptions)
				ss_log.Info("reply=[%v],err=[%v]", err)
				return reply.ResultCode, gin.H{
					"image_str": reply.ImageStr,
				}, nil

			case "2": // 不需要授权
				req := &custProto.UnAuthDownloadImageRequest{
					ImageId: imageId,
				}
				reply, err := CustHandlerInst.Client.UnAuthDownloadImage(context.TODO(), req)
				ss_log.Info("reply=[%v],err=[%v]", err)
				return reply.ResultCode, gin.H{
					"image_url": reply.ImageUrl,
				}, nil
			default:
				ss_log.Error("imgType参数错误")
				return ss_err.ERR_PARAM, nil, nil
			}

		}, "params")
	}
}

func (s *CustHandler) GetClientInfos() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetClientInfos(context.TODO(), &custProto.GetClientInfosRequest{
				Page:       strext.ToInt32(params[0]),
				PageSize:   strext.ToInt32(params[1]),
				ClientType: strext.ToStringNoPoint(params[2]), //cust还是servicer的
				Id:         strext.ToStringNoPoint(params[3]), //
				Platform:   strext.ToStringNoPoint(params[4]), //ios|android
				Account:    strext.ToStringNoPoint(params[5]),
				Uuid:       strext.ToStringNoPoint(params[6]),
			})
			ss_log.Info("err=[%v]", err)
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "client_type", "id", "platform", "account", "uuid")
	}
}

func (s *CustHandler) GetLogAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetLogAccounts(context.TODO(), &custProto.GetLogAccountsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]), //
				EndTime:   strext.ToStringNoPoint(params[3]), //
				LogType:   strext.ToStringNoPoint(params[4]), //web、cli_app、cli_pos
				LogNo:     strext.ToStringNoPoint(params[5]), //
				Type:      strext.ToStringNoPoint(params[6]), //参考：constants.LogAccountWebType_Account
				Account:   strext.ToStringNoPoint(params[7]), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "log_type", "log_no", "type", "account")
	}
}

func (s *CustHandler) GetPushConfs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetPushConfs(context.TODO(), &custProto.GetPushConfsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]), //
				EndTime:   strext.ToStringNoPoint(params[3]), //
				UseStatus: strext.ToStringNoPoint(params[4]), //
				Pusher:    strext.ToStringNoPoint(params[5]), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "use_status", "search")
	}
}

func (s *CustHandler) GetPushConf() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetPushConf(context.TODO(), &custProto.GetPushConfRequest{
				PusherNo: strext.ToStringNoPoint(params[0]), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, err
		}, "pusher_no")
	}
}

func (s *CustHandler) GetPushTemp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			tempNo := strext.ToStringNoPoint(params[0])
			if tempNo == "" {
				ss_log.Error("GetPushTemp 接口,tempNO 为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}
			reply, err := CustHandlerInst.Client.GetPushTemp(context.TODO(), &custProto.GetPushTempRequest{
				TempNo: tempNo, //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, err
		}, "temp_no")
	}
}

func (s *CustHandler) GetPushRecords() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetPushRecords(context.TODO(), &custProto.GetPushRecordsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]), //
				EndTime:   strext.ToStringNoPoint(params[3]), //
				Status:    strext.ToStringNoPoint(params[4]), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "status")
	}
}

func (s *CustHandler) GetPushTemps() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetPushTemps(context.TODO(), &custProto.GetPushTempsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) InsertOrUpdatePushConf() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdatePushConfsRequest{
				PusherNo:       container.GetValFromMapMaybe(params, "pusher_no").ToStringNoPoint(),
				Pusher:         container.GetValFromMapMaybe(params, "pusher").ToStringNoPoint(),
				Config:         container.GetValFromMapMaybe(params, "config").ToStringNoPoint(),
				UseStatus:      container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				ConditionValue: container.GetValFromMapMaybe(params, "condition_value").ToStringNoPoint(),
				ConditionType:  container.GetValFromMapMaybe(params, "condition_type").ToStringNoPoint(),
				LoginUid:       inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.InsertOrUpdatePushConfsVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdatePushConfs(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}
func (s *CustHandler) InsertOrUpdatePushTemp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdatePushTempRequest{
				TempNo:     container.GetValFromMapMaybe(params, "temp_no").ToStringNoPoint(),
				PushNos:    container.GetValFromMapMaybe(params, "push_nos").ToStringNoPoint(),
				TitleKey:   container.GetValFromMapMaybe(params, "title_key").ToStringNoPoint(),
				ContentKey: container.GetValFromMapMaybe(params, "content_key").ToStringNoPoint(),
				LenArgs:    container.GetValFromMapMaybe(params, "len_args").ToStringNoPoint(),
				OpType:     container.GetValFromMapMaybe(params, "op_type").ToInt32(),
				LoginUid:   inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.InsertOrUpdatePushTempVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdatePushTemp(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeletePushConf() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeletePushConfRequest{
				PusherNo: container.GetValFromMapMaybe(params, "pusher_no").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.DeletePushConfVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			var reply, err = CustHandlerInst.Client.DeletePushConf(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) ModifyPushConfStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyPushConfStatusRequest{
				PusherNo:  container.GetValFromMapMaybe(params, "pusher_no").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}

			//if errStr := verify.DeletePushConfVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}

			var reply, err = CustHandlerInst.Client.ModifyPushConfStatus(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

//Event
func (s *CustHandler) GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetEvents(context.TODO(), &custProto.GetEventsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]), //
				EndTime:   strext.ToStringNoPoint(params[3]), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time")
	}
}

func (s *CustHandler) GetEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetEvent(context.TODO(), &custProto.GetEventRequest{
				EventNo: strext.ToStringNoPoint(params[0]), //
			})

			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, err
		}, "event_no")
	}
}

func (s *CustHandler) InsertOrUpdateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateEventRequest{
				EventNo:   container.GetValFromMapMaybe(params, "event_no").ToStringNoPoint(),
				EventName: container.GetValFromMapMaybe(params, "event_name").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.InsertOrUpdateEventVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateEvent(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteEventRequest{
				EventNo: container.GetValFromMapMaybe(params, "event_no").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.EventNo == "" {
				ss_log.Error("EventNo参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteEvent(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

//EvaParam
func (s *CustHandler) GetEvaParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetEvaParams(context.TODO(), &custProto.GetEvaParamsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) InsertOrUpdateEvaParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateEvaParamRequest{
				Key: container.GetValFromMapMaybe(params, "key").ToStringNoPoint(),
				Val: container.GetValFromMapMaybe(params, "val").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.InsertOrUpdateEvaParamVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateEvaParam(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteEvaParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteEvaParamRequest{
				Key: container.GetValFromMapMaybe(params, "key").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.Key == "" {
				ss_log.Error("Key参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteEvaParam(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetGlobalParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetGlobalParam(context.TODO(), &custProto.GetGlobalParamRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) InsertOrUpdateGlobalParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateGlobalParamRequest{
				ParamKey:   container.GetValFromMapMaybe(params, "param_key").ToStringNoPoint(),
				ParamValue: container.GetValFromMapMaybe(params, "param_value").ToStringNoPoint(),
				Remark:     container.GetValFromMapMaybe(params, "remark").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if errStr := verify.InsertOrUpdateGlobalParamVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateGlobalParam(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteGlobalParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteGlobalParamRequest{
				ParamKey: container.GetValFromMapMaybe(params, "param_key").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.ParamKey == "" {
				ss_log.Error("Key参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteGlobalParam(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetLogResults() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetLogResults(context.TODO(), &custProto.GetLogResultsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) GetOps() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetOps(context.TODO(), &custProto.GetOpsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) GetOp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetOp(context.TODO(), &custProto.GetOpRequest{
				OpNo: strext.ToStringNoPoint(params[0]), //
			})

			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, err
		}, "op_no")
	}
}

func (s *CustHandler) InsertOrUpdateOp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateOpRequest{
				OpNo:       container.GetValFromMapMaybe(params, "op_no").ToStringNoPoint(),
				OpName:     container.GetValFromMapMaybe(params, "op_name").ToStringNoPoint(),
				ScriptName: container.GetValFromMapMaybe(params, "script_name").ToStringNoPoint(),
				Param:      container.GetValFromMapMaybe(params, "param").ToStringNoPoint(),
				Score:      container.GetValFromMapMaybe(params, "score").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			//if errStr := verify.InsertOrUpdateEventVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}

			reply, err := CustHandlerInst.Client.InsertOrUpdateOp(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteOp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteOpRequest{
				OpNo: container.GetValFromMapMaybe(params, "op_no").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.OpNo == "" {
				ss_log.Error("EventNo参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteOp(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetRelaApiEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetRelaApiEvents(context.TODO(), &custProto.GetRelaApiEventsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) InsertOrUpdateRelaApiEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateRelaApiEventRequest{
				ApiType: container.GetValFromMapMaybe(params, "api_type").ToStringNoPoint(),
				EventNo: container.GetValFromMapMaybe(params, "event_no").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			//if errStr := verify.InsertOrUpdateEventVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}

			reply, err := CustHandlerInst.Client.InsertOrUpdateRelaApiEvent(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteRelaApiEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteRelaApiEventRequest{
				ApiType: container.GetValFromMapMaybe(params, "api_type").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.ApiType == "" {
				ss_log.Error("ApiType参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteRelaApiEvent(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetRelaEventRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetRelaEventRules(context.TODO(), &custProto.GetRelaEventRulesRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) InsertOrUpdateRelaEventRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateRelaEventRuleRequest{
				RelaNo:  container.GetValFromMapMaybe(params, "rela_no").ToStringNoPoint(),
				EventNo: container.GetValFromMapMaybe(params, "event_no").ToStringNoPoint(),
				RuleNo:  container.GetValFromMapMaybe(params, "rule_no").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			//if errStr := verify.InsertOrUpdateEventVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}

			reply, err := CustHandlerInst.Client.InsertOrUpdateRelaEventRule(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteRelaEventRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteRelaEventRuleRequest{
				RelaNo: container.GetValFromMapMaybe(params, "rela_no").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.RelaNo == "" {
				ss_log.Error("RelaNo参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteRelaEventRule(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetRiskThresholds() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetRiskThresholds(context.TODO(), &custProto.GetRiskThresholdsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) InsertOrUpdateRiskThreshold() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateRiskThresholdRequest{
				RuleNo:        container.GetValFromMapMaybe(params, "rule_no").ToStringNoPoint(),
				RiskThreshold: container.GetValFromMapMaybe(params, "risk_threshold").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			//if errStr := verify.InsertOrUpdateEventVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}

			reply, err := CustHandlerInst.Client.InsertOrUpdateRiskThreshold(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteRiskThreshold() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteRiskThresholdRequest{
				RuleNo: container.GetValFromMapMaybe(params, "rule_no").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.RuleNo == "" {
				ss_log.Error("RuleNo参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteRiskThreshold(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetRules(context.TODO(), &custProto.GetRulesRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *CustHandler) GetRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetRule(context.TODO(), &custProto.GetRuleRequest{
				RuleNo: strext.ToStringNoPoint(params[0]), //
			})

			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, err
		}, "rule_no")
	}
}

func (s *CustHandler) InsertOrUpdateRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateRuleRequest{
				RuleNo:   container.GetValFromMapMaybe(params, "rule_no").ToStringNoPoint(),
				RuleName: container.GetValFromMapMaybe(params, "rule_name").ToStringNoPoint(),
				Rule:     container.GetValFromMapMaybe(params, "rule").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateRule(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeleteRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		//params, _ := c.Get("params")
		//ss_log.Info("%v", params)
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeleteRuleRequest{
				RuleNo: container.GetValFromMapMaybe(params, "rule_no").ToStringNoPoint(),
				//LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.RuleNo == "" {
				ss_log.Error("EventNo参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DeleteRule(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

//用户身份认证
func (*CustHandler) GetAuthMaterials() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAuthMaterials(context.TODO(), &custProto.GetAuthMaterialsRequest{
				Page:           strext.ToInt32(params[0]),
				PageSize:       strext.ToInt32(params[1]),
				StartTime:      strext.ToString(params[2]),
				EndTime:        strext.ToString(params[3]),
				Account:        strext.ToString(params[4]),
				AuthMaterialNo: strext.ToString(params[5]),
				Status:         strext.ToString(params[6]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "account", "auth_material_no", "status")
	}
}

func (*CustHandler) ModifyAuthMaterialStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyAuthMaterialStatusRequest{
				AuthMaterialNo: container.GetValFromMapMaybe(params, "auth_material_no").ToString(),
				Status:         container.GetValFromMapMaybe(params, "status").ToString(),
				AccountUid:     container.GetValFromMapMaybe(params, "account_uid").ToString(), //该条认证信息的账号uid
				LoginUid:       inner_util.GetJwtDataString(c, "account_uid"),                  //后台登陆账号的uid
			}

			if errStr := verify.ModifyAuthMaterialStatusVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.ModifyAuthMaterialStatus(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

//个人商家认证材料列表
func (*CustHandler) GetAuthMaterialBusinessList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAuthMaterialBusinessList(context.TODO(), &custProto.GetAuthMaterialBusinessListRequest{
				Page:           strext.ToInt32(params[0]),
				PageSize:       strext.ToInt32(params[1]),
				StartTime:      strext.ToString(params[2]),
				EndTime:        strext.ToString(params[3]),
				Account:        strext.ToString(params[4]),
				AuthMaterialNo: strext.ToString(params[5]),
				Status:         strext.ToString(params[6]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "account", "auth_material_no", "status")
	}
}

func (*CustHandler) ModifyAuthMaterialBusinessStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyAuthMaterialBusinessStatusRequest{
				AuthMaterialNo: container.GetValFromMapMaybe(params, "auth_material_no").ToString(),
				Status:         container.GetValFromMapMaybe(params, "status").ToString(),
				LoginUid:       inner_util.GetJwtDataString(c, "account_uid"), //后台登陆账号的uid
			}

			if errStr := verify.ModifyAuthMaterialBusinessStatusVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.ModifyAuthMaterialBusinessStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

//企业商家认证材料列表
func (*CustHandler) GetAuthMaterialEnterpriseList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAuthMaterialEnterpriseList(context.TODO(), &custProto.GetAuthMaterialEnterpriseListRequest{
				Page:           strext.ToInt32(params[0]),
				PageSize:       strext.ToInt32(params[1]),
				StartTime:      strext.ToString(params[2]),
				EndTime:        strext.ToString(params[3]),
				Account:        strext.ToString(params[4]),
				AuthMaterialNo: strext.ToString(params[5]),
				Status:         strext.ToString(params[6]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "account", "auth_material_no", "status")
	}
}

//审核企业商家认证材料
func (*CustHandler) ModifyAuthMaterialEnterpriseStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyAuthMaterialEnterpriseStatusRequest{
				AuthMaterialNo: container.GetValFromMapMaybe(params, "auth_material_no").ToString(),
				Status:         container.GetValFromMapMaybe(params, "status").ToString(),
				AccountUid:     container.GetValFromMapMaybe(params, "account_uid").ToString(), //该条认证信息的账号uid
				LoginUid:       inner_util.GetJwtDataString(c, "account_uid"),                  //后台登陆账号的uid
			}

			if errStr := verify.ModifyAuthMaterialEnterpriseStatusVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.ModifyAuthMaterialEnterpriseStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetStatisticUserWithdraws() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserWithdraw(context.TODO(), &custProto.GetStatisticUserWithdrawRequest{
				StartDate: strext.ToStringNoPoint(params[0]),
				EndDate:   strext.ToStringNoPoint(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"datas":           reply.DateList,
				"usd_amount_list": reply.UsdAmountList,
				"usd_fee_list":    reply.UsdFeeList,
				"usd_num_list":    reply.UsdNumList,

				"khr_amount_list": reply.KhrAmountList,
				"khr_fee_list":    reply.KhrFeeList,
				"khr_num_list":    reply.KhrNumList,
			}, 0, err
		}, "start_date", "end_date")
	}
}

func (s *CustHandler) GetStatisticUserWithdrawList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserWithdrawList(context.TODO(), &custProto.GetStatisticUserWithdrawListRequest{
				StartDate:    strext.ToStringNoPoint(params[0]),
				EndDate:      strext.ToStringNoPoint(params[1]),
				Page:         strext.ToInt32(params[2]),
				PageSize:     strext.ToInt32(params[3]),
				WithdrawType: strext.ToStringNoPoint(params[4]),
				CurrencyType: strext.ToStringNoPoint(params[5]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "start_date", "end_date", "page", "page_size", "withdraw_type", "currency_type")
	}
}

func (s *CustHandler) GetStatisticUserRecharges() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserRecharge(context.TODO(), &custProto.GetStatisticUserRechargeRequest{
				StartDate: strext.ToStringNoPoint(params[0]),
				EndDate:   strext.ToStringNoPoint(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"datas":           reply.DateList,
				"usd_amount_list": reply.UsdAmountList,
				"usd_fee_list":    reply.UsdFeeList,
				"usd_num_list":    reply.UsdNumList,

				"khr_amount_list": reply.KhrAmountList,
				"khr_fee_list":    reply.KhrFeeList,
				"khr_num_list":    reply.KhrNumList,
			}, 0, err
		}, "start_date", "end_date")
	}
}

func (s *CustHandler) GetStatisticUserRechargeList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserRechargeList(context.TODO(), &custProto.GetStatisticUserRechargeListRequest{
				StartDate:    strext.ToStringNoPoint(params[0]),
				EndDate:      strext.ToStringNoPoint(params[1]),
				Page:         strext.ToInt32(params[2]),
				PageSize:     strext.ToInt32(params[3]),
				RechargeType: strext.ToStringNoPoint(params[4]),
				CurrencyType: strext.ToStringNoPoint(params[5]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "start_date", "end_date", "page", "page_size", "recharge_type", "currency_type")
	}
}

func (s *CustHandler) GetStatisticUserExchanges() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserExchange(context.TODO(), &custProto.GetStatisticUserExchangeRequest{
				StartDate: strext.ToStringNoPoint(params[0]),
				EndDate:   strext.ToStringNoPoint(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"datas":               reply.DateList,
				"usd2khr_amount_list": reply.Usd2KhrAmountList,
				"usd2khr_fee_list":    reply.Usd2KhrFeeList,
				"usd2khr_num_list":    reply.Usd2KhrNumList,

				"khr2usd_amount_list": reply.Khr2UsdAmountList,
				"khr2usd_fee_list":    reply.Khr2UsdFeeList,
				"khr2usd_num_list":    reply.Khr2UsdNumList,
			}, 0, err
		}, "start_date", "end_date")
	}
}

func (s *CustHandler) GetStatisticUserExchangeList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserExchangeList(context.TODO(), &custProto.GetStatisticUserExchangeListRequest{
				StartDate: strext.ToStringNoPoint(params[0]),
				EndDate:   strext.ToStringNoPoint(params[1]),
				Page:      strext.ToInt32(params[2]),
				PageSize:  strext.ToInt32(params[3]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "start_date", "end_date", "page", "page_size")
	}
}

func (s *CustHandler) GetStatisticUserTransfers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserTransfer(context.TODO(), &custProto.GetStatisticUserTransferRequest{
				StartDate: strext.ToStringNoPoint(params[0]),
				EndDate:   strext.ToStringNoPoint(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"datas":           reply.DateList,
				"usd_amount_list": reply.UsdAmountList,
				"usd_fee_list":    reply.UsdFeeList,
				"usd_num_list":    reply.UsdNumList,

				"khr_amount_list": reply.KhrAmountList,
				"khr_fee_list":    reply.KhrFeeList,
				"khr_num_list":    reply.KhrNumList,
			}, 0, err
		}, "start_date", "end_date")
	}
}

func (s *CustHandler) GetStatisticUserTransferList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserTransferList(context.TODO(), &custProto.GetStatisticUserTransferListRequest{
				StartDate:    strext.ToStringNoPoint(params[0]),
				EndDate:      strext.ToStringNoPoint(params[1]),
				Page:         strext.ToInt32(params[2]),
				PageSize:     strext.ToInt32(params[3]),
				CurrencyType: strext.ToStringNoPoint(params[4]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "start_date", "end_date", "page", "page_size", "currency_type")
	}
}

func (s *CustHandler) GetStatisticDates() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticDate(context.TODO(), &custProto.GetStatisticDateRequest{
				StartDate: strext.ToStringNoPoint(params[0]),
				EndDate:   strext.ToStringNoPoint(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"datas":                 reply.DateList,
				"reg_user_num_list":     reply.RegUserNumList,
				"reg_servicer_num_list": reply.RegServicerNumList,
			}, 0, err
		}, "start_date", "end_date")
	}
}

func (s *CustHandler) GetStatisticDateList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticDateList(context.TODO(), &custProto.GetStatisticDateListRequest{
				StartDate: strext.ToStringNoPoint(params[0]),
				EndDate:   strext.ToStringNoPoint(params[1]),
				Page:      strext.ToInt32(params[2]),
				PageSize:  strext.ToInt32(params[3]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "start_date", "end_date", "page", "page_size")
	}
}

func (s *CustHandler) ReStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ReStatisticRequest{
				Type:      container.GetValFromMapMaybe(params, "type").ToStringNoPoint(),
				StartDate: container.GetValFromMapMaybe(params, "start_date").ToStringNoPoint(),
				EndDate:   container.GetValFromMapMaybe(params, "end_date").ToStringNoPoint(),
			}
			if errStr := verify.ReStatisticVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.ReStatistic(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DeletePushTemp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DeletePushTempRequest{
				TempNo:   container.GetValFromMapMaybe(params, "temp_no").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			}

			if errStr := verify.DeletePushTempVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			var reply, err = CustHandlerInst.Client.DeletePushTemp(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetStatisticUserMoneyDates() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserMoney(context.TODO(), &custProto.GetStatisticUserMoneyRequest{
				StartTime: strext.ToStringNoPoint(params[0]),
				EndTime:   strext.ToStringNoPoint(params[1]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, gin.H{
				"datas":                   reply.TimeList,
				"usd_balance_list":        reply.UsdBalanceList,
				"khr_balance_list":        reply.KhrBalanceList,
				"usd_frozen_balance_list": reply.UsdFrozenBalanceList,
				"khr_frozen_balance_list": reply.KhrFrozenBalanceList,
			}, 0, err
		}, "start_time", "end_time")
	}
}

func (s *CustHandler) GetStatisticUserMoneyList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetStatisticUserMoneyList(context.TODO(), &custProto.GetStatisticUserMoneyListRequest{
				StartTime: strext.ToStringNoPoint(params[0]),
				EndTime:   strext.ToStringNoPoint(params[1]),
				Page:      strext.ToInt32(params[2]),
				PageSize:  strext.ToInt32(params[3]),
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "start_time", "end_time", "page", "page_size")
	}
}

func (*CustHandler) GetBusinessIndustryList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessIndustryList(context.TODO(), &custProto.GetBusinessIndustryListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				Code:      strext.ToStringNoPoint(params[4]),

				NameCh: strext.ToStringNoPoint(params[5]),
				NameEn: strext.ToStringNoPoint(params[6]),
				Level:  strext.ToStringNoPoint(params[7]),
				UpCode: strext.ToStringNoPoint(params[8]),
				Search: strext.ToStringNoPoint(params[9]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "code", "name_ch", "name_en", "level", "up_code", "search")
	}
}

func (*CustHandler) GetBusinessIndustryDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessIndustryDetail(context.TODO(), &custProto.GetBusinessIndustryDetailRequest{
				Code: strext.ToStringNoPoint(params[0]),
			})

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, err
		}, "code")
	}
}

func (*CustHandler) InsertOrUpdateBusinessIndustry() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateBusinessIndustryRequest{
				Code:   container.GetValFromMapMaybe(params, "code").ToString(),    //
				NameCh: container.GetValFromMapMaybe(params, "name_ch").ToString(), //
				NameEn: container.GetValFromMapMaybe(params, "name_en").ToString(), //
				NameKm: container.GetValFromMapMaybe(params, "name_km").ToString(), //
				Level:  container.GetValFromMapMaybe(params, "level").ToString(),   //

				UpCode: container.GetValFromMapMaybe(params, "up_code").ToString(), //
			}

			if errStr := verify.InsertOrUpdateBusinessIndustryVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateBusinessIndustry(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) DelBusinessIndustry() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DelBusinessIndustryRequest{
				Code: container.GetValFromMapMaybe(params, "code").ToString(), //
			}

			reply, err := CustHandlerInst.Client.DelBusinessIndustry(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetCashRechargeOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCashRechargeOrderList(context.TODO(), &custProto.GetCashRechargeOrderListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				LogNo:        strext.ToString(params[2]),
				StartTime:    strext.ToString(params[3]),
				EndTime:      strext.ToString(params[4]),
				OrderStatus:  strext.ToString(params[5]),
				CurrencyType: strext.ToString(params[6]),
				AccAccount:   strext.ToString(params[7]),
				OpAccAccount: strext.ToString(params[8]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "log_no", "start_time", "end_time", "order_status", "currency_type", "acc_account", "op_acc_account")
	}
}

func (*CustHandler) GetBusinessSceneList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessSceneList(context.TODO(), &custProto.GetBusinessSceneListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				SceneName: strext.ToStringNoPoint(params[4]),
				Lang:      constants.LangZhCN,
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "scene_name")
	}
}
func (*CustHandler) GetBusinessSceneDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessSceneDetail(context.TODO(), &custProto.GetBusinessSceneDetailRequest{
				SceneNo: strext.ToStringNoPoint(params[0]),
				//Lang:      constants.LangZhCN,//注意，这里的详情主要用于修改，所以返回key是正确的
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "scene_no")
	}
}

func (*CustHandler) InsertOrUpdateBusinessScene() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			loginUid := inner_util.GetJwtDataString(c, "account_uid")

			if loginUid == "" {
				ss_log.Error("loginUid参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			req := &custProto.InsertOrUpdateBusinessSceneRequest{
				LoginUid:       loginUid,
				ImageStr:       container.GetValFromMapMaybe(params, "img_base64").ToString(),
				ExampleImgs:    container.GetValFromMapMaybe(params, "example_imgs").ToString(),
				ExampleNames:   container.GetValFromMapMaybe(params, "example_names").ToString(),
				SceneNo:        container.GetValFromMapMaybe(params, "scene_no").ToStringNoPoint(),
				SceneName:      container.GetValFromMapMaybe(params, "scene_name").ToStringNoPoint(),
				Notes:          container.GetValFromMapMaybe(params, "notes").ToStringNoPoint(),
				TradeType:      container.GetValFromMapMaybe(params, "trade_type").ToStringNoPoint(),
				PaymentChannel: container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				FloatRate:      container.GetValFromMapMaybe(params, "float_rate").ToStringNoPointReg(`^-?[\d]+$`),
				ApplyType:      container.GetValFromMapMaybe(params, "apply_type").ToStringNoPoint(),
				IsManualSigned: container.GetValFromMapMaybe(params, "is_manual_signed").ToInt32(),
			}

			if req.SceneName == "" {
				ss_log.Error("SceneName参数为空")
				return ss_err.ERR_PARAM, 0, nil
			}

			if req.FloatRate == "" {
				ss_log.Error("FloatRate参数不合法")
				return ss_err.ERR_PARAM, 0, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateBusinessScene(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

//禁用产品
func (*CustHandler) IsEnabled() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_ADMIN && accountType != constants.AccountType_OPERATOR {
				ss_log.Error("账号权限不足，accountType=%v", accountType)
				return ss_err.ERR_ACCOUNT_NO_PERMISSION, nil, nil
			}
			req := &custProto.IsEnabledSceneRequest{
				SceneNo:   container.GetValFromMapMaybe(params, "scene_no").ToString(), //
				IsEnabled: container.GetValFromMapMaybe(params, "is_enabled").ToString(),
			}

			reply, err := CustHandlerInst.Client.IsEnabledScene(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetBusinessSignedList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessSignedList(context.TODO(), &custProto.GetBusinessSignedListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				Status:    strext.ToStringNoPoint(params[4]),
				AppId:     strext.ToStringNoPoint(params[5]),
				Lang:      constants.LangZhCN,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "status", "app_id")
	}
}

func (*CustHandler) GetBusinessAppList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessAppList(context.TODO(), &custProto.GetBusinessAppListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				Account:   strext.ToStringNoPoint(params[4]),
				AppId:     strext.ToStringNoPoint(params[5]),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "account", "app_id")
	}
}

func (*CustHandler) UpdateBusinessAppStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessAppStatusRequest{
				AppId:  container.GetValFromMapMaybe(params, "app_id").ToString(), //
				Status: container.GetValFromMapMaybe(params, "status").ToString(), //
				Notes:  container.GetValFromMapMaybe(params, "notes").ToString(),  //
			}

			reply, err := CustHandlerInst.Client.UpdateBusinessAppStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetBusinessChannels() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessChannels(context.TODO(), &custProto.GetBusinessChannelsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				ChannelName: strext.ToStringNoPoint(params[2]),
				ChannelType: strext.ToStringNoPoint(params[3]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "channel_name", "channel_type")
	}
}

func (s *CustHandler) GetBusinessChannelDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &custProto.GetBusinessChannelDetailRequest{
				Id: container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
			}
			if req.Id == "" {
				ss_log.Error("Id参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}
			reply, err := CustHandlerInst.Client.GetBusinessChannelDetail(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "params")
	}
}

func (s *CustHandler) InsertBusinessChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertBusinessChannelRequest{
				Id:               container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				ChannelNo:        container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				CurrencyType:     container.GetValFromMapMaybe(params, "currency_type").ToStringNoPoint(),
				SupportType:      container.GetValFromMapMaybe(params, "support_type").ToStringNoPoint(),
				SaveRate:         container.GetValFromMapMaybe(params, "save_rate").ToStringNoPoint(),
				SaveSingleMinFee: container.GetValFromMapMaybe(params, "save_single_min_fee").ToStringNoPoint(),

				SaveMaxAmount:        container.GetValFromMapMaybe(params, "save_max_amount").ToStringNoPoint(),
				SaveChargeType:       container.GetValFromMapMaybe(params, "save_charge_type").ToStringNoPoint(),
				WithdrawRate:         container.GetValFromMapMaybe(params, "withdraw_rate").ToStringNoPoint(),
				WithdrawSingleMinFee: container.GetValFromMapMaybe(params, "withdraw_single_min_fee").ToStringNoPoint(),

				WithdrawMaxAmount:  container.GetValFromMapMaybe(params, "withdraw_max_amount").ToStringNoPoint(),
				WithdrawChargeType: container.GetValFromMapMaybe(params, "withdraw_charge_type").ToStringNoPoint(),
				LoginUid:           inner_util.GetJwtDataString(c, "account_uid"),
			}
			//if errStr := verify.InsertPosChannelVerify(req); errStr != "" {
			//	return errStr, nil, nil
			//}

			//开始添加或更新渠道
			response, _ := CustHandlerInst.Client.InsertBusinessChannel(context.TODO(), req)
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) DeleteBusinessChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.DeleteBusinessChannel(context.TODO(), &custProto.DeleteBusinessChannelRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (s *CustHandler) ModifyBusinessChannelStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			response, _ := CustHandlerInst.Client.ModifyBusinessChannelStatus(context.TODO(), &custProto.ModifyBusinessChannelStatusRequest{
				Id:        container.GetValFromMapMaybe(params, "id").ToStringNoPoint(),
				UseStatus: container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"),
			})
			return response.ResultCode, "", nil
		})
	}
}

func (*CustHandler) GetBusinessToHeadList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessToHeadList(context.TODO(), &custProto.GetBusinessToHeadListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				MoneyType: strext.ToStringNoPoint(params[4]),
				LogNo:     strext.ToStringNoPoint(params[5]),
				Account:   strext.ToStringNoPoint(params[6]),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "money_type", "log_no", "account")
	}
}

func (*CustHandler) GetToBusinessList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetToBusinessListRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				MoneyType:   strext.ToStringNoPoint(params[4]),
				LogNo:       strext.ToStringNoPoint(params[5]),
				OrderStatus: strext.ToStringNoPoint(params[6]),
				Account:     strext.ToStringNoPoint(params[7]),
			}

			reply, err := CustHandlerInst.Client.GetToBusinessList(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "money_type", "log_no", "order_status", "account")
	}
}

func (*CustHandler) GetBusinessList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessList(context.TODO(), &custProto.GetBusinessListRequest{
				Page:       strext.ToStringNoPoint(params[0]),
				PageSize:   strext.ToStringNoPoint(params[1]),
				StartTime:  strext.ToStringNoPoint(params[2]),
				EndTime:    strext.ToStringNoPoint(params[3]),
				AuthStatus: strext.ToStringNoPoint(params[4]),
				UseStatus:  strext.ToStringNoPoint(params[5]),
				Account:    strext.ToStringNoPoint(params[6]),
				BusinessId: strext.ToStringNoPoint(params[7]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "auth_status", "use_status", "account", "business_id")
	}
}

func (*CustHandler) ModifyBusinessStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyBusinessStatusRequest{
				BusinessNo: container.GetValFromMapMaybe(params, "business_no").ToString(), //
				Status:     container.GetValFromMapMaybe(params, "status").ToString(),      //

				//状态类型（use_status、income_authorization、outgo_authorization）
				StatusType: container.GetValFromMapMaybe(params, "status_type").ToString(),
			}

			if errStr := verify.ModifyBusinessStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.ModifyBusinessStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

//用户账单
func (*CustHandler) GetChangeBalanceOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetChangeBalanceOrders(context.TODO(), &custProto.GetChangeBalanceOrdersRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				LogNo:       strext.ToStringNoPoint(params[4]),
				Account:     strext.ToStringNoPoint(params[5]),
				AccountNo:   strext.ToStringNoPoint(params[6]),
				AccountType: strext.ToStringNoPoint(params[7]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "log_no", "account", "account_no", "account_type")
	}

}

func (s *CustHandler) GetBulletins() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBulletins(context.TODO(), &custProto.GetBulletinsRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),              //
				EndTime:     strext.ToStringNoPoint(params[3]),              //
				UseStatus:   strext.ToStringNoPoint(params[4]),              //
				Title:       strext.ToStringNoPoint(params[5]),              //
				AccountType: inner_util.GetJwtDataString(c, "account_type"), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "use_status", "title")
	}
}

func (s *CustHandler) InsertOrUpdateBulletin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateBulletinRequest{
				BulletinId: container.GetValFromMapMaybe(params, "bulletin_id").ToStringNoPoint(),
				Title:      container.GetValFromMapMaybe(params, "title").ToStringNoPoint(),
				Content:    container.GetValFromMapMaybe(params, "content").ToStringNoPoint(),
			}
			if errStr := verify.InsertOrUpdateBulletinVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateBulletin(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) DelBulletin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DelBulletinRequest{
				BulletinId: container.GetValFromMapMaybe(params, "bulletin_id").ToStringNoPoint(),
			}

			if req.BulletinId == "" {
				ss_log.Error("BulletinId参数为空")
				return ss_err.ERR_PARAM, 0, nil
			}

			reply, err := CustHandlerInst.Client.DelBulletin(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) UpdateBulletinStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBulletinStatusRequest{
				BulletinId: container.GetValFromMapMaybe(params, "bulletin_id").ToStringNoPoint(),
				Status:     container.GetValFromMapMaybe(params, "status").ToStringNoPoint(),
				StatusType: container.GetValFromMapMaybe(params, "status_type").ToStringNoPoint(), //top_status、use_status
			}

			if errStr := verify.UpdateBulletinStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.UpdateBulletinStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (s *CustHandler) GetBusinessMessagesList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessMessagesList(context.TODO(), &custProto.GetBusinessMessagesListRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]), //
				EndTime:     strext.ToStringNoPoint(params[3]), //
				Account:     strext.ToStringNoPoint(params[4]), //
				AccountType: strext.ToStringNoPoint(params[5]), //
			})
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "account", "account_type")
	}
}

func (*CustHandler) UpdateBusinessSignedInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessSignedInfoRequest{
				SignedId: container.GetValFromMapMaybe(params, "signed_no").ToString(), //
				Cycle:    container.GetValFromMapMaybe(params, "cycle").ToString(),     //
				Rate:     container.GetValFromMapMaybe(params, "rate").ToStringNoPointReg(`^[\d]+$`),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"), //
			}

			if errStr := verify.UpdateBusinessSignedInfoVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.UpdateBusinessSignedInfo(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) UpdateBusinessSignedStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessSignedStatusRequest{
				SignedId: container.GetValFromMapMaybe(params, "signed_no").ToString(), //
				Status:   container.GetValFromMapMaybe(params, "status").ToString(),    //
				Notes:    container.GetValFromMapMaybe(params, "notes").ToString(),     //
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),                //
			}

			if errStr := verify.UpdateBusinessSignedStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.UpdateBusinessSignedStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) UpdateBusinessSceneIdx() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessSceneIdxRequest{
				SceneNo: container.GetValFromMapMaybe(params, "scene_no").ToString(), //
				SwapIdx: container.GetValFromMapMaybe(params, "swap_idx").ToString(), //
			}

			reply, err := CustHandlerInst.Client.UpdateBusinessSceneIdx(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

/**
商家转账订单列表
*/
func (*CustHandler) GetBusinessTransferOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.ManagementGetBusinessTransferOrderListRequest{
				Page:         strext.ToStringNoPoint(params[0]),
				PageSize:     strext.ToStringNoPoint(params[1]),
				StartTime:    strext.ToStringNoPoint(params[2]),
				EndTime:      strext.ToStringNoPoint(params[3]),
				CurrencyType: strext.ToStringNoPoint(params[4]),
				LogNo:        strext.ToStringNoPoint(params[5]),
				OrderStatus:  strext.ToStringNoPoint(params[6]),
				ToAccount:    strext.ToStringNoPoint(params[7]),
				FromAccount:  strext.ToStringNoPoint(params[8]),
				BatchNo:      strext.ToStringNoPoint(params[9]),
				TransferType: strext.ToStringNoPoint(params[10]),
			}

			reply, err := CustHandlerInst.Client.ManagementGetBusinessTransferOrderList(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用服务失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "currency_type", "log_no", "order_status",
			"to_account", "from_account", "batch_no", "transfer_type")
	}
}

/**
商家转账批次列表
*/
func (*CustHandler) GetBusinessTransferBatchList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessTransferBatchListRequest{
				Page:         strext.ToStringNoPoint(params[0]),
				PageSize:     strext.ToStringNoPoint(params[1]),
				StartTime:    strext.ToStringNoPoint(params[2]),
				EndTime:      strext.ToStringNoPoint(params[3]),
				CurrencyType: strext.ToStringNoPoint(params[4]),
				BatchNo:      strext.ToStringNoPoint(params[5]),
				Status:       strext.ToStringNoPoint(params[6]),
				Account:      strext.ToStringNoPoint(params[7]),
			}

			reply, err := CustHandlerInst.Client.GetBusinessTransferBatchList(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用服务失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "currency_type", "batch_no", "status", "account")
	}
}

func (*CustHandler) GetBusinessAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessAccounts(context.TODO(), &custProto.GetBusinessAccountsRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Account:  strext.ToString(params[2]),
				SortType: strext.ToString(params[3]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "account", "sort_type")
	}
}

//商家账户收益
func (*CustHandler) GetBusinessAccountsProfit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessAccountsProfit(context.TODO(), &custProto.GetBusinessAccountsProfitRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Account:  strext.ToString(params[2]),
				SortType: strext.ToString(params[3]),
			})
			return reply.ResultCode, reply.List, reply.Total, err
		}, "page", "page_size", "account", "sort_type")
	}
}

//商家账户收益明细列表
func (*CustHandler) GetBusinessProfitList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessProfitList(context.TODO(), &custProto.GetBusinessProfitListRequest{
				Page:              strext.ToInt32(params[0]),
				PageSize:          strext.ToInt32(params[1]),
				StartTime:         strext.ToString(params[2]),
				EndTime:           strext.ToString(params[3]),
				Reason:            strext.ToString(params[4]),
				BusinessAccountNo: strext.ToString(params[5]),
				CurrencyType:      strext.ToString(params[6]),
				OpType:            strext.ToString(params[7]),
			})
			return reply.ResultCode, reply.List, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "reason", "business_account_no", "currency_type", "op_type")
	}
}

func (*CustHandler) GetBusinessBillList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessBillListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToString(params[2]),
				EndTime:      strext.ToString(params[3]),
				Reason:       strext.ToString(params[4]),
				Uid:          strext.ToString(params[5]),
				LogNo:        strext.ToString(params[6]),
				CurrencyType: strext.ToString(params[7]),
				VaType:       strext.ToString(params[8]),
			}

			if req.Uid == "" {
				ss_log.Error("参数Uid为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessBillList(context.TODO(), req)

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "reason", "uid", "log_no", "currency_type", "va_type")
	}

}

//获取商家修改认证材料列表
func (*CustHandler) GetAuthMaterialBusinessUpdateList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAuthMaterialBusinessUpdateList(context.TODO(), &custProto.GetAuthMaterialBusinessUpdateListRequest{
				Page:           strext.ToInt32(params[0]),
				PageSize:       strext.ToInt32(params[1]),
				StartTime:      strext.ToString(params[2]),
				EndTime:        strext.ToString(params[3]),
				Account:        strext.ToString(params[4]),
				AuthMaterialNo: strext.ToString(params[5]),
				Status:         strext.ToString(params[6]),
			})

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "account", "auth_material_no", "status")
	}
}

func (*CustHandler) ModifyAuthMaterialBusinessUpdateStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.ModifyAuthMaterialBusinessUpdateStatusRequest{
				Id:     container.GetValFromMapMaybe(params, "id").ToString(),
				Status: container.GetValFromMapMaybe(params, "status").ToString(),
				Notes:  container.GetValFromMapMaybe(params, "notes").ToString(),
			}

			if errStr := verify.ModifyAuthMaterialBusinessUpdateStatusVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.ModifyAuthMaterialBusinessUpdateStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetBusinessIndustryRateCycleList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessIndustryRateCycleList(context.TODO(), &custProto.GetBusinessIndustryRateCycleListRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				Code:      strext.ToStringNoPoint(params[4]),

				UpCode:       strext.ToStringNoPoint(params[5]),
				IndustryName: strext.ToStringNoPoint(params[6]),
				ChannelName:  strext.ToStringNoPoint(params[7]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "code", "up_code", "industry_name", "channel_name")
	}
}

func (*CustHandler) InsertOrUpdateBusinessIndustryRateCycle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateBusinessIndustryRateCycleRequest{
				Id:                container.GetValFromMapMaybe(params, "id").ToString(),         //
				Code:              container.GetValFromMapMaybe(params, "code").ToString(),       //
				BusinessChannelNo: container.GetValFromMapMaybe(params, "channel_no").ToString(), //
				Rate:              container.GetValFromMapMaybe(params, "rate").ToString(),       //
				Cycle:             container.GetValFromMapMaybe(params, "cycle").ToString(),      //
				LoginUid:          inner_util.GetJwtDataString(c, "account_uid"),                 //
			}

			if errStr := verify.InsertOrUpdateBusinessIndustryRateCycleVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateBusinessIndustryRateCycle(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) DelBusinessIndustryRateCycle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DelBusinessIndustryRateCycleRequest{
				Id:       container.GetValFromMapMaybe(params, "id").ToString(), //
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),         //
			}

			if req.Id == "" {
				ss_log.Error("必要参数Id为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DelBusinessIndustryRateCycle(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

//获取核销码列表
func (*CustHandler) GetWriteOffCodeList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetWriteOffList(context.TODO(), &custProto.GetWriteOffListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToString(params[2]),
				EndTime:      strext.ToString(params[3]),
				Code:         strext.ToString(params[4]),
				PayerAccount: strext.ToString(params[5]),
				PayeeAccount: strext.ToString(params[6]),
				UseStatus:    strext.ToString(params[7]),
			})
			return reply.ResultCode, reply.List, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "code", "payer_account", "payee_account", "use_status")
	}
}

//处理核销码（freeze冻结、unfreeze解冻、cancel注销）
func (*CustHandler) DisposeWriteOffCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.DisposeWriteOffCode(context.TODO(), &custProto.DisposeWriteOffCodeRequest{
				Code:     container.GetValFromMapMaybe(params, "code").ToStringNoPoint(),
				OpType:   container.GetValFromMapMaybe(params, "op_type").ToStringNoPoint(),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, nil, err
		})
	}
}

func (*CustHandler) UpdateBusinessSceneSignedStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessSceneSignedStatusRequest{
				SignedId: container.GetValFromMapMaybe(params, "signed_no").ToString(), //
				Status:   container.GetValFromMapMaybe(params, "status").ToString(),    //
				Notes:    container.GetValFromMapMaybe(params, "notes").ToString(),     //
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"),                //
			}

			if errStr := verify.UpdateBusinessSceneSignedStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.UpdateBusinessSceneSignedStatus(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) UpdateBusinessSceneSignedInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.UpdateBusinessSceneSignedInfoRequest{
				SignedId: container.GetValFromMapMaybe(params, "signed_no").ToString(), //
				Cycle:    container.GetValFromMapMaybe(params, "cycle").ToString(),     //
				Rate:     container.GetValFromMapMaybe(params, "rate").ToStringNoPointReg(`^[\d]+$`),
				LoginUid: inner_util.GetJwtDataString(c, "account_uid"), //
			}

			if errStr := verify.UpdateBusinessSceneSignedInfoVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.UpdateBusinessSceneSignedInfo(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

func (*CustHandler) GetBusinessSceneSignedList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessSceneSignedList(context.TODO(), &custProto.GetBusinessSceneSignedListRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]),
				EndTime:   strext.ToStringNoPoint(params[3]),
				Status:    strext.ToStringNoPoint(params[4]),
				Lang:      constants.LangZhCN,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "status")
	}
}
func (*CustHandler) GetApiPayLogList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetApiPayLogList(context.TODO(), &custProto.GetApiPayLogListRequest{
				Page:      strext.ToStringNoPoint(params[0]),
				PageSize:  strext.ToStringNoPoint(params[1]),
				StartTime: strext.ToStringNoPoint(params[2]), //创建时间-开始
				EndTime:   strext.ToStringNoPoint(params[3]), //创建时间-结算
				ReqMethod: strext.ToStringNoPoint(params[4]),
				ReqUri:    strext.ToStringNoPoint(params[5]),

				ReqBody:        strext.ToStringNoPoint(params[6]),
				RespData:       strext.ToStringNoPoint(params[7]),
				TrafficStatus:  strext.ToStringNoPoint(params[8]), //通信状态(0失败，1成功)
				BusinessStatus: strext.ToStringNoPoint(params[9]), //业务处理装(0失败，1成功)
				AppId:          strext.ToStringNoPoint(params[10]),

				ReqStartTime: strext.ToStringNoPoint(params[11]), //请求时间开始
				ReqEndTime:   strext.ToStringNoPoint(params[12]), //请求时间结束
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "req_method",
			"req_uri", "req_body", "resp_data", "traffic_status", "business_status",
			"app_id", "req_start_time", "req_end_time")
	}
}
