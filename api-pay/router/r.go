package router

import (
	"a.a/mp-server/api-pay/common"
	"a.a/mp-server/api-pay/handler"
	mw "a.a/mp-server/api-pay/middleware"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
)

var (
	billH         handler.BillHandler
	businessBillH handler.BusinessBillHandler
)

func InitRouter() *gin.Engine {
	gin.SetMode("release")
	router := gin.New()
	router.Use(gin.Logger(),
		mw.RespMwInst.Resp(),       // 输出
		mw.SignerMwInst.Sign(),     // 签名
		mw.JsonPackerMwInst.Pack(), // 组织数据
		// ================== 以上中间件在主要用于组织数据返回 ======================

		mw.GenTraceNoMwInst.GenTraceNo(),        // 生成记录号
		mw.GetParamsMwInst.CheckRequestMethod(), // 检测请求方式

		mw.GetParamsMwInst.CheckRoute(),        // 检测路由
		mw.GetParamsMwInst.GetParamsFromPost(), // 获取post/json信息

		mw.VerifyAppMwInst.Check(),    // 验证app相关信息
		mw.SignerMwInst.Verify(),      // 检查签名
		mw.VerifyMwInst.IpWhiteList(), // ip白名单验证
		mw.VerifyMwInst.RequestTime(), // 验证请求的时间

		mw.RecoveryMiddleWareInst.Recovery(),
	)

	router.NoRoute(func(c *gin.Context) {
		c.Set(common.RET_CODE, ss_err.ACErrSysRouteNotFound)
		return
	})

	payGroup := router.Group("/api")

	payGroup.POST("/prepay", businessBillH.Prepay()) // 预下单
	payGroup.POST("/query", businessBillH.Query())   // 查询订单

	payGroup.POST("/refund", businessBillH.Refund())            // 退款申请
	payGroup.POST("/query_refund", businessBillH.QueryRefund()) // 退款订单查询

	payGroup.POST("/transfer/enterprise", billH.EnterpriseTransfer()) //企业转账
	payGroup.POST("/transfer/query", billH.QueryTransfer())           //企业转账查询

	return router
}
