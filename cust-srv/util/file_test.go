package util

import (
	custProto "a.a/mp-server/common/proto/cust"
	_ "a.a/mp-server/cust-srv/test"
	"testing"
)

func TestCreateBillFile(t *testing.T) {
	filePath := "D:\\a.xlsx"
	titles := []string{
		"支付订单号",
		"订单号",
		"创建时间",
		"币种",
		"金额",
		"手续费",
		"应用名称",
		"商品名称",
		"订单状态",
	}
	orderList := []*custProto.BusinessBillData{
		{
			OrderNo:      "2020082719413862920291",
			Fee:          "350",
			CreateTime:   "2020-08-27 18:41:40",
			Amount:       "10000",
			OrderStatus:  "4",
			OutOrderNo:   "20200827194140773",
			CurrencyType: "USD",
			Subject:      "好吧",
			AppName:      "测速机120",
		},
	}
	queryStr := []string{
		"#起始日期：[2020年08月21日 00:00:00]",
		"终止日期：[2020年08月22日 00:00:00]",
	}
	f := &FileCommonInfo{
		FileName:     "a.xlsx",
		FilePath:     filePath,
		Head:         "#modernpay待结算/已结算订单明细",
		Account:      "#账号：lj15939519169@163.com",
		Total:        "1",
		USDCnt:       "1",
		USDSum:       "10000",
		USDRefundCnt: "0",
		USDRefundSum: "0",
		KHRCnt:       "0",
		KHRSum:       "0",
		KHRRefundCnt: "0",
		KHRRefundSum: "0",
		QueryStr:     queryStr,
		Titles:       titles,
		OrderList:    orderList,
	}
	if err := CreateBillFile(f); err != nil {
		t.Errorf("CreateBillFile() error = %v", err)
		return
	}
}
