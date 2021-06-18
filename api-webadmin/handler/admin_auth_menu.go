package handler

import (
	"context"
	"strings"

	colloection "a.a/cu/container"
	"a.a/cu/strext"
	adminAuthProto "a.a/mp-server/common/proto/admin-auth"
	"a.a/mp-server/common/ss_net"
	"github.com/gin-gonic/gin"
)

func (a *AdminAuthHandler) GetAdminMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(urlId interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAdminMenu(context.TODO(), &adminAuthProto.GetAdminMenuRequest{
				UrlUid: strext.ToString(urlId),
			})
			return reply.ResultCode, reply.Data, err
		}, "url_uid")
	}
}

func (a *AdminAuthHandler) GetAdminMenuList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetAdminMenuList(context.TODO(), &adminAuthProto.GetAdminMenuListRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Search:   params[2],
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "search")
	}
}

func (a *AdminAuthHandler) SaveOrInsertAdminMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.SaveOrInsertAdminMenu(context.TODO(), &adminAuthProto.SaveOrInsertAdminMenuRequest{
				UrlUid:        colloection.GetValFromMapMaybe(params, "url_uid").ToString(),
				UrlName:       colloection.GetValFromMapMaybe(params, "url_name").ToString(),
				Url:           colloection.GetValFromMapMaybe(params, "url").ToString(),
				ParentUid:     colloection.GetValFromMapMaybe(params, "parent_uid").ToString(),
				Title:         colloection.GetValFromMapMaybe(params, "title").ToString(),
				Icon:          colloection.GetValFromMapMaybe(params, "icon").ToString(),
				ComponentName: colloection.GetValFromMapMaybe(params, "component_name").ToString(),
				ComponentPath: colloection.GetValFromMapMaybe(params, "component_path").ToString(),
				Redirect:      colloection.GetValFromMapMaybe(params, "redirect").ToString(),
				Idx:           colloection.GetValFromMapMaybe(params, "idx").ToInt32(),
				IsHidden:      colloection.GetValFromMapMaybe(params, "is_hidden").ToInt32(),
			})
			return reply.ResultCode, reply.Uid, err
		})
	}
}

func (a *AdminAuthHandler) DeleteAdminMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoDelete(c, func(params interface{}) (string, error) {
			reply, err := AdminAuthHandlerInst.Client.DeleteAdminMenu(context.TODO(), &adminAuthProto.DeleteAdminMenuRequest{
				UrlUid: colloection.GetValFromMapMaybe(params, "url_uid").ToString(),
			})
			return reply.ResultCode, err
		})
	}
}

func (a *AdminAuthHandler) AdminMenuRefreshChild() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			urls := colloection.GetValFromMapMaybe(params, "url_nos").ToStringNoPoint()
			urls2 := strings.Split(urls, ",")

			reply, err := AdminAuthHandlerInst.Client.AdminMenuRefreshChild(context.TODO(), &adminAuthProto.AdminMenuRefreshChildRequest{
				UrlNo: urls2,
			})
			return reply.ResultCode, "", err
		})
	}
}

func (a *AdminAuthHandler) GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetSingle(c, func(urlId interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.GetMenu(context.TODO(), &adminAuthProto.GetMenuRequest{
				UrlUid: strext.ToString(urlId),
			})
			return reply.ResultCode, reply.Data, err
		}, "url_uid")
	}
}

func (a *AdminAuthHandler) GetMenuList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoGetList(c, func(params []string) (string, interface{}, int32, error) {
			reply, err := AdminAuthHandlerInst.Client.GetMenuList(context.TODO(), &adminAuthProto.GetMenuListRequest{
				Page:     strext.ToInt32(params[0]),
				PageSize: strext.ToInt32(params[1]),
				Search:   params[2],
			})
			return reply.ResultCode, reply.DataList, reply.Total, err
		}, "page", "page_size", "search")
	}
}

func (a *AdminAuthHandler) SaveOrInsertMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			reply, err := AdminAuthHandlerInst.Client.SaveOrInsertMenu(context.TODO(), &adminAuthProto.SaveOrInsertMenuRequest{
				UrlUid:        colloection.GetValFromMapMaybe(params, "url_uid").ToString(),
				UrlName:       colloection.GetValFromMapMaybe(params, "url_name").ToString(),
				Url:           colloection.GetValFromMapMaybe(params, "url").ToString(),
				ParentUid:     colloection.GetValFromMapMaybe(params, "parent_uid").ToString(),
				Title:         colloection.GetValFromMapMaybe(params, "title").ToString(),
				Icon:          colloection.GetValFromMapMaybe(params, "icon").ToString(),
				ComponentName: colloection.GetValFromMapMaybe(params, "component_name").ToString(),
				ComponentPath: colloection.GetValFromMapMaybe(params, "component_path").ToString(),
				Redirect:      colloection.GetValFromMapMaybe(params, "redirect").ToString(),
				Idx:           colloection.GetValFromMapMaybe(params, "idx").ToInt32(),
				IsHidden:      colloection.GetValFromMapMaybe(params, "is_hidden").ToInt32(),
			})
			return reply.ResultCode, reply.Uid, err
		})
	}
}

func (a *AdminAuthHandler) DeleteMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoDelete(c, func(params interface{}) (string, error) {
			reply, err := AdminAuthHandlerInst.Client.DeleteMenu(context.TODO(), &adminAuthProto.DeleteMenuRequest{
				UrlUid: colloection.GetValFromMapMaybe(params, "url_uid").ToString(),
			})
			return reply.ResultCode, err
		})
	}
}

func (a *AdminAuthHandler) MenuRefreshChild() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss_net.NetU.DoUpdate(c, func(params interface{}) (string, interface{}, error) {
			urls := colloection.GetValFromMapMaybe(params, "url_nos").ToStringNoPoint()
			urls2 := strings.Split(urls, ",")

			reply, err := AdminAuthHandlerInst.Client.MenuRefreshChild(context.TODO(), &adminAuthProto.MenuRefreshChildRequest{
				UrlNo: urls2,
			})
			return reply.ResultCode, "", err
		})
	}
}
