package handler

import (
	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-mobile/common"
	"a.a/mp-server/api-mobile/dao"
	"a.a/mp-server/api-mobile/inner_util"
	"a.a/mp-server/api-mobile/verify"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/proto/auth"
	"a.a/mp-server/common/proto/bill"
	"a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_net"
	"context"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Client go_micro_srv_auth.AuthService
}

var (
	AuthHandlerInst AuthHandler
)

type BillHandler struct {
	Client go_micro_srv_bill.BillService
}

var BillHandlerInst BillHandler

type CustHandler struct {
	Client go_micro_srv_cust.CustService
}

var CustHandlerInst CustHandler

func (s *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNo := inner_util.GetJwtDataString(c, "account_uid")
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			dao.AccDaoInstance.DeleteLoginToken(accountNo)
			return ss_err.ERR_SUCCESS, "", nil
		})
	}
}

/**
 * 移动端登陆处理
 */
func (s *AuthHandler) MobileLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			countryCode := container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint()
			account := container.GetValFromMapMaybe(params, "account").ToStringNoPoint()

			account = ss_func.PrePhone(countryCode, account)

			req := &go_micro_srv_auth.MobileLoginRequest{
				Account:     account,
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				Imei:        container.GetValFromMapMaybe(params, "imei").ToStringNoPoint(),
				Ip:          c.ClientIP(),
				Nonstr:      container.GetValFromMapMaybe(params, "nonstr").ToStringNoPoint(),
				Lang:        ss_net.GetCommonData(c).Lang,
				AccountType: container.GetValFromMapMaybe(params, "account_type").ToStringNoPoint(), // 3-服务商;4-用户;5-收银员
				PosSn:       container.GetValFromMapMaybe(params, "pos_sn").ToStringNoPoint(),
				AppVersion:  container.GetValFromMapMaybe(params, "app_version").ToStringNoPoint(),
				PubKey:      container.GetValFromMapMaybe(params, "pub_key").ToStringNoPoint(), // 前端的pub_key
				Lat:         container.GetValFromMapMaybe(params, "lat").ToStringNoPoint(),     // 维度
				Lng:         container.GetValFromMapMaybe(params, "lng").ToStringNoPoint(),     // 经度
				CountryCode: countryCode,
				Uuid:        container.GetValFromMapMaybe(params, "uuid").ToStringNoPoint(),
			}

			//c.Set(common.Pub_Key, req.PubKey)
			// 参数校验
			if errStr := verify.MobileLoginReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.MobileLogin(context.TODO(), req)
			ss_log.Info("login=[%v],err=[%v]", reply, err)
			if reply.AccountUid != "" {
				replyAccount, err := s.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
					AccountUid:  reply.AccountUid,
					AccountType: reply.AccountType,
				})

				// 强制登录
				//isForce := container.GetValFromMapMaybe(params, "is_force").ToInt32()
				//routesStr, _ := json.Marshal(routes)
				loginToken := strext.NewUUIDNoSplit()
				//errCode := dao.AccDaoInstance.InsertLoginToken(reply.AccountUid, strext.ToStringNoPoint(routesStr), loginToken, isForce, c.ClientIP())
				errCode := dao.AccDaoInstance.InsertLoginToken(reply.AccountUid, "", loginToken, 1, c.ClientIP())
				if errCode != ss_err.ERR_SUCCESS {
					return errCode, gin.H{}, ss_err.ErrLoginFailed
				}

				// 登录成功
				//ss_net.SetJwtAuthentication(c, reply.Jwt)

				return strext.ToString(reply.ResultCode), gin.H{
					"userinfo": gin.H{
						"username":           replyAccount.Nickname,
						"token":              reply.Jwt,
						"access_status":      2,
						"account_uid":        replyAccount.Uid,
						"merchant_uid":       replyAccount.MerchantUid,
						"account_type":       reply.AccountType,
						"gen_key":            reply.GenKey,
						"phone":              reply.Phone,
						"is_set_payment_pwd": reply.IsSetPaymentPwd,
						//"head_portrait_img_no": reply.HeadPortraitImgNo,
						"head_portrait_img_no": replyAccount.HeadPortraitImgNo,
						"app_version":          reply.AppVersion,
						"version_description":  reply.VersionDescription,
						"app_url":              reply.AppUrl,
						//"is_first_login":       reply.IsFirstLogin, // 是否首次登陆,0-首次,1-不是首次
						"is_first_login":         "1", // 由于账号体系做了优化，这一步没有存在的意义了
						"login_token":            loginToken,
						"is_force":               reply.IsForce,              // 是否强制更新;0-否,1-是
						"refresh_token_interval": reply.RefreshTokenInterval, // 客户端刷新token的时间间隔
						"country_code":           countryCode,                // 国家码
					},
				}, err
			} else {
				// 登录失败
				// ss_net.SetJwtAuthentication(c, reply.Jwt)
				return reply.ResultCode, gin.H{}, err
			}
		}, "params")
	}
}

// 注册并且登录
func (s *AuthHandler) RegLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			countryCode := container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint()
			phone := container.GetValFromMapMaybe(params, "phone").ToStringNoPoint()

			phone = ss_func.PrePhone(countryCode, phone)

			req := &go_micro_srv_auth.RegLoginRequest{
				Phone:       phone,
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				Sms:         container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),
				Lang:        ss_net.GetCommonData(c).Lang,
				Ip:          c.ClientIP(),
				AppVersion:  container.GetValFromMapMaybe(params, "app_version").ToStringNoPoint(),
				PubKey:      container.GetValFromMapMaybe(params, "pub_key").ToStringNoPoint(),
				Lat:         container.GetValFromMapMaybe(params, "lat").ToStringNoPoint(),
				Lng:         container.GetValFromMapMaybe(params, "lng").ToStringNoPoint(),
				CountryCode: countryCode,
				Uuid:        container.GetValFromMapMaybe(params, "uuid").ToStringNoPoint(),
				UtmSource:   ss_net.GetCommonData(c).UtmSource,
			}
			// 参数校验
			if errStr := verify.RegLoginReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.RegLogin(context.TODO(), req)
			ss_log.Info("login=[%v]", reply)
			if reply.AccountUid != "" {
				replyAccount, err := s.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
					AccountUid:  reply.AccountUid,
					AccountType: reply.AccountType,
				})

				loginToken := strext.NewUUIDNoSplit()
				errCode := dao.AccDaoInstance.InsertLoginToken(reply.AccountUid, "", loginToken, 1, c.ClientIP())
				if errCode != ss_err.ERR_SUCCESS {
					return errCode, gin.H{}, ss_err.ErrLoginFailed
				}

				// 登录成功
				// ss_net.SetJwtAuthentication(c, reply.Jwt)
				return strext.ToString(reply.ResultCode), gin.H{
					"userinfo": gin.H{
						"username":               replyAccount.Nickname,
						"token":                  reply.Jwt,
						"access_status":          2,
						"account_uid":            replyAccount.Uid,
						"merchant_uid":           replyAccount.MerchantUid,
						"account_type":           reply.AccountType,
						"gen_key":                reply.GenKey,
						"phone":                  reply.Phone,
						"is_set_payment_pwd":     reply.IsSetPaymentPwd,
						"app_version":            reply.AppVersion,
						"version_description":    reply.VersionDescription,
						"app_url":                reply.AppUrl,
						"login_token":            loginToken,
						"is_force":               reply.IsForce,              // 是否强制更新;0-否,1-是
						"refresh_token_interval": reply.RefreshTokenInterval, // 客户端刷新token的时间间隔
						"country_code":           countryCode,                // 国家码
					},
				}, err
			} else {
				// 登录失败
				// ss_net.SetJwtAuthentication(c, reply.Jwt)
				return reply.ResultCode, gin.H{}, err
			}
		}, "params")
	}
}

func (s *AuthHandler) RegSms() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			phone := container.GetValFromMapMaybe(params, "phone").ToStringNoPoint()
			countryCode := container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint()

			phone = ss_func.PrePhone(countryCode, phone)

			if phone == "" {
				phone = inner_util.GetJwtDataString(c, "phone")
			}
			req := &go_micro_srv_cust.RegSmsRequest{
				//Phone:    container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				Phone:       phone,
				Lang:        ss_net.GetCommonData(c).Lang,
				Function:    container.GetValFromMapMaybe(params, "function").ToStringNoPoint(), //功能,注册-reg;找回密码-backpwd;修改支付密码-paypwd;验证手机号-checkphone;更改手机号-modifyphone;
				PubKey:      container.GetValFromMapMaybe(params, "pub_key").ToStringNoPoint(),
				CountryCode: countryCode,
			}

			//参数校验
			if errStr := verify.SmsReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, _ := CustHandlerInst.Client.RegSms(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

func (s *AuthHandler) CheckSms() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_cust.CheckSmsRequest{
				Phone:       container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				Sms:         container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),
				Function:    container.GetValFromMapMaybe(params, "function").ToStringNoPoint(),
				PubKey:      container.GetValFromMapMaybe(params, "pub_key").ToStringNoPoint(),
				CountryCode: container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(),
			}

			c.Set(common.Pub_Key, req.PubKey)
			// 参数校验
			if errStr := verify.CheckSmsReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, _ := CustHandlerInst.Client.CheckSms(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

func (s *AuthHandler) BackPWD() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.MobileBackPwdRequest{
				Phone:       container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				Sms:         container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),
				PubKey:      container.GetValFromMapMaybe(params, "Pub_Key").ToStringNoPoint(),
				CountryCode: container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(),
			}

			req.Phone = ss_func.PrePhone(req.CountryCode, req.Phone)

			//c.Set(common.Pub_Key, req.PubKey)
			//参数校验
			if errStr := verify.BackPWDReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, _ := s.Client.MobileBackPwd(context.TODO(), req)
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

// 修改密码
func (s *AuthHandler) ModifyPWD() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			uid := inner_util.GetJwtDataString(c, "account_uid")
			req := &go_micro_srv_auth.MobileModifyPwdRequest{
				Uid:         uid,
				OldPassword: container.GetValFromMapMaybe(params, "old_password").ToStringNoPoint(),
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),
				NewPassword: container.GetValFromMapMaybe(params, "new_password").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			}
			// 参数校验
			if errStr := verify.ModifyPWDReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := s.Client.MobileModifyPwd(context.TODO(), req)

			if reply.ResultCode == ss_err.ERR_SUCCESS {
				//此处是将修改密码的账户登出（踢出登录）
				ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
					dao.AccDaoInstance.DeleteLoginToken(uid)
					return ss_err.ERR_SUCCESS, "", nil
				})
			}

			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

// 修改支付密码
func (s *AuthHandler) ModifyPayPWD() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.MobileModifyPayPwdRequest{
				Uid:         inner_util.GetJwtDataString(c, "account_uid"),
				Sms:         container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
			}

			// 参数校验
			if errStr := verify.ModifyPayPWDReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := s.Client.MobileModifyPayPwd(context.TODO(), req)

			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

// 修改手机号
func (s *AuthHandler) ModifyPhone() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			phone := container.GetValFromMapMaybe(params, "phone").ToStringNoPoint()
			accNo := inner_util.GetJwtDataString(c, "account_uid")
			accType := inner_util.GetJwtDataString(c, "account_type")
			countryCode := container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint()

			phone = ss_func.PrePhone(countryCode, phone)

			req := &go_micro_srv_auth.MobileModifyPhonedRequest{
				Uid:         accNo,
				Sms:         container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),
				Phone:       phone,
				CountryCode: countryCode,
			}

			// 参数校验
			if errStr := verify.ModifyPhoneReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, _ := s.Client.MobileModifyPhone(context.TODO(), req)

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return reply.ResultCode, nil, nil
			}

			// 修改手机号码后会重新登录，此处不需要进行设置, xiaoyanchun 2020-05-22
			//m, _ := c.Get("decodedJwt")
			//m.(jwt.MapClaims)["account"] = phone
			//c.Set("decodedJwt", m)
			//c.Set("repackJwt", true)
			ss_log.Info("reply=[%v]", reply)

			replyAccount, err := s.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
				AccountUid:  accNo,
				AccountType: accType,
			})

			return strext.ToString(reply.ResultCode), gin.H{
				"userinfo": gin.H{
					"username": replyAccount.Nickname,
					//"token":         reply.Jwt,
					"access_status": 2,
					"account_uid":   replyAccount.Uid,
					"merchant_uid":  replyAccount.MerchantUid,
					"account_type":  accType,
					"gen_key":       replyAccount.GenKey,
					"phone":         phone,
					"country_code":  countryCode,
				},
			}, err

		})
	}
}

// 修改昵称
func (s *AuthHandler) ModifyNickname() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accNo := inner_util.GetJwtDataString(c, "account_uid")
			accType := inner_util.GetJwtDataString(c, "account_type")
			name := container.GetValFromMapMaybe(params, "nickname").ToStringNoPoint()

			req := &go_micro_srv_auth.ModifyNicknameRequest{
				AccountNo:   accNo,
				AccountType: accType,
				Nickname:    name,
			}
			// 参数校验
			if errStr := verify.ModifyNicknameReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.ModifyNickname(context.TODO(), req)

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return strext.ToString(reply.ResultCode), nil, err
			}

			//m, _ := c.Get("decodedJwt")
			//m.(jwt.MapClaims)["account_name"] = name
			//c.Set("decodedJwt", m)
			//c.Set("repackJwt", true)
			return ss_err.ERR_SUCCESS, nil, nil
		})
	}
}

func (s *AuthHandler) Userinfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountUid := inner_util.GetJwtDataString(c, "account_uid")
			accountType := inner_util.GetJwtDataString(c, "account_type")

			replyAccount, err := s.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
				Account:     inner_util.GetJwtDataString(c, "account"),
				AccountUid:  accountUid,
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				AccountType: accountType,
				//AccountName:    inner_util.GetJwtDataString(c, "account_name"),
				LoginAccountNo: inner_util.GetJwtDataString(c, "login_account_no"),
				PubKey:         inner_util.GetJwtDataString(c, "pub_key"),
				JumpIdenNo:     inner_util.GetJwtDataString(c, "jump_iden_no"),
				JumpIdenType:   inner_util.GetJwtDataString(c, "jump_iden_type"),
				MasterAccNo:    inner_util.GetJwtDataString(c, "master_acc_no"),
				IsMasterAcc:    inner_util.GetJwtDataString(c, "is_master_acc"),
				Lang:           strext.ToStringNoPoint(params[0]),
			})
			xLoginToken := c.Request.Header.Get("LoginToken")
			return strext.ToString(replyAccount.ResultCode), gin.H{
				"userinfo": gin.H{
					"head_portrait_img_no": replyAccount.HeadPortraitImgNo,
					"username":             replyAccount.Nickname,
					"login_token":          xLoginToken,
					"access_status":        2,
					"account_uid":          accountUid,
					"merchant_uid":         replyAccount.MerchantUid,
					"account_type":         accountType,
					"gen_key":              replyAccount.GenKey,
					"phone":                replyAccount.Phone,
					"country_code":         replyAccount.CountryCode,
				},
			}, 0, err
		}, "lang")
	}
}

func (s *AuthHandler) VersionInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_auth.GetVersinInfoRequest{
				AppVersion: container.GetValFromMapMaybe(params, "app_version").ToStringNoPoint(),
				VsType:     container.GetValFromMapMaybe(params, "vs_type").ToStringNoPoint(),
				System:     container.GetValFromMapMaybe(params, "system").ToStringNoPoint(),
			}
			if errStr := verify.VersionInfoReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.GetVersinInfo(context.TODO(), req)

			return strext.ToString(reply.ResultCode), gin.H{
				"versionInfo": gin.H{
					"app_version":         reply.AppVersion,
					"version_description": reply.VersionDescription,
					"app_url":             reply.AppUrl,
					"is_force":            reply.IsForce, // 是否强制更新;0-否,1-是
				},
			}, err
		}, "params")
	}
}

// 上传客户端信息
func (s *AuthHandler) UploadClientInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			// 公共参数
			cData := ss_net.GetCommonData(c)

			reply, err := s.Client.UploadClientInfo(context.TODO(), &go_micro_srv_auth.UploadClientInfoRequest{
				DeviceBrand: container.GetValFromMapMaybe(params, "device_brand").ToStringNoPoint(),
				DeviceModel: container.GetValFromMapMaybe(params, "device_model").ToStringNoPoint(),
				Resolution:  container.GetValFromMapMaybe(params, "resolution").ToStringNoPoint(),
				ScreenSize:  container.GetValFromMapMaybe(params, "screen_size").ToStringNoPoint(),
				Imei1:       container.GetValFromMapMaybe(params, "imei1").ToStringNoPoint(),
				Imei2:       container.GetValFromMapMaybe(params, "imei2").ToStringNoPoint(),
				SystemVer:   container.GetValFromMapMaybe(params, "system_ver").ToStringNoPoint(),
				UserAgent:   c.GetHeader("User-Agent"),
				Platform:    cData.Platform,
				AppVer:      cData.AppVersion,
				UploadPoint: container.GetValFromMapMaybe(params, "upload_point").ToStringNoPoint(),
				ClientType:  container.GetValFromMapMaybe(params, "client_type").ToInt32(),
				Account:     container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				Uuid:        container.GetValFromMapMaybe(params, "uuid").ToStringNoPoint(),
			})

			return reply.ResultCode, gin.H{}, err
		}, "params")
	}
}

func (s *AuthHandler) PerfectingInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accNo := inner_util.GetJwtDataString(c, "account_uid")
			name := container.GetValFromMapMaybe(params, "nickname").ToStringNoPoint()
			imageStr := container.GetValFromMapMaybe(params, "image_str").ToStringNoPoint()

			req := &go_micro_srv_auth.PerfectingInfoRequest{
				AccountNo: accNo,
				Nickname:  name,
				ImageStr:  imageStr,
			}

			if errStr := verify.PerfectingInfoReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.PerfectingInfo(context.TODO(), req, global.RequestTimeoutOptions)

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				return strext.ToString(reply.ResultCode), nil, err
			}

			//m, _ := c.Get("decodedJwt")
			//if name != "" {
			//	m.(jwt.MapClaims)["account_name"] = name
			//}

			//if imageStr != "" {
			//	m.(jwt.MapClaims)["head_portrait_img_no"] = reply.HeadPortraitImgNo
			//}
			//c.Set("decodedJwt", m)
			//c.Set("repackJwt", true)
			return ss_err.ERR_SUCCESS, nil, nil
		})
	}
}

func (s *AuthHandler) GetAppCommonHelps() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetAppCommonHelps(context.TODO(), &go_micro_srv_cust.GetAppCommonHelpsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				VsType:    strext.ToStringNoPoint(params[2]),
				Lang:      ss_net.GetCommonData(c).Lang,
				UseStatus: "1", //此为用户查看帮助接口，所以只能看到启用的
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "vs_type")

	}
}

func (s *AuthHandler) GetCommonHelpDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetCommonHelpDetail(context.TODO(), &go_micro_srv_cust.GetCommonHelpDetailRequest{
				HelpNo: strext.ToStringNoPoint(params[0]),
			})
			return reply.ResultCode, reply.Data, 0, err
		}, "help_no")

	}
}

//获取附近服务商列表
func (*AuthHandler) GetNearbyServicerList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := CustHandlerInst.Client.GetNearbyServicerList(context.TODO(), &go_micro_srv_cust.GetNearbyServicerListRequest{
				Page:       strext.ToInt32(params[0]),
				PageSize:   strext.ToInt32(params[1]),
				Lat:        strext.ToStringNoPoint(params[2]),
				Lng:        strext.ToStringNoPoint(params[3]),
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "lat", "lng")
	}
}

func (*AuthHandler) GetAgreementAppDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &go_micro_srv_cust.GetAgreementAppDetailRequest{
				Lang: ss_net.GetCommonData(c).Lang,
				Type: strext.ToStringNoPoint(params[1]),
			}

			reply, err := CustHandlerInst.Client.GetAgreementAppDetail(context.TODO(), req)
			return reply.ResultCode, reply.Data, 0, err
		}, "lang", "type")

	}
}

func (*AuthHandler) GetAppConsultationConfigDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			req := &go_micro_srv_cust.GetAppConsultationConfigDetailRequest{
				Lang: ss_net.GetCommonData(c).Lang,
			}

			reply, err := CustHandlerInst.Client.GetAppConsultationConfigDetail(context.TODO(), req)
			return reply.ResultCode, reply.Datas, 0, err
		}, "")

	}
}

// 刷新token
func (s *AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			req := &go_micro_srv_auth.RefreshTokenRequest{
				JwtIat:         inner_util.GetJwtDataString(c, "iat"),
				Account:        inner_util.GetJwtDataString(c, "account"),
				AccountUid:     inner_util.GetJwtDataString(c, "account_uid"),
				IdenNo:         inner_util.GetJwtDataString(c, "iden_no"),
				AccountType:    inner_util.GetJwtDataString(c, "account_type"),
				LoginAccountNo: inner_util.GetJwtDataString(c, "login_account_no"),
				PubKey:         inner_util.GetJwtDataString(c, "pub_key"),
				JumpIdenNo:     inner_util.GetJwtDataString(c, "jump_iden_no"),
				JumpIdenType:   inner_util.GetJwtDataString(c, "jump_iden_type"),
				MasterAccNo:    inner_util.GetJwtDataString(c, "master_acc_no"),
				IsMasterAcc:    inner_util.GetJwtDataString(c, "is_master_acc"),
				PosSn:          inner_util.GetJwtDataString(c, "pos_sn"),
			}

			reply, err := s.Client.RefreshToken(context.TODO(), req)
			if err != nil {
				ss_log.Error("请求RefreshToken接口失败, err:%v", err)
				return "", gin.H{}, err
			}

			if reply.ResultCode != ss_err.ERR_SUCCESS {
				ss_log.Error("请求RefreshToken接口失败, reply.ResultCode:%v", reply.ResultCode)
				return reply.ResultCode, gin.H{}, err
			}

			// ss_net.SetJwtAuthentication(c, reply.Jwt) // 设置jwt的头信息

			return reply.ResultCode, gin.H{
				"refresh_token_interval": reply.RefreshTokenInterval, // 客户端刷新token的时间间隔
				"token":                  reply.Jwt,
			}, nil
		}, "params")
	}
}

// 上传客户端信息
func (s *AuthHandler) InsertLogAppDot() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			reply, err := s.Client.InsertLogAppDot(context.TODO(), &go_micro_srv_auth.InsertLogAppDotRequest{
				OpType: container.GetValFromMapMaybe(params, "op_type").ToStringNoPoint(),
				Uuid:   container.GetValFromMapMaybe(params, "uuid").ToStringNoPoint(),
			})

			return reply.ResultCode, gin.H{}, err
		}, "params")
	}
}

func (s *AuthHandler) AddAuthMaterialInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {

			req := &go_micro_srv_cust.AddAuthMaterialInfoRequest{
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
				FrontImgNo: container.GetValFromMapMaybe(params, "front_img_no").ToStringNoPoint(),
				BackImgNo:  container.GetValFromMapMaybe(params, "back_img_no").ToStringNoPoint(),
				AuthName:   container.GetValFromMapMaybe(params, "auth_name").ToStringNoPoint(),
				AuthNumber: container.GetValFromMapMaybe(params, "auth_number").ToStringNoPoint(),
			}

			if errStr := verify.AddAuthMaterialInfoVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.AddAuthMaterialInfo(context.TODO(), req, global.RequestTimeoutOptions)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			return reply.ResultCode, nil, nil
		})
	}
}

// 查询实名认证信息
func (s *AuthHandler) GetAuthMaterialInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			reply, err := CustHandlerInst.Client.GetAuthMaterialInfo(context.TODO(), &go_micro_srv_cust.GetAuthMaterialInfoRequest{
				AccountUid: inner_util.GetJwtDataString(c, "account_uid"),
			})

			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return reply.ResultCode, gin.H{
				"data": reply.Data,
			}, err
		}, "params")
	}
}
