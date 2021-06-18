package dao

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"a.a/cu/ss_time"

	"a.a/cu/strext"
)

func TestTransfer_Insert(t *testing.T) {
	InitDB()

	transfer := &Transfer{
		CurrencyType: "USD",
		Amount:       900,
		AppId:        "2020090416495834598604",
		CountryCode:  "855",
		PayeePhone:   "77778888",
		PayeeEmail:   "",
		Remark:       "转账备注",
	}

	if err := TransferInstance.Insert(transfer); err != nil {
		t.Errorf("插入转账订单失败, err:%v, transfer:%+v", err, transfer)
		return
	}

	t.Logf("插入转账订单成功,OutTransferNo:%s", transfer.OutTransferNo)
}

func TestTransfer_GetTransferList(t *testing.T) {
	InitDB()

	page := 1
	pageSize := 5

	list, err := TransferInstance.GetTransferList(page, pageSize)

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("list:%s", strext.ToJson(list))
}

func TestTransfer_UpdateTransferNo(t *testing.T) {
	InitDB()

	transferNo := "2020102819062338267131"            // 支付系统转账单号
	outTransferNo := "transfer2020102819062338267131" // 内部系统转账单号

	err := TransferInstance.UpdateTransferNo(transferNo, outTransferNo)
	if err != nil {
		t.Errorf("更新订单失败, transferNo:%v, outTransferNo:%v, err:%v", transferNo, outTransferNo, err)
		return
	}

	t.Logf("更新订单,transferNo:%s", transferNo)
}

func TestTransfer_GetOneByOutTransferNo(t *testing.T) {
	InitDB()

	outTransferNo := "transfer2020102911423487049170"

	transfer, err := TransferInstance.GetOneByOutTransferNo(outTransferNo)
	if err != nil {
		t.Errorf("查询转账信息失败, outTransferNo:%v, err:%v", outTransferNo, err)
		return
	}

	t.Logf("转账信息:%s", strext.ToJson(transfer))
}

func TestTransfer_UpdateTransferSuccess(t *testing.T) {
	InitDB()

	outTransferNo := "transfer2020102911423487049170"
	transferTime := fmt.Sprintf("%d", time.Now().Unix())

	timestamp, perr := strconv.ParseInt(transferTime, 10, 64)
	if perr != nil {
		t.Errorf("解析时间参数失败, transferTime:%v, err:%v", transferTime, perr)
		return
	}

	transferTimeStr := ss_time.ForPostgres(time.Unix(timestamp, 0))

	err := TransferInstance.UpdateTransferSuccess(outTransferNo, transferTimeStr)
	if err != nil {
		t.Errorf("更新转账成功失败, outTransferNo:%v, transferTimeStr:%v, err:%v", outTransferNo, transferTimeStr, err)
		return
	}

	t.Logf("更新转账成功,outTransferNo:%s", outTransferNo)
}

func TestTransfer_UpdateTransferFail(t *testing.T) {
	InitDB()

	outTransferNo := "transfer2020102910253420749170"

	err := TransferInstance.UpdateTransferFail(outTransferNo)
	if err != nil {
		t.Errorf("更新转账失败失败, outTransferNo:%v, err:%v", outTransferNo, err)
		return
	}

	t.Logf("更新转账失败成功,outTransferNo:%s", outTransferNo)
}
