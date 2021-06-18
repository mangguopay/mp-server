package handler

import (
	"log"
	"net/http"
	"strconv"

	"a.a/mp-server/merchant-mock/conf"

	"a.a/mp-server/merchant-mock/pay"

	"a.a/cu/strext"
	"a.a/mp-server/merchant-mock/dao"

	"github.com/gin-gonic/gin"
)

// 模拟商家扫码用户收款二维码
func MerchantScan(c *gin.Context) {
	c.HTML(http.StatusOK, "merchant/scan.html", gin.H{
		"title": "模拟商家扫码用户收款二维码",
	})
}

// 模拟商家扫码用户收款二维码-下单
func MerchantScanOrder(c *gin.Context) {
	amount := c.PostForm("amount")
	currencyType := c.PostForm("currency_type")
	qrcodeContent := c.PostForm("qrcode_content")

	if amount == "" {
		RedirectTipsError(c, "金额为空")
		return
	}

	if currencyType == "" {
		RedirectTipsError(c, "币种为空")
		return
	}

	if qrcodeContent == "" {
		RedirectTipsError(c, "二维码内容为空")
		return
	}

	log.Printf("amount: %s, currentType: %s, tmpAccountNo: %s\n", amount, currencyType, qrcodeContent)

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
		Title:        "扫码用户",
		CurrencyType: currencyType,
		Amount:       amountInt,
		AppId:        conf.AppId,
	}

	// 插入订单
	if err := dao.OrderInstance.Insert(order); err != nil {
		log.Printf("插入订单失败, err:%v, order:%+v", err, order)
		RedirectTipsError(c, "插入订单失败:"+err.Error())
		return
	}

	// 请求支付系统预下单
	ret, err := pay.ModernPayPreOrder(order, qrcodeContent)
	if err != nil {
		log.Printf("请求支付系统预下单失败, err:%v, order:%+v", err, order)
		RedirectTipsError(c, "请求支付系统预下单失败:"+err.Error())
		return
	}

	// 更新订单状态
	if err := dao.OrderInstance.UpdatePayingStatus(ret.OutOrderNo, ret.OrderNo, ""); err != nil {
		log.Printf("更新订单失败, orderNo:%v, payOrderSn:%v, err:%v", ret.OutOrderNo, ret.OrderNo, err)
		RedirectTipsError(c, "更新订单失败:"+err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/merchant/scan_order_result?order_sn="+order.OrderSn)
}

// 模拟商家扫码用户收款二维码-等待用户支付
func MerchantScanOrderResult(c *gin.Context) {
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

	c.HTML(http.StatusOK, "merchant/scan_result.html", gin.H{
		"title": "等待用户支付",
		"order": order,
	})
}
