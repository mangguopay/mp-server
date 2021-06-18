package p

import (
	"a.a/cu/strext"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func QueryRate() {
	resp := handler.DoSend(UrlBase+"/bill/query_rate", map[string]interface{}{
		"amount":       "1000.03",
		"type":         2,
		"account_type": "4",
		"account_uid":  "ffd703c1-d78e-43ce-864a-c298c2e0f10a",
		"iden_no":      "098f0076-f097-4eac-833e-6097579f7783",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
