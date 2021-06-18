package handler

import (
	"log"
	"net/http"

	"a.a/mp-server/merchant-mock/conf"

	"a.a/mp-server/merchant-mock/dao"
	"github.com/gin-gonic/gin"
)

// 应用列表
func AppList(c *gin.Context) {
	page := 1
	pageSize := 100

	list, err := dao.AppInstance.GetList(page, pageSize)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	c.HTML(http.StatusOK, "app/list.html", gin.H{
		"title":          "应用列表",
		"appList":        list,
		"errMsg":         errMsg,
		"runningAppID":   conf.AppId,
		"runningAppName": conf.AppName,
	})
}

// 应用-修改使用状态
func AppChangeUse(c *gin.Context) {
	appid := c.PostForm("app_id")
	if appid == "" {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "应用ID为空"})
		return
	}

	// 修改app为可用状态
	if err := dao.AppInstance.SetAppUsing(appid); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 103, "msg": "查询订单失败:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	return
}

// 创建应用-显示
func AppAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "app/add.html", gin.H{
		"title": "创建应用",
	})
}

// 应用详情
func AppDetail(c *gin.Context) {
	appId := c.Query("app_id")
	if appId == "" {
		RedirectTipsError(c, "应用ID为空")
		return
	}

	app, err := dao.AppInstance.GetOneByAppId(appId)
	if err != nil {
		RedirectTipsError(c, "查询应用失败:"+err.Error())
		return
	}

	c.HTML(http.StatusOK, "app/detail.html", gin.H{
		"title": "应用详情",
		"app":   app,
	})
}

// 应用订单
func AppCreate(c *gin.Context) {
	appId := c.PostForm("app_id")
	appName := c.PostForm("app_name")
	merchantPrivateKey := c.PostForm("merchant_private_key")
	merchantPublicKey := c.PostForm("merchant_public_key")
	platformPublicKey := c.PostForm("platform_public_key")
	merchantKeyType := c.PostForm("merchant_key_type")

	if appId == "" {
		RedirectTipsError(c, "应用ID为空")
		return
	}

	if appName == "" {
		RedirectTipsError(c, "应用名称为空")
		return
	}

	if merchantPrivateKey == "" {
		RedirectTipsError(c, "商家私钥为空")
		return
	}

	if merchantPublicKey == "" {
		RedirectTipsError(c, "商家公钥为空")
		return
	}

	if platformPublicKey == "" {
		RedirectTipsError(c, "平台公钥为空")
		return
	}

	order := &dao.App{
		AppId:              appId,
		AppName:            appName,
		MerchantPrivateKey: merchantPrivateKey,
		MerchantPublicKey:  merchantPublicKey,
		PlatformPublicKey:  platformPublicKey,
		MerchantKeyType:    merchantKeyType,
	}

	// 插入应用
	if err := dao.AppInstance.Insert(order); err != nil {
		log.Printf("插入应用失败, err:%v, order:%+v", err, order)
		RedirectTipsError(c, "插入应用失败:"+err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/app/list")
}
