package p

import (
	"a.a/cu/strext"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func DeleteBindCard() {
	resp := handler.DoSend(UrlBase+"/info/delete_bind_card", map[string]interface{}{
		"account_uid":  "ffd703c1-d78e-43ce-864a-c298c2e0f10a",
		"money_type":   "usd",
		"car_num":      "1222333",
		"account_type": "3",
		"iden_no":      "b153e70c-47ee-4ad2-9e9f-71d60c2f1f44",
		"channel_name": "中国银行",
		"card_no":      "2e93124d-fc51-45e8-9d33-a61ffa19640d",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
