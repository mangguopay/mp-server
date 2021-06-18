package p

import (
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func DoSweepWithdraw() {
	str := util.RandomDigitStr(6)
	resp := handler.DoSend(UrlBase+"/bill/sweep_withdraw", map[string]interface{}{
		"amount":       "1",
		"password":     doInitPassword3("111111", str),
		"money_type":   "usd",
		"account_uid":  "0e8d24af-bec7-4f95-b038-c48045f51abf",
		"non_str":      str,
		"account_type": "4",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
