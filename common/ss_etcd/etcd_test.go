package ss_etcd

import (
	"a.a/mp-server/common/constants"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetEtcdAddrFromEnv(t *testing.T) {
	err := os.Setenv(constants.ETCDAddrEnvName, "aaa, aaa, bbb, ")

	fmt.Println("err:", err)

	list := GetEtcdAddrFromEnv()

	a := strings.Split(``, ",")

	fmt.Println("a:", len(a))
	fmt.Println("len:", len(list))
	fmt.Println("list:", list)
}
