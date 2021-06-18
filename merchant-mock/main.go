package main

import (
	"database/sql"
	"log"

	"a.a/mp-server/merchant-mock/conf"
	"a.a/mp-server/merchant-mock/dao"
	"a.a/mp-server/merchant-mock/router"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	engine := gin.New()

	dao.InitDB()

	// 获取配置的应用相关信息
	c, err := dao.AppInstance.GetUsingRow()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatal("找不到配置的应用")
		} else {
			log.Fatal("查询配置的应用失败,err:" + err.Error())
		}
	}

	log.Printf("应用ID:%v, 应用名称:%v \n", c.AppId, c.AppName)

	// 设置配置信息
	conf.SetConfig(c.AppId, c.AppName, c.MerchantPrivateKey, c.PlatformPublicKey)

	// 1.添加路由
	router.AddRouter(engine)

	// 2.添加模板文件目录
	engine.LoadHTMLGlob("templates/**/*")

	// 3.添加静态文件目录
	engine.Static("/static", "./static")

	// 4.运行服务
	engine.Run(":" + conf.Port)
}
