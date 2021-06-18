package router

import (
	"a.a/mp-server/merchant-mock/handler"
	"github.com/gin-gonic/gin"
)

// AddRouter 添加路由
func AddRouter(r *gin.Engine) {
	r.GET("/", handler.OrderList)

	// 提示信息页面
	r.GET("/order/tips", handler.TipsPage)

	// 订单列表
	r.GET("/order/list", handler.OrderList)

	// 添加订单-显示页面
	r.GET("/order/add", handler.OrderAdd)

	// 添加订单-下单
	r.POST("/order/create", handler.OrderCreate)

	// 显示支付页面
	r.GET("/order/pay", handler.OrderPay)

	// 查询订单支付状态
	r.GET("/order/pay_query", handler.OrderPayQuery)

	// 异步通知
	r.POST("/order/notify", handler.OrderNotify)

	// 同步跳转
	r.GET("/order/jump_back", handler.OrderJumpBack)

	// 订单退款申请
	r.GET("/order/refund", handler.OrderRefund)

	// 订单退款记录列表
	r.GET("/order/refund/list", handler.OrderRefundList)

	// 模拟商家扫码用户收款二维码
	r.GET("/merchant/scan", handler.MerchantScan)

	// 模拟商家扫码用户收款二维码-下单
	r.POST("/merchant/scan_order", handler.MerchantScanOrder)

	// 模拟商家扫码用户收款二维码-下单结果页
	r.GET("/merchant/scan_order_result", handler.MerchantScanOrderResult)

	// 企业付款-列表
	r.GET("/enterprise_transfer/list", handler.EnterpriseTransferList)

	// 企业付款-显示页面
	r.GET("/enterprise_transfer/index", handler.EnterpriseTransferIndex)

	// 企业付款-转账
	r.POST("/enterprise_transfer/do", handler.EnterpriseTransferDo)

	// 应用列表
	r.GET("/app/list", handler.AppList)
	r.POST("/app/change_use", handler.AppChangeUse)
	r.GET("/app/detail", handler.AppDetail)

	r.GET("/app/add", handler.AppAdd)
	r.POST("/app/create", handler.AppCreate)
}
