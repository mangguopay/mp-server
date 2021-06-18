package dao

import (
	"a.a/mp-server/common/cache"
	_ "a.a/mp-server/cust-srv/test"
	"fmt"
	"testing"
)

func TestGetToCustInfoCount(t *testing.T) {
	//count, err := LogToCustDaoInst.GetToCustInfoCountByDate("2020-05-22", "usd")
	//if err != nil {
	//	fmt.Println("err: ", err.Error())
	//	return
	//}
	//fmt.Printf("%+v", count)
	count1, err := OutgoOrderDaoInst.GetToCustInfoCountWriteOffByDate("2020-06-07", "usd")
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}
	fmt.Printf("%+v", count1)
	if err := StatisticUserWithdrawDaoInst.Insert(count1); err != nil {
		fmt.Println("insert err: ", err.Error())
		return
	}
	//if err := StatisticUserWithdrawDaoInst.Insert(count); err != nil {
	//	fmt.Println("insert err: ", err.Error())
	//	return
	//}
	fmt.Println("======================")
}

func TestRedis(t *testing.T) {
	result, err := cache.RedisClient.Keys("err_pwd_*").Result()
	if err != nil {
		fmt.Println("err: ", err.Error())
		return

	}
	fmt.Println(result)
	err = cache.RedisClient.Del("result...").Err()
	if err != nil {
		fmt.Println("Del err: ", err.Error())
		return
	}
}

func TestRediss(t *testing.T) {
	fmt.Println(cache.RedisClient.HVals("user_1").Val())

}
