package p

import (
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func AddCard() {
	str := util.RandomDigitStr(6)
	resp := handler.DoSend(UrlBase+"/info/add_card", map[string]interface{}{
		"account_uid":  "ffd703c1-d78e-43ce-864a-c298c2e0f10a",
		"money_type":   "usd",
		"rec_car_num":  "222333",
		"rec_name":     "张三",
		"password":     doInitPassword3("111111", str),
		"account_type": "3",
		"iden_no":      "b153e70c-47ee-4ad2-9e9f-71d60c2f1f44",
		"channel_name": "中国银行",
		"is_default":   1,
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
