package service

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/proto/settle"
	"a.a/mp-server/common/ss_err"
	"context"
	"fmt"
	"testing"
	"time"
)

func DoInitDb() {
	l := []string{constants.DB_CRM}
	for _, v := range l {
		host := "10.41.1.241"
		port := "5432"
		user := "postgres"
		password := "123"
		name := "mp_crm"
		alias := strext.ToStringNoPoint(v)
		driver := "postgres"
		switch driver {
		case "postgres":
			db.DoDBInitPostgres(alias, host, port, user, password, name)
		default:
			fmt.Printf("not support database|driver=[%v]\n", driver)
		}
	}
}

func DoInitRedis() {
	err := cache.InitRedis("10.41.1.241", "6379", "123456a")
	fmt.Println("InitRedis-err:", err)
}

func init() {
	DoInitDb()
	DoInitRedis()
}
func TestExchangeService(t *testing.T) {
	req := &go_micro_srv_settle.SettleTransferRequest{
		BillNo:    "2020010216242426143611",
		MoneyType: "usd",
		Fees:      "1000",
		FeesType:  constants.FEES_TYPE_EXCHANGE,
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancel()
	if errStr := ExchangeService(ctx, req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "ExchangeService 单元测试失败")
	}
	t.Logf("----------->%s", "ExchangeService 单元测试成功")
}

func TestTransferService(t *testing.T) {
	req := &go_micro_srv_settle.SettleTransferRequest{
		BillNo:    "2019123114222280589994",
		FeesType:  constants.FEES_TYPE_TRANSFER,
		Fees:      "1000",
		MoneyType: "usd",
	}
	if errStr := TransferService(req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "TransferService 单元测试失败")
	}
	t.Logf("----------->%s", "TransferService 单元测试成功")
}

func TestSavemoneyService(t *testing.T) {
	req := &go_micro_srv_settle.SettleTransferRequest{
		BillNo:    "2019122811522480376489",
		FeesType:  constants.FEES_TYPE_SAVEMONEY,
		Fees:      "2200",
		MoneyType: "khr",
	}

	if errStr := SavemoneyService(req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "SavemoneyService 单元测试失败")
	}
	t.Logf("----------->%s", "SavemoneyService 单元测试成功")
}

func TestMobileNumWithdrawService(t *testing.T) {
	req := &go_micro_srv_settle.SettleTransferRequest{
		BillNo:    "2020010317195130778314",
		Fees:      "1000",
		FeesType:  constants.FEES_TYPE_MOBILE_NUM_WITHDRAW,
		MoneyType: "usd",
	}
	if errStr := MobileNumWithdrawService(req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "MobileNumWithdrawService 单元测试失败")
	}
	t.Logf("----------->%s", "MobileNumWithdrawService 单元测试成功")
}

func TestSweepWithdrawService(t *testing.T) {
	req := &go_micro_srv_settle.SettleTransferRequest{
		BillNo:    "2019123117464842966529",
		FeesType:  constants.FEES_TYPE_SWEEP_WITHDRAW,
		Fees:      "1000",
		MoneyType: "khr",
	}
	if errStr := SweepWithdrawService(req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "SweepWithdrawService 单元测试失败")
	}
	t.Logf("----------->%s", "SweepWithdrawService 单元测试成功")
}

func TestCollectionService(t *testing.T) {
	req := &go_micro_srv_settle.SettleTransferRequest{
		BillNo:    "2019122716553759614581",
		FeesType:  constants.FEES_TYPE_COLLECTION,
		Fees:      "100",
		MoneyType: "usd",
	}
	if errStr := CollectionService(req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "CollectionService 单元测试失败")
	}
	t.Logf("----------->%s", "CollectionService 单元测试成功")
}
