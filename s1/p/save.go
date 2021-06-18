package p

import (
	"a.a/cu/strext"
	"a.a/cu/util"
	"a.a/mp-server/s1/handler"
	"fmt"
)

func QuerySaveReceipt() {
	resp := handler.DoSend(UrlBase+"/bill/query_save_receipt", map[string]interface{}{
		"order_no": "2019121612081171841501",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
func QuerySave() {
	resp := handler.DoSend(UrlBase+"/info/my_data", map[string]interface{}{
		"op_acc_no":    "b153e70c-47ee-4ad2-9e9f-71d60c2f1f44",
		"account_type": "3",
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
func Save(recvPhone, sendPhone, amount, moneyType string) {
	str := util.RandomDigitStr(6)
	resp := handler.DoSend(UrlBase+"/bill/save_money", map[string]interface{}{
		"recv_phone": recvPhone,
		"send_phone": sendPhone,
		"amount":     amount,
		"password":   doInitPassword3("1", str),
		"money_type": moneyType,
		"non_str":    str,
	}, strext.ToStringNoPoint(M["token"]))
	fmt.Println(resp)
}
