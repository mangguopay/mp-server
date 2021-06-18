package handler

import (
	"a.a/cu/ss_log"
	"context"
	"time"

	colloection "a.a/cu/container"
	"a.a/cu/strext"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_net"
	"a.a/mp-server/common/ss_sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
)

func (a *AdminAuthHandler) GetRoleUrlList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(role_uid interface{}) (string, interface{}, error) {
			roleUid := role_uid.(string)
			if roleUid == "" {
				roleUid = ss_sql.UUID
			}

			var opss client.CallOption = func(o *client.CallOptions) {
				o.RequestTimeout = time.Minute * 5
				o.DialTimeout = time.Minute * 5
			}
			reply, err := AdminAuthHandlerInst.Client.GetRoleUrlList(context.TODO(), &adminAuthProto.GetRoleUrlListRequest{
				RoleNo: roleUid,
			}, opss)
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, reply.DataList, err
		}, "role_no")
	}
}

func (a *AdminAuthHandler) GetRoleInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(role_uid interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.GetRoleInfo(context.TODO(), &adminAuthProto.GetRoleInfoRequest{
				RoleNo: role_uid.(string),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, nil
			}
			return reply.ResultCode, reply.Data, err
		}, "role_no")
	}
}

func (a *AdminAuthHandler) GetRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(roleName interface{}) (string, gin.H, error) {
			reply, err := AdminAuthHandlerInst.Client.GetRole(context.TODO(), &adminAuthProto.GetRoleRequest{
				RoleName: roleName.(string),
			})

			if reply == nil {
				return ss_err.ERR_SYS_NETWORK, nil, err
			}

			// 组装属性组
			var lst []map[string]interface{}
			lst = append(lst, map[string]interface{}{"key": "角色uid", "value": reply.RoleNo})
			lst = append(lst, map[string]interface{}{"key": "角色名", "value": reply.RoleName})
			lst = append(lst, map[string]interface{}{"key": "创建时间", "value": reply.CreateTime})
			lst = append(lst, map[string]interface{}{"key": "修改时间", "value": reply.ModifyTime})
			lst = append(lst, map[string]interface{}{"key": "删除时间", "value": reply.DropTime})
			if 1 == reply.UseStatus {
				lst = append(lst, map[string]interface{}{"key": "使用状态", "value": "启用"})
			} else {
				lst = append(lst, map[string]interface{}{"key": "使用状态", "value": "禁用"})
			}
			lst = append(lst, map[string]interface{}{"key": "备注", "value": reply.Remark})

			// 组装权限组
			var urls []map[string]interface{}
			for _, v := range reply.UrlData {
				urls = append(urls, map[string]interface{}{"url_uid": v.UrlUid, "name": v.Name})
			}

			var urls2 []map[string]interface{}
			for _, v := range reply.UrlData_2 {
				urls2 = append(urls2, map[string]interface{}{"url_uid": v.UrlUid, "name": v.Name})
			}

			return strext.ToString(reply.ResultCode), gin.H{
				"role_uid": reply.RoleNo,
				"attrs":    lst,
				"urls_l":   urls2,
				"urls_r":   urls,
			}, err
		}, "role_name")
	}
}

func (a *AdminAuthHandler) GetRoleList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetRoleList(context.TODO(), &adminAuthProto.GetRoleListRequest{
				Page:      strext.ToInt32(params[0]),
				PageSize:  strext.ToInt32(params[1]),
				Search:    params[2],
				AccType:   params[3],
				MasterAcc: params[4],
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil, 0, nil
			}
			return reply.ResultCode, reply.DataList, reply.Len, nil
		}, "page", "page_size", "search", "acc_type", "master_acc")
	}
}

func (a *AdminAuthHandler) UpdateOrInsertRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.UpdateOrInsertRole(context.TODO(), &adminAuthProto.UpdateOrInsertRoleRequest{
				RoleNo:    colloection.GetValFromMapMaybe(params, "role_no").ToStringNoPoint(),
				RoleName:  colloection.GetValFromMapMaybe(params, "role_name").ToStringNoPoint(),
				AccType:   colloection.GetValFromMapMaybe(params, "acc_type").ToStringNoPoint(),
				AccUid:    decoded.(jwt.MapClaims)["account_uid"].(string),
				MasterAcc: colloection.GetValFromMapMaybe(params, "master_acc").ToStringNoPoint(),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) UpdateOrInsertRoleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.UpdateOrInsertRoleAuth(context.TODO(), &adminAuthProto.UpdateOrInsertRoleAuthRequest{
				RoleNo: colloection.GetValFromMapMaybe(params, "role_no").ToString(),
				Urls:   colloection.GetValFromMapMaybe(params, "urls").ToStringList(),
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) AuthRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.AuthRole(context.TODO(), &adminAuthProto.AuthRoleRequest{
				RoleNo:  colloection.GetValFromMapMaybe(params, "role_no").ToString(),
				AccType: colloection.GetValFromMapMaybe(params, "acc_type").ToString(),
				DefType: colloection.GetValFromMapMaybe(params, "def_type").ToString(),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, 0, nil
			}
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) DeleteRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoDelete(c, func(params interface{}) (string, error) {
			reply, err := AdminAuthHandlerInst.Client.DeleteRole(context.TODO(), &adminAuthProto.DeleteRoleRequest{
				RoleNo: colloection.GetValFromMapMaybe(params, "role_no").ToString(),
			})
			if err != nil {
				ss_log.Error("api调用失败,err=[%v]", err)
				return ss_err.ERR_SYS_REMOTE_API_ERR, nil
			}
			return strext.ToString(reply.ResultCode), err
		})
	}
}

//============== admin ====================

func (a *AdminAuthHandler) GetAdminRoleUrlList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(role_uid interface{}) (string, interface{}, error) {
			roleUid := role_uid.(string)
			if roleUid == "" {
				roleUid = ss_sql.UUID
			}

			var opss client.CallOption = func(o *client.CallOptions) {
				o.RequestTimeout = time.Minute * 5
				o.DialTimeout = time.Minute * 5
			}
			reply, err := AdminAuthHandlerInst.Client.GetAdminRoleUrlList(context.TODO(), &adminAuthProto.GetAdminRoleUrlListRequest{
				RoleNo: roleUid,
			}, opss)
			return reply.ResultCode, reply.DataList, err
		}, "role_no")
	}
}

func (a *AdminAuthHandler) GetAdminRoleInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(role_uid interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAdminRoleInfo(context.TODO(), &adminAuthProto.GetAdminRoleInfoRequest{
				RoleNo: role_uid.(string),
			})
			return reply.ResultCode, reply.Data, err
		}, "role_no")
	}
}

func (a *AdminAuthHandler) GetAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGet(c, func(roleName interface{}) (string, gin.H, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAdminRole(context.TODO(), &adminAuthProto.GetAdminRoleRequest{
				RoleName: roleName.(string),
			})

			if reply == nil {
				return ss_err.ERR_SYS_NETWORK, nil, err
			}

			// 组装属性组
			var lst []map[string]interface{}
			lst = append(lst, map[string]interface{}{"key": "角色uid", "value": reply.RoleNo})
			lst = append(lst, map[string]interface{}{"key": "角色名", "value": reply.RoleName})
			lst = append(lst, map[string]interface{}{"key": "创建时间", "value": reply.CreateTime})
			lst = append(lst, map[string]interface{}{"key": "修改时间", "value": reply.ModifyTime})
			lst = append(lst, map[string]interface{}{"key": "删除时间", "value": reply.DropTime})
			if 1 == reply.UseStatus {
				lst = append(lst, map[string]interface{}{"key": "使用状态", "value": "启用"})
			} else {
				lst = append(lst, map[string]interface{}{"key": "使用状态", "value": "禁用"})
			}
			lst = append(lst, map[string]interface{}{"key": "备注", "value": reply.Remark})

			// 组装权限组
			var urls []map[string]interface{}
			for _, v := range reply.UrlData {
				urls = append(urls, map[string]interface{}{"url_uid": v.UrlUid, "name": v.Name})
			}

			var urls2 []map[string]interface{}
			for _, v := range reply.UrlData_2 {
				urls2 = append(urls2, map[string]interface{}{"url_uid": v.UrlUid, "name": v.Name})
			}

			return strext.ToString(reply.ResultCode), gin.H{
				"role_uid": reply.RoleNo,
				"attrs":    lst,
				"urls_l":   urls2,
				"urls_r":   urls,
			}, err
		}, "role_name")
	}
}

func (a *AdminAuthHandler) GetAdminRoleList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAdminRoleList(context.TODO(), &adminAuthProto.GetAdminRoleListRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Search:   params[2],
				AccType:  params[3],
			})
			return reply.ResultCode, reply.DataList, reply.Len, err
		}, "page", "page_size", "search", "acc_type")
	}
}

func (a *AdminAuthHandler) UpdateOrInsertAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded, _ := c.Get("decodedJwt")
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.UpdateOrInsertAdminRole(context.TODO(), &adminAuthProto.UpdateOrInsertAdminRoleRequest{
				RoleNo:    colloection.GetValFromMapMaybe(params, "role_no").ToStringNoPoint(),
				RoleName:  colloection.GetValFromMapMaybe(params, "role_name").ToStringNoPoint(),
				AccType:   colloection.GetValFromMapMaybe(params, "acc_type").ToStringNoPoint(),
				AccUid:    decoded.(jwt.MapClaims)["account_uid"].(string),
				MasterAcc: colloection.GetValFromMapMaybe(params, "master_acc").ToStringNoPoint(),
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) UpdateOrInsertAdminRoleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.UpdateOrInsertAdminRoleAuth(context.TODO(), &adminAuthProto.UpdateOrInsertAdminRoleAuthRequest{
				RoleNo: colloection.GetValFromMapMaybe(params, "role_no").ToString(),
				Urls:   colloection.GetValFromMapMaybe(params, "urls").ToStringList(),
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) AuthAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.AuthAdminRole(context.TODO(), &adminAuthProto.AuthAdminRoleRequest{
				RoleNo:  colloection.GetValFromMapMaybe(params, "role_no").ToString(),
				AccType: colloection.GetValFromMapMaybe(params, "acc_type").ToString(),
				DefType: colloection.GetValFromMapMaybe(params, "def_type").ToString(),
			})
			return reply.ResultCode, 0, err
		})
	}
}

func (a *AdminAuthHandler) DeleteAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoDelete(c, func(params interface{}) (string, error) {
			reply, err := AdminAuthHandlerInst.Client.DeleteAdminRole(context.TODO(), &adminAuthProto.DeleteAdminRoleRequest{
				RoleNo: colloection.GetValFromMapMaybe(params, "role_no").ToString(),
			})
			return strext.ToString(reply.ResultCode), err
		})
	}
}
