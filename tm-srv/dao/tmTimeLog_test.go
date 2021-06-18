package dao

import (
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/global"
	_ "a.a/mp-server/tm-srv/test"
	"fmt"
	"testing"
	"time"
)

func Test_TmTimeLogDao_Insert(t1 *testing.T) {
	start := ss_time.Now(global.Tz)

	time.Sleep(time.Second * 1)
	end := ss_time.Now(global.Tz)

	o := TmTimeLogDao{}
	o.TxNo = strext.GetDailyId()
	o.StartTime = ss_time.ForPostgres(start)
	o.EndTime = ss_time.ForPostgres(end)
	o.Duration = int64(end.Sub(start) / 1000 / 1000) // 单位: 毫秒
	o.EndMode = "commit"
	o.SqlList = ""
	o.SqlNum = 8

	err := o.Insert()

	fmt.Println("err:", err)
}
