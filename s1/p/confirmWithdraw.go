package p

import (
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func ConfirmWithdraw() {
	str := util.RandomDigitStr(6)
	resp := handler.DoSend(UrlBase+"/bill/confirm_withdraw", map[string]interface{}{
		"amount":          "1",
		"password":        doInitPassword3("111111", str),
		"money_type":      "usd",
		"use_account_uid": "0e8d24af-bec7-4f95-b038-c48045f51abf",
		"non_str":         str,
		"account_type":    "5",
		"account_uid":     "0313b149-0eab-45f7-be56-11ed05f1b257",
		"out_order_no":    "2019121616122923857086",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
