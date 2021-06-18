package handler

import (
	"context"
	"encoding/json"

	"a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/common"
	"a.a/mp-server/api-webadmin/dao"
	"a.a/mp-server/api-webadmin/handler/processor"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
	"github.com/wiwii/base64Captcha"
)

var (
	AdminAuthHandlerInst AdminAuthHandler
)

type AdminAuthHandler struct {
	Client adminAuthProto.AdminAuthService
}

func (a *AdminAuthHandler) GetCaptcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			ss_log.Info("======================================================----------------------------")
			response, err := a.Client.GetCaptcha(context.TODO(), &adminAuthProto.GetCaptchaRequet{
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
func (a *AdminAuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(params interface{}) (string, gin.H, error) {
			ss_log.Debug("in=%v\n", params)
			xRealIp := c.Request.Header.Get("X-Real-IP")
			if xRealIp == "" {
				xRealIp = c.ClientIP()
			}
			Account := container.GetValFromMapMaybe(params, "account").ToStringNoPoint()
			c.Set("ip", xRealIp)
			reply, err := a.Client.Login(context.TODO(), &adminAuthProto.LoginRequest{
				Account:   Account,
				Password:  container.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				Verifyid:  container.GetValFromMapMaybe(params, "verifyid").ToStringNoPoint(),
				Verifynum: container.GetValFromMapMaybe(params, "verifynum").ToStringNoPoint(),
				Nonstr:    container.GetValFromMapMaybe(params, "nonstr").ToStringNoPoint(),
				Ip:        xRealIp,
			})
			ss_log.Info("login=[%v]", reply)
			if reply.AccountUid != "" {
				replyAccount, err := a.Client.GetAdminAccount(context.TODO(), &adminAuthProto.GetAdminAccountRequest{
					AccountUid: reply.AccountUid,
					//AccountType: reply.AccountType,
				})

				idenType := reply.AccountType

				routes := processor.MenuProcInst.TransMenu(replyAccount.DataList)

				// 强制登录
				isForce := container.GetValFromMapMaybe(params, "is_force").ToInt32()
				routesStr, _ := json.Marshal(routes)
				loginToken := strext.NewUUIDNoSplit()
				errCode := dao.AccDaoInstance.InsertLoginToken(reply.AccountUid, strext.ToStringNoPoint(routesStr), loginToken, isForce, xRealIp)
				if errCode != ss_err.ERR_SUCCESS {
					return errCode, gin.H{}, ss_err.ErrLoginFailed
				}

				// 登录成功
				// ss_net.SetJwtAuthentication(c, reply.Jwt)
				return strext.ToString(reply.ResultCode), gin.H{
					"userinfo": gin.H{
						"username":         replyAccount.Account,
						"token":            reply.Jwt,
						"access_status":    2,
						"is_update_pass":   1,
						"routes":           routes,
						"account_uid":      replyAccount.Uid,
						"account_type":     idenType,
						"top_status":       "1",
						"side_status":      "1",
						"has_pay_password": "0",
						"login_token":      loginToken,
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
