package dao

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"a.a/cu/ss_time"

	"a.a/cu/strext"
)

func TestRefund_Insert(t *testing.T) {
	InitDB()

	refund := &Refund{
		CurrencyType: "USD",
		Amount:       900,
		AppId:        "2020090416495834598604",
		OrderSn:      "merchant2020102717114359267131",
	}

	if err := RefundInstance.Insert(refund); err != nil {
		t.Errorf("插入退款订单失败, err:%v, refund:%+v", err, refund)
		return
	}

	t.Logf("插入退款订单成功,OutRefundNo:%s", refund.OutRefundNo)
}

func TestRefund_GetRefundList(t *testing.T) {
	InitDB()

	page := 1
	pageSize := 5

	list, err := RefundInstance.GetRefundList(page, pageSize)

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("list:%s", strext.ToJson(list))
}

func TestRefund_UpdateRefundNo(t *testing.T) {
	InitDB()

	outRefundNo := "refund2020102917051006349170" // 本地系统退款单号
	refundNo := "2020102917051006349170"          // 支付系统退款单号

	err := RefundInstance.UpdateRefundNo(outRefundNo, refundNo)
	if err != nil {
		t.Errorf("更新退款单号失败, outRefundNo:%v, refundNo:%v, err:%v", outRefundNo, refundNo, err)
		return
	}

	t.Logf("更新退款单号成功,outRefundNo:%s", outRefundNo)
}

func TestRefund_GetOneByOutRefundNo(t *testing.T) {
	InitDB()

	outRefundNo := "refund2020102917051006349170"

	transfer, err := RefundInstance.GetOneByOutRefundNo(outRefundNo)
	if err != nil {
		t.Errorf("查询退款信息失败, outRefundNo:%v, err:%v", outRefundNo, err)
		return
	}

	t.Logf("退款信息:%s", strext.ToJson(transfer))

}

func TestRefund_UpdateRefundSuccess(t *testing.T) {
	InitDB()

	outRefundferNo := "refund2020102917051006349170"
	refundTime := fmt.Sprintf("%d", time.Now().Unix())

	timestamp, perr := strconv.ParseInt(refundTime, 10, 64)
	if perr != nil {
		t.Errorf("解析时间参数失败, refundTime:%v, err:%v", refundTime, perr)
		return
	}

	refundTimeStr := ss_time.ForPostgres(time.Unix(timestamp, 0))

	err := RefundInstance.UpdateRefundSuccess(outRefundferNo, refundTimeStr)
	if err != nil {
		t.Errorf("更新订单退款成功失败, outRefundferNo:%v, refundTime:%v, err:%v", outRefundferNo, refundTime, err)
		return
	}

	t.Logf("更新订单退款成功,outRefundferNo:%s", outRefundferNo)
}

func TestRefund_UpdateRefundFail(t *testing.T) {
	InitDB()

	outRefundferNo := "refund2020102917051006349170"

	err := RefundInstance.UpdateRefundFail(outRefundferNo)
	if err != nil {
		t.Errorf("更新订单退款失败失败, outRefundferNo:%v, err:%v", outRefundferNo, err)
		return
	}

	t.Logf("更新订单退款失败成功,outRefundferNo:%s", outRefundferNo)
}
