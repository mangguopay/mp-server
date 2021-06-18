package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/inner_util"
	"a.a/mp-server/api-mobile/verify"
	"a.a/mp-server/common/constants"
	authProto "a.a/mp-server/common/proto/auth"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

type Business struct {
}

var BusinessInst Business

//个人商家基本资料
func (Business) GetPersonalBusinessInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(param interface{}) (string, gin.H, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				return ss_err.ERR_PERMISSION_DENIED, nil, nil
			}

			reply, err := CustHandlerInst.Client.GetPersonalBusinessInfo(context.TODO(), &custProto.GetPersonalBusinessInfoRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
			})
			if err != nil {
				ss_log.Error("cust-srv.GetPersonalBusinessInfo()调用失败, err=%v", err)
				return ss_err.ERR_SYSTEM, nil, nil
			}

			data := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				data = gin.H{
					"account":              reply.Data.Account,
					"real_name":            reply.Data.RealName,
					"business_name":        reply.Data.BusinessName,
					"simplify_name":        reply.Data.SimplifyName,
					"business_auth_status": reply.Data.BusinessAuthStatus,
					"operating_period":     reply.Data.OperatingPeriod,
					"organization_code":    reply.Data.OrganizationCode,
					"license_img":          reply.Data.LicenseImg,
				}
			}

			return reply.ResultCode, data, nil
		}, "")
	}
}

//个人商家余额
func (Business) GetPersonalBusinessBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				return ss_err.ERR_PERMISSION_DENIED, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetPersonalBusinessBalance(context.TODO(), &custProto.GetPersonalBusinessBalanceRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
			})
			if err != nil {
				ss_log.Error("cust-srv.GetPersonalBusinessBalance()调用失败, err=%v", err)
				return ss_err.ERR_SYSTEM, nil, 0, nil
			}

			return reply.ResultCode, reply.Data, 2, nil
		})
	}
}

//个人商家收款订单
func (Business) GetPersonalBusinessBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				return ss_err.ERR_PERMISSION_DENIED, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetPersonalBusinessBills(context.TODO(), &custProto.GetPersonalBusinessBillsRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToStringNoPoint(params[2]),
				EndTime:      strext.ToStringNoPoint(params[3]),
				CurrencyType: strext.ToStringNoPoint(params[4]),
				OrderStatus:  strext.ToStringNoPoint(params[5]),
			})
			if err != nil {
				ss_log.Error("cust-srv.GetPersonalBusinessBalance()调用失败, err=%v", err)
				return ss_err.ERR_SYSTEM, nil, 0, nil
			}

			return reply.ResultCode, reply.List, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "currency_type", "order_status")
	}
}

//收款订单详情
func (Business) GetPersonalBusinessBillDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				return ss_err.ERR_PERMISSION_DENIED, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetPersonalBusinessBillDetail(context.TODO(), &custProto.GetPersonalBusinessBillDetailRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				OrderNo:   strext.ToStringNoPoint(params[0]),
			})
			if err != nil {
				ss_log.Error("cust-srv.GetPersonalBusinessBalance()调用失败, err=%v", err)
				return ss_err.ERR_SYSTEM, nil, 0, nil
			}

			return reply.ResultCode, reply.Data, 1, nil
		}, "order_no")
	}
}

//个人商家固定收款码
func (Business) GetPersonalBusinessFixedCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				return ss_err.ERR_PERMISSION_DENIED, nil, nil
			}

			reply, err := CustHandlerInst.Client.GetPersonalBusinessFixedCode(context.TODO(), &custProto.GetPersonalBusinessFixedCodeRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: accountType,
			})
			if err != nil {
				ss_log.Error("cust-srv.GetPersonalBusinessFixedCode()调用失败, err=%v", err)
				return ss_err.ERR_SYSTEM, nil, nil
			}

			data := gin.H{
				"fixed_code":    reply.FixedCode,
				"simplify_name": reply.SimplifyName,
			}
			return reply.ResultCode, data, nil
		}, "")
	}
}

//个人商家收款下单
func (Business) PersonalBusinessPerPay() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type") // 账号类型(用户才可以支付)
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, nil
			}

			subject := inner_util.M(c, "subject") // 商品名称
			if subject == "" {
				lang := ss_net.GetCommonData(c).Lang
				switch lang {
				case constants.LangZhCN:
					subject = "其它"
				case constants.LangEnUS:
					subject = "other"
				case constants.LangKmKH:
					subject = "ផ្សេងៗ"
				default:
					subject = "other"
				}
			}
			req := &businessBillProto.PersonalBusinessPrepayRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				Subject:      subject,                          // 商品名称
				Amount:       inner_util.M(c, "amount"),        // 金额
				CurrencyType: inner_util.M(c, "currency_type"), // 币种
				Remark:       inner_util.M(c, "remark"),        // 备注
				Lang:         ss_net.GetCommonData(c).Lang,
			}

			reply, err := BusinessBillHandlerInst.Client.PersonalBusinessPrepay(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no":   reply.OrderNo,
					"qr_code_id": reply.QrCodeId,
				}
			}

			return reply.ResultCode, respData, nil
		})
	}
}

func (s *Business) GetChannelList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessChannels(context.TODO(), &custProto.GetBusinessChannelsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				ChannelName:  strext.ToStringNoPoint(params[3]),
				UseStatus:    constants.Status_Enable, //只要启用的
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "money_type", "channel_name")
	}
}

func (*Business) GetBusinessCards() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetBusinessCards(context.TODO(), &custProto.GetBusinessCardsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				BalanceType: strext.ToStringNoPoint(params[2]),
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: constants.AccountType_PersonalBusiness,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			for _, data := range reply.DataList {
				//银行卡列表的卡号需要做脱敏处理,只要后4位
				if len(data.CardNumber) > 4 {
					data.CardNumber = data.CardNumber[len(data.CardNumber)-4:]
				}
			}

			return reply.ResultCode, reply.DataList, reply.Total, nil
		}, "page", "page_size", "money_type")
	}

}

func (*Business) GetBusinessCardDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			//app商家服务这边需要验证用户支付密码才能看到详情
			replyCheckPwd, errCheckPwd := AuthHandlerInst.Client.CheckPayPWD(context.TODO(), &authProto.CheckPayPWDRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				Password:    container.GetValFromMapMaybe(params, "pass_word").ToStringNoPoint(),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
			})

			if errCheckPwd != nil {
				ss_log.Error("api调用失败,err=[%v]", errCheckPwd)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if replyCheckPwd.ResultCode != ss_err.ERR_SUCCESS {
				payPasswordErrTips := ""
				if replyCheckPwd.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
					payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, replyCheckPwd.ResultCode, replyCheckPwd.ErrTips)
					replyCheckPwd.ResultCode = ss_err.ERR_DB_PWD
				}
				ss_log.Info("payPasswordErrTips-------------------", payPasswordErrTips)
				c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)
				return replyCheckPwd.ResultCode, nil, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessCardDetail(context.TODO(), &custProto.GetBusinessCardDetailRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: constants.AccountType_PersonalBusiness,
				CardNo:      strext.ToStringNoPoint(params),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, nil
		}, "card_no")
	}

}

func (*Business) AddBusinessCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.InsertOrUpdateBusinessCardRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: constants.AccountType_PersonalBusiness,
				ChannelId:   container.GetValFromMapMaybe(params, "channel_id").ToStringNoPoint(),
				Name:        container.GetValFromMapMaybe(params, "name").ToStringNoPoint(),
				CardNumber:  container.GetValFromMapMaybe(params, "card_number").ToStringNoPointReg(`^[\d]+$`),
				IsDefault:   container.GetValFromMapMaybe(params, "is_default").ToStringNoPoint(),
			}

			if errStr := verify.InsertOrUpdateBusinessCardVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			reply, err := CustHandlerInst.Client.InsertOrUpdateBusinessCard(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (*Business) DelBusinessCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.DelBusinessCardRequest{
				CardNo: container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
			}
			if req.CardNo == "" {
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.DelBusinessCard(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			return reply.ResultCode, 0, err
		})
	}
}

/**
 * 获取总部卡列表
 */
func (s *Business) GetHeadCards() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetHeadquartersCardsRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				MoneyType:   strext.ToStringNoPoint(params[2]),
				AccountType: constants.AccountType_PersonalBusiness,
			}
			if errStr := verify.GetHeadquartersCardsReqVerify(req); errStr != "" {
				return errStr, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetHeadquartersCards(c, req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "money_type")
	}
}

func (*Business) AddBusinessToHead() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &custProto.AddIndividualBusinessToHeadRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				Amount:      container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				PayPwd:      container.GetValFromMapMaybe(params, "pay_pwd").ToStringNoPoint(),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				ImageId:     container.GetValFromMapMaybe(params, "image_id").ToStringNoPoint(),
			}

			if errStr := verify.AddBusinessToHeadVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			var reply, err = CustHandlerInst.Client.AddIndividualBusinessToHead(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}

			payPasswordErrTips := ""
			if reply.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
				payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, reply.ResultCode, reply.PayPasswordErrTips)
				reply.ResultCode = ss_err.ERR_DB_PWD
			}
			ss_log.Info("payPasswordErrTips-------------------", payPasswordErrTips)
			c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)

			return reply.ResultCode, gin.H{
				"log_no": reply.LogNo,
			}, nil
		})
	}
}

func (*Business) GetBusinessToHeadList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessToHeadListRequest{
				Page:        strext.ToStringNoPoint(params[0]),
				PageSize:    strext.ToStringNoPoint(params[1]),
				StartTime:   strext.ToStringNoPoint(params[2]),
				EndTime:     strext.ToStringNoPoint(params[3]),
				MoneyType:   strext.ToStringNoPoint(params[4]),
				LogNo:       strext.ToStringNoPoint(params[5]),
				OrderStatus: strext.ToStringNoPoint(params[6]),
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
			}
			if req.AccountNo == "" {
				ss_log.Error("AccountNo参数为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessToHeadList(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Datas, reply.Total, nil
		}, "page", "page_size", "start_time", "end_time", "money_type", "log_no", "order_status")
	}
}

func (*Business) GetBusinessToHeadDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &custProto.GetBusinessToHeadDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			}

			if req.LogNo == "" {
				ss_log.Error("LogNo为空")
				return ss_err.ERR_PARAM, nil, 0, nil
			}

			reply, err := CustHandlerInst.Client.GetBusinessToHeadDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.Data, 0, nil
		}, "log_no")
	}
}
