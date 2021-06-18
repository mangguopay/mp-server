package handler

import (
	"context"

	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/inner_util"
	"a.a/mp-server/common/constants"
	businessBillProto "a.a/mp-server/common/proto/business-bill"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"

	"github.com/gin-gonic/gin"
)

type BusinessBillHandler struct {
	Client businessBillProto.BusinessBillService
}

var BusinessBillHandlerInst BusinessBillHandler

//获取付款码
func (*BusinessBillHandler) QueryPaymentCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, 0, nil
			}

			req := &businessBillProto.GetPaymentCodeRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				Lang:      ss_net.GetCommonData(c).Lang,
			}
			ss_log.Info("api-mobile.QueryPaymentCode请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.GetPaymentCode(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用business_bill_srv.QueryPaymentCode失败, err=%v", err)
				return ss_err.ERR_SYS_NETWORK, nil, 0, err
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"payment_code": reply.PaymentCode,
				}
			}
			return reply.ResultCode, respData, 1, nil
		})
	}
}

//查询用户待支付订单
func (*BusinessBillHandler) QueryPendingPayOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, 0, nil
			}

			req := &businessBillProto.QueryPendingPayOrderRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				PaymentCode: params[0],
				Lang:        ss_net.GetCommonData(c).Lang,
			}
			ss_log.Info("api-mobile.QueryPendingPayOrder请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.QueryPendingPayOrder(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用business_bill_srv.QueryPendingPayOrder失败, err=%v", err)
				return ss_err.ERR_SYS_NETWORK, nil, 0, err
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			var total int32 = 0
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no":      reply.OrderNo,
					"amount":        reply.Amount,
					"currency_type": reply.CurrencyType,
					"business_name": reply.BusinessName,
					"subject":       reply.Subject,
				}
				if reply.OrderNo != "" {
					total = 1
				}
			}
			return reply.ResultCode, respData, total, nil
		}, "payment_code")
	}
}

//用户付款码支付
func (*BusinessBillHandler) OrderPay() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type") // 账号类型(用户才可以支付)
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, nil
			}

			req := &businessBillProto.OrderPayRequest{
				AccountNo:     inner_util.GetJwtDataString(c, "account_uid"), // 支付账号
				AccountType:   accountType,
				OrderNo:       inner_util.M(c, "order_no"),                      // 订单号
				PaymentPwd:    inner_util.M(c, "payment_pwd"),                   // 支付账号的支付密码
				NonStr:        inner_util.M(c, "non_str"),                       // 和支付密码一起使用的随机字符串
				PaymentMethod: inner_util.GetJwtDataString(c, "payment_method"), //付款方式
				BankCardNo:    inner_util.GetJwtDataString(c, "card_no"),        //银行卡id
				Lang:          ss_net.GetCommonData(c).Lang,
				SignKey:       inner_util.M(c, "sign_key"),    // 指纹支付标识
				DeviceUuid:    inner_util.M(c, "device_uuid"), // 设备id
			}

			ss_log.Info("api-mobile.Pay请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.OrderPay(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no":        reply.OrderNo,
					"order_status":    reply.OrderStatus,
					"create_time":     reply.CreateTime,
					"pay_time":        reply.PayTime,
					"subject":         reply.Subject,
					"user_order_type": reply.UserOrderType,
					"business_name":   reply.BusinessName,
				}
			}

			return reply.ResultCode, respData, nil
		})
	}
}

//查询订单交易信息(金额, 订单标题, 商家名称)
func (*BusinessBillHandler) QueryOrderTransInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, 0, nil
			}

			req := &businessBillProto.QueryTransInfoRequest{
				QrCodeId: params[0],
				Lang:     ss_net.GetCommonData(c).Lang,
			}
			ss_log.Info("api-mobile.QueryOrderInfo请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.QueryTransInfo(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用business_bill_srv.QueryOrderInfo失败, err=%v", err)
				return ss_err.ERR_SYS_NETWORK, nil, 0, err
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"amount":        reply.Amount,
					"subject":       reply.Subject,
					"business_name": reply.BusinessName,
					"currency_type": reply.CurrencyType,
				}
			}
			return reply.ResultCode, respData, 1, nil
		}, "qr")
	}
}

// 扫码支付-用户主动扫商家带金额的二维码(临时可用二维码)
func (*BusinessBillHandler) QrCodeAmountPay() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type") // 账号类型(用户才可以支付)
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, nil
			}

			req := &businessBillProto.QrCodeAmountPayRequest{
				AccountNo:       inner_util.GetJwtDataString(c, "account_uid"), // 支付账号
				AccountType:     accountType,
				QrCodeId:        inner_util.M(c, "qr"),                            // 二维码id
				PaymentPassword: inner_util.M(c, "account_pay_password"),          // 支付账号的支付密码
				NonStr:          inner_util.M(c, "non_str"),                       // 和支付密码一起使用的随机字符串
				PaymentMethod:   inner_util.GetJwtDataString(c, "payment_method"), //付款方式
				BankCardNo:      inner_util.GetJwtDataString(c, "card_no"),        //银行卡id
				Lang:            ss_net.GetCommonData(c).Lang,
				SignKey:         inner_util.M(c, "sign_key"),    // 指纹支付标识
				DeviceUuid:      inner_util.M(c, "device_uuid"), // 设备id
			}

			ss_log.Info("api-mobile.Pay请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.QrCodeAmountPay(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no":        reply.OrderNo,
					"order_status":    reply.OrderStatus,
					"create_time":     reply.CreateTime,
					"pay_time":        reply.PayTime,
					"subject":         reply.Subject,
					"user_order_type": reply.UserOrderType,
				}
			}

			return reply.ResultCode, respData, nil
		})
	}
}

// 扫码支付时查询商家的应用信息
func (*BusinessBillHandler) GetBusinessAppInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, 0, nil
			}

			req := &businessBillProto.GetBusinessAppInfoRequest{
				FixedQrcode: params[0],
				Lang:        ss_net.GetCommonData(c).Lang,
			}
			ss_log.Info("api-mobile.GetBusinessAppInfo请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.GetBusinessAppInfo(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用business_bill_srv.GetBusinessAppInfo失败, err=%v", err)
				return ss_err.ERR_SYS_NETWORK, nil, 0, err
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"app_name": reply.AppName,
				}
			}
			return reply.ResultCode, respData, 1, nil

		}, "qr")
	}
}

//用户APP下单
func (*BusinessBillHandler) QrCodeFixedPrePay() gin.HandlerFunc {
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
			req := &businessBillProto.QrCodeFixedPrePayRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"), // 支付账号
				AccountType:  accountType,
				QrCodeId:     inner_util.M(c, "qr"),            // 二维码id
				Subject:      subject,                          // 商品名称
				Amount:       inner_util.M(c, "amount"),        // 金额
				CurrencyType: inner_util.M(c, "currency_type"), // 币种
				Remark:       inner_util.M(c, "remark"),        // 备注
				Lang:         ss_net.GetCommonData(c).Lang,
			}

			ss_log.Info("api-mobile.QrCodeFixedPay请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.QrCodeFixedPrePay(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no": reply.OrderNo,
				}
			}

			return reply.ResultCode, respData, nil
		})
	}
}

// 扫码支付-用户主动扫商家不带金额的二维码(永久可用的二维码)
func (*BusinessBillHandler) QrCodeFixedPay() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type") // 账号类型(用户才可以支付)
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, nil
			}

			req := &businessBillProto.QrCodeFixedPayRequest{
				AccountNo:       inner_util.GetJwtDataString(c, "account_uid"), // 支付账号
				AccountType:     accountType,
				OrderNo:         inner_util.M(c, "order_no"),                      // 订单号
				PaymentPassword: inner_util.M(c, "account_pay_password"),          // 支付账号的支付密码
				NonStr:          inner_util.M(c, "non_str"),                       // 和支付密码一起使用的随机字符串
				PaymentMethod:   inner_util.GetJwtDataString(c, "payment_method"), //付款方式
				BankCardNo:      inner_util.GetJwtDataString(c, "card_no"),        //银行卡id
				Lang:            ss_net.GetCommonData(c).Lang,
				SignKey:         inner_util.M(c, "sign_key"),    // 指纹支付标识
				DeviceUuid:      inner_util.M(c, "device_uuid"), // 设备id
			}

			ss_log.Info("api-mobile.QrCodeFixedPay请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.QrCodeFixedPay(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no":        reply.OrderNo,
					"order_status":    reply.OrderStatus,
					"create_time":     reply.CreateTime,
					"pay_time":        reply.PayTime,
					"subject":         reply.Subject,
					"user_order_type": reply.UserOrderType,
				}
			}

			return reply.ResultCode, respData, nil
		})
	}
}

// APP支付
func (*BusinessBillHandler) AppPay() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type") // 账号类型(用户才可以支付)
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, nil
			}

			req := &businessBillProto.AppPayRequest{
				AccountNo:     inner_util.GetJwtDataString(c, "account_uid"), // 支付账号
				AccountType:   accountType,
				AppPayContent: inner_util.M(c, "app_pay_content"),               // 支付信息(平台发出，平台回收)
				PaymentPwd:    inner_util.M(c, "account_pay_password"),          // 支付账号的支付密码
				NonStr:        inner_util.M(c, "non_str"),                       // 和支付密码一起使用的随机字符串
				PaymentMethod: inner_util.GetJwtDataString(c, "payment_method"), //付款方式
				BankCardNo:    inner_util.GetJwtDataString(c, "card_no"),        //银行卡id
				Lang:          ss_net.GetCommonData(c).Lang,
				SignKey:       inner_util.M(c, "sign_key"),    // 指纹支付标识
				DeviceUuid:    inner_util.M(c, "device_uuid"), // 设备id
			}

			reply, err := BusinessBillHandlerInst.Client.AppPay(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no":        reply.OrderNo,
					"order_status":    reply.OrderStatus,
					"create_time":     reply.CreateTime,
					"pay_time":        reply.PayTime,
					"subject":         reply.Subject,
					"user_order_type": reply.UserOrderType,
				}
			}

			return reply.ResultCode, respData, nil
		})
	}
}

// 扫码支付时查询个人商家的信息
func (*BusinessBillHandler) GetPersonalBusinessInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			if accountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", accountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, 0, nil
			}

			req := &businessBillProto.GetPersonalBusinessInfoRequest{
				FixedCode: params[0],
				Lang:      ss_net.GetCommonData(c).Lang,
			}

			reply, err := BusinessBillHandlerInst.Client.GetPersonalBusinessInfo(context.TODO(), req)
			if err != nil {
				ss_log.Error("调用business_bill_srv.GetPersonalBusinessInfo失败, err=%v", err)
				return ss_err.ERR_SYS_NETWORK, nil, 0, err
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"business_name": reply.BusinessName,
				}
			}
			return reply.ResultCode, respData, 1, nil

		}, "qr")
	}
}

//用户扫个人商家固码下单
func (*BusinessBillHandler) PersonalFixedCodePrePay() gin.HandlerFunc {
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
			req := &businessBillProto.QrCodeFixedPrePayRequest{
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"), // 支付账号
				AccountType:  accountType,
				QrCodeId:     inner_util.M(c, "qr"),            // 二维码id
				Subject:      subject,                          // 商品名称
				Amount:       inner_util.M(c, "amount"),        // 金额
				CurrencyType: inner_util.M(c, "currency_type"), // 币种
				Remark:       inner_util.M(c, "remark"),        // 备注
				Lang:         ss_net.GetCommonData(c).Lang,
			}

			ss_log.Info("api-mobile.QrCodeFixedPay请求参数=[%v]", strext.ToJson(req))

			reply, err := BusinessBillHandlerInst.Client.QrCodeFixedPrePay(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			reply.ResultCode = ss_err.PayRetCode(reply.ResultCode)
			respData := gin.H{}
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				respData = gin.H{
					"order_no": reply.OrderNo,
				}
			}

			return reply.ResultCode, respData, nil
		})
	}
}
