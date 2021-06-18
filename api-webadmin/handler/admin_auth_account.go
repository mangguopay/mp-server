package handler

import (
	"context"
	"errors"
	"fmt"

	colloection "a.a/cu/container"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webadmin/dao"
	"a.a/mp-server/api-webadmin/inner_util"
	"a.a/mp-server/api-webadmin/verify"
	"a.a/mp-server/common/constants"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (a *AdminAuthHandler) GetAccountList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAccountList(context.TODO(), &adminAuthProto.GetAccountListRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				Search:      strext.ToString(params[2]),
				AccountType: strext.ToString(params[3]),
				IsActived:   strext.ToString(params[4]),
			})
			return reply.ResultCode, reply.AccountData, reply.Total, err
		}, "page", "page_size", "search", "account_type", "is_actived")
	}
}

func (a *AdminAuthHandler) GetAdminAccountList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAdminAccountList(context.TODO(), &adminAuthProto.GetAdminAccountListRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				Search:      strext.ToString(params[2]),
				AccountType: strext.ToString(params[3]),
			})
			return reply.ResultCode, reply.AccountData, reply.Total, err
		}, "page", "page_size", "search", "account_type")
	}
}

//客户管理的查询账户
func (a *AdminAuthHandler) GetAccountList2() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAccountList2(context.TODO(), &adminAuthProto.GetAccountListRequest2{
				Page:          strext.ToInt32(params[0]),
				PageSize:      strext.ToInt32(params[1]),
				StartTime:     strext.ToString(params[2]),
				EndTime:       strext.ToString(params[3]),
				QueryNickname: strext.ToString(params[4]),
				QueryPhone:    strext.ToString(params[5]),
				Account:       strext.ToString(params[6]),
				AuthStatus:    strext.ToString(params[7]),
				IsActived:     strext.ToString(params[8]),
				UseStatus:     strext.ToString(params[9]),
				SortType:      strext.ToString(params[10]), //usd_up正向，usd_down反向，khr_up正向，khr_down反向
			})
			if err != nil {
				ss_log.Error("GetAccountList2 err: %s", err.Error())
				return ss_err.ERR_PARAM, nil, 0, nil
			}
			return reply.ResultCode, reply.AccountData, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "query_nickname", "query_phone", "account", "auth_status", "is_actived", "use_status", "sort_type")
	}
}

//未激活的账号
func (a *AdminAuthHandler) GetUnActivedAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetUnActivedAccounts(context.TODO(), &adminAuthProto.GetUnActivedAccountsRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				StartTime: strext.ToString(params[2]),
				EndTime:   strext.ToString(params[3]),
				Account:   strext.ToString(params[4]),
			})
			if err != nil {
				ss_log.Error("调用api失败 err: %s", err.Error())
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.AccountData, reply.Total, err
		}, "page", "page_size", "start_time", "end_time", "account")
	}
}

/**
 * 获取账户
 */
func (a *AdminAuthHandler) GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(AccountUid interface{}) (string, gin.H, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAccount(context.TODO(), &adminAuthProto.GetAccountRequest{
				AccountUid: strext.ToString(AccountUid),
			})
			return strext.ToString(reply.ResultCode), gin.H{
				"info": reply,
			}, err
		}, "uid")
	}
}

/**
 * 获取账户
 */
func (a *AdminAuthHandler) GetAdminAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(AccountUid interface{}) (string, gin.H, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAdminAccount(context.TODO(), &adminAuthProto.GetAdminAccountRequest{
				AccountUid: strext.ToString(AccountUid),
			})

			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			return strext.ToString(reply.ResultCode), gin.H{
				"info": reply,
			}, err
		}, "uid")
	}
}

func (a *AdminAuthHandler) SaveAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		selfAccountType := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_type"])
		masterAcc := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_uid"])
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			if selfAccountType == constants.AccountType_OPERATOR || selfAccountType == constants.AccountType_ADMIN {
				// 运营或admin才能建指定的
				masterAcc = colloection.GetValFromMapMaybe(params, "master_acc").ToStringNoPoint()
				if selfAccountType == "" {
					return ss_err.ERR_PARAM, "", errors.New("noAccountType")
				}
			}

			req := &adminAuthProto.SaveAccountRequest{
				AccountUid:  colloection.GetValFromMapMaybe(params, "uid").ToStringNoPoint(),
				Nickname:    colloection.GetValFromMapMaybe(params, "nickname").ToStringNoPoint(),
				Account:     colloection.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				UseStatus:   colloection.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				Phone:       colloection.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				Email:       colloection.GetValFromMapMaybe(params, "email").ToStringNoPoint(),
				Password:    colloection.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				AccountType: colloection.GetValFromMapMaybe(params, "account_type").ToStringNoPoint(),
				MasterAcc:   masterAcc,
			}

			reply, err := AdminAuthHandlerInst.Client.SaveAccount(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			loginUid := inner_util.GetJwtDataString(c, "account_uid")
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				if req.AccountUid == "" {
					AccountType := ""
					switch req.AccountType {
					case constants.AccountType_ADMIN:
						AccountType = "管理员"
					case constants.AccountType_OPERATOR:
						AccountType = "运营"
					default:
						ss_log.Error("未知AccountType:[%v]", req.AccountType)
					}
					description := fmt.Sprintf("添加新账号 账号:[%v],昵称[%v],手机号[%v],账号类型[%v]", req.Nickname, req.Account, req.Phone, AccountType)
					errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
					if errAddLog != nil {
						ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
					}
				} else {
					//todo 查询旧的账号
					oldAccount, err := dao.AccDaoInstance.GetAccountByUid(req.AccountUid)
					if err != nil {
						ss_log.Error("查询账号出错，err=[%v]", err)
					}

					description := fmt.Sprintf("修改账号[%v]信息成功,新信息为：账号:[%v],昵称[%v],手机号[%v]", oldAccount, req.Account, req.Nickname, req.Phone)
					errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
					if errAddLog != nil {
						ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
					}

				}

			}

			return reply.ResultCode, reply.Uid, err
		})
	}
}

func (a *AdminAuthHandler) SaveAdminAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		selfAccountType := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_type"])
		masterAcc := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_uid"])
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			if selfAccountType == constants.AccountType_OPERATOR || selfAccountType == constants.AccountType_ADMIN {
				// 运营或admin才能建指定的
				masterAcc = colloection.GetValFromMapMaybe(params, "master_acc").ToStringNoPoint()
				if selfAccountType == "" {
					return ss_err.ERR_PARAM, "", errors.New("noAccountType")
				}
			}

			req := &adminAuthProto.SaveAdminAccountRequest{
				AccountUid:  colloection.GetValFromMapMaybe(params, "uid").ToStringNoPoint(),
				Account:     colloection.GetValFromMapMaybe(params, "account").ToStringNoPoint(),
				UseStatus:   colloection.GetValFromMapMaybe(params, "use_status").ToStringNoPoint(),
				Phone:       colloection.GetValFromMapMaybe(params, "phone").ToStringNoPoint(),
				Email:       colloection.GetValFromMapMaybe(params, "email").ToStringNoPoint(),
				Password:    colloection.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				AccountType: colloection.GetValFromMapMaybe(params, "account_type").ToStringNoPoint(),
			}

			reply, err := AdminAuthHandlerInst.Client.SaveAdminAccount(context.TODO(), req)
			if err != nil {
				ss_log.Error("err=[%v]")
				return ss_err.ERR_SYS_REMOTE_API_ERR, "", nil
			}

			loginUid := inner_util.GetJwtDataString(c, "account_uid")
			if reply.ResultCode == ss_err.ERR_SUCCESS {
				if req.AccountUid == "" {
					AccountType := ""
					switch req.AccountType {
					case constants.AccountType_ADMIN:
						AccountType = "管理员"
					case constants.AccountType_OPERATOR:
						AccountType = "运营"
					default:
						ss_log.Error("未知AccountType:[%v]", req.AccountType)
					}
					description := fmt.Sprintf("添加新管理员账号 账号:[%v],手机号[%v],账号类型[%v]", req.Account, req.Phone, AccountType)
					errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
					if errAddLog != nil {
						ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
					}
				} else {
					//todo 查询旧的账号
					oldAccount, err := dao.AccDaoInstance.GetAccountByUid(req.AccountUid)
					if err != nil {
						ss_log.Error("查询账号出错，err=[%v]", err)
					}

					description := fmt.Sprintf("修改账号[%v]信息成功,新信息为：账号:[%v],手机号[%v]", oldAccount, req.Account, req.Phone)
					errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
					if errAddLog != nil {
						ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
					}

				}

			}

			return reply.ResultCode, reply.Uid, err
		})
	}
}

func (a *AdminAuthHandler) DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoDelete(c, func(params interface{}) (string, error) {
			//reply, err := AuthHandlerInst.Client.DeleteAccount(context.TODO(), &adminAuthProto.DeleteAccountRequest{
			//	AccountUids: colloection.GetValFromMapMaybe(params, "uid").ToStringListSplit(","),
			//})
			//return strext.ToString(reply.ResultCode), err
			return strext.ToString(ss_err.ERR_SYS_REMOTE_API_ERR), nil

		})
	}
}

/**
 * 获取账户
 */
func (a *AdminAuthHandler) GetAccountByNickname() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(nickname interface{}) (string, gin.H, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAccountByNickname(context.TODO(), &adminAuthProto.GetAccountByNicknameRequest{
				Nickname: strext.ToString(nickname),
			})

			// 组装属性组
			var lst []map[string]interface{}
			lst = append(lst, map[string]interface{}{"key": "用户uid", "value": reply.Uid})
			lst = append(lst, map[string]interface{}{"key": "用户名", "value": reply.Nickname})
			lst = append(lst, map[string]interface{}{"key": "账号", "value": reply.Account})
			lst = append(lst, map[string]interface{}{"key": "创建时间", "value": reply.CreateTime})
			lst = append(lst, map[string]interface{}{"key": "修改时间", "value": reply.ModifyTime})
			lst = append(lst, map[string]interface{}{"key": "删除时间", "value": reply.DropTime})
			if "1" == reply.UseStatus {
				lst = append(lst, map[string]interface{}{"key": "使用状态", "value": "启用"})
			} else {
				lst = append(lst, map[string]interface{}{"key": "使用状态", "value": "禁用"})
			}

			// 组装权限组
			var roles []map[string]interface{}
			for _, v := range reply.DataList {
				roles = append(roles, map[string]interface{}{"role_uid": v.RoleNo, "role_name": v.RoleName})
			}

			var roles2 []map[string]interface{}
			for _, v := range reply.DataList_2 {
				roles2 = append(roles2, map[string]interface{}{"role_uid": v.RoleNo, "role_name": v.RoleName})
			}

			return strext.ToString(reply.ResultCode), gin.H{
				"uid":     reply.Uid,
				"attrs":   lst,
				"roles_l": roles2,
				"roles_r": roles,
			}, err
		}, "nickname")
	}
}

func (a *AdminAuthHandler) UpdateOrInsertAccountAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.UpdateOrInsertAccountAuth(context.TODO(), &adminAuthProto.UpdateOrInsertAccountAuthRequest{
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				LoginUid:    inner_util.GetJwtDataString(c, "account_uid"),
				Uid:         colloection.GetValFromMapMaybe(params, "uid").ToString(),
				Roles:       colloection.GetValFromMapMaybe(params, "roles").ToString(),
			})
			//return strext.ToString(reply.ResultCode), 0, err
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) UpdateOrInsertAdminAccountAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.UpdateOrInsertAdminAccountAuth(context.TODO(), &adminAuthProto.UpdateOrInsertAdminAccountAuthRequest{
				AccountType: inner_util.GetJwtDataString(c, "account_type"),
				LoginUid:    inner_util.GetJwtDataString(c, "account_uid"),
				Uid:         colloection.GetValFromMapMaybe(params, "uid").ToString(),
				Roles:       colloection.GetValFromMapMaybe(params, "roles").ToString(),
			})
			//return strext.ToString(reply.ResultCode), 0, err
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) ModifyPw() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		accountType := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_type"])
		accountNo := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_uid"])
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			if accountType == constants.AccountType_ADMIN {
				accountNoT := colloection.GetValFromMapMaybe(params, "uid").ToStringNoPoint()
				if accountNoT != "" {
					accountNo = accountNoT
				}
			}
			if accountType == constants.AccountType_OPERATOR {
				accountNoT := colloection.GetValFromMapMaybe(params, "uid").ToStringNoPoint()
				if accountNoT != "" {
					accountNo = accountNoT
				}
			}
			reply, err := AdminAuthHandlerInst.Client.ModifyPw(context.TODO(), &adminAuthProto.ModifyPwRequest{
				Account: accountNo,
				NewPw:   colloection.GetValFromMapMaybe(params, "new_password").ToStringNoPoint(),
				OldPw:   colloection.GetValFromMapMaybe(params, "old_password").ToStringNoPoint(),
				Cli:     constants.LOGIN_CLI_WEB,
			})

			if reply.ResultCode == ss_err.ERR_SUCCESS {
				//此处是将修改密码的账户登出（踢出登录）
				dao.AccDaoInstance.DeleteLoginToken(accountNo)

				//以下是添加关键操作日志
				loginUid := inner_util.GetJwtDataString(c, "account_uid")
				oldAccount, errGet := dao.AccDaoInstance.GetAccountByUid(accountNo)
				if errGet != nil {
					ss_log.Error("获取账号失败，uid:[%v],err=[%v]", accountNo, errGet)
				}
				description := fmt.Sprintf("修改账号[%v]密码", oldAccount)
				errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
				if errAddLog != nil {
					ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
				}
			}

			return reply.ResultCode, "", err
		})
	}
}

func (a *AdminAuthHandler) ModifyAdminPw() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			switch accountType {
			case constants.AccountType_ADMIN:
			case constants.AccountType_OPERATOR:
			default:
				ss_log.Error("jwt获取到的账号类型为[%v],无权限调用接口", accountType)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			accountNo := colloection.GetValFromMapMaybe(params, "uid").ToStringNoPoint()

			reply, err := AdminAuthHandlerInst.Client.ModifyAdminPw(context.TODO(), &adminAuthProto.ModifyAdminPwRequest{
				AccountUid: accountNo,
				NewPw:      colloection.GetValFromMapMaybe(params, "new_password").ToStringNoPoint(),
				OldPw:      colloection.GetValFromMapMaybe(params, "old_password").ToStringNoPoint(),
			})

			if reply.ResultCode == ss_err.ERR_SUCCESS {
				//此处是将修改密码的账户登出（踢出登录）
				dao.AccDaoInstance.DeleteLoginToken(accountNo)

				//以下是添加关键操作日志
				loginUid := inner_util.GetJwtDataString(c, "account_uid")
				oldAccount, errGet := dao.AccDaoInstance.GetAccountByUid(accountNo)
				if errGet != nil {
					ss_log.Error("获取账号失败，uid:[%v],err=[%v]", accountNo, errGet)
				}
				description := fmt.Sprintf("修改账号[%v]密码", oldAccount)
				errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
				if errAddLog != nil {
					ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
				}
			}

			return reply.ResultCode, "", err
		})
	}
}

func (a *AdminAuthHandler) ResetPw() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		accountType := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_type"])
		accountNo := strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_uid"])
		loginAccNo := accountNo
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			if accountType == constants.AccountType_ADMIN || accountType == constants.AccountType_OPERATOR {
				accountNo = colloection.GetValFromMapMaybe(params, "uid").ToStringNoPoint()
			}
			reply, err := AdminAuthHandlerInst.Client.ResetPw(context.TODO(), &adminAuthProto.ResetPwRequest{
				Account:    accountNo,
				Cli:        constants.LOGIN_CLI_WEB,
				LoginPw:    colloection.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				LoginAccNo: loginAccNo,
			})

			if reply.ResultCode == ss_err.ERR_SUCCESS {
				//此处是将修改密码的账户登出（踢出登录）
				dao.AccDaoInstance.DeleteLoginToken(accountNo)

				//关键操作日志
				loginUid := inner_util.GetJwtDataString(c, "account_uid")
				account, errGet := dao.AccDaoInstance.GetAccountByUid(accountNo)
				if errGet != nil {
					ss_log.Error("获取账号失败，uid=[%v],err=[%v]", accountNo, err)
				}
				description := fmt.Sprintf("重置账号[%v]的登陆密码 ", account)
				errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
				if errAddLog != nil {
					ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
				}
			}

			return reply.ResultCode, "", err
		})
	}
}

func (a *AdminAuthHandler) ResetAdminPw() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			accountType := inner_util.GetJwtDataString(c, "account_type")
			switch accountType {
			case constants.AccountType_ADMIN:
			case constants.AccountType_OPERATOR:
			default:
				ss_log.Error("jwt获取到的账号类型为[%v],无权限调用接口", accountType)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}

			accountNo := colloection.GetValFromMapMaybe(params, "uid").ToStringNoPoint()

			reply, err := AdminAuthHandlerInst.Client.ResetAdminPw(context.TODO(), &adminAuthProto.ResetAdminPwRequest{
				Account:    accountNo,
				LoginPw:    colloection.GetValFromMapMaybe(params, "password").ToStringNoPoint(),
				LoginAccNo: inner_util.GetJwtDataString(c, "account_uid"),
			})

			if reply.ResultCode == ss_err.ERR_SUCCESS {
				//此处是将修改密码的账户登出（踢出登录）
				dao.AccDaoInstance.DeleteLoginToken(accountNo)

				//关键操作日志
				loginUid := inner_util.GetJwtDataString(c, "account_uid")
				account, errGet := dao.AccDaoInstance.GetAdminAccountByUid(accountNo)
				if errGet != nil {
					ss_log.Error("获取账号失败，uid=[%v],err=[%v]", accountNo, err)
				}
				description := fmt.Sprintf("重置账号[%v]的登陆密码 ", account)
				errAddLog := dao.LogDaoInstance.InsertWebAccountLog(description, loginUid, constants.LogAccountWebType_Account_Menu)
				if errAddLog != nil {
					ss_log.Error("插入Web操作日志[%v]失败--------------->err=[%v],", description, errAddLog)
				}
			}

			return reply.ResultCode, "", err
		})
	}
}

func (a *AdminAuthHandler) GetLogLoginList() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetLogLoginList(context.TODO(), &adminAuthProto.GetLogLoginListRequest{
				Page:        strext.ToInt32(params[0]),
				PageSize:    strext.ToInt32(params[1]),
				AccountNo:   strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_uid"]),
				AccountType: strext.ToStringNoPoint(decoded.(jwt.MapClaims)["account_type"]),
				Search:      strext.ToString(params[4]),
				StartTime:   strext.ToString(params[5]),
				EndTime:     strext.ToString(params[6]),
			})
			return reply.ResultCode, reply.Datas, reply.Total, err
		}, "page", "page_size", "account_uid", "account_type", "search", "start_time", "end_time")
	}
}

func (a *AdminAuthHandler) ModifyUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			req := &adminAuthProto.ModifyUserStatusRequest{
				LoginUid:  inner_util.GetJwtDataString(c, "account_uid"), //进行操作的账号uid
				Uid:       colloection.GetValFromMapMaybe(params, "uid").ToString(),
				SetStatus: colloection.GetValFromMapMaybe(params, "set_status").ToString(),
			}
			if errStr := verify.CheckModifyUserStatusVerify(req); errStr != "" {
				return errStr, nil, nil
			}

			reply, err := AdminAuthHandlerInst.Client.ModifyUserStatus(context.TODO(), req)
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) CheckAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(account interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.CheckAccount(context.TODO(), &adminAuthProto.CheckAccountRequest{
				Account: strext.ToString(account),
			})
			return reply.ResultCode, reply.Data, err
		}, "account")
	}
}

func (a *AdminAuthHandler) CheckAdminAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(account interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.CheckAdminAccount(context.TODO(), &adminAuthProto.CheckAdminAccountRequest{
				Account: strext.ToString(account),
			})
			return reply.ResultCode, reply.Data, err
		}, "account")
	}
}

func (a *AdminAuthHandler) GetRoleFromAcc() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			accountType := decoded.(jwt.MapClaims)["account_type"].(string)
			if accountType == constants.AccountType_ADMIN {
			} else if accountType == constants.AccountType_OPERATOR {
			} else {
				return ss_err.ERR_ACCOUNT_NO_PERMISSION, nil, 0, nil
			}
			reply, err := AdminAuthHandlerInst.Client.GetRoleFromAcc(context.TODO(), &adminAuthProto.GetRoleFromAccRequest{
				AccNo: params[0],
			})
			return reply.ResultCode, reply.Datas, 0, err
		}, "acc_no")
	}
}

func (a *AdminAuthHandler) GetAdminRoleFromAcc() gin.HandlerFunc {
	return func(c *gin.Context) {
		//decoded, _ := c.Get("decodedJwt")
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			//accountType := decoded.(jwt.MapClaims)["account_type"].(string)
			accountType := inner_util.GetJwtDataString(c, "account_type")

			switch accountType {
			case constants.AccountType_ADMIN:
			case constants.AccountType_OPERATOR:
			default:
				ss_log.Error("登录的账号角色[%v],无权限调用接口", accountType)
				return ss_err.ERR_ACCOUNT_NO_PERMISSION, nil, 0, nil
			}
			reply, err := AdminAuthHandlerInst.Client.GetAdminRoleFromAcc(context.TODO(), &adminAuthProto.GetAdminRoleFromAccRequest{
				AccNo: params[0],
			})
			return reply.ResultCode, reply.Datas, 0, err
		}, "acc_no")
	}
}
