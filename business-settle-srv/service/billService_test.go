package service

import (
	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
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
		name := "p_crm"
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
	m := map[string]interface{}{}
	cache.RedisCli.OpenWithPassword("10.41.1.241",
		"6379",
		"123456a",
		"a")
	fmt.Printf("not support cache|driver=[%v]\n", m["adapter"])
}

func init() {
	DoInitDb()
	DoInitRedis()
}
func TestBillService(t *testing.T) {
	req := &go_micro_srv_settle.SettleFeesRequest{
		OrderNo:  "50012020030620175656734633",
		FeesType: "1",
		Fees:     "100",
		Amount:   "1000",
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancel()
	if errStr := BillService(ctx, req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "ExchangeService 单元测试失败")
	}
	t.Logf("----------->%s", "ExchangeService 单元测试成功")
}

func TestWithdrawalService(t *testing.T) {
	req := &go_micro_srv_settle.SettleFeesRequest{
		OrderNo:  "50012020030619013432468109",
		FeesType: "2",
		Fees:     "100",
		Amount:   "1000",
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancel()
	if errStr := WithdrawalService(ctx, req); errStr != ss_err.ERR_SUCCESS {
		t.Fatalf("----------->%s", "WithTimeout 单元测试失败")
	}
	t.Logf("----------->%s", "ExchangeService 单元测试成功")
}
