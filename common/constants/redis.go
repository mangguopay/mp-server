package constants

import "fmt"

const (
	DefPoolName = "a"

	QUEUENAME_PAY_ORDER_NO     = "order_id_list"
	QUEUENAME_PRODUCT_ORDER_NO = "prod_order_id_list"

	// 服务使用到的redis库
	BusinessBillSrvRedisDb = 5 // business-bill-srv专用redis库
	NotifySrvRedisDb       = 6 // notify-srv专用redis库

	EXP_GEN_CODE_KEY = "listenExp-"
)

const (
	//用户付款码key
	CustPaymentCode = "payment_code"
)

func GetCustPaymentCodeRedisKey(paymentCode string) string {
	return fmt.Sprintf("%s_%s", CustPaymentCode, paymentCode)
}
