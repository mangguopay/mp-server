package pay

import (
	"fmt"
	"log"
	"testing"

	"a.a/mp-server/common/ss_func"
	"a.a/mp-server/merchant-mock/conf"
	"a.a/mp-server/merchant-mock/myrsa"

	"a.a/cu/encrypt"

	"a.a/cu/strext"
	"a.a/mp-server/merchant-mock/dao"
)

func SetAppInfo() {
	// 广州电商
	conf.AppId = "2020090416495834598604"
	conf.SelfPrivateKey = "MIIEogIBAAKCAQEAwG6NqoQI0rqR999p4rgf6icZzAnNdeBoH0huM+VX+JL7tdaQO9Q4uABDgmulmGBXXHAqzFnRO+HcLMafVk0JDNofkCreOw8VJZ4gkDE7NVvtUbE9kNXMn0DBRC7aOPPbmxTV4A4sm6guoRprV7IuH0SQCjPqsAUjLPnB8yOglTzBGYnKTVT3pxLwrc8eyNB73d3RaY/T2qElPjyrv5lFHgRB7y+3D2r9nm+J0Gai4Q325SEnRBlviL9cV3gfTrZwjmilhWZX7jWUXorDWiPZ8QpXnggyynZsU7PqKycub8VMxAQp1OkO+bUx8wun3QvqCHASrXiHFbUGpz1618qnCwIDAQABAoIBACU9f+cO9FIrzxHkR66bqXl8Ja7p+rwkOKJNNx6N9M9jSpkvM+yQLoXVbzsvL/XkNyPphS7U9vwawqxbp/xgni7Bi7gvw6A0VAhaqLye+nFiH/ReU6bS6W2sb3qNgqfg8Y/6oUViGNnq21NMNJfdILXmY+XFlFaKN/t9Pj6al3op2KFpbC/wZVAMuULMTLaPSfSEVIICslva6EoaSEZwu/ArCj6Sd5aGNnCsEal7ticUCcTyd82huzL44ts81njcgn6JKtlqN1+sbNtI2Qw6vRsqFYq+7vkc9O3oTYqk9zGJXhblAc0qG91bi1DNI5enXlOZzx2AzfE59ko7vI7VjckCgYEA720a4T2qrl8Os3Bsh7gDxPgyEOAQqvW8clAJnimd1pm/ALMuCUqMCzpZnZMoDeadJKnhxtkmYt86dm3EyQYOqLySBZ5acbm3HIA9yLL3Wbx9VXuPCtiyCW/1sqoOCFo5xysNArF+1ACBlRRlOgkvia3JyZNgqa/ugnBPF5EftZ8CgYEAzcCoM7yNx9lsyeUgjmjIqspV2Ku+5Rp6Ut9KrnZAftPmGMW2Big5F26NFdk5BLQniMzvOqFgT8LuAJL5UpNwd8RcqX+j1QzqkNJ/rQWZ6gjY9/HmVvQCT7y97NfHyc+Iqch1JpHYnmm9++6GsOekdzlIJpmuWfnFJ1d8/ti9nxUCgYBcMdMr8KcMxiXPCvd/u2gYLMp6nQ1OB5otGozZjoTs4f8eseCES2Rp3mort0KxN6BDQfdirrONVxAYEmS4U9DJQPOpnjSNnknYe7lO0ztUHrTUeeO29YJ5B9fAmKMWrEebXgSAiQCheiBr25KvKmJXzcfqCwopzUk2iTCrjoJ7WQKBgAWh2njgFnl8CDBWp1d+os+aDlOKEAWxfdu65Q51ijpujoPrFZqBi16K1i3c7mSkkdh606mlNT+5tE4rt8t84b1FkMvLBK8WVW11da1E0/vGFjpjUszajR2lqwVKuttZZQJQzHQr1eQwPNUiqGk7ThM9bC4yUsV/wtfa2z8Wq8opAoGAZ7R3hJwANUor2d6odDE3xGr41koOdBaOFaYlm4A+jtFkP2tQRsAp04V00TklXqDjA7Te0gJEeOJB018Y663Hm9Wn41Zw+8jkSyGM8VHN5a3XQxMidmGMiTaA3/TYHj/Vgde2BK38+K2WQpy6VRhlmsLncFH2HKcM243aUw6elc0="
	conf.PlatformPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5fwtAkc9+7vIUIJEQM48jG97dPTCDjcx1rQeAmiYLk/W0Ygen5bxmF3BqYEWjq3ol+FoesuxoDKDZXdifsC34+C5GJRq7dc5pZgQePKfMm782BmAnbcn3hgEd8fi50tbaPbpwvnHDOP9EPqOHa6boaxNOjFHIvEcQ7EmXFtXu3G/y13/O4gDcuQ4qKHAMPOVU2rlH0X8oGJAKO0x/ul97dtitRTDJzd0l0BLWeTwe8xUNlWOxV9DMlI9J7o1GtwwqJJdskFAykbXjgLEfRW9NbDpO7jqOm2yELpF70RHsFPWB0Gc3oGQ4ATs2XeYIHw6H6lE1sW3q7Vmvarjk7sduwIDAQAB"
}

func TestModernPayPreOrder(t *testing.T) {
	SetAppInfo()

	order := &dao.Order{
		OrderSn:      "merchant" + strext.GetDailyId(),
		Title:        "Huawei mate40",
		CurrencyType: "USD",
		Amount:       700,
		TradeType:    "MANGOPAY_FACE_TO_FACE",
	}
	ret, err := ModernPayPreOrder(order, "")
	if err != nil {
		t.Errorf("ModernPayPreOrder-err: %v", err)
		return
	}

	t.Logf("ret:%+v", ret)
}

func TestQueryPay(t *testing.T) {
	SetAppInfo()

	//payOrderSn := ""
	payOrderSn := "2020091117474451727812"
	//outOrderNo := "merchant2020091117474277248051"
	outOrderNo := ""

	got, err := QueryPay(payOrderSn, outOrderNo)
	if err != nil {
		t.Errorf("QueryPay() error = %v", err)
		return
	}
	t.Logf("返回数据: %v", strext.ToJson(got))
}

func TestModernPayVerifyNotify(t *testing.T) {
	body := []byte(`{"amount":800,"app_id":"2020063014584944251674","currency_type":"USD","order_no":"md2020070814165657213726","out_order_no":"2020070814165656549170","pay_account":"077778888","pay_time":"1594190257","sign":"gvOVEzzEjbUE1bMmLuCb71fd6QL5h+D1EpcPc2Dp5S0DvRR+m7+zt8ZP2/uc2UW0tbSrvVKXEsYKPQn4J+VRFOXJHacdb+anaY7C64aCn47WLzShXfmaK7Ov6UsPLKnViDRNnG3HZcD7bwpQUt9VQ1KzPeL3jwGWEgUuGz3cMKq5xNgaXYe6Crq9I3MTWr+NqTamD75m/AfjkKc5O6TeOHLEo+5XXXT6aNtHxlLRJtafBKNBV/vNm3MulRtEQDtcEXu1HtjBc9YLrUBfU0RQsU9CZHsjEPraUn49GgoRfvn8sNltjUGbdV02ShoNp/Yx+w8tUnqZYu0djvEWKAgkwQ==","sign_type":"RSA2","subject":"购买商品002","timestamp":"1594190257"}`)
	ret, err := NotifyPaymentVerify(body)
	if err != nil {
		t.Errorf("ModernPayVerifyNotify-err: %v", err)
		return
	}

	t.Logf("ret:%+v", ret)
}

func TestEncryptPassword(t *testing.T) {
	originPwd := "111111"
	nonstr := GetRandomString(16)
	databasePwd := "7dcf7a89e67edbcff0695c515aa3d53b"

	enPwd := EncryptPassword(originPwd, nonstr)

	t.Logf("originPwd:%v, nonstr:%v, enPwd:%v", originPwd, nonstr, enPwd)

	if encrypt.DoMd5Salted(databasePwd, nonstr) == enPwd {
		t.Logf("密码验证成功")
	} else {
		t.Logf("密码验证失败")
	}

	//参数: password: 1e459f74e2ea79a547139e58e25a6e19
	//参数: non_str: 1dAs7qLY3PgO5Rbw
	//数据库支付密码字段: 7dcf7a89e67edbcff0695c515aa3d53b
	//判断：encrypt.DoMd5Salted("数据库支付密码字段", "参数non_str") == "参数password"
}

func TestRSA2Sign(t *testing.T) {
	SetAppInfo()

	data := `{"amount":"700","app_id":"2020080510551083814320","attach":"123456","currency_type":"USD","notify_url":"http://www.xxx.com/order/notify","out_order_no":"sn2020082016365543549170","sign_type":"RSA2","subject":"华为mate40","time_expire":"1597913215","timestamp":"1597912615"}`

	reqData := strext.Json2Map(data)
	log.Printf("reqData:%v \n", reqData)

	dataReqStr := ss_func.ParamsMapToString(reqData, SignKey)
	log.Printf("dataReqStr:%v \n", dataReqStr)

	reqSign, err := myrsa.RSA2Sign(dataReqStr, conf.SelfPrivateKey)

	log.Printf("reqSign:%v \n", reqSign)
	log.Printf("err:%v \n", err)
}

func TestRSA2Sign222(t *testing.T) {
	SetAppInfo()

	stringA := "amount=700&app_id=2020090416495834598604&attach=123456&currency_type=USD&lang=zh_CN&notify_url=http://www.xxx.com/order/notify&out_order_no=merchant2020102716543254349170&sign_type=RSA2&subject=Huawei mate40&time_expire=1603789472&timestamp=1603788872&trade_type=MANGOPAY_FACE_TO_FACE"
	reqSign, err := myrsa.RSA2Sign(stringA, conf.SelfPrivateKey)

	log.Printf("reqSign:%v \n", reqSign)
	log.Printf("err:%v \n", err)
}

func TestRSA2VerifySign(t *testing.T) {
	platformPublicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuf5Vr7NpuAqaRYaPJWXWfEN1mCJMz+fgtNwZ6Thg7OzBJ0q3JDGO4Qb8fRISMUk5llCoE2s9Wfp/7IN6qVuCJ0lzgNOj7Sk1XQFQ+knq3mbUE649M8hEgXyJNN6uKtKdOHm+RvKrNgYBzajAqhyiDnn0WyXGA1sUxP4Acu5J2STG92/ovgEF+ZfbXqkZQqAYdjCeSJQsZrjcQCQvlfuR6H6uDQJylpptJFhIZ3HFXYUcG84o27oTHj0BlrsKCxbug5QX2KLWAsvR2VZLnGtUyETJf0gOIVKKlk19l2RAlszk6uQmdEX9fR7nayAaqIiIE0NVn0CjmvdVXacEZrn+MwIDAQAB"

	data := `{"attach":"123456","code":"0","msg":"Success","order_no":"2020082016365544264845","out_order_no":"sn2020082016365543549170","qr_code":"mp://pay/bizpay?qr=2446df9906828e9b39f38f13dc512f18"}`
	respSign := "I7K6QpPhcwrqdEsWfY0W7k9cQraQHA8h2GcPDa6INHulSFhQs/ko1S6RINY1rIt1Vfqe5brLLARNIcLvRGdYzESZTx+AkCowYql+X9eSNjaICLIC5dMI3JkN6dUp9ei30zYMjwoE9ko0Tbo6LNAHtL9naxTQ+rzAfmvdaDjKGdqiTcCGJgYxhAN7i/e6oOFaMc+JVqEHIYG9sLMyOpzBOm5YZc42KFQTFgr7UvNR8N7TSrYTgk3RPr1hCZIGOHk8FUE1uzsz4mD3K66Izueg+XB3r2V2klN1Eo30uZWQzGVRQBNh525Cezqrar/LwjsBTL1wSOesqFmlbp6l0yQQ2Q=="

	respData := strext.Json2Map(data)
	log.Printf("reqData:%v \n", respData)

	dataRespStr := ss_func.ParamsMapToString(respData, SignKey)
	log.Printf("dataRespStr:%v \n", dataRespStr)

	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), platformPublicKey)

	log.Printf("verifyErr:%v \n", verifyErr)
}

func TestRSA2VerifySignBody(t *testing.T) {
	platformPublicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuf5Vr7NpuAqaRYaPJWXWfEN1mCJMz+fgtNwZ6Thg7OzBJ0q3JDGO4Qb8fRISMUk5llCoE2s9Wfp/7IN6qVuCJ0lzgNOj7Sk1XQFQ+knq3mbUE649M8hEgXyJNN6uKtKdOHm+RvKrNgYBzajAqhyiDnn0WyXGA1sUxP4Acu5J2STG92/ovgEF+ZfbXqkZQqAYdjCeSJQsZrjcQCQvlfuR6H6uDQJylpptJFhIZ3HFXYUcG84o27oTHj0BlrsKCxbug5QX2KLWAsvR2VZLnGtUyETJf0gOIVKKlk19l2RAlszk6uQmdEX9fR7nayAaqIiIE0NVn0CjmvdVXacEZrn+MwIDAQAB"

	dataRespStr := "amount=700&app_id=2020080510551083814320&currency_type=USD&order_no=2020082015593742639023&order_status=2&out_order_no=sn2020082015593742149170&pay_account=0855077778888&pay_time=1597935613&sign_type=RSA2&subject=华为mate40&timestamp=1597911862"
	respSign := "bx+3hyqq00SDT3hCZ7vqsd5ObVNKl8t28EpxgoSge1vErqotIdcbikpPljDpqL4X7UNxnP17rHYpppEloJPuXAMOww366DYHMThbeELPZxbwEajecSanTnM9WGng8ygx8MlocLt3pdxXRh+2FW4s6eChMLI7frMU1Hmrh0ESDiqMfGUbvzeDOr27w7gi9lS68PAFW29wpA+2Od1wgiKmloYfLP1uk4ajGL1YKdEMToHGdfa6ArWs9S6a5cYXaMQU433IGegsNjvqhch43yBc+RmZVubGgh+mZaPfUn8ZeWhCWZ5bXw7ssVav7+IRSWDd9pBsP3uskoh4M64C8d5WDA=="

	verifyErr := myrsa.RSA2Verify(dataRespStr, fmt.Sprintf("%s", respSign), platformPublicKey)

	log.Printf("verifyErr:%v \n", verifyErr)
}

func TestModernPayRefundOrder(t *testing.T) {
	SetAppInfo()

	payOrderSn := "2020091117474451727812"
	outOrderNo := ""
	amount := "2000"
	outrefundNo := "refund" + strext.GetDailyId()

	got, err := ModernPayRefundOrder(payOrderSn, outOrderNo, amount, outrefundNo)
	if err != nil {
		t.Errorf("ModernPayRefundQueryOrder() error = %v", err)
		return
	}
	t.Logf("返回数据: %v", strext.ToJson(got))
}

func TestModernPayRefundQueryOrder(t *testing.T) {
	SetAppInfo()

	refundNo := "2020091117571576810935"
	// outRefundNo := "refund2020091117571576249170"
	outRefundNo := ""

	got, err := ModernPayRefundQueryOrder(refundNo, outRefundNo)
	if err != nil {
		t.Errorf("ModernPayRefundQueryOrder() error = %v", err)
		return
	}
	t.Logf("返回数据: %v", strext.ToJson(got))
}

func TestModernPayEnterpriseTransfer(t *testing.T) {
	SetAppInfo()

	request := &dao.Transfer{
		OutTransferNo: "etransfer" + strext.GetDailyId(),
		CurrencyType:  "USD",
		Amount:        700,
		CountryCode:   "0855",
		PayeePhone:    "77778888",
		PayeeEmail:    "",
		Remark:        "奖金",
	}
	ret, err := ModernPayEnterpriseTransfer(request)
	if err != nil {
		t.Errorf("ModernPayEnterpriseTransfer-err: %v", err)
		return
	}

	t.Logf("ret:%+v", ret)
}

func TestModernPayEnterpriseTransferQuery(t *testing.T) {
	SetAppInfo()

	transferNo := "2020103019293708910168"
	//outTransferNo := "transfer2020103015171343849170"
	outTransferNo := ""

	got, err := ModernPayEnterpriseTransferQuery(outTransferNo, transferNo)
	if err != nil {
		t.Errorf("ModernPayEnterpriseTransferQuery() error = %v", err)
		return
	}
	t.Logf("返回数据: %v", strext.ToJson(got))
}
