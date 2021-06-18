package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/inner_util"
	"a.a/mp-server/api-mobile/verify"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/proto/bill"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_net"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

/**
 * 货币兑换
 */
func (AuthHandler) Exchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			lat := container.GetValFromMapMaybe(params, "lat").ToStringNoPoint()
			lng := container.GetValFromMapMaybe(params, "lng").ToStringNoPoint()
			lat = fmt.Sprintf("%.8f", strext.ToFloat64(lat))
			lng = fmt.Sprintf("%.8f", strext.ToFloat64(lng))

			req := &go_micro_srv_bill.ExchangeRequest{
				// 转入类型
				InType: container.GetValFromMapMaybe(params, "in_type").ToStringNoPoint(),
				// 转出类型
				OutType: container.GetValFromMapMaybe(params, "out_type").ToStringNoPoint(),
				// 金额
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				// 账号
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				Lang:      ss_net.GetCommonData(c).Lang,

				TransFrom:   container.GetValFromMapMaybe(params, "trans_from").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				// 操作员账号
				OpAccNo:    inner_util.GetJwtDataString(c, "iden_no"),
				NonStr:     container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				Ip:         c.ClientIP(),
				Password:   container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				Lat:        lat,
				Lng:        lng,
				SignKey:    container.GetValFromMapMaybe(params, "sign_key").ToStringNoPoint(),
				DeviceUuid: container.GetValFromMapMaybe(params, "device_uuid").ToStringNoPoint(),
			}

			// 参数校验
			if errStr := verify.ExchangeReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.Exchange(context.TODO(), req)
			if err != nil {
				ss_log.Error("BillHandlerInst.Client.Exchange err: %s", err.Error())
				return ss_err.ERR_PARAM, nil, nil
			}
			ss_log.Info("reply=[%v]", reply)
			payPasswordErrTips := ""
			if reply.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
				payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, reply.ResultCode, reply.PayPasswordErrTips)
				reply.ResultCode = ss_err.ERR_DB_PWD
			}
			ss_log.Info("payPasswordErrTips-------------------", payPasswordErrTips)
			c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)

			return reply.ResultCode, gin.H{
				"order_no": reply.OrderNo,
				"risk_no":  reply.RiskNo,
			}, nil
		})
	}
}

/**
 * 转账
 */
func (AuthHandler) Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			lat := container.GetValFromMapMaybe(params, "lat").ToStringNoPoint()
			lng := container.GetValFromMapMaybe(params, "lng").ToStringNoPoint()
			lat = fmt.Sprintf("%.8f", strext.ToFloat64(lat))
			lng = fmt.Sprintf("%.8f", strext.ToFloat64(lng))
			toPhone := container.GetValFromMapMaybe(params, "to_phone").ToStringNoPoint()
			countryCode := container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint()

			toPhone = ss_func.PrePhone(countryCode, toPhone)

			req := &go_micro_srv_bill.TransferRequest{
				// 货币类型
				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				// 转出账号
				FromAccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:   inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:        inner_util.GetJwtDataString(c, "iden_no"),
				Lang:          ss_net.GetCommonData(c).Lang,
				// 转入手机号
				ToPhone: toPhone,
				// 金额
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				// 转账类型
				ExchangeType:  container.GetValFromMapMaybe(params, "exchange_type").ToStringNoPoint(),
				NonStr:        container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				Password:      container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				CountryCode:   countryCode,
				Ip:            c.ClientIP(),
				Lat:           lat,
				Lng:           lng,
				PaymentMethod: container.GetValFromMapMaybe(params, "payment_method").ToStringNoPoint(),
				CardNo:        container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				SignKey:       container.GetValFromMapMaybe(params, "sign_key").ToStringNoPoint(),
				DeviceUuid:    container.GetValFromMapMaybe(params, "device_uuid").ToStringNoPoint(),
			}

			if errStr := verify.TransferReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.Transfer(context.TODO(), req)
			if err != nil {
				ss_log.Error("BillHandlerInst.Client.Transfer err: %s", err.Error())
				return ss_err.ERR_PARAM, nil, nil
			}
			ss_log.Info("reply=[%v]", reply)
			payPasswordErrTips := ""
			if reply.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
				payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, reply.ResultCode, reply.PayPasswordErrTips)
				reply.ResultCode = ss_err.ERR_DB_PWD
			}
			c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)
			reply2, err := AuthHandlerInst.Client.AddAccountCollect(context.TODO(), &go_micro_srv_auth.AddAccountCollectRequest{
				// 转入手机号
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				ToPhone:     toPhone,
				CountryCode: countryCode,
			})
			if err != nil {
				ss_log.Error("err=[%v]", err)
			}

			ss_log.Info("reply2=[%v]", reply2)

			return reply.ResultCode, gin.H{
				"order_no": reply.OrderNo,
				"risk_no":  reply.RiskNo,
			}, nil
		})
	}
}

/**
 * 获取收款码
 */
func (AuthHandler) GenRecvCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_bill.GenRecvCodeRequest{
				// 货币类型
				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				// 金额
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				// 账号
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
			}

			// 参数校验
			if errStr := verify.GenRecvCodeReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := BillHandlerInst.Client.GenRecvCode(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, gin.H{
				"code": reply.Code,
			}, nil
		})
	}
}

/**
 * 收款码被扫，获取信息
 */
func (AuthHandler) ScanRecvCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle2(c, func(params []string) (string, interface{}, error) {
			reply, _ := BillHandlerInst.Client.ScanRecvCode(context.TODO(), &go_micro_srv_bill.ScanRecvCodeRequest{
				// 码
				Code: strext.ToStringNoPoint(params[0]),
			})
			ss_log.Info("reply=[%v]", reply)
			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, nil
			}

			return reply.ResultCode, gin.H{
				"amount":          reply.Data.Amount,
				"recv_phone":      reply.Data.RecvPhone,
				"fee_rate":        reply.Data.FeeRate,
				"recv_account_no": reply.Data.AccountNo,
				"money_type":      reply.Data.MoneyType,
			}, nil
		}, "code")
	}
}

// pos存款
//func (a *AuthHandler) SaveMoney() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
//			req := &go_micro_srv_bill.SaveMoneyRequest{
//				// 收款人手机号
//				RecvPhone: container.GetValFromMapMaybe(params, "recv_phone").ToStringNoPoint(),
//				// 存款人手机号
//				SendPhone: container.GetValFromMapMaybe(params, "send_phone").ToStringNoPoint(),
//				// 金额
//				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
//				// 支付密码
//				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
//				// 币种
//				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
//
//				// 操作
//				AccountType: inner_util.GetJwtDataString(c, "account_type"),
//				// 操作员账号
//				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
//				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
//				CountryCode: container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(),
//
//				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
//
//				Lang: ss_net.GetCommonData(c).Lang,
//				Ip:   c.ClientIP(),
//			}
//
//			// 参数校验
//			if errStr := verify.SaveMoneyReqVerify(req); errStr != "" {
//				return errStr, nil, nil
//			}
//
//			reply, err := BillHandlerInst.Client.SaveMoney(context.TODO(), req)
//			ss_log.Info("SaveMoney=[%v],err=[%v]", reply, err)
//			return reply.ResultCode, gin.H{
//				"order_no": reply.OrderNo,
//				"risk_no":  reply.RiskNo,
//			}, nil
//		})
//	}
//}

// 收款
func (a *AuthHandler) Collection() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate3(c, func(params interface{}) (string, interface{}, error) {
			lat := container.GetValFromMapMaybe(params, "lat").ToStringNoPoint()
			lng := container.GetValFromMapMaybe(params, "lng").ToStringNoPoint()
			lat = fmt.Sprintf("%.8f", strext.ToFloat64(lat))
			lng = fmt.Sprintf("%.8f", strext.ToFloat64(lng))
			req := &go_micro_srv_bill.CollectionRequest{
				// 扫码人的id
				SweepAccountUid: inner_util.GetJwtDataString(c, "account_uid"),
				// 语言
				Lang: ss_net.GetCommonData(c).Lang,
				// 收款人的id,从二维码中获得
				RecAccountUid: container.GetValFromMapMaybe(params, "rec_account_uid").ToStringNoPoint(),
				// 币种
				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				// 金额
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),

				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),

				NonStr: container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),

				AccountType: inner_util.GetJwtDataString(c, "account_type"),

				OpAccNo: inner_util.GetJwtDataString(c, "iden_no"),

				GenCode: container.GetValFromMapMaybe(params, "gen_code").ToStringNoPoint(),
				Lat:     lat,
				Lng:     lng,
				Ip:      c.ClientIP(),
			}

			// 参数校验
			if errStr := verify.CollectionReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := BillHandlerInst.Client.Collection(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, gin.H{
				"order_no": reply.OrderNo,
				"risk_no":  reply.RiskNo,
			}, nil
		})
	}
}

// 扫一扫取款码
func (AuthHandler) GenWithdrawCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_bill.GenWithdrawCodeRequest{
				AccountNo:   inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
			}

			reply, _ := BillHandlerInst.Client.GenWithdrawCode(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, gin.H{
				"code": reply.Code,
			}, nil
		})
	}
}

// 修改用户扫一扫取款吗状态
func (*AuthHandler) ModityGenCodeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_bill.ModifyGenCodeStatusRequest{
				GenKey:      container.GetValFromMapMaybe(params, "gen_key").ToStringNoPoint(),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				Status:      container.GetValFromMapMaybe(params, "status").ToInt32(),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"), // 用户(cust)的accountid
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			}

			if errStr := verify.ModityGenCodeStatusReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := BillHandlerInst.Client.ModifyGenCodeStatus(context.TODO(), req)

			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

// 获取扫一扫取款码状态
func (AuthHandler) QuerySweepCodeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (s string, hs gin.H, e error) {
			req := &go_micro_srv_bill.QuerySweepCodeStatusRequest{
				GenCode:   strext.ToStringNoPoint(params),
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				IdenNo:    inner_util.GetJwtDataString(c, "iden_no"),
			}

			if errStr := verify.QuerySweepCodeStatusReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := BillHandlerInst.Client.QuerySweepCodeStatus(context.TODO(), req)
			return reply.ResultCode, gin.H{
				"status":           reply.Status,
				"order_no":         reply.OrderNo,
				"sweep_account_no": reply.SweepAccountNo,
				"nick_name":        reply.NickName,
				"head_url":         reply.HeadUrl,
			}, nil

		}, "gen_code")
	}
}

// 扫一扫取款
func (a *AuthHandler) SweepWithdrawal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			lat := container.GetValFromMapMaybe(params, "lat").ToStringNoPoint()
			lng := container.GetValFromMapMaybe(params, "lng").ToStringNoPoint()
			lat = fmt.Sprintf("%.8f", strext.ToFloat64(lat))
			lng = fmt.Sprintf("%.8f", strext.ToFloat64(lng))
			req := &go_micro_srv_bill.SweepWithdrawRequest{
				// 金额
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				// 支付密码
				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				// 币种
				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				// 账号id
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),

				NonStr: container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),

				AccountType:   inner_util.GetJwtDataString(c, "account_type"),
				OpAccNo:       inner_util.GetJwtDataString(c, "iden_no"),
				Lang:          ss_net.GetCommonData(c).Lang,
				GenCode:       container.GetValFromMapMaybe(params, "gen_code").ToStringNoPoint(),
				SwithdrawType: container.GetValFromMapMaybe(params, "swithdraw_type").ToInt32(),
				Ip:            c.ClientIP(),
				Lat:           lat,
				Lng:           lng,
			}

			if errStr := verify.SweepWithdrawalReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
			defer cancel()

			reply, err := BillHandlerInst.Client.SweepWithdrawal(ctx, req)
			ss_log.Info("reply=[%v],err=[%v]", reply, err)
			payPasswordErrTips := ""
			if reply.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
				payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, reply.ResultCode, reply.PayPasswordErrTips)
				reply.ResultCode = ss_err.ERR_DB_PWD
			}
			c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)

			return reply.ResultCode, gin.H{
				"order_no": reply.OrderNo,
			}, nil
		})
	}
}

// pos端确认提现操作
func (a *AuthHandler) ConfirmpWithdrawal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_bill.ConfirmWithdrawRequest{
				// 金额
				Amount: container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				// 支付密码
				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				// 币种
				MoneyType: container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				// 用户账号id
				UseAccountUid: container.GetValFromMapMaybe(params, "use_account_uid").ToStringNoPoint(),

				NonStr: container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),

				AccountType: inner_util.GetJwtDataString(c, "account_type"),

				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),

				OpAccNo: inner_util.GetJwtDataString(c, "iden_no"),

				Lang: ss_net.GetCommonData(c).Lang,

				OutOrderNo: container.GetValFromMapMaybe(params, "out_order_no").ToStringNoPoint(),
				GenCode:    container.GetValFromMapMaybe(params, "gen_code").ToStringNoPoint(),
			}

			// 参数校验
			if errStr := verify.ConfirmpWithdrawalReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := BillHandlerInst.Client.ConfirmWithdrawal(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

// pos 取消提现操作
func (a *AuthHandler) CancelWithdrawal() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_bill.CancelWithdrawRequest{
				// 订单号
				OrderNo:      container.GetValFromMapMaybe(params, "order_no").ToStringNoPoint(),
				CancelReason: container.GetValFromMapMaybe(params, "cancel_reason").ToStringNoPoint(),
			}
			// 参数校验
			if errStr := verify.CancelWithdrawalReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := BillHandlerInst.Client.CancelWithdraw(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

// 上传图片
func (a *AuthHandler) UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			imageStr := container.GetValFromMapMaybe(params, "image_str").ToStringNoPoint()
			if len(imageStr) > constants.UploadImgBase64LengthMax {
				return ss_err.ERR_ACCOUNT_IMAGE_BIG, nil, nil
			}

			addWatermark := ""
			imgType := container.GetValFromMapMaybe(params, "img_type").ToStringNoPoint() //1上传头像图片，2身份认证图片，3凭证图片

			switch imgType { //当前是身份认证的图片和上传凭证图片是要加水印的
			case "1":
			case "2":
				addWatermark = constants.AddWatermark_True
			case "3":
				addWatermark = constants.AddWatermark_True
			default:
			}

			req := &go_micro_srv_cust.UploadImageRequest{
				ImageStr:     imageStr,
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
				Type:         container.GetValFromMapMaybe(params, "type").ToInt32(), //1上传为需要授权才可访问的图片，2不需要授权
				AddWatermark: addWatermark,                                           //1为图片添加水印，为空则不添加
			}

			if errStr := verify.UploadImageReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.UploadImage(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("upload err:--------> %s", err.Error())
			}
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, gin.H{
				"image_id":   reply.ImageId,
				"image_name": reply.ImageName,
			}, nil
		})
	}
}

// 下载图片
/*
func (a *AuthHandler) DownloadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_cust.DownloadImageRequest{
				ImageId: container.GetValFromMapMaybe(params, "image_id").ToStringNoPoint(),
			}

			if errStr := verify.DownloadImageReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, _ := CustHandlerInst.Client.DownloadImage(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, gin.H{
				"image_str": reply.ImageStr,
			}, nil
		})
	}
}
*/

func (a *AuthHandler) UnAuthDownloadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_cust.UnAuthDownloadImageRequest{
				ImageId: container.GetValFromMapMaybe(params, "image_id").ToStringNoPoint(),
				PubKey:  container.GetValFromMapMaybe(params, "pub_key").ToStringNoPoint(),
			}
			//参数校验
			if errStr := verify.UnAuthDownloadImageReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := CustHandlerInst.Client.UnAuthDownloadImage(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, gin.H{
				"image_str": reply.ImageUrl,
			}, nil
		})
	}
}

// 转账到总部
func (AuthHandler) TransferToHeadquarters() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_bill.TransferToHeadquartersRequest{
				ImageId:     container.GetValFromMapMaybe(params, "image_id").ToStringNoPoint(),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				Amount:      container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				RecName:     container.GetValFromMapMaybe(params, "rec_name").ToStringNoPoint(),
				RecCarNum:   container.GetValFromMapMaybe(params, "rec_car_num").ToStringNoPoint(),
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
			}
			// 参数校验
			if errStr := verify.TransferToHeadquartersReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, _ := BillHandlerInst.Client.TransferToHeadquarters(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

/**
// 测试参数
{
    "image_id":"ssss.com",
    "account_uid":"e9586425-bfb7-4054-88b2-f1dfa47bdfa3",
    "money_type":"usd",
    "amount":"10000",
    "rec_name":"财务部小米",
    "rec_car_num":"442255331",
    "password":"adf",
    "account_type":"4",
    "non_str":"adfa",
    "op_acc_no":"664aa230-0665-475d-8f5c-d8a94f22923a",
    "card_no":"bc107ca7-cf81-44ea-94c9-39bd37243b9b",
    "ip":"127.2.0"
}
*/
func (AuthHandler) CustTransferToHeadquarters() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			lat := container.GetValFromMapMaybe(params, "lat").ToStringNoPoint()
			lng := container.GetValFromMapMaybe(params, "lng").ToStringNoPoint()
			lat = fmt.Sprintf("%.8f", strext.ToFloat64(lat))
			lng = fmt.Sprintf("%.8f", strext.ToFloat64(lng))
			req := &go_micro_srv_bill.CustTransferToHeadquartersRequest{
				ImageId:     container.GetValFromMapMaybe(params, "image_id").ToStringNoPoint(),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				Amount:      container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				RecName:     container.GetValFromMapMaybe(params, "rec_name").ToStringNoPoint(),
				RecCarNum:   container.GetValFromMapMaybe(params, "rec_car_num").ToStringNoPoint(),
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
				CardNo:      container.GetValFromMapMaybe(params, "card_no").ToStringNoPoint(),
				Lat:         lat,
				Lng:         lng,
				Ip:          c.ClientIP(),
				SignKey:     container.GetValFromMapMaybe(params, "sign_key").ToStringNoPoint(),
				DeviceUuid:  container.GetValFromMapMaybe(params, "device_uuid").ToStringNoPoint(),
			}
			// 参数校验
			if errStr := verify.CustTransferToHeadquartersReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := BillHandlerInst.Client.CustTransferToHeadquarters(context.TODO(), req)
			if err != nil {
				ss_log.Error("BillHandlerInst.Client.CustTransferToHeadquarters err: %s", err.Error())
				return ss_err.ERR_PARAM, nil, nil
			}
			ss_log.Info("reply=[%v]", reply)

			payPasswordErrTips := ""
			if reply.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
				payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, reply.ResultCode, reply.PayPasswordErrTips)
				reply.ResultCode = ss_err.ERR_DB_PWD
			}
			c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)

			return reply.ResultCode, gin.H{
				"order_no": reply.OrderNo,
			}, nil
		})
	}
}

// 向总部请款
func (*AuthHandler) ApplyMoney() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_bill.ApplyMoneyRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				MoneyType:   container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				Amount:      container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				RecCarNum:   container.GetValFromMapMaybe(params, "rec_car_num").ToStringNoPoint(),
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				OpAccNo:     inner_util.GetJwtDataString(c, "iden_no"),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				ChannelName: container.GetValFromMapMaybe(params, "channel_name").ToStringNoPoint(),
			}

			// 参数校验
			if errStr := verify.ApplyMoneyReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := BillHandlerInst.Client.ApplyMoney(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

// app端客户直接从总部提现
func (*AuthHandler) CustWithdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			lat := container.GetValFromMapMaybe(params, "lat").ToStringNoPoint()
			lng := container.GetValFromMapMaybe(params, "lng").ToStringNoPoint()
			lat = fmt.Sprintf("%.8f", strext.ToFloat64(lat))
			lng = fmt.Sprintf("%.8f", strext.ToFloat64(lng))
			req := &go_micro_srv_bill.CustWithdrawRequest{
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
				MoneyType:    container.GetValFromMapMaybe(params, "money_type").ToStringNoPoint(),
				Amount:       container.GetValFromMapMaybe(params, "amount").ToStringNoPoint(),
				RecCarNo:     container.GetValFromMapMaybe(params, "rec_car_no").ToStringNoPoint(),
				Password:     container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				AccountType:  inner_util.GetJwtDataString(c, "account_type"),
				NonStr:       container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				WithdrawType: container.GetValFromMapMaybe(params, "withdraw_type").ToStringNoPoint(), //1-普通提现;2-全部提现
				IdenNo:       inner_util.GetJwtDataString(c, "iden_no"),
				Lat:          lat,
				Lng:          lng,
				Ip:           c.ClientIP(),
				Lang:         ss_net.GetCommonData(c).Lang,
				SignKey:      container.GetValFromMapMaybe(params, "sign_key").ToStringNoPoint(),
				DeviceUuid:   container.GetValFromMapMaybe(params, "device_uuid").ToStringNoPoint(),
			}

			// 参数校验
			if errStr := verify.CustWithdrawReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.CustWithdraw(context.TODO(), req)
			if err != nil {
				ss_log.Error("BillHandlerInst.Client.CustWithdraw err: %s", err.Error())
				return ss_err.ERR_PARAM, nil, nil
			}
			ss_log.Info("reply=[%v]", reply)
			payPasswordErrTips := ""
			if reply.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
				payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, reply.ResultCode, reply.PayPasswordErrTips)
				reply.ResultCode = ss_err.ERR_DB_PWD
			}
			c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)
			return reply.ResultCode, gin.H{
				"order_no":              reply.OrderNo,
				"pay_password_err_tips": payPasswordErrTips,
			}, nil
		})
	}
}

func (*AuthHandler) GetTransferToHeadquartersLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.GetTransferToHeadquartersLog(context.TODO(), &go_micro_srv_bill.GetTransferToHeadquartersLogRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToString(params[2]),
				EndTime:      strext.ToString(params[3]),
				OrderStatus:  strext.ToString(params[4]),
				CurrencyType: strext.ToString(params[5]),
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:  inner_util.GetJwtDataString(c, "account_type"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "order_status", "money_type")

	}
}

func (*AuthHandler) GetTransferToServicerLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.GetTransferToServicerLogs(context.TODO(), &go_micro_srv_bill.GetTransferToServicerLogsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				StartTime:    strext.ToStringNoPoint(params[2]),
				EndTime:      strext.ToStringNoPoint(params[3]),
				OrderStatus:  strext.ToStringNoPoint(params[4]),
				CurrencyType: strext.ToStringNoPoint(params[5]),
				AccountNo:    inner_util.GetJwtDataString(c, "account_uid"),
				AccountType:  inner_util.GetJwtDataString(c, "account_type"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "order_status", "money_type")

	}
}

func (s *AuthHandler) CustIncomeBillsDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.CustIncomeBillsDetail(context.TODO(), &go_micro_srv_bill.CustIncomeBillsDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "log_no")
	}
}

func (s *AuthHandler) CustOutgoBillsDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.CustOutgoBillsDetail(context.TODO(), &go_micro_srv_bill.CustOutgoBillsDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "log_no")
	}
}

func (s *AuthHandler) CustTransferBillsDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.CustTransferBillsDetail(context.TODO(), &go_micro_srv_bill.CustTransferBillsDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "log_no")
	}
}

func (s *AuthHandler) CustCollectionBillsDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := BillHandlerInst.Client.CustCollectionBillsDetail(context.TODO(), &go_micro_srv_bill.CustCollectionBillsDetailRequest{
				LogNo: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "log_no")
	}
}

func (s *AuthHandler) GetAccountCollect() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AuthHandlerInst.Client.GetAccountCollect(context.TODO(), &go_micro_srv_auth.GetAccountCollectRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
			})

			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size")
	}
}

func (s *AuthHandler) CustOrderBillDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (s string, hs gin.H, e error) {
			req := &go_micro_srv_bill.CustOrderBillDetailRequest{
				AccountNo: inner_util.GetJwtDataString(c, "account_uid"),
				OrderNo:   container.GetValFromMapMaybe(params, "order_no").ToStringNoPoint(),
				OrderType: container.GetValFromMapMaybe(params, "order_type").ToStringNoPoint(),
				LogNo:     container.GetValFromMapMaybe(params, "log_no").ToStringNoPoint(),
			}

			// 参数校验
			if errStr := verify.CustOrderBillDetailReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := BillHandlerInst.Client.CustOrderBillDetail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, nil
		}, "params")
	}
}

//app客户获取银行卡提现账单
func (s *AuthHandler) GetLogToCusts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetLogToCusts(context.TODO(), &go_micro_srv_cust.GetLogToCustsRequest{
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToString(params[2]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "currency_type")
	}
}

func (*AuthHandler) GetLogCustToHeadquarters() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetLogCustToHeadquarters(context.TODO(), &go_micro_srv_cust.GetLogCustToHeadquartersRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "currency_type")
	}
}

func (*AuthHandler) CustOutgoBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.CustOutgoBills(context.TODO(), &go_micro_srv_cust.CustOutgoBillsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "currency_type")
	}
}

func (*AuthHandler) CustIncomeBills() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.CustIncomeBills(context.TODO(), &go_micro_srv_cust.CustIncomeBillsRequest{
				Page:         strext.ToInt32(params[0]),
				PageSize:     strext.ToInt32(params[1]),
				CurrencyType: strext.ToStringNoPoint(params[2]),
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "currency_type")
	}
}

func (*AuthHandler) AddFingerprint() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			//1.验证支付密码
			req := &go_micro_srv_auth.CheckPayPWDRequest{
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),

				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				NonStr:   container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
			}

			if req.AccountType != constants.AccountType_USER {
				ss_log.Error("登陆的账号角色为[%v],无权限调用此api", req.AccountType)
				return ss_err.ERR_SYS_NO_API_AUTH, nil, nil
			}

			if errStr := verify.CheckPayPWDReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := AuthHandlerInst.Client.CheckPayPWD(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				payPasswordErrTips := ""
				if reply.ResultCode == ss_err.ERR_PAY_FAILED_COUNT {
					payPasswordErrTips = ss_err.GetMsgAddArgs(ss_net.GetCommonData(c).Lang, reply.ResultCode, reply.ErrTips)
					reply.ResultCode = ss_err.ERR_DB_PWD
				}
				ss_log.Info("payPasswordErrTips-------------------", payPasswordErrTips)
				c.Set(ss_net.RET_CUSTOM_MSG, payPasswordErrTips)
				return reply.ResultCode, nil, nil
			}

			//2.添加指纹标识
			reqAdd := &go_micro_srv_auth.AddAppFingerprintRequest{
				AccountNo:  inner_util.GetJwtDataString(c, "account_uid"),
				DeviceUuid: container.GetValFromMapMaybe(params, "device_uuid").ToStringNoPoint(),
			}

			if errStr := verify.AddAppFingerprintVerify(reqAdd); errStr != "" {
				return errStr, nil, nil
			}

			replyAdd, errAdd := AuthHandlerInst.Client.AddAppFingerprint(context.TODO(), reqAdd)
			if errAdd != nil {
				ss_log.Error("api调用失败,err=[%v]", errAdd)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"sign_key": replyAdd.SignKey,
			}, nil
		})
	}
}

//关闭指纹支付
func (*AuthHandler) CloseFingerprint() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.CloseAppFingerprintRequest{
				AccountNo:  inner_util.GetJwtDataString(c, "account_uid"),
				DeviceUuid: container.GetValFromMapMaybe(params, "device_uuid").ToStringNoPoint(),
			}

			if errStr := verify.CloseAppFingerprintVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := AuthHandlerInst.Client.CloseAppFingerprint(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, nil, nil
		})
	}
}

func (*AuthHandler) GetMainIndustryCascaderDatas() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetMainIndustryCascaderDatas(context.TODO(), &go_micro_srv_cust.GetMainIndustryCascaderDatasRequest{
				Lang: ss_net.GetCommonData(c).Lang,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}

			return reply.ResultCode, reply.Datas, 0, nil
		}, "key")
	}
}

//申请个人商家
func (*AuthHandler) AddAuthMaterialBusiness() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_cust.AddAuthMaterialBusinessRequest{
				ImgId:        container.GetValFromMapMaybe(params, "img_id").ToString(),        //营业执照图片id
				AuthName:     container.GetValFromMapMaybe(params, "auth_name").ToString(),     //公司名称
				AuthNumber:   container.GetValFromMapMaybe(params, "auth_number").ToString(),   //注册号/组织机构代码
				StartDate:    container.GetValFromMapMaybe(params, "start_date").ToString(),    //营业期限起始时间
				EndDate:      container.GetValFromMapMaybe(params, "end_date").ToString(),      //营业期限结束时间
				TermType:     container.GetValFromMapMaybe(params, "term_type").ToString(),     //营业期限类型（1-短期，2长期）
				SimplifyName: container.GetValFromMapMaybe(params, "simplify_name").ToString(), //简称
				IndustryNo:   container.GetValFromMapMaybe(params, "industry_no").ToString(),   //经营类目id
				AccountUid:   inner_util.GetJwtDataString(c, "account_uid"),                    //登陆账号的uid
			}

			if errStr := verify.AddAuthMaterialBusinessVerify(req); errStr != "" {
				return errStr, 0, nil
			}

			req.AuthName = strings.Trim(req.AuthName, " ")
			req.AuthNumber = strings.Trim(req.AuthNumber, " ")
			req.SimplifyName = strings.Trim(req.SimplifyName, " ")

			var reply, err = CustHandlerInst.Client.AddAuthMaterialBusiness(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, nil
		})
	}
}

func (s *AuthHandler) GetAuthMaterialDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			accountUid := inner_util.GetJwtDataString(c, "account_uid")

			if accountUid == "" {
				ss_log.Error("AccountUid参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}

			reply, err := CustHandlerInst.Client.GetAuthMaterialBusinessDetail(context.TODO(), &go_micro_srv_cust.GetAuthMaterialBusinessDetailRequest{
				AccountUid: accountUid,
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err

		}, "params")
	}
}
