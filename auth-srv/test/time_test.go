package test

import (
	"a.a/cu/ss_time"
	"fmt"
	"testing"
)

func TestTime(t *testing.T) {

	fmt.Println(ss_time.CheckDateTimeIsRight(ss_time.DateTimeSlashFormat, "2020/07/30 23:59:59"))
}
