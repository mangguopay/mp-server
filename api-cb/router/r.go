package router

import (
	"a.a/mp-server/api-cb/handler"
	mw "a.a/mp-server/api-cb/middleware"
	"github.com/gin-gonic/gin"
)

var (
	billH handler.BusinessBillHandler
)

func InitRouter() *gin.Engine {
	gin.SetMode("release")
	router := gin.New()
	router.Use(gin.Logger(),
		mw.GenTraceNoMwInst.GenTraceNo(),           // 生成记录号
		mw.GetParamsMwInst.GetPostJsonBodyParams(), // 获取post/json信息

		mw.RecoveryMiddleWareInst.Recovery())

	router.POST("/cb/p/:supplierCode", billH.PayCallback())      // 支付
	router.POST("/cb/a/:supplierCode", billH.TransferCallback()) // 代付

	return router
}
