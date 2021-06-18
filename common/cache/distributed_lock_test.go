package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDistributedLock(t *testing.T) {
	InitRedis("10.41.1.241", "6379", "123456a")

	ok := GetDistributedLock("kkkkkkkkk", "vvvvvvvvv", time.Second*3000)

	fmt.Println("ok:", ok)
}

func TestReleaseDistributedLock(t *testing.T) {
	InitRedis("10.41.1.241", "6379", "123456a")

	ok := ReleaseDistributedLock("kkkkkkkkk", "vvvvvvvvvsss")

	fmt.Println("ok:", ok)
}
