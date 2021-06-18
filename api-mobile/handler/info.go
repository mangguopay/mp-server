package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/inner_util"
	"a.a/mp-server/api-mobile/verify"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

/**
 * 获取功能列表
 */
func (s *AuthHandler) GetFuncList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetFuncList(context.TODO(), &go_micro_srv_cust.GetFuncListRequest{
				// 功能类型
				ApplicationType: strext.ToInt32(params[0]),
			})

			return reply.ResultCode, reply.Datas, 0, err
		}, "application_type")
	}
}

/**
 * 获取用户余额
 */
func (s *AuthHandler) GetRemain() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			uid := inner_util.GetJwtDataString(c, "account_uid")
			if uid == "" {
				ss_log.Error("err=[GetRemain 接口------>%s]", " account_uid  为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}
			reply, err := s.Client.GetRemain(context.TODO(), &go_micro_srv_auth.GetRemainRequest{
				// 账号
				AccountNo: uid,
			})
			return reply.ResultCode, reply.Data, 0, err
		})
	}
}

// 获取pos机余额
func (s *AuthHandler) GetPosRemain() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.GetPosRemainRequest{
				// 账号
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			}
			if errStr := verify.GetPosRemainReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := s.Client.GetPosRemain(context.TODO(), req)
			return reply.ResultCode, reply.Data, err
		})
	}
}

/**
 * 获取兑换费率
 */
func (s *AuthHandler) GetExchangeRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetExchangeRate(context.TODO(), &go_micro_srv_cust.GetExchangeRateRequest{})
			return reply.ResultCode, reply.Datas, 0, err
		})
	}
}

/**
 * 获取服务商信息
 */
func (s *AuthHandler) GetServicer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServicer(context.TODO(), &go_micro_srv_cust.GetServicerRequest{
				// 服务商id
				IdenNo:     inner_util.GetJwtDataString(c, "iden_no"),
				ServicerNo: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, gin.H{
				"servicer_no":    reply.Data.ServicerNo,
				"contact_person": reply.Data.ContactPerson,
				"contact_phone":  reply.Data.ContactPhone,
				"contact_addr":   reply.Data.ContactAddr,
				"addr":           reply.Data.Addr,
				"servicer_name":  reply.Data.ServicerName, //网点名称
				"business_time":  reply.Data.BusinessTime, //网点营业时间
				"lat":            reply.Data.Lat,          //纬度
				"lng":            reply.Data.Lng,          //经度
			}, 0, err
		}, "servicer_no")
	}
}

/**
 * 获取总部卡列表
 */
func (s *AuthHandler) GetHeadquartersCards() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &go_micro_srv_cust.GetHeadquartersCardsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				MoneyType:   strext.ToStringNoPoint(params[2]),
				AccountType: inner_util.GetJwtDataString(c, "account_type"), //3服务商，4用户
			}
			if errStr := verify.GetHeadquartersCardsReqVerify(req); errStr != "" {
				return errStr, nil, 0, nil
			}
			reply, err := CustHandlerInst.Client.GetHeadquartersCards(c, req)

			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "money_type", "account_type")
	}
}

/**
 * 获取服务商或者用户的卡列表
 */
func (s *AuthHandler) GetServicerOrCustCards() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetUserCards(context.TODO(), &go_micro_srv_cust.GetUserCardsRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				BalanceType: strext.ToString(params[0]),
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "balance_type")
	}

}

/**
 * 获取服务商或者用户的银行卡详细信息
 */
func (s *AuthHandler) GetServicerOrCustCardDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetUserCardDetail(context.TODO(), &go_micro_srv_cust.GetUserCardDetailRequest{
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				CardNo:      strext.ToString(params[0]),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, nil
		}, "card_no")
	}
}

func (s *AuthHandler) ModifyCardsDefalut() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.ModifyCardsDefalut(context.TODO(), &go_micro_srv_cust.ModifyCardsDefalutRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				IsDefalut:   container.GetValFromMapMaybe(params, "is_defalut").ToStringNoPoint(),
				BalanceType: container.GetValFromMapMaybe(params, "balance_type").ToStringNoPoint(),
			})
			return reply.ResultCode, 0, err
		})
	}
}

//app用户账单
func (s *AuthHandler) CustBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.CustBills(context.TODO(), &go_micro_srv_cust.CustBillsRequest{
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
				QueryTime:    strext.ToString(params[0]),
				CurrencyType: strext.ToString(params[1]),
				Page:         strext.ToInt32(params[2]),
				PageSize:     strext.ToInt32(params[3]),
			})
			return reply.ResultCode, gin.H{
				"datas":        reply.DataList,
				"incom_sum":    reply.IncomeSum,
				"spending_sum": reply.SpendingSum,
			}, reply.Total, err
		}, "query_time", "currency_type", "page", "page_size")
	}
}

// // pos端获取交易明细 post服务商账单
func (s *AuthHandler) ServicerBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {

			reply, err := CustHandlerInst.Client.ServicerBills(context.TODO(), &go_micro_srv_cust.ServicerBillsRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				//EndTime:      strext.ToStringNoPoint(params[0]),
				CurrencyType: strext.ToStringNoPoint(params[3]),
				BillType:     strext.ToStringNoPoint(params[4]),
				OpNo:         strext.ToStringNoPoint(params[5]),
			})
			return reply.ResultCode, gin.H{
				"datas": reply.Datas,
			}, reply.Total, err
		}, "page", "page_size", "start_time", "money_type", "bill_type", "op_no")
	}
}

//pos获取收款额度
func (s *AuthHandler) ServicerCollectLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.ServicerCollectLimit(context.TODO(), &go_micro_srv_cust.ServicerCollectLimitRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			})
			return reply.ResultCode, gin.H{
				"usd_auth_collect_limit":     reply.Datas.UsdAuthCollectLimit,
				"usd_no_spent_collect_limit": reply.Datas.UsdNoSpentCollectLimit,
				"khr_auth_collect_limit":     reply.Datas.KhrAuthCollectLimit,
				"khr_no_spent_collect_limit": reply.Datas.KhrNoSpentCollectLimit,
			}, 0, err
		})
	}
}

func (s *AuthHandler) AddRecvCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_cust.InsertOrUpdateCardRequest{
				CardNo:        container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				AccountNo:     inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:   inner_util.GetJwtDataString(c, "account_type"),
				ChannelNo:     container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				CardAccName:   container.GetValFromMapMaybe(params, "card_acc_name").ToStringNoPoint(),
				CardNumber:    container.GetValFromMapMaybe(params, "card_number").ToStringNoPointReg(`^[\d]+$`),
				BalanceType:   container.GetValFromMapMaybe(params, "balance_type").ToStringNoPoint(),
				IsDefault:     container.GetValFromMapMaybe(params, "is_default").ToStringNoPoint(),
				CollectStatus: "1",
				AuditStatus:   "1",
			}
			if req.CardNumber == "" {
				ss_log.Error("卡号不符合规则或为空")
				return ss_err.ERR_Card_Number_FAILD, nil, nil
			}
			reply, err := CustHandlerInst.Client.InsertOrUpdateCard(context.TODO(), req)

			return reply.ResultCode, 0, err
		})
	}
}

func (s *AuthHandler) GetPosChannelList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetPosChannelList(context.TODO(), &go_micro_srv_cust.GetPosChannelListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				//RoleType:     inner_util.GetJwtDataString(c, "account_type"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "currency_type")
	}
}

// pos机新增银行卡
func (s *AuthHandler) AddCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.AddCardRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				RecCarNum:   container.GetValFromMapMaybe(params, "rec_car_num").ToStringNoPoint(),
				RecName:     container.GetValFromMapMaybe(params, "rec_name").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
				ChannelName: container.GetValFromMapMaybe(params, "channel_name").ToStringNoPoint(),
				IsDefault:   container.GetValFromMapMaybe(params, "is_default").ToInt32(),
			}

			// 参数校验
			if errStr := verify.AddCardReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := s.Client.AddCard(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

// 修改为默认卡
func (s *AuthHandler) ModifyDefaultCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.ModifyDefaultCardRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
			}
			// 参数校验
			if errStr := verify.ModifyDefaultCardReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.ModifyDefaultCard(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

// 解绑银行卡
func (s *AuthHandler) DeleteBindCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.DeleteBindCardRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				CarNum:      container.GetValFromMapMaybe(params, "car_num").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
				ChannelName: container.GetValFromMapMaybe(params, "channel_name").ToStringNoPoint(),
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
			}

			if errStr := verify.DeleteBindCardReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.DeleteBindCard(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

// 计算和查询费率
func (s *AuthHandler) QueryRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {

			//................
			idenNo := inner_util.GetJwtDataString(c, "iden_no")
			accountType := inner_util.GetJwtDataString(c, "account_type")
			pwdReq := &go_micro_srv_bill.QueryCustHasPwdRequest{
				AccountType: accountType,
				IdenNo:      idenNo,
			}

			pwdReply, err := BillHandlerInst.Client.QueryCustHasPwd(context.TODO(), pwdReq)
			if pwdReply.ResultCode != ss_err.ERR_SUCCESS {
				return pwdReply.ResultCode, nil, err
			}

			req := &go_micro_srv_bill.QeuryRateRequest{
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				Type:   container.GetValFromMapMaybe(params, "type").ToInt32(),

				AccountType: accountType,

				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),

				IdenNo: idenNo,
			}

			if errStr := verify.QueryRateReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.QeuryFees(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "")
	}
}
func (s *AuthHandler) CustQueryRate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {

			//................
			idenNo := inner_util.GetJwtDataString(c, "iden_no")
			accountType := inner_util.GetJwtDataString(c, "account_type")
			pwdReq := &go_micro_srv_bill.QueryCustHasPwdRequest{
				AccountType: accountType,
				IdenNo:      idenNo,
			}

			pwdReply, err := BillHandlerInst.Client.QueryCustHasPwd(context.TODO(), pwdReq)
			if pwdReply.ResultCode != ss_err.ERR_SUCCESS {
				return pwdReply.ResultCode, nil, err
			}

			req := &go_micro_srv_bill.CustQeuryRateRequest{
				Amount:    container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				OpType:    container.GetValFromMapMaybe(params, "op_type").ToInt32(),
				ChannelNo: container.GetValFromMapMaybe(params, "channel_no").ToStringNoPoint(),
				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),

				AccountType: accountType,

				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),

				//IdenNo: idenNo,
			}

			if errStr := verify.CustQueryRateReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.CustQeuryFees(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "")
	}
}

// 查询最大最小限额
func (s *AuthHandler) QueryMinMaxAmount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {

			req := &go_micro_srv_bill.QueryMinMaxAmountRequest{
				Type:      container.GetValFromMapMaybe(params, "type").ToStringNoPoint(),       //1-手机号取款;2-存款;3-转正;4-扫码取款
				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(), // usd;khr

			}

			reply, err := BillHandlerInst.Client.QueryMinMaxAmount(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "")
	}
}

func (s *AuthHandler) GetServicerBillingDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetServicerBillingDetails(context.TODO(), &go_micro_srv_cust.GetServicerBillingDetailsRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:  inner_util.GetJwtDataString(c, "account_type"),
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
			})

			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "money_type")
	}
}

// 手机号,扫一扫取款打印小票查询
func (s *AuthHandler) WithdrawReceipt() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_bill.WithdrawReceiptRequest{
				OrderNo: strext.ToString(params),
			}

			if errStr := verify.WithdrawReceiptReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := BillHandlerInst.Client.WithdrawReceipt(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "order_no")
	}
}

func (s *AuthHandler) QuerySaveReceipt() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_bill.QuerySaveReceiptRequest{
				OrderNo: strext.ToString(params),
			}

			if errStr := verify.QuerySaveReceiptReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := BillHandlerInst.Client.QuerySaveReceipt(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "order_no")
	}
}

// 扫码提现后的订单详情
func (s *AuthHandler) SweepWithdrawDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_bill.SweepWithdrawDetailRequest{
				OrderNo: strext.ToString(params),
			}

			if errStr := verify.SweepWithdrawDetailReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.SweepWithdrawDetail(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "order_no")
	}
}

// 存款后的订单详情
func (s *AuthHandler) SaveMoneyDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_bill.SaveMoneyDetailRequest{
				OrderNo: strext.ToString(params),
			}

			if errStr := verify.SaveMoneyDetailReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.SaveMoneyDetail(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "order_no")
	}
}

func (s *AuthHandler) SaveWithdrawDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_bill.SaveDetailRequest{
				OrderNo: container.GetValFromMapMaybe(params, "order_no").ToStringNoPoint(),
				Type:    container.GetValFromMapMaybe(params, "type").ToInt32(),
			}

			if errStr := verify.SaveWithdrawDetailReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.SaveDetail(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "params")
	}
}

// 兑换金额查询
func (s *AuthHandler) ExchangeAmount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_bill.ExchangeAmountRequest{
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				Type:   container.GetValFromMapMaybe(params, "type").ToInt32(), // 1-usd-->khr;2-khr-->usd
			}

			reply, err := BillHandlerInst.Client.ExchangeAmount(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"amount": reply.Amount,
			}, err
		}, "params")
	}
}

func (s *AuthHandler) MyData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			reply, err := s.Client.MyData(context.TODO(), &go_micro_srv_auth.MyDataRequest{
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			})
			return reply.ResultCode, gin.H{
				"open_idx":             reply.OpenIdx,
				"contact_person":       reply.ContactPerson,
				"contact_phone":        reply.ContactPhone,
				"contact_addr":         reply.ContactAddr,
				"addr":                 reply.Addr,
				"income_authorization": reply.IncomeAuthorization,
				"outgo_authorization":  reply.OutgoAuthorization,
				"create_time":          reply.CreateTime,
			}, err
		}, "params")
	}
}

// pos端获取账单列表
func (s *AuthHandler) GetServicerCheckList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.GetServicerCheckList(context.TODO(), &go_micro_srv_bill.GetServicerCheckListRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				StartTime:   strext.ToStringNoPoint(params[0]),
				EndTime:     strext.ToStringNoPoint(params[1]),
				Page:        strext.ToInt32(params[2]),
				PageSize:    strext.ToInt32(params[3]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "start_time", "end_time", "page", "page_size")
	}
}

func (s *AuthHandler) GetServicerProfitLedgers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.GetServicerProfitLedgers(context.TODO(), &go_micro_srv_bill.GetServicerProfitLedgersRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:  inner_util.GetJwtDataString(c, "account_type"),
				StartTime:    strext.ToStringNoPoint(params[0]),
				EndTime:      strext.ToStringNoPoint(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				Page:         strext.ToInt32(params[3]),
				PageSize:     strext.ToInt32(params[4]),
			})
			return reply.ResultCode, gin.H{
				"datas":         reply.Datas,
				"khr_count_sum": reply.KhrCountSum,
				"usd_count_sum": reply.UsdCountSum,
			}, reply.Total, err
		}, "start_time", "end_time", "currency_type", "page", "page_size")
	}
}

func (s *AuthHandler) GetServicerProfitLedgerDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.GetServicerProfitLedgerDetail(context.TODO(), &go_micro_srv_bill.GetServicerProfitLedgerDetailRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				LogNo:       params[0],
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "log_no")
	}
}

func (s *AuthHandler) GetLogAppMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetLogAppMessages(context.TODO(), &go_micro_srv_cust.GetLogAppMessagesRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *AuthHandler) CheckPayPWD() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.CheckPayPWDRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),

				NonStr: container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
			}

			if errStr := verify.CheckPayPWDReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := s.Client.CheckPayPWD(context.TODO(), req)
			return reply.ResultCode, "", err
		})
	}
}

func (s *AuthHandler) RealTimeCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_bill.RealTimeCountRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			}

			reply, err := BillHandlerInst.Client.RealTimeCount(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"datas": reply.Data,
			}, err
		}, "params")
	}
}

func (s *AuthHandler) CustPayment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type") // 账号类型(用户才可以支付)
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, nil
			}

			req := &go_micro_srv_auth.CustPaymentRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				CustNo:      inner_util.GetJwtDataString(c, "iden_no"),
				AccountType: accountType,
			}
			reply, err := s.Client.CustPayment(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"def_pay_no": reply.DefPayNo,
				"data":       reply.Data,
			}, err
		}, "params")
	}
}

// pos机围栏
func (s *AuthHandler) PosFence() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			lat := container.GetValFromMapMaybe(params, "lat").ToStringNoPoint()
			lng := container.GetValFromMapMaybe(params, "lng").ToStringNoPoint()
			lat = fmt.Sprintf("%.8f", strext.ToFloat64(lat))
			lng = fmt.Sprintf("%.8f", strext.ToFloat64(lng))
			req := &go_micro_srv_cust.PosFenceRequest{
				PosSn: container.GetValFromMapMaybe(params, "pos_sn").ToStringNoPoint(),
				Lat:   lat,
				Lng:   lng,
			}
			// 参数校验
			if errStr := verify.PosFenceReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := CustHandlerInst.Client.PosFence(context.TODO(), req)
			return reply.ResultCode, "", err
		})
	}
}

// 切换语言获得新 token
//func (s *AuthHandler) CheckLang() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
//			req := &go_micro_srv_auth.CheckLangRequest{
//				Account:        inner_util.GetJwtDataString(c, "account"),
//				AccountUid:     inner_util.GetJwtDataString(c, "account_uid"),
//				IdenNo:         inner_util.GetJwtDataString(c, "iden_no"),
//				AccountType:    inner_util.GetJwtDataString(c, "account_type"),
//				AccountName:    inner_util.GetJwtDataString(c, "account_name"),
//				LoginAccountNo: inner_util.GetJwtDataString(c, "login_account_no"),
//				PubKey:         inner_util.GetJwtDataString(c, "pub_key"),
//				JumpIdenNo:     inner_util.GetJwtDataString(c, "jump_iden_no"),
//				JumpIdenType:   inner_util.GetJwtDataString(c, "jump_iden_type"),
//				MasterAccNo:    inner_util.GetJwtDataString(c, "master_acc_no"),
//				IsMasterAcc:    inner_util.GetJwtDataString(c, "is_master_acc"),
//			}
//			reply, err := AuthHandlerInst.Client.CheckLang(context.TODO(), req)
//			return reply.ResultCode, gin.H{
//				"userinfo": gin.H{
//					"token": reply.Jwt,
//				},
//			}, err
//		})
//	}
//}

// 修改用户默认卡
func (s *AuthHandler) UpdateDefCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.UpdateDefCard(context.TODO(), &go_micro_srv_cust.UpdateDefCardRequest{
				CustNo:   inner_util.GetJwtDataString(c, "iden_no"),
				DefPayNo: container.GetValFromMapMaybe(params, "def_pay_no").ToStringNoPoint(),
			})
			return reply.ResultCode, "", err
		})
	}
}

func (s *AuthHandler) GetSerPos() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetSerPos(context.TODO(), &go_micro_srv_cust.GetSerPosRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

// 修改pos状态
func (s *AuthHandler) ModifySerPosStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.ModifySerPosStatus(context.TODO(), &go_micro_srv_cust.ModifySerPosStatusRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				TerminalNo:  container.GetValFromMapMaybe(params, "terminal_no").ToStringNoPoint(),
				UseStatus:   container.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
			})
			return reply.ResultCode, "", err
		})
	}
}

func (s *AuthHandler) GetLogAppMessagesCnt() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			reply, err := BillHandlerInst.Client.GetLogAppMessagesCnt(context.TODO(), &go_micro_srv_bill.GetLogAppMessagesCntRequest{
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			if err != nil {
				ss_log.Error("err=[%v]", err)
				return ss_err.ERR_PARAM, gin.H{
					"total": "0",
				}, err
			}
			return reply.ResultCode, gin.H{
				"total": reply.Total,
			}, err
		}, "params")
	}
}

// 修改pos状态
func (s *AuthHandler) ModifyAppMessagesIsRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := CustHandlerInst.Client.ModifyAppMessagesIsRead(context.TODO(), &go_micro_srv_cust.ModifyAppMessagesIsReadRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, "", err
		})
	}
}

func (s *AuthHandler) GetUseChannelList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetUseChannelList(context.TODO(), &go_micro_srv_cust.GetUseChannelListRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				//RoleType:     inner_util.GetJwtDataString(c, "account_type"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "currency_type")
	}
}
