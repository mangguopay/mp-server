package p

import (
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func LogToHeadquarters() {
	str := util.RandomDigitStr(6)
	resp := handler.DoSend(UrlBase+"/bill/transfer_to_headquarters", map[string]interface{}{
		"image_url":    "1",
		"account_uid":  "ffd703c1-d78e-43ce-864a-c298c2e0f10a",
		"money_type":   "usd",
		"amount":       "1",
		"rec_name":     "23",
		"rec_car_num":  "43",
		"password":     doInitPassword3("111111", str),
		"account_type": "3",
		"non_str":      str,
		"iden_no":      "b153e70c-47ee-4ad2-9e9f-71d60c2f1f44",
		"card_no":      "d6e8eafa-ed2b-45fe-97a9-89ba319b30b4",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
