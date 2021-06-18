package dao

import (
	"a.a/cu/ss_time"
	_ "a.a/mp-server/riskctrl-srv/test"
	"fmt"
	"testing"
	"time"
)

func TestUserLoginLogDao_GetLastN(t *testing.T) {
	uid := "49c71695-e29a-4309-91a2-27ebe1547563"

	sd, _ := time.ParseDuration("-24h")
	nTime := time.Now().Add(sd * 7)
	lastTime := ss_time.ForPostgres(nTime)

	fmt.Println("lastTime:", lastTime)

	list, err := UserLoginLogDaoInstance.GetLastByTime(uid, lastTime)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("list:%+v", list)
}
