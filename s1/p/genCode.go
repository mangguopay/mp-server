package p

import (
	"a.a/cu/strext"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func GenCode() {
	resp := handler.DoSend(UrlBase+"/bill/gen_recv_code", map[string]interface{}{
		"account_uid": "0e8d24af-bec7-4f95-b038-c48045f51abf",
		"money_type":  "usd",
		"amount":      "100",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
