package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"a.a/mp-server/merchant-mock/conf"

	"a.a/cu/strext"

	"a.a/mp-server/merchant-mock/dao"
	"a.a/mp-server/merchant-mock/pay"

	"github.com/gin-gonic/gin"
)

type QrCode struct {
	OrderNo    string
	OutOrderNo string
	Amount     int64
}

// 订单列表
func OrderList(c *gin.Context) {
	page := 1
	pageSize := 100

	list, err := dao.OrderInstance.GetOrderList(page, pageSize)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	c.HTML(http.StatusOK, "order/list.html", gin.H{
		"title":     "订单列表",
		"orderList": list,
		"errMsg":    errMsg,
	})
}

// 创建订单-显示
func OrderAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "order/add.html", gin.H{
		"title": "添加订单",
	})
}

// 创建订单
func OrderCreate(c *gin.Context) {
	tradeType := c.PostForm("trade_type")
	title := c.PostForm("title")
	currencyType := c.PostForm("currency_type")
	amount := c.PostForm("amount")

	if tradeType == "" {
		RedirectTipsError(c, "交易类型为空")
		return
	}

	if title == "" {
		RedirectTipsError(c, "订单标题为空")
		return
	}

	if currencyType == "" {
		RedirectTipsError(c, "币种为空")
		return
	}

	if amount == "" {
		RedirectTipsError(c, "金额为空")
		return
	}

	log.Printf("title: %s; currentType: %s, amount: %s \n", title, currencyType, amount)

	amountInt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Printf("解析金额参数失败,amount:%s, err:%v", amount, err)
		RedirectTipsError(c, "解析金额参数失败")
		return
	}

	if currencyType == "USD" {
		amountInt = amountInt * 100
	}

	order := &dao.Order{
		OrderSn:      strext.GetDailyId(),
		Title:        title,
		CurrencyType: currencyType,
		Amount:       amountInt,
		AppId:        conf.AppId,
		TradeType:    tradeType,
	}

	// 插入订单
	if err := dao.OrderInstance.Insert(order); err != nil {
		log.Printf("插入订单失败, err:%v, order:%+v", err, order)
		RedirectTipsError(c, "插入订单失败:"+err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/order/list")
}

// 支付订单-显示界面
func OrderPay(c *gin.Context) {
	orderSn := c.Query("order_sn")
	if orderSn == "" {
		RedirectTipsError(c, "订单号为空")
		return
	}

	order, err := dao.OrderInstance.GetOneByOrderSn(orderSn)
	if err != nil {
		log.Printf("查询订单失败, orderSn:%v, err:%v", orderSn, err)
		RedirectTipsError(c, "查询订单失败:"+err.Error())
		return
	}

	if order.Status != dao.OrderStatusUnpay && order.Status != dao.OrderStatusPaying {
		RedirectTipsError(c, "订单状态不是“未支付”或“支付中”状态,不能进行支付操作")
		return
	}

	if order.PayOrderSn == "" { // 没有支付系统订单号，才需要请求支付系统下单
		// 请求支付系统预下单
		ret, err := pay.ModernPayPreOrder(order, "")
		if err != nil {
			log.Printf("请求支付系统预下单失败, err:%v, order:%+v", err, order)
			RedirectTipsError(c, "请求支付系统预下单失败:"+err.Error())
			return
		}
		// 更新订单状态
		if err := dao.OrderInstance.UpdatePayingStatus(ret.OutOrderNo, ret.OrderNo, ret.QrCode); err != nil {
			log.Printf("更新订单失败, orderNo:%v, payOrderSn:%v, err:%v", ret.OutOrderNo, ret.OrderNo, err)
			RedirectTipsError(c, "更新订单失败:"+err.Error())
			return
		}

		order.QrCode = ret.QrCode
		order.PayOrderSn = ret.OrderNo
	}

	c.HTML(http.StatusOK, "order/pay.html", gin.H{
		"title":      "支付订单",
		"order":      order,
		"qrCodeText": order.QrCode, //pay.CreatePayData(order),
	})
}

func OrderPayQuery(c *gin.Context) {
	orderSn := c.Query("order_sn")
	isPaid := 0 // 是否已支付

	if orderSn == "" {
		c.JSON(http.StatusOK, gin.H{"code": 100, "msg": "缺少order_sn参数", "is_paid": isPaid})
		return
	}

	// 查询订单信息
	order, err := dao.OrderInstance.GetOneByOrderSn(orderSn)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 101, "msg": "查询订单失败:" + err.Error(), "is_paid": isPaid})
		return
	}

	if order.Status == dao.OrderStatusPaid {
		isPaid = 1
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "is_paid": isPaid})
	return
}

// 跳转到错误提示
func RedirectTipsError(c *gin.Context, errMsg string) {
	redirectUri := "/order/tips?err_msg=" + errMsg
	c.Redirect(http.StatusFound, redirectUri)
}

// 跳转到成功提示
func RedirectTipsSuccess(c *gin.Context, errMsg string) {
	redirectUri := "/order/tips?success_msg=" + errMsg
	c.Redirect(http.StatusFound, redirectUri)
}

// 提示信息页面
func TipsPage(c *gin.Context) {
	errMsg := c.Query("err_msg")
	successMsg := c.Query("success_msg")

	c.HTML(http.StatusOK, "order/tips.html", gin.H{
		"title":      "提示信息",
		"errMsg":     errMsg,
		"successMsg": successMsg,
	})
}

// 订单退款申请
func OrderRefund(c *gin.Context) {
	orderSn := c.Query("order_sn")
	if orderSn == "" {
		RedirectTipsError(c, "订单号为空")
		return
	}

	order, err := dao.OrderInstance.GetOneByOrderSn(orderSn)
	if err != nil {
		log.Printf("查询订单失败, orderSn:%v, err:%v", orderSn, err)
		RedirectTipsError(c, "查询订单失败:"+err.Error())
		return
	}

	if order.Status != dao.OrderStatusPaid {
		RedirectTipsError(c, "订单状态不是“已支付”状态,不能进行退款操作")
		return
	}

	refund := &dao.Refund{
		CurrencyType: order.CurrencyType,
		Amount:       order.Amount,
		AppId:        order.AppId,
		OrderSn:      order.OrderSn,
	}

	if err := dao.RefundInstance.Insert(refund); err != nil {
		log.Printf("插入退款订单失败, err:%v, refund:%+v", err, refund)
		RedirectTipsError(c, "插入退款订单失败:"+err.Error())
		return
	}

	// 请求支付系统申请退款
	ret, err := pay.ModernPayRefundOrder(order.PayOrderSn, order.OrderSn, fmt.Sprintf("%v", order.Amount), refund.OutRefundNo)
	if err != nil {
		log.Printf("请求支付系统申请退款失败, err:%v, order:%+v", err, order)
		RedirectTipsError(c, "请求支付系统申请退款失败:"+err.Error())
		return
	}

	// 更新退款单号
	if err := dao.RefundInstance.UpdateRefundNo(refund.OutRefundNo, ret.RefundNo); err != nil {
		log.Printf("更新退款单号失败, OutRefundNo:%v, RefundNo:%v, err:%v", refund.OutRefundNo, ret.RefundNo, err)
		RedirectTipsError(c, "更新退款单号失败:"+err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/order/refund/list")
}

// 订单退款列表
func OrderRefundList(c *gin.Context) {
	page := 1
	pageSize := 100

	list, err := dao.RefundInstance.GetRefundList(page, pageSize)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	c.HTML(http.StatusOK, "order/refund_list.html", gin.H{
		"title":      "退款列表",
		"refundList": list,
		"errMsg":     errMsg,
	})
}
