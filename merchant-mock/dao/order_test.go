package dao

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"a.a/cu/ss_time"

	"a.a/cu/strext"
)

func TestOrder_GetOrderList(t *testing.T) {
	InitDB()

	page := 2
	pageSize := 5

	list, err := OrderInstance.GetOrderList(page, pageSize)

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("list:%s", strext.ToJson(list))
}

func TestOrder_Insert(t *testing.T) {
	InitDB()

	order := &Order{
		OrderSn:      strext.GetDailyId(),
		Title:        "测试商品001",
		CurrencyType: "USD",
		Amount:       500,
		AppId:        "2020081016292373211673",
	}

	if err := OrderInstance.Insert(order); err != nil {
		t.Errorf("插入订单失败, err:%v, order:%+v", err, order)
		return
	}

	t.Logf("插入订单成功,orderSn:%s", order.OrderSn)
}

func TestOrder_UpdatePayingStatus(t *testing.T) {
	InitDB()

	orderNo := "2020070719484342249170"      // 内部系统id号
	payOrderSn := "md2020070719525187253061" // 支付系统订单号
	qrCode := "xxxxxx"

	err := OrderInstance.UpdatePayingStatus(orderNo, payOrderSn, qrCode)
	if err != nil {
		t.Errorf("更新订单失败, orderNo:%v, payOrderSn:%v, err:%v", orderNo, payOrderSn, err)
		return
	}

	t.Logf("更新订单,orderSn:%s", orderNo)
}

func TestOrder_GetOneByOrderSn(t *testing.T) {
	InitDB()

	paramsOrderSn := "merchant2020103011093683974792"

	order, err := OrderInstance.GetOneByOrderSn(paramsOrderSn)
	if err != nil {
		t.Errorf("查询订单失败, paramsOrderSn:%v, err:%v", paramsOrderSn, err)
		return
	}

	t.Logf("订单信息:%s", strext.ToJson(order))
}

func TestOrder_UpdateStatus(t *testing.T) {
	InitDB()

	orderNo := "2020070720000292274792" // 内部系统id号
	status := OrderStatusPaid           // 支付系统订单号

	err := OrderInstance.UpdateStatus(orderNo, status)
	if err != nil {
		t.Errorf("更新订单状态失败, orderNo:%v, status:%v, err:%v", orderNo, status, err)
		return
	}

	t.Logf("更新订单状态成功,orderSn:%s, status:%d", orderNo, status)
}

func TestOrder_UpdatePaidOk(t *testing.T) {
	InitDB()

	orderNo := "2020070814165656549170"
	payTime := fmt.Sprintf("%d", time.Now().Unix())
	payAccount := "077778888"

	timestamp, perr := strconv.ParseInt(payTime, 10, 64)
	if perr != nil {
		t.Errorf("解析时间参数失败, payTime:%v, err:%v", payTime, perr)
		return
	}

	payTimeStr := ss_time.ForPostgres(time.Unix(timestamp, 0))

	err := OrderInstance.UpdatePaidOk(orderNo, payTimeStr, payAccount)
	if err != nil {
		t.Errorf("更新订单支付成功失败, orderNo:%v, payTime:%v, payAccount:%v, err:%v", orderNo, payTime, payAccount, err)
		return
	}

	t.Logf("更新订单支付成功,orderSn:%s", orderNo)

}
