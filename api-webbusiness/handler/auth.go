package handler

import (
	"context"
	"encoding/json"
	"strings"

	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webbusiness/common"
	"a.a/mp-server/api-webbusiness/dao"
	"a.a/mp-server/api-webbusiness/handler/processor"
	"a.a/mp-server/api-webbusiness/inner_util"
	"a.a/mp-server/api-webbusiness/verify"
	"a.a/mp-server/common/constants"
	go_micro_srv_auth "a.a/mp-server/common/proto/auth"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/common/ss_net"
	"a.a/mp-server/common/ss_sql"
	"github.com/gin-gonic/gin"
	"github.com/wiwii/base64Captcha"
)

type AuthHandler struct {
	Client go_micro_srv_auth.AuthService
}

var (
	AuthHandlerInst AuthHandler
)

func (s *AuthHandler) GetCaptcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			ss_log.Info("======================================================----------------------------")
			response, err := s.Client.GetCaptcha(context.TODO(), &go_micro_srv_auth.GetCaptchaRequet{
				Strlen: 4,
			})
			if err != nil {
				ss_log.Info("GetCaptcha----------------->\nerr=[%v],resp=[%v]", err, response)
				return response.ResultCode, gin.H{}, nil
			}

			common.ConfigC.Content = response.Base64Png
			_, digitCap := base64Captcha.GenerateCaptcha(response.Base64Png, common.ConfigC)
			base64Png := base64Captcha.CaptchaWriteToBase64Encoding(digitCap)

			return ss_err.ERR_SUCCESS, gin.H{
				"verifyid":  response.Verifyid,
				"base64png": base64Png,
			}, nil
		}, "params")
	}
}

/**
 * 登陆处理
 */
func (s *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			ss_log.Debug("in=%v\n", params)
			xRealIp := c.Request.Header.Get("X-Real-IP")
			if xRealIp == "" {
				xRealIp = c.ClientIP()
			}
			c.Set("ip", xRealIp)

			account := container.GetValFromMapMaybe(params, "account").ToStringNoPoint()
			if account == "" {
				ss_log.Error("account参数为空")
				return ss_err.ERR_PARAM, nil, nil
			}
			account = strings.ToLower(account)
			reply, err := s.Client.BusinessLogin(context.TODO(), &go_micro_srv_auth.BusinessLoginRequest{
				Account:   account,
				Password:  container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				Verifyid:  container.GetValFromMapMaybe(params, "verifyid").ToStringNoPoint(),
				Verifynum: container.GetValFromMapMaybe(params, "verifynum").ToStringNoPoint(),
				//LoginAccType: container.GetValFromMapMaybe(params, "login_acc_type").ToStringNoPoint(),
				Nonstr: container.GetValFromMapMaybe(params, "nonstr").ToStringNoPoint(),
				Client: constants.LOGIN_CLI_WEB,
				Ip:     xRealIp,
			})
			ss_log.Info("login------------------------------------->[%v]", reply)
			if reply.AccountUid != "" {
				replyAccount, err := s.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
					AccountUid:  reply.AccountUid,
					AccountType: reply.AccountType,
				})

				idenNo := replyAccount.MerchantUid
				idenType := reply.AccountType
				isMasterAcc := "1"
				if replyAccount.MasterAcc != "" && replyAccount.MasterAcc != ss_sql.UUID {
					getAccMaster, err := s.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
						AccountUid: replyAccount.MasterAcc,
					})

					if getAccMaster == nil {
						return reply.ResultCode, gin.H{}, err
					}

					idenNo = getAccMaster.MerchantUid
					idenType = getAccMaster.AccountType
					isMasterAcc = "0"
				}

				routes := processor.MenuProcInst.TransMenu(replyAccount.DataList)
				ss_log.Info("MerchantUid=[%v]|AccountType=[%v]", replyAccount.MerchantUid, replyAccount.AccountType)

				// 强制登录
				//isForce := container.GetValFromMapMaybe(params, "is_force").ToInt32()
				routesStr, _ := json.Marshal(routes)
				loginToken := strext.NewUUIDNoSplit()
				//errCode := dao.AccDaoInstance.InsertLoginToken(reply.AccountUid, strext.ToStringNoPoint(routesStr), loginToken, isForce, xRealIp)
				errCode := dao.AccDaoInstance.InsertLoginToken(reply.AccountUid, strext.ToStringNoPoint(routesStr), loginToken, 1, xRealIp)
				if errCode != ss_err.ERR_SUCCESS {
					return errCode, gin.H{}, ss_err.ErrLoginFailed
				}

				phone := ""
				switch idenType {
				case constants.AccountType_PersonalBusiness:
					phone = replyAccount.Phone
				case constants.AccountType_EnterpriseBusiness:
					phone = replyAccount.BusinessPhone

				}

				// 登录成功
				// ss_net.SetJwtAuthentication(c, reply.Jwt)
				return strext.ToString(reply.ResultCode), gin.H{
					"userinfo": gin.H{
						"username":            replyAccount.Account,
						"email":               replyAccount.Email,
						"phone":               phone,
						"country_code":        replyAccount.CountryCode,
						"token":               reply.Jwt,
						"init_pay_pwd_status": reply.InitPayPwdStatus,
						"access_status":       2,
						"is_update_pass":      1,
						"routes":              routes,
						"account_uid":         replyAccount.Uid,
						"merchant_uid":        idenNo,
						"account_type":        idenType,
						"top_status":          "1",
						"side_status":         "1",
						"has_pay_password":    reply.InitPayPwdStatus,
						"is_master_acc":       isMasterAcc,
						"login_token":         loginToken,
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

func (a *AuthHandler) SaveBusinessAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			xRealIp := c.Request.Header.Get("X-Real-IP")
			if xRealIp == "" {
				xRealIp = c.ClientIP()
			}
			c.Set("ip", xRealIp)

			//先确认输入的邮箱验证码是否正确
			email := container.GetValFromMapMaybe(params, "email").ToStringNoPoint()
			reqCheckMailCode := &go_micro_srv_cust.CheckMailCodeRequest{
				Mail:     strings.ToLower(email), //邮箱不区分大小写
				MailCode: container.GetValFromMapMaybe(params, "mail_code").ToStringNoPoint(),
				Function: constants.Reg_By_Mail, //功能,注册-reg_mail;
			}
			if errStr := verify.CheckMailCodeVerify(reqCheckMailCode); errStr != "" {
				return ss_err.ERR_PARAM, nil, nil
			}
			replyCode, errCode := CustHandlerInst.Client.CheckMailCode(context.TODO(), reqCheckMailCode)
			if errCode != nil {
				ss_log.Error("调用验证邮箱验证码Api出错，err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}
			if replyCode.ResultCode != ss_err.ERR_SUCCESS {
				return replyCode.ResultCode, "", nil
			}

			//保存账号
			req := &go_micro_srv_auth.SaveBusinessAccountRequest{
				Email:       email,
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				Phone:       container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				CountryCode: container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(),
				Client:      constants.LOGIN_CLI_WEB,
				Ip:          xRealIp,
			}
			// 参数校验
			if errStr := verify.SaveBusinessAccountVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			if errStr := ss_func.CheckCountryCode(req.CountryCode); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("国家码[%v]不合法", req.CountryCode)
				return errStr, nil, nil
			}

			req.Phone = ss_func.PrePhone(req.CountryCode, req.Phone)

			reply, err := a.Client.SaveBusinessAccount(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			ss_log.Info("login------------------------------------->[%v]", reply)
			if reply.Uid != "" {
				replyAccount, err := a.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
					AccountUid:  reply.Uid,
					AccountType: reply.AccountType,
				})

				idenNo := replyAccount.MerchantUid
				idenType := reply.AccountType
				isMasterAcc := "1"
				if replyAccount.MasterAcc != "" && replyAccount.MasterAcc != ss_sql.UUID {
					getAccMaster, err := a.Client.GetAccount(context.TODO(), &go_micro_srv_auth.GetAccountRequest{
						AccountUid: replyAccount.MasterAcc,
					})

					if getAccMaster == nil {
						return reply.ResultCode, gin.H{}, err
					}

					idenNo = getAccMaster.MerchantUid
					idenType = getAccMaster.AccountType
					isMasterAcc = "0"
				}

				routes := processor.MenuProcInst.TransMenu(replyAccount.DataList)
				ss_log.Info("MerchantUid=[%v]|AccountType=[%v]", replyAccount.MerchantUid, replyAccount.AccountType)

				// 强制登录
				isForce := container.GetValFromMapMaybe(params, "is_force").ToInt32()
				routesStr, _ := json.Marshal(routes)
				loginToken := strext.NewUUIDNoSplit()
				errCode := dao.AccDaoInstance.InsertLoginToken(reply.Uid, strext.ToStringNoPoint(routesStr), loginToken, isForce, xRealIp)
				if errCode != ss_err.ERR_SUCCESS {
					return errCode, gin.H{}, ss_err.ErrLoginFailed
				}

				phone := ""
				switch idenType {
				case constants.AccountType_PersonalBusiness:
					phone = replyAccount.Phone
				case constants.AccountType_EnterpriseBusiness:
					phone = replyAccount.BusinessPhone
				}

				// 登录成功
				return strext.ToString(reply.ResultCode), gin.H{
					"userinfo": gin.H{
						"username":            replyAccount.Account,
						"email":               replyAccount.Email,
						"phone":               phone,
						"country_code":        replyAccount.CountryCode,
						"token":               reply.Jwt,
						"init_pay_pwd_status": reply.InitPayPwdStatus,
						"access_status":       2,
						"is_update_pass":      1,
						"routes":              routes,
						"account_uid":         replyAccount.Uid,
						"merchant_uid":        idenNo,
						"account_type":        idenType,
						"top_status":          "1",
						"side_status":         "1",
						"has_pay_password":    reply.InitPayPwdStatus,
						"is_master_acc":       isMasterAcc,
						"login_token":         loginToken,
					},
				}, err
			} else {
				// 登录失败
				// ss_net.SetJwtAuthentication(c, reply.Jwt)
				return reply.ResultCode, gin.H{}, err
			}

		})
	}
}

func (a *AuthHandler) SendMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			function := container.GetValFromMapMaybe(params, "function").ToStringNoPoint() //参考
			email := ""

			switch function {
			case constants.Reg_By_Mail: //注册
				fallthrough
			case constants.Backpwd_By_Mail: //修改登录密码(现阶段使用邮箱验证码修改登录密码是在忘记密码那修改)
				fallthrough
			case constants.ModifyEmail_By_NewEmail: //修改邮箱，发送邮件到新邮箱
				email = container.GetValFromMapMaybe(params, "email").ToStringNoPoint()
			default:
				email = inner_util.GetJwtDataString(c, "email")
			}
			email = strings.ToLower(email)
			req := &go_micro_srv_cust.SendMailRequest{
				Email:    email,
				Lang:     constants.LangEnUS,
				Function: function,
			}

			//参数校验
			if errStr := verify.SendMailReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := CustHandlerInst.Client.SendMail(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			return reply.ResultCode, nil, err
		})
	}
}

func (a *AuthHandler) SendSms() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			function := container.GetValFromMapMaybe(params, "function").ToStringNoPoint()

			phone := ""
			countryCode := ""

			switch function {
			case "":
				ss_log.Error("function参数为空")
				return ss_err.ERR_PARAM, nil, nil
			case constants.BACKPWD_Business: //找回登录密码
				fallthrough
			case constants.MODIFYPHONE_Business: //修改手机号，向新手机发送验证码
				phone = container.GetValFromMapMaybe(params, "phone").ToStringNoPoint()              //找回登录密码 手机号
				countryCode = container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint() //找回登录密码 国家码

				phone = ss_func.PrePhone(countryCode, phone)

			default:
				phone = inner_util.GetJwtDataString(c, "phone")
				countryCode = inner_util.GetJwtDataString(c, "country_code")
			}

			req := &go_micro_srv_cust.RegSmsRequest{
				Phone: phone,
				//Lang:        ss_net.GetCommonData(c).Lang,
				Lang:        constants.LangEnUS, //前端未做多语言，暂时先写死
				Function:    function,
				CountryCode: countryCode,
			}

			//参数校验
			if errStr := verify.SmsReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := CustHandlerInst.Client.RegSms(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

func (s *AuthHandler) CheckSms() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			function := container.GetValFromMapMaybe(params, "function").ToStringNoPoint()
			phone := ""
			switch function {
			case constants.BACKPWD_Business:
				phone = container.GetValFromMapMaybe(params, "phone").ToStringNoPoint()               //找回登录密码 手机号
				countryCode := container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint() //找回登录密码 国家码

				phone = ss_func.PrePhone(countryCode, phone)

			default:
				phone = inner_util.GetJwtDataString(c, "phone")
			}

			req := &go_micro_srv_cust.CheckSmsRequest{
				Phone:    phone,
				Sms:      container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),
				Function: function,
			}
			// 参数校验
			if errStr := verify.CheckSmsReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}
			reply, err := CustHandlerInst.Client.CheckSms(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

func (s *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNo := inner_util.GetJwtDataString(c, "account_uid")
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			dao.AccDaoInstance.DeleteLoginToken(accountNo)
			return ss_err.ERR_SUCCESS, "", nil
		})
	}
}

// 修改登录密码
func (s *AuthHandler) ModifyLoginPWD() gin.HandlerFunc {
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

			reply, err := s.Client.MobileModifyPwd(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			if reply.ResultCode == ss_err.ERR_SUCCESS {
				//此处是将修改密码的账户登出（踢出登录）
				//ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
				dao.AccDaoInstance.DeleteLoginToken(uid)
				//	return ss_err.ERR_SUCCESS, "", nil
				//})
			}

			return reply.ResultCode, "", nil
		})
	}
}

// 使用手机验证码方式修改支付密码
func (s *AuthHandler) ModifyPayPWD() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.MobileModifyPayPwdRequest{
				Uid:         inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				Sms:         container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),      //短信验证密码
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(), //新支付密码
			}

			// 参数校验
			if errStr := verify.ModifyPayPWDReqVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.MobileModifyPayPwd(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			ss_log.Info("reply=[%v]", reply)
			return reply.ResultCode, "", nil
		})
	}
}

//使用邮箱验证码方式修改支付密码
func (s *AuthHandler) ModifyPayPWDByMailCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.ModifyPayPWDByMailCodeRequest{
				Uid:         inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				MailCode:    container.GetValFromMapMaybe(params, "mail_code").ToStringNoPoint(),
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(), //新密码
			}

			// 参数校验
			if errStr := verify.ModifyPayPWDByMailCodeVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.ModifyPayPWDByMailCode(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

//使用原支付密码方式修改支付密码
func (s *AuthHandler) ModifyPayPWDByOldPwd() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.ModifyPayPWDByOldPwdRequest{
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				AccountUid:  inner_util.GetJwtDataString(c, "account_uid"),
				OldPayPwd:   container.GetValFromMapMaybe(params, "old_pay_pwd").ToStringNoPoint(), //旧支付密码
				NewPayPwd:   container.GetValFromMapMaybe(params, "new_pay_pwd").ToStringNoPoint(), //新支付密码
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),     //加密旧支付密码的随机数
			}

			// 参数校验
			if errStr := verify.ModifyPayPWDByOldPwdVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.ModifyPayPWDByOldPwd(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

//使用登录密码修改邮箱
func (s *AuthHandler) ModifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountUid := inner_util.GetJwtDataString(c, "account_uid")
			accountType := inner_util.GetJwtDataString(c, "account_type")
			req := &go_micro_srv_auth.ModifyEmailRequest{
				Uid:         accountUid,
				AccountType: accountType,
				MailCode:    container.GetValFromMapMaybe(params, "mail_code").ToStringNoPoint(), //新邮箱收到的验证码
				Email:       container.GetValFromMapMaybe(params, "email").ToStringNoPoint(),     //新邮箱
				Password:    container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),  //登录密码
				NonStr:      container.GetValFromMapMaybe(params, "non_str").ToStringNoPoint(),   //加密登录密码的随机数

				Account:     inner_util.GetJwtDataString(c, "account"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				Phone:       inner_util.GetJwtDataString(c, "phone"),
				CountryCode: inner_util.GetJwtDataString(c, "country_code"),
			}

			// 参数校验
			if errStr := verify.ModifyEmailVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.ModifyEmail(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			//企业商家的账号就是邮箱，修改邮箱后会修改账号，所以当前账号会被踢下线
			if reply.ResultCode == ss_err.ERR_SUCCESS && accountType == constants.AccountType_EnterpriseBusiness {
				if err := dao.AccDaoInstance.DeleteLoginToken(accountUid); err != nil {
					ss_log.Error("删除账号uid[%v]的LoginToken失败，err=[%v]", accountUid, err)
				}
			}

			return reply.ResultCode, gin.H{
				"token": reply.Token,
			}, nil
		})
	}
}

//修改手机号
func (s *AuthHandler) ModifyPhone() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountUid := inner_util.GetJwtDataString(c, "account_uid")
			accountType := inner_util.GetJwtDataString(c, "account_type")
			req := &go_micro_srv_auth.BusinessModifyPhoneRequest{
				Uid:         accountUid,
				AccountType: accountType,
				Account:     inner_util.GetJwtDataString(c, "account"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				Email:       inner_util.GetJwtDataString(c, "email"),

				Sms:         container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),          //新手机号收到的验证码
				Phone:       container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),        //新手机号
				CountryCode: container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint(), //国家码
			}

			// 参数校验
			if errStr := verify.BusinessModifyPhoneRequestVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			if errStr := ss_func.CheckCountryCode(req.CountryCode); errStr != ss_err.ERR_SUCCESS {
				ss_log.Error("国家码[%v]不合法", req.CountryCode)
				return errStr, nil, nil
			}

			req.Phone = ss_func.PrePhone(req.CountryCode, req.Phone)

			reply, err := s.Client.BusinessModifyPhone(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			//个人商家的账号就是国家码加手机号，修改手机号后会修改账号，所以当前账号会被踢下线
			if reply.ResultCode == ss_err.ERR_SUCCESS && accountType == constants.AccountType_PersonalBusiness {
				if err := dao.AccDaoInstance.DeleteLoginToken(accountUid); err != nil {
					ss_log.Error("删除账号uid[%v]的LoginToken失败，err=[%v]", accountUid, err)
				}
			}

			return reply.ResultCode, gin.H{
				"token": reply.Token,
			}, nil
		})
	}
}

// 使用手机验证码方式修改登录密码（NoToken）
func (s *AuthHandler) NoTokenModifyPWD() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {

			req := &go_micro_srv_auth.BusinessModifyPWDBySmsRequest{
				Phone:    container.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),    //手机号
				Account:  container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),  //账号
				Sms:      container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),      //短信验证密码
				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(), //新登录密码
			}

			// 参数校验
			if errStr := verify.BusinessNotokenModifyPWDBySmsVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			countryCode := container.GetValFromMapMaybe(params, "country_code").ToStringNoPoint() //手机号
			req.Phone = ss_func.PrePhone(countryCode, req.Phone)

			reply, err := s.Client.BusinessModifyPWDBySms(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

// 使用手机验证码方式修改登录密码
func (s *AuthHandler) ModifyPWD() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.BusinessModifyPWDBySmsRequest{
				Uid:      inner_util.GetJwtDataString(c, "account_uid"),
				Phone:    inner_util.GetJwtDataString(c, "phone"),
				Sms:      container.GetValFromMapMaybe(params, "sms").ToStringNoPoint(),      //短信验证密码
				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(), //新登录密码
			}

			// 参数校验
			if errStr := verify.BusinessModifyPWDBySmsVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.BusinessModifyPWDBySms(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

//使用邮箱验证码方式修改登录密码(NoToken)
func (s *AuthHandler) NoTokenModifyPWDByMailCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.ModifyPWDByMailCodeRequest{
				Account: container.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				//Email:    container.GetValFromMapMaybe(params, "email").ToStringNoPoint(),
				MailCode: container.GetValFromMapMaybe(params, "mail_code").ToStringNoPoint(),
				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(), //新密码
			}

			req.Email = req.Account

			// 参数校验
			if errStr := verify.NoTokenModifyPWDByMailCodeVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.ModifyPWDByMailCode(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

//使用邮箱验证码方式修改登录密码
func (s *AuthHandler) ModifyPWDByMailCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.ModifyPWDByMailCodeRequest{
				Uid:      inner_util.GetJwtDataString(c, "account_uid"),
				Email:    inner_util.GetJwtDataString(c, "email"),
				MailCode: container.GetValFromMapMaybe(params, "mail_code").ToStringNoPoint(),
				Password: container.GetValFromMapMaybe(params, "password").ToStringNoPoint(), //新密码
			}

			// 参数校验
			if errStr := verify.ModifyPWDByMailCodeVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := s.Client.ModifyPWDByMailCode(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})
	}
}

func (s *AuthHandler) InitPayPwd() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &go_micro_srv_auth.InitBusinessPayPwdRequest{
				Uid:         inner_util.GetJwtDataString(c, "account_uid"),
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				IdenNo:      inner_util.GetJwtDataString(c, "iden_no"),
				PayPwd:      container.GetValFromMapMaybe(params, "pay_pwd").ToStringNoPoint(), //新密码
			}

			// 参数校验
			if errStr := verify.InitBusinessPayPwdVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := AuthHandlerInst.Client.InitBusinessPayPwd(context.TODO(), req)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, "", nil
		})

	}
}
