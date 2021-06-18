package dao

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"testing"
)

func TestBusinessTransferOrderDao_Insert(t *testing.T) {
	d := &BusinessTransferOrderDao{}
	d.FromAccountNo = "951f1dcf-31d1-4fb4-a992-e204a62eca0a"
	d.FromBusinessNo = "b7733819-70c6-479b-bdba-25d369255481"
	d.ToBusinessNo = "0e05f94b-dad4-4293-8f10-70dbae1830bc"
	d.ToAccountNo = "951f1dcf-31d1-4fb4-a992-e204a62eca0a"
	d.Amount = "100"
	d.RealAmount = "100" //外扣
	d.Fee = "0"
	d.Rate = "0"
	d.CurrencyType = "USD"
	d.OrderStatus = "0"
	d.WrongReason = ""
	d.PaymentType = "2"
	d.BatchNo = ""
	d.BatchRowNum = "0"
	d.TransferType = constants.BusinessTransferOrderTypeOrdinary
	got, err := BusinessTransferOrderDaoInst.Insert(d)
	if err != nil {
		t.Errorf("GetBusinessStatusInfo() error = %v", err)
		return
	}
	t.Logf(strext.ToJson(got))
}
func TestBusinessTransferOrderDao_Add(t *testing.T) {
	d := &BusinessTransferOrderDao{}
	d.FromAccountNo = "951f1dcf-31d1-4fb4-a992-e204a62eca0a"
	d.FromBusinessNo = "b7733819-70c6-479b-bdba-25d369255481"
	d.ToBusinessNo = ""
	d.ToAccountNo = "951f1dcf-31d1-4fb4-a992-e204a62eca0a"
	d.Amount = "100"
	d.RealAmount = "100" //外扣
	d.Fee = "0"
	d.Rate = "0"
	d.CurrencyType = "USD"
	d.OrderStatus = "0"
	d.WrongReason = ""
	d.PaymentType = "2"
	d.TransferType = constants.BusinessTransferOrderTypeOrdinary
	got, err := BusinessTransferOrderDaoInst.Add(d)
	if err != nil {
		t.Errorf("GetBusinessStatusInfo() error = %v", err)
		return
	}
	t.Logf(strext.ToJson(got))
}

func TestBusinessTransferOrderDao_Insert2(t *testing.T) {
	d := &BusinessTransferOrderDao{}
	d.FromAccountNo = "951f1dcf-31d1-4fb4-a992-e204a62eca0a"
	d.FromBusinessNo = "b7733819-70c6-479b-bdba-25d369255481"
	d.ToBusinessNo = "0e05f94b-dad4-4293-8f10-70dbae1830bc"
	d.ToAccountNo = "951f1dcf-31d1-4fb4-a992-e204a62eca0a"
	d.Amount = "100"
	d.RealAmount = "100" //外扣
	d.Fee = "0"
	d.Rate = "0"
	d.CurrencyType = "USD"
	d.OrderStatus = "0"
	d.WrongReason = ""
	d.PaymentType = "2"
	d.BatchNo = ""
	d.BatchRowNum = "0"
	d.TransferType = constants.BusinessTransferOrderTypeOrdinary
	gotLogNoT, err := BusinessTransferOrderDaoInst.Insert2(d)
	if err != nil {
		t.Errorf("Insert2() error = %v", err)
		return
	}
	t.Logf(strext.ToJson(gotLogNoT))
}
