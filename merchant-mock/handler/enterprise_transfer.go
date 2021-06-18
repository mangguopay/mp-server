package handler

import (
	"log"
	"net/http"
	"strconv"

	"a.a/mp-server/merchant-mock/conf"
	"a.a/mp-server/merchant-mock/dao"
	"a.a/mp-server/merchant-mock/pay"
	"github.com/gin-gonic/gin"
)

// 企业付款-列表
func EnterpriseTransferList(c *gin.Context) {
	page := 1
	pageSize := 100

	list, err := dao.TransferInstance.GetTransferList(page, pageSize)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	c.HTML(http.StatusOK, "enterprise_transfer/list.html", gin.H{
		"title":        "企业付款列表",
		"transferList": list,
		"errMsg":       errMsg,
	})
}

// 企业付款页面显示
func EnterpriseTransferIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "enterprise_transfer/index.html", gin.H{
		"title": "企业付款页面显示",
	})
}

// 进行企业付款
func EnterpriseTransferDo(c *gin.Context) {
	amount := c.PostForm("amount")
	currencyType := c.PostForm("currency_type")
	countryCode := c.PostForm("country_code")
	payeePhone := c.PostForm("payee_phone")
	payeeEmail := c.PostForm("payee_email")
	remark := c.PostForm("remark")

	if amount == "" {
		RedirectTipsError(c, "金额为空")
		return
	}

	if currencyType == "" {
		RedirectTipsError(c, "币种为空")
		return
	}

	if payeePhone == "" && payeeEmail == "" {
		RedirectTipsError(c, "付款人手机号和邮箱不能同时为空")
		return
	}

	log.Printf("amount:%s; currentType:%s, countryCode:%s, payeePhone:%s, payeeEmail:%s, remark:%s \n", amount, currencyType,
		countryCode, payeePhone, payeeEmail, remark)

	amountInt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Printf("解析金额参数失败,amount:%s, err:%v", amount, err)
		RedirectTipsError(c, "解析金额参数失败")
		return
	}

	if currencyType == "USD" {
		amountInt = amountInt * 100
	}

	transfer := &dao.Transfer{
		CurrencyType: currencyType,
		Amount:       amountInt,
		AppId:        conf.AppId,
		CountryCode:  countryCode,
		PayeePhone:   payeePhone,
		PayeeEmail:   payeeEmail,
		Remark:       remark,
	}

	if err := dao.TransferInstance.Insert(transfer); err != nil {
		log.Printf("插入转账订单失败, err:%v, transfer:%+v", err, transfer)
		RedirectTipsError(c, "插入转账订单失败:"+err.Error())
		return
	}

	// 请求支付系统转账
	ret, err := pay.ModernPayEnterpriseTransfer(transfer)
	if err != nil {
		log.Printf("请求支付系统转账失败, err:%v, transfer:%+v", err, transfer)
		RedirectTipsError(c, "请求支付系统转账失败:"+err.Error())
		return
	}

	log.Printf("ret:%+v\n", ret)

	// 更新转账单号
	if err := dao.TransferInstance.UpdateTransferNo(ret.TransferNo, transfer.OutTransferNo); err != nil {
		log.Printf("更新转账订单失败, TransferNo:%v, OutTransferNo:%v, err:%v", ret.TransferNo, transfer.OutTransferNo, err)
		RedirectTipsError(c, "更新转账订单失败:"+err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/enterprise_transfer/list")
}
