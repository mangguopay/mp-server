package dao

import (
	"a.a/cu/strext"
	_ "a.a/mp-server/tm-srv/test"
	"fmt"
	"testing"
)

func Test_TmTxLogDao_Insert(t1 *testing.T) {
	o := TmTxLogDao{}
	o.Id = strext.GetDailyId()
	o.UnFinishNo = 8
	err := o.Insert()

	fmt.Println("err:", err)
}
